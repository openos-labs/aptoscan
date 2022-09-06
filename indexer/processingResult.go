package indexer

type ProcessResult struct {
	Name         string
	StartVersion uint64
	EndVersion   uint64
}

func NewProcessResult(name string, start, end uint64) *ProcessResult {
	return &ProcessResult{
		Name:         name,
		StartVersion: start,
		EndVersion:   end,
	}
}
