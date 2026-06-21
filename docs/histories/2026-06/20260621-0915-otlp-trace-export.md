## [2026-06-21 09:15] | Task: OTLP Trace Export

## Request

Continue hardening the Go/Cobra/Gin/GORM/Postgres ARD registry and toolkit with real
verification, milestone commits, and strict alignment with the self-hosted enterprise
direction.

## Changes

- Added an optional OTLP/HTTP JSON trace exporter for completed registry server spans.
- Added `ARD_OTLP_TRACES_ENDPOINT` and `--otlp-traces-endpoint` for `ard serve` and
  `ard-server`; trace export remains disabled by default.
- Preserved inbound W3C parent span IDs so exported server spans can point back to the
  caller span.
- Extended the real artifact E2E flow with a local OTLP capture endpoint and verified a
  live registry `/health` request exports a trace span.
- Added `docs/OBSERVABILITY.md` and updated reliability, deployment, security,
  architecture, quality, README, and release-note docs.

## Intent

Enterprises need a way to connect registry activity to existing observability backends
without making tracing a hard runtime dependency. The implementation keeps the existing
lightweight trace propagation model and adds a narrow OTLP/HTTP export path only when an
operator configures an endpoint.

## Files

- `internal/traceexporter/`
- `internal/tracecontext/`
- `internal/httpapi/`
- `internal/cli/`
- `internal/config/`
- `scripts/test-e2e-artifacts.sh`
- `docs/OBSERVABILITY.md`
