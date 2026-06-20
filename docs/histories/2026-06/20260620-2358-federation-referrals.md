## [2026-06-20 23:58] | Task: Federation Referrals

### User Request

> Continue building the Go/Cobra/GORM/Gin/Postgres ARD registry toward a B2B neutral
> registry and management toolkit with real verification and milestone commits.

### Changes

- Added `SearchResponse.referrals` support for `POST /search` when
  `federation=referrals`.
- Added store support for active registry referral entries using
  `application/ai-registry+json` and `application/ai-registry`.
- Added `ard search --federation none|referrals|auto`.
- Kept CLI default federation mode as `none` to avoid changing current scripted behavior.
- Added Postgres-backed HTTP integration coverage for referrals mode.
- Extended E2E to import a local upstream registry referral and verify it through
  `ardctl search --federation referrals --json`.
- Updated README, architecture, and quality docs.

### Design Notes

This implements the spec's client-followed federation path. `referrals` mode returns
local results plus registry catalog entries the client may query next. It does not yet
implement `federation=auto`, where the registry would query upstream registries and merge
results server-side.

### Verification

- Passed: `make fmt-check`
- Passed: `make test`
- Passed: `make build`
- Passed: `make test-integration`
- Passed: `make test-e2e`
- `make test-e2e` verified referrals mode while also running real MCP, Skill, OpenAPI,
  and A2A artifact onboarding.
- Upstream `ard-spec` manifest and registry conformance passed. The manifest check still
  reports the expected OpenAPI extension media type warning.
