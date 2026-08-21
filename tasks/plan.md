# Implementation plan: architecture normalization

## Authority and scope

[`docs/plans/architecture-normalization.md`](../docs/plans/architecture-normalization.md) is the normative source of truth. This file only decomposes its accepted roadmap into executable work. A conflict is resolved in favor of the master document; a new architectural choice must be recorded and approved there before implementation continues.

Implementation was approved on 2026-08-06 for the accepted plan and dependency set. No task below authorizes a tracker breaking change, destructive data migration, or unplanned direct dependency.

## Execution rules

- Complete tasks in dependency order and keep the repository green after every task.
- Limit one implementation batch to at most five handwritten files. Repeat a batch task when a feature needs more files.
- Generated files do not count toward the batch limit, but must only change through their generator.
- Add characterization coverage before changing risky behavior and a regression test with every defect fix.
- Run project commands inside the existing Dev Container.
- At each checkpoint, record results in the master document. Continue under the approved plan unless evidence introduces a new decision or expands compatibility scope.

## Dependency graph

```text
accepted decisions and ADRs
  -> reproducible checks and characterization tests
    -> correctness and runtime-base-path fixes
      -> dependency/toolchain normalization
        -> vertical frontend feature migrations
          -> vertical backend feature migrations
            -> API/data/query simplification
              -> final hardening and release evidence
```

## Phase 0 — Plan and decision freeze

### Task 0.1 — Record expensive architecture decisions as ADRs

**Acceptance criteria:** ADRs cover frontend ownership, Go modular-monolith boundaries/composition, GraphQL and selective persistence mapping, and quality/delivery policy; each links back to decision IDs in the master document.

**Verification:** Documentation links resolve; no proposal is presented as accepted without a recorded decision.

**Dependencies:** None.
**Likely files:** `docs/decisions/*`, master plan.
**Scope:** Medium, split by ADR.

### Task 0.2 — Freeze dependency disposition and official sources

**Acceptance criteria:** Every current and planned direct dependency has keep/add/remove/replace status, current official migration evidence, and owning phase; registry versions are refreshed immediately before implementation.

**Verification:** `go list -u -m all`, npm registry audit, import/reference search, official-source links.

**Dependencies:** Task 0.1.
**Likely files:** master plan, `tasks/todo.md`.
**Scope:** Small.

### Task 0.3 — Obtain implementation approval

**Acceptance criteria:** The maintainer explicitly approves this phase plan and dependency set. Any requested correction is reflected in the master plan before code changes begin.

**Verification:** Approval is recorded in the interview history and document status.

**Dependencies:** Tasks 0.1–0.2.
**Likely files:** master plan.
**Scope:** Small.

### Checkpoint 0

- All accepted expensive decisions have ADRs.
- No blocking clarification remains.
- Implementation and dependency changes have explicit approval.

## Phase 1 — Enforceable safety net and baseline

### Task 1.1 — Create one non-mutating quality command

**Acceptance criteria:** One root command runs full Biome formatting/import/lint checks, strict dashboard/tracker TypeScript, Go format/lint/tests, generation freshness, and builds with zero warnings.

**Verification:** Run the new root check command twice; the second run produces no worktree changes.

**Dependencies:** Task 0.3.
**Likely files:** root/dashboard/server Taskfiles, package scripts, CI workflow.
**Scope:** Medium.

### Task 1.2 — Enable the accepted TypeScript strictness baseline

**Acceptance criteria:** Dashboard enables `exactOptionalPropertyTypes`; browser and tooling configs have explicit environments; `skipLibCheck` remains limited to external declarations; application errors are fixed without unsafe assertions.

**Verification:** Dashboard and tracker typechecks pass.

**Dependencies:** Task 1.1.
**Likely files:** dashboard and tracker TypeScript configs plus affected focused files.
**Scope:** Repeatable medium batches.

### Task 1.3 — Add frontend unit and browser harnesses

**Acceptance criteria:** Bun runs pure unit tests; Playwright runs the real dashboard against the application; test configuration works locally and in CI without a second JavaScript test runner.

**Verification:** One unit smoke test and one browser smoke test pass in the Dev Container and CI configuration validates.

