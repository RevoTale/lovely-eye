# Lovely Eye Dashboard

A privacy-friendly web analytics dashboard built with React, Vite, and Tailwind CSS.

## Features

- 📊 Real-time analytics dashboard
- 🔒 Privacy-first analytics
- 📱 Responsive design with shadcn/ui components
- 🎨 Dark mode support
- 📈 GraphQL API integration with code generation

## Tech Stack

- **React 18** - UI framework
- **Vite** - Build tool (static export)
- **TypeScript** - Type safety with strict mode
- **Tailwind CSS** - Styling
- **shadcn/ui** - UI component library
- **Apollo Client** - GraphQL client
- **graphql-codegen** - Type-safe GraphQL operations
- **React Router** - Client-side routing

## Development

```bash
# Install dependencies
npm install

# Generate GraphQL types
npm run codegen

# Start development server
npm run dev

# Type check
npm run typecheck

# Lint
npm run lint

# Build for production
npm run build
```

## Static Export

The dashboard is built as a static export and served by the Go backend server. The build output goes to the `dist/` directory.

### Runtime Configuration

The app supports runtime configuration through `config.js`:

- `BASE_PATH` - Base URL path for the dashboard (e.g., `/dashboard`)
- `API_URL` - Backend API URL
- `GRAPHQL_URL` - GraphQL endpoint URL

These can be configured by the Go server when serving the dashboard, allowing the same build to be deployed to different environments.

## Project Structure

```
src/
├── components/     # React components
│   └── ui/         # shadcn/ui components
├── config/         # Runtime configuration
├── generated/      # Generated GraphQL types
├── graphql/        # GraphQL operations
├── hooks/          # Custom React hooks
├── layouts/        # Layout components
├── lib/            # Utility functions and Apollo client
└── pages/          # Page components
```

## Adding shadcn/ui Components

This project uses shadcn/ui. To add new components:

```bash
npx shadcn@latest add [component-name]
```

Or manually copy components from the shadcn/ui documentation.
