# ADR-0003: Source Text Selection

## Status
Accepted

## Date
2026-01-11

## Context
We need to select authoritative, public-domain source texts for:

1. English Bible (baseline reference)
2. Hebrew Old Testament (translation source)
3. Greek New Testament (translation source)

Requirements:
- Must be Public Domain (no licensing restrictions on derived works)
- Must have stable, downloadable digital editions
- Must align historically with traditional Protestant canon

## Decision

### English: King James Version (1769)
**Source:** eBible.org eng-kjv2006

- Standardized 1769 text
- Public Domain (except UK print restrictions, which don't apply digitally)
- USFM format with Strong's numbers available
- Well-established digital lineage

### Hebrew OT: Westminster Leningrad Codex (WLC)
**Source:** Open Scriptures Hebrew Bible (morphhb)

- Masoretic Text tradition
- Base consonantal text is Public Domain
- Morphology data (CC BY 4.0) kept separate as optional enrichment
- OSIS XML format, actively maintained

### Greek NT: Scrivener 1894 Textus Receptus
**Source:** Byzantine Text Project (byztxt/greektext-scrivener)

- Explicitly reconstructs the Greek underlying the KJV
- Public Domain ("Copy freely")
- Plain text format, easy to parse
- Maintained by Dr. Maurice A. Robinson

## Rationale

### Why KJV?
- Most widely recognized English translation
- Unambiguous Public Domain status
- Historical alignment with TR Greek text
- Strong's numbers enable Hebrew/Greek cross-referencing

### Why WLC over BHS?
- BHS (Biblia Hebraica Stuttgartensia) has copyright restrictions
- WLC is the same underlying text, freely available
- OSHB project provides high-quality digital edition

### Why TR 1894 over NA28/UBS5?
- NA28 and UBS5 are copyrighted critical editions
- TR 1894 is Public Domain
- Historical consistency: TR underlies KJV, enabling direct comparison
- Suitable for traditional/KJV-aligned translation profiles

## Consequences

### Positive
- All source texts are unambiguously Public Domain
- Derived translations inherit no licensing restrictions
- Historical consistency across English/Hebrew/Greek

### Negative
- TR differs from modern critical Greek texts (intentional trade-off)
- Some scholars prefer NA28/UBS5 for academic work
- WLC morphology requires separate attribution if used

## Alternatives Considered

### NIV/ESV/NASB for English baseline
- All copyrighted; rejected for licensing reasons

### BHS for Hebrew
- Copyright restrictions on apparatus; WLC base text preferred

### NA28/UBS5 for Greek
- Copyrighted; incompatible with Public Domain goals

## References
- See `docs/references.md` for full URLs and download instructions
