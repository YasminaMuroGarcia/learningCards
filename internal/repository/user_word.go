package repository

import (
	"fmt"
	"learning-cards/internal/models"
	"log"
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
	startOfDay := time.Now().Truncate(24 * time.Hour)
	endOfDay := startOfDay.Add(24 * time.Hour)

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
func (ur *UserWordRepository) UpdateLearningStatus(wordID uint, learned bool) error {
	var userWord models.UserWord

	if err := ur.db.Where("word_id = ?", wordID).First(&userWord).Error; err != nil {
		return err
	}

	now := time.Now()
	userWord.LastReview = now

	if !learned {
		// When the user failed a word, it goes directly to the first box
		userWord.IncorrectAttempts++
		userWord.BoxNumber = 1
		userWord.NextReview = now.Add(24 * time.Hour)
	} else {
		if userWord.BoxNumber < 5 {
			userWord.BoxNumber++
		}
		userWord.NextReview = calculateNextReview(userWord.BoxNumber, now)
		userWord.CorrectAttempts++
	}

	if err := ur.db.Save(&userWord).Error; err != nil {
		fmt.Printf("Error updating user word: %v", err)
		return err
	}
	return nil
}

func (ur *UserWordRepository) CheckUserWordExists(wordID uint) (bool, error) {
	var count int64
	err := ur.db.Model(&models.UserWord{}).Where("word_id = ?", wordID).Count(&count).Error
	if err != nil {
		log.Printf("Error checking existence of word %d: %v", wordID, err)
		return false, err
	}

	return count > 0, nil
}

func calculateNextReview(boxNumber uint, currentTime time.Time) time.Time {
	switch boxNumber {
	case 1:
		return currentTime.Add(24 * time.Hour)
	case 2:
		return currentTime.Add(3 * 24 * time.Hour)
	case 3:
		return currentTime.Add(7 * 24 * time.Hour)
	case 4:
		return currentTime.Add(14 * 24 * time.Hour)
	case 5:
		return currentTime.Add(30 * 24 * time.Hour)
	default:
		return currentTime
	}
}
