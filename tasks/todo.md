# Architecture normalization checklist

Normative decisions: [`docs/plans/architecture-normalization.md`](../docs/plans/architecture-normalization.md). Task acceptance and verification details: [`tasks/plan.md`](plan.md).

## Phase 0 — Plan and decision freeze

- [x] 0.1 Record expensive architecture decisions as ADRs.
- [x] 0.2 Refresh and freeze the dependency disposition matrix.
- [x] 0.3 Obtain explicit implementation approval.
- [x] Checkpoint 0: decisions, ADRs, dependencies, and implementation scope approved.

## Phase 1 — Safety net and baseline

- [x] 1.1 Create one non-mutating zero-warning quality command.
- [x] 1.2 Enable the accepted TypeScript strictness baseline.
- [x] 1.3 Add Node unit and Playwright browser harnesses.
- [x] Checkpoint 1A: quality and test harnesses green.
- [x] 1.4 Characterize runtime base paths and auth isolation.
- [x] 1.5 Characterize critical admin state transitions.
- [x] 1.6 Characterize backend contracts.
- [x] 1.7 Record reproducible performance baselines and numeric guardrails.
- [x] Checkpoint 1B: risky boundaries protected and baseline recorded.

## Phase 2 — Correctness

- [x] 2.1 Implement the runtime base-path integration boundary.
- [x] 2.2 Isolate auth cookies by base path.
- [x] 2.3 Correct site and persistence error semantics.
- [x] 2.4 Correct auth persistence and timestamp behavior.
- [x] 2.5 Correct GeoIP synchronization and refresh failures.
- [x] Checkpoint 2: known defects closed with regression evidence.

## Phase 3 — Dependencies and mechanical normalization

- [x] 3.1 Upgrade frontend compiler and core runtime.
- [x] 3.2 Remove redundant frontend tooling dependencies.
- [x] 3.3 Modernize Go dependencies and executable tools.
- [x] 3.4 Remove unnecessary Go runtime packages.
- [x] 3.5 Remove proven dead code and normalize formatting.
- [x] Checkpoint 3: approved dependency graph and fresh size comparison.

## Phase 4 — Frontend normalization

- [x] 4.1 Establish `app`, `features`, and `shared` boundaries.
- [x] 4.2 Rebuild the canonical current shadcn foundation.
- [x] 4.3 Migrate the auth vertical slice.
- [x] 4.4 Migrate the sites vertical slice.
- [x] 4.5 Migrate analytics and typed URL state.
- [x] 4.6 Migrate site settings.
- [x] 4.7 Migrate event definitions.
- [x] 4.8 Enforce the no-flicker data lifecycle.
- [x] 4.9 Delete legacy frontend structure.
- [x] Checkpoint 4: accepted UX and no frontend legacy layer.

## Phase 5 — Backend normalization

- [x] 5.1 Separate application composition and lifecycle.
- [x] 5.2 Normalize auth feature ownership.
- [x] 5.3 Normalize site feature and selective persistence mapping.
- [x] 5.4 Normalize analytics feature ownership.
- [x] 5.5 Normalize events and GeoIP feature ownership.
- [x] 5.6 Split GraphQL ownership by feature.
- [x] 5.7 Normalize platform and seed ownership.
- [x] 5.8 Delete legacy backend layers.
- [x] Checkpoint 5: explicit dependency direction and no backend legacy layer.

## Phase 6 — API, data, and performance

- [x] 6.1 Normalize GraphQL contracts and error categories.
- [x] 6.2 Reduce measured dashboard query fan-out.
- [x] 6.3 Normalize polling, cache, pagination, and stale data.
- [x] 6.4 Enforce measured performance guardrails.
- [x] Checkpoint 6: correctness, compatibility, and performance evidence complete.

## Phase 7 — Hardening and closure

- [x] 7.1 Run security and privacy hardening.
- [x] 7.2 Run the complete release matrix.
- [x] 7.3 Finalize maintained documentation and examples.
- [x] 7.4 Close every definition-of-done item with evidence.
- [ ] Final checkpoint: maintainer review and no temporary or unintended changes.