**Dependencies:** Task 1.1.
**Likely files:** package/Taskfile, Playwright config, test helpers, workflow.
**Scope:** Medium.

### Checkpoint 1A

- Zero-warning checks pass.
- Unit and real-browser harnesses run reproducibly.
- No generated artifact is stale.

### Task 1.4 — Characterize runtime base paths and auth isolation

**Acceptance criteria:** One dashboard build is exercised at `/`, `/lovely-eye`, and `/tools/lovely-eye`; nested refresh, assets, GraphQL, login/logout, and same-origin multi-instance cookie isolation have failing-before/fixed-after coverage.

**Verification:** Base-path browser matrix and focused Go HTTP tests pass.

**Dependencies:** Task 1.3.
**Likely files:** browser specs, server dashboard/auth tests, fixtures.
**Scope:** Repeatable medium batches.

### Task 1.5 — Characterize critical admin state transitions

**Acceptance criteria:** Tests cover auth resolution, first-load skeletons, retained background data, mutation feedback, layout stability, site switching, multi-domain editing, URL state, back/forward, and refresh.

**Verification:** Focused Playwright scenarios pass without arbitrary sleeps.

**Dependencies:** Task 1.3.
**Likely files:** browser specs and helpers.
**Scope:** Repeatable medium batches.

### Task 1.6 — Characterize backend contracts

**Acceptance criteria:** Focused tests cover GraphQL error categories, server construction, configuration failures, dashboard fallback, site/repository errors, and current aggregate multi-domain collection.

**Verification:** Focused Go packages and existing e2e suite pass.

**Dependencies:** Task 1.1.
**Likely files:** Go test files by owning package.
**Scope:** Repeatable medium batches.

### Task 1.7 — Record reproducible performance baselines

**Acceptance criteria:** Fresh compressed initial/route bundle sizes, Go binary sizes, critical browser timings, dashboard query count/fan-out, and representative database timings are recorded with commands and environment.

**Verification:** A second run produces comparable measurements; numeric guardrails are added to the master plan.

**Dependencies:** Tasks 1.3–1.6.
**Likely files:** measurement scripts/config and master plan.
**Scope:** Medium.

### Checkpoint 1B

- Characterization coverage protects every high-risk boundary.
- Baseline and numeric guardrails are recorded.
- Full repository checks are green.

## Phase 2 — Correctness and deployment invariants

### Task 2.1 — Implement the runtime base-path integration boundary

**Acceptance criteria:** Vite emits relative assets; Go applies normalized `BASE_PATH` once; runtime config and TanStack `basepath` agree; SPA fallback serves nested routes; application routes remain unprefixed.

**Verification:** ADR-0001 matrix passes with one build.

**Dependencies:** Task 1.4.
**Likely files:** Vite config, runtime config, router, Go dashboard/server boundary, tests.
**Scope:** Medium batches.

### Task 2.2 — Isolate auth cookies by base path

**Acceptance criteria:** Cookie names and `Path` derive deterministically from normalized `BASE_PATH`; set/read/refresh/clear use one policy; security attributes remain intact; same-origin instances do not interfere.

**Verification:** Focused auth tests and multi-instance browser scenario pass.

**Dependencies:** Tasks 1.4 and 2.1.
**Likely files:** auth cookie policy, composition config, tests.
**Scope:** Medium.

### Task 2.3 — Correct site and persistence error semantics

**Acceptance criteria:** Domain uniqueness checks never discard repository errors; not-found is distinct from infrastructure failure; create/update/delete return stable feature errors; multi-domain aggregation remains unchanged.

**Verification:** Table-driven service/repository tests pass on SQLite and PostgreSQL where dialect behavior matters.

**Dependencies:** Task 1.6.
**Likely files:** site service/repository and tests.
**Scope:** Medium batches.

### Task 2.4 — Correct auth persistence and timestamp behavior

**Acceptance criteria:** Created users return persisted timestamps; cookie handling is outside the authentication application service; transport maps stable errors without leaking persistence types.

**Verification:** Auth unit, GraphQL, and e2e tests pass.

**Dependencies:** Tasks 1.6 and 2.2.
**Likely files:** auth feature, transport adapter, tests.
**Scope:** Medium batches.

### Task 2.5 — Correct GeoIP synchronization and refresh failures

