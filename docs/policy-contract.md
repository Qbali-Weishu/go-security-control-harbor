# Compatibility Policy Contract

`as_of = 2026-07-01T00:00:00Z`

This service is used by a defensive engineering team that pre-screens multi-control deployment bundles before they enter regulated or internal network segments. The service must produce a deterministic answer from the bundled policy files only.

## Mandatory evaluation axes

### 1. Platform support

Every selected component must support:

- the profile OS
- the profile runtime
- the profile kernel version

Kernel ranges follow this rule exactly:

- `min` is inclusive
- `max_exclusive` is exclusive

A component must be rejected if the profile kernel is equal to `max_exclusive`.

### 2. Conditional prerequisites

A component can declare:

- `requires`: every matching entry is mandatory when its condition matches the profile
- `requires_any`: at least one candidate must be present when its condition matches the profile

The evaluator must check all matching conditions, not only unconditional entries.

### 3. Directional data-path protection

If the selected bundle contains any `raw_payload_emitter` and the profile zone is `regulated`, the request must include an accepted sanitizer before the first collector in `data_path`.

Accepted sanitizers depend on incident state:

- `steady`: `content-sanitizer` or `payload-redactor`
- `containment`: `content-sanitizer` or `payload-redactor`
- `eradication`: `content-sanitizer` only

Presence alone is insufficient. Order matters:

- the accepted sanitizer must appear before `central-collector`
- if `egress_mode = restricted`, `telemetry-relay` must appear after `central-collector`
- if `egress_mode = restricted` and `fips_mode` is listed in `flow.auditor_required_fips_modes`, `egress-auditor` must appear after `telemetry-relay`

### 4. Conflicts and waivers

Conflicts are bidirectional. If component `A` declares a conflict against `B`, selecting `B` with `A` is still invalid even when `B` does not repeat the declaration.

A blocker can be waived only when all of these are true:

- the rule is marked `waivable`
- the request includes the approval ticket
- the approval `rule_code` matches
- the approval component pair matches, regardless of order
- the profile ID, zone, and incident state are within the approval scope
- if the approval declares `fips_modes`, the profile FIPS mode is in scope
- if the approval declares `egress_modes`, the profile egress mode is in scope
- `profile.as_of < approval.expires_at`

If any one of those checks fails, the waiver is invalid and the blocker remains active.

### 5. Budget aggregation, score, and warning advisories

Total component overhead is the sum across all selected components for every budget dimension:

- `cpu_milli`
- `memory_mb`
- `hook_units`

A profile is incompatible when any summed dimension exceeds its budget.

Role capacity is a separate hard limit. `flow.role_capacity` maps a component role to the
maximum number of selected components that may declare that role. Count every selected
component whose `roles` list contains the limited role; if the count exceeds the configured
limit, the bundle is incompatible. This is independent of the numeric budgets above and applies
regardless of zone or egress mode. Roles without an entry in `flow.role_capacity` are unlimited.

If the bundle remains compatible but any utilization is greater than or equal to `flow.warning_utilization`, the response must still include one advisory `required_action` per over-threshold dimension:

- `review cpu headroom before rollout`
- `review memory headroom before rollout`
- `review hook headroom before rollout`

Score semantics are fixed:

- blocked bundles return `score = 0`
- compatible bundles return `score = 1 - average(cpu_ratio, memory_ratio, hook_ratio)`
- the score is rounded to two decimals

### 6. Deterministic output

The response must be deterministic:

- blockers sorted by `code`, then by component list
- `required_actions` deduplicated and sorted, including warning advisories
- trace totals must reflect the same budgets used for pass/fail and score math, even when multiple axes fail at once
