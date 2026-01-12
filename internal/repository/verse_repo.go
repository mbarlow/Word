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