**Acceptance criteria:** All synchronization/download/refresh failures are preserved or explicitly classified as best-effort; UI receives actionable state; no error is silently discarded.

**Verification:** Focused GeoIP/service/GraphQL tests pass.

**Dependencies:** Task 1.6.
**Likely files:** GeoIP service/adapter/resolver and tests.
**Scope:** Medium batches.

### Checkpoint 2

- Every known correctness defect has a regression test.
- Base-path, cookie, tracker, and multi-domain invariants pass.
- Full repository checks are green.

## Phase 3 — Dependency and mechanical normalization

### Task 3.1 — Upgrade frontend compiler and core runtime

**Acceptance criteria:** React, TypeScript 7, Vite, SWC plugin, Biome, GraphQL client/codegen, and TanStack packages are latest stable; official migration changes are applied; no deprecated config remains.

**Verification:** Generation, full checks, production build, unit tests, and browser smoke pass.

**Dependencies:** Checkpoint 2.
**Likely files:** package/lock/config files plus focused compatibility fixes.
**Scope:** Split by ecosystem into medium batches.

### Task 3.2 — Remove redundant frontend tooling dependencies

**Acceptance criteria:** Babel React plugin, PostCSS path, direct router generator, and unused direct `react-is` are absent; the active Vite/SWC/Tailwind path builds identically.

**Verification:** Reference search, frozen install, production build, and bundle comparison pass.

**Dependencies:** Task 3.1.
**Likely files:** package/lock, configs, docs.
**Scope:** Medium.

### Task 3.3 — Modernize Go dependencies and executable tools

**Acceptance criteria:** Direct modules are latest stable; `tools.go` is replaced with Go `tool` directives; generation uses pinned tools; `go mod tidy`, verify, and vulnerability scan are clean.

**Verification:** Generate twice, `go mod verify`, `govulncheck ./...`, lint, tests, and builds pass.

**Dependencies:** Checkpoint 2.
**Likely files:** `go.mod`, `go.sum`, Taskfile, generation directive, tools file removal.
**Scope:** Medium.

### Task 3.4 — Remove unnecessary Go runtime packages

**Acceptance criteria:** Migration command uses explicit standard-library dispatch; `urfave/cli` is removed; unconditional `bundebug` is removed; lifecycle and command errors remain explicit and tested.

**Verification:** Migration command tests, SQLite/PostgreSQL migration suites, dependency graph, binary build and size comparison pass.

**Dependencies:** Task 3.3.
**Likely files:** migration command, database setup, tests, module files.
**Scope:** Medium batches.

### Task 3.5 — Remove proven dead code and normalize formatting

**Acceptance criteria:** Confirmed unused repository methods, handler code, frontend barrels, TODO-only UI, and stale config are removed; full Biome/Go formatting is applied in reviewable batches.

**Verification:** Reference/dead-code checks and full repository checks pass after each batch.

**Dependencies:** Tasks 3.2 and 3.4.
**Likely files:** bounded cleanup batches.
**Scope:** Repeatable medium batches.

### Checkpoint 3

- Frozen installs and module verification pass.
- Dependency graph matches the approved matrix.
- Before/after bundle and binary measurements are recorded.

## Phase 4 — Vertical frontend normalization

### Task 4.1 — Establish `app`, `features`, and `shared` boundaries

**Acceptance criteria:** App providers/router/runtime config and canonical `shared/api/generated`, `shared/config`, `shared/lib`, and `shared/ui` entry points exist without compatibility barrels or feature imports.

**Verification:** Boundary checks, typecheck, build, and base-path smoke pass.

**Dependencies:** Checkpoint 3.
**Likely files:** new boundary files and import rules.
**Scope:** Medium batches.

### Task 4.2 — Rebuild the canonical shadcn foundation

**Acceptance criteria:** Current Base Vega/Base UI components are generated into `shared/ui`; `@base-ui/react`, `react-day-picker`, and `tw-animate-css` provide the adopted component foundation; Tailwind v4 tokens/config follow current docs; adopted components use `data-slot` and retain product theming/accessibility.

**Verification:** shadcn diff/source review, component browser smoke, visual/layout checks, build, and dependency audit pass.

