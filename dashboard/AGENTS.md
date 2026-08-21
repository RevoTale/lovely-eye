# Dashboard instructions

## Ownership

- Use TypeScript only. Use `.tsx` for React components and `.ts` otherwise.
- Keep composition, layouts, and thin route files in `src/app`.
- Keep product API/state/screens/UI in the owning `src/features/<feature>` module.
- Keep feature-agnostic runtime config, API infrastructure, utilities, and adopted shadcn sources in
  `src/shared`.
- Do not recreate legacy top-level `components`, `config`, `gql`, `hooks`, `layouts`, `lib`,
  `pages`, or `routes` folders.
- Respect the executable Biome and boundary-check limits; split code by responsibility, not merely
  to satisfy line counts.

## UI and data state

- Prefer current shadcn/ui components from `src/shared/ui` and preserve their documented structure.
- Treat the semantic light/dark tokens in `src/index.css` as the visual theme source of truth. Do
  not add feature-local brand colors, external font requests, or font packages.
- Keep the admin responsive and accessible across desktop, tablet, and mobile layouts.
- Preserve shared page gutters and center bounded task forms on wide screens; do not replace the
  shell contract with page-specific outer padding.
- After changing layout or feature composition, verify critical routes and interactive states at
  320, 768, 1024, and 1440 px with no horizontal page overflow.
- Show structure-matching skeletons only on first load. During refresh, retain existing data and use
  subtle progress without layout shifts or flicker.
- Keep mutation feedback visible and never clear valid input after a failed mutation.
- Use generated GraphQL types and stable `extensions.code`; never branch on error message text.

## Runtime base path

- One static build must work at `/`, `/lovely-eye`, and nested runtime paths such as
  `/tools/lovely-eye`.
- Keep Vite asset URLs relative. Configure `BASE_PATH` only at the Go mount/runtime-config boundary
  and TanStack Router `basepath`.
- Keep route definitions, links, and navigation targets unprefixed.
- Read backend and GraphQL URLs from validated runtime configuration.
- Preserve nested-route direct loads and SPA refresh fallback; verify
  `docs/decisions/0001-runtime-base-path.md` after routing or static-serving changes.

## Tooling

- Use the pinned Node.js and pnpm versions from `package.json`; do not introduce Bun, npm, Yarn,
  Vitest, or parallel build/CSS paths.
- Generated `src/app/route-tree.gen.ts` and `src/shared/api/generated` files are never hand-edited.
- Run GraphQL generation after operation changes and the complete dashboard checks before handoff.
