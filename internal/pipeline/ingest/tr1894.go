package ingest

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/mbarlow/word/internal/pipeline"
	"golang.org/x/text/unicode/norm"
)

var (
	// Verse reference pattern: chapter:verse at start of line
	trVerseRe = regexp.MustCompile(`^\s*(\d+):(\d+)\s+(.*)`)
)

// TR1894Parser parses Scrivener TR 1894 plain text files.
type TR1894Parser struct {
	sourceDir string
}

// NewTR1894Parser creates a new TR1894 parser.
func NewTR1894Parser(sourceDir string) *TR1894Parser {
	return &TR1894Parser{sourceDir: sourceDir}
}

// ParseBook parses a single book and returns verses.
func (p *TR1894Parser) ParseBook(osis string) ([]pipeline.Verse, error) {
	// TR files use OSIS-like names with .SCV extension
	filename := p.osisToFilename(osis)
	path := filepath.Join(p.sourceDir, filename)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("no TR1894 file found for %s at %s", osis, path)
	}

	return p.parseFile(path, osis)
}

func (p *TR1894Parser) osisToFilename(osis string) string {
	// TR files use abbreviated names
	mapping := map[string]string{
		"MAT": "MT.SCV", "MRK": "MR.SCV", "LUK": "LU.SCV", "JHN": "JOH.SCV",
		"ACT": "AC.SCV", "ROM": "RO.SCV", "1CO": "1CO.SCV", "2CO": "2CO.SCV",
		"GAL": "GA.SCV", "EPH": "EPH.SCV", "PHP": "PHP.SCV", "COL": "COL.SCV",
		"1TH": "1TH.SCV", "2TH": "2TH.SCV", "1TI": "1TI.SCV", "2TI": "2TI.SCV",
		"TIT": "TIT.SCV", "PHM": "PHM.SCV", "HEB": "HEB.SCV", "JAS": "JAS.SCV",
		"1PE": "1PE.SCV", "2PE": "2PE.SCV", "1JN": "1JO.SCV", "2JN": "2JO.SCV",
		"3JN": "3JO.SCV", "JUD": "JUDE.SCV", "REV": "RE.SCV",
	}
	if fname, ok := mapping[strings.ToUpper(osis)]; ok {
		return fname
	}
	return strings.ToUpper(osis) + ".SCV"
}

func (p *TR1894Parser) parseFile(path string, osis string) ([]pipeline.Verse, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var verses []pipeline.Verse
	var currentChapter, currentVerse int
	var textBuf strings.Builder

	slug := pipeline.OSISToSlug(osis)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Check for verse reference
		if m := trVerseRe.FindStringSubmatch(line); m != nil {
			// Save previous verse
			if currentVerse > 0 && textBuf.Len() > 0 {
				verses = append(verses, p.makeVerse(osis, slug, currentChapter, currentVerse, textBuf.String(), filepath.Base(path)))
				textBuf.Reset()
			}

			currentChapter, _ = strconv.Atoi(m[1])
			currentVerse, _ = strconv.Atoi(m[2])
			text := p.cleanText(m[3])
			if text != "" {
				textBuf.WriteString(text)
			}
		} else {
			// Continuation line
			text := p.cleanText(line)
			if currentVerse > 0 && text != "" {
				if textBuf.Len() > 0 {
					textBuf.WriteString(" ")
				}
				textBuf.WriteString(text)
			}
		}
	}

	// Save last verse
	if currentVerse > 0 && textBuf.Len() > 0 {
		verses = append(verses, p.makeVerse(osis, slug, currentChapter, currentVerse, textBuf.String(), filepath.Base(path)))
	}

	return verses, scanner.Err()
}

func (p *TR1894Parser) cleanText(text string) string {
	// Remove bracketed headers like [EUAGGELION TO KATA MATYAION]
	if strings.HasPrefix(strings.TrimSpace(text), "[") && strings.HasSuffix(strings.TrimSpace(text), "]") {
		return ""
	}

	// Normalize Unicode
	text = norm.NFC.String(text)

	// Clean whitespace
	text = strings.TrimSpace(text)
	text = collapseWhitespace(text)

	return text
}

func (p *TR1894Parser) makeVerse(osis, slug string, chapter, verse int, text, artifact string) pipeline.Verse {
	return pipeline.Verse{
		ID:       fmt.Sprintf("grc-tr1894/%s/%d/%d", osis, chapter, verse),
		Work:     "grc-tr1894",
		OSIS:     osis,
		BookSlug: slug,
		Chapter:  chapter,
		Verse:    verse,
		Lang:     "grc",
		Text:     text,
		Source: pipeline.Source{
			Upstream: "byztxt-scrivener",
			Artifact: artifact,
		},
	}
}
