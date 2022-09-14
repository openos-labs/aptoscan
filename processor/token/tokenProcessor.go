package token

import (
	"apotscan/logger"
	"apotscan/types"
	"apotscan/types/token"
	"fmt"
	mapset "github.com/deckarep/golang-set"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"strconv"
)

type TokenTransactionProcessor struct {
	db            *gorm.DB
	redisCli      *redis.Client
	chainId       uint8
	name          string
	logger        *logger.Logger
	indexTokenUri bool
}

func New(name string, redisCli *redis.Client, db *gorm.DB, chainId uint8, logConf *logger.Config, indexTokenUri bool) (*TokenTransactionProcessor, error) {
	_logger, err := logger.New(logConf)
	if err != nil {
		return nil, err
	}
	return &TokenTransactionProcessor{
		db:            db,
		redisCli:      redisCli,
		chainId:       chainId,
		name:          name,
		logger:        _logger,
		indexTokenUri: indexTokenUri,
	}, nil
}

func (tp *TokenTransactionProcessor) Name() string {
	return tp.name
}

func (tp *TokenTransactionProcessor) ChainId() uint8 {
	return tp.chainId
}

func (tp *TokenTransactionProcessor) GetDB() *gorm.DB {
	return tp.db
}

func (tp *TokenTransactionProcessor) GetRedis() *redis.Client {
	return tp.redisCli
}

func (tp *TokenTransactionProcessor) GetLogger() *logger.Logger {
	return tp.logger
}

func (tp *TokenTransactionProcessor) ProcessTransactions(txs []types.Transaction, startVersion, endVersion int64) (*types.ProcessResult, error) {
	var tokenUris = make(map[string]string)
	txsWithTokenEvent, err := token.GetTransactionsWithTokenEvent(txs)
	if err != nil {
		return nil, err
	}
	if err = processTokenOnChainData(tp.db, txsWithTokenEvent, &tokenUris); err != nil {
		return nil, err
	}
	if tp.indexTokenUri {
		//todo: deal with metadata
	}
	return &types.ProcessResult{
		Name:         tp.Name(),
		StartVersion: startVersion,
		EndVersion:   endVersion,
	}, nil
}

