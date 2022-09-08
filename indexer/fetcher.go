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
	StartingVersion     int64
	CurrentVersion      int64
	HighestKnownVersion int64
	//todo:
	//TransactionsSenders []
}

func NewFetcher(client aptos.API, currentVersion int64) *Fetcher {
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
	f.HighestKnownVersion = int64(version)
	f.ChainId = uint8(res.ChainID)
	return nil
}

func (f *Fetcher) FetchNextBatch(batch int) ([]types.Transaction, error) {
	txs, err := f.Client.GetTransactions(int(f.CurrentVersion), batch)
	if err != nil {
		return nil, err
	}
	var transactions []types.Transaction
	for _, tx := range txs {
		transaction := new(types.Transaction)
		if err = transaction.FromAptos(tx); err != nil {
			return nil, err
		}
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
	if err = transaction.FromAptos(*tx); err != nil {
		return nil, err
	}
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

func (f *Fetcher) SetVersion(version int64) {
	if f.StartingVersion != math.MaxInt64 {
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
	FetchNextBatch(batch int) ([]types.Transaction, error)
	FetchVersion(version uint64) (*types.Transaction, error)
	FetchLedgerInfo() (*types.LedgerInfo, error)
	SetVersion(version int64)
	GetChainId() uint8
	Start()
}
