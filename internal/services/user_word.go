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
func (s *UserWordService) GetUserDueToday() ([]models.UserWord, error) {
	return s.repo.GetWordsDueToday()
}
func (s *UserWordService) AddUserWord(wordID uint) error {
	return s.repo.AddUserWord(wordID)
}
func (s *UserWordService) GetAllWords() ([]models.Word, error) {
	return s.repo.GetAllWords()
}