//processTokenOnChainData todo: add logic
func processTokenOnChainData(db *gorm.DB, txsWithEvents []*token.TransactionWithTokenEvents, uris *map[string]string) error {
	var collections []*token.CollectionInDB
	var tokensData []*token.TokenDataInDB
	var tokenTransferEvents []*token.TokenTransferEventInDB
	var tokenActivities []*token.TokenActivityInDB
	var ownershipChanges []*token.OwnershipInDB
	var tokenDataChange map[string]int64

	ownershipStr := fmt.Sprintf("select * from %s where (token_id,owner) in (", token.OwnershipInDB{}.TableName())
	ownershipSet := mapset.NewSet()

	for _, tx := range txsWithEvents {
		for _, event := range tx.TokenEvents {
			switch event.TokenEventData.EventType() {
			case token.TypeWithdrawEvent:
				//todo: 增加withdraw和deposit的event
				tokenId := event.TokenEventData.(token.WithdrawEvent).Id.ToString()

				if !ownershipSet.Contains(tokenId + tx.Tx.Sender) {
					ownershipSet.Add(tokenId + tx.Tx.Sender)
					ownershipStr = ownershipStr + fmt.Sprintf("(%s,%s),", tokenId, tx.Tx.Sender)
				}

				ownershipChanges = append(ownershipChanges, &token.OwnershipInDB{
					TokenId: tokenId,
					Owner:   tx.Tx.Sender,
					Amount:  -event.TokenEventData.(token.WithdrawEvent).Amount,
					Version: tx.Tx.Version,
				})

			case token.TypeDepositEvent:
				tokenId := event.TokenEventData.(token.WithdrawEvent).Id.ToString()

				if !ownershipSet.Contains(tokenId + tx.Tx.Sender) {
					ownershipSet.Add(tokenId + tx.Tx.Sender)
					ownershipStr = ownershipStr + fmt.Sprintf("(%s,%s),", tokenId, tx.Tx.Sender)
				}

				ownershipChanges = append(ownershipChanges, &token.OwnershipInDB{
					TokenId: tokenId,
					Owner:   tx.Tx.Sender,
					Amount:  event.TokenEventData.(token.WithdrawEvent).Amount,
					Version: tx.Tx.Version,
				})

			case token.TypeCreateTokenDataEvent:
				tokenDataId := event.TokenEventData.(token.CreateTokenDataEvent).Id.ToString()
				(*uris)[tokenDataId] = event.TokenEventData.(token.CreateTokenDataEvent).Uri
				tokenInDb, err := getTokenData(event.TokenEventData.(token.CreateTokenDataEvent), &tx.Tx)
				if err != nil {
					return err
				}
				tokensData = append(tokensData, tokenInDb)

			case token.TypeCollectionCreationEvent:
				collectionInDb, err := getCollection(event.TokenEventData.(token.CollectionCreationEvent), &tx.Tx)
				if err != nil {
					return err
				}
				collections = append(collections, collectionInDb)

			case token.TypeBurnTokenEvent:
				sequenceNum, err := strconv.ParseInt(event.SequenceNumber, 10, 64)
				if err != nil {
					return err
				}
				timestamp, err := strconv.ParseInt(tx.Tx.Timestamp, 10, 64)
				if err != nil {
					return nil
				}
				tokenId := event.TokenEventData.(token.BurnTokenEvent).Id.ToString()
				amount := int64(event.TokenEventData.(token.BurnTokenEvent).Amount)
				tokenActivities = append(tokenActivities, &token.TokenActivityInDB{
					EventKey:       event.Key,
					SequenceNumber: sequenceNum,
					Account:        tx.Tx.Sender,
					TokenId:        tokenId,
					EventType:      event.Type,
					Amount:         amount,
					Version:        tx.Tx.Version,
					Timestamp:      timestamp,
				})
				if tokenData, ok := tokenDataChange[tokenId]; !ok {
					tokenDataChange[tokenId] -= amount
				} else {
					tokenData -= amount
					tokenDataChange[tokenId] = tokenData
				}

			case token.TypeMutateTokenPropertyMapEvent:

			case token.TypeMintTokenEvent:
				sequenceNum, err := strconv.ParseInt(event.SequenceNumber, 10, 64)
				if err != nil {
					return err
				}
				tokenId := event.TokenEventData.(token.BurnTokenEvent).Id.ToString()
				amount := int64(event.TokenEventData.(token.BurnTokenEvent).Amount)
				tokenActivities = append(tokenActivities, &token.TokenActivityInDB{
					EventKey:       event.Key,
					SequenceNumber: sequenceNum,
					Account:        tx.Tx.Sender,
					TokenId:        tokenId,
					EventType:      event.Type,
					Amount:         amount,
					Version:        tx.Tx.Version,
					Timestamp:      0,
				})
				if tokenData, ok := tokenDataChange[tokenId]; !ok {
					tokenDataChange[tokenId] -= amount
				} else {
					tokenData += amount
					tokenDataChange[tokenId] = tokenData
				}

			case token.TypeTokenListingEvent:
				//todo: deal with listing
			case token.TypeTokenSwapEvent:
				sequenceNum, err := strconv.ParseInt(event.SequenceNumber, 10, 64)
				if err != nil {
					return err
				}
				timestamp, err := strconv.ParseInt(tx.Tx.Timestamp, 10, 64)
				if err != nil {
					return nil
				}

				coinType := event.TokenEventData.(token.TokenSwapEvent).CoinTypeInfo
				tokenTransferEvents = append(tokenTransferEvents, &token.TokenTransferEventInDB{
					Version:        tx.Tx.Version,
					EventKey:       event.Key,
					SequenceNumber: sequenceNum,
					Caller:         tx.Tx.Sender,
					To:             event.TokenEventData.(token.TokenSwapEvent).TokenBuyer,
					EventType:      event.Type,
					//todo:tokenId
					TokenId:     event.TokenEventData.(token.TokenSwapEvent).TokenId.ToString(),
					CoinType:    coinType.ToString(),
					TokenAmount: event.TokenEventData.(token.TokenSwapEvent).TokenAmount,
					CoinAmount:  event.TokenEventData.(token.TokenSwapEvent).CoinAmount,
					Timestamp:   timestamp,
				})

			//todo: token transfer event 后面统一处理一下吧
			case token.TypeTokenOfferEvent:
				sequenceNum, err := strconv.ParseInt(event.SequenceNumber, 10, 64)
				if err != nil {
					return err
				}
				timestamp, err := strconv.ParseInt(tx.Tx.Timestamp, 10, 64)
				if err != nil {
					return nil
				}
				tokenId := event.TokenEventData.(token.TokenOfferEvent).TokenId.ToString()
				amount := int64(event.TokenEventData.(token.TokenOfferEvent).Amount)
				tokenActivities = append(tokenActivities, &token.TokenActivityInDB{
					EventKey:       event.Key,
					SequenceNumber: sequenceNum,
					Account:        tx.Tx.Sender,
					TokenId:        tokenId,
					EventType:      event.Type,
					Amount:         amount,
					Version:        tx.Tx.Version,
					Timestamp:      timestamp,
				})

			case token.TypeTokenClaimEvent:
				sequenceNum, err := strconv.ParseInt(event.SequenceNumber, 10, 64)
				if err != nil {
					return err
				}
				timestamp, err := strconv.ParseInt(tx.Tx.Timestamp, 10, 64)
				if err != nil {
					return nil
				}
				tokenId := event.TokenEventData.(token.TokenClaimEvent).TokenId.ToString()
				amount := int64(event.TokenEventData.(token.TokenClaimEvent).Amount)
				tokenActivities = append(tokenActivities, &token.TokenActivityInDB{
					EventKey:       event.Key,
					SequenceNumber: sequenceNum,
					Account:        tx.Tx.Sender,
					TokenId:        tokenId,
					EventType:      event.Type,
					Amount:         amount,
					Version:        tx.Tx.Version,
					Timestamp:      timestamp,
				})

			case token.TypeTokenCancelOfferEvent:
				sequenceNum, err := strconv.ParseInt(event.SequenceNumber, 10, 64)
				if err != nil {
					return err
				}
				timestamp, err := strconv.ParseInt(tx.Tx.Timestamp, 10, 64)
				if err != nil {
					return nil
				}
				tokenId := event.TokenEventData.(token.TokenCancelOfferEvent).TokenId.ToString()
				amount := int64(event.TokenEventData.(token.TokenCancelOfferEvent).Amount)
				tokenActivities = append(tokenActivities, &token.TokenActivityInDB{
					EventKey:       event.Key,
					SequenceNumber: sequenceNum,
					Account:        tx.Tx.Sender,
					TokenId:        tokenId,
					EventType:      event.Type,
					Amount:         amount,
					Version:        tx.Tx.Version,
					Timestamp:      timestamp,
				})

			default:
				continue
			}
		}
	}

	return nil
}

