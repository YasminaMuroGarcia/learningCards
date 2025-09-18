package handlers

import (
	"learning-cards/internal/services"
	"log"
	"net/http"
	"strconv"

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

func (h *UserWordHandler) UpdateUserWord(c *gin.Context) {
	wordID := c.Param("wordID")

	id, err := strconv.ParseUint(wordID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid word ID"})
		return
	}

	var requestBody struct {
		Learned bool `json:"learned"`
	}
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.service.UpdateUserWord(uint(id), requestBody.Learned)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update word"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Word updated successfully"})
}
