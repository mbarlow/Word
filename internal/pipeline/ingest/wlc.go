package ingest

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/mbarlow/word/internal/pipeline"
	"golang.org/x/text/unicode/norm"
)

// WLC OSIS XML structures
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
	Verses []verseElem `xml:"verse"`
}

type verseElem struct {
	OsisID  string `xml:"osisID,attr"`
	Content string `xml:",innerxml"`
}

// WLCParser parses Westminster Leningrad Codex OSIS XML files.
type WLCParser struct {
	sourceDir string
}

// NewWLCParser creates a new WLC parser.
func NewWLCParser(sourceDir string) *WLCParser {
	return &WLCParser{sourceDir: sourceDir}
}

// ParseBook parses a single book and returns verses.
func (p *WLCParser) ParseBook(osis string) ([]pipeline.Verse, error) {
	// WLC uses different file naming (Gen.xml, Exod.xml, etc.)
	filename := p.osisToFilename(osis)
	path := filepath.Join(p.sourceDir, filename)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("no WLC file found for %s at %s", osis, path)
	}

	return p.parseFile(path, osis)
}

func (p *WLCParser) osisToFilename(osis string) string {
	// Map standard OSIS codes to WLC filenames
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

func (p *WLCParser) parseFile(path string, osis string) ([]pipeline.Verse, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var doc osisDoc
	if err := xml.Unmarshal(data, &doc); err != nil {
		return nil, fmt.Errorf("failed to parse XML: %w", err)
	}

	var verses []pipeline.Verse
	slug := pipeline.OSISToSlug(osis)

	for _, div := range doc.OsisText.Divs {
		if div.Type != "book" {
			continue
		}

		for _, chapter := range div.Chapters {
			chapterNum := p.extractChapterNum(chapter.OsisID)

			for _, verse := range chapter.Verses {
				verseNum := p.extractVerseNum(verse.OsisID)
				text := p.extractText(verse.Content)

				verses = append(verses, pipeline.Verse{
					ID:       fmt.Sprintf("heb-wlc/%s/%d/%d", osis, chapterNum, verseNum),
					Work:     "heb-wlc",
					OSIS:     osis,
					BookSlug: slug,
					Chapter:  chapterNum,
					Verse:    verseNum,
					Lang:     "he",
					Text:     text,
					Source: pipeline.Source{
						Upstream: "oshb-morphhb",
						Artifact: filepath.Base(path),
					},
				})
			}
		}
	}

	return verses, nil
}

func (p *WLCParser) extractChapterNum(osisID string) int {
	// Format: "Gen.1"
	parts := strings.Split(osisID, ".")
	if len(parts) >= 2 {
		n, _ := strconv.Atoi(parts[1])
		return n
	}
	return 0
}

func (p *WLCParser) extractVerseNum(osisID string) int {
	// Format: "Gen.1.1"
	parts := strings.Split(osisID, ".")
	if len(parts) >= 3 {
		n, _ := strconv.Atoi(parts[2])
		return n
	}
	return 0
}

func (p *WLCParser) extractText(content string) string {
	// Extract Hebrew text from <w> elements
	// The content contains XML like: <w lemma="..." morph="...">בְּ/רֵאשִׁ֖ית</w>

	var words []string
	inWord := false
	var wordBuf strings.Builder

	// Simple state machine to extract text between <w> tags
	i := 0
	for i < len(content) {
		if strings.HasPrefix(content[i:], "<w ") || strings.HasPrefix(content[i:], "<w>") {
			// Skip to >
			for i < len(content) && content[i] != '>' {
				i++
			}
			i++ // skip >
			inWord = true
			wordBuf.Reset()
		} else if strings.HasPrefix(content[i:], "</w>") {
			if inWord && wordBuf.Len() > 0 {
				words = append(words, wordBuf.String())
			}
			inWord = false
			i += 4
		} else if strings.HasPrefix(content[i:], "<seg") {
			// Skip seg elements but include their content if it's end-of-verse marker
			for i < len(content) && content[i] != '>' {
				i++
			}
			i++ // skip >
		} else if strings.HasPrefix(content[i:], "</seg>") {
			i += 6
		} else if strings.HasPrefix(content[i:], "<") {
			// Skip other tags
			for i < len(content) && content[i] != '>' {
				i++
			}
			i++
		} else {
			if inWord {
				wordBuf.WriteByte(content[i])
			}
			i++
		}
	}

	text := strings.Join(words, " ")
	text = norm.NFC.String(text)
	text = strings.TrimSpace(text)

	return text
}
