#!/usr/bin/env bash
set -euo pipefail

# Download upstream source texts for Word pipeline
# Fetches: KJV (eBible.org), WLC (OSHB), TR1894 (Byzantine Text)

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(dirname "$SCRIPT_DIR")"
DATA_DIR="${WORD_DATA_DIR:-$REPO_ROOT/data}"
SOURCE_DIR="$DATA_DIR/source"

log() { echo "[$(date '+%Y-%m-%d %H:%M:%S')] $*"; }
error() { log "ERROR: $*" >&2; exit 2; }

mkdir -p "$SOURCE_DIR"/{kjv_raw,wlc_raw,tr1894_raw}

# --- KJV from eBible.org ---
log "Downloading KJV from eBible.org..."
KJV_URL="https://ebible.org/Scriptures/eng-kjv2006_usfm.zip"
KJV_ZIP="$SOURCE_DIR/kjv_raw/eng-kjv2006_usfm.zip"

if [[ ! -f "$KJV_ZIP" ]]; then
    curl -fsSL -o "$KJV_ZIP" "$KJV_URL" || error "Failed to download KJV"
    log "Extracting KJV USFM files..."
    unzip -q -o "$KJV_ZIP" -d "$SOURCE_DIR/kjv_raw/"
else
    log "KJV already downloaded, skipping"
fi

# --- WLC from OSHB ---
log "Downloading WLC from Open Scriptures Hebrew Bible..."
WLC_REPO="https://github.com/openscriptures/morphhb"
WLC_DIR="$SOURCE_DIR/wlc_raw/morphhb"

if [[ ! -d "$WLC_DIR" ]]; then
    git clone --depth 1 "$WLC_REPO" "$WLC_DIR" || error "Failed to clone OSHB"
else
    log "WLC already downloaded, skipping"
fi

# --- TR1894 from Byzantine Text ---
log "Downloading TR1894 from Byzantine Text Project..."
TR_REPO="https://github.com/byztxt/greektext-scrivener"
TR_DIR="$SOURCE_DIR/tr1894_raw/greektext-scrivener"

if [[ ! -d "$TR_DIR" ]]; then
    git clone --depth 1 "$TR_REPO" "$TR_DIR" || error "Failed to clone TR1894"
else
    log "TR1894 already downloaded, skipping"
fi

log "All sources downloaded to $SOURCE_DIR"
log "KJV:    $SOURCE_DIR/kjv_raw/"
log "WLC:    $SOURCE_DIR/wlc_raw/morphhb/wlc/"
log "TR1894: $SOURCE_DIR/tr1894_raw/greektext-scrivener/textonly/"
