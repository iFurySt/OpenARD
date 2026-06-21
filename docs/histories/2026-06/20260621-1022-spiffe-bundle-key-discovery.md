## [2026-06-21 10:22] | Task: SPIFFE bundle key discovery

### Execution Context

- Agent ID: `codex`
- Base Model: `gpt-5`
- Runtime: `local CLI`

### User Query

> Continue building the Go/Cobra/GORM/Gin/Postgres ARD registry and toolkit with real
> verification, milestone commits, and strict alignment with the self-hosted enterprise
> ARD objective.

### Changes Overview

- Area: Security
- Key actions:
  - Added opt-in `ard verify catalog --jws-discover-spiffe`.
  - Discovered OKP/Ed25519 keys from `SPIFFE-X509` attestation HTTPS JWKS URIs.
  - Required SPIFFE bundle URI hosts to match the `spiffe://` trust domain.
  - Reused the existing detached compact JWS verification path.
  - Updated CLI public surface checks and trust documentation.

### Design Intent

ARD spec examples publish SPIFFE bundle material through
`trustManifest.attestations[].uri`, so the implementation uses that existing field
instead of adding an implementation-specific catalog extension. Discovery remains
explicitly opt-in and only proves that a host-matched bundle advertised the signing key
at verification time. It does not validate X.509-SVID chains, workload identity runtime
state, revocation, federation bundle sequence, or claim truth.

### Files Modified

- `internal/verify/signature.go`
- `internal/verify/signature_test.go`
- `internal/cli/verify.go`
- `internal/cli/verify_test.go`
- `internal/tools/publicsurface/main.go`
- `README.md`
- `docs/TRUST.md`
- `docs/SECURITY.md`
- `docs/ARCHITECTURE.md`
- `docs/QUALITY_SCORE.md`
- `docs/releases/feature-release-notes.md`
