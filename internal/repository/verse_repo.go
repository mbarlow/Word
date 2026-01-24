package repository

import (
	"github.com/mbarlow/word/internal/model"
	"gorm.io/gorm"
)

type VerseRepo struct {
	db *gorm.DB
}

func NewVerseRepo(db *gorm.DB) *VerseRepo {
	return &VerseRepo{db: db}
}

func (r *VerseRepo) GetByID(id string) (*model.Verse, error) {
	var v model.Verse
	if err := r.db.First(&v, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &v, nil
}

func (r *VerseRepo) GetChapter(work, osis string, chapter int) ([]model.Verse, error) {
	var verses []model.Verse
	err := r.db.Where("work = ? AND osis = ? AND chapter = ?", work, osis, chapter).
		Order("verse ASC").
		Find(&verses).Error
	return verses, err
}

func (r *VerseRepo) Search(work, query string, limit int) ([]model.Verse, error) {
	var verses []model.Verse
	q := r.db.Where("text LIKE ?", "%"+query+"%")
	if work != "" {
		q = q.Where("work = ?", work)
	}
	err := q.Limit(limit).Find(&verses).Error
	return verses, err
}

func (r *VerseRepo) ListWorks() ([]model.Work, error) {
	var works []model.Work
	err := r.db.Find(&works).Error
	return works, err
}

func (r *VerseRepo) ListBooks() ([]model.Book, error) {
	var books []model.Book
	err := r.db.Order("\"order\" ASC").Find(&books).Error
	return books, err
}

// GetRandomVerse returns a random verse with optional filtering.
// work: filter by work code (e.g., "kjv")
// book: filter by OSIS book code (e.g., "PSA")
// testament: filter by testament ("OT" or "NT")
func (r *VerseRepo) GetRandomVerse(work, book, testament string) (*model.Verse, error) {
	var v model.Verse
	q := r.db.Model(&model.Verse{})

	if work != "" {
		q = q.Where("verses.work = ?", work)
	}
	if book != "" {
		q = q.Where("verses.osis = ?", book)
	}
	if testament != "" {
		// Use subquery to filter by testament from books table
		q = q.Where("verses.osis IN (SELECT osis FROM books WHERE testament = ?)", testament)
	}

	if err := q.Order("RANDOM()").Limit(1).First(&v).Error; err != nil {
		return nil, err
	}
	return &v, nil
}
