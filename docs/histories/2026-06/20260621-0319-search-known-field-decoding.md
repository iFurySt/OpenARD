# Search Known Field Decoding

## Request

Continue hardening the Go/Cobra/GORM/Gin/Postgres ARD registry and toolkit with real
verification, milestone commits, and alignment with the ARD specification.

## Changes

- Added JSON decode-time known-field checks for `SearchRequest` root fields.
- Added JSON decode-time known-field checks for shared search/explore `query` fields.
- Preserved dynamic filter keys inside `query.filter`.
- Added model tests for unsupported root/query fields and existing scalar/array filters.
- Added Postgres-backed HTTP integration coverage for `/search` returning `400` on
  unsupported root/query fields.
- Updated architecture, quality, and feature release notes.

## Intent

The ARD OpenAPI schema marks `SearchRequest` and `QueryModel` with
`additionalProperties: false`, while `filter` remains intentionally dynamic. The
registry now rejects misspelled request fields such as `federations` and unsupported
query fields such as `sort` before running search or federation work.

## Files

- `internal/ard/models.go`
- `internal/ard/models_test.go`
- `internal/httpapi/router_integration_test.go`
- `docs/ARCHITECTURE.md`
- `docs/QUALITY_SCORE.md`
- `docs/releases/feature-release-notes.md`
