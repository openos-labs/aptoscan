package indexer

import (
	"apotscan/types"
	aptos "github.com/portto/aptos-go-sdk/client"
	"math"
	"strconv"
)

type Fetcher struct {
	Client              aptos.API
	ChainId             uint8
	StartingVersion     uint64
	CurrentVersion      uint64
	HighestKnownVersion uint64
	//todo:
	//TransactionsSenders []
}

func NewFetcher(client aptos.API, currentVersion uint64) *Fetcher {
	//todo:
	//TransactionsSenders []
	return &Fetcher{
		Client:              client,
		CurrentVersion:      currentVersion,
		HighestKnownVersion: currentVersion,
		StartingVersion:     math.MaxUint64,
	}
}

func (f *Fetcher) setHighestKnownVersion() error {
	res, err := f.Client.LedgerInformation()
	if err != nil {
		return err
	}
	version, err := strconv.Atoi(res.LedgerVersion)
	if err != nil {
		return err
	}
	f.HighestKnownVersion = uint64(version)
	f.ChainId = uint8(res.ChainID)
	return nil
}

func (f *Fetcher) FetchNextBatch(batch uint64) ([]types.Transaction, error) {
	txs, err := f.Client.GetTransactions(int(f.CurrentVersion), int(batch))
	if err != nil {
		return nil, err
	}
	var transactions []types.Transaction
	for _, tx := range txs {
		transaction := new(types.Transaction)
		transaction.FromAptos(tx)
		transactions = append(transactions, *transaction)
	}
	return transactions, nil
}

func (f *Fetcher) FetchVersion(v uint64) (*types.Transaction, error) {
	tx, err := f.Client.GetTransactionByVersion(v)
	if err != nil {
		return nil, err
	}
	transaction := new(types.Transaction)
	transaction.FromAptos(*tx)
	return transaction, nil
}

func (f *Fetcher) FetchLedgerInfo() (*types.LedgerInfo, error) {
	info, err := f.Client.LedgerInformation()
	if err != nil {
		return nil, err
	}
	ledgerInfo := new(types.LedgerInfo)
	ledgerInfo.FromAptos(*info)
	return ledgerInfo, nil
}

func (f *Fetcher) SetVersion(version uint64) {
	if f.StartingVersion != math.MaxUint64 {
		panic("TransactionFetcher already started!")
	}
	f.StartingVersion = version
}

func (f *Fetcher) GetChainId() uint8 {
	return f.ChainId
}

//todo:
func (f *Fetcher) Start() {

}

type TransactionFetcher interface {
	FetchNextBatch(batch uint64) ([]types.Transaction, error)
	FetchVersion(version uint64) (*types.Transaction, error)
	FetchLedgerInfo() (*types.LedgerInfo, error)
	SetVersion(version uint64)
	GetChainId() uint8
	Start()
}
