# Lovely Eye architecture normalization

## Document status

- Status: Implementation complete — awaiting maintainer review
- Implementation approval: Granted on 2026-08-06 for the accepted phase plan and dependency set
- Last updated: 2026-08-17
- Owner: Project maintainer
- Purpose: Living source of truth for the full frontend and backend normalization update

This document stores confirmed intent, constraints, decisions, open clarifications, audit evidence, planned work, and completion evidence. It is intentionally updated throughout the interview and implementation process.

Execution details and the derived checklist live in [`tasks/plan.md`](../../tasks/plan.md) and [`tasks/todo.md`](../../tasks/todo.md). They are subordinate to this document and introduce no independent architectural decisions.

## Operating rules

1. Discuss one clarification at a time using the interview format: hypothesis, confidence, one question, and a stated guess.
2. Record a user answer in the decision log only after explicit confirmation.
3. Keep proposals labeled `Proposed`; do not treat them as requirements.
4. Separate behavior-preserving structural work from behavior changes.
5. Complete and verify each phase before starting a dependent phase.
6. Do not perform a big-bang rewrite. Every implementation step must remain reviewable and leave the repository green.
7. Do not add a dependency until the decision log explains why the existing stack or standard library is insufficient.
8. Do not remove uncertain code or dependencies without confirming references, runtime behavior, and tests.
9. Attach verification evidence when changing an item to `Done`.
10. Create a dedicated ADR for any accepted decision that is expensive to reverse; link it from this document.
11. Do not classify every difference from historical behavior as a regression. Preserve confirmed value and invariants, not accidental behavior.
12. Document architectural reasons and operational consequences. Add code comments where they protect a non-obvious invariant or explain why the obvious alternative is unsafe; do not narrate self-explanatory code.

## Confirmed intent

Intent: Rebuild Lovely Eye's foundations around a reliable personal admin panel, strict end-to-end correctness, current tooling, understandable data flow, and feature ownership so future expansion is safe instead of fragile.

Confidence: 100% — explicitly confirmed by the user on 2026-08-06. Target architecture and implementation details remain separate decisions.

- Outcome: A fully normalized frontend and backend with more accurate data and understandable ownership.
- User: The maintainer's personal, stable, and convenient admin panel.
- Why now: Chaotic boundaries, unnecessary dependencies, and insufficient quality guarantees make further development risky.
- Success: No known bugs or flicker; strict types; current dependencies; modular UI; idiomatic Go; safe quality gates.
- Binding constraints: Preserve all data, arbitrary runtime `BASE_PATH`, direct nested-route loading, and approved migration control for external breaking changes.
- Out of scope: Permanent legacy layers and a public/private-token external API in this update.

## Confirmed goals

- Normalize the entire frontend and backend, not only individual files.
- Make module ownership and dependency direction clear.
- Remove unnecessary direct dependencies and accidental technologies.
- Make the actual stack concise and understandable.
- Make project rules enforceable by tooling rather than documentation alone.
- Preserve arbitrary runtime `BASE_PATH` support with one static frontend build.
- Use this document and one-question-at-a-time interviews for discovery and decisions.
- Deliver a personal admin panel without known bugs, conflicting state, layout instability, or loading flicker.
- Make future feature development safe through explicit quality guarantees and regression gates.
- Use strict typing throughout the application so invalid states and data-shape mistakes are caught at their owning boundary.
- Update libraries and frameworks to current versions with the required migrations and proper refactoring, not compatibility shims.
- Prefer Rust- and SWC-based tooling when it fits the requirement and reduces complexity.
- Use deterministic code generation when it removes repetitive code or strengthens type safety; do not reject it merely because generated files exist.
- Improve analytics data accuracy through a more understandable data and responsibility architecture.
- Preserve first-class support for assigning multiple domains to one site.
- Make the admin UI more convenient and coherent.
- Rebuild the shadcn/ui integration from current recommended patterns and refresh all adopted components accordingly.
- Use current shadcn/ui components and styling patterns as the default wherever they satisfy the product need; custom UI is the exception.
- Leave no legacy architecture or transitional compatibility layers in the final state.
- Use clean idiomatic Go and explicit frontend feature modularity.
- Keep architecture, migration, and operational documentation current, including rationale and focused examples. Explain non-obvious invariants at the relevant code boundary where documentation alone would not prevent misuse.

## Confirmed invariants

- One dashboard artifact must work at `/`, one-level subpaths, and nested runtime subpaths without rebuilding.
- Database migrations must preserve all existing analytics, configuration, and user data; deleting data or requiring reimport is not allowed.
- Application routes must not manually prepend `BASE_PATH`.
- Project tooling uses Bun; do not introduce Node/npm workflows.
- Go remains the backend language and React with TypeScript remains the dashboard foundation unless a later explicit decision supersedes this.
- No Git commit is created without an explicit user request.
- Generated files are not hand-edited.

See [ADR-0001](../decisions/0001-runtime-base-path.md) for the accepted runtime base-path decision.

## Remaining architecture and delivery decisions

No blocking architecture clarification remains. Phase 6 resolved the chunk-budget mismatch by enforcing stable transfer boundaries: initial JavaScript, analytics increment, total JavaScript, total CSS, and Go binaries. Login/register/sites route attribution remains observational because Vite assigns shared shadcn chunks to whichever lazy route first owns them; eagerly moving those chunks would make the number smaller without reducing transferred code. The CSS ceiling is 12,500 gzip bytes, leaving measured headroom above the stable 11,412-byte shadcn/Tailwind v4 output while still blocking material growth.

## Change classification

Use these categories during planning, review, and testing:

- **Behavior-preserving refactor:** the accepted contract and observable outcomes remain unchanged while implementation and ownership improve.
- **Design improvement:** an old contract is deliberately replaced because the new design improves correctness, data accuracy, usability, simplicity, or maintainability. It requires a written rationale, explicit acceptance criteria, before/after evidence, and tests for the new contract. Obsolete characterization tests are intentionally replaced, not silently weakened.
- **Regression:** an unintended loss of a confirmed invariant, accepted outcome, compatibility commitment, or user value. Examples include incorrect data, flicker, broken navigation, lost functionality, or an undocumented contract break.

Historical behavior is evidence, not automatically a requirement. Before calling a difference a regression, identify the accepted contract or value it violates. Before calling it an improvement, demonstrate the intended new contract and compare correctness, data accuracy, UX, simplicity, and compatibility.

API contracts may change when the product needs new display data or when a demonstrably better design replaces the old contract. “Best practice” alone is not sufficient evidence: the proposal must identify the concrete correctness, type-safety, data-accuracy, usability, or maintainability gain. Schema, generated types, application consumers, documentation, and tests change together. Compatibility and rollout still follow the separately accepted compatibility policy.

The current GraphQL API is an internal dashboard/backend contract with no known external consumers. Backend schema and bundled dashboard consumers may therefore evolve atomically without preserving obsolete GraphQL shapes. A possible future private-token API is not a current compatibility requirement; if introduced, its externally consumable contract and versioning policy must be designed and accepted explicitly.

The tracker payload or script API may receive a breaking change, but implementation is blocked until the user explicitly approves a migration plan. That plan must define affected integrations, new contract, rollout and coexistence period, snippet update path, data-loss risk, verification, and rollback. Until then, existing tracker behavior remains the compatibility baseline.

Environment variables may be renamed or removed without permanent legacy aliases. Every such change must include migration documentation and a precise startup error that identifies the obsolete variable and its replacement or removal. Any temporary alias must have an explicit removal point.

Internal dashboard route URLs may change when a better routing or UX structure requires it; legacy redirects are not required. Runtime `BASE_PATH`, direct loading and refresh of nested routes, and externally used tracker URLs remain protected contracts unless their separate migration policy is followed.

## Baseline evidence

Recorded during the 2026-08-05 audit:

- Current dashboard lint, source-size, and TypeScript gate passes.
- Full `biome check` reports 192 formatter and import-organization errors that the current gate does not enforce.
- No frontend unit or component tests exist.
- Nineteen handwritten frontend files exceed the documented 120-line limit.
- The source-size gate allows 220 lines per file and 160 lines per function.
- Go tests pass and golangci-lint reports zero issues.
- Go coverage is uneven: `graph`, `server`, `dashboard`, and `database` are at 0%; repository is 16.0%; services are 46.8%.
- [x] `go mod tidy -diff` reports no unused Go module requirements.
- The current production Vite configuration emits root-absolute asset URLs and does not yet satisfy ADR-0001.
- GraphQL currently serves one bundled dashboard through 11 queries and 10 mutations. The frontend uses fixed generated operations rather than runtime-composed queries, while Apollo defaults mostly favor network refreshes and explicit refetches over normalized-cache behavior.
- The GraphQL stack adds schema and client code generation plus approximately 11,230 lines of generated Go executor code. The user explicitly accepted GraphQL with `gqlgen`; generated size alone is not a removal reason. Frontend client-library cost remains independently reviewable.
- A site already accepts multiple allowed domains through one public site key. Analytics records are site-scoped and do not retain a domain dimension, so all configured domains intentionally contribute to one aggregate analytics property.

## Known findings inventory

### Correctness and regression risks

- [x] Fix runtime subpath assets, component asset URLs, and deep-link coverage.
- [x] Stop discarding repository errors during site domain uniqueness checks.
- [x] Distinguish not-found errors from infrastructure/database failures.
- [x] Stop discarding GeoIP synchronization and refresh errors.
- [x] Return persisted user creation timestamps instead of request-time timestamps.
- [x] Implement accepted base-path-specific auth cookie names and paths for same-origin instance isolation.

### Frontend structure

