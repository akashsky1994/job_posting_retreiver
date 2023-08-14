package model

import (
	"time"

	"gorm.io/gorm"
)

type JobListing struct {
	gorm.Model
	ID        uint      `gorm:"primary_key;autoIncrement"`
	JobLink   string    `gorm:"unique;column:job_link;not_null;type:varchar(255);"`
	JobTitle  string    `gorm:"column:job_title;not_null;type:varchar(100);"`
	Location  []string  `gorm:"column:location;type:varchar[];"`
	Remote    bool      `gorm:"column:remote;"`
	Status    string    `gorm:"column:status;default:active;not_null"`
	CreatedAt time.Time `gorm:"column:created_at;not_null"`
	UpdatedAt time.Time `gorm:"column:updated_at;not_null"`
	CompanyID uint
	Company   Company `gorm:"foreignKey:CompanyID;references:ID"`
	OrgName   string  `gorm:"column:org_name;not_null;type:varchar(100);"`
	Source    string
}

type Company struct {
	ID   uint   `gorm:"primary_key;not_null;autoIncrement"`
	Name string `gorm:"uniqueIndex:cname,sort:desc"`
	Flag bool   `gorm:"default:1"`
}

type Link struct {
	ID uint `gorm:"primary_key;not_null;autoIncrement"`
}
