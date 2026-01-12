package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mbarlow/word/internal/model"
	"github.com/mbarlow/word/internal/pipeline"
	"github.com/mbarlow/word/internal/pipeline/ingest"
	"github.com/mbarlow/word/internal/pipeline/render"
	"github.com/mbarlow/word/internal/pipeline/store"
	"github.com/mbarlow/word/internal/pipeline/validate"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: pipeline <command> [args]")
		fmt.Println("Commands:")
		fmt.Println("  ingest <work> <book>  - Ingest a book (work: kjv, heb-wlc, grc-tr1894)")
		fmt.Println("  render                - Render all JSONL to Markdown")
		fmt.Println("  validate              - Validate rendered output")
		fmt.Println("  genesis1              - Run full pipeline for Genesis 1 (all works)")
		fmt.Println("  load                  - Load sample data into SQLite")
		os.Exit(1)
	}

	cmd := os.Args[1]

	switch cmd {
	case "ingest":
		if len(os.Args) < 4 {
			fmt.Println("Usage: pipeline ingest <work> <book>")
			os.Exit(1)
		}
		runIngest(os.Args[2], os.Args[3])

	case "render":
		runRender()

	case "validate":
		runValidate()

	case "genesis1":
		runGenesis1Pipeline()

	case "load":
		runLoad()

	default:
		fmt.Printf("Unknown command: %s\n", cmd)
		os.Exit(1)
	}
}

func runIngest(work, book string) {
	if _, err := pipeline.LoadCanon("data/canon/books.json"); err != nil {
		fmt.Printf("Failed to load canon: %v\n", err)
		os.Exit(1)
	}

	var verses []pipeline.Verse
	var err error

	switch work {
	case "kjv":
		parser := ingest.NewKJVParser("data/source/kjv_raw")
		verses, err = parser.ParseBook(book)

	case "heb-wlc":
		parser := ingest.NewWLCParser("data/source/wlc_raw/morphhb/wlc")
		verses, err = parser.ParseBook(book)

	case "grc-tr1894":
		parser := ingest.NewTR1894Parser("data/source/tr1894_raw/greektext-scrivener/textonly")
		verses, err = parser.ParseBook(book)

	default:
		fmt.Printf("Unknown work: %s\n", work)
		os.Exit(1)
	}

	if err != nil {
		fmt.Printf("Failed to parse: %v\n", err)
		os.Exit(1)
	}

	for _, v := range verses {
		data, _ := json.Marshal(v)
		fmt.Println(string(data))
	}

	fmt.Fprintf(os.Stderr, "Ingested %d verses from %s %s\n", len(verses), work, book)
}

func runRender() {
	fmt.Println("Render not yet implemented - use genesis1 for sample")
}

func runValidate() {
	fmt.Println("Validate not yet implemented - use genesis1 for sample")
}

