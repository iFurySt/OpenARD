# DID Web Key Discovery

## User Request

Continue building the Go/Cobra/GORM/Gin/Postgres ARD registry and toolkit with real
verification and milestone commits.

## Changes

- Added opt-in `ard verify catalog --jws-discover-did-web`.
- Added DID web document resolution from `trustManifest.identity` values shaped as
  `did:web`.
- Extracted `verificationMethod[].publicKeyJwk` keys from DID documents and accepted
  only OKP/Ed25519 keys for detached JWS verification.
- Required DID document `id` and verification method `controller` to match the entry
  `trustManifest.identity` when those fields are present.
- Kept automatic discovery disabled by default and scoped to `did:web`; SPIFFE,
  certificate, OIDC, and non-`did:web` DID discovery remain out of scope.
- Updated README, trust, security, architecture, quality, and release notes.

## Design Intent

`did:web` is the first practical automatic key-discovery slice because the identity
itself deterministically maps to an HTTPS DID document. The implementation reuses the
existing Ed25519/JWKS detached JWS verifier instead of adding a parallel signature
path, preserving the same canonical `trustManifest` payload semantics and trust-domain
guardrails.

## Important Files

- `internal/verify/signature.go`
- `internal/verify/signature_test.go`
- `internal/cli/verify.go`
- `internal/cli/verify_test.go`
- `docs/TRUST.md`
- `docs/SECURITY.md`
