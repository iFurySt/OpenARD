# Search

`ard` search is intentionally simple for the first public release. The implementation
prioritizes deterministic, inspectable behavior over opaque ranking.

## Recall

- `POST /search`, `ard search`, and `pkg/client.Search` use the same registry search
  path.
- Query text is split on whitespace and matched case-insensitively against normalized
  catalog entry text.
- The indexed text includes identifier, publisher, display name, media type,
  description, tags, capabilities, and representative queries.
- Query terms are recall-oriented: an entry can match when any term is present.
- Request filters still apply after recall, and inactive lifecycle states are excluded
  from public search.

## Score

`score` is semantic relevance only. It is not a trust, safety, governance, policy, or
compliance signal.

The first-release score is a deterministic approximation:

- Empty query text scores matching entries as `50`.
- A query with no matched terms scores `0`.
- Any matched term starts at `50`.
- Additional matched query terms increase the score linearly up to `100`.

## Ordering

Local search results are ordered by:

1. `score` descending.
2. `identifier` ascending.
3. `displayName` ascending.
4. `source` ascending.

This keeps page-token behavior stable and makes local search compatible with
score-ranked auto federation. Future ranking improvements should preserve the contract
that `score` remains relevance-only and should keep trust and policy signals separate.
