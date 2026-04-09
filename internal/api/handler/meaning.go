package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/mbarlow/word/internal/pipeline/meaning"
)

// GetMeaning godoc
// @Summary      Get chapter meaning graph
// @Description  Returns the semantic meaning graph for a chapter — tokens with Hebrew roots, POS, tags, weights, and book-level repetition data.
// @Tags         meaning
// @Produce      json
// @Param        work    path string true "Work code (only heb-wlc supported)"
// @Param        book    path string true "OSIS book code (e.g. ECC, JER)"
// @Param        chapter path int    true "Chapter number"
// @Success      200 {object} meaning.MeaningGraph
// @Failure      404 {object} ErrorResponse
// @Router       /v1/meaning/{work}/{book}/{chapter} [get]
func (h *Handler) GetMeaning(c echo.Context) error {
	work := c.Param("work")
	book := strings.ToUpper(c.Param("book"))
	chapter, err := strconv.Atoi(c.Param("chapter"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid chapter")
	}

	path := filepath.Join("data", "meaning", work, book, fmt.Sprintf("ch-%d.json", chapter))
	data, err := os.ReadFile(path)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("meaning graph not found for %s/%s/%d — run: go run ./cmd/pipeline meaning %s", work, book, chapter, book))
	}

	var graph meaning.MeaningGraph
	if err := json.Unmarshal(data, &graph); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to parse meaning data")
	}

	return c.JSON(http.StatusOK, graph)
}

// GetMeaningSummary godoc
// @Summary      Get book meaning summary
// @Description  Returns book-level repetition statistics — root counts, glosses, verse references for all roots appearing 3+ times.
// @Tags         meaning
// @Produce      json
// @Param        work path string true "Work code (only heb-wlc supported)"
// @Param        book path string true "OSIS book code (e.g. ECC, JER)"
// @Success      200 {object} meaning.BookSummary
// @Failure      404 {object} ErrorResponse
// @Router       /v1/meaning/{work}/{book} [get]
func (h *Handler) GetMeaningSummary(c echo.Context) error {
	work := c.Param("work")
	book := strings.ToUpper(c.Param("book"))

	path := filepath.Join("data", "meaning", work, book, "summary.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("meaning summary not found for %s/%s — run: go run ./cmd/pipeline meaning %s", work, book, book))
	}

	var summary meaning.BookSummary
	if err := json.Unmarshal(data, &summary); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to parse meaning summary")
	}

	return c.JSON(http.StatusOK, summary)
}
