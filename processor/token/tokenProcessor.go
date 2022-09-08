package token

import (
	"apotscan/logger"
	"apotscan/types"
	"apotscan/types/token"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type TokenTransactionProcessor struct {
	db       *gorm.DB
	redisCli *redis.Client
	chainId  uint8
	name     string
	logger   *logger.Logger
}

func New(name string, redisCli *redis.Client, db *gorm.DB, chainId uint8, logConf *logger.Config) (*TokenTransactionProcessor, error) {
	token.Init()
	_logger, err := logger.New(logConf)
	if err != nil {
		return nil, err
	}
	return &TokenTransactionProcessor{
		db:       db,
		redisCli: redisCli,
		chainId:  chainId,
		name:     name,
		logger:   _logger,
	}, nil
}

func (tp *TokenTransactionProcessor) Name() string {
	return tp.name
}

func (tp *TokenTransactionProcessor) ChainId() uint8 {
	return tp.chainId
}

func (tp *TokenTransactionProcessor) GetDB() *gorm.DB {
	return tp.db
}

func (tp *TokenTransactionProcessor) GetRedis() *redis.Client {
	return tp.redisCli
}

func (tp *TokenTransactionProcessor) GetLogger() *logger.Logger {
	return tp.logger
}

func (tp *TokenTransactionProcessor) ProcessTransactions(txs []types.Transaction, startVersion, endVersion int64) (*types.ProcessResult, error) {
	var tokenUris = make(map[string]string)
	txsWithTokenEvent :=
}
