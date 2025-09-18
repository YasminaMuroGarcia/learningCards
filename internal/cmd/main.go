package main

import (
	"learning-cards/config"
	"learning-cards/internal/handlers"
	"learning-cards/internal/models"
	"learning-cards/internal/repository"
	"learning-cards/internal/services"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	db, err := initializeDatabase()
	if err != nil {
		log.Fatal("failed to connect to the database:", err)
	}

	migrate(db)

	words := []models.Word{
		{Word: "der Supermarkt", Translation: "el supermercado", Category: "basic, shopping"},
		{Word: "das Geld", Translation: "el dinero", Category: "basic, shopping"},
		{Word: "die Karte (Bankkarte)", Translation: "la tarjeta bancaria", Category: "basic, shopping"},
		{Word: "die Tasche", Translation: "el bolso", Category: "basic, shopping"},
		{Word: "der Tampon", Translation: "el tampón", Category: "basic, shopping"},
		{Word: "die Binde", Translation: "la compresa", Category: "basic, shopping"},
		{Word: "das Kondom", Translation: "el preservativo", Category: "basic, shopping"},
		{Word: "die Zahnbürste", Translation: "el cepillo de dientes", Category: "basic, shopping"},
		{Word: "die Hilfe", Translation: "la ayuda", Category: "basic, shopping"},
		{Word: "die Polizei", Translation: "la policía", Category: "basic, shopping"},
	}

	insertData(db, words)

	userWordRepo := repository.NewUserWordRepository(db)
	userWordService := services.NewUserWordService(userWordRepo)
	userWordHandler := handlers.NewUserWordHandler(userWordService)
	r := gin.Default()
	r.GET("/words/daily", userWordHandler.GetUserWordDueToday)
	r.PUT("/words/update/:wordID", userWordHandler.UpdateUserWord)

	// Set up the cron job
	setupCronJobs(userWordHandler)

	if err := r.Run(); err != nil {
		log.Fatal("failed to start server:", err)
	}
}

// initializeDatabase initializes the database connection
func initializeDatabase() (*gorm.DB, error) {
	dbConfig := config.LoadDBConfig()
	dsn := "host=" + dbConfig.Host +
		" user=" + dbConfig.User +
		" password=" + dbConfig.Password +
		" port=" + dbConfig.Port +
		" sslmode=" + dbConfig.SSLMode +
		" TimeZone=Europe/Berlin" +
		" lc_messages=en_US" // Set error messages to English
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

// migrate does the migration for the database
func migrate(db *gorm.DB) {
	var allModels = []interface{}{&models.Word{}, &models.UserWord{}}
	if err := db.AutoMigrate(allModels...); err != nil {
		log.Println("Migration failed:", err)
	}
}

func insertData(db *gorm.DB, words []models.Word) {
	var count int64
	db.Model(&models.Word{}).Count(&count)
	if count == 0 {
		for _, word := range words {
			if err := db.Create(&word).Error; err != nil {
				log.Printf("failed to insert word %s: %v", word.Word, err)
			}
		}
	}
}

// setupCronJobs sets up the cron jobs based on the environment
func setupCronJobs(handler *handlers.UserWordHandler) {
	appConfig := config.LoadAppConfig()
	c := cron.New()
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatalf("Error getting hostname: %v", err)
	}
	if hostname == "localhost" || hostname == "127.0.0.1" || hostname == appConfig.Hostname {
		// Runs every minute if on localhost
		_, err = c.AddFunc("@every 1m", func() {
			if err := handler.SyncUserWords(); err != nil {
				log.Printf("Error syncing user words: %v", err)
			}
		})
		if err != nil {
			log.Fatalf("Error setting up cron job for localhost: %v", err)
		}
	} else {
		// Runs daily at midnight if not on localhost
		_, err = c.AddFunc("0 0 * * *", func() {
			if err := handler.SyncUserWords(); err != nil {
				log.Printf("Error syncing user words: %v", err)
			}
		})
		if err != nil {
			log.Fatalf("Error setting up cron job for production: %v", err)
		}
	}

	c.Start()
}
