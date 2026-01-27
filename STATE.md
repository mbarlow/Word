# Project State

## Current Focus

- Translation profile experimentation (5 profiles: techdoc, kids, structural, lisp, yaml)
- Cognate detection system for Hebrew text analysis
- Greek text integration (TR 1894 recently merged)

## Active Decisions

- SQLite for canonical verse storage (simple, portable)
- Two-pass translation workflow (draft + QA critic)
- Provenance tracking on all generated translations

## Known Risks

- Translation consistency across key theological terms needs monitoring
- Full-text search performance not yet benchmarked
- Hebrew DSL experiment is exploratory (not production)

## Next Concrete Actions

- [ ] Extend cognate detection to Greek texts
- [ ] Add consistency tracking for theological terms
- [ ] Test translation profiles on more chapters
- [ ] Document hot verse flagging criteria

## Last Updated

2026-01-26
