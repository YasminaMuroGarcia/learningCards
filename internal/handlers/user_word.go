package handlers

import (
	"learning-cards/internal/services"
	"log"
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

	allWords, err := h.service.GetAllWords()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve all words."})
		return
	}

	existingUserWords := make(map[uint]struct{})
	for _, userWord := range userWords {
		existingUserWords[userWord.WordID] = struct{}{}
	}

	for _, word := range allWords {
		if _, exists := existingUserWords[word.ID]; !exists {
			err := h.service.AddUserWord(word.ID)
			if err != nil {
				log.Printf("Error adding word %d: %v", word.ID, err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add word to user words."})
				return
			}
		}
	}

	c.JSON(http.StatusOK, userWords)
}