- [x] Remove the duplicate `routes` plus `pages` ownership model.
- [x] Remove the router-to-route-to-page-to-router dependency loop.
- [x] Split analytics, site creation, and site settings into explicit routes.
- [x] Make route search validation return final typed values instead of strings and `Record<string, unknown>` casts.
- [x] Centralize authentication routing decisions in one boundary.
- [x] Group auth, sites, analytics, site settings, and event definitions by feature.
- [x] Move GraphQL operations to their owning features and keep generated artifacts isolated.
- [x] Replace large dashboard prop chains with feature-level view models or focused composition.
- [x] Use one normalized runtime configuration API instead of direct `window.__ENV__` access.
- [x] Remove unused global barrel files and reduce broad UI exports.
- [x] Define justified exceptions for generated and third-party-derived UI files.
- [x] Rebuild the shadcn/ui foundation from current recommended installation and component patterns.
- [x] Re-adopt only the shadcn components the product actually uses; do not carry unused component surface forward.
- [x] Define and test stable loading, stale-data, navigation, mutation, and layout behavior with no avoidable flicker.
- [x] Review the admin information architecture and interaction flow for personal daily use.

### Backend structure

- [x] Replace generic layer buckets with clear feature ownership where it reduces cross-package changes.
- [x] Stop GraphQL models from depending on service or repository types.
- [x] Stop application services from exporting repository aliases and persistence result types.
- [x] Separate HTTP cookie handling from the authentication application service.
- [x] Split GraphQL schema and resolver ownership by feature and ensure handwritten resolver logic is linted.
- [x] Separate application bootstrap, dependency construction, migration, startup tasks, and HTTP route construction.
- [x] Make configuration parsing explicit, validated, and error-returning.
- [x] Apply accepted selective persistence mapping at feature boundaries without mechanical model duplication.
- [x] Move application-only packages out of public `pkg/` unless a public API is intentional.
- [x] Move example-data business logic out of `cmd/load-example-data/main.go`.

### Dead code and dependency cleanup

- [x] Verify and remove unused repository CRUD methods and `AnalyticsHandler.Event`.
- [x] Remove unused frontend barrel files.
- [x] Remove confirmed redundant frontend packages:
  - `@vitejs/plugin-react`
  - `@tailwindcss/postcss`
  - direct `postcss`
  - `autoprefixer`
  - direct `@tanstack/router-generator`
- [x] Replace the direct `urfave/cli/v3` migration dependency with explicit standard-library dispatch; its remaining indirect edge belongs to Atlas tooling.
- [x] Remove `bundebug` from the production server binary.
- [x] Isolate Atlas/provider imports behind the `tools` build tag while retaining their required migration-schema capability.
- [x] Retain justified peer, test, and tooling dependencies whose owning capability remains.

### Data and API behavior

- [x] Remove unused site loading from the authentication query.
- [x] Define stable GraphQL error categories and frontend handling.
- [x] Evolve the internal GraphQL schema atomically with the bundled dashboard under the accepted compatibility policy.
- [x] Measure dashboard database query fan-out before consolidating aggregate queries.
- [x] Review polling ownership, cache policy, pagination, and stale-data behavior by feature.
- [x] Preserve tracker compatibility unless an explicit breaking-change migration plan is approved before implementation.
- [x] Preserve aggregate multi-domain collection: every configured domain is authorized for the same site key and contributes to the same site-level analytics dataset.

### Quality gates

- [x] Replace lint-only frontend validation with a non-mutating full Biome check.
- [x] Align documented and enforced file/function-size policies.
- [x] Separate browser and tooling TypeScript configurations.
- [x] Enable accepted dashboard `exactOptionalPropertyTypes` while keeping `skipLibCheck` limited to external declarations.
- [x] Ensure handwritten gqlgen resolver implementations do not escape lint rules as generated code.
- [x] Use executable ownership/freshness gates plus compiler, import/reference, and module-graph review; do not claim an unsound universal dead-code oracle.
- [x] Validate untrusted runtime data at boundaries rather than relying on TypeScript assertions.
- [x] Establish a strict-type policy for generated code, application code, GraphQL, configuration, and persistence mappings.
- [x] Add browser-level regression coverage for layout stability, navigation, stale-data refresh, and critical admin workflows.

## Dependency disposition plan

Audit snapshot: 2026-08-06. Versions are candidates observed from the official npm registry and `go list -u -m`; implementation must re-check the latest stable release and its official migration notes immediately before changing a lockfile. This section approves the planned dependency set, not an unreviewed bulk upgrade.

### Frontend

| Capability | Keep or add | Remove or replace | Reason |
| --- | --- | --- | --- |
| React runtime | `react`, `react-dom` | Direct `react-is` unless a fresh build proves a direct API requirement | React remains the UI foundation; undeclared transitive needs stay transitive. |
| GraphQL client | `@apollo/client`, `graphql`, generated client preset; add direct `rxjs` | — | Apollo remains useful for shared query lifecycle and stable refresh behavior. Its official installation requires `rxjs` as a top-level peer. |
| Routing | `@tanstack/react-router`, `@tanstack/router-plugin` | Direct `@tanstack/router-generator` | The Vite plugin owns file-route generation; the generator need not also be a direct dependency. |
| Build and CSS | `vite`, `@vitejs/plugin-react-swc`, `tailwindcss`, `@tailwindcss/vite`, `lightningcss` | `@vitejs/plugin-react`, `@tailwindcss/postcss`, `postcss`, `autoprefixer` | Keep the active SWC/Vite and Tailwind Vite paths; remove unused parallel Babel/PostCSS paths. |
| shadcn foundation | Add pinned `shadcn`, unified `radix-ui`, `react-day-picker`, and `tw-animate-css`; keep `class-variance-authority`, `clsx`, `tailwind-merge`, and `lucide-react` | Individual `@radix-ui/react-*`, `@daypicker/react`, `tailwindcss-animate` | Rebuild from the current shadcn Vite/new-york/Radix source, using its unified Radix package and current animation/calendar dependencies. |
| Product UI/data | `date-fns`, `recharts`, `zod` | — | Each owns a used product capability: dates, charts, and runtime boundary validation. |
| Code quality and generation | `@biomejs/biome`, `@graphql-codegen/cli`, `@graphql-codegen/client-preset`, TypeScript and React/Node type packages | — | These enforce the accepted strictness and deterministic generation policies. |
| Frontend tests | Bun's built-in test runner for pure units; add `@playwright/test` for browser contracts | Do not add Vitest initially | One unit runner plus real-browser coverage satisfies the risk-based policy with less tooling overlap. |

Registry candidates at audit time:

- Core/API: React and React DOM `19.2.8`, Apollo `4.2.10`, GraphQL `17.0.2`, RxJS `7.8.2`, GraphQL Codegen CLI `7.2.0`, client preset `6.1.1`.
- Routing/build: TanStack Router `1.170.20`, router plugin `1.168.25`, Vite `8.2.0`, React SWC plugin `4.3.3`, TypeScript `7.0.2`, Biome `2.5.7`, Tailwind CSS and its Vite plugin `4.3.3`, Lightning CSS `1.33.0`.
- UI: shadcn `4.16.1`, unified Radix UI `1.6.7`, React DayPicker `10.0.1`, `tw-animate-css` `1.4.0`, Lucide React `1.28.0`, Recharts `3.10.1`.
- Browser tests: Playwright `1.62.1`.

Migration evidence:

