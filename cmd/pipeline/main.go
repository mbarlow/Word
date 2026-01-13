package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/mbarlow/word/internal/model"
	"github.com/mbarlow/word/internal/pipeline"
	"github.com/mbarlow/word/internal/pipeline/ingest"
	"github.com/mbarlow/word/internal/pipeline/render"
	"github.com/mbarlow/word/internal/pipeline/store"
	"github.com/mbarlow/word/internal/pipeline/validate"
	"github.com/mbarlow/word/internal/translation"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: pipeline <command> [args]")
		fmt.Println("Commands:")
		fmt.Println("  ingest <work> <book>     - Ingest a book (work: kjv, heb-wlc, grc-tr1894)")
		fmt.Println("  render                   - Render all JSONL to Markdown")
		fmt.Println("  validate                 - Validate rendered output")
		fmt.Println("  genesis1                 - Run full pipeline for Genesis 1 (all works)")
		fmt.Println("  load                     - Load sample data into SQLite")
		fmt.Println("  profiles                 - List available translation profiles")
		fmt.Println("  translate <profile> <vid> - Translate a verse (e.g., translate techdoc_en heb-wlc/GEN/1/1)")
		fmt.Println("  translate-chapter <profile> <work> <book> <chapter> - Translate a full chapter")
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

	case "profiles":
		runListProfiles()

	case "translate":
		if len(os.Args) < 4 {
			fmt.Println("Usage: pipeline translate <profile> <vid>")
			fmt.Println("Example: pipeline translate techdoc_en heb-wlc/GEN/1/1")
			os.Exit(1)
		}
		runTranslate(os.Args[2], os.Args[3])

	case "translate-chapter":
		if len(os.Args) < 6 {
			fmt.Println("Usage: pipeline translate-chapter <profile> <work> <book> <chapter>")
			fmt.Println("Example: pipeline translate-chapter techdoc_en heb-wlc GEN 1")
			os.Exit(1)
		}
		chapter, _ := strconv.Atoi(os.Args[5])
		runTranslateChapter(os.Args[2], os.Args[3], os.Args[4], chapter)

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

