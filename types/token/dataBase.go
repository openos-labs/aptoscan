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
	OwnershipId string //todo: 去掉？
	TokenId     string `gorm:"column:token_id"`
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
	TokenDataId              string
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
	PropertyKeys    string
	PropertyValues  string
	PropertyTypes   string
	Version         int64

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
	Caller         string
	To             string
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
	Account        string
	TokenId        string
	EventType      string
	Amount         int64
	Timestamp      int64

	CreatedAt *time.Time `gorm:"autoCreateTime"`
	UpdatedAt *time.Time `gorm:"autoUpdateTime;not null"`
}

func (TokenActivityInDB) TableName() string {
	return "token_activitys"
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
	return nil
}
