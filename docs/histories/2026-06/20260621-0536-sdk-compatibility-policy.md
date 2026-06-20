## [2026-06-21 05:36] | Task: Go SDK Compatibility Policy

### Execution Context

- Agent ID: `codex`
- Base Model: `GPT-5`
- Runtime: `Codex CLI`

### User Query

> Continue building the Go/Cobra/GORM/Gin/Postgres ARD registry and toolkit with real
> verification, milestone commits, and strict alignment with the ARD self-hosted
> registry direction.

### Changes Overview

- Area: public Go SDK, docs, repository routing.
- Key actions:
  - Added `docs/SDK_COMPATIBILITY.md` for public SDK import paths, pre-1.0 expectations,
    compatibility boundaries, and validation gates.
  - Added package documentation for `pkg/ard` and `pkg/client`.
  - Linked the policy from README, architecture docs, AGENTS routing, release notes, and
    quality scoring.

### Design Intent

External enterprise consumers need to know which APIs they can depend on before adopting
the SDK. The policy keeps `pkg/ard` and `pkg/client` as the only public Go import paths,
keeps `internal/` explicitly unstable, and preserves room for pre-1.0 ARD draft changes
while documenting expectations before public tags.

### Files Modified

- `docs/SDK_COMPATIBILITY.md`
- `pkg/ard/doc.go`
- `pkg/ard/models.go`
- `pkg/client/doc.go`
- `pkg/client/client.go`
- `README.md`
- `AGENTS.md`
- `docs/ARCHITECTURE.md`
- `docs/QUALITY_SCORE.md`
- `docs/releases/feature-release-notes.md`
