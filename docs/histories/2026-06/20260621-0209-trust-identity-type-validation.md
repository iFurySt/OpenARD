# Trust Identity Type Validation

## Request

Continue building the Go/Cobra/GORM/Gin/Postgres ARD registry and toolkit with
spec-aligned validation, real verification, and milestone commits.

## Changes

- Added `trustManifest.identityType` validation against the ARD schema values:
  `spiffe`, `did`, `https`, and `other`.
- Rejected non-string `trustManifest.identityType` values.
- Rejected non-string `trustManifest.sourceDigest` values so malformed metadata cannot
  bypass existing digest format checks.
- Updated trust, security, architecture, quality, and release documentation.

## Design Notes

This is a schema-alignment guard, not a new identity proof mechanism. The implementation
still leaves DID, SPIFFE, certificate, key-resolution, and detached JWS verification to
future trust verification work.

## Files

- `internal/ard/models.go`
- `internal/ard/models_test.go`
- `docs/TRUST.md`
- `docs/SECURITY.md`
- `docs/ARCHITECTURE.md`
- `docs/QUALITY_SCORE.md`
- `docs/releases/feature-release-notes.md`
