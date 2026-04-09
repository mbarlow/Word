package meaning

// MeaningGraph holds the semantic analysis for a single chapter.
type MeaningGraph struct {
	SchemaVersion string         `json:"schemaVersion"`
	Work          string         `json:"work"`
	OSIS          string         `json:"osis"`
	Book          string         `json:"book"`
	Chapter       int            `json:"chapter"`
	Verses        []VerseMeaning `json:"verses"`
	Repetitions   []Repetition   `json:"repetitions"`
}

// VerseMeaning holds token-level meaning data for a single verse.
type VerseMeaning struct {
	VerseID string         `json:"verseId"`
	Chapter int            `json:"chapter"`
	Verse   int            `json:"verse"`
	Tokens  []MeaningToken `json:"tokens"`
}

// MeaningToken represents a single Hebrew word with its semantic data.
type MeaningToken struct {
	Position int      `json:"position"`
	Hebrew   string   `json:"hebrew"`
	Lemma    string   `json:"lemma"`
	Root     string   `json:"root"`
	POS      string   `json:"pos"`
	Morph    string   `json:"morph"`
	Tags     []string `json:"tags,omitempty"`
	Weight   float64  `json:"weight"`
	RepCount int      `json:"repCount"`
}

// Repetition tracks a Hebrew root that appears multiple times across a book.
type Repetition struct {
	Root   string   `json:"root"`
	Lemma  string   `json:"lemma"`
	Gloss  string   `json:"gloss"`
	Count  int      `json:"count"`
	Verses []string `json:"verses"`
}

// BookSummary holds book-level repetition and semantic statistics.
type BookSummary struct {
	SchemaVersion string       `json:"schemaVersion"`
	Work          string       `json:"work"`
	OSIS          string       `json:"osis"`
	Book          string       `json:"book"`
	TotalTokens   int          `json:"totalTokens"`
	UniqueRoots   int          `json:"uniqueRoots"`
	Chapters      int          `json:"chapters"`
	Repetitions   []Repetition `json:"repetitions"`
}
