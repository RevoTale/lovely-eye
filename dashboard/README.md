# Lovely Eye Dashboard

React dashboard for Lovely Eye analytics.

## Stack

- React + TypeScript
- Vite (static export)
- Tailwind CSS + shadcn/ui
- Apollo Client + graphql-codegen
- TanStack Router

## Development

```bash
bun install
bun run codegen   # generate GraphQL types
bun run dev       # start dev server
bun run build     # production build
```

## Build

Static export to `dist/`, served by Go backend. Go server dynamically generates `config.js` per request:

- `BASE_PATH` - dashboard URL path
- `API_URL` - backend API URL
- `GRAPHQL_URL` - GraphQL endpoint

The same build must work without rebuilding at `/`, `/lovely-eye`, or any nested runtime path. Vite assets stay relative, while the Go server and TanStack Router apply `BASE_PATH` at runtime. Do not hardcode or manually prepend a deployment path in application routes.

See [ADR-0001](../docs/decisions/0001-runtime-base-path.md) for the full invariant and verification matrix.

## Project Structure

```
src/
├── app/             # composition, layouts, and thin TanStack route files
├── features/        # feature-owned API, state, screens, and product UI
├── shared/
│   ├── api/         # Apollo integration and isolated generated GraphQL output
│   ├── config/      # validated runtime configuration
│   ├── lib/         # feature-agnostic utilities
│   └── ui/          # current shadcn/ui sources and shared visual primitives
├── index.css        # Tailwind v4 and shadcn design tokens
└── main.tsx         # browser entry point
```

The boundary checker rejects legacy top-level feature buckets and invalid inward imports. Route and
generated files are infrastructure; never hand-edit `app/route-tree.gen.ts` or
`shared/api/generated/`.

## Adding Components

```bash
bunx shadcn@latest add [component-name]
```

Run `bunx shadcn@latest diff` after component updates and then the complete dashboard checks. Prefer
adopted shadcn components for shared UI; keep product behavior in its owning feature.
