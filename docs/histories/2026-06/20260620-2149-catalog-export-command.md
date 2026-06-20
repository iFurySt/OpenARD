## [2026-06-20 21:49] | Task: Add Catalog Export Command

### User Request

> Continue building the Go/Cobra/GORM/Gin/Postgres ARD registry and management toolkit
> with spec alignment, real verification, and milestone commits.

### Changes

- Added `ard export catalog` / `ardctl export catalog`.
- Added `Store.ExportCatalog` to read all persisted entries in deterministic order.
- Exported catalogs use `specVersion: "1.0"` and validate through the same ARD catalog
  validator before writing.
- Added host metadata flags:
  - `--host-display-name`
  - `--host-identifier`
  - `--documentation-url`
- Added `--output` / `-o`, with stdout as the default.
- Updated README, architecture notes, quality score, and integration coverage.

### Design Notes

Catalog export closes the ingestion loop: a registry can now import catalogs and artifacts
into Postgres, then emit a portable `ai-catalog.json` for backup, migration, review, or
publication behind `/.well-known/ai-catalog.json`.

### Verification

- Passed: `make fmt`
- Passed: `go test ./...`
- Passed: `make build`
- Passed: `make test-integration`
- Passed: Postgres E2E export flow:
  - `ardctl add catalog` imported the checked-in ARD catalog fixture.
  - `ardctl add mcp` imported the checked-in MCP server card fixture.
  - `ardctl export catalog -o <file>` wrote a catalog containing all imported entries.
  - `ard verify catalog <file> --json` reported the exported manifest as valid.
  - Upstream `ard-spec` `conformance-test manifest <file>` passed with 0 critical
    specification errors and 0 warnings.