- [shadcn Vite installation](https://ui.shadcn.com/docs/installation/vite)
- [shadcn Tailwind v4 migration](https://ui.shadcn.com/docs/tailwind-v4)
- [shadcn unified Radix migration](https://ui.shadcn.com/docs/changelog/2026-02-radix-ui)
- [shadcn `components.json`](https://ui.shadcn.com/docs/components-json)
- [TanStack Router manual setup](https://tanstack.com/router/latest/docs/installation/manual)
- [Vite official plugins](https://vite.dev/plugins/)
- [TypeScript 7 announcement](https://devblogs.microsoft.com/typescript/announcing-typescript-7-0/)
- [Apollo Client installation](https://www.apollographql.com/docs/react/get-started)

Dependency refresh snapshot: 2026-08-15. The approved stack was rechecked against the npm and Go
registries before changing either lockfile. Retained frontend packages were updated to their latest
stable releases, including Apollo Client `4.2.12`, TanStack Router `1.170.29` and plugin `1.168.32`,
Lucide `1.31.0`, Biome `2.5.8`, GraphQL Codegen client preset `6.1.3`, shadcn `4.18.0`, and Vite
`8.2.1`. `shadcn diff` reports no component-source updates. The client-preset regeneration removes
its obsolete partial-fragment overload without requiring an application compatibility layer. A
fresh lockfile regeneration resolves both the direct `^6.1.3` requirement and the CLI's `^6.1.0`
requirement to one `6.1.3` package without an override.

The full Bun audit exposed newly patched build-tool edges. The frozen lockfile now resolves their
patched versions through the upstream-compatible ranges without project-level overrides. One
dev-only `micromatch -> picomatch@2.3.1` advisory remains documented in
[`docs/security.md`](../security.md): Bun 1.3.14 does not apply its documented scoped override forms,
a direct pin does not replace the nested edge, and a global override would incorrectly force v2 on
v4 consumers. Production dependencies remain clean.

### Go

| Capability | Keep or add | Remove or replace | Reason |
| --- | --- | --- | --- |
| Schema and migrations | Atlas, Atlas Bun provider, Bun and the PostgreSQL/SQLite dialects and drivers | — | Both supported databases and data-preserving migrations remain invariants. Atlas/provider ownership must stay tooling-only. |
| GraphQL | `gqlgen`, `gqlparser`, and test-only `genqlient` | Legacy blank-import `tools.go` pattern | GraphQL/gqlgen is accepted; Go 1.26 tool dependencies belong in `go.mod` `tool` directives. |
| Auth and validation | `jwt/v5`, `x/crypto` | — | Used for signed tokens, password hashing, and analytics identity cryptography. |
| Analytics dimensions and GeoIP | `useragent`, `geoip2-golang/v2`, `maxminddb-golang/v2` | — | These provide active tracker/analytics behavior not covered by the standard library. |
| Database runtime | Bun runtime modules and `modernc.org/sqlite` | `bundebug` from the production path | The current debug hook is always installed despite having no runtime debug policy. Add structured query diagnostics only if the observability design demonstrates a need. |
| Migration CLI | Standard-library command dispatch | `urfave/cli/v3` | Four fixed subcommands do not justify a runtime CLI framework. |
| Tests and tools | `testify`; add a pinned `govulncheck` tool; declare gqlgen/genqlient executables with `tool` directives | — | Preserve readable tests and make generation/security checks reproducible. |

Go upgrade candidates at audit time:

- Atlas `1.2.2` → `1.3.0`.
- gqlgen `0.17.92` → `0.17.94`.
- maxminddb `2.4.0` → `2.4.1`.
- gqlparser `2.5.35` → `2.5.36`.
- `golang.org/x/crypto` `0.53.0` → `0.54.0`.
- `modernc.org/sqlite` `1.53.0` → `1.56.0`.
- Remove `urfave/cli/v3`; other direct modules were already latest stable in the audit snapshot.

Use `go mod tidy`, `go mod verify`, full tests, migration tests, and `govulncheck ./...` after the dependency change. Go's official module guidance documents both update discovery and Go 1.24+ [`tool` directives](https://go.dev/doc/modules/managing-dependencies#tools).

The 2026-08-15 refresh moved the minimum toolchain to security-fixed Go `1.26.6`, updated GeoIP2 to
`2.3.0`, MaxMind DB to `2.5.0`, `x/crypto` to `0.55.0`, the pinned `govulncheck` tool to `1.7.0`,
and the static esbuild module to `0.28.2`. `go get -t -u ./...`, explicit tool updates, and
`go mod tidy` refreshed retained indirect requirements without promoting transitive modules to the
application API. The paired GeoIP update preserves the existing reader API; cache policy remains a
separate measured decision. The first scan on Go 1.26.5 found six reachable standard-library
advisories fixed by 1.26.6; the refreshed scan reports zero reachable vulnerabilities.

## Accepted target shape

### Frontend

Status: Implemented and protected by automated boundary checks.

```text
dashboard/src/
  app/
    providers/
    router/
    routes/
  features/
    auth/
    sites/
    analytics/
    site-settings/
    event-definitions/
  shared/
    api/generated/
    config/
    ui/
    lib/
```

Boundary proposal:

- Routes own path definitions, typed URL parsing, guards, loaders, and feature composition only.
- Features own GraphQL operations, business state, screens, and feature-specific UI.
- Shared modules contain only code with multiple real feature consumers.
- Features do not import concrete route definitions or other feature internals.
- Generated GraphQL code is infrastructure, not a feature ownership boundary.
- Generated TanStack route-tree files are never hand-edited.
- `shared/ui` contains the canonical shadcn source components; features compose them directly.
- Imports from shadcn's underlying primitive libraries stay inside `shared/ui`.
- A project-specific wrapper is justified only by repeated product behavior or policy, not by renaming or hiding a shadcn API.

### Backend

Status: Feature-oriented modular-monolith architecture, explicit idiomatic Go composition, and selective persistence mapping accepted.

```text
server/internal/
  app/
  auth/
  site/
  analytics/
  event/
  geoip/
  transport/
    graphql/
    http/
    dashboard/
  platform/
    config/
    database/
  seed/
```

Boundary proposal:

- Use a feature-oriented modular monolith with explicit constructor composition at the application root; wiring may be handwritten or generated according to complexity.
- Feature packages own application types and narrow consumer-side storage interfaces.
- Transport packages map HTTP/GraphQL values to feature inputs and outputs.
- Persistence row structs and framework-specific tags stay inside persistence adapters. Features introduce their own types only when business invariants or a materially different contract require them; do not mechanically duplicate every database model.
- Persistence types never leak into feature contracts, GraphQL, or other transport boundaries.
- `cmd` packages only load configuration, build the application, and run a command.
- Do not add a clean-architecture framework or runtime DI container. Code generation remains allowed under the project generation policy.

## Phase 1 performance baseline and guardrails

The reproducible command is `./scripts/measure-baseline.sh`. Run it from the repository root inside the Dev Container with no other project workload. It performs a fresh dashboard build, compressed manifest analysis, trimmed Go builds, ten serial Go benchmark samples, and one real Chromium critical-flow sample. Two consecutive runs on 2026-08-06 used Go 1.26.4, Bun 1.3.14, Linux arm64/aarch64, and the same container state.

| Measurement | Run 1 | Run 2 | Phase 6 guardrail |
| --- | ---: | ---: | ---: |
| Initial dashboard, gzip | 241,209 B | 241,209 B | <= 265,000 B |
| Login route increment, gzip | 1,074 B | 1,074 B | Observe; Vite shared-chunk attribution is not a transfer boundary |
| Register route increment, gzip | 1,221 B | 1,221 B | Observe; Vite shared-chunk attribution is not a transfer boundary |
| Sites route increment, gzip | 1,620 B | 1,620 B | Observe; Vite shared-chunk attribution is not a transfer boundary |
| Site analytics route increment, gzip | 123,334 B | 123,334 B | <= 136,000 B |
| All dashboard JavaScript, gzip | 362,261 B | 362,261 B | <= 400,000 B |
| All dashboard CSS, gzip | 8,966 B | 8,966 B | <= 12,500 B after the current shadcn/Tailwind v4 foundation |
| Server binary | 26,444,141 B | 26,444,141 B | <= 29,100,000 B |
| Migration binary | 19,456,283 B | 19,456,283 B | <= 21,400,000 B |
| Login-to-ready, Chromium | 857 ms | 866 ms | <= 1,500 ms |
| Dashboard reload-to-data-ready, Chromium | 902 ms | 873 ms | <= 1,500 ms |
| Dashboard GraphQL operations | 9 | 9 | <= 9; Task 6.2 must reduce rather than merely preserve this fan-out |
| SQLite visitor-count read | 25,695–26,160 ns/op | 25,841–26,787 ns/op | <= 35,000 ns/op, <= 6,000 B/op, <= 30 allocs/op |
| SQLite top-pages read | 70,370–72,354 ns/op | 70,915–73,009 ns/op | <= 95,000 ns/op, <= 8,600 B/op, <= 80 allocs/op |

The browser and database timing limits are investigation thresholds for comparable local/container conditions, not claims about production latency. A threshold failure requires a same-environment repeat and an evidence-backed explanation or correction; noisy timing alone must not motivate a clarity-reducing optimization. Task 6.4 makes bundle, binary, browser GraphQL request-count, and server SQL fan-out limits executable blocking checks. Route-attribution increments and timing samples remain diagnostic because they are sensitive to chunk ownership and host load rather than transferred capability or data correctness.

## Accepted phased roadmap

The ordering follows the accepted dependency graph and rollout policy. Implementation remains blocked until the detailed task plan receives explicit approval.

### Phase 0 — Plan and decision freeze

Entry: Audit, intent interview, target architecture, compatibility policy, dependency disposition, and detailed task plan exist.

- Convert accepted high-cost decisions into ADRs.
- Refresh registry versions and official migration sources immediately before implementation.
- Obtain explicit approval for the phase plan and planned direct dependency changes.

Exit: ADRs are linked, no blocking clarification remains, and implementation is explicitly approved.

### Phase 1 — Safety net and enforceable baseline

- Add runtime base-path integration coverage.
- Add frontend smoke/component tests for auth, routing, sites, analytics state, and settings.
- Add backend characterization tests around GraphQL, server construction, dashboard serving, and configuration.
- Make formatting, import organization, generation consistency, types, lint, and tests enforceable gates.

Exit: Existing behavior is covered sufficiently for structural changes, and all gates are green.

### Phase 2 — Correctness fixes

- Resolve the known swallowed-error, timestamp, GeoIP, and base-path defects.
- Keep behavior changes isolated from structural moves.
- Add one regression test per corrected defect.

Exit: Known correctness findings are closed with tests.

### Phase 3 — Mechanical cleanup

- Apply formatter/import organization.
- Remove confirmed dead code, unused barrels, and redundant dependencies.
- Align TypeScript and Go quality policies with documented rules.
- Regenerate and verify artifacts.

Exit: Smaller dependency graph and no known dead compatibility shims.

### Phase 4 — Frontend normalization

- Establish `app`, `features`, and `shared` boundaries.
- Normalize routes and typed URL state.
- Move auth, sites, analytics, settings, and event definitions one feature at a time.
- Simplify data orchestration and component contracts after each move.
- Rebuild the shared shadcn/ui layer from the accepted current reference implementation and migrate consumers deliberately.
- Verify critical admin workflows and loading transitions in a real browser after each feature migration.

Exit: A feature can be understood and changed without traversing generic global folders.

### Phase 5 — Backend normalization

- Separate bootstrap and transport concerns.
- Establish feature-owned types and interfaces.
- Move GraphQL resolver ownership by feature.
- Remove repository types from service and transport contracts.
- Relocate application-only `pkg` code and seed logic.

Exit: Dependency direction is explicit and feature changes have a contained package blast radius.

### Phase 6 — API, data, and performance simplification

- Remove redundant GraphQL fields/queries only under the compatibility decision.
- Measure and reduce database query fan-out where justified.
- Normalize caching, polling, pagination, and error semantics.
- Measure bundle and binary impact before and after dependency changes.

Exit: Targets defined during discovery are met without undocumented compatibility breaks.

### Phase 7 — Final hardening and release

- Run the complete test, lint, generation, build, migration, and runtime-path matrix.
- Review security, privacy, tracker behavior, SQLite/PostgreSQL parity, and deployment docs.
- Update architecture documentation and supersede outdated ADRs instead of deleting history.
- Produce a final before/after stack and dependency inventory.

Exit: All definition-of-done items have linked evidence and no unresolved blocking clarification remains.

## Clarification record

All blocking discovery clarifications are answered. Reopen an item only when new evidence invalidates an accepted premise.

| ID | Topic | Status | Blocks |
| --- | --- | --- | --- |
| C-001 | Primary outcome and optimization target | Answered | All architecture decisions |
| C-002 | User-visible behavior changes and authority to replace contracts | Answered | Correctness and frontend phases |
| C-003 | Compatibility commitments | Answered | GraphQL, migrations, config, tracker |
| C-004 | Frontend ownership boundaries | Answered | Frontend normalization |
| C-005 | Backend architecture and composition style | Answered | Backend normalization |
| C-006 | ORM model versus domain model boundary | Answered | Backend normalization |
| C-007 | GraphQL scope and future | Answered | API normalization |
| C-008 | Strictness and CI policy | Answered | Safety net and quality gates |
| C-009 | Test depth and acceptable refactor risk | Answered | Safety net |
| C-010 | Performance and size targets | Answered | Optimization |
| C-011 | Dependency freshness, Rust/SWC preference, and admission/removal policy | Answered | Dependency cleanup |
| C-012 | Rollout, migration, and stopping point | Answered | Release |
| C-013 | Admin UX stability and flicker acceptance | Answered | Frontend safety net and normalization |
| C-014 | Runtime-base-path auth cookie isolation | Answered | Auth, runtime base path, browser coverage |

## Decision log

| ID | Decision | Status | Date | Evidence or rationale |
| --- | --- | --- | --- | --- |
| D-001 | One static dashboard build supports arbitrary runtime base paths. | Accepted | 2026-08-05 | [ADR-0001](../decisions/0001-runtime-base-path.md) |
| D-002 | Use feature-owned frontend modules with thin TanStack Router route files limited to paths, typed URL state, guards, loaders, and composition. Generated route trees are never hand-edited. | Accepted principle | 2026-08-06 | [ADR-0002](../decisions/0002-feature-owned-frontend.md); top-level structure is recorded in D-015. |
| D-003 | Use a feature-oriented Go modular monolith. Feature packages own application behavior and consumer-side interfaces; transport and platform code stay at the edges; `cmd` only composes and runs the application. Do not impose formal Clean Architecture layers. | Accepted architecture | 2026-08-06 | [ADR-0003](../decisions/0003-feature-oriented-go-modular-monolith.md); composition style is recorded in D-016. |
| D-004 | Preserve GraphQL as the dashboard API and use `gqlgen` for the Go implementation. Keep the schema and generated contracts strict, deterministic, and feature-owned at the handwritten source boundary. Frontend GraphQL client choice remains independently reviewable. | Accepted architecture | 2026-08-06 | [ADR-0004](../decisions/0004-graphql-and-code-generation.md). |
| D-005 | Keep persistence row structs and ORM/database tags inside persistence adapters. Introduce separate feature types only when business invariants or a materially different contract justify mapping; do not mechanically duplicate every database model. Persistence types must not leak into feature, GraphQL, or transport contracts. | Accepted architecture | 2026-08-06 | Confirmed during the C-006 interview as selective mapping. |
| D-006 | Quality includes a stable personal admin UI, strict type safety, current dependencies, better data accuracy, refreshed shadcn/ui, idiomatic Go, and no final legacy layer. | Accepted outcome | 2026-08-05 | Confirmed during C-001 interview response. |
| D-007 | An intentional, evidenced replacement of an old contract with a better design is a design improvement; a regression is an unintended loss of a confirmed invariant, outcome, compatibility commitment, or user value. | Accepted principle | 2026-08-05 | Confirmed during the C-002 interview. |
| D-008 | API changes are allowed for required display data or a demonstrably better design when the concrete gain is documented and schema, generated types, consumers, documentation, and tests evolve together. | Accepted principle | 2026-08-05 | Confirmed during the C-002 interview; rollout remains subject to C-003. |
| D-009 | The current GraphQL API has no external consumers and may evolve atomically with the bundled dashboard. A future private-token API will require an explicit external contract and versioning decision when introduced. | Accepted principle | 2026-08-05 | Confirmed during the C-003 interview. |
| D-010 | Tracker payload and script APIs may break only after explicit user approval of a documented migration and rollback plan. Existing tracker behavior remains compatible until that approval. | Accepted principle | 2026-08-06 | Confirmed during the C-003 interview. |
| D-011 | Database migrations must preserve all existing data. Destructive migration or mandatory data reimport is not allowed. | Accepted invariant | 2026-08-06 | Confirmed during the C-003 interview. |
| D-012 | Environment variables may be renamed or removed without permanent aliases when migration documentation and a precise startup diagnostic identify the required operator action. | Accepted principle | 2026-08-06 | Confirmed during the C-003 interview. |
| D-013 | Internal dashboard URLs may change without legacy redirects. Runtime base-path behavior, nested-route refresh, and external tracker URLs remain protected contracts. | Accepted principle | 2026-08-06 | Confirmed during the C-003 interview. |
| D-014 | Current shadcn/ui components and styling patterns are the default UI foundation. Canonical components live in `shared/ui`; features compose them directly; underlying primitives stay inside that boundary; wrappers require repeated product-specific behavior. | Accepted principle | 2026-08-06 | Confirmed during the C-004 interview. |
| D-015 | Frontend source uses only `app`, `features`, and `shared` at the architectural top level. Global `pages`, `components`, and `hooks` buckets are removed; feature-internal folders exist only when justified. | Accepted architecture | 2026-08-06 | Confirmed during the C-004 interview. |
| D-016 | Use explicit idiomatic Go constructor composition with consumer-owned interfaces, concrete constructor results, and explicit lifecycle cleanup. Wiring may be handwritten or generated according to graph complexity. Do not use globals, service locators, or hidden runtime reflection. | Accepted architecture | 2026-08-06 | Confirmed during the C-005 interview after clarifying that “manual DI” was the wrong term. |
| D-017 | Code generation is allowed when it is deterministic, has a clear handwritten source of truth, is reproducible by project commands, keeps generated output isolated, and is checked for staleness. Generated files are never hand-edited. | Accepted principle | 2026-08-06 | User explicitly clarified that codegen should not be avoided. |
| D-018 | CI uses a zero-warning blocking policy for handwritten project code: full formatting, import organization, lint, strict type checks, tests, and generated-artifact freshness must pass. Enable `exactOptionalPropertyTypes` for dashboard code while retaining `skipLibCheck` solely for external declaration files. Exceptions must be narrow, documented, and limited to generated or third-party-derived code. | Accepted quality policy | 2026-08-06 | [ADR-0005](../decisions/0005-quality-and-incremental-delivery.md); replaces the current lint-only frontend gate. |
| D-019 | Use risk-based testing rather than a global coverage-percentage target. Add characterization coverage before high-risk structural changes, a regression test for every corrected defect, real-browser coverage for critical admin and runtime-base-path flows, and Go unit/integration coverage for business and data behavior including both supported databases where dialect behavior matters. Every incremental step must keep all gates green. | Accepted quality policy | 2026-08-06 | Confirmed during C-009. |
| D-020 | Prioritize perceived admin-UI stability, critical-flow responsiveness, and dashboard-query efficiency over bundle or binary minimization. Treat frontend bundle and Go binary size as measured guardrails, establish a reproducible fresh baseline before setting numeric budgets, and require evidence for any material regression. | Accepted performance policy | 2026-08-06 | Confirmed during C-010; exact budgets are a Phase 1 measurement output rather than an arbitrary discovery-time number. |
| D-021 | “Latest” means the latest stable release at the time of the planned upgrade; prerelease, RC, beta, and nightly versions are excluded. Prefer Rust/SWC-based tooling when it is at least as correct, maintainable, and architecture-appropriate, but never trade correctness or maintainability merely for its implementation language. | Accepted dependency policy | 2026-08-06 | Confirmed during the first part of C-011. |
| D-022 | Prefer the standard library and already accepted stack before adding a package. Every direct dependency must justify its capability, maintenance, security, and bundle/binary cost; remove unused, redundant, or wrongly production-scoped dependencies. New direct dependencies and replacements require an approved phase plan. Transitive and shadcn-managed packages need no separate approval when they belong to an already approved dependency set. | Accepted dependency policy | 2026-08-06 | Confirmed during the second part of C-011; applies to Go and frontend dependencies. |
| D-023 | Deliver the normalization as one coordinated update through individually reviewable, green phases. Temporary migration bridges may exist only inside their owning phase and must be removed before final completion. Stop only when the full definition of done is satisfied. Implementation requires separate approval of the completed phase plan, and any tracker breaking change still requires its own approved migration and rollback plan. | Accepted delivery policy | 2026-08-06 | Confirmed during C-012. |
| D-024 | Keep architectural rationale, migration consequences, and focused examples in maintained documentation or ADRs. Add code comments where they explain a non-obvious reason, invariant, or safety constraint at its enforcement boundary; do not add comments that merely restate self-explanatory code. | Accepted documentation policy | 2026-08-06 | Explicitly added by the user while confirming C-012. |
| D-025 | Use a stable loading contract: skeletons appear only for first load; background refresh retains prior data with subtle progress; protected UI is not rendered before auth resolution; loading and refresh do not blank the screen or shift established layout geometry; mutations must not expose contradictory stale state. | Accepted UX invariant | 2026-08-06 | Confirmed during the first part of C-013. |
| D-026 | Use site-centric admin navigation: after login open the last still-accessible site, or the only site directly; when no site exists, lead to creation. Keep a site switcher readily available. Give analytics and settings separate routes rather than a `view` query toggle, while filters, date ranges, and pagination remain strictly typed URL state that survives refresh and browser navigation. | Accepted UX architecture | 2026-08-06 | Confirmed during the second part of C-013. |
| D-027 | A site supports multiple domains as a first-class product capability. Architecture normalization, API changes, persistence mapping, validation, tracker behavior, and UI redesign must preserve that capability. | Accepted product invariant | 2026-08-06 | Explicitly added by the user while confirming site-centric navigation; aggregate semantics are recorded in D-028. |
| D-028 | Treat every configured domain of a site as an alias of one analytics property: all use the same site key and contribute to the same site-level dataset. Do not add a per-domain analytics dimension or filter as part of this update. | Accepted product invariant | 2026-08-06 | Confirmed during the final part of C-013; matches current site-scoped storage semantics. |
| D-029 | Isolate authentication per runtime base path. Derive deterministic cookie names from the normalized `BASE_PATH`, scope cookie `Path` to that base path, and preserve `HttpOnly`, `Secure`, and `SameSite` protections so login or logout in one same-origin instance does not affect another. | Accepted security and deployment invariant | 2026-08-06 | [ADR-0001](../decisions/0001-runtime-base-path.md); the current fixed names with `Path=/` do not provide instance isolation. |
| D-030 | Approve the detailed seven-phase implementation plan and its dependency keep/add/remove/replace matrix. Work may proceed incrementally under that scope; a new dependency, tracker break, destructive migration, or material scope expansion still requires its separately defined approval. | Accepted delivery authorization | 2026-08-06 | [ADR-0005](../decisions/0005-quality-and-incremental-delivery.md) and explicit user approval of the plan. |
| D-031 | Keep every admin composition mobile- and wide-screen friendly. Pages must retain usable mobile gutters and avoid horizontal page overflow at 320 px; bounded forms such as Add New Site remain centered on wide screens. Verify critical routes and interactive states at 320, 768, 1024, and 1440 px. | Accepted UX invariant | 2026-08-06 | Explicit user requirement; protected by `responsive-layout.spec.ts`. |
| D-032 | Use the Tweakcn Zen Inspired theme as the admin visual baseline through centralized shadcn semantic tokens. Preserve its light/dark palette, radius, tracking, and shadow character; use built-in system font stacks without external font requests or font packages. | Accepted UX invariant | 2026-08-06 | User selected the theme and explicitly approved built-in fonts. Protected by `theme-contract.spec.ts`. |

## Definition of done

Status: Accepted; numeric measurement budgets were established from two Phase 1 baseline runs.

- [x] Confirmed intent includes outcome, user, why now, success, binding constraint, and explicit non-goals.
- [x] Every accepted expensive decision has an ADR.
- [x] Runtime base-path matrix passes using one dashboard build.
- [x] Frontend and backend boundaries match the accepted architecture.
- [x] No confirmed dead code or redundant direct dependency remains.
- [x] Dependency direction is checked by tooling where practical.
- [x] Frontend full Biome check, typecheck, tests, and production build pass.
- [x] Go formatting, lint, tests, race-relevant tests, generation, and builds pass.
- [x] SQLite and PostgreSQL migration/integration coverage passes.
- [x] GraphQL, tracker, environment-variable, URL, and database compatibility match the accepted policy.
- [x] Bundle, binary, query, and runtime targets are measured where targets were defined.
- [x] Documentation describes the final stack and module ownership without stale alternatives.
- [x] Accepted architectural reasons, migrations, and operational behavior have current documentation and focused examples; non-obvious code-level invariants are explained at their enforcement boundary.
- [x] Final worktree contains no unintended changes and completion evidence is recorded below.

## Completion evidence

Add dated entries only after verification.

| Date | Phase or item | Evidence | Result |
| --- | --- | --- | --- |
| 2026-08-06 | Phase 0 — plan and decision freeze | ADR-0001 through ADR-0005, refreshed npm/Go registry snapshot, approved [`tasks/plan.md`](../../tasks/plan.md), and D-030 | Passed; implementation authorized and no blocking clarification remains. |
| 2026-08-06 | Task 1.1 — non-mutating zero-warning quality command | `task check` run twice in the Dev Container; full Biome check, dashboard/tracker typechecks, source-size policy, Go tests twice, `gofmt`-enabled golangci-lint, isolated generation freshness, temporary dashboard build, and `go build ./...`; `git diff --check` and generated-path diff clean afterward | Passed; both runs exited 0, the second produced no workspace artifacts, and the Vite runtime-config warning was removed at the HTML integration boundary. |
| 2026-08-06 | Task 1.2 — TypeScript strictness baseline | Enabled `exactOptionalPropertyTypes`; separated browser (`tsconfig.json`, no ambient Node types) from tooling (`tsconfig.node.json`, explicit Node types); dashboard and tracker typechecks, full frontend check, and warning-free temporary production build | Passed; optional values are now omitted or normalized deliberately, and the unused custom date-picker ref facade was replaced by the native button ref contract. |
| 2026-08-06 | Tasks 1.3 and Checkpoint 1A — frontend test harnesses | Bun `node:test` smoke coverage for refresh-state retention; Playwright `1.62.1` Chromium smoke against a temporary production dashboard build, real Go server, SQLite database, migrations, seeded admin, login, and authenticated Sites UI; CI installs Go and Chromium dependencies and runs both harnesses | Passed; 2 unit tests and 1 browser flow pass, strict test types/full Biome are green, the accessibility heading exposed by the first browser RED-run was fixed, and port `4173` is released after cleanup. |
| 2026-08-06 | Tasks 1.4, 2.1, and 2.2 — runtime base path and auth isolation | Playwright matrix at `/`, `/lovely-eye`, and `/tools/lovely-eye` covers runtime config, relative assets, unprefixed-asset 404, GraphQL login, cookie scope, nested route direct load/refresh, and logout; a test-only Go proxy mounts `/instance-a` and `/instance-b` on one origin from one dashboard build; focused dashboard HTTP and auth cookie Go tests | Passed after RED runs exposed and fixed Vite root-absolute assets, missing root `<base href="/">`, and root-scoped/shared auth cookies; logging out of one same-origin instance leaves the other authenticated. |
| 2026-08-06 | Task 1.5 — critical admin state transitions | Controlled Playwright GraphQL gates (no sleeps) cover unresolved auth without protected/login flicker, first-load skeletons, disabled mutation feedback with retained form values, real multi-domain creation, fresh-page skeletons, retained data and stable heading geometry during background refresh, site switching, typed preset URL state, browser back/forward, and refresh | Passed; 2 focused state scenarios pass against the real Go+SQLite application, alongside the existing admin/base-path/isolation browser flows. |
| 2026-08-06 | Task 1.6 — backend contract characterization | Focused tests cover construction failures (tracker, DB driver, trusted proxies), runtime base-path config normalization, dashboard SPA/static fallback, scoped auth cookies, current GraphQL semantic errors (`unauthorized`, invalid domain, site not found), and two configured domains collecting into one site-level visitor/page-view dataset; existing e2e/repository/config suites cover the remaining server contracts | Passed; new focused server/config/e2e tests and golangci-lint are green. Stable typed GraphQL extension codes remain intentionally scheduled for Task 6.1. |
| 2026-08-06 | Task 1.7 and Checkpoint 1B — reproducible baseline | `./scripts/measure-baseline.sh` run twice serially in the same Go 1.26.4/Bun 1.3.14 Linux arm64 Dev Container; fresh manifest/gzip report, trimmed Go binaries, ten-sample repository benchmarks with allocation counts, real Chromium readiness, and GraphQL operation fan-out; numeric Phase 6 guardrails recorded above | Passed; deterministic size results were identical, timing ranges were comparable, full `task check` passed, and the complete Bun/Playwright unit, admin, three-base-path, and same-origin isolation matrix passed. |
| 2026-08-06 | Task 2.3 — site and persistence error semantics | RED/GREEN service tests prove domain lookup failures are propagated at the lookup boundary, missing rows map to `ErrSiteNotFound`, infrastructure failures remain distinct, and create/update/delete/key-regeneration expose stable feature errors; repository existence queries use `EXISTS`, update/delete enforce affected-row semantics, and focused GraphQL/site/multi-domain e2e remains green | Passed; focused tests and golangci-lint passed, followed by full `task check` with Go tests run twice. The query/update behavior is dialect-neutral, so no dialect-specific PostgreSQL branch was introduced. |
| 2026-08-06 | Task 2.4 — auth persistence and HTTP boundary | RED/GREEN tests cover persisted `CreatedAt` across register/login/me, storage failures distinct from invalid-credential/not-found feature errors, and deterministic base-path cookie behavior through a dedicated `CookieManager`; `auth.Service` no longer imports HTTP or exposes cookie methods, middleware no longer asserts a concrete service type, and generated e2e GraphQL types were refreshed | Passed; auth unit and focused GraphQL tests, full `task check`, four critical Playwright flows, and same-origin multi-instance auth isolation are green. |
| 2026-08-06 | Task 2.5 and Checkpoint 2 — GeoIP failure semantics | RED/GREEN coverage fixes empty-site `EXISTS` semantics, preserves lookup-close/download/refresh/country-sync causes in `GeoIPStatus`, classifies startup and post-site-update GeoIP work as observable best-effort, returns actionable refresh status without null GraphQL data, refetches status after admin actions, and isolates Playwright from external GeoIP downloads | Passed; focused service/GraphQL tests, full `task check`, five critical browser tests including visible enable/retry failure state, the three-path runtime matrix, and same-origin auth isolation are green. All Phase 2 known findings are closed with regression evidence. |
| 2026-08-06 | Task 3.1 — frontend compiler and core runtime | Updated React/DOM 19.2.8, Apollo 4.2.10 with its direct RxJS peer, GraphQL 17.0.2, TanStack Router 1.170.20/plugin 1.168.25, Vite 8.2.0, SWC plugin 4.3.3, TypeScript 7.0.2, Biome 2.5.7, and GraphQL Codegen CLI 7.2.0/client preset 6.1.1; replaced the unavailable TypeScript 7 compiler-API size checker with Biome's native file/function limits; isolated the codegen partial-fragment overload compatibility at one typed `readFragment` boundary | Passed; generation-compatible strict TypeScript 7 checks, Biome, Bun unit tests, warning-free Vite 8 production build, real admin login smoke, and runtime base-path browser coverage are green. No TypeScript 6 compatibility bridge or deprecated Vite configuration was added. |
| 2026-08-06 | Task 3.2 — redundant frontend tooling removal | Reference search found only the active SWC React plugin, Tailwind Vite plugin, and Lightning CSS transformer; removed direct Babel React plugin, Tailwind PostCSS adapter, PostCSS, Autoprefixer, router generator, and `react-is`; upgraded the retained Tailwind/Vite and Lightning CSS packages; verified `bun install --frozen-lockfile` | Passed; full Biome/TypeScript/unit checks and a fresh manifest production build are green. Compared with the Phase 1 baseline, initial gzip is 241,808 B (+599 B, +0.25%), site-view route gzip is 123,285 B (-49 B), total JavaScript gzip is 362,960 B (+699 B), and CSS gzip is 8,974 B (+8 B), all inside the accepted guardrails. |
| 2026-08-06 | Task 3.3 — Go dependencies and executable tools | Rechecked the module graph, upgraded all outdated retained direct modules, moved the module baseline to Go 1.26.5, replaced `tools.go` with pinned `tool` directives for gqlgen, genqlient, and govulncheck, and changed generation to `go tool`; generation was run twice and freshness verified | Passed; `go mod tidy`, `go mod verify`, full tests/lint/build, generated-artifact comparison, and `govulncheck ./...` are green. The first scan correctly caught GO-2026-5856 in Go 1.26.4 and the 1.26.5 upgrade removed it; verbose output reports no reachable/package vulnerabilities, with only the non-imported and unfixable `x/crypto/openpgp` module advisory remaining. |
| 2026-08-06 | Task 3.4 — Go runtime package removal | Replaced the migration CLI framework with explicit stdlib dispatch, help/error handling, lock lifecycle, output errors, and SQLite lifecycle tests; removed unconditional Bun query logging; propagated caller context through database connect/ping, migrations, GeoIP initialization, and initial-admin creation; fixed migration scripts to clean Compose state on every exit | Passed; focused and full Go tests/lint/builds are green, SQLite and PostgreSQL suites completed init/up/down/reapply successfully, and neither `urfave/cli` nor `bundebug` appears in the compiled server/migrate dependency graph. `urfave/cli` remains only as a transitive dependency of pinned build-time gqlgen. Server size is 26,517,128 B (+72,987 B, +0.28%); migrate is 19,068,251 B (-388,032 B, -1.99%). |
| 2026-08-06 | Task 3.5 and Checkpoint 3 — dead code and mechanical normalization | Removed 19 repository/handler methods with no production or test references, three unused frontend barrels, and commented TODO-only menu UI; applied full Go/Biome formatting; verified the frozen Bun install, tidied/verified Go graph, generated artifact freshness, complete tests/lint/build, and two fresh measurements | Passed; deterministic results were identical across both measurements: initial gzip 241,800 B (+591 B, +0.24%), site-view gzip 123,286 B (-48 B), total JavaScript gzip 362,939 B (+678 B), CSS gzip 8,974 B (+8 B), server 26,437,835 B (-6,306 B), and migrate 19,068,251 B (-388,032 B). Browser readiness remained 864–904 ms and GraphQL fan-out 8–9. Database allocations remained unchanged; one run matched the baseline timing while the immediate repeat was host-load noisy and is retained as an investigation sample for Task 6.4 rather than treated as a deterministic regression. |
| 2026-08-06 | Tasks 4.1–4.2 — frontend boundaries and shadcn foundation | Moved runtime composition and thin routes to `app`, feature behavior to five `features`, reusable runtime/API/lib/UI to `shared`, and generated GraphQL output to `shared/api/generated`; added an executable import/legacy-root boundary gate; rebuilt the new-york/Tailwind v4/OKLCH foundation on unified `radix-ui`, current shadcn sources, `data-slot`, and only adopted components; official `shadcn diff` reported no updates | Passed; boundary check, strict TypeScript, 16 unit tests, generated freshness, official component diff, warning-free production build, and three-path runtime matrix are green. No compatibility barrels or primitive imports outside `shared/ui` remain. |
| 2026-08-06 | Tasks 4.3–4.5 — auth, sites, and analytics vertical slices | Auth owns its operations/state/redirect decision; Sites owns list/create/switch and normalized multi-domain input; `/` selects the last accessible/only site or the explicit list/create destination; analytics has an explicit route, feature view models, final typed URL search values, deterministic history/refresh behavior, and a lazy chart plot with stable skeleton geometry | Passed; auth/site/analytics unit tests, no-flash first-load and retained-refresh browser assertions, multi-domain create, last-site selection, back/forward/refresh, full 7-test admin suite, three runtime base paths, and same-origin auth isolation are green. The auth query no longer loads sites. |
| 2026-08-06 | Tasks 4.6–4.9 and Checkpoint 4 — settings, event definitions, no-flicker lifecycle, and legacy removal | Settings owns domains/tracking/blocking/GeoIP/danger-zone behavior and composes event definitions only through its public feature entry; event definitions own operations, validation, accessible shadcn field controls, and explicit mutation failures; loading/error/stale/mutation states follow D-025; legacy `components`, `config`, `gql`, `hooks`, `layouts`, `lib`, `pages`, and `routes` contain no source files and are forbidden by the boundary gate | Passed; focused browser contracts cover settings first load, persisted multi-domain update, site deletion, and event-definition create/edit/delete with in-flight state. RED runs exposed and fixed two real FK defects: definition deletion now removes fields transactionally, and site deletion atomically removes its complete analytics/configuration graph. Repository regression tests, full `task check`, `task generate-check`, 7 Playwright flows, base-path matrix, and isolation all pass. |
| 2026-08-06 | Phase 4 performance checkpoint | Fresh `./scripts/measure-baseline.sh` on Go 1.26.5/Bun 1.3.14: initial 189,227 B gzip, analytics increment 79,010 B, all JavaScript 393,935 B, CSS 11,412 B, server 26,439,172 B, migrate 19,068,251 B; Chromium login 1,381 ms, dashboard 837 ms, 8 GraphQL operations; repository benchmarks remained within time/allocation targets | Passed for Phase 4 runtime and total-transfer acceptance. Initial payload fell 21.7% from the Phase 1 baseline and analytics fell 35.9%; the old login/register/sites attribution and 10,000 B CSS guardrail are intentionally left visible for Task 6.4 because modular shadcn chunk ownership changed and CSS is 1,412 B over target. |
| 2026-08-06 | Task 5.1 — application composition and lifecycle | Added an explicit `internal/app` composition root for configuration, migrations, repositories/features, startup tasks, HTTP transport, serve lifecycle, bounded graceful shutdown, and cleanup; `cmd/server` now only loads config, establishes signal context, runs, and reports errors; configuration parsing returns validated aggregated errors instead of terminating or silently accepting malformed values | Passed; construction failure, cancellation, occupied-port, command-package, config table, e2e, full `task check`, lint, generation freshness, and all builds are green. No runtime container, global service locator, or hidden lifecycle remains. |
| 2026-08-06 | Task 5.2 — auth feature ownership | Auth now owns a narrow `UserStore`, stored authentication shape, feature errors, concrete service, and claims context; its service imports no SQL, Bun, repository, persistence model, GraphQL, or HTTP package. The repository adapter maps Bun rows/errors to the auth contract, while cookie policy and refresh middleware live in `transport/http`; GraphQL owns only the consumer interfaces it calls | Passed; focused auth/repository/cookie/transport/app tests, backend e2e, and golangci-lint are green. Tests preserve first-user admin policy, disabled registration, storage-error classification, persisted timestamps, and deterministic runtime-base-path cookie isolation. |
| 2026-08-06 | Task 5.3 — site feature and persistence mapping | Moved site behavior to `internal/site`, where the feature owns site/domain/blocking types, validation, authorization, multi-domain invariants, stable errors, and a narrow `Store`; the Bun repository now maps database rows and missing-row errors at the adapter boundary. GraphQL, collection HTTP handling, analytics access checks, app composition, and seed callers consume the feature type rather than Bun models | Passed; site/repository/GraphQL/handler/services/app focused tests, full backend tests including GraphQL aggregate multi-domain contracts, golangci-lint, and SQLite/PostgreSQL migration up/down/reapply suites are green. Domain order, alias aggregation, authorization, blocked lists, key regeneration, and complete site deletion remain protected. |
| 2026-08-06 | Task 5.4 — analytics feature ownership | Moved collection, identity rotation, session policy, filters, aggregates, query behavior, bot detection, and their tests from the generic service bucket into `internal/analytics`; added feature-owned query/filter/stat/event contracts and explicit repository-to-feature mapping so Bun result types no longer reach GraphQL or HTTP; renamed the concrete constructor/type idiomatically to `analytics.NewService`/`analytics.Service` | Passed; analytics unit/integration tests, handlers, GraphQL, app, backend E2E, full Go tests, and golangci-lint are green. Tracker request fields and collection behavior are unchanged, including UTC-day-skipped identity, duplicate/exit semantics, event sanitization, and site-level multi-domain aggregation. |
| 2026-08-06 | Task 5.5 — event, country, and GeoIP feature ownership | Moved event-definition validation/types/store/service to `internal/event`, country behavior/types/store to `internal/country`, and GeoIP application lifecycle to `internal/geoip/service`; retained core GeoIP contracts in `internal/geoip` with downloader/lookup as sibling adapters to preserve an acyclic Go import graph. Analytics now consumes event and GeoIP contracts instead of their repositories/models; the generic `internal/services` package has no files | Passed; focused feature/repository/handler/GraphQL tests, GeoIP failure/refresh/close tests, backend E2E, complete Go tests, and zero-warning lint are green. Event field validation, persisted country normalization, best-effort analytics behavior, and actionable GeoIP status causes remain protected. |
| 2026-08-06 | Task 5.6 — feature-owned GraphQL transport | Split the monolithic SDL and resolver file into `auth`, `site`, `analytics`, `event`, `geoip`, and shared schema sources with matching handwritten feature resolver adapters; configured gqlgen to preserve those adapters while generated executor/model output remains isolated; all Go and TypeScript generators consume the schema glob. The freshness gate now rejects generated resolver markers and newly generated stubs until they are adopted as linted handwritten adapters | Passed; codegen completed twice, generated freshness is green, GraphQL E2E and all Go tests passed twice, Biome/strict TypeScript/unit tests/builds passed, and golangci-lint inspected the handwritten resolvers with zero issues. The schema contract and runtime base-path boundary are unchanged. |
| 2026-08-06 | Task 5.7 — platform, internal utility, and seed ownership | Moved environment parsing and database lifecycle under `internal/platform`, moved trusted-proxy client-IP resolution under `internal`, absorbed site validation into the site feature, and removed the application-only public `pkg` tree. Extracted reusable example-data loading into `internal/seed`; the command now only loads configuration, invokes the operation, and reports its result. Added a blocking Go boundary check for legacy roots/imports and business-layer imports from `cmd` | Passed; seed integration coverage proves first-run creation and repeat-run reuse on migrated SQLite, all Go tests and zero-warning lint are green, all command binaries build, the boundary gate passes, and SQLite/PostgreSQL init/up/down/reapply migration suites pass. No runtime, tracker, data, or base-path contract changed. |
| 2026-08-06 | Task 5.8 and Checkpoint 5 — backend legacy-layer removal | Replaced generic `models`, `repository`, `services`, `handlers`, and `middleware` buckets with feature-owned persistence adapters and explicit HTTP transport packages; database row types remain private to persistence, while GraphQL and HTTP consume feature contracts. Strengthened the executable boundary gate to reject legacy roots/imports, persistence leakage into GraphQL/HTTP, and inward imports from persistence adapters | Passed; full `task check` exited 0, including Go tests twice, E2E, boundary checks, zero-warning lint, generation freshness, frontend checks, and all builds. SQLite and PostgreSQL migration suites each completed init/up/down/reapply successfully; `git diff --check` is clean. |
| 2026-08-06 | Tasks 6.1–6.4 and Checkpoint 6 — API, data, and performance | Added six stable GraphQL error codes with masked/logged internal failures and typed frontend handling; removed unused SDL; replaced untyped GeoIP and optional pagination contracts with generated enums/required paging; made event-count totals exact; removed permanent GeoIP polling, hidden-tab polling, and redundant Site refetches; consolidated dashboard reads with overview and SQL-window queries; added executable bundle/binary/request/fan-out guardrails | Passed after RED/GREEN coverage for error categories and out-of-range exact totals. A representative dashboard operation fell from 25 to 11 SQL queries; Chromium uses 8 GraphQL requests. Fresh guarded results: initial 189,318 B gzip, analytics 79,026 B, all JavaScript 394,043 B, CSS 11,412 B, server 26,609,018 B, migrate 19,074,606 B, login 1,374 ms, dashboard 901 ms. SQLite and isolated PostgreSQL query tests, both migration init/up/down/reapply suites, 7 Playwright admin flows, and full `task check` pass. |
| 2026-08-06 | Task 7.1 — security and privacy hardening | Threat-boundary review; CSP/HSTS/nosniff/frame/referrer/permissions headers; bounded HTTP server and GraphQL body/complexity; production playground removal; exact-one-value and persistence-length collect validation; HTTP(S)-only stored referrers; raw-IP log removal; bounded/redacted GeoIP downloads; dependency overrides and residual-risk register | Passed; RED/GREEN security, config, GraphQL, collect, referrer, and GeoIP tests are green. `govulncheck` reports zero reachable/package vulnerabilities; `bun audit --prod` reports zero, while the full dev-tool audit matches the documented 2 high, 2 moderate, and 1 low deferred records. |
| 2026-08-06 | Task 7.2 — complete release matrix | Frozen Bun install; `go mod tidy -diff`; final `task check`; race tests for auth, analytics, persistence, and HTTP transport; 7 admin flows; one-build `/`, `/lovely-eye`, `/tools/lovely-eye` matrix; same-origin auth isolation; isolated SQLite/PostgreSQL init/up/down/reapply; real example-data command | Passed; 18 Bun units, all Go tests twice, zero-warning lint, generation freshness, both builds, all browser suites, and migrations are green. Seed command created 80 clients, 169 sessions, 710 page views, and 272 predefined events in an ephemeral DB; all ephemeral DBs were removed. |
| 2026-08-06 | Tasks 7.3–7.4 — documentation and implementation closure | Updated root stack/config guide, feature-owned frontend/server module maps, dashboard AI instructions, analytics/security operations, dependency rationale, and phase evidence; stale legacy-path reference scan, tracked-secret scan, worktree/status review, and `git diff --check` | Passed; all definition-of-done items have evidence, no generated/stale/temp artifact remains, and the approved normalization implementation is complete. Final maintainer review remains the explicit handoff action. |
| 2026-08-06 | Responsive admin follow-up | Centralized bounded Add New Site and Site Settings content; established shared 16/24/32 px mobile-to-desktop gutters and a wide-screen shell bound; made responsive grids explicitly shrinkable; hardened long domains, analytics labels, event metadata, tracking controls, GeoIP paths, editor actions, and code snippets | Passed; real Chromium coverage verifies critical routes at 320/768/1024/1440 px plus custom date picker, key regeneration, and event-editor states without horizontal page overflow. Dashboard static checks, 18 unit tests, and all 9 primary browser flows pass. |
| 2026-08-06 | Site-creation recovery regression | Clear stale form errors whenever the user edits the site name or domain collection after a rejected submission | Passed; RED/GREEN Chromium coverage proves an invalid domain can be corrected and submitted successfully. Dashboard checks and all 10 primary browser flows pass. |
| 2026-08-06 | Zen Inspired admin theme | Replaced the global shadcn semantic palette, charts, radius, tracking, and runtime-switchable shadows with the selected Tweakcn theme; used system sans/serif/mono stacks; normalized destructive foreground usage; documented the theme boundary for future AI changes | Passed; visual QA of Add New Site and Analytics in both modes found no console errors or external requests. The light/dark contract, 18 unit tests, all 11 primary browser flows, the `/`, `/lovely-eye`, and `/tools/lovely-eye` deployment matrix, zero-warning lint/type checks, generation freshness, Go tests, and production builds are green. |
| 2026-08-15 | Full stable dependency and CI refresh | Rechecked npm, Go, Docker Hub, MCR/GHCR, and official GitHub Action releases; upgraded every outdated direct frontend dependency, every retained `go.mod` requirement and tool, Go 1.26.6, PostgreSQL 18.6, production base images, checkout/setup-go v7, generated GraphQL output, and security overrides; updated current deployment examples and the residual-risk register | Passed; `bun outdated` and the per-requirement Go proxy audit are empty, frozen Bun install and both Go module verifications are clean, `shadcn diff` and generation freshness report no pending component/artifact changes, actionlint exits zero, production Bun audit and reachable/package Go vulnerability results are empty, and the documented dev-only picomatch/module-only openpgp residuals are unchanged in runtime reachability. Full `task check`, targeted race tests, 11 admin flows, the three-path runtime matrix, same-origin isolation, SQLite/PostgreSQL init/up/down/reapply, and the exact-image production health/dashboard check pass. The fresh performance baseline also passed: initial JavaScript 189,521 B gzip, analytics increment 78,765 B, all JavaScript 394,042 B, CSS 11,659 B, server 27,000,882 B, migrate 19,152,362 B, Chromium login 1,366 ms, dashboard 903 ms, and 8 GraphQL operations. |
| 2026-08-15 | Major-release and container-upgrade handoff | Audited the Release Please and GHCR flow against SemVer, GitHub Release, Docker tag/digest, Docker Metadata Action, and OCI annotation guidance; added v2 consumer release notes, a permanent backup/upgrade/rollback guide, a maintainer release runbook, controlled SemVer container channels, OCI image metadata, automatic publication of the multi-architecture image digest, and a bounded non-secret Docker build context | Passed; Release Please configuration validates against the official schema, actionlint exits zero, Biome checks all 199 files, `task check` passes, and the production image builds for the supported target with a 1.91 MB context and the expected OCI labels. An isolated v1.7.0 → v2 schema → v1.7.0 SQLite exercise proves the documented code-only rollback contract. The GitHub Release stays draft until CI passes and the exact multi-architecture image plus digest are available. No runtime, database, tracker, API, or base-path implementation changed. |
| 2026-08-15 | Go collection RAM and allocation optimization | Profiled the production SQLite container, Go GC, collect handler, site persistence, and rate limiter under idle and repeatable load; removed the duplicate site graph lookup across the HTTP/analytics boundary; retained live configuration reads while replacing per-request Bun relation builders with narrow raw queries; added a collect benchmark, allocation/query guards, relation coverage, a profiling runbook, and Go test-binary Docker exclusion | Passed; ten-sample warm collect results moved from 198–206 us/op, 95.8–96.0 KiB/op, and 900 allocs/op to 120–124 us/op, 45.2–45.3 KiB/op, and 565 allocs/op. On identical cloned SQLite data and 2,000 requests, GC cycles fell from 91 to 39 and VmHWM from 23.5 to 21.5 MiB while elapsed time remained 3,044 ms and cgroup memory remained about 28 MiB. A real PostgreSQL 18.6 image accepted the optimized collect path. The value-map limiter experiment was reverted after worsening rotating-key memory from 138 to 220 B/op; hard limiter capacity remains a separately documented traffic-policy clarification rather than an assumed behavior change. |
| 2026-08-15 | Deterministic frontend GraphQL generation | Traced the CI-only stale `fragment-masking.ts` result to coexisting client preset `6.1.3` and CLI-nested `6.1.1` entries introduced by the dependency refresh; globally overrode the accepted preset version and regenerated the frozen Bun lockfile | Passed; a clean frozen install contains only preset `6.1.3`, the native `codegen --check` gate passes without a workflow workaround, `task check` exits zero, and `git diff --check` is clean. |
| 2026-08-17 | Transitive dependency override cleanup | Audited every direct dashboard package and each top-level override; removed all nine overrides after regenerating a coherent lockfile | Passed; parent ranges retain patched Babel and other transitive versions, `js-yaml` returns from the forced incompatible v5 override to supported `4.3.1`, a fresh install resolves one client preset `6.1.3`, development audit findings are unchanged, production audit is clean, shadcn reports no updates, generation is current, and full `task check` exits zero. |
| 2026-08-17 | Release Please action normalization | Replaced the stale `gvillo` fork with the official `googleapis/release-please-action@v5.0.0`, granted the documented issue-label permission, and serialized release-branch updates per target ref | Passed; official v5 metadata accepts the existing manifest inputs, actionlint `1.7.12` exits zero, and `git diff --check` is clean. Repository rules must still allow the release token to force-update the generated `release-please--branches--*` branch. |

## Interview history

Record concise confirmed answers, not full conversation transcripts.

| Date | Clarification | Confirmed answer | Resulting decision or action |
| --- | --- | --- | --- |
| 2026-08-05 | C-001: primary outcome | Personal bug-free and flicker-free admin panel; safe future development with quality guarantees; strict typing; current dependencies with proper refactoring; Rust/SWC tool preference; clearer architecture for data accuracy; improved admin UX; current from-scratch shadcn/ui integration; no legacy; idiomatic Go and modular UI. | Added D-006 and expanded frontend, dependency, quality, and acceptance workstreams. Dependency-conflict rules remain a separate delivery decision. |
| 2026-08-05 | C-002: regression versus design improvement | Do not mistake an intentional better-by-design contract for a regression merely because it differs from legacy behavior. Internal design and clear UX defects may improve without separate approval; API changes are allowed for required display data or a demonstrably better practice. | Added D-007, D-008, and the change-classification policy. Compatibility limits moved to C-003. |
| 2026-08-05 | C-003: GraphQL consumers | There are no external GraphQL consumers today; private-token access may be added in the future. | Added D-009. Current GraphQL may evolve with the dashboard; future external API design is deferred until it is needed. |
| 2026-08-06 | C-003: tracker compatibility | Tracker contracts may break, but only after prior explicit approval of the migration plan. | Added D-010 and made compatibility the default until an approved plan exists. |
| 2026-08-06 | C-003: database migration compatibility | Existing data may not be deleted and may not require reimport during migration. | Added D-011 as a data-preservation invariant. |
| 2026-08-06 | C-003: environment compatibility | Environment variables may be renamed or removed without permanent legacy aliases when migration documentation and a precise startup error are provided. | Added D-012. |
| 2026-08-06 | C-003: browser URL compatibility | Internal dashboard URLs may change without legacy redirects; runtime base-path behavior, nested-route refresh, and tracker URLs remain protected. | Added D-013 and completed C-003. |
| 2026-08-06 | C-001: final intent confirmation | Explicitly confirmed the Outcome/User/Why now/Success/Constraint/Out of scope restatement without refinements. | Promoted the intent hypothesis to confirmed intent and opened C-004. |
| 2026-08-06 | C-004: route and feature ownership | TanStack route files are thin routing infrastructure; feature modules own GraphQL, state, screens, and business UI. Use shadcn/ui as much as practical for styling and components. | Accepted D-002 and the default-use portion of D-014. |
| 2026-08-06 | C-004: shadcn ownership | Canonical shadcn components live in `shared/ui`; features compose them directly; primitive imports stay inside that boundary; wrappers exist only for repeated product behavior. | Completed the ownership details in D-014. Exact top-level structure remains open. |
| 2026-08-06 | C-004: frontend top-level structure | Use only `app`, `features`, and `shared` as architectural top-level buckets; remove global `pages`, `components`, and `hooks`; add feature-internal folders only when needed. | Added D-015, completed C-004, and opened C-005. |
| 2026-08-06 | C-005: backend architecture | Use a feature-oriented modular monolith with feature-owned behavior and interfaces, edge-owned transport/platform adapters, minimal `cmd`, and no ceremonial Clean Architecture layers. | Accepted D-003. DI mechanism remains open. |
| 2026-08-06 | C-005: dependency injection and codegen | Initially accepted constructor injection while allowing useful deterministic code generation instead of treating codegen as undesirable. | Added D-016 and D-017; terminology required clarification. |
| 2026-08-06 | C-005: composition clarification | “Manual DI” was misleading. Use explicit idiomatic Go composition; wiring may be handwritten or generated based on complexity, while globals, service locators, and hidden runtime reflection remain excluded. | Corrected D-016, completed C-005, and opened C-006. |
| 2026-08-06 | C-006: persistence boundary | Keep persistence rows and Bun/database tags inside adapters; define feature-owned types only for business invariants or materially different shapes; avoid mechanical one-to-one model duplication. | Accepted D-005, completed C-006, and opened C-007. |
| 2026-08-06 | C-007: GraphQL scope | Do not replace GraphQL with a generated HTTP/JSON API; retain `gqlgen`. | Accepted D-004, completed C-007, and opened C-008. Apollo remains independently reviewable as a frontend dependency. |
| 2026-08-06 | C-008: strictness and CI | Use blocking zero-warning checks for format, imports, lint, strict types, tests, and generated-code freshness. Enable dashboard `exactOptionalPropertyTypes`; retain `skipLibCheck` for external declarations; allow only narrow documented generated/third-party exceptions. | Added D-018, completed C-008, and opened C-009. |
| 2026-08-06 | C-009: test depth and refactor risk | Use risk-based tests without a global coverage quota: characterize risky behavior first, add regression tests for fixes, cover critical admin and base-path flows in a real browser, test data/business logic at unit and integration levels, and exercise both databases where relevant. Keep every incremental step green. | Added D-019, completed C-009, and opened C-010. |
| 2026-08-06 | C-010: performance and size priorities | Prioritize UI stability, critical admin responsiveness, and dashboard-query behavior. Bundle and binary size are guardrails, not primary optimization goals; set numeric budgets only after a reproducible baseline. | Added D-020, completed C-010, and opened C-011. |
| 2026-08-06 | C-011: freshness and Rust/SWC preference | Upgrade to latest stable releases only; exclude prereleases. Treat Rust/SWC as a strong tie-breaker while correctness, maintainability, and architectural fit remain controlling. | Added D-021. Dependency admission and removal criteria remain active under C-011. |
| 2026-08-06 | C-011: dependency admission and removal | Prefer the standard library and accepted stack; justify direct dependencies by capability and lifetime cost; remove unused, duplicate, and incorrectly scoped packages. Require an approved phase plan for new direct dependencies or replacements; covered transitive and shadcn-managed packages need no separate approval. | Added D-022, completed C-011, and opened C-012. |
| 2026-08-06 | C-012: rollout and stopping point | Deliver the entire roadmap as one coordinated update through reviewable green phases; remove temporary migration bridges before completion; stop only at the full definition of done; require separate implementation approval and a dedicated approved plan for tracker breaks. | Added D-023 and completed C-012. |
| 2026-08-06 | C-012: documentation addition | Cover the update with maintained documentation, rationale, and examples; explain relevant non-obvious reasons in code where appropriate. | Added D-024 and opened C-013 for the remaining admin-UX acceptance criteria. |
| 2026-08-06 | C-013: loading and flicker contract | Use first-load-only skeletons, retain data during background refresh, resolve auth before protected rendering, preserve established layout geometry, and prevent contradictory mutation state. | Added D-025. Admin navigation and convenience criteria remain active under C-013. |
| 2026-08-06 | C-013: site-centric navigation | Open the last accessible or only site after login, lead to creation when none exists, keep site switching available, use separate analytics/settings routes, and preserve typed URL-owned analytics state. | Added D-026. |
| 2026-08-06 | C-013: multi-domain site capability | One site must support multiple domains. | Added D-027 and opened verification of the current cross-domain validation, tracker, and analytics semantics before completing C-013. |
| 2026-08-06 | C-013: multi-domain analytics semantics | Configured domains are aliases of one analytics property and aggregate into one site-level dataset; no per-domain dimension is required. | Added D-028 and completed C-013. The subsequent consistency pass opened C-014 for unresolved cookie isolation. |
| 2026-08-06 | C-014: runtime-base-path auth isolation | Use deterministic base-path-specific cookie names and scope each cookie to its normalized runtime base path while preserving existing security attributes. | Added D-029 and completed C-014. No blocking discovery clarification remains. |
| 2026-08-06 | Implementation approval | Approved the detailed phase plan and dependency set as authorization to begin implementation. | Added D-030 and started Phase 0. Tracker breaking changes and out-of-plan dependencies remain separately gated. |
