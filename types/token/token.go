package token

import (
	"apotscan/types"
	"encoding/json"
	"fmt"
)

func GetTransactionsWithTokenEvent(txs []types.Transaction) ([]*TransactionWithTokenEvents, error) {
	var txsWithTokenEvent []*TransactionWithTokenEvents
	for _, tx := range txs {
		events, err := getTransactionWithTokenEvent(tx)
		if err != nil {
			return nil, err
		}
		if len(events) != 0 {
			txsWithTokenEvent = append(txsWithTokenEvent, &TransactionWithTokenEvents{
				tx,
				events,
			})
		}
	}
	return txsWithTokenEvent, nil
}

func getTransactionWithTokenEvent(tx types.Transaction) ([]TokenEvent, error) {
	var events []TokenEvent
	for _, event := range tx.Events {
		if tx.Type != types.UserTransaction {
			continue
		}
		data, err := json.Marshal(event.Data)
		if err != nil {
			return nil, fmt.Errorf("tx %d event %s can not be marshal with error %v", tx.Version, event.Key, err)
		}
		switch event.Type {
		case TypeWithdrawEvent:
			var e WithdrawEvent
			if err = json.Unmarshal(data, &e); err != nil {
				return nil, fmt.Errorf("tx %d event %s can not be unmarshal with error %v", tx.Version, event.Key, err)
			}
			events = append(events, TokenEvent{
				Key:            event.Key,
				SequenceNumber: event.SequenceNumber,
				Type:           event.Type,
				TokenEventData: e,
			})
		case TypeDepositEvent:
			var e DepositEvent
			if err = json.Unmarshal(data, &e); err != nil {
				return nil, fmt.Errorf("tx %d event %s can not be unmarshal with error %v", tx.Version, event.Key, err)
			}
			events = append(events, TokenEvent{
				Key:            event.Key,
				SequenceNumber: event.SequenceNumber,
				Type:           event.Type,
				TokenEventData: e,
			})
		case TypeCreateTokenDataEvent:
			var e CreateTokenDataEvent
			if err = json.Unmarshal(data, &e); err != nil {
				return nil, fmt.Errorf("tx %d event %s can not be unmarshal with error %v", tx.Version, event.Key, err)
			}
			events = append(events, TokenEvent{
				Key:            event.Key,
				SequenceNumber: event.SequenceNumber,
				Type:           event.Type,
				TokenEventData: e,
			})
		case TypeCollectionCreationEvent:
			var e CollectionCreationEventRaw
			if err = json.Unmarshal(data, &e); err != nil {
				return nil, fmt.Errorf("tx %d event %s can not be unmarshal with error %v", tx.Version, event.Key, err)
			}
			events = append(events, TokenEvent{
				Key:            event.Key,
				SequenceNumber: event.SequenceNumber,
				Type:           event.Type,
				TokenEventData: e.GetEvent(),
			})
		case TypeBurnTokenEvent:
			var e BurnTokenEvent
			if err = json.Unmarshal(data, &e); err != nil {
				return nil, fmt.Errorf("tx %d event %s can not be unmarshal with error %v", tx.Version, event.Key, err)
			}
			events = append(events, TokenEvent{
				Key:            event.Key,
				SequenceNumber: event.SequenceNumber,
				Type:           event.Type,
				TokenEventData: e,
			})
		case TypeMutateTokenPropertyMapEvent:
			var e MutateTokenPropertyMapEvent
			if err = json.Unmarshal(data, &e); err != nil {
				return nil, fmt.Errorf("tx %d event %s can not be unmarshal with error %v", tx.Version, event.Key, err)
			}
			events = append(events, TokenEvent{
				Key:            event.Key,
				SequenceNumber: event.SequenceNumber,
				Type:           event.Type,
				TokenEventData: e,
			})
		case TypeMintTokenEvent:
			var e MintTokenEvent
			if err = json.Unmarshal(data, &e); err != nil {
				return nil, fmt.Errorf("tx %d event %s can not be unmarshal with error %v", tx.Version, event.Key, err)
			}
			events = append(events, TokenEvent{
				Key:            event.Key,
				SequenceNumber: event.SequenceNumber,
				Type:           event.Type,
				TokenEventData: e,
			})
		case TypeTokenListingEvent:
			var e TokenListingEvent
			if err = json.Unmarshal(data, &e); err != nil {
				return nil, fmt.Errorf("tx %d event %s can not be unmarshal with error %v", tx.Version, event.Key, err)
			}
			events = append(events, TokenEvent{
				Key:            event.Key,
				SequenceNumber: event.SequenceNumber,
				Type:           event.Type,
				TokenEventData: e,
			})
		case TypeTokenSwapEvent:
			var e TokenSwapEvent
			if err = json.Unmarshal(data, &e); err != nil {
				return nil, fmt.Errorf("tx %d event %s can not be unmarshal with error %v", tx.Version, event.Key, err)
			}
			events = append(events, TokenEvent{
				Key:            event.Key,
				SequenceNumber: event.SequenceNumber,
				Type:           event.Type,
				TokenEventData: e,
			})
		default:
			continue
		}
	}
	return events, nil
}

type TransactionWithTokenEvents struct {
	Tx          types.Transaction
	TokenEvents []TokenEvent
}
