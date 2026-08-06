# ADR-0005: Enforce zero-warning quality gates and incremental delivery

## Status

Accepted

## Date

2026-08-06

## Context

The normalization affects routing, UI state, GraphQL, persistence, migrations, deployment paths, and most package boundaries. A big-bang rewrite would make regressions hard to locate and would leave no reliable intermediate state. Existing checks also lint only part of the frontend, omit browser coverage, and do not fully enforce documented formatting and strictness.

The desired outcome is not merely reorganized code. It is a stable personal admin panel, strict end-to-end contracts, current dependencies, preserved data, and a repeatable guarantee that future changes remain safe.

This ADR records decisions D-006, D-007, D-018 through D-024, and D-030.

## Decision

- Implement the update through the approved seven-phase plan in small, reviewable increments.
- Every increment leaves all relevant checks green; temporary migration bridges exist only inside their owning phase and are removed before completion.
- CI blocks on warnings or errors in formatting, import organization, lint, strict application types, tests, builds, and generated-artifact freshness.
- Dashboard application code enables `exactOptionalPropertyTypes`; `skipLibCheck` remains limited to external declaration files.
- Exceptions are narrow, documented, and limited to generated or third-party-derived code.
- Use risk-based testing rather than a global coverage percentage: characterize high-risk behavior first, add a regression test for each fixed defect, use real-browser tests for critical admin/base-path behavior, and exercise both databases where dialect behavior matters.
- Establish reproducible performance, query, bundle, and binary baselines before setting numeric guardrails. UI stability, critical-flow responsiveness, and dashboard-query behavior take priority over size alone.
- Upgrade to latest stable dependencies only. Prereleases are excluded; Rust/SWC-based tooling is a strong tie-breaker, while correctness, maintainability, and architecture fit control.
- New direct dependencies and replacements require an approved phase plan. Standard-library or accepted-stack solutions are preferred.
- Architecture and migration rationale belongs in maintained documentation or ADRs. Code comments explain non-obvious reasons and invariants at their enforcement boundary rather than narrating code.

## Alternatives considered

### Big-bang rewrite followed by final testing

Rejected because failures would compound across layers and rollback or review would be impractical.

### Global coverage-percentage target

Rejected because a percentage can reward low-value tests while missing risky user and data contracts.

### Adopt every latest or Rust-based tool regardless of fit

Rejected because implementation language and version recency do not outweigh correctness or maintenance cost.

### Permanent compatibility layer

Rejected because the final state must contain no legacy architecture. Compatibility is preserved only where explicitly accepted, with migration plans for approved breaks.

## Consequences

- Phase checkpoints and completion evidence are part of delivery, not optional reporting.
- Tooling and browser-test dependencies are justified by the accepted phase plan.
- A new material decision pauses implementation until it is recorded and approved.
- The update is complete only when the full definition of done has linked evidence and no temporary or legacy layer remains.
