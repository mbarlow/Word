package pipeline

// Verse represents a verse in the canonical intermediate format.
type Verse struct {
	ID       string `json:"id"`        // VID: work/osis/chapter/verse
	Work     string `json:"work"`      // e.g., "kjv", "heb-wlc", "grc-tr1894"
	OSIS     string `json:"osis"`      // 3-letter book code
	BookSlug string `json:"book_slug"` // e.g., "genesis"
	Chapter  int    `json:"chapter"`
	Verse    int    `json:"verse"`
	Lang     string `json:"lang"` // en, he, grc
	Text     string `json:"text"`
	Source   Source `json:"source"`
}

// Source tracks provenance of the verse.
type Source struct {
	Upstream string `json:"upstream"` // e.g., "ebible-kjv"
	Artifact string `json:"artifact"` // filename
}

// BookInfo contains metadata about a canonical book.
type BookInfo struct {
	OSIS      string `json:"osis"`
	Slug      string `json:"slug"`
	Name      string `json:"name"`
	Testament string `json:"testament"`
	Order     int    `json:"order"`
	Chapters  int    `json:"chapters"`
}

// Canon holds the book list.
type Canon struct {
	Books []BookInfo `json:"books"`
}
