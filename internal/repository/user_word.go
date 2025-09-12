package repository

import (
	"learning-cards/internal/models"

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
	if err := ur.db.Find(&userWords).Error; err != nil {
		return nil, err
	}
	return userWords, nil
}
