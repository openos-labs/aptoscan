package module

import (
	"apotscan/logger"
	"apotscan/types"
	"apotscan/types/module"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type ModuleTransactionProcessor struct {
	db       *gorm.DB
	redisCli *redis.Client
	chainId  uint8
	name     string
	logger   *logger.Logger
}

func New(name string, redisCli *redis.Client, db *gorm.DB, chainId uint8, logConf *logger.Config) (*ModuleTransactionProcessor, error) {
	_logger, err := logger.New(logConf)
	if err != nil {
		return nil, err
	}
	return &ModuleTransactionProcessor{
		db:       db,
		redisCli: redisCli,
		chainId:  chainId,
		name:     name,
		logger:   _logger,
	}, nil
}

func (mp *ModuleTransactionProcessor) Name() string {
	return mp.name
}

func (mp *ModuleTransactionProcessor) ChainId() uint8 {
	return mp.chainId
}

func (mp *ModuleTransactionProcessor) GetDB() *gorm.DB {
	return mp.db
}

func (mp *ModuleTransactionProcessor) GetRedis() *redis.Client {
	return mp.redisCli
}

func (mp *ModuleTransactionProcessor) GetLogger() *logger.Logger {
	return mp.logger
}

func (mp *ModuleTransactionProcessor) ProcessTransactions(txs []types.Transaction, startVersion, endVersion int64) (*types.ProcessResult, error) {
	var modulse []module.Module
	for _, tx := range txs {
		if tx.Type != types.UserTransaction || tx.Payload.Type != types.ModuleBundlePayload {
			continue
		}
	}
	return &types.ProcessResult{
		Name:         mp.Name(),
		StartVersion: startVersion,
		EndVersion:   endVersion,
	}, nil
}
