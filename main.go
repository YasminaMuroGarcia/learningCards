package main

import (
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

	migrate(db) // Call the migration function

	// Your application logic here
}

func migrate(db *gorm.DB) {
	err := db.AutoMigrate(&Word{})
	if err != nil {
		log.Println("Migration failed:", err)
	} else {
		log.Println("Migration successful: User table created.")
	}
}
