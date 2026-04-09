package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/mbarlow/word/internal/model"
	"github.com/mbarlow/word/internal/repository"
	"gorm.io/gorm"
)

// ChapterResponse represents a chapter with its verses.
type ChapterResponse struct {
	Work    string        `json:"work"`
	Book    string        `json:"book"`
	Chapter int           `json:"chapter"`
	Verses  []model.Verse `json:"verses"`
}

// SearchResponse represents search results.
type SearchResponse struct {
	Query   string        `json:"query"`
	Count   int           `json:"count"`
	Results []model.Verse `json:"results"`
}

// ErrorResponse represents an API error.
type ErrorResponse struct {
	Message string `json:"message"`
}

type Handler struct {
	verses *repository.VerseRepo
}

func New(db *gorm.DB) *Handler {
	return &Handler{
		verses: repository.NewVerseRepo(db),
	}
}

// ListWorks godoc
// @Summary      List available works
// @Description  Returns all Bible works/translations in the database.
// @Tags         works
// @Produce      json
// @Success      200 {array}  model.Work
// @Failure      500 {object} ErrorResponse
// @Router       /v1/works [get]
func (h *Handler) ListWorks(c echo.Context) error {
	works, err := h.verses.ListWorks()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, works)
}

// GetChapter godoc
// @Summary      Get chapter text
// @Description  Returns all verses for a given work, book, and chapter.
// @Tags         text
// @Produce      json
// @Param        work    path string true "Work code (e.g. kjv, heb-wlc)"
// @Param        book    path string true "OSIS book code (e.g. GEN, PSA)"
// @Param        chapter path int    true "Chapter number"
// @Success      200 {object} ChapterResponse
// @Failure      400 {object} ErrorResponse
// @Failure      404 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /v1/text/{work}/{book}/{chapter} [get]
func (h *Handler) GetChapter(c echo.Context) error {
	work := c.Param("work")
	book := strings.ToUpper(c.Param("book"))
	chapter, err := strconv.Atoi(c.Param("chapter"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid chapter")
	}

	verses, err := h.verses.GetChapter(work, book, chapter)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if len(verses) == 0 {
		return echo.NewHTTPError(http.StatusNotFound, "chapter not found")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"work":    work,
		"book":    book,
		"chapter": chapter,
		"verses":  verses,
	})
}

// GetVerse godoc
// @Summary      Get a single verse
// @Description  Returns a single verse by work, book, chapter, and verse number.
// @Tags         text
// @Produce      json
// @Param        work    path string true "Work code"
// @Param        book    path string true "OSIS book code"
// @Param        chapter path string true "Chapter number"
// @Param        verse   path string true "Verse number"
// @Success      200 {object} model.Verse
// @Failure      404 {object} ErrorResponse
// @Router       /v1/verse/{work}/{book}/{chapter}/{verse} [get]
func (h *Handler) GetVerse(c echo.Context) error {
	work := c.Param("work")
	book := strings.ToUpper(c.Param("book"))
	chapter := c.Param("chapter")
	verse := c.Param("verse")

	id := work + "/" + book + "/" + chapter + "/" + verse
	v, err := h.verses.GetByID(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "verse not found")
	}

	return c.JSON(http.StatusOK, v)
}

// Compare godoc
// @Summary      Compare verse across works
// @Description  Returns a verse from multiple works/translations side by side.
// @Tags         text
// @Produce      json
// @Param        works query string true "Comma-separated work codes (e.g. kjv,heb-wlc)"
// @Param        ref   query string true "Verse reference as BOOK.CHAPTER.VERSE (e.g. GEN.1.1)"
// @Success      200 {object} map[string]model.Verse
// @Failure      400 {object} ErrorResponse
// @Router       /v1/compare [get]
func (h *Handler) Compare(c echo.Context) error {
	worksParam := c.QueryParam("works")
	ref := c.QueryParam("ref")

	if worksParam == "" || ref == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "works and ref required")
	}

	works := strings.Split(worksParam, ",")
	parts := strings.Split(ref, ".")
	if len(parts) != 3 {
		return echo.NewHTTPError(http.StatusBadRequest, "ref format: BOOK.CHAPTER.VERSE")
	}

	book := strings.ToUpper(parts[0])
	chapter := parts[1]
	verse := parts[2]

	result := make(map[string]interface{})
	for _, work := range works {
		id := work + "/" + book + "/" + chapter + "/" + verse
		v, err := h.verses.GetByID(id)
		if err == nil {
			result[work] = v
		}
	}

	return c.JSON(http.StatusOK, result)
}

// Search godoc
// @Summary      Search verses
// @Description  Full-text search across verses, optionally filtered by work.
// @Tags         search
// @Produce      json
// @Param        q     query string true  "Search query"
// @Param        work  query string false "Filter by work code"
// @Param        limit query int    false "Max results (1-100, default 50)"
// @Success      200 {object} SearchResponse
// @Failure      400 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /v1/search [get]
func (h *Handler) Search(c echo.Context) error {
	query := c.QueryParam("q")
	work := c.QueryParam("work")
	limitStr := c.QueryParam("limit")

	if query == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "q required")
	}

	limit := 50
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	verses, err := h.verses.Search(work, query, limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"query":   query,
		"count":   len(verses),
		"results": verses,
	})
}

// ListBooks godoc
// @Summary      List canonical books
// @Description  Returns all 66 canonical books in order.
// @Tags         books
// @Produce      json
// @Success      200 {array}  model.Book
// @Failure      500 {object} ErrorResponse
// @Router       /v1/books [get]
func (h *Handler) ListBooks(c echo.Context) error {
	books, err := h.verses.ListBooks()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, books)
}

// GetRandomVerse godoc
// @Summary      Get a random verse
// @Description  Returns a random verse with optional filtering by work, book, or testament.
// @Tags         text
// @Produce      json
// @Param        work      query string false "Filter by work code"
// @Param        book      query string false "Filter by OSIS book code"
// @Param        testament query string false "Filter by testament (OT or NT)"
// @Success      200 {object} model.Verse
// @Failure      400 {object} ErrorResponse
// @Failure      404 {object} ErrorResponse
// @Router       /v1/random-verse [get]
func (h *Handler) GetRandomVerse(c echo.Context) error {
	work := c.QueryParam("work")
	book := strings.ToUpper(c.QueryParam("book"))
	testament := strings.ToUpper(c.QueryParam("testament"))

	// Validate testament if provided
	if testament != "" && testament != "OT" && testament != "NT" {
		return echo.NewHTTPError(http.StatusBadRequest, "testament must be OT or NT")
	}

	v, err := h.verses.GetRandomVerse(work, book, testament)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "no verses found matching criteria")
	}

	return c.JSON(http.StatusOK, v)
}
