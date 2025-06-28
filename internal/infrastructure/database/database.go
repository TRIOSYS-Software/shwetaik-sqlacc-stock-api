package database

import (
	"shwetaik-sqlacc-stock-api/internal/config"

	firebird "github.com/flylink888/gorm-firebird"
	"gorm.io/gorm"
)

func NewConnection(config *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(firebird.Open(config.DBString), &gorm.Config{
		PrepareStmt: true,
	})
	return db, err
}

func Migrate(db *gorm.DB) {
	db.AutoMigrate()
}
