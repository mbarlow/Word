package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mbarlow/word/internal/pipeline"
	"github.com/mbarlow/word/internal/pipeline/ingest"
	"github.com/mbarlow/word/internal/pipeline/render"
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

	default:
		fmt.Printf("Unknown command: %s\n", cmd)
		os.Exit(1)
	}
}

func runIngest(work, book string) {
	// Load canon
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

	// Output as JSONL
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

func runGenesis1Pipeline() {
	fmt.Println("=== Word Pipeline: Genesis 1 ===")
	fmt.Println()

	// Load canon
	if _, err := pipeline.LoadCanon("data/canon/books.json"); err != nil {
		fmt.Printf("Failed to load canon: %v\n", err)
		os.Exit(1)
	}

	var allVerses []pipeline.Verse
	validator := validate.NewValidator()
	renderer := render.NewMarkdownRenderer("texts")

	// Genesis 1 expected verse count
	validator.ExpectedCounts["kjv/GEN/1"] = 31
	validator.ExpectedCounts["heb-wlc/GEN/1"] = 31
	validator.ExpectedCounts["grc-tr1894/MAT/1"] = 25 // Matthew 1 for Greek

	// Ingest KJV Genesis
	fmt.Print("Ingesting KJV Genesis... ")
	kjvParser := ingest.NewKJVParser("data/source/kjv_raw")
	kjvVerses, err := kjvParser.ParseBook("GEN")
	if err != nil {
		fmt.Printf("FAILED: %v\n", err)
	} else {
		// Filter to chapter 1
		var ch1 []pipeline.Verse
		for _, v := range kjvVerses {
			if v.Chapter == 1 {
				ch1 = append(ch1, v)
			}
		}
		allVerses = append(allVerses, ch1...)
		fmt.Printf("OK (%d verses in chapter 1)\n", len(ch1))

		// Validate
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

	// Ingest TR1894 Matthew (Greek NT doesn't have Genesis)
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

	// Show output files
	fmt.Println()
	fmt.Println("Generated files:")
	files, _ := filepath.Glob("texts/*/*/*")
	for _, f := range files {
		fmt.Printf("  %s\n", f)
	}

	// Sample output
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
	fmt.Println("Pipeline complete!")
}
