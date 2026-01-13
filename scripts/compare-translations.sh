#!/bin/bash
# Compare translations across models
# Usage: ./scripts/compare-translations.sh [BOOK] [CHAPTER] [VERSE]
#        ./scripts/compare-translations.sh GEN 1 1      # Single verse
#        ./scripts/compare-translations.sh GEN 1        # Full chapter
#        ./scripts/compare-translations.sh              # Genesis 1:1-11

set -e

BOOK="${1:-GEN}"
CHAPTER="${2:-1}"
VERSE="${3:-}"

# Colors for better readability
BLUE='\033[0;34m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
MAGENTA='\033[0;35m'
BOLD='\033[1m'
DIM='\033[2m'
NC='\033[0m' # No Color

# Paths
DRAFTS_DIR="translations/drafts"
DB_PATH="data/word.db"

# Find all model directories for a profile
get_models() {
    local profile="$1"
    if [ -d "$DRAFTS_DIR/t_$profile" ]; then
        ls -1 "$DRAFTS_DIR/t_$profile" 2>/dev/null | grep -v '\.jsonl'
    fi
}

# Get verse from JSONL file
get_translation() {
    local profile="$1"
    local model="$2"
    local book="$3"
    local chapter="$4"
    local verse="$5"
    local file="$DRAFTS_DIR/t_$profile/$model/$book/$chapter.jsonl"

    if [ -f "$file" ]; then
        jq -r "select(.id == \"t_${profile}/${book}/${chapter}/${verse}\") | .text" "$file" 2>/dev/null
    fi
}

# Get notes from JSONL file
get_notes() {
    local profile="$1"
    local model="$2"
    local book="$3"
    local chapter="$4"
    local verse="$5"
    local file="$DRAFTS_DIR/t_$profile/$model/$book/$chapter.jsonl"

    if [ -f "$file" ]; then
        jq -r "select(.id == \"t_${profile}/${book}/${chapter}/${verse}\") | .notes // [] | .[]" "$file" 2>/dev/null
    fi
}

# Get flags from critic file
get_flags() {
    local profile="$1"
    local model="$2"
    local book="$3"
    local chapter="$4"
    local verse="$5"
    local file="$DRAFTS_DIR/t_$profile/$model/$book/$chapter.critic.jsonl"

    if [ -f "$file" ]; then
        jq -r "select(.id == \"t_${profile}/${book}/${chapter}/${verse}\") | .flags[] | \"[\(.type)] \(.detail)\"" "$file" 2>/dev/null
    fi
}

# Get source text from database
get_source() {
    local work="$1"
    local book="$2"
    local chapter="$3"
    local verse="$4"
    local vid="$work/$book/$chapter/$verse"

    if [ -f "$DB_PATH" ]; then
        sqlite3 "$DB_PATH" "SELECT text FROM verses WHERE id = '$vid'" 2>/dev/null
    fi
}

# Print a single verse comparison
print_verse() {
    local book="$1"
    local chapter="$2"
    local verse="$3"

    echo ""
    echo -e "${BOLD}═══════════════════════════════════════════════════════════════${NC}"
    echo -e "${BOLD}  $book $chapter:$verse${NC}"
    echo -e "${BOLD}═══════════════════════════════════════════════════════════════${NC}"

    # Hebrew source
    local hebrew=$(get_source "heb-wlc" "$book" "$chapter" "$verse")
    if [ -n "$hebrew" ]; then
        echo -e "\n${MAGENTA}[Hebrew WLC]${NC}"
        echo -e "${DIM}$hebrew${NC}"
    fi

    # KJV reference
    local kjv=$(get_source "kjv" "$book" "$chapter" "$verse")
    if [ -n "$kjv" ]; then
        echo -e "\n${BLUE}[KJV Reference]${NC}"
        echo "$kjv"
    fi

    # Get all models for techdoc_en profile
    local models=$(get_models "techdoc_en")

    if [ -z "$models" ]; then
        echo -e "\n${YELLOW}No translations found. Run translate-chapter first.${NC}"
        return
    fi

    # Show each model's translation
    for model in $models; do
        local text=$(get_translation "techdoc_en" "$model" "$book" "$chapter" "$verse")
        if [ -n "$text" ]; then
            echo -e "\n${GREEN}[$model]${NC}"
            echo "$text"

            # Show notes if any
            local notes=$(get_notes "techdoc_en" "$model" "$book" "$chapter" "$verse")
            if [ -n "$notes" ]; then
                echo -e "${DIM}  Notes:${NC}"
                echo "$notes" | while read -r note; do
                    echo -e "${DIM}    - $note${NC}"
                done
            fi

            # Show flags if any
            local flags=$(get_flags "techdoc_en" "$model" "$book" "$chapter" "$verse")
            if [ -n "$flags" ]; then
                echo -e "${YELLOW}  Flags:${NC}"
                echo "$flags" | while read -r flag; do
                    echo -e "${YELLOW}    $flag${NC}"
                done
            fi
        fi
    done
}

# Main
echo -e "${BOLD}Translation Comparison Tool${NC}"
echo -e "${DIM}Profile: techdoc_en${NC}"

# Determine verse range
if [ -n "$VERSE" ]; then
    # Single verse
    print_verse "$BOOK" "$CHAPTER" "$VERSE"
elif [ "$BOOK" = "GEN" ] && [ "$CHAPTER" = "1" ] && [ -z "$1" ]; then
    # Default: Genesis 1:1-11
    echo -e "${DIM}Showing: $BOOK $CHAPTER:1-11 (default)${NC}"
    for v in $(seq 1 11); do
        print_verse "$BOOK" "$CHAPTER" "$v"
    done
else
    # Full chapter - detect verse count from database
    if [ -f "$DB_PATH" ]; then
        max_verse=$(sqlite3 "$DB_PATH" "SELECT MAX(verse) FROM verses WHERE osis = '$BOOK' AND chapter = $CHAPTER" 2>/dev/null)
        if [ -n "$max_verse" ]; then
            echo -e "${DIM}Showing: $BOOK $CHAPTER:1-$max_verse${NC}"
            for v in $(seq 1 "$max_verse"); do
                print_verse "$BOOK" "$CHAPTER" "$v"
            done
        else
            echo "No verses found for $BOOK $CHAPTER"
        fi
    else
        echo "Database not found at $DB_PATH"
    fi
fi

echo ""
echo -e "${DIM}───────────────────────────────────────────────────────────────${NC}"
echo -e "${DIM}Models: $(get_models techdoc_en | tr '\n' ' ')${NC}"
