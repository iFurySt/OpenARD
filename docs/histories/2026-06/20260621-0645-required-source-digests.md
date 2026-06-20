# Required Source Digests

## Request

Continue hardening the Go/Cobra/GORM/Gin/Postgres ARD registry and toolkit toward a
verified public release with real security checks.

## Changes

- Added `verify.SourceDigestOptions` and `VerifySourceDigestsWithOptions`.
- Added `ard verify catalog --require-source-digests`.
- Made strict mode require every URL-delivered catalog entry to carry
  `trustManifest.sourceDigest` and verify all pinned digests.
- Kept embedded `data` entries exempt because they do not have a retrievable source URL.
- Added unit, CLI, and E2E coverage for strict source digest requirements.
- Updated README, trust/security docs, architecture, quality notes, and release notes.

## Design Intent

The existing `--source-digests` mode verifies pinned artifacts but intentionally skips
URL entries that have no pin. Enterprise adoption needs a stricter gate for curated
catalogs so unpinned remote artifacts cannot pass verification unnoticed.

## Main Files

- `internal/verify/source_digest.go`
- `internal/cli/verify.go`
- `internal/cli/verify_test.go`
- `internal/verify/source_digest_test.go`
- `scripts/test-e2e-artifacts.sh`
- `docs/TRUST.md`
- `docs/SECURITY.md`