**Dependencies:** Task 4.1.
**Likely files:** components config, CSS, shared UI files in bounded batches.
**Scope:** Repeatable medium batches.

### Task 4.3 — Migrate the auth vertical slice

**Acceptance criteria:** Auth owns operations, state, screens, and redirects; app owns only provider/guard composition; protected content never flashes before auth resolution.

**Verification:** Auth unit/browser/e2e tests and base-path cookie matrix pass.

**Dependencies:** Tasks 4.1–4.2.
**Likely files:** auth feature batches, routes, app provider.
**Scope:** Repeatable medium batches.

### Task 4.4 — Migrate the sites vertical slice

**Acceptance criteria:** Sites owns list/create/switch operations and UI; login opens last accessible/only site or creation; multi-domain entry is first-class and validated; routes stay thin.

**Verification:** Site navigation, create/edit multi-domain, refresh, back/forward, and direct-route browser tests pass.

**Dependencies:** Task 4.3.
**Likely files:** sites feature batches and thin routes.
**Scope:** Repeatable medium batches.

### Task 4.5 — Migrate analytics and typed URL state

**Acceptance criteria:** Analytics owns operations/view models/screens; route validation returns final typed filters/date/pagination; prop chains are replaced by focused feature composition; URL behavior is deterministic.

**Verification:** Analytics unit/browser tests, generated types, and query baseline comparison pass.

**Dependencies:** Task 4.4.
**Likely files:** analytics feature batches and route schema.
**Scope:** Repeatable medium batches.

### Task 4.6 — Migrate site settings

**Acceptance criteria:** Settings has explicit routes and owns domain, tracking, blocking, GeoIP, and danger-zone behavior; it composes canonical shadcn components and preserves multi-domain/site-key semantics.

**Verification:** Settings mutation, failure, loading, navigation, and destructive-action browser tests pass.

**Dependencies:** Tasks 4.4–4.5.
**Likely files:** settings feature batches and routes.
**Scope:** Repeatable medium batches.

### Task 4.7 — Migrate event definitions

**Acceptance criteria:** Event definitions own operations, validation, screens, and mutation state; no imports reach another feature's internals.

**Verification:** Unit/browser/GraphQL e2e tests pass.

**Dependencies:** Tasks 4.5–4.6.
**Likely files:** event-definition feature batches and routes.
**Scope:** Repeatable medium batches.

### Task 4.8 — Enforce the no-flicker data lifecycle

**Acceptance criteria:** First load, background refresh, polling, mutation, errors, and stale data follow D-025 across all features; Apollo cache/refetch policy is explicit by feature; layout geometry remains stable.

**Verification:** Playwright transition assertions pass without screenshots masking transient blank states.

**Dependencies:** Tasks 4.3–4.7.
**Likely files:** feature query hooks/view models and browser tests in batches.
**Scope:** Repeatable medium batches.

### Task 4.9 — Delete legacy frontend structure

**Acceptance criteria:** Global `pages`, `components`, `hooks`, old layouts, broad barrels, and transitional imports are gone; only `app`, `features`, and `shared` remain at the architectural top level.

**Verification:** Boundary/reference checks, full frontend checks, browser suite, production build, and base-path matrix pass.

**Dependencies:** Tasks 4.3–4.8.
**Likely files:** legacy removals and final imports.
**Scope:** Repeatable medium batches.

### Checkpoint 4

- All critical admin flows use the accepted site-centric UX.
- shadcn and frontend dependencies match the approved current pattern.
- No legacy frontend ownership layer remains.

## Phase 5 — Vertical backend normalization

### Task 5.1 — Separate application composition and lifecycle

**Acceptance criteria:** `internal/app` explicitly composes configuration, database, features, transports, startup tasks, and cleanup; `cmd` only selects a command and runs it; no runtime container or globals exist.

**Verification:** Construction/lifecycle tests, all command builds, and shutdown/error tests pass.

**Dependencies:** Checkpoint 4.
**Likely files:** app composition, server command, focused tests.
**Scope:** Medium batches.

### Task 5.2 — Normalize auth feature ownership

**Acceptance criteria:** Auth owns tokens/credentials/application behavior and consumer-side interfaces; HTTP/GraphQL adapters own cookies and transport mapping; no Bun or GraphQL models leak into feature contracts.

