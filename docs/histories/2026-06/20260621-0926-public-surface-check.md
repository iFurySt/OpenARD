## [2026-06-21 09:26] | Task: Public Surface Check

## Request

Continue hardening the Go/Cobra/Gin/GORM/Postgres ARD registry and toolkit with real
verification, milestone commits, and strict alignment with the self-hosted enterprise
direction.

## Changes

- Added `internal/tools/publicsurface`, a repository-native compatibility checker for
  exported `pkg/ard` and `pkg/client` symbols.
- Added CLI surface checks for the direct commands and root flags exposed by `ard`,
  `ardctl`, and `ard-server`.
- Added `make check-public-surface`.
- Wired the public surface check into GitHub Actions CI and made the workflow checker
  require that CI step.
- Updated SDK compatibility, collaboration, architecture, quality, README, and release
  notes.

## Intent

The project is nearing a first public tag. The public Go SDK and CLI command shape now
have a mechanical compatibility gate so accidental surface drift is caught before a
release rather than discovered by downstream adopters.

## Files

- `internal/tools/publicsurface/`
- `internal/tools/workflowcheck/`
- `.github/workflows/ci.yml`
- `Makefile`
- `docs/SDK_COMPATIBILITY.md`
- `docs/QUALITY_SCORE.md`
