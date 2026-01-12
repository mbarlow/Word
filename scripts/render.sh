#!/usr/bin/env bash
set -euo pipefail

# Render canonical dataset to Markdown

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(dirname "$SCRIPT_DIR")"

log() { echo "[$(date '+%Y-%m-%d %H:%M:%S')] $*"; }

log "Running render pipeline..."

cd "$REPO_ROOT"
go run ./cmd/pipeline render

log "Render complete"
