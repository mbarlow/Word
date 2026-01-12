package model

// Book represents a canonical book of the Bible.
type Book struct {
	OSIS       string `gorm:"primaryKey" json:"osis"` // 3-letter OSIS code
	Slug       string `gorm:"uniqueIndex" json:"slug"`
	Name       string `json:"name"`
	Testament  string `json:"testament"` // OT or NT
	Order      int    `json:"order"`     // Canonical order (1-66)
	Chapters   int    `json:"chapters"`  // Number of chapters
}

func (Book) TableName() string {
	return "books"
}
