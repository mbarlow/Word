package store

import (
	"fmt"

	"github.com/glebarez/sqlite"
	"github.com/mbarlow/word/internal/model"
	"github.com/mbarlow/word/internal/pipeline"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// SQLiteStore manages the canonical verse database.
type SQLiteStore struct {
	db *gorm.DB
}

// NewSQLiteStore opens or creates the SQLite database.
func NewSQLiteStore(path string) (*SQLiteStore, error) {
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Auto-migrate schema
	if err := db.AutoMigrate(&model.Verse{}, &model.Work{}, &model.Book{}); err != nil {
		return nil, fmt.Errorf("failed to migrate schema: %w", err)
	}

	return &SQLiteStore{db: db}, nil
}

// DB returns the underlying gorm.DB for use by repository layer.
func (s *SQLiteStore) DB() *gorm.DB {
	return s.db
}

// LoadVerses inserts or updates verses in the database.
func (s *SQLiteStore) LoadVerses(verses []pipeline.Verse) error {
	for _, v := range verses {
		dbVerse := model.Verse{
			ID:       v.ID,
			Work:     v.Work,
			OSIS:     v.OSIS,
			BookSlug: v.BookSlug,
			Chapter:  v.Chapter,
			Verse:    v.Verse,
			Lang:     v.Lang,
			Text:     v.Text,
			Source:   fmt.Sprintf(`{"upstream":"%s","artifact":"%s"}`, v.Source.Upstream, v.Source.Artifact),
		}

		// Upsert: update if exists, create if not
		result := s.db.Where("id = ?", v.ID).First(&model.Verse{})
		if result.Error == gorm.ErrRecordNotFound {
			if err := s.db.Create(&dbVerse).Error; err != nil {
				return fmt.Errorf("failed to insert verse %s: %w", v.ID, err)
			}
		} else {
			if err := s.db.Model(&model.Verse{}).Where("id = ?", v.ID).Updates(&dbVerse).Error; err != nil {
				return fmt.Errorf("failed to update verse %s: %w", v.ID, err)
			}
		}
	}
	return nil
}

// LoadWork inserts or updates a work record.
func (s *SQLiteStore) LoadWork(work model.Work) error {
	result := s.db.Where("id = ?", work.ID).First(&model.Work{})
	if result.Error == gorm.ErrRecordNotFound {
		return s.db.Create(&work).Error
	}
	return s.db.Model(&model.Work{}).Where("id = ?", work.ID).Updates(&work).Error
}

// LoadBook inserts or updates a book record.
func (s *SQLiteStore) LoadBook(book model.Book) error {
	result := s.db.Where("osis = ?", book.OSIS).First(&model.Book{})
	if result.Error == gorm.ErrRecordNotFound {
		return s.db.Create(&book).Error
	}
	return s.db.Model(&model.Book{}).Where("osis = ?", book.OSIS).Updates(&book).Error
}

// LoadBooksFromCanon loads all books from the canon into the database.
func (s *SQLiteStore) LoadBooksFromCanon(canon *pipeline.Canon) error {
	for _, b := range canon.Books {
		book := model.Book{
			OSIS:      b.OSIS,
			Slug:      b.Slug,
			Name:      b.Name,
			Testament: b.Testament,
			Order:     b.Order,
			Chapters:  b.Chapters,
		}
		if err := s.LoadBook(book); err != nil {
			return err
		}
	}
	return nil
}

// GetVerseCount returns the total number of verses in the database.
func (s *SQLiteStore) GetVerseCount() (int64, error) {
	var count int64
	err := s.db.Model(&model.Verse{}).Count(&count).Error
	return count, err
}

// GetWorkCount returns the total number of works in the database.
func (s *SQLiteStore) GetWorkCount() (int64, error) {
	var count int64
	err := s.db.Model(&model.Work{}).Count(&count).Error
	return count, err
}
