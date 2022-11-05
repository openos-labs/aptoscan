package types

import (
	"encoding/json"
	"gorm.io/gorm"
	"time"
)

type MoveRecourse struct {
	TransactionVersion     int64
	WriteSetChangeIndex    int64
	TransactionBlockHeight int64
	Name                   string
	Type                   string
	Address                string
	Module                 string
	GenericTypeParams      []byte
	Data                   []byte
	IsDeleted              bool

	CreatedAt *time.Time `gorm:"autoCreateTime"`
	UpdatedAt *time.Time `gorm:"autoUpdateTime;not null"`
}

func (MoveRecourse) TableName() string {
	return "move_resources"
}

type MoveStructTag struct {
	Module            string
	Name              string
	GenericTypeParams []byte
}

func (m *MoveRecourse) fromWriteResource(change *Change, writeSetChangeIndex, transactionVersion, transactionBlockHeight int64) error {
	data, err := json.Marshal(&change.Data)
	if err != nil {
		return err
	}
	GenericTypeParams := MoveStructTag{}
	m.TransactionVersion = transactionVersion
	m.TransactionBlockHeight = transactionBlockHeight
	m.WriteSetChangeIndex = writeSetChangeIndex

	m.Type = change.Type
	m.Address = change.Address
	m.Module = change.Module
	m.GenericTypeParams = GenericTypeParams
	m.Data = nil

}

func AutoCreateMoveResourceTable(db *gorm.DB) error {
	err := db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").AutoMigrate(&MoveRecourse{})
	if err != nil {
		return err
	}
	return nil
}