func runListProfiles() {
	profiles, err := translation.ListProfiles("translations/profiles")
	if err != nil {
		fmt.Printf("Failed to list profiles: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Available translation profiles:")
	for _, p := range profiles {
		profile, _ := translation.LoadProfileByName("translations/profiles", p)
		if profile != nil {
			fmt.Printf("  %s - %s (%s)\n", p, profile.Audience, profile.TargetLang)
		} else {
			fmt.Printf("  %s\n", p)
		}
	}
}

func runTranslate(profileName, vid string) {
	// Load profile
	profile, err := translation.LoadProfileByName("translations/profiles", profileName)
	if err != nil {
		fmt.Printf("Failed to load profile: %v\n", err)
		os.Exit(1)
	}

	// Get source verse from database
	db, err := store.NewSQLiteStore("data/word.db")
	if err != nil {
		fmt.Printf("Failed to open database: %v\n", err)
		os.Exit(1)
	}

	// Look up source verse
	sourceText, kjvRef, err := lookupVerse(db, vid)
	if err != nil {
		fmt.Printf("Failed to look up verse %s: %v\n", vid, err)
		os.Exit(1)
	}

	// Get LLM client
	client := getLLMClient()
	if client == nil {
		fmt.Println("No LLM configured. Set OLLAMA_HOST or ANTHROPIC_API_KEY")
		os.Exit(1)
	}

	// Create translator
	translator := translation.NewTranslator(client, profile, 0.2)

	// Translate
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	source := translation.SourceVerse{
		VID:    vid,
		Text:   sourceText,
		KJVRef: kjvRef,
	}

	fmt.Printf("Translating %s with profile %s...\n", vid, profileName)
	fmt.Printf("Hot verse: %v\n", translation.IsHotVerse(vid))
	fmt.Println()

	draft, err := translator.TranslateVerse(ctx, source)
	if err != nil {
		fmt.Printf("Translation failed: %v\n", err)
		os.Exit(1)
	}

	// Output draft JSON
	output, _ := json.MarshalIndent(draft, "", "  ")
	fmt.Println("=== Draft Translation ===")
	fmt.Println(string(output))

	// Run critic pass
	fmt.Println()
	fmt.Println("=== Critic Pass ===")
	critic, err := translator.CriticVerse(ctx, source, draft)
	if err != nil {
		fmt.Printf("Critic pass failed: %v\n", err)
	} else {
		criticOutput, _ := json.MarshalIndent(critic, "", "  ")
		fmt.Println(string(criticOutput))
	}
}

func runTranslateChapter(profileName, work, book string, chapter int) {
	// Load profile
	profile, err := translation.LoadProfileByName("translations/profiles", profileName)
	if err != nil {
		fmt.Printf("Failed to load profile: %v\n", err)
		os.Exit(1)
	}

	// Get verses from database
	db, err := store.NewSQLiteStore("data/word.db")
	if err != nil {
		fmt.Printf("Failed to open database: %v\n", err)
		os.Exit(1)
	}

	// Get all verses in chapter
	verses, err := getChapterVerses(db, work, book, chapter)
	if err != nil {
		fmt.Printf("Failed to get chapter: %v\n", err)
		os.Exit(1)
	}

	if len(verses) == 0 {
		fmt.Printf("No verses found for %s/%s/%d\n", work, book, chapter)
		os.Exit(1)
	}

	// Get LLM client
	client := getLLMClient()
	if client == nil {
		fmt.Println("No LLM configured. Set OLLAMA_HOST or ANTHROPIC_API_KEY")
		os.Exit(1)
	}

	// Create translator
	translator := translation.NewTranslator(client, profile, 0.2)

	fmt.Printf("Translating %s %s chapter %d (%d verses) with profile %s\n",
		work, book, chapter, len(verses), profileName)
	fmt.Println()

	ctx := context.Background()

	var drafts []*translation.Draft
	var critics []*translation.CriticResult

	for i, v := range verses {
		vid := v.VID
		fmt.Printf("[%d/%d] %s... ", i+1, len(verses), vid)

		// Get KJV reference if available
		kjvRef := ""
		if work != "kjv" {
			kjvVID := fmt.Sprintf("kjv/%s/%d/%d", book, chapter, v.Verse)
			if ref, _, err := lookupVerse(db, kjvVID); err == nil {
				kjvRef = ref
			}
		}

		source := translation.SourceVerse{
			VID:    vid,
			Text:   v.Text,
			KJVRef: kjvRef,
		}

		// Translate
		draft, err := translator.TranslateVerse(ctx, source)
		if err != nil {
			fmt.Printf("FAILED: %v\n", err)
			continue
		}
		drafts = append(drafts, draft)

		// Critic
		critic, err := translator.CriticVerse(ctx, source, draft)
		if err != nil {
			fmt.Printf("translated (critic failed)\n")
		} else {
			critics = append(critics, critic)
			flagCount := len(critic.Flags)
			if flagCount > 0 {
				fmt.Printf("OK (%d flags)\n", flagCount)
			} else {
				fmt.Println("OK")
			}
		}
	}

	// Output results
	fmt.Println()
	fmt.Printf("=== Translation Complete: %d verses ===\n", len(drafts))

	// Write drafts to JSONL (include model name for comparison)
	modelSlug := sanitizeModelName(client.ModelName())
	outputDir := fmt.Sprintf("translations/drafts/t_%s/%s/%s", profileName, modelSlug, book)
	os.MkdirAll(outputDir, 0755)
	outputFile := filepath.Join(outputDir, fmt.Sprintf("%d.jsonl", chapter))

	f, err := os.Create(outputFile)
	if err != nil {
		fmt.Printf("Failed to create output file: %v\n", err)
	} else {
		defer f.Close()
		for _, d := range drafts {
			line, _ := json.Marshal(d)
			f.WriteString(string(line) + "\n")
		}
		fmt.Printf("Drafts written to: %s\n", outputFile)
	}

	// Write critic results
	criticFile := filepath.Join(outputDir, fmt.Sprintf("%d.critic.jsonl", chapter))
	cf, err := os.Create(criticFile)
	if err != nil {
		fmt.Printf("Failed to create critic file: %v\n", err)
	} else {
		defer cf.Close()
		for _, c := range critics {
			line, _ := json.Marshal(c)
			cf.WriteString(string(line) + "\n")
		}
		fmt.Printf("Critic flags written to: %s\n", criticFile)
	}

	// Summary
	totalFlags := 0
	for _, c := range critics {
		totalFlags += len(c.Flags)
	}
	fmt.Printf("\nTotal flags: %d\n", totalFlags)
}

// getLLMClient returns an LLM client based on environment configuration.
func getLLMClient() translation.LLMClient {
	// Try Ollama first
	if host := os.Getenv("OLLAMA_HOST"); host != "" {
		model := os.Getenv("OLLAMA_MODEL")
		if model == "" {
			model = "llama3.2"
		}
		return translation.NewOllamaClient(host, model)
	}

	// Try Anthropic
	if apiKey := os.Getenv("ANTHROPIC_API_KEY"); apiKey != "" {
		model := os.Getenv("ANTHROPIC_MODEL")
		if model == "" {
			model = "claude-sonnet-4-20250514"
		}
		return translation.NewAnthropicClient(apiKey, model)
	}

	return nil
}

// ChapterVerse is a simple verse holder for chapter translation.
type ChapterVerse struct {
	VID   string
	Verse int
	Text  string
}

// lookupVerse gets verse text from the database.
func lookupVerse(db *store.SQLiteStore, vid string) (text string, kjvRef string, err error) {
	// Query verse directly from underlying DB
	type result struct {
		Text string
	}
	var r result
	if err := db.DB().Raw("SELECT text FROM verses WHERE id = ?", vid).Scan(&r).Error; err != nil {
		return "", "", err
	}
	if r.Text == "" {
		return "", "", fmt.Errorf("verse not found: %s", vid)
	}

	// Try to get KJV reference
	// Parse VID to build KJV VID
	// vid format: work/BOOK/CHAPTER/VERSE
	parts := filepath.SplitList(vid)
	if len(parts) == 1 {
		// Split by /
		import_parts := splitVID(vid)
		if len(import_parts) >= 4 && import_parts[0] != "kjv" {
			kjvVID := fmt.Sprintf("kjv/%s/%s/%s", import_parts[1], import_parts[2], import_parts[3])
			var kjvResult result
			if err := db.DB().Raw("SELECT text FROM verses WHERE id = ?", kjvVID).Scan(&kjvResult).Error; err == nil {
				kjvRef = kjvResult.Text
			}
		}
	}

	return r.Text, kjvRef, nil
}

// sanitizeModelName converts model name to filesystem-safe slug.
func sanitizeModelName(model string) string {
	// Replace colons and slashes with dashes
	result := ""
	for _, c := range model {
		if c == ':' || c == '/' || c == '\\' {
			result += "-"
		} else {
			result += string(c)
		}
	}
	return result
}

// splitVID splits a verse ID into parts.
func splitVID(vid string) []string {
	var parts []string
	current := ""
	for _, c := range vid {
		if c == '/' {
			parts = append(parts, current)
			current = ""
		} else {
			current += string(c)
		}
	}
	if current != "" {
		parts = append(parts, current)
	}
	return parts
}

// getChapterVerses returns all verses for a chapter.
func getChapterVerses(db *store.SQLiteStore, work, book string, chapter int) ([]ChapterVerse, error) {
	var verses []ChapterVerse
	rows, err := db.DB().Raw(
		"SELECT id, verse, text FROM verses WHERE work = ? AND osis = ? AND chapter = ? ORDER BY verse",
		work, book, chapter,
	).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var v ChapterVerse
		if err := rows.Scan(&v.VID, &v.Verse, &v.Text); err != nil {
			return nil, err
		}
		verses = append(verses, v)
	}
	return verses, nil
}
