package model

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type JobListing struct {
	ID        uint           `gorm:"primary_key;column:id;autoIncrement" json:"id,omitempty"`
	JobLink   string         `gorm:"unique;column:job_link;not_null;type:text;" json:"job_link,omitempty"`
	JobTitle  string         `gorm:"column:job_title;not_null;type:text;" json:"job_title,omitempty"`
	Locations pq.StringArray `gorm:"column:locations;type:text[];" json:"locations,omitempty"`
	Remote    bool           `gorm:"column:remote;" json:"remote,omitempty"`
	Status    string         `gorm:"column:status;default:active;not_null" json:"status,omitempty"`
	CreatedAt time.Time      `gorm:"column:created_at;not_null;type:timestamptz" json:"-"`
	UpdatedAt time.Time      `gorm:"column:updated_at;not_null;type:timestamptz" json:"-"`
	CompanyID uint           `gorm:"column:company_id;not_null" json:"-"`
	Company   Company        `gorm:"foreignKey:CompanyID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"company,omitempty"`
	OrgName   string         `gorm:"column:org_name;not_null;type:varchar(100);" json:"org_name,omitempty"`
	Source    string         `gorm:"column:source;not_null" json:"source,omitempty"`
	FileLogID uint           `gorm:"column:filelog_id;not_null" json:"-"`
	FileLog   FileLog        `gorm:"foreignKey:FileLogID;references:ID" json:"-"`
	DeletedAt gorm.DeletedAt
}

type Company struct {
	ID   uint   `gorm:"column:id;primary_key;not_null;autoIncrement" json:"id,omitempty"`
	Name string `gorm:"column:name;uniqueIndex:cname,sort:desc" json:"name,omitempty"`
	Flag bool   `gorm:"column:flag;default:1" json:"flag,omitempty"`
}

type Link struct {
	ID uint `gorm:"primary_key;not_null;autoIncrement"`
}

type FileLog struct {
	ID       uint   `gorm:"primary_key;autoIncrement"`
	Source   string `gorm:"column:source;not_null"`
	FilePath string `gorm:"unique;not_null"`
	Status   string `gorm:"column:status;not_null"`
}
