#!/bin/bash
# example-run-all-profiles.sh
# Translate verses using all available profiles and display side-by-side comparison
#
# Usage:
#   ./scripts/example-run-all-profiles.sh              # Default: Genesis 1:1-5
#   ./scripts/example-run-all-profiles.sh GEN 1 11     # Specific verse
#   ./scripts/example-run-all-profiles.sh GEN 1 1-5    # Verse range

set -e

# Colors
RESET="\033[0m"
BOLD="\033[1m"
DIM="\033[2m"
CYAN="\033[36m"
GREEN="\033[32m"
YELLOW="\033[33m"
MAGENTA="\033[35m"
BLUE="\033[34m"

# Configuration
OLLAMA_HOST="${OLLAMA_HOST:-http://localhost:11434}"
OLLAMA_MODEL="${OLLAMA_MODEL:-gemma3:4b}"
export OLLAMA_HOST OLLAMA_MODEL

# Default to Genesis 1:1-5 if no arguments
BOOK="${1:-GEN}"
CHAPTER="${2:-1}"
VERSE_ARG="${3:-1-5}"

# Parse verse range
if [[ "$VERSE_ARG" == *-* ]]; then
    START_VERSE="${VERSE_ARG%-*}"
    END_VERSE="${VERSE_ARG#*-}"
else
    START_VERSE="$VERSE_ARG"
    END_VERSE="$VERSE_ARG"
fi

# Profiles to compare
PROFILES=("techdoc_en" "structural_en" "lisp_en" "yaml_en")

echo -e "${BOLD}${CYAN}"
echo "╔══════════════════════════════════════════════════════════════════════════╗"
echo "║               WORD — Multi-Profile Translation Comparison                ║"
echo "╚══════════════════════════════════════════════════════════════════════════╝"
echo -e "${RESET}"

echo -e "${DIM}Model: ${OLLAMA_MODEL}${RESET}"
echo -e "${DIM}Reference: ${BOOK} ${CHAPTER}:${START_VERSE}-${END_VERSE}${RESET}"
echo ""

# Function to translate a single verse with a profile
translate_verse() {
    local profile="$1"
    local vid="$2"

    go run ./cmd/pipeline translate "$profile" "$vid" 2>/dev/null | \
        jq -r '.text // "ERROR: Translation failed"' 2>/dev/null || \
        echo "ERROR: Translation failed"
}

# Function to get Hebrew source
get_hebrew() {
    local vid="$1"
    sqlite3 data/word.db "SELECT text FROM verses WHERE vid = '$vid'" 2>/dev/null || echo "(Hebrew not available)"
}

# Function to get KJV reference
get_kjv() {
    local book="$1"
    local chapter="$2"
    local verse="$3"
    sqlite3 data/word.db "SELECT text FROM verses WHERE vid = 'kjv/${book}/${chapter}/${verse}'" 2>/dev/null || echo "(KJV not available)"
}

