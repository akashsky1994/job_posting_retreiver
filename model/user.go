package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID     uint   `gorm:"primary_key;column:id;autoIncrement" json:"id,omitempty"`
	Name   string `gorm:"column:name;not_null"`
	Active bool   `gorm:"column:active;not_null;default:true"`
}

type UserJob struct {
	gorm.Model
	ID         uint       `gorm:"primary_key;column:id;autoIncrement" json:"id,omitempty"`
	UserID     uint       `gorm:"not_null;index:user_job_ids,unique,type:btree"`
	User       User       `gorm:"foreignKey:UserID"`
	JobID      uint       `gorm:"not_null;index:user_job_ids,unique,type:btree"`
	JobListing JobListing `gorm:"foreignKey:JobID"`
	Applied    int8       `gorm:"default:0"`
}
