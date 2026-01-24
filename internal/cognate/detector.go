package cognate

// Detector identifies cognate pairs within verses.
type Detector struct {
	parser *OSHBParser
}

// NewDetector creates a new cognate detector.
func NewDetector(sourceDir string) *Detector {
	return &Detector{
		parser: NewOSHBParser(sourceDir),
	}
}

// AnalyzeChapter analyzes a chapter for cognate patterns.
func (d *Detector) AnalyzeChapter(osis string, chapter int) (*ChapterAnalysis, error) {
	verses, err := d.parser.ParseChapter(osis, chapter)
	if err != nil {
		return nil, err
	}

	analysis := &ChapterAnalysis{
		Work:    "heb-wlc",
		OSIS:    osis,
		Chapter: chapter,
		Verses:  make([]VerseAnalysis, 0, len(verses)),
		Summary: CognateSummary{
			TotalVerses:        len(verses),
			VersesWithCognates: 0,
			TotalCognates:      0,
			PatternCounts:      make(map[string]int),
		},
	}

	for _, verse := range verses {
		// Detect cognates in this verse
		cognates := d.detectCognates(verse.Words)
		verse.Cognates = cognates

		if len(cognates) > 0 {
			analysis.Summary.VersesWithCognates++
			analysis.Summary.TotalCognates += len(cognates)
			for _, c := range cognates {
				analysis.Summary.PatternCounts[c.Pattern]++
			}
		}

		analysis.Verses = append(analysis.Verses, verse)
	}

	return analysis, nil
}

// detectCognates finds cognate pairs within a verse's words.
func (d *Detector) detectCognates(words []WordAnalysis) []CognatePair {
	var pairs []CognatePair

	// Group words by base lemma
	lemmaGroups := make(map[string][]WordAnalysis)
	for _, w := range words {
		if w.BaseLemma == "" {
			continue
		}
		lemmaGroups[w.BaseLemma] = append(lemmaGroups[w.BaseLemma], w)
	}

	// Also check for adjacent lemmas (cognate roots often have sequential Strong's numbers)
	// Build a map of all lemma numbers for quick lookup
	lemmaNumbers := make(map[int][]WordAnalysis)
	for _, w := range words {
		if w.BaseLemma == "" {
			continue
		}
		num := parseInt(ExtractNumber(w.BaseLemma))
		if num > 0 {
			lemmaNumbers[num] = append(lemmaNumbers[num], w)
		}
	}

	// Process exact matches (same lemma)
	for lemma, group := range lemmaGroups {
		if len(group) < 2 {
			continue
		}

		// Find pairs within the group
		for i := 0; i < len(group); i++ {
			for j := i + 1; j < len(group); j++ {
				pattern := classifyPattern(group[i], group[j])
				pairs = append(pairs, CognatePair{
					Root:    lemma,
					Word1:   group[i],
					Word2:   group[j],
					Pattern: pattern,
				})
			}
		}
	}

	// Process adjacent lemmas (e.g., 1876 and 1877)
	processed := make(map[string]bool) // track processed pairs
	for num, group1 := range lemmaNumbers {
		// Check for adjacent number
		if group2, exists := lemmaNumbers[num+1]; exists {
			for _, w1 := range group1 {
				for _, w2 := range group2 {
					// Create unique key for this pair
					pairKey := w1.ID + "|" + w2.ID
					if processed[pairKey] {
						continue
					}
					processed[pairKey] = true

					pattern := classifyPattern(w1, w2)
					// Use the lower number as the root identifier
					root := ExtractNumber(w1.BaseLemma) + "/" + ExtractNumber(w2.BaseLemma)
					pairs = append(pairs, CognatePair{
						Root:    root,
						Word1:   w1,
						Word2:   w2,
						Pattern: pattern,
					})
				}
			}
		}
	}

	return pairs
}

// classifyPattern determines the type of cognate relationship.
func classifyPattern(w1, w2 WordAnalysis) string {
	pos1 := w1.POS
	pos2 := w2.POS

	// Same word repeated (same text)
	if w1.Text == w2.Text {
		return "repetition"
	}

	// Verb + Noun = cognate accusative (most common Hebrew pattern)
	if (pos1 == "V" && pos2 == "N") || (pos1 == "N" && pos2 == "V") {
		return "cognate_accusative"
	}

	// Two nouns from same root
	if pos1 == "N" && pos2 == "N" {
		return "noun_pair"
	}

	// Two verbs from same root
	if pos1 == "V" && pos2 == "V" {
		return "verb_pair"
	}

	// Verb/Noun with adjective
	if pos1 == "A" || pos2 == "A" {
		return "adjectival_cognate"
	}

	// Default: related_forms
	return "related_forms"
}

// AnalyzeVerse analyzes a single verse for cognates.
// This is useful for testing individual verses.
func (d *Detector) AnalyzeVerse(osis string, chapter, verse int) (*VerseAnalysis, error) {
	verses, err := d.parser.ParseChapter(osis, chapter)
	if err != nil {
		return nil, err
	}

	for _, v := range verses {
		if v.Verse == verse {
			v.Cognates = d.detectCognates(v.Words)
			return &v, nil
		}
	}

	return nil, nil
}
