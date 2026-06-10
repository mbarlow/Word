# Running Word on seed

How the Word API is deployed on seed and wired into the shared edge. The
container has been running this way since June 2026; this doc makes the path
reproducible. See `~/git/github.com/mbarlow/CLAUDE.md` (workspace root) for
the wider fern/seed picture.

## The shape

- **One container**, `word-service`, built from `deploy/Dockerfile` by
  `compose.seed.yml`. No host port — Caddy reaches `word-service:8080` over
  the shared `fern-pi-hole_default` network as `http://word.lan`.
- **State is bind-mounted from the repo**: `./data` (the SQLite DB) and
  `./texts` (source texts). The DB outlives rebuilds because it lives in the
  working tree, not the image or a named volume.
- **Ollama stays on fern** — translation work needs the NVIDIA GPU. When fern
  is down, `llama.lan` (seed-llm) is the OpenAI-compat fallback.

## Bring-up (fresh box)

```bash
cd ~/git/github.com/mbarlow/Word

# 1. Source texts + DB. Either restore data/word.db from seed-backups
#    (it's in the nightly restic set), or rebuild from scratch:
./scripts/download-sources.sh     # fetch kjv / heb-wlc / grc-tr1894 into texts/
./scripts/pipeline.sh             # ingest -> validate -> render into data/

# 2. Run it
docker compose -f compose.seed.yml up -d --build

# 3. Verify
curl -s word-service:8080/health        # from another container, or:
curl -s -H 'Host: word.lan' http://127.0.0.1/v1/works
```

## Edge wiring (already in place)

The standard seed recipe — only needed again on a from-scratch rebuild:

1. `seed-caddy/Caddyfile`: `http://word.lan { reverse_proxy word-service:8080 }`
2. Pi-hole `dns.hosts`: `word.lan` → `192.168.0.119`
3. `docker compose restart caddy` (not `caddy reload` — bind-mount inode gotcha)

## Already covered elsewhere

- **Backups**: `data/word.db` is in seed-backups' nightly restic set
  (sqlite `.backup` hot copy — no downtime).
- **Metrics**: `/metrics` (echoprometheus) is scraped by Prometheus on the
  shared network (`seed-observability/prometheus/prometheus.yml`, job
  `word-service`). Logs are picked up by Promtail like every container.
- **Health**: compose healthcheck curls `/health`; `docker ps` shows it.

## Upgrading

```bash
cd ~/git/github.com/mbarlow/Word
git pull
docker compose -f compose.seed.yml up -d --build   # DB persists (bind mount)
```

Dev on this repo happens via Tilt (`tilt up`, see README) — that's a separate
local loop and does not touch the seed container.
