package indexer

import aptos "github.com/portto/aptos-go-sdk/client"

type TransactionProcessor interface {
	Name() string
	ProcessTransaction(transactions []aptos.Transactions, startVersion, endVersion uint64) (ProcessResult, error)
	ConnectionPool()
}
