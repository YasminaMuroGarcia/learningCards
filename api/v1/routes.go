package v1

import (
	"learning-cards/internal/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, userWordHandler *handlers.UserWordHandler) {
	r.GET("/v1/words/daily", userWordHandler.GetUserWordDueToday)
	r.GET("/v1/words/category/:category", userWordHandler.GetUserWordsByCategory)
	r.PUT("/v1/words/update/:wordID", userWordHandler.UpdateUserWord)

}
