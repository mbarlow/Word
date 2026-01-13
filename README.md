# Word

Canonical Bible text repository with reproducible pipeline and translation generation system.

## Overview

Word is a production-grade system for:

1. **Canonical Source Repository** — Public-domain Bible texts (KJV, Hebrew WLC, Greek TR 1894)
2. **Reproducible Pipeline** — Download, normalize, validate, and render texts to Markdown
3. **Translation Generation** — LLM-powered draft translations with provenance tracking
4. **API Backend** — RESTful API for verse lookup, search, and cross-reference

## Quick Start

```bash
# Start local development environment
tilt up

# Download source texts
./scripts/download-sources.sh

# Run full pipeline (ingest → normalize → render → validate)
./scripts/pipeline.sh

# API available at http://localhost:8080
```

## Project Structure

```
Word/
├── cmd/
│   ├── server/           # API server entry point
│   └── pipeline/         # CLI for ingest, translate, render
├── internal/
│   ├── api/              # HTTP handlers, middleware, routing
│   ├── config/           # Environment-based configuration
│   ├── model/            # GORM models (Verse, Work, Book)
│   ├── repository/       # Database queries
│   ├── service/          # Business logic
│   ├── pipeline/         # Ingest, normalize, render, validate stages
│   └── translation/      # LLM translation profiles, prompts, clients
├── data/
│   ├── canon/            # Canonical metadata (books.json, versification)
│   └── source/           # Downloaded raw files (gitignored)
├── texts/                # Rendered Markdown output
│   ├── kjv/
│   ├── heb-wlc/
│   └── grc-tr1894/
├── translations/
│   ├── profiles/         # Translation profile configs (JSON)
│   └── drafts/           # Generated translations (JSONL)
├── experiments/          # DSL experiments (Hebrew as code)
│   ├── genesis.dsl       # Genesis 1 as Lisp-like executable
│   └── dsl/main.go       # Interpreter that "runs" creation
├── scripts/              # Bash scripts for pipeline operations
├── docs/
│   ├── adr/              # Architecture Decision Records
│   └── references.md     # Source URLs and licenses
├── deploy/               # Dockerfile
├── docker-compose.yml
├── Tiltfile
└── LICENSES/
```

## Source Texts

