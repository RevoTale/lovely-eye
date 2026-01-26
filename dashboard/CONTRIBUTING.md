# Lovely Eye Dashboard - Project Rules

## Stack

- Bun for package management
- Vite for static export (no SSR)
- TanStack Router for type-safe routing
- Apollo Client for GraphQL
- shadcn/ui + Tailwind CSS for UI

## Build

- Output goes to `dist/`, served by Go server
- Runtime config via `config.js` (not bundled)
- Go server dynamically generates `config.js` per request

## Code Generation

- Run `bun run codegen` after modifying GraphQL operations in `src/gql/`

## Commands

```bash
bun run build    # Type check + production build
bun run codegen  # Generate GraphQL types
```
