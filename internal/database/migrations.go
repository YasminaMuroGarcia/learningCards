package database

import (
	"learning-cards/internal/models"
	"log"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	modelsList := []interface{}{
		&models.Word{},
		&models.UserWord{},
	}
	if err := db.AutoMigrate(modelsList...); err != nil {
		log.Println("migration failed:", err)
	}
}
