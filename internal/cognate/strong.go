package cognate

import (
	"regexp"
	"strings"
)

var (
	// prefixRe matches Hebrew prefixes like "b/", "c/", "d/", "l/", "m/" before Strong's numbers
	prefixRe = regexp.MustCompile(`^[bcdlmw]/`)

	// suffixRe matches letter suffixes like "a", "b" after Strong's numbers (e.g., "1254 a")
	suffixRe = regexp.MustCompile(`\s*[a-z]$`)

	// numberRe extracts just the numeric portion of a Strong's number
	numberRe = regexp.MustCompile(`\d+`)
)

// NormalizeLemma strips prefixes and suffixes from a Strong's number to get the base lemma.
// Examples:
//   "b/7225" -> "7225"
//   "1254 a" -> "1254"
//   "d/776" -> "776"
//   "c/d/776" -> "776"
func NormalizeLemma(lemma string) string {
	if lemma == "" {
		return ""
	}

	// Strip all prefix patterns (can be multiple, e.g., "c/d/776")
	result := lemma
	for prefixRe.MatchString(result) {
		result = prefixRe.ReplaceAllString(result, "")
	}

	// Strip letter suffixes like " a"
	result = suffixRe.ReplaceAllString(result, "")

	// Clean up whitespace
	result = strings.TrimSpace(result)

	return result
}

// ExtractNumber extracts just the numeric portion of a Strong's number.
// Examples:
//   "b/7225" -> "7225"
//   "1254 a" -> "1254"
//   "H1254" -> "1254"
func ExtractNumber(lemma string) string {
	match := numberRe.FindString(lemma)
	return match
}

// AreCognates checks if two lemmas share the same root.
// Uses normalized base lemmas for comparison.
func AreCognates(lemma1, lemma2 string) bool {
	base1 := NormalizeLemma(lemma1)
	base2 := NormalizeLemma(lemma2)

	if base1 == "" || base2 == "" {
		return false
	}

	// Direct match
	if base1 == base2 {
		return true
	}

	// Check if numbers are adjacent (cognate roots often have sequential Strong's numbers)
	// e.g., 1876 (verb) and 1877 (noun) from the same root
	num1 := ExtractNumber(base1)
	num2 := ExtractNumber(base2)

	if num1 != "" && num2 != "" {
		// Check adjacency (within 1)
		n1 := parseInt(num1)
		n2 := parseInt(num2)
		if n1 > 0 && n2 > 0 && abs(n1-n2) <= 1 {
			return true
		}
	}

	return false
}

// parseInt is a simple string to int converter
func parseInt(s string) int {
	result := 0
	for _, c := range s {
		if c >= '0' && c <= '9' {
			result = result*10 + int(c-'0')
		}
	}
	return result
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
