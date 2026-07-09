# Security Control Compatibility Gate

This service evaluates whether a proposed security control bundle can be deployed on a managed Linux host profile without violating platform, data-path, waiver, or resource constraints.

## Runtime contract

- `GET /healthz` returns a basic liveness document.
- `POST /v1/compatibility/assess` accepts a bundle request and returns a deterministic compatibility decision.
- Policy files are loaded from `testdata/policies/` by default, or from `COMPAT_POLICY_ROOT` if it is set.

## Local run

```bash
go run ./cmd/compatibilityd
```

Once the server is running, submit any JSON request under `testdata/requests/`
with an HTTP client available in your environment.

## Key engineering constraint

The evaluator must reason across multiple axes at once:

- platform support: OS, runtime, kernel range
- prerequisite closure: direct and conditional requirements
- directional data-flow safety in regulated zones
- conflict handling with ticket-scoped waivers
- summed CPU, memory, and hook budgets
- deterministic trace output for downstream auditors
