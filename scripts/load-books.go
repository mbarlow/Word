//go:build ignore

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/mbarlow/word/internal/pipeline/ingest"
	"github.com/mbarlow/word/internal/pipeline/store"
)

var ntBooks = map[string]bool{
	"MAT": true, "MRK": true, "LUK": true, "JHN": true,
	"ACT": true, "ROM": true, "1CO": true, "2CO": true,
	"GAL": true, "EPH": true, "PHP": true, "COL": true,
	"1TH": true, "2TH": true, "1TI": true, "2TI": true,
	"TIT": true, "PHM": true, "HEB": true, "JAS": true,
	"1PE": true, "2PE": true, "1JN": true, "2JN": true,
	"3JN": true, "JUD": true, "REV": true,
}

var allOT = []string{
	"GEN", "EXO", "LEV", "NUM", "DEU", "JOS", "JDG", "RUT",
	"1SA", "2SA", "1KI", "2KI", "1CH", "2CH", "EZR", "NEH",
	"EST", "JOB", "PSA", "PRO", "ECC", "SNG", "ISA", "JER",
	"LAM", "EZK", "DAN", "HOS", "JOL", "AMO", "OBA", "JON",
	"MIC", "NAM", "HAB", "ZEP", "HAG", "ZEC", "MAL",
}

var allNT = []string{
	"MAT", "MRK", "LUK", "JHN", "ACT", "ROM", "1CO", "2CO",
	"GAL", "EPH", "PHP", "COL", "1TH", "2TH", "1TI", "2TI",
	"TIT", "PHM", "HEB", "JAS", "1PE", "2PE", "1JN", "2JN",
	"3JN", "JUD", "REV",
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("Usage: go run scripts/load-books.go [--all | --ot | --nt | BOOK ...]")
		fmt.Println("Examples:")
		fmt.Println("  go run scripts/load-books.go --all        # Load everything")
		fmt.Println("  go run scripts/load-books.go --ot         # All OT (KJV + Hebrew)")
		fmt.Println("  go run scripts/load-books.go --nt         # All NT (KJV + Greek)")
		fmt.Println("  go run scripts/load-books.go JER ECC SNG  # Specific books")
		os.Exit(1)
	}

	var books []string
	for _, arg := range args {
		switch strings.ToLower(arg) {
		case "--all":
			books = append(books, allOT...)
			books = append(books, allNT...)
		case "--ot":
			books = append(books, allOT...)
		case "--nt":
			books = append(books, allNT...)
		default:
			books = append(books, strings.ToUpper(arg))
		}
	}

	db, err := store.NewSQLiteStore("data/word.db")
	if err != nil {
		fmt.Printf("Failed to open DB: %v\n", err)
		os.Exit(1)
	}

	kjvParser := ingest.NewKJVParser("data/source/kjv_raw")
	wlcParser := ingest.NewWLCParser("data/source/wlc_raw/morphhb/wlc")
	trParser := ingest.NewTR1894Parser("data/source/tr1894_raw/greektext-scrivener/textonly")

	kjvTotal, hebTotal, grcTotal := 0, 0, 0

	for _, book := range books {
		// KJV — all books
		fmt.Printf("KJV %-4s ", book)
		verses, err := kjvParser.ParseBook(book)
		if err != nil {
			fmt.Printf("SKIP (%v)\n", err)
		} else if err := db.LoadVerses(verses); err != nil {
			fmt.Printf("DB ERR (%v)\n", err)
		} else {
			fmt.Printf("%4d verses  ", len(verses))
			kjvTotal += len(verses)
		}

		// Hebrew WLC — OT only
		if !ntBooks[book] {
			fmt.Printf("WLC ")
			hVerses, err := wlcParser.ParseBook(book)
			if err != nil {
				fmt.Printf("SKIP (%v)", err)
			} else if err := db.LoadVerses(hVerses); err != nil {
				fmt.Printf("DB ERR (%v)", err)
			} else {
				fmt.Printf("%4d verses", len(hVerses))
				hebTotal += len(hVerses)
			}
		}

		// Greek TR 1894 — NT only
		if ntBooks[book] {
			fmt.Printf("TR  ")
			gVerses, err := trParser.ParseBook(book)
			if err != nil {
				fmt.Printf("SKIP (%v)", err)
			} else if err := db.LoadVerses(gVerses); err != nil {
				fmt.Printf("DB ERR (%v)", err)
			} else {
				fmt.Printf("%4d verses", len(gVerses))
				grcTotal += len(gVerses)
			}
		}

		fmt.Println()
	}

	count, _ := db.GetVerseCount()
	fmt.Println()
	fmt.Printf("Loaded:  KJV %d  |  Hebrew %d  |  Greek %d\n", kjvTotal, hebTotal, grcTotal)
	fmt.Printf("Total verses in DB: %d\n", count)
}
