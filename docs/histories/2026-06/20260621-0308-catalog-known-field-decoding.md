# Catalog Known Field Decoding

## Request

Continue hardening the Go/Cobra/GORM/Gin/Postgres ARD registry and toolkit with real
verification, milestone commits, and alignment with the ARD specification.

## Changes

- Added JSON decode-time known-field checks for catalog root fields.
- Added JSON decode-time known-field checks for catalog `host` fields.
- Added tests that reject deprecated root `collections` and unsupported host fields.
- Added a regression test that catalog entry extension fields are not rejected by this
  root/host closed-shape check.
- Updated architecture, security, quality, and feature release notes.

## Intent

The upstream ARD catalog schema marks the root manifest object and `host` object with
`additionalProperties: false`, and the conformance guide explicitly flags the deprecated
root `collections` field. Go struct decoding previously ignored these fields silently.
The decoder now fails early before validation and persistence, while leaving catalog entry
extension fields available for a future preservation/filtering design.

## Files

- `internal/ard/models.go`
- `internal/ard/models_test.go`
- `docs/ARCHITECTURE.md`
- `docs/SECURITY.md`
- `docs/QUALITY_SCORE.md`
- `docs/releases/feature-release-notes.md`
