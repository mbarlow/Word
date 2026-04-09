# Pipeline

The pipeline transforms raw source texts into normalized verse data, rendered Markdown, and a queryable SQLite database.

## Stages

```
raw source files → ingest → normalize → render → validate → SQLite
```

### 1. Ingest

Parses upstream source files into the canonical `Verse` intermediate format (defined in `types.go`).

Each source has its own parser in `ingest/`:

| Parser | Source format | Input path |
|--------|-------------|------------|
| `kjv.go` | USFM | `data/source/kjv_raw/` |
| `wlc.go` | OSHB XML | `data/source/wlc_raw/morphhb/wlc/` |
| `tr1894.go` | Plain text | `data/source/tr1894_raw/greektext-scrivener/textonly/` |

Usage: `go run ./cmd/pipeline ingest <work> <book>` — outputs verse JSON to stdout.

### 2. Normalize

Applies Unicode normalization (NFC), whitespace cleanup, and verse marker standardization. Located in `normalize/`.

### 3. Render

Converts verse data into Markdown chapter files under `texts/`:

```
texts/kjv/genesis/1.md
texts/heb-wlc/genesis/1.md
texts/grc-tr1894/matthew/1.md
```

Implementation in `render/markdown.go`.

### 4. Validate

Verifies verse counts and chapter counts against expected values. Implementation in `validate/validate.go`.

### 5. SQLite Store

Loads verses, works, and canonical book metadata into `data/word.db` via GORM. Implementation in `store/sqlite.go`.

## CLI Commands

| Command | Description |
|---------|-------------|
| `ingest <work> <book>` | Parse a single book from raw sources |
| `render` | Render all JSONL to Markdown (stub) |
| `validate` | Validate rendered output (stub) |
| `genesis1` | Full working demo: ingest + validate + render + SQLite for Genesis 1 and Matthew 1 |
| `load` | Load sample chapters (Genesis 1, Matthew 1, John 1) into SQLite |
| `translate <profile> <vid>` | Two-pass LLM translation (draft + critic) of a single verse |
| `translate-chapter <profile> <work> <book> <ch>` | Translate every verse in a chapter, write JSONL to `translations/drafts/` |
| `cognates <work> <book> <ch>` | Detect Hebrew cognate patterns from OSHB morphology data |
| `profiles` | List available translation profiles |

## Current State

The standalone `render` and `validate` commands are stubs. The working path is `genesis1` or `load`, which run ingest + render + validate + SQLite load together for sample chapters. The full-corpus pipeline is not yet wired up.

## Data Types

The canonical intermediate format is `pipeline.Verse` (see `types.go`), which carries the verse text, its VID (`work/osis/chapter/verse`), language, and source provenance.
