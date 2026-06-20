# Well-Known Active Catalog

## Request

Continue the Go/Cobra/Gin/GORM/Postgres ARD implementation with real verification,
small milestones, and spec-aligned registry publication behavior.

## Changes

- Changed `GET /.well-known/ai-catalog.json` to publish the registry self entry plus
  active persisted catalog entries.
- Kept pending and disabled entries out of the public well-known catalog by reusing the
  same active export path as public discovery.
- Added Postgres-backed HTTP integration coverage proving active entries are published,
  pending entries are hidden, and the response validates as an ARD catalog.
- Extended real E2E and compose verification to fetch the well-known catalog after
  imports and assert it contains active entries.
- Updated README, architecture, quality, and release notes.

## Intent

Self-hosted registries should be crawlable through the ARD well-known catalog location,
not just searchable through `/search`. Publishing active entries there makes the registry
usable as a standards-shaped catalog source while preserving lifecycle governance.

## Files

- `internal/httpapi/router.go`
- `internal/httpapi/router_integration_test.go`
- `scripts/test-e2e-artifacts.sh`
- `scripts/test-compose.sh`
- `README.md`
- `docs/ARCHITECTURE.md`
- `docs/QUALITY_SCORE.md`
- `docs/releases/feature-release-notes.md`
