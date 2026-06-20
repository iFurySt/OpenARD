## [2026-06-21 04:52] | Task: Richer list filter operators

### Execution Context

- Agent ID: `codex`
- Base Model: `GPT-5`
- Runtime: `local CLI`

### User Query

> Continue building the Go/Cobra/GORM/Gin/Postgres ARD registry and toolkit with real verification, milestone commits, and real tests for MCP/Skill/A2A-related flows.

### Changes Overview

- Area: public browse and local registry inventory filters
- Key actions:
  - Added parsed filter clauses for `!=`, `contains`, and `>=` while preserving existing field-oriented `ListFilter` usage.
  - Extended Postgres list filtering for text fields, publisher IDs, JSON array fields, metadata scalar fields, and timestamp boundaries.
  - Added parser, Postgres integration, HTTP integration, and real E2E checks for the richer operators.

### Design Intent

B2B registry operators need inventory queries that can exclude resource classes and
search within names, publishers, tags, capabilities, and metadata without exporting the
catalog. The implementation keeps the existing simple `AND` grammar and adds operators
that map cleanly to indexed or bounded Postgres predicates; grouped boolean expressions
remain a separate future design.

### Files Modified

- `internal/store/list_query.go`
- `internal/store/postgres.go`
- `internal/store/list_query_test.go`
- `internal/store/postgres_integration_test.go`
- `internal/httpapi/router_integration_test.go`
- `scripts/test-e2e-artifacts.sh`
- `README.md`
- `docs/ARCHITECTURE.md`
- `docs/QUALITY_SCORE.md`
- `docs/releases/feature-release-notes.md`
