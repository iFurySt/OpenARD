# Agents Filter And Order

## Request

Continue hardening the Go/Cobra/GORM/Gin/Postgres ARD registry and toolkit with real
verification, milestone commits, and alignment with the ARD specification.

## Changes

- Implemented public `/agents` deterministic filter parsing for common spec fields:
  `displayName`, `type`, `publisherId`, `createdAfter`, and `updatedAfter`.
- Implemented whitelisted `orderBy` parsing for display name, type, publisher, creation
  time, and update time.
- Extended Postgres list queries with parameterized filter predicates and safe order
  column mapping.
- Kept unsupported filter and order fields as `400 INVALID_ARGUMENT` responses.
- Added parser unit coverage and Postgres-backed HTTP integration coverage for filtered
  and ordered browse requests.
- Updated architecture, quality, and feature release notes.

## Intent

The ARD List API is intended for deterministic browsing in developer portals. The
registry now supports a practical subset of the spec-defined filter and order surface
without accepting arbitrary SQL or silently ignoring unsupported fields.

## Files

- `internal/httpapi/router.go`
- `internal/httpapi/router_test.go`
- `internal/httpapi/router_integration_test.go`
- `internal/store/postgres.go`
- `docs/ARCHITECTURE.md`
- `docs/QUALITY_SCORE.md`
- `docs/releases/feature-release-notes.md`
