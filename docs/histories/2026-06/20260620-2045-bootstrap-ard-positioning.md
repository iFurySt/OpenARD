## [2026-06-20 20:45] | Task: Bootstrap ARD Positioning

### Execution Context

- Agent ID: `codex`
- Base Model: `GPT-5`
- Runtime: `local workspace`

### User Query

> Capture the ARD opportunity, positioning, and planning conclusions into this repository
> so future development can continue here.

### Changes Overview

- Area: product and architecture documentation.
- Key actions:
  - Reframed the repository from a base template into the `ard` project.
  - Captured the neutral self-hosted registry positioning.
  - Added target product surfaces, MVP scope, non-goals, and initial architecture.
  - Added a durable ARD market-context reference with links to Google and Hugging Face
    launch articles.

### Design Intent

This change preserves the planning context needed for future implementation without
depending on chat history. The chosen direction is a self-hosted, vendor-neutral ARD
registry distribution with client and publishing kit, rather than a public marketplace or
platform-specific registry.

### Files Modified

- `AGENTS.md`
- `README.md`
- `docs/PRODUCT_SENSE.md`
- `docs/ARCHITECTURE.md`
- `docs/references/README.md`
- `docs/references/ard-market-context.md`
- `docs/histories/2026-06/20260620-2045-bootstrap-ard-positioning.md`
