# Explore Request Validation

## Request

Continue hardening the Go/Cobra/GORM/Gin/Postgres ARD registry and toolkit with real
verification, milestone commits, and alignment with the ARD specification.

## Changes

- Added JSON decode-time known-field checks for `ExploreRequest` root fields.
- Added JSON decode-time known-field checks for `ExploreFacetRequest` fields.
- Reused shared `query` known-field validation for explore queries while preserving
  dynamic filter keys.
- Added shared explore request validation for required facets, required facet fields,
  and non-negative facet `limit` / `minCount` values.
- Updated `/explore` to return `400` through the shared validator before store work.
- Added model tests and Postgres-backed HTTP integration coverage for invalid explore
  request shapes.
- Updated architecture, quality, and feature release notes.

## Intent

The ARD OpenAPI schema marks `ExploreRequest`, `QueryModel`, and `ExploreFacetRequest`
as closed shapes while leaving filter keys dynamic. The registry now rejects mistyped
explore fields such as root `sort`, query `sort`, and facet `order` before running
facet aggregation.

## Files

- `internal/ard/models.go`
- `internal/ard/models_test.go`
- `internal/httpapi/router.go`
- `internal/httpapi/router_integration_test.go`
- `docs/ARCHITECTURE.md`
- `docs/QUALITY_SCORE.md`
- `docs/releases/feature-release-notes.md`
