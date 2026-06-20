## [2026-06-20 21:59] | Task: Add Token-Protected Admin API

### User Request

> Continue building the Go/Cobra/GORM/Gin/Postgres ARD registry and management toolkit
> with real verification, milestone commits, and enterprise-oriented registry behavior.

### Changes

- Added optional Gin `/admin/*` management routes:
  - `GET /admin/entries`
  - `POST /admin/entries`
  - `POST /admin/catalogs`
  - `GET /admin/catalog`
  - `DELETE /admin/entries/:identifier`
- Added `--admin-token` to `ard serve` and `ard-server`.
- Added `ARD_ADMIN_TOKEN` fallback.
- Admin routes are not registered unless a token is configured.
- Admin routes require `Authorization: Bearer <token>` when enabled.
- Updated README, architecture notes, security notes, quality score, and integration
  tests.

### Design Notes

The ARD public discovery API stays unauthenticated for local and federated read flows.
Write/delete/export management over HTTP is implementation-specific and token-protected.
This is intentionally an MVP control plane guard, not a replacement for TLS, ingress
policy, RBAC, audit logs, or enterprise identity.

### Verification

- Passed: `make fmt`
- Passed: `go test ./...`
- Passed: `make build`
- Passed: `make test-integration`
- Passed: HTTP admin API E2E with Postgres:
  - `/admin/entries` returned `404` when no admin token was configured.
  - `/admin/entries` returned `401` when a token was configured but no bearer was sent.
  - `POST /admin/entries` with bearer token created a valid MCP catalog entry.
  - `GET /admin/entries?kind=mcp` listed the created entry.
  - `GET /admin/catalog` exported the created entry.
  - Public `/search` found the created entry.
  - `DELETE /admin/entries/:identifier` removed the entry.
  - Follow-up admin list confirmed the entry was gone.
  - Upstream `ard-spec` registry conformance passed against the token-enabled server.
