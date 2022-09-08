package types

import "time"

type ProcessorStatus struct {
	Name    string
	Version int64
	Success bool
	Detail  string

	CreatedAt *time.Time `gorm:"autoCreateTime"`
	UpdatedAt *time.Time `gorm:"autoUpdateTime;not null"`
}

func (ProcessorStatus) TableName() string {
	return "processor_statuses"
}

func NewProcessorStatus(name string, version int64, success bool, detail string) *ProcessorStatus {
	return &ProcessorStatus{
		Name:    name,
		Version: version,
		Success: success,
		Detail:  detail,
	}
}

func ProcessorStatusFromVersions(name string, startVersion, endVersion int64, success bool, detail string) []ProcessorStatus {
	var status []ProcessorStatus
	newProcessorStatus := NewProcessorStatus(name, 0, success, detail)
	for i := startVersion; i < endVersion; i++ {
		newProcessorStatus.Version = i
		status = append(status, *newProcessorStatus)
	}
	return status
}
