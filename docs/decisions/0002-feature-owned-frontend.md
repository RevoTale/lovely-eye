# ADR-0002: Use feature-owned frontend modules

## Status

Accepted

## Date

2026-08-06

Last updated: 2026-08-21

## Context

The dashboard currently divides ownership between global routes, pages, components, hooks, and layout files. A single change can cross several generic folders, while route files delegate to pages that depend back on router state. This obscures feature boundaries, encourages broad barrels and prop chains, and makes loading or navigation behavior difficult to reason about.

The dashboard must also preserve the runtime-base-path contract in [ADR-0001](0001-runtime-base-path.md), use current shadcn patterns as its default UI foundation, and keep strict URL state for analytics navigation.

This ADR records decisions D-002, D-014, D-015, D-025, and D-026 from the architecture normalization plan.

## Decision

Use only these architectural top-level source areas:

```text
dashboard/src/
  app/
  features/
  shared/
```

- `app` owns provider composition, the router integration, thin route files, and application startup.
- A route owns only its path, final typed search parsing, guards/loaders, and feature composition.
- `features` own GraphQL operations, product state, view models, screens, and feature-specific UI.
- `shared` contains only infrastructure or UI with multiple real feature consumers.
- Features do not import another feature's internals or concrete route definitions.
- Generated route and GraphQL files remain isolated and are never edited by hand.
- Canonical current shadcn components live in `shared/ui`; features compose them directly.
- Underlying primitive imports stay inside `shared/ui`. The accepted current setup uses shadcn's Base UI registry and `@base-ui/react`; application and feature modules never import Base UI directly.
- Product wrappers require repeated behavior or policy; renaming or hiding a shadcn API is not sufficient.
- Analytics filters, date ranges, and pagination remain final typed URL state.
- The application uses the stable loading contract from D-025: first-load-only skeletons, retained data during refresh, no protected-content flash, stable geometry, and coherent mutation state.
- Navigation is site-centric: open the last accessible or only site, lead to creation when none exists, keep site switching available, and give analytics and settings explicit routes.

## Alternatives considered

### Keep global `pages`, `components`, and `hooks`

Rejected because technical-role buckets scatter feature changes and preserve the current dependency loops.

### Let routes own complete screens and business state

Rejected because router infrastructure would become the product ownership boundary and make features harder to test or reuse.

### Add a project-specific design-system wrapper over every shadcn component

Rejected because it creates an unnecessary abstraction and prevents direct use of the documented component APIs.

## Consequences

- Frontend normalization proceeds as vertical feature slices, not a bulk folder move.
- Import-direction checks must prevent feature-to-feature internals and shared-to-feature dependencies.
- Existing global folders and transitional barrels must be removed before completion.
- UI migration follows current official shadcn Base UI sources and preserves only adopted product components.
- Browser tests protect routing, base paths, loading transitions, layout stability, and critical admin workflows.
