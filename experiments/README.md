# Experiments: Hebrew Scripture as Executable Code

This folder explores the hypothesis that Hebrew scripture can be understood as a declarative/functional programming language.

## Genesis DSL — The Creation Executable

Genesis 1 represented as a Lisp-like domain-specific language where divine speech acts are function calls that return observable state changes.

### Quick Start

```bash
# Run the Genesis executable (auto-locates genesis.dsl)
go run ./experiments/dsl

# Skip the dramatic pauses
go run ./experiments/dsl -fast

# Run a specific program
go run ./experiments/dsl path/to/program.dsl
```

### Files

| File | Description |
|------|-------------|
| `genesis.dsl` | Genesis 1 as Lisp-like Hebrew DSL (source program) |
| `dsl/reader.go` | Lexer + parser — source text to s-expression AST |
| `dsl/interp.go` | Evaluator — dispatches verbs, instantiates the Universe |
| `dsl/main.go` | Entry point — source resolution, flags, render |

## The DSL Syntax

### Core Verbs (Functions)

| Hebrew | Transliteration | DSL Form | Meaning |
|--------|-----------------|----------|---------|
| וַיֹּאמֶר | va-yomer | `(yomer ELOHIM ...)` | God said (speech act) |
| יְהִי | yehi | `(yehi X)` | Let there be X (jussive command) |
| וַיְהִי כֵן | va-yehi ken | `(va-yehi KEN)` | And it was so → returns TRUE |
| בָּרָא | bara | `(bara ELOHIM ...)` | Create ex nihilo |
| עָשָׂה | asah | `(va-yaas ELOHIM ...)` | Make/produce |
| וַיַּרְא | va-yar | `(va-yar ELOHIM ki TOV)` | God saw that it was good |
| וַיְבָרֶךְ | va-yevarekh | `(va-yevarekh ELOHIM ...)` | God blessed |
| וַיִּקְרָא | va-yikra | `(va-yikra ELOHIM ...)` | God called/named |

### Naming Convention

- **CAPS** for divine/proper names: `ELOHIM`, `YHWH`, `HA-SHAMAYIM`, `HA-ARETZ`
- **Keywords** for attributes: `:state`, `:result`, `:kind`, `:le-mino`
- **Comments** with `;;` for Hebrew text and explanations

### Example: Day 1 (Light)

```lisp
;; 1:3-5 — DAY 1: LIGHT

(yomer ELOHIM                         ;; God said
  (yehi OR))                          ;; "Let there be light"
(va-yehi OR)                          ;; And there was light → OUTPUT

(va-yar ELOHIM et HA-OR ki TOV)       ;; God saw light: GOOD
(va-yavdel ELOHIM                     ;; God separated
  :between HA-OR
  :and     HA-CHOSHEKH)

(va-yikra ELOHIM                      ;; God called/named
  (la-or    YOM)                      ;; light → "Day"
  (la-choshekh LAILAH))               ;; darkness → "Night"

(va-yehi EREV va-yehi VOKER           ;; evening, morning
  :day 1)                             ;; Day One (אֶחָד)
```

### Example: Day 3 (Cognate Accusative)

```lisp
;; COGNATE STRUCTURE: תַּדְשֵׁא...דֶּשֶׁא (vegetate vegetation)
(yomer ELOHIM
  (tadshe HA-ERETZ DESHE              ;; earth.vegetate(vegetation)
    (ESEV mazria ZERA)                ;; plant.seed(seed)
    (ETZ PRI oseh PRI                 ;; tree.fruit(fruit)
      :le-mino   SELF                 ;; after its kind
      :zaro-vo   SELF)))              ;; seed in itself
(va-yehi KEN)
```

This demonstrates **cognate accusative** (figura etymologica) — verb and noun from same root:
- תַּדְשֵׁא דֶּשֶׁא (tadshe deshe) = "vegetate vegetation"
- מַזְרִיעַ זֶרַע (mazria zera) = "seeding seed"

In programming terms: `function.call(function.result)` — recursion/self-reference.

## How the Interpreter Works

It is a real interpreter, not a print script:

1. **`reader.go`** — lexes the source (discarding `;` comments, handling
   strings and parens) and parses it into a tree of s-expression `Node`s.
2. **`interp.go`** — walks the top-level forms and dispatches each on its
   head verb (`yomer`, `bara`, `yehi`, `va-yehi`, `va-yar`, ...).
3. Creative verbs scan their form for known Hebrew nouns (`OR`, `RAKIA`,
   `DESHE`, `ADAM`, ...) and set the matching field on a `Universe` struct.
   Proclitic prefixes (`HA-`, `ve-`, `u-`) are peeled before lookup.
