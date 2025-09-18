package database

import (
	"fmt"
	"learning-cards/config"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Open() (*gorm.DB, error) {
	dbConfig := config.LoadDBConfig()
	dsn := fmt.Sprintf("host=%s user=%s password=%s port=%s sslmode=%s TimeZone=Europe/Berlin lc_messages=en_US",
		dbConfig.Host, dbConfig.User, dbConfig.Password, dbConfig.Port, dbConfig.SSLMode)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	// optionally ping
	sqlDB, err := db.DB()
	if err == nil {
		if err = sqlDB.Ping(); err != nil {
			log.Printf("warning: ping DB failed: %v", err)
		}
	}
	return db, nil
}