**Verification:** Auth unit/e2e/cookie/base-path tests pass.

**Dependencies:** Task 5.1.
**Likely files:** auth feature and transport batches.
**Scope:** Repeatable medium batches.

### Task 5.3 — Normalize site feature and persistence mapping

**Acceptance criteria:** Site behavior, multi-domain invariants, and narrow storage interfaces are feature-owned; Bun rows stay in the adapter; mapping exists only where shapes/invariants differ.

**Verification:** Site service/repository/GraphQL and both-database tests pass.

**Dependencies:** Tasks 5.1–5.2.
**Likely files:** site feature and persistence batches.
**Scope:** Repeatable medium batches.

### Task 5.4 — Normalize analytics feature

**Acceptance criteria:** Collection, identity, session, filters, aggregates, and query contracts have explicit ownership; repository result types do not escape adapters; tracker behavior remains compatible.

**Verification:** Analytics unit/repository/e2e and tracker compatibility suites pass.

**Dependencies:** Tasks 5.1 and 5.3.
**Likely files:** analytics feature and persistence batches.
**Scope:** Repeatable medium batches.

### Task 5.5 — Normalize events and GeoIP features

**Acceptance criteria:** Event definitions and GeoIP each own application behavior and interfaces; download/lookup/persistence are adapters; error and best-effort semantics remain explicit.

**Verification:** Focused unit/integration/GraphQL tests pass.

**Dependencies:** Tasks 5.3–5.4.
**Likely files:** event and GeoIP feature batches.
**Scope:** Repeatable medium batches.

### Task 5.6 — Split GraphQL ownership by feature

**Acceptance criteria:** Feature SDL and resolver adapters map only feature contracts; generated gqlgen output remains isolated; handwritten resolvers are linted; transport models do not alias service/repository types.

**Verification:** Generate twice, lint handwritten resolvers, GraphQL e2e, and client codegen checks pass.

**Dependencies:** Tasks 5.2–5.5.
**Likely files:** schema/resolver batches and gqlgen config.
**Scope:** Repeatable medium batches.

### Task 5.7 — Normalize platform and seed ownership

**Acceptance criteria:** Database/config/bootstrap adapters are platform/app concerns; application-only public `pkg` code moves internal; example-data logic is reusable and command shells remain minimal.

**Verification:** Package boundary checks, seed command tests, builds, and migration suites pass.

**Dependencies:** Tasks 5.1–5.6.
**Likely files:** platform, seed, command batches.
**Scope:** Repeatable medium batches.

### Task 5.8 — Delete legacy backend layers

**Acceptance criteria:** Generic service/repository/model buckets and transitional aliases are removed when their feature migration completes; final dependency direction matches accepted architecture.

**Verification:** Import graph, full Go checks, e2e, migrations, and generation pass.

**Dependencies:** Tasks 5.2–5.7.
**Likely files:** legacy removals and imports in batches.
**Scope:** Repeatable medium batches.

### Checkpoint 5

- Feature ownership contains backend change blast radius.
- Persistence and transport types do not leak.
- No legacy backend layer or runtime DI mechanism remains.

## Phase 6 — API, data, and performance simplification

### Task 6.1 — Normalize GraphQL contracts and error categories

**Acceptance criteria:** Only used product fields remain; new display data has documented rationale; stable typed error categories map consistently to Apollo/UI behavior; schema and both generated clients evolve together.

**Verification:** Schema diff review, generation freshness, GraphQL e2e, frontend typecheck and browser failures pass.

**Dependencies:** Checkpoint 5.
**Likely files:** feature SDL/resolvers/operations and tests in vertical batches.
**Scope:** Repeatable medium batches.

### Task 6.2 — Reduce measured dashboard query fan-out

**Acceptance criteria:** Before/after query counts and timings prove each consolidation; data accuracy is unchanged or improved; no speculative aggregate abstraction is introduced.

**Verification:** Repository integration tests, representative benchmarks, and recorded query traces pass.

**Dependencies:** Tasks 1.7, 5.4, and 6.1.
**Likely files:** analytics query adapters and tests in batches.
**Scope:** Repeatable medium batches.

