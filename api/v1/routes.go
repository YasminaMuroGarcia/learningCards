package v1

import (
	"learning-cards/internal/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, userWordHandler *handlers.UserWordHandler) {
	r.GET("/words/daily", userWordHandler.GetUserWordDueToday)
	r.PUT("/words/update/:wordID", userWordHandler.UpdateUserWord)
}
