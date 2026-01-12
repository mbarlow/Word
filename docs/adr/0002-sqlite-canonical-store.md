# ADR-0002: SQLite as Canonical Data Store (Phase 1)

## Status
Accepted

## Date
2026-01-11

## Context
The Word project requires a canonical data store for verse-level data. The store must:

- Support fast lookups by verse ID
- Handle ~31,000 verses per work (KJV has 31,102 verses)
- Be portable and easy to version control (for small datasets)
- Support future migration to PostgreSQL

## Decision
We will use **SQLite** as the canonical data store for Phase 1, with JSONL as an interchange format for pipeline stages.

### Data Flow
```
Upstream Sources → Ingest (JSONL) → SQLite (canonical) → Render (Markdown)
                                  → API (queries)
```

## Rationale

### SQLite
- Zero-configuration, single-file database
- Excellent read performance for our query patterns
- GORM supports SQLite and PostgreSQL with identical model code
- Easy to inspect with standard tooling (`sqlite3` CLI)
- Portable: can ship the entire database as an artifact

### JSONL for Interchange
- Human-readable, line-by-line processing
- Git-friendly diffs for small changes
- Easy to pipe through Unix tools
- Natural format for LLM input/output

### Why Not Parquet
- Overkill for ~100k rows (all works combined)
- Adds dependency complexity (Arrow libraries)
- Better suited for analytics workloads, not transactional API queries

## Consequences

### Positive
- Simple local development (no database container required for basic work)
- Easy to reset: delete `word.db` and re-run ingest
- Can embed database in Docker image for read-only deployments

### Negative
- Single-writer limitation (not an issue for our read-heavy workload)
- No full-text search without extensions (we'll add this in Phase 2 with PostgreSQL)

## Migration Path to PostgreSQL
1. GORM models remain unchanged
2. Update connection string and dialect in config
3. Run migrations
4. Re-import from JSONL canonical export

## File Locations
- **Canonical DB:** `data/word.db`
- **JSONL exports:** `data/export/*.jsonl`
- **Raw downloads:** `data/source/*_raw/` (gitignored)

## References
- SQLite: https://sqlite.org/
- GORM SQLite Driver: https://gorm.io/docs/connecting_to_the_database.html#SQLite
