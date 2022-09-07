package types

import aptos "github.com/portto/aptos-go-sdk/client"

type LedgerInfo struct {
	ChainID             uint64 `json:"chain_id"`
	Epoch               string `json:"epoch"`
	LedgerVersion       string `json:"ledger_version"`
	LedgerTimestamp     string `json:"ledger_timestamp"`
	OldestLedgerVersion string `json:"oldest_ledger_version"`
	OldestBlockHeight   string `json:"oldest_block_height"`
	BlockHeight         string `json:"block_height"`
	NodeRole            string `json:"node_role"`
}

func (l *LedgerInfo) FromAptos(ledger aptos.LedgerInfo) {
	*l = LedgerInfo{
		ChainID:             ledger.ChainID,
		Epoch:               ledger.Epoch,
		LedgerVersion:       ledger.LedgerVersion,
		LedgerTimestamp:     ledger.LedgerTimestamp,
		OldestLedgerVersion: ledger.OldestLedgerVersion,
		OldestBlockHeight:   ledger.OldestBlockHeight,
		BlockHeight:         ledger.BlockHeight,
		NodeRole:            ledger.NodeRole,
	}
}
