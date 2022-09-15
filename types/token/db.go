package token

import (
	"gorm.io/gorm"
	"time"
)

type CollectionInDB struct {
	CollectionId string
	Creator      string
	Name         string
	Description  string
	MaxAmount    int64
	Uri          string
	InsertAt     int64
	Version      int64

	CreatedAt *time.Time `gorm:"autoCreateTime"`
	UpdatedAt *time.Time `gorm:"autoUpdateTime;not null"`
}

func (CollectionInDB) TableName() string {
	return "collections"
}

type MetaDataInDB struct {
	TokenId              string
	Name                 string
	Symbol               string
	SellerFeeBasisPoints int64
	Description          string
	Image                string
	ExternalUrl          string
	AnimationUrl         string
	Attributes           []byte
	Properties           []byte
	Version              int64

	CreatedAt *time.Time `gorm:"autoCreateTime"`
	UpdatedAt *time.Time `gorm:"autoUpdateTime;not null"`
}

func (MetaDataInDB) TableName() string {
	return "metadatas"
}

type OwnershipInDB struct {
	OwnershipId string
	TokenId     string `gorm:"column:token_id"`
	TokenDataId string `gorm:"column:token_data_id"`
	Owner       string `gorm:"column:owner"`
	Amount      int64
	Version     int64

	CreatedAt *time.Time `gorm:"autoCreateTime"`
	UpdatedAt *time.Time `gorm:"autoUpdateTime;not null"`
}

func (OwnershipInDB) TableName() string {
	return "ownerships"
}

type TokenDataInDB struct {
	TokenDataId              string `gorm:"column:token_data_id"`
	Creator                  string
	Collection               string
	Name                     string
	Description              string
	MaxAmount                int64
	Supply                   int64
	Uri                      string
	RoyaltyPayeeAddress      string
	RoyaltyPointsDenominator int64
	RoyaltyPointsNumerator   int64
	PropertyKey              []byte
	PropertyValues           []byte
	PropertyTypes            []byte
	MintedAt                 int64
	LastMintedAt             int64
	Version                  int64

	CreatedAt *time.Time `gorm:"autoCreateTime"`
	UpdatedAt *time.Time `gorm:"autoUpdateTime;not null"`
}

func (TokenDataInDB) TableName() string {
	return "token_datas"
}

type TokenPropertyInDB struct {
	TokenId         string
	PreviousTokenId string
	PropertyKeys    []byte
	PropertyValues  []byte
	PropertyTypes   []byte
	Version         int64
	Timestamp       int64

	CreatedAt *time.Time `gorm:"autoCreateTime"`
	UpdatedAt *time.Time `gorm:"autoUpdateTime;not null"`
}

func (TokenPropertyInDB) TableName() string {
	return "token_propertys"
}

type EventInDB struct {
	TransactionHash string
	Key             string
	SequenceNumber  int64
	Type            string
	Data            []byte

	CreatedAt *time.Time `gorm:"autoCreateTime"`
	UpdatedAt *time.Time `gorm:"autoUpdateTime;not null"`
}

func (EventInDB) TableName() string {
	return "events"
}

type TokenTransferEventInDB struct {
	Version        int64
	EventKey       string
	SequenceNumber int64
	TokenSeller    string
	TokenBuyer     string
	EventType      string
	TokenId        string
	CoinType       string
	TokenAmount    int64
	CoinAmount     int64
	Timestamp      int64

	CreatedAt *time.Time `gorm:"autoCreateTime"`
	UpdatedAt *time.Time `gorm:"autoUpdateTime;not null"`
}

func (TokenTransferEventInDB) TableName() string {
	return "token_transfer_events"
}

type TokenActivityInDB struct {
	Version        int64
	EventKey       string
	SequenceNumber int64

	EventType string
	Amount    int64
	Timestamp int64

	From    string
	To      string
	TokenId string
	Caller  string

	CreatedAt *time.Time `gorm:"autoCreateTime"`
	UpdatedAt *time.Time `gorm:"autoUpdateTime;not null"`
}

func (TokenActivityInDB) TableName() string {
	return "token_activitys"
}

type PendingTransfer struct {
	PendingId string `gorm:"column:pending_id"`
	TokenId   string
	From      string
	To        string
	Version   int64
	Timestamp int64
	Amount    int64

	CreatedAt *time.Time `gorm:"autoCreateTime"`
	UpdatedAt *time.Time `gorm:"autoUpdateTime;not null"`
}

func (PendingTransfer) TableName() string {
	return "pending_tokens"
}

func AutoCreateTokensTable(db *gorm.DB) error {
	err := db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").AutoMigrate(&CollectionInDB{})
	if err != nil {
		return err
	}
	err = db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").AutoMigrate(&MetaDataInDB{})
	if err != nil {
		return err
	}
	err = db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").AutoMigrate(&OwnershipInDB{})
	if err != nil {
		return err
	}
	err = db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").AutoMigrate(&TokenDataInDB{})
	if err != nil {
		return err
	}
	err = db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").AutoMigrate(&TokenPropertyInDB{})
	if err != nil {
		return err
	}
	err = db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").AutoMigrate(&EventInDB{})
	if err != nil {
		return err
	}
	err = db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").AutoMigrate(&TokenActivityInDB{})
	if err != nil {
		return err
	}
	err = db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").AutoMigrate(&TokenTransferEventInDB{})
	if err != nil {
		return err
	}
	err = db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").AutoMigrate(&PendingTransfer{})
	if err != nil {
		return err
	}
	return nil
}
