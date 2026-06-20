# Build Version Metadata

## Request

Continue hardening the Go/Cobra/GORM/Gin/Postgres ARD registry and toolkit with real
verification, milestone commits, and release-ready operational surfaces before the first
public tag.

## Changes

- Added shared build metadata for version, commit, and build date.
- Added `version` commands for `ard`, `ardctl`, and `ard-server`, plus Cobra
  `--version` output for binary inventory.
- Embedded build metadata through `make build`, `make package`, Docker builds, and the
  Docker Compose build path.
- Exposed build metadata in server startup logs and `GET /health`.
- Extended the public Go client health response with optional build metadata fields.
- Updated README, architecture, deployment, reliability, SDK compatibility, release
  notes, and quality docs for the new operational surface.

## Design Intent

Release consumers and self-hosted operators need to identify the exact binary serving a
registry without inspecting a package externally. Embedding metadata at build time keeps
the runtime behavior simple while making CLI, Docker, release archive, and health-check
surfaces agree.

## Main Files

- `internal/buildinfo/`
- `internal/cli/version.go`
- `internal/httpapi/router.go`
- `pkg/client/client.go`
- `Makefile`
- `Dockerfile`
- `scripts/package-release.sh`
- `docs/DEPLOYMENT.md`
