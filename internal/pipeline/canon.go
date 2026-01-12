package pipeline

import (
	"encoding/json"
	"os"
	"strings"
)

var booksByOSIS map[string]BookInfo
var booksBySlug map[string]BookInfo

// LoadCanon loads the canonical book list from books.json.
func LoadCanon(path string) (*Canon, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var canon Canon
	if err := json.Unmarshal(data, &canon); err != nil {
		return nil, err
	}

	// Build lookup maps
	booksByOSIS = make(map[string]BookInfo)
	booksBySlug = make(map[string]BookInfo)
	for _, b := range canon.Books {
		booksByOSIS[b.OSIS] = b
		booksBySlug[b.Slug] = b
	}

	return &canon, nil
}

// GetBookByOSIS returns book info by OSIS code.
func GetBookByOSIS(osis string) (BookInfo, bool) {
	b, ok := booksByOSIS[strings.ToUpper(osis)]
	return b, ok
}

// GetBookBySlug returns book info by slug.
func GetBookBySlug(slug string) (BookInfo, bool) {
	b, ok := booksBySlug[strings.ToLower(slug)]
	return b, ok
}

// OSISToSlug converts OSIS code to slug.
func OSISToSlug(osis string) string {
	if b, ok := GetBookByOSIS(osis); ok {
		return b.Slug
	}
	return strings.ToLower(osis)
}
