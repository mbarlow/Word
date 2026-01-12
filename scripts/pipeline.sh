#!/usr/bin/env bash
set -euo pipefail

# Run full Word pipeline: ingest → normalize → render → validate

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(dirname "$SCRIPT_DIR")"

log() { echo "[$(date '+%Y-%m-%d %H:%M:%S')] $*"; }

log "Starting Word pipeline..."

log "Stage 1: Ingest"
"$SCRIPT_DIR/ingest.sh"

log "Stage 2: Render"
"$SCRIPT_DIR/render.sh"

log "Stage 3: Validate"
"$SCRIPT_DIR/validate.sh"

log "Pipeline complete"
