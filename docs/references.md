# Word — Source References

This document catalogs all upstream source materials, their licensing status, and retrieval methods.

---

## 1. King James Version (KJV) — English

### Source
**eBible.org — King James Version (eng-kjv2006)**

- **URL:** https://ebible.org/find/show.php?id=eng-kjv2006
- **Download:** https://ebible.org/Scriptures/eng-kjv2006_usfm.zip
- **Format:** USFM (Unified Standard Format Markers)
- **Content:** Protocanon only (66 books), standardized 1769 text with Strong's numbers

### License
**Public Domain**

The King James Version (1769 standardized text) is in the Public Domain worldwide, with one exception: letters patent issued by King James I apply within the United Kingdom, requiring permission to print or import printed copies. This restriction has no effect outside the UK and does not apply to digital distribution.

### Attribution
No attribution required. Optional credit: "King James Version from eBible.org"

---

## 2. Westminster Leningrad Codex (WLC) — Hebrew Old Testament

### Source
**Open Scriptures Hebrew Bible (OSHB / morphhb)**

- **Repository:** https://github.com/openscriptures/morphhb
- **Text Files:** https://github.com/openscriptures/morphhb/tree/master/wlc/
- **Format:** OSIS XML
- **NPM Package:** https://www.npmjs.com/package/morphhb (JSON export)

### License
**Dual License — Important Distinction**

| Component | License | Notes |
|-----------|---------|-------|
| Base consonantal text (WLC) | **Public Domain** | No restrictions |
| Lemma and morphology data | **CC BY 4.0** | Requires attribution |

For this project, we use the **base text only** to maintain clean Public Domain status. Morphological data may be added as an optional enrichment layer with proper attribution.

### Attribution (if using morphology)
"Lemma and morphology data from the Open Scriptures Hebrew Bible Project (https://hb.openscriptures.org/), licensed under CC BY 4.0."

---

## 3. Scrivener 1894 Textus Receptus — Greek New Testament

### Source
**Byzantine Text Project — greektext-scrivener**

- **Repository:** https://github.com/byztxt/greektext-scrivener
- **Text Files:** https://github.com/byztxt/greektext-scrivener/tree/master/textonly
- **Format:** Plain text (UTF-8 Greek)
- **Website:** https://www.byzantinetext.com

### Background
F. H. A. Scrivener (1813–1891) produced this edition to reproduce the Greek text underlying the 1611 King James Version. He used the Beza 1598 edition as his starting point, examining eighteen editions of the Textus Receptus to determine the correct Greek rendering for each KJV reading.

### License
**Public Domain**

Explicitly marked: "Copy freely." No restrictions on use or distribution.

### Attribution
No attribution required. Optional credit: "Greek text edited by Dr. Maurice A. Robinson, maintained by Dr. Ulrik Sandborg-Petersen."

### Alternative Sources
- **Internet Archive (PDF scan):** https://archive.org/details/greek-new-testament-scrivener-textus-receptus-1894
- **With morphology:** https://archive.org/details/textusreceptus

---

## 4. Supplementary Resources

### OSIS Book Codes
- **Specification:** http://www.crosswire.org/osis/
- **Book list:** Standard 3-letter codes (GEN, EXO, LEV, ... MAT, MRK, LUK, JHN, ...)

### Versification References
- **CCEL Versification:** https://www.ccel.org/
- **Paratext Versification:** Standard verse mappings for Bible software

---

## 5. Version Tracking

| Source | Version/Date Retrieved | SHA256 (if applicable) |
|--------|------------------------|------------------------|
| eng-kjv2006 | TBD | TBD |
| morphhb/wlc | v2.0 | TBD |
| greektext-scrivener | main branch | TBD |

*Checksums will be recorded after initial download.*

---

## 6. License Summary

| Work | License | Commercial Use | Attribution Required |
|------|---------|----------------|---------------------|
| KJV (eng-kjv2006) | Public Domain | Yes | No |
| WLC base text | Public Domain | Yes | No |
| WLC morphology | CC BY 4.0 | Yes | Yes |
| TR 1894 | Public Domain | Yes | No |

All source texts selected for this project are **Public Domain** to ensure maximum reusability and no legal encumbrances on derived works, including LLM-generated translations.
