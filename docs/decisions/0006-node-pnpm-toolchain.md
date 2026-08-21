# ADR-0006: Standardize frontend tooling on Node.js and pnpm

## Status

Accepted

## Date

2026-08-21

## Context

The dashboard previously used Bun as package manager, script runtime, test runner, container build
runtime, Dev Container feature, and CI bootstrap. This duplicated an ecosystem-specific runtime
choice across every development and delivery boundary and made the repository's `node:test` suites
depend on a non-Node runner.

The maintainer explicitly requested one current Node.js and pnpm toolchain with no Bun runtime. The
production application remains a Go binary serving static dashboard assets; Node.js is a build and
development runtime and is not added to the final container image.

## Decision

- Use Node.js `26.7.0` and pnpm `11.22.0`, pinned in the dashboard manifest, Dev Container, CI, and
  dashboard build stage. A mismatched local runtime is an error; pnpm must not silently download a
  second project-local Node.js. Update these pins together after official release and compatibility
  review.
- Commit `pnpm-lock.yaml` and require `pnpm install --frozen-lockfile` in CI and container builds.
- Use Node's built-in `node:test` runner for unit tests. Use `tsx` only as the TypeScript loader
  because Node intentionally does not apply `tsconfig` path aliases such as `@/` at runtime.
- Keep dependency lifecycle scripts default-deny and explicitly permit only the reviewed
  `@swc/core` and `esbuild` build steps through pnpm `allowBuilds`.
- Allow GraphQL 17 for the stale `graphql-config` peer declaration without overriding package
  versions. This exception remains valid only while deterministic generation and freshness checks
  pass and is removed when the upstream peer range includes GraphQL 17.
- Use the official `pnpm/setup` action in CI. It reads the pinned package manager and Node runtime
  from `dashboard/package.json`; workflow commands still install explicitly so frozen-lockfile
  behavior remains visible.
- Do not add Bun, npm, Yarn, or a second frontend test runner as a parallel project workflow. npm is
  used only to bootstrap the pinned pnpm binary in the isolated Node.js Docker build stage.

## Alternatives considered

### Keep Bun only as the unit-test runner

Rejected because it would retain the Bun runtime and a second toolchain solely for tests.

### Use Node's TypeScript stripping without a loader

Rejected because Node deliberately ignores `tsconfig.json` path mappings. Rewriting product imports
to a second alias convention would be a larger application-level migration with no product benefit.

### Add Vitest

Rejected because the existing tests already use `node:test`; adding another test framework would
increase the dependency and configuration surface without adding required behavior.

### Use the latest LTS Node.js release

Rejected for this build-only toolchain because the accepted goal is the latest available stable
version where compatible. The exact Current release is pinned and the full build/browser matrix is
required before each update.

## Consequences

- Developers must rebuild the Dev Container once to replace the removed Bun feature with Node.js
  and pnpm.
- Every frontend command runs through the same pnpm dependency graph locally, in CI, and in Docker.
- Toolchain upgrades intentionally update every declared version pin and the lockfile in one
  reviewed change.
- Bun ORM references in the Go server remain unchanged; they name the retained Go persistence
  library and are unrelated to the removed JavaScript runtime.
