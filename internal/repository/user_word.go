package repository

import (
	"fmt"
	"learning-cards/internal/models"
	"time"

	"gorm.io/gorm"
)

type UserWordRepository struct {
	db *gorm.DB
}

func NewUserWordRepository(db *gorm.DB) *UserWordRepository {
	return &UserWordRepository{db: db}
}
func (ur *UserWordRepository) GetUserWords() ([]models.UserWord, error) {
	var userWords []models.UserWord
	if err := ur.db.Preload("Word").Find(&userWords).Error; err != nil {
		return nil, err
	}
	return userWords, nil
}

func (ur *UserWordRepository) GetWordsDueToday() ([]models.UserWord, error) {
	var userWords []models.UserWord
	startOfDay := time.Now().Truncate(24 * time.Hour) // Start of today
	endOfDay := startOfDay.Add(24 * time.Hour)        // End of today

	if err := ur.db.Preload("Word").Where("next_review >= ? AND next_review < ?", startOfDay, endOfDay).Find(&userWords).Error; err != nil {
		return nil, err
	}
	return userWords, nil
}

// GetAllWords Get all the words that are on the word table
func (ur *UserWordRepository) GetAllWords() ([]models.Word, error) {
	var words []models.Word
	if err := ur.db.Find(&words).Error; err != nil {
		return nil, err
	}
	return words, nil
}

// AddUserWord Add a new user word to the user_word table
func (ur *UserWordRepository) AddUserWord(wordID uint) error {
	userWord := models.UserWord{
		WordID:            wordID,
		BoxNumber:         1,
		LastReview:        time.Now(),
		NextReview:        time.Now(),
		CorrectAttempts:   0,
		IncorrectAttempts: 0,
	}
	fmt.Printf("Wer are here here hrer")
	err := ur.db.Create(&userWord).Error
	if err != nil {
		fmt.Printf("Error creating user word: %v", err)
	}
	return err
}
