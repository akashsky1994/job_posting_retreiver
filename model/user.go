package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID     uint   `gorm:"primary_key;column:id;autoIncrement" json:"id,omitempty"`
	Name   string `gorm:"column:name;not_null"`
	Active bool   `gorm:"column:active;not_null;default:true"`
}

type UserJob struct {
	gorm.Model
	ID         uint       `gorm:"primary_key;column:id;autoIncrement" json:"id,omitempty"`
	UserID     uint       `gorm:"not_null"`
	User       User       `gorm:"foreignKey:UserID"`
	JobID      uint       `gorm:"not_null"`
	JobListing JobListing `gorm:"foreignKey:JobID"`
}
