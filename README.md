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
├── cmd/server/           # API server entry point
├── internal/
│   ├── api/              # HTTP handlers, middleware, routing
│   ├── config/           # Environment-based configuration
│   ├── model/            # GORM models (Verse, Work, Book)
│   ├── repository/       # Database queries
│   ├── service/          # Business logic
│   └── pipeline/         # Ingest, normalize, render, validate stages
├── data/
│   ├── canon/            # Canonical metadata (books.json, versification)
│   └── source/           # Downloaded raw files (gitignored)
├── texts/                # Rendered Markdown output
│   ├── kjv/
│   ├── heb-wlc/
│   └── grc-tr1894/
├── translations/         # LLM translation profiles and output
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
