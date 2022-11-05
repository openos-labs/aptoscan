package token

import (
	"gorm.io/gorm"
	"time"
)

type Token struct {
	CreatorAddress     string
	CollectionNameHash string
	NameHash           string
	CollectionName     string
	Name               string
	PropertyVersion    int64
	TransactionVersion int64
	TokenProperties    []byte

	CreatedAt *time.Time `gorm:"autoCreateTime"`
	UpdatedAt *time.Time `gorm:"autoUpdateTime;not null"`
}

func (Token) TableName() string {
	return "tokens"
}

type CollectionData struct {
	CreatorAddress     string
	CollectionNameHash string
	CollectionName     string
	Description        string
	TransactionVersion int64
	MetadataUri        string
	Supply             int64
	Maximum            int64
	MaximumMutable     bool
	UriMutable         bool
	DescriptionMutable bool

	CreatedAt *time.Time `gorm:"autoCreateTime"`
	UpdatedAt *time.Time `gorm:"autoUpdateTime;not null"`
}

func (CollectionData) TableName() string {
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

type Ownership struct {
	CreatorAddress     string
	CollectionNameHash string
	NameHash           string
	CollectionName     string
	Name               string
	PropertyVersion    int64
	TransactionVersion int64
	OwnerAddress       string
	Amount             int64
	TableHandle        string
	TableType          string

	CreatedAt *time.Time `gorm:"autoCreateTime"`
	UpdatedAt *time.Time `gorm:"autoUpdateTime;not null"`
}

func (Ownership) TableName() string {
	return "ownerships"
}

type TokenData struct {
	CreatorAddress         string
	CollectionNameHash     string
	NameHash               string
	CollectionName         string
	Name                   string
	TransactionVersion     int64
	Maximum                int64
	Supply                 int64
	LargestPropertyVersion int64
	MetadataUri            string

	PayeeAddress             string
	RoyaltyPointsDenominator int64
	RoyaltyPointsNumerator   int64

	MaximumMutable     bool
	UriMutable         bool
	DescriptionMutable bool
	PropertiesMutable  bool
	RoyaltyMutable     bool
	DefaultProperties  []byte

	CreatedAt *time.Time `gorm:"autoCreateTime"`
	UpdatedAt *time.Time `gorm:"autoUpdateTime;not null"`
}

func (TokenData) TableName() string {
	return "token_datas"
}

type TokenProperty struct {
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

func (TokenProperty) TableName() string {
	return "token_propertys"
}

//type EventInDB struct {
//	TransactionHash string
//	Key             string
//	SequenceNumber  int64
//	Type            string
//	Data            []byte
//
//	CreatedAt *time.Time `gorm:"autoCreateTime"`
//	UpdatedAt *time.Time `gorm:"autoUpdateTime;not null"`
//}

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
	err := db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").AutoMigrate(&CollectionData{})
	if err != nil {
		return err
	}
	err = db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").AutoMigrate(&MetaDataInDB{})
	if err != nil {
		return err
	}
	err = db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").AutoMigrate(&Ownership{})
	if err != nil {
		return err
	}
	err = db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").AutoMigrate(&TokenData{})
	if err != nil {
		return err
	}
	err = db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").AutoMigrate(&TokenProperty{})
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
	err = db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").AutoMigrate(&Token{})
	if err != nil {
		return err
	}
	return nil
}
