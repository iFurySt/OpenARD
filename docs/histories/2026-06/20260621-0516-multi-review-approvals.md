## [2026-06-21 05:16] | Task: Multi-Reviewer Approval Thresholds

### Execution Context

- Agent ID: `codex`
- Base Model: `GPT-5`
- Runtime: `Codex CLI`

### User Query

> Continue building the Go/Cobra/GORM/Gin/Postgres ARD registry and toolkit with real
> verification, milestone commits, and strict alignment with the ARD self-hosted
> registry direction.

### Changes Overview

- Area: governance, admin API, CLI, SDK, persistence, tests, docs.
- Key actions:
  - Added `requiredApprovals` to ingestion policy, normalized empty and zero values to
    one approval, and rejected negative values.
  - Added persisted review approval records keyed by entry identifier and reviewer
    token name.
  - Changed approval handling so a pending entry remains pending until the distinct
    reviewer threshold is reached; duplicate reviewer approvals return a conflict.
  - Kept rejection immediate: rejecting a pending entry still disables it and records the
    audit reason.
  - Exposed approval counts through the admin API, `pkg/client`, and `ardctl admin
    review approve`.
  - Extended integration and E2E coverage to verify two-reviewer approval flows against
    policy-pending Skill imports and updates.

### Design Intent

Self-hosted enterprise registries need more than a single-person approval switch for
resource publication. The implementation keeps ARD catalog entries spec-shaped by
storing approval state in a registry-owned table, while using existing role token names
as stable local reviewer identifiers. Empty policy configuration stays backward
compatible by requiring one approval.

### Files Modified

- `internal/policy/policy.go`
- `internal/policy/policy_test.go`
- `internal/store/postgres.go`
- `internal/httpapi/router.go`
- `internal/httpapi/router_integration_test.go`
- `internal/cli/admin.go`
- `pkg/client/admin.go`
- `pkg/client/client_test.go`
- `scripts/test-e2e-artifacts.sh`
- `docs/POLICY.md`
- `docs/ADMIN_AUTH.md`
- `docs/ARCHITECTURE.md`
- `docs/QUALITY_SCORE.md`
- `docs/releases/feature-release-notes.md`
