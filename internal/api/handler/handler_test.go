package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/mbarlow/word/internal/model"
	"github.com/mbarlow/word/internal/repository"
)

// newTestServer builds an echo instance with the v1 routes registered
// (mirroring cmd/server/main.go) over a seeded temp-file SQLite DB.
func newTestServer(t *testing.T) *echo.Echo {
	t.Helper()
	db, err := repository.NewDB(filepath.Join(t.TempDir(), "test.db"))
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
		{ID: "heb-wlc/GEN/1/1", Work: "heb-wlc", OSIS: "GEN", BookSlug: "genesis", Chapter: 1, Verse: 1, Lang: "he", Text: "בראשית ברא אלהים"},
		{ID: "kjv/JHN/1/1", Work: "kjv", OSIS: "JHN", BookSlug: "john", Chapter: 1, Verse: 1, Lang: "en", Text: "In the beginning was the Word."},
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

	e := echo.New()
	h := New(db)
	v1 := e.Group("/v1")
	v1.GET("/works", h.ListWorks)
	v1.GET("/books", h.ListBooks)
	v1.GET("/text/:work/:book/:chapter", h.GetChapter)
	v1.GET("/verse/:work/:book/:chapter/:verse", h.GetVerse)
	v1.GET("/compare", h.Compare)
	v1.GET("/search", h.Search)
	v1.GET("/random-verse", h.GetRandomVerse)
	v1.GET("/meaning/:work/:book/:chapter", h.GetMeaning)
	v1.GET("/meaning/:work/:book", h.GetMeaningSummary)
	return e
}

// get performs a GET against the test server and returns the recorder.
func get(e *echo.Echo, path string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodGet, path, nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	return rec
}

// decode unmarshals a recorder body into v, failing the test on error.
func decode(t *testing.T, rec *httptest.ResponseRecorder, v interface{}) {
	t.Helper()
	if err := json.Unmarshal(rec.Body.Bytes(), v); err != nil {
		t.Fatalf("unmarshal response %q: %v", rec.Body.String(), err)
	}
}

func TestListWorks(t *testing.T) {
	e := newTestServer(t)

	rec := get(e, "/v1/works")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}

	var works []model.Work
	decode(t, rec, &works)
	if len(works) != 2 {
		t.Fatalf("got %d works, want 2", len(works))
	}
	ids := map[string]bool{}
	for _, w := range works {
		ids[w.ID] = true
	}
	if !ids["kjv"] || !ids["heb-wlc"] {
		t.Errorf("missing expected work IDs, got %v", ids)
	}
}

func TestListBooks(t *testing.T) {
	e := newTestServer(t)

	rec := get(e, "/v1/books")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}

	var books []model.Book
	decode(t, rec, &books)
	if len(books) != 2 {
		t.Fatalf("got %d books, want 2", len(books))
	}
	if books[0].OSIS != "GEN" {
		t.Errorf("books[0].OSIS = %q, want GEN (canonical order)", books[0].OSIS)
	}
}

func TestGetChapter(t *testing.T) {
	e := newTestServer(t)

	tests := []struct {
		name       string
		path       string
		wantStatus int
		wantVerses int
	}{
		{"ok", "/v1/text/kjv/GEN/1", http.StatusOK, 2},
		{"lowercase book ok", "/v1/text/kjv/gen/1", http.StatusOK, 2},
		{"invalid chapter", "/v1/text/kjv/GEN/notanumber", http.StatusBadRequest, 0},
		{"missing chapter", "/v1/text/kjv/GEN/99", http.StatusNotFound, 0},
		{"missing work", "/v1/text/nope/GEN/1", http.StatusNotFound, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := get(e, tt.path)
			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d (body: %s)", rec.Code, tt.wantStatus, rec.Body.String())
			}
			if tt.wantStatus != http.StatusOK {
				return
			}
			var resp ChapterResponse
			decode(t, rec, &resp)
			if resp.Work != "kjv" || resp.Book != "GEN" || resp.Chapter != 1 {
				t.Errorf("unexpected chapter envelope: %+v", resp)
			}
			if len(resp.Verses) != tt.wantVerses {
				t.Errorf("got %d verses, want %d", len(resp.Verses), tt.wantVerses)
			}
		})
	}
}

