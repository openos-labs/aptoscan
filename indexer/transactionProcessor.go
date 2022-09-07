package indexer

import (
	"apotscan/types"
	"fmt"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"math"
	"strconv"
)

type TransactionProcessor interface {
	Name() string
	ProcessTransaction(transactions [][]types.Transaction, startVersion, endVersion uint64) (ProcessResult, error)
	GetDB() *gorm.DB
	GetRedis() *redis.Client
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

func (p *Processor) ProcessTransactionsWithStatus(txns []types.Transaction) error {
	if len(txns) == 0 {

	}
	return nil
}