# Process each verse
for ((v = START_VERSE; v <= END_VERSE; v++)); do
    VID="heb-wlc/${BOOK}/${CHAPTER}/${v}"

    echo -e "${BOLD}${YELLOW}"
    echo "════════════════════════════════════════════════════════════════════════════"
    echo "  ${BOOK} ${CHAPTER}:${v}"
    echo "════════════════════════════════════════════════════════════════════════════"
    echo -e "${RESET}"

    # Hebrew Source
    HEBREW=$(get_hebrew "$VID")
    echo -e "${MAGENTA}[Hebrew Source]${RESET}"
    echo -e "${DIM}${HEBREW}${RESET}"
    echo ""

    # KJV Reference
    KJV=$(get_kjv "$BOOK" "$CHAPTER" "$v")
    echo -e "${BLUE}[KJV Reference]${RESET}"
    echo "$KJV"
    echo ""

    # Translate with each profile
    for profile in "${PROFILES[@]}"; do
        case "$profile" in
            techdoc_en)
                COLOR="$GREEN"
                LABEL="Prose (techdoc_en)"
                ;;
            structural_en)
                COLOR="$CYAN"
                LABEL="Structural (structural_en)"
                ;;
            lisp_en)
                COLOR="$YELLOW"
                LABEL="Lisp (lisp_en)"
                ;;
            yaml_en)
                COLOR="$MAGENTA"
                LABEL="YAML (yaml_en)"
                ;;
            *)
                COLOR="$RESET"
                LABEL="$profile"
                ;;
        esac

        echo -e "${COLOR}[${LABEL}]${RESET}"

        TRANSLATION=$(translate_verse "$profile" "$VID")

        # Pretty print based on profile type
        if [[ "$profile" == "lisp_en" || "$profile" == "yaml_en" ]]; then
            # Code-style output - preserve newlines
            echo "$TRANSLATION" | sed 's/\\n/\n/g'
        elif [[ "$profile" == "structural_en" ]]; then
            # For structural, also show extra fields if available
            FULL_OUTPUT=$(go run ./cmd/pipeline translate "$profile" "$VID" 2>/dev/null)
            echo "$TRANSLATION"

            # Show structural pseudo-code if present
            STRUCTURAL=$(echo "$FULL_OUTPUT" | jq -r '.structural // empty' 2>/dev/null)
            if [[ -n "$STRUCTURAL" && "$STRUCTURAL" != "null" ]]; then
                echo -e "${DIM}  Structural: ${STRUCTURAL}${RESET}"
            fi

            # Show literary devices if present
            DEVICES=$(echo "$FULL_OUTPUT" | jq -r '.literary_devices // [] | join(", ")' 2>/dev/null)
            if [[ -n "$DEVICES" ]]; then
                echo -e "${DIM}  Devices: ${DEVICES}${RESET}"
            fi
        else
            echo "$TRANSLATION"
        fi
        echo ""
    done

    # Also show the DSL version from genesis.dsl if this is Genesis 1
    if [[ "$BOOK" == "GEN" && "$CHAPTER" == "1" ]]; then
        echo -e "${BOLD}${GREEN}[Genesis DSL]${RESET}"
        # Extract the relevant section from genesis.dsl based on verse/day
        case "$v" in
            1|2)
                echo "(bara ELOHIM"
                echo "  (et HA-SHAMAYIM)"
                echo "  (ve-et HA-ARETZ))"
                ;;
            3)
                echo "(yomer ELOHIM (yehi OR))"
                echo "(va-yehi OR)"
                ;;
            4|5)
                echo "(va-yar ELOHIM et HA-OR ki TOV)"
                echo "(va-yavdel ELOHIM :between HA-OR :and HA-CHOSHEKH)"
                echo "(va-yikra ELOHIM (la-or YOM) (la-choshekh LAILAH))"
                ;;
            6|7|8)
                echo "(yomer ELOHIM (yehi RAKIA be-tokh HA-MAYIM))"
                echo "(va-yehi KEN)"
                ;;
            9|10)
                echo "(yomer ELOHIM (yikavu HA-MAYIM :to (MAKOM ECHAD)))"
                echo "(va-yikra ELOHIM (la-yabashah ERETZ) (u-le-mikveh YAMIM))"
                ;;
            11|12|13)
                echo "(yomer ELOHIM"
                echo "  (tadshe HA-ERETZ DESHE        ;; vegetate(vegetation)"
                echo "    (ESEV mazria ZERA)          ;; plant.seed(seed)"
                echo "    (ETZ PRI oseh PRI)))        ;; tree.fruit(fruit)"
                echo "(va-yehi KEN)"
                ;;
            *)
                echo ";; See experiments/genesis.dsl for full source"
                ;;
        esac
        echo ""
    fi
done

echo -e "${BOLD}${CYAN}"
echo "════════════════════════════════════════════════════════════════════════════"
echo "  Comparison Complete"
echo "════════════════════════════════════════════════════════════════════════════"
echo -e "${RESET}"

echo -e "${DIM}To run the Genesis executable:${RESET}"
echo "  go run ./experiments/dsl/main.go"
echo ""
echo -e "${DIM}To generate full chapter translations:${RESET}"
echo "  go run ./cmd/pipeline translate-chapter <profile> heb-wlc GEN 1"
echo ""
