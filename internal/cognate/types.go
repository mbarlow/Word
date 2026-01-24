package cognate

// WordAnalysis represents a single word with morphological analysis.
type WordAnalysis struct {
	ID       string `json:"id"`        // OSHB word ID
	Text     string `json:"text"`      // Hebrew text
	Lemma    string `json:"lemma"`     // Strong's number (raw from source)
	BaseLemma string `json:"base_lemma"` // Strong's number (normalized, prefixes stripped)
	Morph    string `json:"morph"`     // Full morphology code
	POS      string `json:"pos"`       // Part of speech (N=noun, V=verb, etc.)
	Position int    `json:"position"`  // Position in verse (0-indexed)
}

// CognatePair represents two words that share the same root/lemma.
type CognatePair struct {
	Root    string       `json:"root"`    // Base Strong's number linking them
	Word1   WordAnalysis `json:"word1"`   // First word (often verb)
	Word2   WordAnalysis `json:"word2"`   // Second word (often noun)
	Pattern string       `json:"pattern"` // cognate_accusative, repetition, verb_noun, etc.
}

// VerseAnalysis contains the morphological analysis of a verse.
type VerseAnalysis struct {
	VID      string         `json:"vid"`
	OSIS     string         `json:"osis"`
	Chapter  int            `json:"chapter"`
	Verse    int            `json:"verse"`
	Words    []WordAnalysis `json:"words"`
	Cognates []CognatePair  `json:"cognates"`
}

// ChapterAnalysis contains the cognate analysis for an entire chapter.
type ChapterAnalysis struct {
	Work    string          `json:"work"`
	OSIS    string          `json:"osis"`
	Chapter int             `json:"chapter"`
	Verses  []VerseAnalysis `json:"verses"`
	Summary CognateSummary  `json:"summary"`
}

// CognateSummary provides statistics about cognates in a chapter.
type CognateSummary struct {
	TotalVerses   int            `json:"total_verses"`
	VersesWithCognates int       `json:"verses_with_cognates"`
	TotalCognates int            `json:"total_cognates"`
	PatternCounts map[string]int `json:"pattern_counts"`
}
