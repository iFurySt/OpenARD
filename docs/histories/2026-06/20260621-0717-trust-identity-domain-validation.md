# Trust Identity Domain Validation

## Request

Continue building the Go/Cobra/Gin/GORM/Postgres ARD registry and toolkit with
spec-aligned trust validation, real tests, and milestone commits.

## Changes

- Added `trustManifest.identityType` shape checks for `https`, `spiffe`, and `did`.
- Extended publisher-domain alignment from HTTP(S) identities to SPIFFE trust domains and
  `did:web` domains.
- Preserved non-web DID identities for future resolvers while still checking basic DID
  syntax when `identityType` is `did`.
- Added focused catalog model tests for matching and mismatched HTTPS, SPIFFE, and
  `did:web` identities.
- Updated trust, security, architecture, quality, and release documentation.

## Design Intent

The ARD spec says the `urn:air:` publisher acts as the organizational trust anchor and
must align with the cryptographic workload identity in `trustManifest`. This change makes
that alignment explicit for identity formats whose trust domain can be parsed locally.

This is still metadata consistency, not full identity proof. It does not resolve DID
documents, validate SPIFFE SVIDs, verify certificates, or prove control of the publisher
domain.

## Files

- `internal/ard/models.go`
- `internal/ard/models_test.go`
- `docs/TRUST.md`
- `docs/SECURITY.md`
- `docs/ARCHITECTURE.md`
- `docs/QUALITY_SCORE.md`
- `docs/releases/feature-release-notes.md`
