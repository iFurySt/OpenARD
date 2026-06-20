## [2026-06-20 23:40] | Task: Role-Scoped Admin Tokens

### User Request

> Continue building the Go/Cobra/GORM/Gin/Postgres ARD registry toward a B2B neutral
> registry and management toolkit with real verification and milestone commits.

### Changes

- Added optional role-scoped admin token files through `--admin-tokens-file` and
  `ARD_ADMIN_TOKENS_FILE`.
- Kept `--admin-token` and `ARD_ADMIN_TOKEN` as the compatible full-access admin path.
- Added `reader`, `publisher`, `reviewer`, `operator`, and `admin` roles.
- Bound admin routes to read, publish, review, and operate permissions.
- Added token file parsing, role normalization, and permission tests.
- Added Postgres-backed admin RBAC integration coverage.
- Extended the real artifact E2E flow to run with a token file and verify role behavior
  through `ardctl admin`.
- Added `docs/ADMIN_AUTH.md` and updated README, architecture, security, and quality docs.

### Design Notes

This is intentionally not a full identity provider integration. It gives self-hosted
teams a practical least-privilege step while keeping the deployment model simple and
offline-friendly. Existing single-token deployments keep working unchanged.

### Verification

- Passed: `make fmt-check`
- Passed: `make test`
- Passed: `make build`
- Passed: `make test-integration`
- Passed: `make test-e2e`
- `make test-e2e` verified role-scoped tokens with real MCP, Skill, OpenAPI, and A2A
  artifact onboarding.
- Upstream `ard-spec` manifest and registry conformance passed. The manifest check still
  reports the expected OpenAPI extension media type warning.
