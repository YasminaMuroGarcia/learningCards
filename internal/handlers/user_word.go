package handlers

import (
	"learning-cards/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserWordHandler struct {
	service *services.UserWordService
}

func NewUserWordHandler(service *services.UserWordService) *UserWordHandler {
	return &UserWordHandler{
		service: service,
	}
}

func (h *UserWordHandler) GetUserWords(c *gin.Context) {
	userWords, err := h.service.GetUserWords()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user words."})
		return
	}
	c.JSON(http.StatusOK, userWords)

}

func (h *UserWordHandler) GetUserWordDueToday(c *gin.Context) {
	userWords, err := h.service.GetUserDueToday()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user words for today."})
		return
	}
	c.JSON(http.StatusOK, userWords)
}
