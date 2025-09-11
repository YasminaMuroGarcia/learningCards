package main

import (
	"learning-cards/internal/models"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "host=" + os.Getenv("DB_HOST") + " user=" + os.Getenv("DB_USER") + " password=" + os.Getenv("DB_PASSWORD") + " port=5432 sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to the database:", err)
	}

	migrate(db)

}

// migrate does the migration for the database
func migrate(db *gorm.DB) {
	var allModels = []interface{}{&models.Word{}, &models.UserWord{}}
	err := db.AutoMigrate(allModels...)
	if err != nil {
		log.Println("Migration failed:", err)
	}
}
