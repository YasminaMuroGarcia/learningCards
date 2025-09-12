package main

import (
	"learning-cards/config"
	"learning-cards/internal/models"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {

	dbConfig := config.LoadDBConfig()
	dsn := "host=" + dbConfig.Host + " user=" + dbConfig.User + " password=" + dbConfig.Password + " port=" + dbConfig.Port + " sslmode=" + dbConfig.SSLMode
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
