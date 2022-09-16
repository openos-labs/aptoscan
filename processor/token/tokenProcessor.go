package token

import (
	"apotscan/logger"
	"apotscan/types"
	"apotscan/types/token"
	"fmt"
	mapset "github.com/deckarep/golang-set"
	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
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
	var ownershipIds []string
	ownershipSet := mapset.NewSet()

	var tokenDataChanges []*TokenDataAmountChange
	var tokenDataChangeIds []string
	tokenDataChangeSet := mapset.NewSet()

	var pendingTransferIds []string
	var pendingTransfers []*TokenTransferEvent
	pendingTransferSet := mapset.NewSet()

	for _, tx := range txsWithEvents {
		for _, event := range tx.TokenEvents {
			switch event.TokenEventData.EventType() {
			case token.TypeWithdrawEvent:
				//todo: 增加withdraw和deposit的event
				tokenId := event.TokenEventData.(token.WithdrawEvent).Id.ToString()
				tokenDataId := event.TokenEventData.(token.WithdrawEvent).Id.TokenDataId.ToString()
				ownershipId := fmt.Sprintf("%s::%s,", tokenId, tx.Tx.Sender)
				if !ownershipSet.Contains(ownershipId) {
					ownershipSet.Add(ownershipId)
					ownershipIds = append(ownershipIds, ownershipId)
				}

				ownershipChanges = append(ownershipChanges, &token.OwnershipInDB{
					TokenId:     tokenId,
					TokenDataId: tokenDataId,
					Owner:       tx.Tx.Sender,
					Amount:      -event.TokenEventData.(token.WithdrawEvent).Amount,
					Version:     tx.Tx.Version,
				})

			case token.TypeDepositEvent:
				tokenId := event.TokenEventData.(token.WithdrawEvent).Id.ToString()
				tokenDataId := event.TokenEventData.(token.WithdrawEvent).Id.TokenDataId.ToString()
				ownershipId := fmt.Sprintf("%s::%s,", tokenId, tx.Tx.Sender)
				if !ownershipSet.Contains(ownershipId) {
					ownershipSet.Add(ownershipId)
					ownershipIds = append(ownershipIds, ownershipId)
				}
				ownershipChanges = append(ownershipChanges, &token.OwnershipInDB{
					TokenId:     tokenId,
					TokenDataId: tokenDataId,
					Owner:       tx.Tx.Sender,
					Amount:      event.TokenEventData.(token.WithdrawEvent).Amount,
					Version:     tx.Tx.Version,
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
				tokenDataId := event.TokenEventData.(token.BurnTokenEvent).Id.TokenDataId.ToString()
				amount := int64(event.TokenEventData.(token.BurnTokenEvent).Amount)
				tokenActivities = append(tokenActivities, &token.TokenActivityInDB{
					EventKey:       event.Key,
					SequenceNumber: sequenceNum,
					Caller:         tx.Tx.Sender,
					TokenId:        tokenId,
					EventType:      event.Type,
					Amount:         amount,
					Version:        tx.Tx.Version,
					Timestamp:      timestamp,
					From:           tx.Tx.Sender,
				})
				if !tokenDataChangeSet.Contains(tokenDataId) {
					tokenDataChangeSet.Add(tokenDataId)
					tokenDataChangeIds = append(tokenDataChangeIds, tokenDataId)
					tokenDataChanges = append(tokenDataChanges, &TokenDataAmountChange{
						Amount:      -amount,
						TokenDataId: tokenDataId,
						Version:     tx.Tx.Version,
					})
				} else {
					tokenDataChanges = append(tokenDataChanges, &TokenDataAmountChange{
						Amount:      -amount,
						TokenDataId: tokenDataId,
						Version:     tx.Tx.Version,
					})
				}

			case token.TypeMutateTokenPropertyMapEvent:
				//todo:
				if err := insertTokenProperties(db, event.TokenEventData.(token.MutateTokenPropertyMapEvent), &tx.Tx); err != nil {
					return err
				}

			case token.TypeMintTokenEvent:
				sequenceNum, err := strconv.ParseInt(event.SequenceNumber, 10, 64)
				if err != nil {
					return err
				}
				tokenId := event.TokenEventData.(token.BurnTokenEvent).Id.ToString()
				tokenDataId := event.TokenEventData.(token.BurnTokenEvent).Id.TokenDataId.ToString()

				amount := int64(event.TokenEventData.(token.BurnTokenEvent).Amount)
				tokenActivities = append(tokenActivities, &token.TokenActivityInDB{
					EventKey:       event.Key,
					SequenceNumber: sequenceNum,
					Caller:         tx.Tx.Sender,
					TokenId:        tokenId,
					EventType:      event.Type,
					Amount:         amount,
					Version:        tx.Tx.Version,
					Timestamp:      0,
					//To:             tx.Tx.Sender,
				})
				if !tokenDataChangeSet.Contains(tokenDataId) {
					tokenDataChangeSet.Add(tokenDataId)
					tokenDataChangeIds = append(tokenDataChangeIds, tokenDataId)
					tokenDataChanges = append(tokenDataChanges, &TokenDataAmountChange{
						Amount:      amount,
						TokenDataId: tokenDataId,
						Version:     tx.Tx.Version,
					})
				} else {
					tokenDataChanges = append(tokenDataChanges, &TokenDataAmountChange{
						Amount:      amount,
						TokenDataId: tokenDataId,
						Version:     tx.Tx.Version,
					})
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
					TokenSeller:    tx.Tx.Sender,
					TokenBuyer:     event.TokenEventData.(token.TokenSwapEvent).TokenBuyer,
					EventType:      event.Type,
					//todo:tokenId
					TokenId:     event.TokenEventData.(token.TokenSwapEvent).TokenId.ToString(),
					CoinType:    coinType.ToString(),
					TokenAmount: event.TokenEventData.(token.TokenSwapEvent).TokenAmount,
					CoinAmount:  event.TokenEventData.(token.TokenSwapEvent).CoinAmount,
					Timestamp:   timestamp,
				})

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
				to := event.TokenEventData.(token.TokenOfferEvent).ToAddress
				tokenActivities = append(tokenActivities, &token.TokenActivityInDB{
					EventKey:       event.Key,
					SequenceNumber: sequenceNum,
					Caller:         tx.Tx.Sender,
					TokenId:        tokenId,
					EventType:      event.Type,
					Amount:         amount,
					Version:        tx.Tx.Version,
					Timestamp:      timestamp,

					From: tx.Tx.Sender,
					To:   to,
				})
				tokenTransferEvent := &TokenTransferEvent{
					TokenId:   tokenId,
					From:      tx.Tx.Sender,
					To:        to,
					Amount:    amount,
					Version:   tx.Tx.Version,
					Timestamp: timestamp,
				}

				pendingId := tokenTransferEvent.GetId()
				if !pendingTransferSet.Contains(pendingId) {
					pendingTransferSet.Add(pendingTransferSet)
					pendingTransfers = append(pendingTransfers, tokenTransferEvent)
					pendingTransferIds = append(pendingTransferIds, pendingId)
				}

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
				from := event.TokenEventData.(token.TokenOfferEvent).ToAddress
				tokenActivities = append(tokenActivities, &token.TokenActivityInDB{
					EventKey:       event.Key,
					SequenceNumber: sequenceNum,
					Caller:         tx.Tx.Sender,
					TokenId:        tokenId,
					EventType:      event.Type,
					Amount:         amount,
					Version:        tx.Tx.Version,
					Timestamp:      timestamp,

					From: from,
					To:   tx.Tx.Sender,
				})
				tokenTransferEvent := &TokenTransferEvent{
					TokenId:   tokenId,
					From:      from,
					To:        tx.Tx.Sender,
					Amount:    -amount,
					Version:   tx.Tx.Version,
					Timestamp: timestamp,
				}

				pendingId := tokenTransferEvent.GetId()
				if !pendingTransferSet.Contains(pendingId) {
					pendingTransferSet.Add(pendingTransferSet)
					pendingTransfers = append(pendingTransfers, tokenTransferEvent)
					pendingTransferIds = append(pendingTransferIds, pendingId)
				}

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
				to := event.TokenEventData.(token.TokenOfferEvent).ToAddress
				tokenActivities = append(tokenActivities, &token.TokenActivityInDB{
					EventKey:       event.Key,
					SequenceNumber: sequenceNum,
					Caller:         tx.Tx.Sender,
					TokenId:        tokenId,
					EventType:      event.Type,
					Amount:         amount,
					Version:        tx.Tx.Version,
					Timestamp:      timestamp,

					From: tx.Tx.Sender,
					To:   to,
				})
				tokenTransferEvent := &TokenTransferEvent{
					TokenId:   tokenId,
					From:      tx.Tx.Sender,
					To:        to,
					Amount:    -amount,
					Version:   tx.Tx.Version,
					Timestamp: timestamp,
				}

				pendingId := tokenTransferEvent.GetId()
				if !pendingTransferSet.Contains(pendingId) {
					pendingTransferSet.Add(pendingTransferSet)
					pendingTransfers = append(pendingTransfers, tokenTransferEvent)
					pendingTransferIds = append(pendingTransferIds, pendingId)
				}

			default:
				continue
			}
		}
	}

	if err := db.Save(&collections).Error; err != nil {
		return nil
	}

	if err := db.Save(&tokensData).Error; err != nil {
		return nil
	}

	if err := db.Save(&tokenTransferEvents).Error; err != nil {
		return nil
	}

	if err := db.Save(&tokenActivities).Error; err != nil {
		return nil
	}
	//todo: 并发
	if err := dealWithOwnerShips(db, ownershipChanges, ownershipIds); err != nil {
		return nil
	}

	if err := dealWithTokenDataChanges(db, tokenDataChanges, tokenDataChangeIds); err != nil {
		return nil
	}

	if err := dealWithPendingTransfers(db, pendingTransfers, pendingTransferIds); err != nil {
		return nil
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

func insertTokenProperties(db *gorm.DB, event token.MutateTokenPropertyMapEvent, tx *types.Transaction) error {
	keys, err := event.Keys.Marshal()
	if err != nil {
		return err
	}

	values, err := event.Values.Marshal()
	if err != nil {
		return err
	}

	_types, err := event.Types.Marshal()

	timestamp, err := strconv.ParseInt(tx.Timestamp, 10, 64)
	if err != nil {
		return nil

	}
	oldTokenId := event.OldId.ToString()
	newTokenId := event.NewID.ToString()
	if oldTokenId != newTokenId {
		var oldOwnerships []*token.OwnershipInDB
		var newOwnerships []*token.OwnershipInDB

		if err := db.Where("token_id = ?", oldTokenId).Find(&oldOwnerships).Error; err != nil {
			return err
		}

		for _, oldOwnership := range oldOwnerships {
			newOwnership := &token.OwnershipInDB{
				OwnershipId: fmt.Sprintf("%s::%s,", newTokenId, oldOwnership.OwnershipId),
				TokenId:     newTokenId,
				TokenDataId: oldOwnership.TokenDataId,
				Owner:       oldOwnership.Owner,
				Amount:      oldOwnership.Amount,
			}
			oldOwnership.Version = tx.Version
			oldOwnership.Amount = 0
			newOwnerships = append(newOwnerships, newOwnership)
		}
		if err := db.Save(&oldOwnerships).Error; err != nil {
			return err
		}

		if err := db.Save(&newOwnerships).Error; err != nil {
			return err
		}
	}

	var tokenProperty = token.TokenPropertyInDB{
		TokenId:         oldTokenId,
		PreviousTokenId: newTokenId,
		PropertyKeys:    keys,
		PropertyValues:  values,
		PropertyTypes:   _types,
		Version:         tx.Version,
		Timestamp:       timestamp,
	}

	return db.Save(&tokenProperty).Error
}

func dealWithOwnerShips(db *gorm.DB, ownershipChanges []*token.OwnershipInDB, ownershipIds []string) error {
	var ownerShipsInDb []*token.OwnershipInDB
	if err := db.Where("OwnershipId IN (?)", ownershipIds).Find(&ownerShipsInDb).Error; err != nil {
		return err
	}
	ownerShipInDbMap := make(map[string]*token.OwnershipInDB)
	for _, ownership := range ownerShipsInDb {
		ownerShipInDbMap[ownership.OwnershipId] = ownership
	}

	for _, ownership := range ownershipChanges {
		if ownershipInDb, ok := ownerShipInDbMap[ownership.OwnershipId]; ok {
			if ownershipInDb.Version >= ownership.Version {
				continue
			}
			ownershipInDb.Amount += ownership.Amount
			ownershipInDb.Version = ownership.Version
		} else {
			ownerShipInDbMap[ownership.OwnershipId] = ownership
		}
	}

	var newOwnerships []*token.OwnershipInDB
	for _, ownership := range newOwnerships {
		newOwnerships = append(newOwnerships, ownership)
	}

	return db.Save(&newOwnerships).Error
}

func dealWithTokenDataChanges(db *gorm.DB, tokenDataChanges []*TokenDataAmountChange, tokenDataChangeIds []string) error {
	var tokenDatasInDb []*token.TokenDataInDB
	if err := db.Where("token_data_id IN (?)", tokenDataChangeIds).Find(&tokenDatasInDb).Error; err != nil {
		return err
	}

	tokenDataInDbMap := make(map[string]*token.TokenDataInDB)
	for _, tokenDataInDb := range tokenDatasInDb {
		tokenDataInDbMap[tokenDataInDb.TokenDataId] = tokenDataInDb
	}

	for _, tokenDataAmountChange := range tokenDataChanges {
		if tokenDataInDb, ok := tokenDataInDbMap[tokenDataAmountChange.TokenDataId]; ok {
			if tokenDataInDb.Version >= tokenDataAmountChange.Version {
				continue
			}
			tokenDataInDb.Supply += tokenDataAmountChange.Amount
			tokenDataInDb.Version = tokenDataAmountChange.Version
			if tokenDataInDb.Supply < 0 {
				return fmt.Errorf("version:%d,tokenId:%s:creator:%s,collection:%s,token:%s supply less than 0", tokenDataAmountChange.Version, tokenDataInDb.TokenDataId, tokenDataInDb.Creator, tokenDataInDb.Collection, tokenDataInDb.Name)
			}

		} else {
			return fmt.Errorf("version:%d,tokenId:%s:creator:%s,collection:%s,token:%s; mint or burn un-created token", tokenDataAmountChange.Version, tokenDataInDb.TokenDataId, tokenDataInDb.Creator, tokenDataInDb.Collection, tokenDataInDb.Name)

		}
	}

	var newTokensData []*token.TokenDataInDB
	for _, tokenData := range tokenDataInDbMap {
		newTokensData = append(newTokensData, tokenData)
	}

	return db.Save(&newTokensData).Error
}

func dealWithPendingTransfers(db *gorm.DB, pendingTransfers []*TokenTransferEvent, pendingTransferIds []string) error {
	var pendingTransfersInDb []*token.PendingTransfer
	if err := db.Where("pending_id IN (?)", pendingTransferIds).Find(&pendingTransfersInDb).Error; err != nil {
		return err
	}
	pendingTransferInDbMap := make(map[string]*token.PendingTransfer)
	for _, pendingTransferInDb := range pendingTransfersInDb {
		pendingId := pendingTransferInDb.PendingId
		pendingTransferInDbMap[pendingId] = &token.PendingTransfer{
			PendingId: pendingId,
			TokenId:   pendingTransferInDb.TokenId,
			From:      pendingTransferInDb.From,
			To:        pendingTransferInDb.To,
			Version:   pendingTransferInDb.Version,
			Timestamp: pendingTransferInDb.Timestamp,
			Amount:    pendingTransferInDb.Amount,
		}
	}

	for _, pendingTransfer := range pendingTransfers {
		if pendingTransferInDb, ok := pendingTransferInDbMap[pendingTransfer.Id]; ok {
			if pendingTransferInDb.Version >= pendingTransfer.Version {
				continue
			}

			pendingTransferInDb.Amount += pendingTransfer.Amount
			pendingTransferInDb.Version = pendingTransfer.Version
			if pendingTransferInDb.Amount < 0 {
				return fmt.Errorf("version%d, pending transfer id : %s, pending transfer amount should > 0", pendingTransfer.Version, pendingTransfer.Id)
			}
		} else {
			if pendingTransfer.Amount < 0 {
				return fmt.Errorf("version%d, pending transfer id : %s, pending transfer amount should > 0", pendingTransfer.Version, pendingTransfer.Id)
			}
			pendingTransferInDbMap[pendingTransfer.Id] = &token.PendingTransfer{
				PendingId: pendingTransfer.GetId(),
				TokenId:   pendingTransfer.TokenId,
				From:      pendingTransfer.From,
				To:        pendingTransfer.To,
				Version:   pendingTransfer.Version,
				Timestamp: pendingTransfer.Timestamp,
				Amount:    pendingTransfer.Amount,
			}
		}
	}

	var newPendingToken []*token.PendingTransfer
	for _, pendingToken := range pendingTransferInDbMap {
		newPendingToken = append(newPendingToken, pendingToken)
	}
	return db.Save(&newPendingToken).Error
}

func getAllMetadata(c *http.Client, uris map[string]string, logger *logger.Logger) []*token.MetaDataInDB {
	var metadatas []*token.MetaDataInDB
	for tokenId, uri := range uris {
		metadata := getMetadata(c, tokenId, uri, logger)
		if metadata != nil {
			attributes, err := metadata.Attributes.Marshal()
			if err != nil {
				logger.WithFields(log.Fields{
					"error":   err,
					"uri":     uri,
					"tokenId": tokenId,
					"type":    ARWEAVE,
				}).Error("can not get marshal attributes")
				continue
			}

			properties, err := metadata.Properties.Marshal()
			if err != nil {
				logger.WithFields(log.Fields{
					"error":   err,
					"uri":     uri,
					"tokenId": tokenId,
					"type":    ARWEAVE,
				}).Error("can not get marshal properties")
				continue
			}
			metadatas = append(metadatas, &token.MetaDataInDB{
				TokenId:              tokenId,
				Name:                 metadata.Name,
				Symbol:               metadata.Symbol,
				SellerFeeBasisPoints: metadata.SellerFeeBasisPoints,
				Description:          metadata.Description,
				Image:                metadata.Image,
				ExternalUrl:          metadata.ExternalUrl,
				AnimationUrl:         metadata.AnimationUrl,
				Attributes:           attributes,
				Properties:           properties,
				CreatedAt:            nil,
				UpdatedAt:            nil,
			})
		}
	}
	return metadatas
}