func runLoad() {
	fmt.Println("=== Loading Sample Data into SQLite ===")
	fmt.Println()

	// Load canon
	canon, err := pipeline.LoadCanon("data/canon/books.json")
	if err != nil {
		fmt.Printf("Failed to load canon: %v\n", err)
		os.Exit(1)
	}

	// Open database
	fmt.Print("Opening database... ")
	db, err := store.NewSQLiteStore("data/word.db")
	if err != nil {
		fmt.Printf("FAILED: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("OK")

	// Load books
	fmt.Print("Loading canonical books... ")
	if err := db.LoadBooksFromCanon(canon); err != nil {
		fmt.Printf("FAILED: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("OK (66 books)")

	// Load works
	fmt.Print("Loading works metadata... ")
	works := []model.Work{
		{ID: "kjv", Name: "King James Version", Lang: "en", Description: "1769 standardized text", License: "Public Domain", SourceURL: "https://ebible.org/find/show.php?id=eng-kjv2006"},
		{ID: "heb-wlc", Name: "Westminster Leningrad Codex", Lang: "he", Description: "Masoretic Hebrew text", License: "Public Domain", SourceURL: "https://github.com/openscriptures/morphhb"},
		{ID: "grc-tr1894", Name: "Scrivener 1894 Textus Receptus", Lang: "grc", Description: "Greek NT underlying KJV", License: "Public Domain", SourceURL: "https://github.com/byztxt/greektext-scrivener"},
	}
	for _, w := range works {
		if err := db.LoadWork(w); err != nil {
			fmt.Printf("FAILED: %v\n", err)
			os.Exit(1)
		}
	}
	fmt.Println("OK (3 works)")

	// Ingest sample chapters
	var allVerses []pipeline.Verse

	// KJV Genesis 1
	fmt.Print("Ingesting KJV Genesis 1... ")
	kjvParser := ingest.NewKJVParser("data/source/kjv_raw")
	kjvVerses, err := kjvParser.ParseBook("GEN")
	if err != nil {
		fmt.Printf("FAILED: %v\n", err)
	} else {
		for _, v := range kjvVerses {
			if v.Chapter == 1 {
				allVerses = append(allVerses, v)
			}
		}
		fmt.Printf("OK (%d verses)\n", countByWork(allVerses, "kjv"))
	}

	// WLC Genesis 1
	fmt.Print("Ingesting WLC Genesis 1... ")
	wlcParser := ingest.NewWLCParser("data/source/wlc_raw/morphhb/wlc")
	wlcVerses, err := wlcParser.ParseBook("GEN")
	if err != nil {
		fmt.Printf("FAILED: %v\n", err)
	} else {
		for _, v := range wlcVerses {
			if v.Chapter == 1 {
				allVerses = append(allVerses, v)
			}
		}
		fmt.Printf("OK (%d verses)\n", countByWork(allVerses, "heb-wlc"))
	}

	// TR1894 Matthew 1
	fmt.Print("Ingesting TR1894 Matthew 1... ")
	trParser := ingest.NewTR1894Parser("data/source/tr1894_raw/greektext-scrivener/textonly")
	trVerses, err := trParser.ParseBook("MAT")
	if err != nil {
		fmt.Printf("FAILED: %v\n", err)
	} else {
		for _, v := range trVerses {
			if v.Chapter == 1 {
				allVerses = append(allVerses, v)
			}
		}
		fmt.Printf("OK (%d verses)\n", countByWork(allVerses, "grc-tr1894"))
	}

	// Load verses into database
	fmt.Print("Loading verses into database... ")
	if err := db.LoadVerses(allVerses); err != nil {
		fmt.Printf("FAILED: %v\n", err)
		os.Exit(1)
	}
	count, _ := db.GetVerseCount()
	fmt.Printf("OK (%d total verses)\n", count)

	fmt.Println()
	fmt.Println("Database ready at data/word.db")
	fmt.Println("Run the API server with: go run ./cmd/server")
}

func countByWork(verses []pipeline.Verse, work string) int {
	count := 0
	for _, v := range verses {
		if v.Work == work {
			count++
		}
	}
	return count
}

func runGenesis1Pipeline() {
	fmt.Println("=== Word Pipeline: Genesis 1 ===")
	fmt.Println()

	canon, err := pipeline.LoadCanon("data/canon/books.json")
	if err != nil {
		fmt.Printf("Failed to load canon: %v\n", err)
		os.Exit(1)
	}

	var allVerses []pipeline.Verse
	validator := validate.NewValidator()
	renderer := render.NewMarkdownRenderer("texts")

	validator.ExpectedCounts["kjv/GEN/1"] = 31
	validator.ExpectedCounts["heb-wlc/GEN/1"] = 31
	validator.ExpectedCounts["grc-tr1894/MAT/1"] = 25

	// Ingest KJV Genesis
	fmt.Print("Ingesting KJV Genesis... ")
	kjvParser := ingest.NewKJVParser("data/source/kjv_raw")
	kjvVerses, err := kjvParser.ParseBook("GEN")
	if err != nil {
		fmt.Printf("FAILED: %v\n", err)
	} else {
		var ch1 []pipeline.Verse
		for _, v := range kjvVerses {
			if v.Chapter == 1 {
				ch1 = append(ch1, v)
			}
		}
		allVerses = append(allVerses, ch1...)
		fmt.Printf("OK (%d verses in chapter 1)\n", len(ch1))
		report := validator.ValidateChapter(ch1)
		fmt.Printf("  %s\n", validate.FormatReport(report))
	}

	// Ingest WLC Genesis
	fmt.Print("Ingesting WLC Genesis... ")
	wlcParser := ingest.NewWLCParser("data/source/wlc_raw/morphhb/wlc")
	wlcVerses, err := wlcParser.ParseBook("GEN")
	if err != nil {
		fmt.Printf("FAILED: %v\n", err)
	} else {
		var ch1 []pipeline.Verse
		for _, v := range wlcVerses {
			if v.Chapter == 1 {
				ch1 = append(ch1, v)
			}
		}
		allVerses = append(allVerses, ch1...)
		fmt.Printf("OK (%d verses in chapter 1)\n", len(ch1))
		report := validator.ValidateChapter(ch1)
		fmt.Printf("  %s\n", validate.FormatReport(report))
	}

	// Ingest TR1894 Matthew
	fmt.Print("Ingesting TR1894 Matthew... ")
	trParser := ingest.NewTR1894Parser("data/source/tr1894_raw/greektext-scrivener/textonly")
	trVerses, err := trParser.ParseBook("MAT")
	if err != nil {
		fmt.Printf("FAILED: %v\n", err)
	} else {
		var ch1 []pipeline.Verse
		for _, v := range trVerses {
			if v.Chapter == 1 {
				ch1 = append(ch1, v)
			}
		}
		allVerses = append(allVerses, ch1...)
		fmt.Printf("OK (%d verses in chapter 1)\n", len(ch1))
		report := validator.ValidateChapter(ch1)
		fmt.Printf("  %s\n", validate.FormatReport(report))
	}

	// Render to Markdown
	fmt.Println()
	fmt.Print("Rendering to Markdown... ")
	if err := renderer.RenderAll(allVerses); err != nil {
		fmt.Printf("FAILED: %v\n", err)
	} else {
		fmt.Println("OK")
	}

	// Load into SQLite
	fmt.Println()
	fmt.Print("Opening database... ")
	db, err := store.NewSQLiteStore("data/word.db")
	if err != nil {
		fmt.Printf("FAILED: %v\n", err)
	} else {
		fmt.Println("OK")

		fmt.Print("Loading canonical books... ")
		if err := db.LoadBooksFromCanon(canon); err != nil {
			fmt.Printf("FAILED: %v\n", err)
		} else {
			fmt.Println("OK")
		}

		fmt.Print("Loading works... ")
		works := []model.Work{
			{ID: "kjv", Name: "King James Version", Lang: "en", Description: "1769 standardized text", License: "Public Domain", SourceURL: "https://ebible.org/find/show.php?id=eng-kjv2006"},
			{ID: "heb-wlc", Name: "Westminster Leningrad Codex", Lang: "he", Description: "Masoretic Hebrew text", License: "Public Domain", SourceURL: "https://github.com/openscriptures/morphhb"},
			{ID: "grc-tr1894", Name: "Scrivener 1894 Textus Receptus", Lang: "grc", Description: "Greek NT underlying KJV", License: "Public Domain", SourceURL: "https://github.com/byztxt/greektext-scrivener"},
		}
		for _, w := range works {
			db.LoadWork(w)
		}
		fmt.Println("OK")

		fmt.Print("Loading verses... ")
		if err := db.LoadVerses(allVerses); err != nil {
			fmt.Printf("FAILED: %v\n", err)
		} else {
			count, _ := db.GetVerseCount()
			fmt.Printf("OK (%d verses)\n", count)
		}
	}

	// Show output files
	fmt.Println()
	fmt.Println("Generated files:")
	files, _ := filepath.Glob("texts/*/*/*")
	for _, f := range files {
		fmt.Printf("  %s\n", f)
	}

	fmt.Println()
	fmt.Println("=== Sample: KJV Genesis 1:1 ===")
	for _, v := range allVerses {
		if v.Work == "kjv" && v.Chapter == 1 && v.Verse == 1 {
			fmt.Println(v.Text)
			break
		}
	}

	fmt.Println()
	fmt.Println("=== Sample: WLC Genesis 1:1 (Hebrew) ===")
	for _, v := range allVerses {
		if v.Work == "heb-wlc" && v.Chapter == 1 && v.Verse == 1 {
			fmt.Println(v.Text)
			break
		}
	}

	fmt.Println()
	fmt.Println("=== Sample: TR1894 Matthew 1:1 (Greek) ===")
	for _, v := range allVerses {
		if v.Work == "grc-tr1894" && v.Chapter == 1 && v.Verse == 1 {
			fmt.Println(v.Text)
			break
		}
	}

	fmt.Println()
	fmt.Println("Pipeline complete! Database ready at data/word.db")
}
