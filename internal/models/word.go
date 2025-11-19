package models

import (
	"time"
)

type Word struct {
	ID          uint      `gorm:"primary_key"`
	Word        string    `gorm:"size:255"`
	Translation string    `gorm:"size:255"`
	Category    string    `gorm:"size:255"`
	CreatedAt   time.Time `gorm:"DEFAULT:CURRENT_TIMESTAMP"`
}
