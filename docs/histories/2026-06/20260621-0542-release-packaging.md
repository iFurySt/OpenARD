## [2026-06-21 05:42] | Task: Release Packaging

## Request

> Continue building the Go/Cobra/GORM/Gin/Postgres ARD registry and toolkit with real
> verification, milestone commits, and strict alignment with the ARD self-hosted
> enterprise distribution goal.

## Changes

- Added `scripts/package-release.sh` to build versioned release archives for `ard`,
  `ardctl`, and `ard-server`.
- Added `make package`, producing Linux and macOS amd64/arm64 archives by default.
- Generated `dist/checksums.txt` with SHA-256 hashes for every archive.
- Added release packaging to GitHub Actions CI.
- Documented binary archives, checksum verification posture, and remaining SBOM/provenance
  gaps.

## Design Notes

Release packaging is intentionally boring and repository-native: it uses the existing Go
toolchain plus shell utilities instead of adding a release framework before public tags
exist. The archives give operators a concrete distribution path today while leaving
signed checksums, SBOMs, and provenance attestations as the next supply-chain milestone.

## Files

- `.github/workflows/ci.yml`
- `Makefile`
- `scripts/package-release.sh`
- `docs/DEPLOYMENT.md`
- `docs/SUPPLY_CHAIN_SECURITY.md`
- `docs/QUALITY_SCORE.md`
- `docs/releases/feature-release-notes.md`
