## [2026-06-20 23:50] | Task: Source Digest Verification

### User Request

> Continue building the Go/Cobra/GORM/Gin/Postgres ARD registry toward a B2B neutral
> registry and management toolkit with real verification and milestone commits.

### Changes

- Added `--pin-source-digest` to local and remote artifact add commands.
- Added `trustManifest.identity` and `trustManifest.sourceDigest` when pinning URL
  artifacts.
- Added URL artifact fetch retry in the adapter source loader.
- Added `ard verify catalog --source-digests` to fetch URL artifacts and verify pinned
  sha256 digests.
- Added minimal `trustManifest` validation for identity and `sourceDigest` format.
- Added unit coverage for adapter pinning, trust validation, and source digest
  verification.
- Extended E2E to pin and verify a real Agentmemory MCP server card URL.
- Added `docs/TRUST.md` and updated README, architecture, security, and quality docs.

### Design Notes

This is deliberately narrower than full trust verification. Source digest verification
proves URL artifact byte integrity at verification time. It does not prove publisher
identity, validate detached signatures, resolve DIDs, verify SPIFFE identities, or assess
runtime safety. The narrower implementation gives enterprises a useful first integrity
control while keeping future signature and identity verification paths open.

### Verification

- Passed: `make fmt-check`
- Passed: `make test`
- Passed: `make build`
- Passed: `make test-integration`
- Passed: `make test-e2e`
- `make test-e2e` pinned and verified `trustManifest.sourceDigest` for the real
  Agentmemory MCP server card URL.
- Upstream `ard-spec` manifest and registry conformance passed. The manifest check still
  reports the expected OpenAPI extension media type warning.
