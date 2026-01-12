package ingest

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/mbarlow/word/internal/pipeline"
	"golang.org/x/text/unicode/norm"
)

var (
	// USFM markers
	chapterRe = regexp.MustCompile(`^\\c\s+(\d+)`)
	verseRe   = regexp.MustCompile(`^\\v\s+(\d+)\s+(.*)`)

	// Inline markers to strip
	wordRe     = regexp.MustCompile(`\\w\s+([^|]+)\|[^*]+\\\w\*`)
	addRe      = regexp.MustCompile(`\\add\s*(.*?)\\add\*`)
	footnoteRe = regexp.MustCompile(`\\f\s+.*?\\f\*`)
	ndRe       = regexp.MustCompile(`\\nd\s*(.*?)\\nd\*`)
	plusWRe    = regexp.MustCompile(`\\\+w\s+([^|]+)\|[^*]+\\\+w\*`)

	// Other markers to remove
	markerRe = regexp.MustCompile(`\\[a-z]+\d*\s*`)
)

// KJVParser parses KJV USFM files.
type KJVParser struct {
	sourceDir string
}

// NewKJVParser creates a new KJV parser.
func NewKJVParser(sourceDir string) *KJVParser {
	return &KJVParser{sourceDir: sourceDir}
}

// ParseBook parses a single book and returns verses.
func (p *KJVParser) ParseBook(osis string) ([]pipeline.Verse, error) {
	// Find the USFM file for this book
	pattern := fmt.Sprintf("*-%seng-kjv2006.usfm", osis)
	matches, err := filepath.Glob(filepath.Join(p.sourceDir, pattern))
	if err != nil || len(matches) == 0 {
		return nil, fmt.Errorf("no USFM file found for %s", osis)
	}

	return p.parseFile(matches[0], osis)
}

func (p *KJVParser) parseFile(path string, osis string) ([]pipeline.Verse, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var verses []pipeline.Verse
	var currentChapter int
	var verseBuffer strings.Builder
	var currentVerse int

	slug := pipeline.OSISToSlug(osis)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Check for chapter marker
		if m := chapterRe.FindStringSubmatch(line); m != nil {
			// Save any pending verse
			if currentVerse > 0 && verseBuffer.Len() > 0 {
				verses = append(verses, p.makeVerse(osis, slug, currentChapter, currentVerse, verseBuffer.String(), filepath.Base(path)))
				verseBuffer.Reset()
			}
			currentChapter, _ = strconv.Atoi(m[1])
			currentVerse = 0
			continue
		}

		// Check for verse marker
		if m := verseRe.FindStringSubmatch(line); m != nil {
			// Save previous verse
			if currentVerse > 0 && verseBuffer.Len() > 0 {
				verses = append(verses, p.makeVerse(osis, slug, currentChapter, currentVerse, verseBuffer.String(), filepath.Base(path)))
				verseBuffer.Reset()
			}
			currentVerse, _ = strconv.Atoi(m[1])
			text := p.cleanText(m[2])
			if text != "" {
				verseBuffer.WriteString(text)
			}
			continue
		}

		// Continuation of current verse (lines without markers that contain text)
		if currentVerse > 0 {
			text := p.cleanText(line)
			if text != "" {
				if verseBuffer.Len() > 0 {
					verseBuffer.WriteString(" ")
				}
				verseBuffer.WriteString(text)
			}
		}
	}

	// Save last verse
	if currentVerse > 0 && verseBuffer.Len() > 0 {
		verses = append(verses, p.makeVerse(osis, slug, currentChapter, currentVerse, verseBuffer.String(), filepath.Base(path)))
	}

	return verses, scanner.Err()
}

func (p *KJVParser) cleanText(text string) string {
	// Extract text from \w markers (keep the word, strip Strong's)
	text = wordRe.ReplaceAllString(text, "$1")
	text = plusWRe.ReplaceAllString(text, "$1")

	// Extract text from \add markers
	text = addRe.ReplaceAllString(text, "$1")

	// Extract text from \nd markers (divine name)
	text = ndRe.ReplaceAllString(text, "$1")

	// Remove footnotes
	text = footnoteRe.ReplaceAllString(text, "")

	// Remove remaining markers
	text = markerRe.ReplaceAllString(text, "")

	// Normalize Unicode (NFC)
	text = norm.NFC.String(text)

	// Clean up whitespace
	text = strings.TrimSpace(text)
	text = collapseWhitespace(text)

	return text
}

func (p *KJVParser) makeVerse(osis, slug string, chapter, verse int, text, artifact string) pipeline.Verse {
	return pipeline.Verse{
		ID:       fmt.Sprintf("kjv/%s/%d/%d", osis, chapter, verse),
		Work:     "kjv",
		OSIS:     osis,
		BookSlug: slug,
		Chapter:  chapter,
		Verse:    verse,
		Lang:     "en",
		Text:     text,
		Source: pipeline.Source{
			Upstream: "ebible-kjv2006",
			Artifact: artifact,
		},
	}
}

func collapseWhitespace(s string) string {
	var result strings.Builder
	lastWasSpace := false
	for _, r := range s {
		if unicode.IsSpace(r) {
			if !lastWasSpace {
				result.WriteRune(' ')
				lastWasSpace = true
			}
		} else {
			result.WriteRune(r)
			lastWasSpace = false
		}
	}
	return result.String()
}
