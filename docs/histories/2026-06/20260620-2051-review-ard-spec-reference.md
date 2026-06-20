## [2026-06-20 20:51] | Task: Review ARD Spec Reference

### Execution Context

- Agent ID: `codex`
- Base Model: `GPT-5`
- Runtime: `local workspace`

### User Query

> Review the local `ard-spec` checkout and upstream repository, verify whether previous
> project assumptions still hold, and capture the reference strategy in this repo.

### Changes Overview

- Area: product, architecture, and reference documentation.
- Key actions:
  - Added upstream ARD spec working notes.
  - Corrected project docs to prefer `urn:air:` over older `urn:ai:` references.
  - Recorded `application/mcp-server-card+json`, conformance, and optional endpoint
    implications.
  - Documented why a git submodule is not recommended yet.

### Design Intent

This change makes `ards-project/ard-spec` the explicit implementation reference while
avoiding premature submodule overhead. The project should track upstream spec behavior
strictly, but only vendor or submodule artifacts when implementation and CI actually need
them.

### Files Modified

- `README.md`
- `docs/PRODUCT_SENSE.md`
- `docs/ARCHITECTURE.md`
- `docs/references/README.md`
- `docs/references/ard-market-context.md`
- `docs/references/ard-spec-working-notes.md`
- `docs/histories/2026-06/20260620-2051-review-ard-spec-reference.md`
