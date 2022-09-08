package types

type ProcessResult struct {
	Name         string
	StartVersion int64
	EndVersion   int64
	Error        error
}

func NewProcessResult(name string, start, end int64) *ProcessResult {
	return &ProcessResult{
		Name:         name,
		StartVersion: start,
		EndVersion:   end,
	}
}
