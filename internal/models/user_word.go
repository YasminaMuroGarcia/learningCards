package models

import "time"

type UserWord struct {
	ID                uint      `gorm:"primary_key,auto_increment"`
	WordID            uint      `gorm:"not null;index;unique"`
	BoxNumber         uint      `gorm:"default:1"`
	LastReview        time.Time `gorm:"DEFAULT:CURRENT_TIMESTAMP"`
	NextReview        time.Time `gorm:"DEFAULT:CURRENT_TIMESTAMP"`
	CorrectAttempts   uint      `gorm:"default:0"`
	IncorrectAttempts uint      `gorm:"default:0"`
	Word              Word      `gorm:"foreignKey:WordID"` // Specify the foreign key relationship
}
