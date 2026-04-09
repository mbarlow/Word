package meaning

import (
	"fmt"
	"math"
	"sort"

	"github.com/mbarlow/word/internal/cognate"
)

// Builder constructs meaning graphs using the OSHB morphology data.
type Builder struct {
	parser *cognate.OSHBParser
}

// NewBuilder creates a meaning graph builder with access to Hebrew source data.
func NewBuilder(oshbDir string) *Builder {
	return &Builder{
		parser: cognate.NewOSHBParser(oshbDir),
	}
}

// BuildChapter generates a meaning graph for a single chapter.
// bookRepCounts should be pre-computed via CountBookRepetitions for accurate rep counts.
func (b *Builder) BuildChapter(osis string, chapter int, bookRepCounts map[string]int) (*MeaningGraph, error) {
	verses, err := b.parser.ParseChapter(osis, chapter)
	if err != nil {
		return nil, fmt.Errorf("parse chapter %s.%d: %w", osis, chapter, err)
	}

	graph := &MeaningGraph{
		SchemaVersion: "1.0.0",
		Work:          "heb-wlc",
		OSIS:          osis,
		Chapter:       chapter,
		Verses:        make([]VerseMeaning, 0, len(verses)),
	}

	for _, va := range verses {
		vm := VerseMeaning{
			VerseID: fmt.Sprintf("heb-wlc/%s/%d/%d", osis, va.Chapter, va.Verse),
			Chapter: va.Chapter,
			Verse:   va.Verse,
			Tokens:  make([]MeaningToken, 0, len(va.Words)),
		}

		for _, w := range va.Words {
			root := cognate.NormalizeLemma(w.Lemma)
			repCount := bookRepCounts[root]

			token := MeaningToken{
				Position: w.Position,
				Hebrew:   w.Text,
				Lemma:    w.BaseLemma,
				Root:     root,
				POS:      w.POS,
				Morph:    w.Morph,
				Tags:     inferTags(root, w.POS),
				Weight:   computeWeight(root, repCount, w.POS),
				RepCount: repCount,
			}
			vm.Tokens = append(vm.Tokens, token)
		}

		graph.Verses = append(graph.Verses, vm)
	}

	return graph, nil
}

// CountBookRepetitions counts how many times each Hebrew root appears across all chapters of a book.
func (b *Builder) CountBookRepetitions(osis string, chapters int) (map[string]int, map[string][]string, error) {
	rootCounts := make(map[string]int)
	rootVerses := make(map[string][]string)

	for ch := 1; ch <= chapters; ch++ {
		verses, err := b.parser.ParseChapter(osis, ch)
		if err != nil {
			continue // skip chapters that fail to parse
		}

		for _, va := range verses {
			vid := fmt.Sprintf("heb-wlc/%s/%d/%d", osis, va.Chapter, va.Verse)
			seen := make(map[string]bool) // track unique roots per verse

			for _, w := range va.Words {
				root := cognate.NormalizeLemma(w.Lemma)
				if root == "" {
					continue
				}
				rootCounts[root]++
				if !seen[root] {
					rootVerses[root] = append(rootVerses[root], vid)
					seen[root] = true
				}
			}
		}
	}

	return rootCounts, rootVerses, nil
}

// BuildBookSummary generates a book-level summary of repetitions.
func (b *Builder) BuildBookSummary(osis, bookName string, chapters int) (*BookSummary, error) {
	rootCounts, rootVerses, err := b.CountBookRepetitions(osis, chapters)
	if err != nil {
		return nil, err
	}

	// Build sorted repetition list (only roots appearing 3+ times)
	var reps []Repetition
	totalTokens := 0
	for root, count := range rootCounts {
		totalTokens += count
		if count >= 3 {
			reps = append(reps, Repetition{
				Root:   root,
				Lemma:  root,
				Gloss:  Gloss(root),
				Count:  count,
				Verses: rootVerses[root],
			})
		}
	}

	sort.Slice(reps, func(i, j int) bool {
		return reps[i].Count > reps[j].Count
	})

	return &BookSummary{
		SchemaVersion: "1.0.0",
		Work:          "heb-wlc",
		OSIS:          osis,
		Book:          bookName,
		TotalTokens:   totalTokens,
		UniqueRoots:   len(rootCounts),
		Chapters:      chapters,
		Repetitions:   reps,
	}, nil
}

// BuildRepetitions generates the repetition list for a chapter graph.
func BuildRepetitions(rootCounts map[string]int, rootVerses map[string][]string) []Repetition {
	var reps []Repetition
	for root, count := range rootCounts {
		if count >= 3 {
			reps = append(reps, Repetition{
				Root:   root,
				Lemma:  root,
				Gloss:  Gloss(root),
				Count:  count,
				Verses: rootVerses[root],
			})
		}
	}
	sort.Slice(reps, func(i, j int) bool {
		return reps[i].Count > reps[j].Count
	})
	return reps
}

// inferTags assigns semantic tags based on the root and POS.
func inferTags(root, pos string) []string {
	var tags []string

	// Tag by known semantic families
	semanticFamilies := map[string][]string{
		"1892": {"emptiness", "breath", "temporary", "vanity"},
		"2451": {"wisdom", "skill", "insight"},
		"5999": {"labor", "toil", "suffering"},
		"3504": {"profit", "gain", "advantage"},
		"8121": {"sun", "cosmic", "natural-order"},
		"7307": {"spirit", "wind", "breath"},
		"430":  {"divine", "God"},
		"3068": {"divine", "LORD", "covenant-name"},
		"1285": {"covenant", "agreement", "promise"},
		"2319": {"new", "fresh", "renewal"},
		"8451": {"law", "instruction", "teaching"},
		"4941": {"justice", "judgment", "ordinance"},
		"6666": {"righteousness", "justice"},
		"7725": {"return", "repentance", "restoration"},
		"157":  {"love", "desire"},
		"160":  {"love", "affection"},
		"1730": {"beloved", "love", "desire"},
		"1588": {"garden", "enclosed-space"},
		"3303": {"beautiful", "lovely"},
	}

	if family, ok := semanticFamilies[root]; ok {
		tags = append(tags, family...)
	}

	// Tag by POS
	switch pos {
	case "V":
		tags = append(tags, "action")
	case "N":
		tags = append(tags, "entity")
	case "A":
		tags = append(tags, "quality")
	}

	return tags
}

// computeWeight calculates semantic intensity (0-1) for a token.
// Higher weight = more semantically significant.
func computeWeight(root string, repCount int, pos string) float64 {
	weight := 0.3 // base weight

	// Repetition boosts weight (log scale to avoid domination by ultra-common words)
	if repCount > 1 {
		weight += math.Min(0.4, math.Log2(float64(repCount))*0.08)
	}

	// Known significant roots get a boost
	significantRoots := map[string]bool{
		"1892": true, "2451": true, "5999": true, "3504": true, // Ecclesiastes
		"430": true, "3068": true, "1254": true,                 // Theological
		"1285": true, "8451": true, "4941": true, "6666": true,  // Prophetic
		"157": true, "1730": true, "1588": true,                  // Song of Songs
	}
	if significantRoots[root] {
		weight += 0.2
	}

	// Nouns and verbs are more semantically loaded than particles
	switch pos {
	case "V", "N":
		weight += 0.1
	case "A":
		weight += 0.05
	}

	return math.Min(1.0, weight)
}
