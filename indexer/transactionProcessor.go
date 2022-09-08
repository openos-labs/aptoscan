package indexer

import (
	"apotscan/logger"
	"apotscan/types"
	"fmt"
	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"math"
	"strconv"
)

type TransactionProcessor interface {
	Name() string
	ProcessTransaction(transactions [][]types.Transaction, startVersion, endVersion uint64) (ProcessResult, error)
	GetDB() *gorm.DB
	GetRedis() *redis.Client
	GetLogger() *logger.Logger
	ProcessTransactionsWithStatus(transactions []types.Transaction) (ProcessResult, error)
}
type Processor struct {
	TransactionProcessor TransactionProcessor
}

func (p *Processor) GetMaxVersion(chainId uint8) (uint64, error) {
	key := fmt.Sprintf(types.MaxVersionKey, p.TransactionProcessor.Name(), chainId)
	result, err := p.TransactionProcessor.GetRedis().Get(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	latestVersion, err := strconv.Atoi(result)
	if err != nil {
		return math.MaxUint64, err
	}
	return uint64(latestVersion), nil
}

func (p *Processor) processTransactionsWithStatus(txns []types.Transaction) error {
	if len(txns) == 0 {
		p.TransactionProcessor.GetLogger().Warning("must provide at least one transaction to this funtion")
	}
	//todo: PROCESSOR_INVOCATIONS
	p.markVersionStarted(txns[0].Version, txns[len(txns)-1].Version)
	return nil
}

func (p *Processor) markVersionStarted(startVersion, endVersion int64) {
	p.TransactionProcessor.GetLogger().WithFields(log.Fields{
		"start version": startVersion,
		"end version":   endVersion,
	}).Debug("marking processing versions started from")
	psms := types.ProcessorStatusFromVersions(p.TransactionProcessor.Name(), startVersion, endVersion, false, "")
}

func (p *Processor) applyProcessorStatus(psms []types.ProcessorStatus) error {
	db := p.TransactionProcessor.GetDB()
	chunks := getChunks()
}
