package indexer

import (
	"apotscan/logger"
	"apotscan/types"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	aptos "github.com/portto/aptos-go-sdk/client"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"math"
	"sync"
)

var (
	ctx = context.Background()
)

type Tailor struct {
	TransactionFetcher TransactionFetcher
	processors         []*Processor
	db                 *gorm.DB
	redisCli           *redis.Client
	logger             *logger.Logger
}

func NewTailor(nodeUrl string, db *gorm.DB, config *logger.Config, redisCli *redis.Client) *Tailor {
	_logger, err := logger.New(config)
	if err != nil {
		panic(err)
	}
	transactionFetcher := NewFetcher(aptos.New(nodeUrl), 0)
	return &Tailor{
		TransactionFetcher: transactionFetcher,
		db:                 db,
		logger:             _logger,
		processors:         []*Processor{},
		redisCli:           redisCli,
	}
}

func (t *Tailor) RunMigrations() {

}

//CheckOrUpdateChainId If chain id doesn't exist, save it. Otherwise make sure that we're indexing the same chain
func (t *Tailor) CheckOrUpdateChainId() error {
	t.logger.Info("Checking if chain id is correct")
	newLedgerInfo, err := t.TransactionFetcher.FetchLedgerInfo()
	if err != nil {
		return err
	}

	result, err := t.redisCli.Get(ctx, types.LedgerInfoKey).Result()
	if err != nil && err != redis.Nil {
		return err
	} else if err == redis.Nil {
		t.logger.WithFields(log.Fields{
			"chain id": newLedgerInfo.ChainID,
		}).Info("Adding chain id into redis, continue indexing")
		data, err := json.Marshal(&newLedgerInfo)
		if err != nil {
			return err
		}

		if err = t.redisCli.Set(ctx, types.LedgerInfoKey, data, -1).Err(); err != nil {
			return err
		}
		return nil
	}

	var currentLedgerInfo types.LedgerInfo
	if err = json.Unmarshal([]byte(result), &currentLedgerInfo); err != nil {
		return err
	}

	if newLedgerInfo.ChainID != currentLedgerInfo.ChainID {
		t.logger.WithFields(log.Fields{
			"try to index chain": currentLedgerInfo.ChainID,
			"exist chain":        newLedgerInfo.ChainID,
		}).Panic("Wrong chain detected!")
	}
	return nil
}

func (t *Tailor) AddProcessor(processor *Processor) {
	t.logger.WithFields(log.Fields{
		"name": processor.TransactionProcessor.Name(),
	}).Info("Adding processor to indexer")
	t.processors = append(t.processors, processor)
}

//HandlePreviousErrors For all versions which have an `success=false` in the `processor_status` table, re-run them
//TODO: also handle gaps in sequence numbers (pg query for this is super easy)
func (t *Tailor) HandlePreviousErrors() error {
	//todo:
	return nil
}

//SetFetcherToLowestProcessorVersion Sets the version of the fetcher to the lowest version among all processors
func (t *Tailor) SetFetcherToLowestProcessorVersion() (uint64, error) {
	var lowest uint64
	lowest = math.MaxUint64
	for _, processor := range t.processors {
		maxVersion, err := processor.GetMaxVersion(t.TransactionFetcher.GetChainId())
		if err != nil {
			return lowest, err
		}
		t.logger.WithFields(log.Fields{
			"chain id":    t.TransactionFetcher.GetChainId(),
			"processor":   processor.TransactionProcessor.Name(),
			"max version": maxVersion,
		}).Debug("Process max version")
		if lowest > maxVersion {
			lowest = maxVersion
		}
	}
	return t.SetFetcherVersion(lowest)
}

func (t *Tailor) SetFetcherVersion(version uint64) (uint64, error) {
	t.TransactionFetcher.SetVersion(version)
	t.logger.Info(fmt.Sprintf("Will start fetching from version %d", version))
	return version, nil
}

func (t *Tailor) ProcessVersion(version uint64) ([]ProcessResult, error) {
	tx, err := t.GetTxn(version)
	if err != nil {
		return nil, err
	}
	return t.ProcessTransactions([]types.Transaction{*tx})
}

func (t *Tailor) ProcessNextBatch(batchSize uint8) (uint64, [][]ProcessResult, error) {
	return 0, nil, nil
}

func (t *Tailor) ProcessTransactions(transactions []types.Transaction) ([]ProcessResult, error) {
	var wg sync.WaitGroup
	for _, processor := range t.processors {
		wg.Add(1)
		go func(processor *Processor, txns []types.Transaction) {
			defer wg.Done()
			err := processor.processTransactionsWithStatus(transactions)
		}(processor, transactions)
	}

	return nil, nil
}

func (t *Tailor) GetTxn(version uint64) (*types.Transaction, error) {
	return t.TransactionFetcher.FetchVersion(version)
}
