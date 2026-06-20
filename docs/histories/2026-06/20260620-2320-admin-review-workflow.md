## [2026-06-20 23:20] | Task: Admin Review Workflow

### User Request

> Continue building the Go/Cobra/GORM/Gin/Postgres ARD registry toward a B2B neutral
> registry and management toolkit, with real verification and milestone commits.

### Changes

- Added `GET /admin/reviews` to list pending entries.
- Added `POST /admin/reviews/:identifier/approve` to activate a pending entry.
- Added `POST /admin/reviews/:identifier/reject` to disable a pending entry.
- Added `ardctl admin review list`.
- Added `ardctl admin review approve IDENTIFIER`.
- Added `ardctl admin review reject IDENTIFIER`.
- Added dedicated review audit actions for approve and reject decisions.
- Extended Postgres-backed HTTP integration tests for review list, approve, and reject.
- Extended the real artifact E2E flow to approve and reject a policy-pending real Skill.
- Added E2E remote artifact prefetch retries so live MCP, Skill, and OpenAPI checks are
  less sensitive to transient EOFs from upstream hosts.
- Updated README, policy, architecture, and quality docs.

### Design Notes

Review is intentionally a thin workflow over lifecycle status. Policy can create pending
entries, public discovery hides them, and admins can explicitly approve or reject them
without remembering the lower-level lifecycle status command. This keeps the B2B review
surface legible while preserving the same storage and public filtering behavior.
Review decisions only accept entries that are still pending; the lower-level status
command remains the escape hatch for general lifecycle changes.

### Verification

- Passed: `make fmt-check`
- Passed: `make test`
- Passed: `make build`
- Passed: `make test-integration`
- Passed: `make test-e2e`
- `make test-e2e` verified `ardctl admin review list`, `approve`, and `reject` against
  a policy-pending real Open Browser Use Skill.
- Upstream `ard-spec` manifest and registry conformance passed. The manifest check still
  reports the expected OpenAPI extension media type warning.
