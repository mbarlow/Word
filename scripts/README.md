# Word Scripts

Bash scripts for pipeline operations and development tasks.

## Script Inventory

| Script | Purpose |
|--------|---------|
| `download-sources.sh` | Download upstream source texts (KJV, WLC, TR1894) |
| `pipeline.sh` | Run full pipeline: ingest → normalize → render → validate |
| `ingest.sh` | Parse raw sources into canonical JSONL |
| `render.sh` | Generate Markdown from canonical dataset |
| `validate.sh` | Verify verse/chapter counts and run spot checks |

## Usage

All scripts should be run from the repository root:

```bash
./scripts/download-sources.sh
./scripts/pipeline.sh
```

## Download Sources

Downloads upstream artifacts to `data/source/`:

```bash
./scripts/download-sources.sh

# Output:
# data/source/kjv_raw/       - eBible.org KJV USFM files
# data/source/wlc_raw/       - OSHB WLC XML files
# data/source/tr1894_raw/    - Byzantine Text Greek files
```

Checksums are verified against known values. See `docs/references.md` for source URLs.

## Pipeline

Run individual stages or the full pipeline:

```bash
# Full pipeline
./scripts/pipeline.sh

# Individual stages
./scripts/ingest.sh      # Raw → JSONL
./scripts/render.sh      # JSONL → Markdown
./scripts/validate.sh    # Verify output
```

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `WORD_DATA_DIR` | `./data` | Base directory for data files |
| `WORD_TEXTS_DIR` | `./texts` | Output directory for Markdown |
| `WORD_LOG_LEVEL` | `info` | Log verbosity (debug, info, warn, error) |

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | General error |
| 2 | Download failed |
| 3 | Checksum mismatch |
| 4 | Validation failed |
