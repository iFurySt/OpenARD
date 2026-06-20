## [2026-06-20 21:20] | Task: Verify And Crawl CLI

### Execution Context

- Agent ID: `codex`
- Base Model: `GPT-5`
- Runtime: `local workspace`

### User Query

> Continue implementing the Go ARD toolkit with real verification, small milestones, and
> concrete validation.

### Changes Overview

- Area: CLI and catalog management.
- Key actions:
  - Added `ard verify catalog SOURCE`.
  - Added `ard crawl URL` for well-known `/.well-known/ai-catalog.json` discovery and
    import.
  - Added catalog URL discovery helper and tests.
  - Updated README and architecture docs.

### Design Intent

This makes the CLI more like a management toolkit instead of only an import/search demo.
`verify catalog` gives users a quick local/remote validation path, while `crawl` proves
the publisher-side ARD discovery flow.

### Verification

- `go test ./...`
- `make test-integration`
- `make build`
- Started a local HTTP site serving `/.well-known/ai-catalog.json`.
- Ran `bin/ard verify catalog http://127.0.0.1:18081/.well-known/ai-catalog.json --json`
  and received a valid result with 2 entries.
- Ran `bin/ard crawl http://127.0.0.1:18081/` against a temporary Postgres 16 Docker
  container.
- Started the registry server and confirmed `ard search "weather forecast" --kind mcp`
  returned the crawled `Weather Data Node`.

### Files Modified

- `README.md`
- `docs/ARCHITECTURE.md`
- `docs/histories/2026-06/20260620-2120-verify-crawl-cli.md`
- `internal/catalog/discovery.go`
- `internal/catalog/loader_test.go`
- `internal/cli/crawl.go`
- `internal/cli/root.go`
- `internal/cli/verify.go`
