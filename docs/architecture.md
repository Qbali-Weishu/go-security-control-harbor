# Architecture

## Packages

- `internal/policy`: loads JSON policy catalogs and approval records.
- `internal/version`: semver-style comparison used for kernel range checks.
- `internal/compat`: compatibility engine, budget evaluation, flow checks, and waiver logic.
- `internal/api`: HTTP handlers and request/response mapping.
- `internal/domain`: shared domain types.

## Execution model

1. Load component, profile, global rule, and approval catalogs at startup.
2. Resolve the request profile and component set.
3. Evaluate platform support, conditional requirements, conflicts, directional flow, and budgets.
4. Emit a stable JSON response with blockers sorted by code, then by component list.
