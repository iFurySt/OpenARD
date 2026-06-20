# Local Search Ranking

## Request

Continue hardening the Go/Cobra/GORM/Gin/Postgres ARD registry and toolkit toward a
verified first public release.

## Changes

- Ordered local search results by descending relevance score before pagination.
- Added stable tie-breakers by identifier, display name, and source.
- Added unit coverage for the ranking helper and relevance score calculation.
- Added Postgres integration coverage proving higher-score local results are returned
  first.
- Added `docs/SEARCH.md` to document the first-release search contract.
- Updated architecture, quality, and release-note docs.

## Design Intent

The ARD spec treats `score` as semantic relevance only. The first release should make
that behavior deterministic and auditable while avoiding hidden trust or governance
signals in ranking. Stable ordering also keeps page-token behavior predictable.

## Main Files

- `internal/store/postgres.go`
- `internal/store/search_ranking_test.go`
- `internal/store/postgres_integration_test.go`
- `docs/SEARCH.md`
- `docs/ARCHITECTURE.md`
