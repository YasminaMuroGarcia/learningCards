package main

import (
	_ "gorm.io/gorm"
)

type Word struct {
	ID          uint   `gorm:"primary_key"`
	Name        string `gorm:"size:255"`
	Repetitions uint   `gorm:"default:1"`
}
