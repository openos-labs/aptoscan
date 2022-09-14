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
		"name": processor.Name(),
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
func (t *Tailor) SetFetcherToLowestProcessorVersion() (int64, error) {
	var lowest int64
	lowest = math.MaxInt64
	for _, processor := range t.processors {
		maxVersion, err := processor.getMaxVersion()
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

func (t *Tailor) SetFetcherVersion(version int64) (int64, error) {
	t.TransactionFetcher.SetVersion(version)
	t.logger.Info(fmt.Sprintf("Will start fetching from version %d", version))
	return version, nil
}

func (t *Tailor) ProcessVersion(version uint64) ([]processResult, error) {
	tx, err := t.GetTxn(version)
	if err != nil {
		return nil, err
	}
	return t.ProcessTransactions([]types.Transaction{*tx}), nil
}

//ProcessNextBatch todo: 这他吗的能并行执行么
func (t *Tailor) ProcessNextBatch(batchSize uint8, singleFetchTxs int) (uint64, []processResult, error) {
	txs, err := t.TransactionFetcher.FetchNextBatch(singleFetchTxs)
	if err != nil {
		return 0, nil, err
	}
	txsAmount := len(txs)
	singleGoroutineTxAmount := txsAmount / int(batchSize)
	if singleGoroutineTxAmount < 0 {
		results := t.ProcessTransactions(txs)
		return uint64(txsAmount), results, nil
	}
	var remainBarch = batchSize
	resultCh := make(chan []processResult)
	var wg sync.WaitGroup
	var results []processResult
	for i := 0; i < int(batchSize); i++ {
		wg.Add(1)
		var txs2Process []types.Transaction
		if i != int(batchSize)-1 {
			txs2Process = txs[i*singleGoroutineTxAmount : (i+1)*singleGoroutineTxAmount]
		} else {
			txs2Process = txs[i*singleGoroutineTxAmount:]
		}
		go func(i int) {
			defer wg.Done()
			resultCh <- t.ProcessTransactions(txs2Process)
		}(i)
	}
	for {
		select {
		case singleBatchResult := <-resultCh:
			remainBarch--
			results = append(results, singleBatchResult...)
			if remainBarch == 0 {
				return uint64(txsAmount), results, nil
			}
		}
	}
}

func (t *Tailor) ProcessTransactions(transactions []types.Transaction) []processResult {
	var txs []types.Transaction
	for _, tx := range transactions {
		if tx.Success == true {
			txs = append(txs, tx)
		}
	}

	var results []processResult
	resultCh := make(chan processResult)
	var remainingTasks = len(t.processors)
	for _, processor := range t.processors {
		go func(processor *Processor, txs []types.Transaction) {
			result, err := processor.processTransactionsWithStatus(txs)
			resultCh <- processResult{
				result: *result,
				error:  err,
			}
		}(processor, txs)
	}

	for {
		select {
		case result := <-resultCh:
			remainingTasks--
			results = append(results, result)
			if remainingTasks == 0 {
				return results
			}
		}
	}
}

func (t *Tailor) GetTxn(version uint64) (*types.Transaction, error) {
	return t.TransactionFetcher.FetchVersion(version)
}

type processResult struct {
	result types.ProcessResult
	error  error
}
