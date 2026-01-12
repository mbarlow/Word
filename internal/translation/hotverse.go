package translation

import (
	"strings"
)

// HotVerseRanges defines verses requiring extra care and notes.
// Format: "BOOK/CHAPTER" or "BOOK/CHAPTER/VERSE" or "BOOK/CHAPTER/START-END"
var HotVerseRanges = []string{
	// Genesis
	"GEN/1/1", "GEN/1/2", "GEN/1/3",
	// Exodus
	"EXO/3/14", "EXO/3/15",
	// Deuteronomy
	"DEU/6/4",
	// Isaiah
	"ISA/7/14", "ISA/9/6",
	"ISA/53/1", "ISA/53/2", "ISA/53/3", "ISA/53/4", "ISA/53/5",
	"ISA/53/6", "ISA/53/7", "ISA/53/8", "ISA/53/9", "ISA/53/10",
	"ISA/53/11", "ISA/53/12",
	// Psalm 22
	"PSA/22/1", "PSA/22/2", "PSA/22/3", "PSA/22/4", "PSA/22/5",
	"PSA/22/6", "PSA/22/7", "PSA/22/8", "PSA/22/9", "PSA/22/10",
	"PSA/22/11", "PSA/22/12", "PSA/22/13", "PSA/22/14", "PSA/22/15",
	"PSA/22/16", "PSA/22/17", "PSA/22/18", "PSA/22/19", "PSA/22/20",
	"PSA/22/21", "PSA/22/22", "PSA/22/23", "PSA/22/24", "PSA/22/25",
	"PSA/22/26", "PSA/22/27", "PSA/22/28", "PSA/22/29", "PSA/22/30",
	"PSA/22/31",
	// Daniel 7
	"DAN/7/1", "DAN/7/2", "DAN/7/3", "DAN/7/4", "DAN/7/5",
	"DAN/7/6", "DAN/7/7", "DAN/7/8", "DAN/7/9", "DAN/7/10",
	"DAN/7/11", "DAN/7/12", "DAN/7/13", "DAN/7/14", "DAN/7/15",
	"DAN/7/16", "DAN/7/17", "DAN/7/18", "DAN/7/19", "DAN/7/20",
	"DAN/7/21", "DAN/7/22", "DAN/7/23", "DAN/7/24", "DAN/7/25",
	"DAN/7/26", "DAN/7/27", "DAN/7/28",
	// Matthew
	"MAT/28/19",
	// John
	"JHN/1/1", "JHN/1/2", "JHN/1/3", "JHN/1/4", "JHN/1/5",
	"JHN/1/6", "JHN/1/7", "JHN/1/8", "JHN/1/9", "JHN/1/10",
	"JHN/1/11", "JHN/1/12", "JHN/1/13", "JHN/1/14", "JHN/1/15",
	"JHN/1/16", "JHN/1/17", "JHN/1/18",
	"JHN/3/16",
	"JHN/8/58",
	// Romans 3-5
	"ROM/3", "ROM/4", "ROM/5",
	// Philippians
	"PHP/2/5", "PHP/2/6", "PHP/2/7", "PHP/2/8", "PHP/2/9",
	"PHP/2/10", "PHP/2/11",
	// Colossians
	"COL/1/15", "COL/1/16", "COL/1/17", "COL/1/18", "COL/1/19", "COL/1/20",
	// Hebrews 1
	"HEB/1",
	// Revelation 1
	"REV/1",
}

// hotVerseSet is a map for O(1) lookup.
var hotVerseSet map[string]bool

func init() {
	hotVerseSet = make(map[string]bool)
	for _, v := range HotVerseRanges {
		hotVerseSet[v] = true
	}
}

// IsHotVerse checks if a verse ID requires extra care.
// vid format: "work/BOOK/CHAPTER/VERSE" - we strip the work prefix.
func IsHotVerse(vid string) bool {
	// Strip work prefix: "kjv/GEN/1/1" -> "GEN/1/1"
	parts := strings.SplitN(vid, "/", 2)
	if len(parts) < 2 {
		return false
	}
	ref := parts[1] // "GEN/1/1"

	// Check exact match
	if hotVerseSet[ref] {
		return true
	}

	// Check chapter-level match (e.g., "ROM/3" matches "ROM/3/23")
	refParts := strings.Split(ref, "/")
	if len(refParts) >= 2 {
		chapterRef := refParts[0] + "/" + refParts[1]
		if hotVerseSet[chapterRef] {
			return true
		}
	}

	return false
}

// GetHotVerseWarning returns a warning message if the verse is hot.
func GetHotVerseWarning(vid string) string {
	if IsHotVerse(vid) {
		return "This is a theologically significant verse requiring extra care and conservative translation choices."
	}
	return ""
}