### Task 6.3 — Normalize polling, cache, pagination, and stale data

**Acceptance criteria:** One owner per polling cycle; cache keys and invalidation are explicit; pagination is typed and consistent; background updates satisfy D-025 without redundant requests.

**Verification:** Request-count assertions, browser transition tests, and performance comparison pass.

**Dependencies:** Tasks 4.8, 6.1, and 6.2.
**Likely files:** feature data hooks/view models and tests in batches.
**Scope:** Repeatable medium batches.

### Task 6.4 — Enforce measured performance guardrails

**Acceptance criteria:** Critical flow, query, bundle, and binary results meet Phase 1 budgets or have an approved evidence-backed exception; no size-only optimization harms clarity or correctness.

**Verification:** Re-run the baseline commands and record comparison in completion evidence.

**Dependencies:** Tasks 6.2–6.3.
**Likely files:** measurement config and master evidence.
**Scope:** Small.

### Checkpoint 6

- Data/API changes have measured correctness and performance evidence.
- No undocumented compatibility break exists.
- Performance guardrails pass.

## Phase 7 — Hardening, documentation, and release evidence

### Task 7.1 — Run security and privacy hardening

**Acceptance criteria:** Dependency vulnerabilities, auth/cookie policy, tracker input, CORS/CSRF, secrets, analytics privacy, and database handling are reviewed; validated findings are fixed with tests.

**Verification:** `govulncheck`, linters, focused security tests, and documented review pass.

**Dependencies:** Checkpoint 6.
**Likely files:** focused fixes/tests and security documentation.
**Scope:** Repeatable medium batches.

### Task 7.2 — Run the complete release matrix

**Acceptance criteria:** Frozen install, format/lint/types, generation freshness, unit/browser/e2e, race-relevant Go tests, builds, SQLite/PostgreSQL migrations, seed command, and runtime base-path matrix all pass from a clean state.

**Verification:** Commands and results are recorded verbatim in completion evidence.

**Dependencies:** Task 7.1.
**Likely files:** master completion evidence only unless a failure exposes a defect.
**Scope:** Medium.

### Task 7.3 — Finalize maintained documentation

**Acceptance criteria:** Final stack, module ownership, API/data behavior, configuration migrations, base-path examples, dependency rationale, development commands, and operational behavior are current; obsolete docs are superseded; non-obvious invariants are commented at enforcement boundaries.

**Verification:** Link/reference review and fresh-maintainer walkthrough pass.

**Dependencies:** Task 7.2.
**Likely files:** README, dashboard/server docs, ADRs, master plan.
**Scope:** Repeatable medium batches.

### Task 7.4 — Close the normalization update

**Acceptance criteria:** Every definition-of-done item links evidence; no temporary bridge, legacy layer, stale dependency, unintended worktree change, or unresolved blocking clarification remains.

**Verification:** Final worktree audit and maintainer review.

**Dependencies:** Tasks 7.2–7.3.
**Likely files:** master plan and task checklist.
**Scope:** Small.

### Final checkpoint

- Full definition of done is satisfied with linked evidence.
- The final implementation matches all accepted decisions D-001 through D-029.
- The maintainer receives a concise before/after stack and migration summary.

## Primary risks and mitigations

| Risk | Impact | Mitigation |
| --- | --- | --- |
| Big-bang cross-layer rewrite | High | Vertical feature slices, five-file handwritten batches, green checkpoints. |
| Mistaking design improvement for regression | High | Apply D-007 classification and explicit before/after acceptance tests. |
| Runtime subpath or cookie regression | High | One-build base-path and multi-instance auth matrix before structural changes. |
| Data loss or changed analytics meaning | High | Data-preserving migrations, both-database tests, aggregate multi-domain invariant. |
| Latest-version incompatibility | Medium | Latest stable only, official migration notes, ecosystem-bounded upgrade batches. |
| UI flicker hidden by happy-path tests | High | Transition assertions in real browsers; first-load/background/mutation states tested separately. |
| Generated code masking handwritten defects | Medium | Isolated generated paths, stale checks, lint all handwritten resolver code. |

## Open questions

None blocking. Phase 1 must replace measurement placeholders with numeric guardrails. Any new contract or dependency outside this plan reopens planning and requires explicit approval.
