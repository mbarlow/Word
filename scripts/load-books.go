//go:build ignore

package main

import (
	"fmt"
	"os"

	"github.com/mbarlow/word/internal/pipeline/ingest"
	"github.com/mbarlow/word/internal/pipeline/store"
)

func main() {
	books := os.Args[1:]
	if len(books) == 0 {
		fmt.Println("Usage: go run scripts/load-books.go JER ECC SNG ...")
		os.Exit(1)
	}

	db, err := store.NewSQLiteStore("data/word.db")
	if err != nil {
		fmt.Printf("Failed to open DB: %v\n", err)
		os.Exit(1)
	}

	kjvParser := ingest.NewKJVParser("data/source/kjv_raw")
	wlcParser := ingest.NewWLCParser("data/source/wlc_raw/morphhb/wlc")

	for _, book := range books {
		fmt.Printf("Loading KJV %s... ", book)
		verses, err := kjvParser.ParseBook(book)
		if err != nil {
			fmt.Printf("FAILED: %v\n", err)
			continue
		}
		if err := db.LoadVerses(verses); err != nil {
			fmt.Printf("DB FAILED: %v\n", err)
			continue
		}
		fmt.Printf("OK (%d verses)\n", len(verses))

		fmt.Printf("Loading WLC %s... ", book)
		hVerses, err := wlcParser.ParseBook(book)
		if err != nil {
			fmt.Printf("FAILED: %v (skipping Hebrew)\n", err)
			continue
		}
		if err := db.LoadVerses(hVerses); err != nil {
			fmt.Printf("DB FAILED: %v\n", err)
			continue
		}
		fmt.Printf("OK (%d verses)\n", len(hVerses))
	}

	count, _ := db.GetVerseCount()
	fmt.Printf("\nTotal verses in DB: %d\n", count)
}
