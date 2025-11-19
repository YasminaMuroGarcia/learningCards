package services

import (
	"learning-cards/internal/models"
	"learning-cards/internal/repository"
	"math/rand"
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
func (s *UserWordService) GetUserWordsDueToday() ([]models.UserWord, error) {
	words, err := s.repo.GetWordsDueToday()
	if err != nil {
		return nil, err
	}
	// Shuffle the words
	rand.Shuffle(len(words), func(i, j int) {
		words[i], words[j] = words[j], words[i]
	})
	return words, nil
}
func (s *UserWordService) AddUserWord(wordID uint) error {
	return s.repo.AddUserWord(wordID)
}
func (s *UserWordService) GetAllWords() ([]models.Word, error) {
	return s.repo.GetAllWords()
}
func (s *UserWordService) UpdateUserWord(wordID uint, learned bool) error {
	return s.repo.UpdateLearningStatus(wordID, learned)
}
func (s *UserWordService) CheckUserWordExists(wordID uint) (bool, error) {
	return s.repo.CheckUserWordExists(wordID)
}
func (s *UserWordService) GetUserWordByCategory(category string) ([]models.UserWord, error) {
	wordByCategory, err := s.repo.GetUserWordsByCategory(category)
	if err != nil {
		return nil, err
	}
	rand.Shuffle(len(wordByCategory), func(i, j int) {
		wordByCategory[i], wordByCategory[j] = wordByCategory[j], wordByCategory[i]
	})
	return wordByCategory, nil
}
