# ADR-0001: Go with Echo Framework and GORM

## Status
Accepted

## Date
2026-01-11

## Context
We need to choose a technology stack for the Word API server. Requirements include:

- Production-grade HTTP API
- Database access with ORM capabilities
- Prometheus metrics instrumentation
- JSON structured logging
- Docker/Tilt-based local development
- Long-term maintainability

## Decision
We will use:

- **Go** as the primary language
- **Echo** as the HTTP framework
- **GORM** as the ORM
- **Prometheus client** for metrics via Echo middleware

## Rationale

### Go
- Single binary deployment simplifies Docker images
- Excellent concurrency model for API workloads
- Strong standard library for text processing (critical for pipeline stages)
- Static typing catches errors at compile time

### Echo
- Minimal, high-performance HTTP framework
- Built-in middleware ecosystem (logging, recovery, CORS)
- First-class Prometheus instrumentation via `echo-contrib`
- Clean routing API

### GORM
- Mature, well-documented ORM
- Supports SQLite (Phase 1) and PostgreSQL (Phase 2) with same model definitions
- Auto-migration for development convenience
- Raw SQL escape hatch when needed

## Consequences

### Positive
- Fast compilation and startup times
- Low memory footprint suitable for local development
- Clear separation between API layer and pipeline tooling
- Easy to add CLI commands for pipeline stages

### Negative
- Go's error handling is verbose compared to exceptions
- GORM's magic can obscure query behavior (mitigated by logging queries in dev)

## Alternatives Considered

### Python + FastAPI
- Faster prototyping but slower runtime
- Better for ML/LLM integration but we can call Go from Python or vice versa
- Rejected: performance and deployment simplicity favor Go

### Rust + Axum
- Maximum performance but slower development velocity
- Rejected: overkill for this workload; Go is sufficient

## References
- Echo: https://echo.labstack.com/
- GORM: https://gorm.io/
- Go Project Layout: https://github.com/golang-standards/project-layout
