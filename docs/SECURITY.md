# Security

Use this document to make secure defaults explicit and legible to agents.

## Admin API

- Public ARD discovery routes do not require authentication in the local registry.
- Implementation-specific `/admin/*` routes are disabled by default.
- Set `ARD_ADMIN_TOKEN` or pass `--admin-token` to `ard serve` / `ard-server` to enable
  admin routes.
- Admin requests must send `Authorization: Bearer <token>`.
- Do not log, commit, export, or echo admin tokens.
- Run admin routes behind TLS and a trusted ingress in shared environments. The built-in
  bearer token is an MVP management guard, not a full enterprise identity layer.

## Current Gaps

- No role-based authorization yet.
- No token rotation workflow yet.
- No request audit log yet.
- No signature or trust manifest verification beyond schema-level validation yet.

## Scope

Dependency, SBOM, and provenance integration guidance lives in `docs/SUPPLY_CHAIN_SECURITY.md`.
