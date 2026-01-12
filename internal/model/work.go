package model

import "time"

// Work represents a Bible work/translation (e.g., kjv, heb-wlc, grc-tr1894).
type Work struct {
	ID          string `gorm:"primaryKey" json:"id"` // e.g., "kjv", "heb-wlc"
	Name        string `json:"name"`
	Lang        string `json:"lang"`
	Description string `json:"description"`
	License     string `json:"license"`
	SourceURL   string `json:"source_url"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (Work) TableName() string {
	return "works"
}
