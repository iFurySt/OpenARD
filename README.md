# ard

Neutral, self-hosted Agentic Resource Discovery registry and toolkit.

## Intro

`ard` is an independent open-source implementation of the Agentic Resource Discovery
ecosystem. It is intended to become the registry distribution that enterprises can fork,
deploy, and run inside their own environments to discover, verify, and govern agentic
resources.

It targets the discovery layer that sits in front of MCP servers, Skills, A2A agents,
OpenAPI tools, and future agentic artifacts. The goal is not to create another public
marketplace. The goal is to provide a neutral implementation of ARD catalogs,
registries, clients, verification, and publishing workflows that any organization can
self-host.

## Positioning

MCP, Skills, A2A, and OpenAPI define how capabilities are used. ARD defines how they are
found.

`ard` aims to provide:

- A self-hosted ARD registry server.
- A standards-aligned ARD client.
- A CLI and publishing kit for adding catalogs and artifacts quickly.
- A crawler and indexer for `/.well-known/ai-catalog.json` catalogs.
- Verification and policy primitives for trust-aware enterprise adoption.

## Target Developer Experience

The intended first-run flow should feel like this:

```sh
ard serve
ard add catalog https://example.com/.well-known/ai-catalog.json
ard add mcp https://example.com/mcp/server.json
ard add openapi ./openapi.yaml
ard crawl
ard search "query observability logs" --kind mcp
ard verify https://example.com/.well-known/ai-catalog.json
```

## Early Scope

The initial project should focus on the smallest useful enterprise distribution:

- `POST /search` ARD registry endpoint.
- `GET /.well-known/ai-catalog.json` discovery document.
- Catalog ingestion for remote ARD catalogs and local artifacts.
- Exact schema, media type, and domain-anchored `urn:air:` validation.
- Lightweight verification for publisher domain, trust metadata, and pinned artifacts.
- Embedded storage and search for easy local and internal deployment.

## Non-Goals

- A hosted SaaS marketplace.
- A replacement for MCP, Skills, A2A, or OpenAPI.
- An agent execution runtime.
- A heavy enterprise governance suite in the first release.
- A platform-specific registry tied to one model provider or agent framework.

## Project Notes

This repository is currently in product and architecture bootstrap. The detailed working
context is captured in:

- [Product Sense](docs/PRODUCT_SENSE.md)
- [Architecture](docs/ARCHITECTURE.md)
- [ARD Market Context](docs/references/ard-market-context.md)
- [ARD Spec Working Notes](docs/references/ard-spec-working-notes.md)

## License

[Apache-2.0](LICENSE)
