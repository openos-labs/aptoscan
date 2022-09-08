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
	ChainId() uint8
	ProcessTransactions(transactions []types.Transaction, startVersion, endVersion int64) (*types.ProcessResult, error)
	GetDB() *gorm.DB
	GetRedis() *redis.Client
	GetLogger() *logger.Logger
}
type Processor struct {
	TransactionProcessor
}

//getMaxVersion Gets the highest version for this `TransactionProcessor` from the DB
//This is so we know where to resume from on restarts
func (p *Processor) getMaxVersion() (int64, error) {
	key := fmt.Sprintf(types.MaxVersionKey, p.Name(), p.ChainId())
	result, err := p.GetRedis().Get(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	latestVersion, err := strconv.ParseInt(result, 10, 64)
	if err != nil {
		return math.MaxUint64, err
	}
	return latestVersion, nil
}

//todo: 到时候还要重新看啥时候setMaxVersion合适
func (p *Processor) setMaxVersion(version int64) error {
	currentMaxVersion, err := p.getMaxVersion()
	if err != nil {
		return err
	}
	if currentMaxVersion >= version {
		return nil
	}
	versionStr := strconv.FormatInt(version, 64)
	return p.GetRedis().Set(ctx, fmt.Sprintf(types.MaxVersionKey, p.Name(), p.ChainId()), versionStr, -1).Err()
}

//processTransactionsWithStatus This is a helper method, tying together the other helper methods to allow tracking status in the DB
func (p *Processor) processTransactionsWithStatus(txns []types.Transaction) (*types.ProcessResult, error) {
	if len(txns) == 0 {
		p.GetLogger().Warning("must provide at least one transaction to this funtion")
	}
	startVersion := txns[0].Version
	endVersion := txns[len(txns)-1].Version
	//todo: PROCESSOR_INVOCATIONS
	if err := p.markVersionStarted(startVersion, endVersion); err != nil {
		return nil, err
	}
	result, err := p.ProcessTransactions(txns, startVersion, endVersion)
	if err != nil {
		return nil, err
	}
	if err = p.updateStatus(result); err != nil {
		return nil, err
	}
	return result, nil
}

//markVersionStarted Writes that a version has been started for this `TransactionProcessor` to the DB
func (p *Processor) markVersionStarted(startVersion, endVersion int64) error {
	p.GetLogger().WithFields(log.Fields{
		"name":          p.Name(),
		"start version": startVersion,
		"end version":   endVersion,
	}).Debug("marking processing versions started from")
	psms := []types.ProcessorStatus{{
		Name:         p.Name(),
		StartVersion: startVersion,
		EndVersion:   endVersion,
		Success:      false,
		Detail:       "",
	}}
	return p.applyProcessorStatus(psms)
}

//applyProcessorStatus Actually performs the write for a `ProcessorStatusModel` change set
func (p *Processor) applyProcessorStatus(psms []types.ProcessorStatus) error {
	db := p.GetDB()
	return db.Save(&psms).Error
}

//updateStatus Writes that a version has been completed successfully for this `TransactionProcessor` to the DB
func (p *Processor) updateStatus(result *types.ProcessResult) error {
	//todo: PROCESSOR_SUCCESSES, PROCESSOR_ERRORS
	p.GetLogger().WithFields(log.Fields{
		"name":          p.Name(),
		"start version": result.StartVersion,
		"end version":   result.EndVersion,
	}).Debug("marking processing versions started from")
	var psms []types.ProcessorStatus
	if result.Error == nil {
		//psms = types.ProcessorStatusFromVersions(p.Name(), result.StartVersion, result.EndVersion, true, "")
		psms = []types.ProcessorStatus{{
			Name:         p.Name(),
			StartVersion: result.StartVersion,
			EndVersion:   result.EndVersion,
			Success:      true,
			Detail:       "",
		}}
	} else {
		//psms = types.ProcessorStatusFromVersions(p.Name(), result.StartVersion, result.EndVersion, false, result.Error.Error())
		psms = []types.ProcessorStatus{{
			Name:         p.Name(),
			StartVersion: result.StartVersion,
			EndVersion:   result.EndVersion,
			Success:      false,
			Detail:       result.Error.Error(),
		}}
	}
	if err := p.setMaxVersion(result.EndVersion); err != nil {
		return err
	}
	return p.applyProcessorStatus(psms)
}

func (p *Processor) getErrorVersion() ([]int64, error) {
	db := p.GetDB()
	var statuses []int64
	if err := db.Select("version").Where("success = ?", false).Scan(&statuses).Error; err != nil {
		return nil, err
	}
	return statuses, nil
}
