package module

import "gorm.io/gorm"

type Module struct {
	Creator     string
	Module      string
	Name        string
	Description string
	Version     int64
}

func (Module) TableName() string {
	return "modules"
}

func AutoCreateTokensTable(db *gorm.DB) error {
	err := db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").AutoMigrate(&Module{})
	if err != nil {
		return err
	}
	return nil
}
