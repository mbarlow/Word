package validate

import (
	"fmt"
	"strings"

	"github.com/mbarlow/word/internal/pipeline"
)

// Report contains validation results.
type Report struct {
	Work       string
	Book       string
	Chapter    int
	Expected   int
	Actual     int
	Missing    []int
	Duplicates []int
	Errors     []string
}

// Validator validates verse data.
type Validator struct {
	// Expected verse counts per chapter (optional)
	ExpectedCounts map[string]int
}

// NewValidator creates a new validator.
func NewValidator() *Validator {
	return &Validator{
		ExpectedCounts: make(map[string]int),
	}
}

// ValidateChapter validates a chapter's verses.
func (v *Validator) ValidateChapter(verses []pipeline.Verse) Report {
	if len(verses) == 0 {
		return Report{Errors: []string{"no verses provided"}}
	}

	first := verses[0]
	report := Report{
		Work:    first.Work,
		Book:    first.BookSlug,
		Chapter: first.Chapter,
		Actual:  len(verses),
	}

	// Check for duplicates and find verse numbers
	seen := make(map[int]int)
	for _, verse := range verses {
		seen[verse.Verse]++
	}

	for num, count := range seen {
		if count > 1 {
			report.Duplicates = append(report.Duplicates, num)
		}
	}

	// Check for gaps (missing verses)
	maxVerse := 0
	for num := range seen {
		if num > maxVerse {
			maxVerse = num
		}
	}

	for i := 1; i <= maxVerse; i++ {
		if _, ok := seen[i]; !ok {
			report.Missing = append(report.Missing, i)
		}
	}

	// Check expected count if available
	key := fmt.Sprintf("%s/%s/%d", first.Work, first.OSIS, first.Chapter)
	if expected, ok := v.ExpectedCounts[key]; ok {
		report.Expected = expected
		if report.Actual != expected {
			report.Errors = append(report.Errors, fmt.Sprintf("expected %d verses, got %d", expected, report.Actual))
		}
	}

	// Validate individual verses
	for _, verse := range verses {
		if verse.Text == "" {
			report.Errors = append(report.Errors, fmt.Sprintf("verse %d has empty text", verse.Verse))
		}
		if verse.ID == "" {
			report.Errors = append(report.Errors, fmt.Sprintf("verse %d has empty ID", verse.Verse))
		}
	}

	return report
}

// FormatReport formats a validation report as a string.
func FormatReport(r Report) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("%s %s %d: %d verses", r.Work, r.Book, r.Chapter, r.Actual))

	if len(r.Missing) > 0 {
		sb.WriteString(fmt.Sprintf(" | MISSING: %v", r.Missing))
	}
	if len(r.Duplicates) > 0 {
		sb.WriteString(fmt.Sprintf(" | DUPLICATES: %v", r.Duplicates))
	}
	if len(r.Errors) > 0 {
		sb.WriteString(fmt.Sprintf(" | ERRORS: %v", r.Errors))
	}

	if len(r.Missing) == 0 && len(r.Duplicates) == 0 && len(r.Errors) == 0 {
		sb.WriteString(" [OK]")
	}

	return sb.String()
}
