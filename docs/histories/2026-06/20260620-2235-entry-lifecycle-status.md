## [2026-06-20 22:35] | Task: Entry Lifecycle Status

### User Request

> Continue building the Go/Cobra/GORM/Gin/Postgres ARD registry and management toolkit
> toward a B2B neutral registry, with real verification and milestone commits.

### Changes

- Added persisted entry lifecycle status with `active`, `pending`, and `disabled` states.
- Kept imported entries `active` by default.
- Prevented catalog re-imports from overwriting an existing entry lifecycle status.
- Filtered public search, browse, explore, health count, and catalog export to `active`
  entries only.
- Added remote admin status management through `PATCH /admin/entries/:identifier/status`.
- Added `ardctl admin status IDENTIFIER STATUS`.
- Extended admin list filtering with `--status` and lifecycle metadata in admin list
  responses.
- Updated README, architecture, security, and quality docs.
- Extended the real artifact E2E script to disable and reactivate the real Open Browser
  Use Skill through `ardctl admin status`.

### Design Notes

Lifecycle state is an implementation-owned governance control, not an upstream ARD
catalog schema field. Public ARD responses stay spec-shaped and only expose active
entries. Admin list responses include `metadata["ard.status"]` so operators can inspect
state without changing the persisted artifact metadata.

### Verification

- Passed: `make fmt-check`
- Passed: `make test`
- Passed: `make build`
- Passed: `make test-integration`
- Passed: `make test-e2e`
- `make test-e2e` disabled and reactivated the real Open Browser Use Skill through
  `ardctl admin status`, confirmed disabled entries are hidden from public search, and
  confirmed upstream `ard-spec` manifest and registry conformance still pass.
