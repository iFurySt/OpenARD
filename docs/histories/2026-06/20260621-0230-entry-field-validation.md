# Entry Field Validation

## Request

Continue building the Go/Cobra/GORM/Gin/Postgres ARD registry and toolkit with
spec-aligned validation, real verification, and milestone commits.

## Changes

- Added catalog entry `updatedAt` validation using RFC3339 date-time parsing.
- Added catalog entry `metadata` value validation so extension values must be strings,
  numbers, booleans, or null.
- Accepted Go-native numeric values and JSON-decoded numeric values for metadata.
- Added focused model tests and adapter tests to ensure existing generated metadata still
  validates.
- Updated security, trust, architecture, quality, and release documentation.

## Design Notes

This is ARD schema alignment and input hygiene. It rejects malformed timestamps and
complex metadata extension values before persistence while leaving embedded artifact
`data` unconstrained for protocol-specific payloads.

## Files

- `internal/ard/models.go`
- `internal/ard/models_test.go`
- `docs/TRUST.md`
- `docs/SECURITY.md`
- `docs/ARCHITECTURE.md`
- `docs/QUALITY_SCORE.md`
- `docs/releases/feature-release-notes.md`
