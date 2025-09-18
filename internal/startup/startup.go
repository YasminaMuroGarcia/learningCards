package startup

import (
	v1 "learning-cards/api/v1"
	"learning-cards/internal/database"
	"learning-cards/internal/handlers"
	"learning-cards/internal/models"
	"learning-cards/internal/repository"
	"learning-cards/internal/services"
	"log"

	"github.com/gin-gonic/gin"
)

func Run() error {
	// DB
	db, err := database.Open()
	if err != nil {
		return err
	}
	database.Migrate(db)

	// repositories / services / handlers
	userWordRepo := repository.NewUserWordRepository(db)
	userWordService := services.NewUserWordService(userWordRepo)
	userWordHandler := handlers.NewUserWordHandler(userWordService)

	// insert sample words set is done via cron job; if you want immediate insert, call insert util here
	words := predefinedWords()

	// router
	r := gin.Default()
	v1.RegisterRoutes(r, userWordHandler)

	// cron
	if err := setupCron(userWordHandler, db, words); err != nil {
		log.Println("cron setup warning:", err)
	}

	return r.Run()
}

func predefinedWords() []models.Word {
	return []models.Word{
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
}
