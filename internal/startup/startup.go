package startup

import (
	v1 "learning-cards/api/v1"
	"learning-cards/internal/database"
	"learning-cards/internal/handlers"
	"learning-cards/internal/repository"
	"learning-cards/internal/services"
	"learning-cards/internal/utils"
	"log"
	"path/filepath"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Run() error {
	db, err := database.Open()
	if err != nil {
		return err
	}
	database.Migrate(db)

	userWordRepo := repository.NewUserWordRepository(db)
	userWordService := services.NewUserWordService(userWordRepo)
	userWordHandler := handlers.NewUserWordHandler(userWordService)

	csvPath := filepath.Join("data", "words.csv")
	csv, err := utils.ReadCSV(csvPath)
	if err != nil {
		return err
	}
	words, err := utils.ConvertToWords(csv)
	if err != nil {
		log.Printf("warning converting words to words: %v", err)
	}

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "PUT", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type"},
		AllowCredentials: true,
	}))
	v1.RegisterRoutes(r, userWordHandler)

	if err := setupCron(userWordHandler, db, words); err != nil {
		log.Println("cron setup warning:", err)
	}

	return r.Run()
}
