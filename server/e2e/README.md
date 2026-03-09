# End-to-end testing

No source code as test dependencies. Exceptions:
- root production server
- root config
- basic types/constants for data integrity
- generated operations: `import operations "github.com/lovely-eye/server/e2e/generated"`

Analytics e2e tests should use a fixed `ANALYTICS_IDENTITY_SECRET` so visitor identity stays deterministic across test runs, including the UTC-day-skipped `today`/`yesterday` client reuse path.
