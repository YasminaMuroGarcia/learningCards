package services

import (
	"learning-cards/internal/models"
	"learning-cards/internal/repository"
)

type UserWordService struct {
	repo *repository.UserWordRepository
}

func NewUserWordService(repo *repository.UserWordRepository) *UserWordService {
	return &UserWordService{repo: repo}
}

func (s *UserWordService) GetUserWords() ([]models.UserWord, error) {
	return s.repo.GetUserWords()
}
