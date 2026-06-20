# Public SDK Surface Check

## Request

Continue preparing the Go/Cobra/GORM/Gin/Postgres ARD registry and toolkit for a
verified public release with real checks instead of smoke tests.

## Changes

- Expanded `make test-public-go-client` so it creates an external Go module and uses all
  public `pkg/client` methods.
- Covered public discovery, browse, explore, catalog, health, admin list, admin reviews,
  admin catalog export, upsert, status updates, review decisions, audit listing,
  audit verification, deletion, validation helpers, publisher helpers, and `HTTPError`
  handling from outside the repository.
- Updated SDK compatibility, quality, and release-note docs to describe the stronger
  external import gate.

## Design Intent

The first public tag should not rely only on in-repository tests for the SDK. An external
module catches accidental reliance on private package boundaries and proves that the
documented public import paths can compile and exercise the intended client surface.

## Main Files

- `scripts/test-public-go-client.sh`
- `docs/SDK_COMPATIBILITY.md`
- `docs/QUALITY_SCORE.md`
- `docs/releases/feature-release-notes.md`
