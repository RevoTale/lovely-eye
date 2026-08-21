# Lovely Eye Dashboard - Project Rules

## Stack

- Node.js and pnpm for package management and scripts
- Vite for static export (no SSR)
- TanStack Router for type-safe routing
- Apollo Client for GraphQL
- shadcn/ui + Tailwind CSS for UI

## Build

- Output goes to `dist/`, served by Go server
- Runtime config via `config.js` (not bundled)
- Go server dynamically generates `config.js` per request

## Code Generation

- Run `pnpm run codegen` after modifying GraphQL operations. Generated artifacts live in
  `src/shared/api/generated/` and must not be edited by hand.

## Commands

```bash
pnpm run build    # Production build
pnpm run typecheck
pnpm run codegen  # Generate GraphQL types
```