| Work | Language | Source | License |
|------|----------|--------|---------|
| KJV (1769) | English | [eBible.org](https://ebible.org/find/show.php?id=eng-kjv2006) | Public Domain |
| WLC | Hebrew | [OSHB/morphhb](https://github.com/openscriptures/morphhb) | Public Domain |
| TR 1894 | Greek | [byztxt](https://github.com/byztxt/greektext-scrivener) | Public Domain |

See [docs/references.md](docs/references.md) for complete source documentation.

## Pipeline Stages

### 1. Ingest
Download and parse upstream artifacts into normalized verse objects.

```bash
./scripts/download-sources.sh  # Fetch upstream files
go run ./cmd/pipeline ingest   # Parse into JSONL
```

### 2. Normalize
Apply Unicode normalization (NFC), whitespace cleanup, and verse marker standardization.

### 3. Render
Generate Markdown chapter files from canonical dataset.

```
texts/kjv/genesis/1.md
texts/heb-wlc/genesis/1.md
texts/grc-tr1894/matthew/1.md
```

### 4. Validate
Verify verse counts, chapter counts, and spot-check against known references.

## API Endpoints

```
GET /v1/works                              # List available works
GET /v1/text/{work}/{book}/{chapter}       # Get chapter text
GET /v1/verse/{work}/{book}/{chapter}/{verse}  # Get single verse
GET /v1/compare?works=kjv,heb-wlc&ref=GEN.1.1  # Cross-reference
GET /v1/search?q=faith&work=kjv            # Full-text search
GET /metrics                               # Prometheus metrics
GET /health                                # Health check
```

## Verse ID Schema

### Canonical ID (VID)
```
<work>/<osis>/<chapter>/<verse>
```
Examples: `kjv/GEN/1/1`, `heb-wlc/PSA/23/1`, `grc-tr1894/JHN/3/16`

### Human Reference
```
<book-slug>.<chapter>.<verse>
```
Examples: `genesis.1.1`, `john.3.16`, `1-peter.1.3`

## Translation Generation

Word supports LLM-powered translation drafting with:

- **Profile-driven generation** — Configure audience, reading level, terminology policies
- **Two-pass workflow** — Draft translator + critic/QA pass
- **Provenance tracking** — Model, prompt hash, timestamp, source verse IDs
- **Hot verse flagging** — Extra review for doctrinally significant passages

### Translation Profiles

| Profile | Format | Description |
|---------|--------|-------------|
| `techdoc_en` | Prose | Clear, accurate English for technical readers |
| `kids_en` | Prose | Simplified language for younger audiences |
| `structural_en` | 6-layer scholarly | Root analysis, pseudo-code, cognates, polysemy, literary devices |
| `lisp_en` | S-expressions | Hebrew as functional programming: `(yomer ELOHIM (yehi OR))` |
| `yaml_en` | Declarative YAML | Structured data: `speaker: ELOHIM, command: {action: yehi}` |

### Running Translations

```bash
# Set up Ollama (or use ANTHROPIC_API_KEY for Claude)
export OLLAMA_HOST=http://localhost:11434
export OLLAMA_MODEL=gemma3:4b

# Translate a single verse
go run ./cmd/pipeline translate structural_en heb-wlc/GEN/1/11

# Translate a full chapter
go run ./cmd/pipeline translate-chapter structural_en heb-wlc GEN 1

# Compare translations across models
./scripts/compare-translations.sh GEN 1 11
```

### Structural Translation Output

The `structural_en` profile produces 6 layers of scholarly analysis:

```json
{
  "text": "And God said, 'Let the earth **vegetate¹ vegetation¹**—plants **seeding² seed²**'",
  "root_analysis": "[1] דשא (d-sh-a): תַּדְשֵׁא (hiphil) ↔ דֶּשֶׁא (noun) — cognate_accusative",
  "structural": "god.say(earth.vegetate<דשא>(vegetation<דשא>)) => TRUE",
  "cognates": [
    {"marker": 1, "root": "דשא", "verb": "תַּדְשֵׁא", "noun": "דֶּשֶׁא", "gloss": "vegetate/vegetation", "pattern": "cognate_accusative"}
  ],
  "polysemy": {"אֱלֹהִים": "God|gods|divine-council", "אֶרֶץ": "earth|land|ground"},
  "literary_devices": ["cognate_accusative", "merism", "inclusio"]
}
```

**Literary devices detected:**
- **Cognate Accusative** (Figura Etymologica) — Verb + noun from same root for emphasis
- **Merism** — Two extremes representing totality ("heavens and earth" = everything)
- **Chiasm** — ABBA inverted parallelism
- **Inclusio** — Bookend repetition for closure

## Experiments: Hebrew as Executable Code

The `experiments/` folder explores treating Hebrew scripture as a programming language.

### Genesis DSL — The Creation Executable

```bash
# Run Genesis 1 as an executable program
go run ./experiments/dsl/main.go
```

This interprets `experiments/genesis.dsl` — a Lisp-like representation of Genesis 1 where:

- `(yomer ELOHIM ...)` = God.say(command)
- `(yehi OR)` = "Let there be light"
- `(va-yehi KEN)` = returns TRUE ("and it was so")
- `(bara ELOHIM ...)` = God.create(ex nihilo)

**Sample output:**
```
╔══════════════════════════════════════════════════════════════╗
║          GENESIS.DSL — The Creation Executable               ║
║              בְּרֵאשִׁית בָּרָא אֱלֹהִים                              ║
╚══════════════════════════════════════════════════════════════╝

► EXEC: (yomer ELOHIM (yehi OR))
  וַיֹּאמֶר אֱלֹהִים יְהִי אוֹר
  ✓ (va-yehi OR) → Light instantiated

──────────────────────────────────────────────────────────────
[DAY 1] LIGHT separated from DARKNESS
──────────────────────────────────────────────────────────────

...

UNIVERSE := {
  exists:  true
  light:   true
  sky:     true
  land:    true
  ...
  humans:  true  ← image_of(ELOHIM)
  status:  COMPLETE
  eval:    TOV MEOD (Very Good)
}

// Process exited. Universe running.
// To inspect: use conscience, prayer, or telescope.
```

See [experiments/README.md](experiments/README.md) for full documentation.

### The Infinite Light Loop 💡

During development, an LLM got stuck translating Genesis into Lisp and produced:

```lisp
(yehi OR)      ;; "Let there be light"
(va-yehi OR)   ;; "And there was light"
(yehi OR)      ;; "Let there be light"
(va-yehi OR)   ;; "And there was light"
... (forever)
```

**Translation:** An eternal creation loop — endlessly commanding light into existence and watching it appear. A recursive function without a base case. The universe's first infinite loop.

The LLM forgot to return `(va-yehi KEN)` (TRUE) to exit the function!

## Development

### Prerequisites
- Go 1.21+
- Docker & Docker Compose
- Tilt

### Local Development
```bash
tilt up        # Start all services with hot reload
tilt down      # Stop services
```

### Running Tests
```bash
go test ./...
```

### Building
```bash
docker build -t word-service -f deploy/Dockerfile .
```

## Architecture Decisions

See [docs/adr/](docs/adr/) for Architecture Decision Records:

- [ADR-0001: Go with Echo and GORM](docs/adr/0001-go-echo-gorm.md)
- [ADR-0002: SQLite as Canonical Store](docs/adr/0002-sqlite-canonical-store.md)
- [ADR-0003: Source Text Selection](docs/adr/0003-source-texts.md)

## License

Source texts are Public Domain. See [LICENSES/](LICENSES/) for individual text licenses.

Project code is MIT licensed.
