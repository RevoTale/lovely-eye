# Project Instructions

## Architecture normalization update

- Read `docs/plans/architecture-normalization.md` before planning or changing architecture, dependencies, routing, GraphQL, persistence, or project structure.
- Treat it as the living source of truth for confirmed intent, constraints, decisions, open clarifications, work order, and completion evidence.
- Record new confirmed decisions and material findings there. Do not promote a proposal to an accepted decision without explicit user confirmation.
- During discovery, use one-question-at-a-time interviews. Keep implementation blocked until the relevant clarification and phase entry criteria are resolved.

## Runtime base path

- Treat deployment under an arbitrary runtime `BASE_PATH` as an architectural invariant.
- The same dashboard build must work at `/`, `/lovely-eye`, and nested paths such as `/tools/lovely-eye`; never require a rebuild for the deployment path.
- Do not hardcode `/dashboard` or another mount path in assets, routes, navigation, API URLs, redirects, or server handlers.
- Build dashboard assets with relative URLs. Apply `BASE_PATH` once at the integration boundary: the Go mount, runtime config, and TanStack Router `basepath`. Keep application route definitions unprefixed.
- Preserve direct loading and refresh of nested client routes through the server's SPA fallback.
- Any change to Vite configuration, routing, static serving, runtime config, authentication redirects, or public URLs must verify the base-path matrix documented in `docs/decisions/0001-runtime-base-path.md`.
