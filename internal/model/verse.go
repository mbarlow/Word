package model

import "time"

// Verse represents a single verse in the canonical dataset.
type Verse struct {
	ID        string `gorm:"primaryKey" json:"id"` // VID: work/osis/chapter/verse
	Work      string `gorm:"index" json:"work"`
	OSIS      string `gorm:"index" json:"osis"`
	BookSlug  string `json:"book_slug"`
	Chapter   int    `gorm:"index" json:"chapter"`
	Verse     int    `json:"verse"`
	Lang      string `json:"lang"`
	Text      string `json:"text"`
	Source    string `json:"source"`    // JSON: upstream, artifact, sha256
	Meta      string `json:"meta"`      // JSON: optional metadata
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Verse) TableName() string {
	return "verses"
}
