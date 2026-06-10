package repository

import (
	"path/filepath"
	"testing"

	"github.com/mbarlow/word/internal/model"
	"gorm.io/gorm"
)

// newTestDB opens a temp-file SQLite DB with the schema migrated and a small
// fixture set seeded: two works, two books (one OT, one NT), five verses.
func newTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := NewDB(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("NewDB: %v", err)
	}

	works := []model.Work{
		{ID: "kjv", Name: "King James Version", Lang: "en"},
		{ID: "heb-wlc", Name: "Westminster Leningrad Codex", Lang: "he"},
	}
	books := []model.Book{
		{OSIS: "GEN", Slug: "genesis", Name: "Genesis", Testament: "OT", Order: 1, Chapters: 50},
		{OSIS: "JHN", Slug: "john", Name: "John", Testament: "NT", Order: 43, Chapters: 21},
	}
	verses := []model.Verse{
		{ID: "kjv/GEN/1/1", Work: "kjv", OSIS: "GEN", BookSlug: "genesis", Chapter: 1, Verse: 1, Lang: "en", Text: "In the beginning God created the heaven and the earth."},
		{ID: "kjv/GEN/1/2", Work: "kjv", OSIS: "GEN", BookSlug: "genesis", Chapter: 1, Verse: 2, Lang: "en", Text: "And the earth was without form, and void."},
		{ID: "kjv/GEN/1/3", Work: "kjv", OSIS: "GEN", BookSlug: "genesis", Chapter: 1, Verse: 3, Lang: "en", Text: "And God said, Let there be light: and there was light."},
		{ID: "kjv/JHN/1/1", Work: "kjv", OSIS: "JHN", BookSlug: "john", Chapter: 1, Verse: 1, Lang: "en", Text: "In the beginning was the Word."},
		{ID: "heb-wlc/GEN/1/1", Work: "heb-wlc", OSIS: "GEN", BookSlug: "genesis", Chapter: 1, Verse: 1, Lang: "he", Text: "בראשית ברא אלהים"},
	}

	if err := db.Create(&works).Error; err != nil {
		t.Fatalf("seed works: %v", err)
	}
	if err := db.Create(&books).Error; err != nil {
		t.Fatalf("seed books: %v", err)
	}
	if err := db.Create(&verses).Error; err != nil {
		t.Fatalf("seed verses: %v", err)
	}
	return db
}

func TestVerseRepo_GetByID(t *testing.T) {
	repo := NewVerseRepo(newTestDB(t))

	v, err := repo.GetByID("kjv/GEN/1/1")
	if err != nil {
		t.Fatalf("GetByID: %v", err)
	}
	if v.Work != "kjv" || v.OSIS != "GEN" || v.Chapter != 1 || v.Verse != 1 {
		t.Errorf("unexpected verse: %+v", v)
	}

	if _, err := repo.GetByID("kjv/GEN/99/99"); err == nil {
		t.Error("expected error for missing verse, got nil")
	}
}

func TestVerseRepo_GetChapter(t *testing.T) {
	repo := NewVerseRepo(newTestDB(t))

	verses, err := repo.GetChapter("kjv", "GEN", 1)
	if err != nil {
		t.Fatalf("GetChapter: %v", err)
	}
	if len(verses) != 3 {
		t.Fatalf("got %d verses, want 3", len(verses))
	}
	for i, v := range verses {
		if v.Verse != i+1 {
			t.Errorf("verses[%d].Verse = %d, want %d (not ordered)", i, v.Verse, i+1)
		}
	}

	empty, err := repo.GetChapter("kjv", "GEN", 99)
	if err != nil {
		t.Fatalf("GetChapter missing: %v", err)
	}
	if len(empty) != 0 {
		t.Errorf("got %d verses for missing chapter, want 0", len(empty))
	}
}

func TestVerseRepo_Search(t *testing.T) {
	repo := NewVerseRepo(newTestDB(t))

	tests := []struct {
		name  string
		work  string
		query string
		limit int
		want  int
	}{
		{"match across books", "", "beginning", 50, 2},
		{"filtered by work", "heb-wlc", "בראשית", 50, 1},
		{"work filter excludes", "heb-wlc", "beginning", 50, 0},
		{"limit applies", "", "beginning", 1, 1},
		{"no match", "", "zebra", 50, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			verses, err := repo.Search(tt.work, tt.query, tt.limit)
			if err != nil {
				t.Fatalf("Search: %v", err)
			}
			if len(verses) != tt.want {
				t.Errorf("got %d results, want %d", len(verses), tt.want)
			}
		})
	}
}

func TestVerseRepo_ListWorks(t *testing.T) {
	repo := NewVerseRepo(newTestDB(t))

	works, err := repo.ListWorks()
	if err != nil {
		t.Fatalf("ListWorks: %v", err)
	}
	if len(works) != 2 {
		t.Fatalf("got %d works, want 2", len(works))
	}
}

func TestVerseRepo_ListBooks(t *testing.T) {
	repo := NewVerseRepo(newTestDB(t))

	books, err := repo.ListBooks()
	if err != nil {
		t.Fatalf("ListBooks: %v", err)
	}
	if len(books) != 2 {
		t.Fatalf("got %d books, want 2", len(books))
	}
	if books[0].OSIS != "GEN" || books[1].OSIS != "JHN" {
		t.Errorf("books not in canonical order: %s, %s", books[0].OSIS, books[1].OSIS)
	}
}

func TestVerseRepo_GetRandomVerse(t *testing.T) {
	repo := NewVerseRepo(newTestDB(t))

	tests := []struct {
		name      string
		work      string
		book      string
		testament string
		check     func(t *testing.T, v *model.Verse)
	}{
		{"by work", "heb-wlc", "", "", func(t *testing.T, v *model.Verse) {
			if v.Work != "heb-wlc" {
				t.Errorf("work = %q, want heb-wlc", v.Work)
			}
		}},
		{"by book", "", "JHN", "", func(t *testing.T, v *model.Verse) {
			if v.OSIS != "JHN" {
				t.Errorf("osis = %q, want JHN", v.OSIS)
			}
		}},
		{"by testament", "", "", "NT", func(t *testing.T, v *model.Verse) {
			if v.OSIS != "JHN" {
				t.Errorf("osis = %q, want JHN (only NT book seeded)", v.OSIS)
			}
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, err := repo.GetRandomVerse(tt.work, tt.book, tt.testament)
			if err != nil {
				t.Fatalf("GetRandomVerse: %v", err)
			}
			tt.check(t, v)
		})
	}

	if _, err := repo.GetRandomVerse("nonexistent", "", ""); err == nil {
		t.Error("expected error for no matching verses, got nil")
	}
}
