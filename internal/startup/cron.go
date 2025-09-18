package startup

import (
	"learning-cards/config"
	"learning-cards/internal/handlers"
	"learning-cards/internal/models"
	"log"
	"os"

	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

func setupCron(handler *handlers.UserWordHandler, db *gorm.DB, words []models.Word) error {
	appCfg := config.LoadAppConfig()
	c := cron.New()
	hostname, err := os.Hostname()
	if err != nil {
		return err
	}

	isLocal := hostname == "localhost" || hostname == "127.0.0.1" || hostname == appCfg.Hostname

	if isLocal {
		addCron(c, "@every 1m", func() {
			if err := handler.SyncUserWords(); err != nil {
				log.Printf("Error syncing user words: %v", err)
			}
		}, "localhost")
		addCron(c, "@every 1m", func() {
			insertData(db, words)
			log.Println("Inserted predefined words into the database.")
		}, "localhost")
	} else {
		addCron(c, "0 0 * * *", func() {
			if err := handler.SyncUserWords(); err != nil {
				log.Printf("Error syncing user words: %v", err)
			}
		}, "production")
		addCron(c, "0 1 * * *", func() {
			insertData(db, words)
			log.Println("Inserted predefined words into the database.")
		}, "production")
	}

	c.Start()
	return nil
}

func addCron(c *cron.Cron, schedule string, job func(), env string) {
	if _, err := c.AddFunc(schedule, job); err != nil {
		log.Fatalf("Error setting up cron job for %s: %v", env, err)
	}
}

func insertData(db *gorm.DB, words []models.Word) {
	var count int64
	db.Model(&models.Word{}).Count(&count)
	if count == 0 && len(words) > 0 {
		for _, w := range words {
			if err := db.Create(&w).Error; err != nil {
				log.Printf("failed to insert word %s: %v", w.Word, err)
			}
		}
	}
}
