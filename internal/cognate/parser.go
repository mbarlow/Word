package cognate

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// OSIS XML structures for parsing OSHB morphology data
type osisDoc struct {
	XMLName  xml.Name `xml:"osis"`
	OsisText osisText `xml:"osisText"`
}

type osisText struct {
	Divs []bookDiv `xml:"div"`
}

type bookDiv struct {
	Type     string       `xml:"type,attr"`
	OsisID   string       `xml:"osisID,attr"`
	Chapters []chapterDiv `xml:"chapter"`
}

type chapterDiv struct {
	OsisID string     `xml:"osisID,attr"`
	Verses []verseRaw `xml:"verse"`
}

type verseRaw struct {
	OsisID  string `xml:"osisID,attr"`
	Content string `xml:",innerxml"`
}

// wordElem represents a <w> element with morphology
type wordElem struct {
	ID    string // word ID
	Lemma string // Strong's number(s)
	Morph string // morphology code
	Text  string // Hebrew text
}

// OSHBParser parses OSHB XML files with morphology data.
type OSHBParser struct {
	sourceDir string
}

// NewOSHBParser creates a new OSHB parser.
func NewOSHBParser(sourceDir string) *OSHBParser {
	return &OSHBParser{sourceDir: sourceDir}
}

// ParseChapter parses a single chapter and returns verse analyses with morphology.
func (p *OSHBParser) ParseChapter(osis string, chapter int) ([]VerseAnalysis, error) {
	filename := p.osisToFilename(osis)
	path := filepath.Join(p.sourceDir, filename)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("no OSHB file found for %s at %s", osis, path)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var doc osisDoc
	if err := xml.Unmarshal(data, &doc); err != nil {
		return nil, fmt.Errorf("failed to parse XML: %w", err)
	}

	var verses []VerseAnalysis

	for _, div := range doc.OsisText.Divs {
		if div.Type != "book" {
			continue
		}

		for _, ch := range div.Chapters {
			chapterNum := extractChapterNum(ch.OsisID)
			if chapterNum != chapter {
				continue
			}

			for _, v := range ch.Verses {
				verseNum := extractVerseNum(v.OsisID)
				words := p.parseWords(v.Content)

				va := VerseAnalysis{
					VID:     fmt.Sprintf("heb-wlc/%s/%d/%d", osis, chapterNum, verseNum),
					OSIS:    osis,
					Chapter: chapterNum,
					Verse:   verseNum,
					Words:   words,
				}
				verses = append(verses, va)
			}
		}
	}

	return verses, nil
}

func (p *OSHBParser) osisToFilename(osis string) string {
	mapping := map[string]string{
		"GEN": "Gen.xml", "EXO": "Exod.xml", "LEV": "Lev.xml", "NUM": "Num.xml",
		"DEU": "Deut.xml", "JOS": "Josh.xml", "JDG": "Judg.xml", "RUT": "Ruth.xml",
		"1SA": "1Sam.xml", "2SA": "2Sam.xml", "1KI": "1Kgs.xml", "2KI": "2Kgs.xml",
		"1CH": "1Chr.xml", "2CH": "2Chr.xml", "EZR": "Ezra.xml", "NEH": "Neh.xml",
		"EST": "Esth.xml", "JOB": "Job.xml", "PSA": "Ps.xml", "PRO": "Prov.xml",
		"ECC": "Eccl.xml", "SNG": "Song.xml", "ISA": "Isa.xml", "JER": "Jer.xml",
		"LAM": "Lam.xml", "EZK": "Ezek.xml", "DAN": "Dan.xml", "HOS": "Hos.xml",
		"JOL": "Joel.xml", "AMO": "Amos.xml", "OBA": "Obad.xml", "JON": "Jonah.xml",
		"MIC": "Mic.xml", "NAM": "Nah.xml", "HAB": "Hab.xml", "ZEP": "Zeph.xml",
		"HAG": "Hag.xml", "ZEC": "Zech.xml", "MAL": "Mal.xml",
	}
	if fname, ok := mapping[strings.ToUpper(osis)]; ok {
		return fname
	}
	return osis + ".xml"
}

// parseWords extracts word elements with morphology from verse content.
func (p *OSHBParser) parseWords(content string) []WordAnalysis {
	var words []WordAnalysis

	// Regex to match <w> elements with attributes
	wordRe := regexp.MustCompile(`<w\s+([^>]*)>([^<]*)</w>`)
	attrRe := regexp.MustCompile(`(\w+)="([^"]*)"`)

	matches := wordRe.FindAllStringSubmatch(content, -1)
	for i, m := range matches {
		if len(m) < 3 {
			continue
		}

		attrs := m[1]
		text := m[2]

		word := WordAnalysis{
			Text:     text,
			Position: i,
		}

		// Parse attributes
		attrMatches := attrRe.FindAllStringSubmatch(attrs, -1)
		for _, am := range attrMatches {
			if len(am) < 3 {
				continue
			}
			key, value := am[1], am[2]
			switch key {
			case "id":
				word.ID = value
			case "lemma":
				word.Lemma = value
				word.BaseLemma = NormalizeLemma(value)
			case "morph":
				word.Morph = value
				word.POS = extractPOS(value)
			}
		}

		words = append(words, word)
	}

	return words
}

// extractPOS extracts the part of speech from a morphology code.
// Hebrew morph codes start with "H" followed by POS letter:
// N=noun, V=verb, A=adjective, R=preposition, C=conjunction, etc.
func extractPOS(morph string) string {
	if len(morph) < 2 {
		return ""
	}

	// Strip leading "H" for Hebrew
	code := morph
	if strings.HasPrefix(morph, "H") {
		code = morph[1:]
	}

	// Handle compound morphology (e.g., "HR/Ncfsa" for preposition + noun)
	// Return the first POS component
	parts := strings.Split(code, "/")
	if len(parts) > 0 && len(parts[0]) > 0 {
		return string(parts[0][0])
	}

	return ""
}

func extractChapterNum(osisID string) int {
	// Format: "Gen.1"
	parts := strings.Split(osisID, ".")
	if len(parts) >= 2 {
		n, _ := strconv.Atoi(parts[1])
		return n
	}
	return 0
}

func extractVerseNum(osisID string) int {
	// Format: "Gen.1.1"
	parts := strings.Split(osisID, ".")
	if len(parts) >= 3 {
		n, _ := strconv.Atoi(parts[2])
		return n
	}
	return 0
}
