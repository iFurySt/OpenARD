# Agents Query Validation

## Request

Continue hardening the Go/Cobra/GORM/Gin/Postgres ARD registry and toolkit with real
verification, milestone commits, and alignment with the ARD specification.

## Changes

- Added public `/agents` query parameter validation before store access.
- Rejected non-integer, zero, negative, and over-100 `pageSize` values with `400`.
- Preserved opaque `pageToken` validation through the existing pagination layer.
- Recognized spec-defined `filter` and `orderBy` parameters but returned `400` until
  deterministic filtering and ordering are implemented.
- Rejected unknown public `/agents` query parameters instead of ignoring typos.
- Added Postgres-backed HTTP integration coverage for invalid browse query inputs.
- Updated architecture, quality, and feature release notes.

## Intent

The ARD OpenAPI schema defines public browse query parameters for deterministic listing.
The implementation currently supports pagination only, so unsupported deterministic
filtering and ordering now fail explicitly rather than appearing to work while being
ignored.

## Files

- `internal/httpapi/router.go`
- `internal/httpapi/router_integration_test.go`
- `docs/ARCHITECTURE.md`
- `docs/QUALITY_SCORE.md`
- `docs/releases/feature-release-notes.md`
