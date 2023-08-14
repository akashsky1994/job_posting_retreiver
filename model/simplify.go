package model

import "time"

type SimplifyRecord struct {
	JobTitle  string    `json:"title"`
	JobLink   string    `json:"url"`
	Company   string    `json:"company_name"`
	Locations []string  `json:"locations"`
	JobType   string    `json:"type"`
	Remote    bool      `json:"remote"`
	UpdatedAt time.Time `json:"updated_at"`
}
