# ADR-0004: Retain GraphQL with gqlgen and deterministic generation

## Status

Accepted

## Date

2026-08-06

## Context

Lovely Eye has one bundled dashboard client and no current external GraphQL consumers. The API uses gqlgen on the server and generated typed documents on the frontend. The generated executor is large, and most dashboard operations are fixed, so replacing GraphQL with a smaller HTTP API was considered during normalization.

The maintainer explicitly chose to retain gqlgen. The value is a schema-first, strictly typed contract that can evolve atomically with the bundled dashboard. A possible future private-token API is outside this update and must not silently turn the current internal contract into an externally versioned API.

This ADR records decisions D-004, D-008, D-009, and D-017.

## Decision

- GraphQL remains the dashboard API and gqlgen remains its Go implementation.
- Handwritten SDL and resolver adapters are owned by features; generated gqlgen output stays isolated under transport infrastructure.
- GraphQL transport types do not alias feature storage interfaces, repository results, or persistence rows.
- Frontend operations live with their owning features; generated artifacts live under shared API infrastructure.
- Server and dashboard generation are deterministic, pinned, reproducible by project commands, checked for staleness, and never hand-edited.
- Apollo remains the frontend query lifecycle/cache client under the approved dependency plan. Cache, polling, refetch, pagination, and error policy must be explicit by feature and satisfy the no-flicker contract.
- The current GraphQL schema may change atomically with the bundled dashboard when the change has documented correctness, data, UX, type-safety, or maintainability value.
- Schema, generated Go and TypeScript types, consumers, documentation, and tests change together.
- A future externally consumable private-token API requires a separate contract and versioning ADR.

## Alternatives considered

### Replace GraphQL with generated HTTP/JSON endpoints

Rejected by explicit maintainer decision after reviewing the single-client operation footprint and generated-code cost.

### Preserve every historical GraphQL field indefinitely

Rejected because there are no current external consumers and permanent legacy shapes conflict with the normalization goal.

### Let one global graph package own all business mapping

Rejected because it recreates a horizontal layer and lets transport concerns become the application architecture.

## Consequences

- Generated line count alone is not a reason to remove GraphQL.
- Handwritten resolvers must not be excluded from linting merely because they share a directory with generated files.
- GraphQL error categories and frontend handling require an explicit cross-layer contract.
- Query fan-out and Apollo behavior must be measured and simplified based on evidence.
