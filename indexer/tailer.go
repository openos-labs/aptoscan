package indexer

import (
	"apotscan/logger"
	aptos "github.com/portto/aptos-go-sdk/client"
	"gorm.io/gorm"
)

type Tailor struct {
	transactionFetcher TransactionFetcher
	processors         []TransactionProcessor
	db                 *gorm.DB
	logger             *logger.Logger
}

func NewTailor(nodeUrl string, db *gorm.DB, config *logger.Config) *Tailor {
	logger, err := logger.New(config)
	if err != nil {
		panic(err)
	}
	transactionFetcher := NewFetcher(aptos.New(nodeUrl), 0)
	return &Tailor{
		transactionFetcher: transactionFetcher,
		db:                 db,
		logger:             logger,
		processors:         []TransactionProcessor{},
	}
}

func (t *Tailor) RunMigrations() {

}

func (t *Tailor) CheckOrUpdateChainId() {
	t.logger.Info("Checking if chain id is correct")
}
