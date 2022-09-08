package token

import "apotscan/types"

func GetTransactionsWithTokenEvent(txs []types.Transaction) []*TransactionWithTokenEvent {
	var txsWithTokenEvent []*TransactionWithTokenEvent
	for _, tx := range txs {
		event := getTransactionWithTokenEvent(tx)
		if event != nil {
			txsWithTokenEvent = append(txsWithTokenEvent, &TransactionWithTokenEvent{
				tx,
				event,
			})
		}
	}
	return txsWithTokenEvent
}

func getTransactionWithTokenEvent(tx types.Transaction) []TokenEvent {
	var events []TokenEvent
	for _, event := range tx.Events {
		switch event.Type {
		case TypeWithdrawEvent:

		case TypeDepositEvent:
		case TypeCreateTokenDataEvent:
		case TypeCollectionCreationEvent:
		case TypeBurnTokenEvent:
		case TypeMutateTokenPropertyMapEvent:
		case TypeMintTokenEvent:
		default:
			continue
		}
	}
	return events
}

type TransactionWithTokenEvent struct {
	Tx         types.Transaction
	TokenEvent []TokenEvent
}
