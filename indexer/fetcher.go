package indexer

import (
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

func (f *Fetcher) FetchNextBatch(batch uint64) ([]aptos.TransactionResp, error) {
	return f.Client.GetTransactions(int(f.CurrentVersion), int(batch))
}

func (f *Fetcher) FetchVersion(v uint64) (*aptos.TransactionResp, error) {
	return f.Client.GetTransactionByVersion(v)
}

func (f *Fetcher) FetchLedgerInfo() (*aptos.LedgerInfo, error) {
	return f.Client.LedgerInformation()
}

func (f *Fetcher) SetVersion(version uint64) {
	if f.StartingVersion != math.MaxUint64 {
		panic("TransactionFetcher already started!")
	}
	f.StartingVersion = version
}

//todo:
func (f *Fetcher) Start() {

}

type TransactionFetcher interface {
	FetchNextBatch(batch uint64) ([]aptos.TransactionResp, error)
	FetchVersion(version uint64) (*aptos.TransactionResp, error)
	FetchLedgerInfo() (*aptos.LedgerInfo, error)
	SetVersion(version uint64)
	Start()
}
