package token

import (
	"apotscan/types/token"
	"encoding/json"
	mapset "github.com/deckarep/golang-set"
	aptos "github.com/portto/aptos-go-sdk/client"
	"io/ioutil"
	"testing"
)

func TestTokenTransactionProcessor_GetTransactions(t *testing.T) {
	transactionTypes := mapset.NewSet()
	transactionTypes.Add(token.TypeWithdrawEvent)
	transactionTypes.Add(token.TypeDepositEvent)
	//transactionTypes.Add(token.TypeCollectionCreationEvent)
	transactionTypes.Add(token.TypeBurnTokenEvent)
	transactionTypes.Add(token.TypeMintTokenEvent)
	transactionTypes.Add(token.TypeTokenListingEvent)
	transactionTypes.Add(token.TypeTokenOfferEvent)
	transactionTypes.Add(token.TypeTokenClaimEvent)
	transactionTypes.Add(token.TypeTokenCancelOfferEvent)

	var collectionCreationEvent token.CollectionCreationEvent
	var createTokenDataEvent token.CreateTokenDataEvent
	var transactions []aptos.TransactionResp

	cli := aptos.New("https://fullnode.devnet.aptoslabs.com")
	offset := 1000
	start := 0
	//start := 202000
	//start := 43000
getCollection:
	for {
		t.Logf("start:%d", start)
		txs, _ := cli.GetTransactions(start, offset)
		for _, tx := range txs {
			if tx.Success && tx.Type == "user_transaction" {
				for _, event := range tx.Events {
					if event.Type == token.TypeCollectionCreationEvent {
						data, _ := json.Marshal(&event.Data)
						var collectionCreationEventRaw token.CollectionCreationEventRaw
						if err := json.Unmarshal(data, &collectionCreationEventRaw); err != nil {
							t.Fatal(err)
						}
						collectionCreationEvent = collectionCreationEventRaw.GetEvent()
						t.Logf("%v", event)
						t.Logf("get new collection creator:%s,name:%s\n", collectionCreationEvent.Creator, collectionCreationEvent.CollectionName)
						transactions = append(transactions, tx)
						break getCollection
					}
				}
			}
		}
		start += offset
	}
getToken:
	for {
		t.Logf("start:%d", start)
		txs, _ := cli.GetTransactions(start, offset)
		for _, tx := range txs {
			if tx.Success && tx.Type == "user_transaction" {
				for _, event := range tx.Events {
					if event.Type == token.TypeCreateTokenDataEvent {
						data, _ := json.Marshal(&event.Data)
						if err := json.Unmarshal(data, &createTokenDataEvent); err != nil {
							t.Fatal(err)
						}
						t.Logf("%s\n", event)
						t.Logf("get new token creator:%s,name:%s,collection:%s\n", createTokenDataEvent.Id.Creator, createTokenDataEvent.Id.Name, createTokenDataEvent.Id.Collection)
						transactions = append(transactions, tx)
						break getToken
					}
				}
			}
		}
		start += offset
	}
	tokenDataId := createTokenDataEvent.Id.ToString()

	for transactionTypes.Cardinality() > 0 {
		t.Logf("start:%d", start)
		txs, _ := cli.GetTransactions(start, offset)
		for _, tx := range txs {
			if tx.Success && tx.Type == "user_transaction" {
				for _, event := range tx.Events {
					if transactionTypes.Contains(event.Type) {
						switch event.Type {
						case token.TypeWithdrawEvent:
							var withdrawEvent token.WithdrawEvent
							data, _ := json.Marshal(&event)
							if err := json.Unmarshal(data, &withdrawEvent); err != nil {
								t.Fatal(err)
							}
							if tokenDataId == withdrawEvent.Id.TokenDataId.ToString() {
								transactions = append(transactions, tx)
							}
							transactionTypes.Remove(token.TypeWithdrawEvent)
							t.Log("get withdraw event")

						case token.TypeDepositEvent:
							var depositEvent token.DepositEvent
							data, _ := json.Marshal(&event)
							if err := json.Unmarshal(data, &depositEvent); err != nil {
								t.Fatal(err)
							}
							if tokenDataId == depositEvent.Id.TokenDataId.ToString() {
								transactions = append(transactions, tx)
							}
							transactionTypes.Remove(token.TypeDepositEvent)
							t.Log("get deposit event")

						case token.TypeBurnTokenEvent:
							var burnTokenEvent token.BurnTokenEvent
							data, _ := json.Marshal(&event)
							if err := json.Unmarshal(data, &burnTokenEvent); err != nil {
								t.Fatal(err)
							}
							if tokenDataId == burnTokenEvent.Id.TokenDataId.ToString() {
								transactions = append(transactions, tx)
							}
							transactionTypes.Remove(token.TypeBurnTokenEvent)
							t.Log("get burn event")

						case token.TypeMintTokenEvent:
							var mintTokenEvent token.MintTokenEvent
							data, _ := json.Marshal(&event)
							if err := json.Unmarshal(data, &mintTokenEvent); err != nil {
								t.Fatal(err)
							}
							if tokenDataId == mintTokenEvent.Id.TokenDataId.ToString() {
								transactions = append(transactions, tx)
							}
							transactionTypes.Remove(token.TypeMintTokenEvent)
							t.Log("get mint event")

						case token.TypeTokenListingEvent:
							var tokenListingEvent token.TokenListingEvent
							data, _ := json.Marshal(&event)
							if err := json.Unmarshal(data, &tokenListingEvent); err != nil {
								t.Fatal(err)
							}
							if tokenDataId == tokenListingEvent.TokenId.TokenDataId.ToString() {
								transactions = append(transactions, tx)
							}
							transactionTypes.Remove(token.TypeTokenListingEvent)
							t.Log("get token listing event")

						case token.TypeTokenOfferEvent:
							var tokenOfferEvent token.TokenOfferEvent
							data, _ := json.Marshal(&event)
							if err := json.Unmarshal(data, &tokenOfferEvent); err != nil {
								t.Fatal(err)
							}
							if tokenDataId == tokenOfferEvent.TokenId.TokenDataId.ToString() {
								transactions = append(transactions, tx)
							}
							transactionTypes.Remove(token.TypeTokenOfferEvent)
							t.Log("get token offer event")

						case token.TypeTokenClaimEvent:
							var tokenClaimEvent token.TokenClaimEvent
							data, _ := json.Marshal(&event)
							if err := json.Unmarshal(data, &tokenClaimEvent); err != nil {
								t.Fatal(err)
							}
							if tokenDataId == tokenClaimEvent.TokenId.TokenDataId.ToString() {
								transactions = append(transactions, tx)
							}
							transactionTypes.Remove(token.TypeTokenClaimEvent)
							t.Log("get token claim event")

						case token.TypeTokenCancelOfferEvent:
							var tokenCancelOfferEvent token.TokenCancelOfferEvent
							data, _ := json.Marshal(&event)
							if err := json.Unmarshal(data, &tokenCancelOfferEvent); err != nil {
								t.Fatal(err)
							}
							if tokenDataId == tokenCancelOfferEvent.TokenId.TokenDataId.ToString() {
								transactions = append(transactions, tx)
							}
							transactionTypes.Remove(token.TypeTokenCancelOfferEvent)
							t.Log("get token claim event")
						}
					}
				}
			}
		}
		start += offset
	}
	data, _ := json.Marshal(&transactions)
	t.Log(ioutil.WriteFile("test_transactions.json", data, 0777))
}
