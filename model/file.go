package model

import (
	"gorm.io/gorm"
)

type FileLogs struct {
	gorm.Model
	ID       uint `gorm:"primary_key;autoIncrement"`
	Source   string
	FilePath string
	Status   string
}
