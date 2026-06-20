# Architecture

This file is the top-level map for `ard`. The project is currently in bootstrap, so the
sections below describe the target architecture and boundaries that future implementation
should follow.

## Product Surfaces

- Registry server: self-hosted ARD registry exposing discovery, search, health, and
  catalog endpoints.
- CLI: operational entry point for serve, add, crawl, verify, search, and export flows.
- Client library: standards-aligned ARD client for registry search and catalog traversal.
- Publisher kit: helpers that convert existing MCP, Skills, A2A, OpenAPI, and URL
  artifacts into ARD catalog entries.
- Verification engine: reusable validation and trust checks shared by server and CLI.

## Intended Repository Shape

- `cmd/ard/`: CLI entry point when implementation begins.
- `apps/registry/`: deployable registry server if the server is separated from the CLI
  binary.
- `packages/ard/`: shared ARD models, schemas, request/response handling, and client code.
- `packages/catalog/`: catalog parsing, generation, crawling, and federation traversal.
- `packages/index/`: indexing and retrieval abstractions.
- `packages/verify/`: schema validation, URN checks, publisher-domain checks, trust
  metadata validation, and artifact pinning.
- `packages/conformance/`: wrappers or integration glue for upstream ARD conformance
  checks, if needed after implementation begins.
- `packages/adapters/`: MCP, Skills, A2A, OpenAPI, and future resource adapters.
- `infra/`: Docker, deployment, and environment definitions.
- `scripts/`: repository automation that agents can run directly.
- `docs/`: repository knowledge base and system of record.

The exact language and package layout can change, but the boundaries above should remain
visible. A Go implementation is a strong default because it supports a single binary,
simple containers, and enterprise-friendly deployment.

## Runtime Topology

The smallest useful deployment is one registry process with embedded persistence:

```text
catalog URLs / local artifacts
        |
        v
crawler + adapter layer
        |
        v
validation + verification
        |
        v
metadata store + search index
        |
        v
ARD /search API + CLI client
```

The first storage target should be embedded and operationally boring, such as SQLite with
FTS. The architecture should leave room for Postgres, OpenSearch, Meilisearch, or vector
backends later, but those should not be required for a basic deployment.

## Core Data Flow

1. A user adds a catalog, registry, or artifact with the CLI or API.
2. The crawler fetches `/.well-known/ai-catalog.json` or a direct artifact URL.
3. The adapter layer normalizes supported artifacts into ARD catalog entries.
4. The verification layer validates schema, media type, `url`/`data` exclusivity,
   domain-anchored `urn:air:` identifiers, publisher domains, and trust metadata.
5. The index layer stores normalized entries and searchable fields.
6. `POST /search` accepts an ARD `SearchRequest` and returns a ranked `SearchResponse`.
7. Clients fetch the selected artifact and execute it through its native protocol.

## Boundary Rules

- Keep ARD models and protocol handling independent from transport, storage, and CLI
  code.
- Keep verification pure and reusable; both server-side ingestion and CLI validation
  should call the same logic.
- Keep adapters narrow. MCP, Skills, A2A, and OpenAPI adapters should translate metadata,
  not execute tools.
- Search and ranking should consume normalized catalog entries, not protocol-specific
  objects.
- Federation traversal should be bounded by depth, registry count, response size, and
  timeout controls.
- Secrets and tokens may be used during request scope only; they must not be stored or
  emitted in plain text.
- Specification behavior should be derived from `ards-project/ard-spec`, especially
  `spec/ard.md`, `spec/schemas/`, ADRs, and `conformance/`.

## Initial API Targets

- `GET /.well-known/ai-catalog.json`: advertise this registry and any configured catalog
  entries.
- `POST /search`: ARD search endpoint.
- `POST /explore`: optional; may initially return `501 Not Implemented`.
- `GET /agents`: optional deterministic browse endpoint; useful for B2B portals but not
  required for the first conformance pass.
- `GET /health`: deployment health.
- CLI equivalents: `serve`, `add`, `crawl`, `verify`, `search`, and `export`.

## Specification Alignment

The upstream specification source is:

- Repository: `https://github.com/ards-project/ard-spec`
- Rendered spec: `https://agenticresourcediscovery.org/spec/`
- Current observed draft: v0.9
- Current observed local checkout commit during planning: `a78be70`

Implementation decisions should prefer the upstream main spec, schemas, ADRs, and
conformance tool over older reference implementations. In particular:

- Use `urn:air:` identifiers, not the older `urn:ai:` form.
- Treat `application/mcp-server-card+json` as the MCP discovery media type.
- Keep `score` strictly as semantic relevance, not a trust or safety signal.
- Support web ingestion of `ai-catalog.json` catalogs as a required registry capability.
- Keep `/explore` local-only and optional; if unsupported, return `501`.
- Keep federation controlled by root-level `SearchRequest.federation`.

Do not vendor or fork the upstream spec content casually. If the implementation needs
schemas or conformance tools in-repo, add a pinned, documented copy under a clearly named
third-party or generated directory and record the source commit.

## Open Decisions

- Final implementation language and framework.
- Embedded storage schema.
- Ranking strategy for the first release.
- Trust manifest verification depth for MVP.
- Whether the server and CLI ship as one binary or separate packages.
- Whether to vendor selected upstream spec artifacts, use a git submodule, or fetch pinned
  artifacts during development.

When these decisions are made, update this file in the same task as the code.
