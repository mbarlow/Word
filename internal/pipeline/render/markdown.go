package render

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/mbarlow/word/internal/pipeline"
)

// MarkdownRenderer renders verses to Markdown chapter files.
type MarkdownRenderer struct {
	outputDir string
}

// NewMarkdownRenderer creates a new renderer.
func NewMarkdownRenderer(outputDir string) *MarkdownRenderer {
	return &MarkdownRenderer{outputDir: outputDir}
}

// RenderChapter renders a chapter's verses to a Markdown file.
func (r *MarkdownRenderer) RenderChapter(verses []pipeline.Verse) error {
	if len(verses) == 0 {
		return nil
	}

	// Sort verses by verse number
	sort.Slice(verses, func(i, j int) bool {
		return verses[i].Verse < verses[j].Verse
	})

	first := verses[0]

	// Determine output path
	dir := filepath.Join(r.outputDir, first.Work, first.BookSlug)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	path := filepath.Join(dir, fmt.Sprintf("%d.md", first.Chapter))

	// Build content
	var sb strings.Builder

	// Frontmatter
	sb.WriteString("---\n")
	sb.WriteString(fmt.Sprintf("work: %s\n", first.Work))
	sb.WriteString(fmt.Sprintf("book: %s\n", first.BookSlug))
	sb.WriteString(fmt.Sprintf("chapter: %d\n", first.Chapter))
	sb.WriteString(fmt.Sprintf("lang: %s\n", first.Lang))
	sb.WriteString("---\n\n")

	// Chapter heading
	bookName := strings.Title(strings.ReplaceAll(first.BookSlug, "-", " "))
	sb.WriteString(fmt.Sprintf("# %s %d\n\n", bookName, first.Chapter))

	// Verses as numbered list
	for _, v := range verses {
		sb.WriteString(fmt.Sprintf("%d. %s\n", v.Verse, v.Text))
	}

	return os.WriteFile(path, []byte(sb.String()), 0644)
}

// RenderAll renders all verses grouped by work/book/chapter.
func (r *MarkdownRenderer) RenderAll(verses []pipeline.Verse) error {
	// Group by work/book/chapter
	groups := make(map[string][]pipeline.Verse)

	for _, v := range verses {
		key := fmt.Sprintf("%s/%s/%d", v.Work, v.BookSlug, v.Chapter)
		groups[key] = append(groups[key], v)
	}

	for _, chapterVerses := range groups {
		if err := r.RenderChapter(chapterVerses); err != nil {
			return err
		}
	}

	return nil
}
