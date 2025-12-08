package database

import (
	"shwetaik-sqlacc-stock-api/internal/config"
	"time"

	firebird "github.com/flylink888/gorm-firebird"
	"gorm.io/gorm"
)

func NewConnection(config *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(firebird.Open(config.DBString), &gorm.Config{
		PrepareStmt: true,
	})
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(5)                   // max connections to Firebird
	sqlDB.SetMaxIdleConns(5)                   // keep some idle for reuse
	sqlDB.SetConnMaxLifetime(30 * time.Minute) // recycle connections before Firebird kills them
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)  // close idle connections
	return db, nil
}

func Migrate(db *gorm.DB) {
	db.AutoMigrate()
}
