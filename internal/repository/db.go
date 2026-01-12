package repository

import (
	"github.com/glebarez/sqlite"
	"github.com/mbarlow/word/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDB(path string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}

	// Auto-migrate schema
	if err := db.AutoMigrate(&model.Verse{}, &model.Work{}, &model.Book{}); err != nil {
		return nil, err
	}

	return db, nil
}
