#!/usr/bin/env bash
set -euo pipefail

# Build meaning graphs for specified book(s) or all OT books
# Usage:
#   ./scripts/meaning.sh ECC          # Single book
#   ./scripts/meaning.sh ECC JER SNG  # Multiple books
#   ./scripts/meaning.sh --all        # All 39 OT books (Hebrew only)

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
ROOT="$(dirname "$SCRIPT_DIR")"
cd "$ROOT"

ALL_OT="GEN EXO LEV NUM DEU JOS JDG RUT 1SA 2SA 1KI 2KI 1CH 2CH EZR NEH EST JOB PSA PRO ECC SNG ISA JER LAM EZK DAN HOS JOL AMO OBA JON MIC NAM HAB ZEP HAG ZEC MAL"

if [ $# -eq 0 ]; then
    echo "Usage: $0 [--all | BOOK ...]"
    echo "Examples:"
    echo "  $0 ECC              # Ecclesiastes"
    echo "  $0 ECC JER SNG      # Multiple books"
    echo "  $0 --all            # All 39 OT books"
    exit 1
fi

books=""
for arg in "$@"; do
    case "$arg" in
        --all) books="$ALL_OT" ;;
        *)     books="$books $arg" ;;
    esac
done

for book in $books; do
    echo ""
    echo "=== $book ==="
    go run ./cmd/pipeline meaning "$book" 2>&1 >/dev/null
done

echo ""
echo "Done. Output in data/meaning/heb-wlc/"
