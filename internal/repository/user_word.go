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
	err := ur.db.Create(&userWord).Error
	if err != nil {
		fmt.Printf("Error creating user word: %v", err)
	}
	return err
}

func (ur *UserWordRepository) MarkAsLearned(wordID uint) error {
	var userWord models.UserWord

	if err := ur.db.Where("word_id = ?", wordID).First(&userWord).Error; err != nil {
		return err
	}

	if userWord.BoxNumber < 5 {
		userWord.BoxNumber++
	}
	switch userWord.BoxNumber {
	case 1:
		userWord.NextReview = time.Now().Add(24 * time.Hour) // 1 day
	case 2:
		userWord.NextReview = time.Now().Add(3 * 24 * time.Hour) // 3 days
	case 3:
		userWord.NextReview = time.Now().Add(7 * 24 * time.Hour) // 7 days
	case 4:
		userWord.NextReview = time.Now().Add(14 * 24 * time.Hour) // 2 weeks
	case 5:
		userWord.NextReview = time.Now().Add(30 * 24 * time.Hour) // 1 month
	}

	userWord.CorrectAttempts++
	userWord.LastReview = time.Now()
	err := ur.db.Save(&userWord).Error
	if err != nil {
		fmt.Printf("Error updating user word: %v", err)
		return err
	}
	return nil
}
