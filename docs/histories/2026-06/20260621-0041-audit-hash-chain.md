## [2026-06-21 00:41] | Task: Audit Hash Chain

### User Request

> Continue the Go/Cobra/GORM/Gin/Postgres ARD implementation with real verification and
> milestone commits.

### Changes

- Added `previousHash` and `hash` fields to persisted admin audit events.
- Hash-chained audit event writes in the Postgres store.
- Added store-level audit chain verification.
- Added migration-time backfill for missing audit hashes without repairing non-empty
  tampered hashes.
- Added `GET /admin/audit/verify` behind reader admin permission.
- Added `ardctl admin audit --verify-chain`.
- Expanded Postgres integration coverage for valid chains and direct DB tampering.
- Expanded admin API integration and real artifact E2E coverage for audit hashes and
  chain verification.
- Updated README, architecture, security, trust, quality, and release notes.

### Design Intent

The audit trail now detects in-database event tampering without introducing an external
signing or immutable-log dependency. This is still an MVP integrity layer: the hash chain
is useful for operator checks, but it does not replace external anchoring, signatures, or
database access control.

### Files Touched

- `internal/store/postgres.go`
- `internal/store/postgres_integration_test.go`
- `internal/httpapi/router.go`
- `internal/httpapi/router_integration_test.go`
- `internal/cli/admin.go`
- `internal/cli/admin_test.go`
- `scripts/test-e2e-artifacts.sh`
- `README.md`
- `docs/ARCHITECTURE.md`
- `docs/QUALITY_SCORE.md`
- `docs/SECURITY.md`
- `docs/TRUST.md`
- `docs/releases/feature-release-notes.md`
