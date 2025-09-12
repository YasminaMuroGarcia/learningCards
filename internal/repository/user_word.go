package repository

import (
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
