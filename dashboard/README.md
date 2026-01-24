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

Same build works across environments.

## Project Structure

```
src/
├── components/     # React components
│   └── ui/         # shadcn/ui components (auto installed via command. Do not edit manually)
├── config/         # runtime configuration
├── gql/            # generated GraphQL types
├── hooks/          # custom React hooks
├── layouts/        # layout components
├── lib/            # utilities that are app logic agnositc
└── pages/          # page components for TanStack router
```

## Adding Components

```bash
bunx shadcn@latest add [component-name]
```
