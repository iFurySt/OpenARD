# JWKS Trust Anchors

## Request

Continue building the Go/Cobra/Gin/GORM/Postgres ARD registry and toolkit with
spec-aligned trust verification, real tests, and milestone commits.

## Changes

- Extended `ard verify catalog --jws-trust-anchors` to accept local JWKS files.
- Added support for JWKS OKP/Ed25519 public keys using `kty`, `crv`, `kid`, `alg`, and
  `x`.
- Preserved the existing ard-native `publicKey` trust-anchor format.
- Added unit coverage for JWKS signature verification and unsupported JWKS curve
  rejection.
- Updated the CLI signature verification test to use a JWKS trust-anchor file.
- Updated trust, security, architecture, quality, and release documentation.

## Design Intent

Enterprise key material is often distributed as JWKS. Supporting local JWKS files makes
signed trust-manifest verification easier to adopt without introducing remote key
discovery or identity-provider coupling.

This is still explicit operator trust-anchor verification. `ard` does not fetch remote
JWKS documents, run OIDC discovery, prove key ownership, or validate signed claim truth.

## Files

- `internal/verify/signature.go`
- `internal/verify/signature_test.go`
- `internal/cli/verify_test.go`
- `docs/TRUST.md`
- `docs/SECURITY.md`
- `docs/ARCHITECTURE.md`
- `docs/QUALITY_SCORE.md`
- `docs/releases/feature-release-notes.md`
- `README.md`