func getCollection(event token.CollectionCreationEvent, tx *types.Transaction) (*token.CollectionInDB, error) {
	timestamp, err := strconv.ParseInt(tx.Timestamp, 10, 64)
	if err != nil {
		return nil, err
	}
	collection := &token.CollectionInDB{
		CollectionId: fmt.Sprintf("%s:%s", event.Creator, event.CollectionName),
		Creator:      event.Creator,
		Name:         event.CollectionName,
		Description:  event.Description,
		MaxAmount:    int64(event.Maximum),
		Uri:          event.Uri,
		Version:      tx.Version,

		InsertAt: timestamp,
	}
	return collection, err
}

func getTokenData(event token.CreateTokenDataEvent, tx *types.Transaction) (*token.TokenDataInDB, error) {
	timestamp, err := strconv.ParseInt(tx.Timestamp, 10, 64)
	if err != nil {
		return nil, err
	}

	royaltyPointsDenominator, err := strconv.ParseInt(event.RoyaltyPointsDenominator, 10, 64)
	if err != nil {
		return nil, err
	}

	propertyKeys, err := event.PropertyKeys.Marshal()
	if err != nil {
		return nil, err
	}
	propertyValues, err := event.PropertyValues.Marshal()
	if err != nil {
		return nil, err
	}
	propertyTypes, err := event.PropertyTypes.Marshal()
	if err != nil {
		return nil, err
	}

	tokenData := &token.TokenDataInDB{
		TokenDataId:              event.Id.ToString(),
		Creator:                  event.Id.Creator,
		Collection:               event.Id.Collection,
		Name:                     event.Id.Name,
		Description:              event.Description,
		MaxAmount:                int64(event.Maximum),
		Supply:                   0,
		Uri:                      event.Uri,
		RoyaltyPayeeAddress:      event.RoyaltyPayeeAddress,
		RoyaltyPointsDenominator: royaltyPointsDenominator,
		RoyaltyPointsNumerator:   event.RoyaltyPointsNumerator,
		PropertyKey:              propertyKeys,
		PropertyValues:           propertyValues,
		PropertyTypes:            propertyTypes,
		MintedAt:                 timestamp,
		LastMintedAt:             timestamp,
		Version:                  tx.Version,
		CreatedAt:                nil,
		UpdatedAt:                nil,
	}
	return tokenData, nil
}