4. **`va-yehi ... :day N`** commits a day boundary and prints running state.
5. **`(return UNIVERSE)`** renders the final struct. Exit code is `0` only
   if every field was actually instantiated and status is `COMPLETE`.

Two analyses run as forms evaluate:

- **Cognate detection** — verb/noun pairs whose consonant skeletons nest
  (`tadshe ↔ DESHE`, `yishretzu ↔ SHERETZ`) are flagged as figura
  etymologica. No longer hardcoded — derived from the parsed form.
- **Bara-root emphasis** — counts occurrences of the create-root in a form
  (`va-yivra ... bara ... bara` on day 6 → "used 3×").

```
► EXEC  (yomer ELOHIM (tadshe HA-ERETZ DESHE (ESEV mazria ZERA) ...))
  ELOHIM speaks — evaluating command(s):

► EXEC  (tadshe HA-ERETZ DESHE (ESEV mazria ZERA) (ETZ PRI oseh PRI ...))
  ✓ instantiated: vegetation, plants, trees
  ≈ cognate: tadshe ↔ DESHE (verb/noun share a root)

──────────────────────────────────────────────────────────
[DAY 3]  evaluation: TOV (good)
──────────────────────────────────────────────────────────
  ✓ light    true
  ✓ sky      true
  ✓ land     true
  ✓ seas     true
  ✓ plants   true
  · sun      false
  ...

UNIVERSE := {
  exists:  true   light:  true   sky:   true   land:    true
  seas:    true   plants: true   sun:   true   moon:    true
  stars:   true   fish:   true   birds: true   animals: true
  humans:  true  ← image_of(ELOHIM)
  days:    7      status: COMPLETE      eval: TOV MEOD (very good)
}

// process exited 0 — universe running.
```

## Literary Devices as Programming Patterns

| Hebrew Device | Programming Analog | Example |
|---------------|-------------------|---------|
| Cognate Accusative | `fn(fn.result)` — recursion | vegetate(vegetation) |
| Merism | `range(MIN, MAX)` — bounds | heavens + earth = all |
| Chiasm (ABBA) | Stack push/pop | inverted parallelism |
| Inclusio | `try { } finally { }` | bookend repetition |
| Repetition | `assert(x == x)` | identity emphasis |

## Future Work

### Hebrew DSL Profile (`genesis_dsl`)

A translation profile that generates DSL output for any Old Testament passage:

```bash
# Future command
go run ./cmd/pipeline translate genesis_dsl heb-wlc/EXO/20/1-17
```

### Greek DSL (Koine)

Extending the DSL concept to Greek New Testament:

| Greek | Transliteration | DSL Form | Meaning |
|-------|-----------------|----------|---------|
| Ἐν ἀρχῇ | en arche | `(en ARCHE ...)` | In the beginning |
| ἦν | en | `(en LOGOS)` | was (existence) |
| ἐγένετο | egeneto | `(egeneto X)` | became/came into being |
| λέγω | lego | `(lego ...)` | I say |

**John 1:1-3 as Greek DSL:**

```lisp
;; JOHN 1:1-3 — The Logos Prologue

(en ARCHE                             ;; In the beginning
  (en HO-LOGOS))                      ;; was the Word

(kai HO-LOGOS                         ;; And the Word
  (en PROS TON-THEON))                ;; was with God

(kai THEOS                            ;; And God
  (en HO-LOGOS))                      ;; was the Word

(panta di AUTOU egeneto               ;; All things through him came into being
  (kai choris AUTOU                   ;; and apart from him
    (egeneto OUDE-HEN)))              ;; came into being nothing

;; Return: LOGOS = { divine: true, creator: true, with_god: true }
```

## The Infinite Light Loop

During development, an LLM translating Genesis into Lisp got stuck in an infinite loop:

```lisp
(yehi OR)      ;; "Let there be light"
(va-yehi OR)   ;; "And there was light"
(yehi OR)      ;; "Let there be light"
(va-yehi OR)   ;; "And there was light"
... (forever)
```

**What happened:** The model got caught in an eternal creation loop — endlessly commanding light into existence and watching it appear. A recursive function without a base case.

**The fix:** Return `(va-yehi KEN)` — "and it was so" — which evaluates to TRUE and exits the function.

**Theological interpretation:** Even in code, creation needs a Sabbath rest. Without `(va-yishbot)` (he rested), the process never terminates!

## Related Files

- [Main README](../README.md) — Project overview
- [Translation Profiles](../translations/profiles/) — JSON profile configs
- [Structural Profile Spec](../../obsidian/plans/projects/Word/structural_en.md) — 6-layer output format
