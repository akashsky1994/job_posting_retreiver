package model

import "time"

type TrueUpRecord struct {
	JobTitle  string    `json:"title"`
	JobLink   string    `json:"url"`
	Company   string    `json:"company"`
	Location  string    `json:"location"`
	Remote    int       `json:"location_remote"`
	UpdatedAt time.Time `json:"updated_at"`
}
