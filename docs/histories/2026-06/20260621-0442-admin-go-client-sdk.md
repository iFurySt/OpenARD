## [2026-06-21 04:42] | Task: Admin Go client SDK

### Execution Context

- Agent ID: `codex`
- Base Model: `GPT-5`
- Runtime: `local CLI`

### User Query

> Continue building the Go/Cobra/GORM/Gin/Postgres ARD registry and toolkit with real verification, milestone commits, and real tests for MCP/Skill/A2A-related flows.

### Changes Overview

- Area: public Go SDK and E2E verification
- Key actions:
  - Added token-protected admin client methods to `pkg/client`.
  - Covered entry list/upsert/delete, catalog import/export, review decisions, status changes, audit listing, and audit chain verification.
  - Extended external import checks and real E2E to exercise the admin SDK against a live registry.

### Design Intent

The project is positioned as a self-hosted registry and management toolkit, so B2B
operators should be able to embed management workflows directly instead of shelling out
to `ardctl admin`. The SDK remains a thin HTTP wrapper over existing admin JSON
contracts and keeps internal store types out of public packages.

### Files Modified

- `pkg/client/admin.go`
- `pkg/client/client.go`
- `pkg/client/client_test.go`
- `scripts/test-public-go-client.sh`
- `scripts/test-e2e-artifacts.sh`
- `README.md`
- `docs/ARCHITECTURE.md`
- `docs/QUALITY_SCORE.md`
- `docs/releases/feature-release-notes.md`
