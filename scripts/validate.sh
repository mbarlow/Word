#!/usr/bin/env bash
set -euo pipefail

# Validate rendered output against canonical expectations

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(dirname "$SCRIPT_DIR")"

log() { echo "[$(date '+%Y-%m-%d %H:%M:%S')] $*"; }

log "Running validation..."

cd "$REPO_ROOT"
go run ./cmd/pipeline validate

log "Validation complete"