func TestGetVerse(t *testing.T) {
	e := newTestServer(t)

	rec := get(e, "/v1/verse/kjv/GEN/1/1")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	var v model.Verse
	decode(t, rec, &v)
	if v.ID != "kjv/GEN/1/1" {
		t.Errorf("verse ID = %q, want kjv/GEN/1/1", v.ID)
	}
	if v.Text == "" {
		t.Error("verse text is empty")
	}

	if rec := get(e, "/v1/verse/kjv/GEN/99/99"); rec.Code != http.StatusNotFound {
		t.Errorf("missing verse status = %d, want 404", rec.Code)
	}
}

func TestCompare(t *testing.T) {
	e := newTestServer(t)

	tests := []struct {
		name       string
		path       string
		wantStatus int
		wantWorks  []string
	}{
		{"two works", "/v1/compare?works=kjv,heb-wlc&ref=GEN.1.1", http.StatusOK, []string{"kjv", "heb-wlc"}},
		{"missing params", "/v1/compare", http.StatusBadRequest, nil},
		{"bad ref format", "/v1/compare?works=kjv&ref=GEN.1", http.StatusBadRequest, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := get(e, tt.path)
			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d (body: %s)", rec.Code, tt.wantStatus, rec.Body.String())
			}
			if tt.wantStatus != http.StatusOK {
				return
			}
			var result map[string]model.Verse
			decode(t, rec, &result)
			if len(result) != len(tt.wantWorks) {
				t.Fatalf("got %d works in result, want %d", len(result), len(tt.wantWorks))
			}
			for _, w := range tt.wantWorks {
				if _, ok := result[w]; !ok {
					t.Errorf("result missing work %q", w)
				}
			}
		})
	}
}

func TestSearch(t *testing.T) {
	e := newTestServer(t)

	tests := []struct {
		name       string
		path       string
		wantStatus int
		wantCount  int
	}{
		{"matches across works", "/v1/search?q=beginning", http.StatusOK, 2},
		{"filtered by work", "/v1/search?q=beginning&work=kjv", http.StatusOK, 2},
		{"limit applies", "/v1/search?q=beginning&limit=1", http.StatusOK, 1},
		{"no results", "/v1/search?q=zebra", http.StatusOK, 0},
		{"missing q", "/v1/search", http.StatusBadRequest, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := get(e, tt.path)
			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d (body: %s)", rec.Code, tt.wantStatus, rec.Body.String())
			}
			if tt.wantStatus != http.StatusOK {
				return
			}
			var resp SearchResponse
			decode(t, rec, &resp)
			if resp.Count != tt.wantCount {
				t.Errorf("count = %d, want %d", resp.Count, tt.wantCount)
			}
			if len(resp.Results) != tt.wantCount {
				t.Errorf("got %d results, want %d", len(resp.Results), tt.wantCount)
			}
		})
	}
}

func TestGetRandomVerse(t *testing.T) {
	e := newTestServer(t)

	rec := get(e, "/v1/random-verse?work=kjv&testament=NT")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200 (body: %s)", rec.Code, rec.Body.String())
	}
	var v model.Verse
	decode(t, rec, &v)
	if v.OSIS != "JHN" {
		t.Errorf("osis = %q, want JHN (only NT book seeded)", v.OSIS)
	}

	if rec := get(e, "/v1/random-verse?testament=XX"); rec.Code != http.StatusBadRequest {
		t.Errorf("invalid testament status = %d, want 400", rec.Code)
	}
	if rec := get(e, "/v1/random-verse?work=nonexistent"); rec.Code != http.StatusNotFound {
		t.Errorf("no-match status = %d, want 404", rec.Code)
	}
}

func TestGetMeaning_NotFound(t *testing.T) {
	e := newTestServer(t)

	// The handler reads from data/meaning relative to the working directory,
	// which doesn't exist under the test's CWD — must 404, not 500.
	if rec := get(e, "/v1/meaning/heb-wlc/GEN/1"); rec.Code != http.StatusNotFound {
		t.Errorf("meaning status = %d, want 404", rec.Code)
	}
	if rec := get(e, "/v1/meaning/heb-wlc/GEN"); rec.Code != http.StatusNotFound {
		t.Errorf("meaning summary status = %d, want 404", rec.Code)
	}
	if rec := get(e, "/v1/meaning/heb-wlc/GEN/notanumber"); rec.Code != http.StatusBadRequest {
		t.Errorf("invalid chapter status = %d, want 400", rec.Code)
	}
}
