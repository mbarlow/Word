package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/mbarlow/word/internal/repository"
	"gorm.io/gorm"
)

type Handler struct {
	verses *repository.VerseRepo
}

func New(db *gorm.DB) *Handler {
	return &Handler{
		verses: repository.NewVerseRepo(db),
	}
}

func (h *Handler) ListWorks(c echo.Context) error {
	works, err := h.verses.ListWorks()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, works)
}

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

func (h *Handler) ListBooks(c echo.Context) error {
	books, err := h.verses.ListBooks()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, books)
}
