# Public Browse Client

## Request

Continue the Go/Cobra/Gin/GORM/Postgres ARD implementation with real verification,
small milestones, and spec-aligned client/server surfaces.

## Changes

- Added `ardctl browse` for unauthenticated public registry browsing through `GET /agents`.
- Supported remote browse flags for registry URL, kind, filter, order, limit, page token,
  and raw JSON output.
- Covered public browse request construction, CLI output, and limit validation with unit
  tests.
- Extended the real E2E artifact flow to exercise public browse filtering, ordering, and
  pagination against a running registry.
- Updated README, architecture, quality, and release notes so the new client surface is
  visible in-repo.

## Intent

Public `/agents` browsing is useful only if clients can call it without writing custom
HTTP glue. The CLI command mirrors the registry's deterministic browse controls and keeps
admin-only management separate from public discovery.

## Files

- `internal/cli/browse.go`
- `internal/cli/browse_test.go`
- `internal/cli/root.go`
- `internal/cli/root_test.go`
- `scripts/test-e2e-artifacts.sh`
- `README.md`
- `docs/ARCHITECTURE.md`
- `docs/QUALITY_SCORE.md`
- `docs/releases/feature-release-notes.md`
