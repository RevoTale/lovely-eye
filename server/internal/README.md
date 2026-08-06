# Server module ownership

Lovely Eye is a feature-oriented modular monolith. Feature packages own behavior, types, errors, and
consumer-side store interfaces. Their `persistence` subpackages own Bun rows, SQL, and database error
mapping; persistence types do not cross feature or transport boundaries.

- `analytics`, `auth`, `country`, `event`, and `site` own product behavior.
- `geoip` owns shared GeoIP contracts; `downloader`, `lookup`, and `service` are explicit adapters.
- `graph` owns gqlgen transport mapping, stable public error codes, request-body and operation-complexity limits.
- `transport/http` owns routes, cookies, client-IP trust, collect handling, CORS, rate limits, logging, and security headers.
- `platform/config` and `platform/database` own environment parsing and database lifecycle.
- `dashboard` serves the one-build static SPA under the runtime `BASE_PATH` invariant.
- `seed` owns example-data creation without becoming an application feature dependency.
- `app` is the explicit composition root; commands only load configuration and run an application action.

Do not recreate generic `models`, `repository`, `services`, `handlers`, `middleware`, or `server`
layers. See [ADR-0003](../../docs/decisions/0003-feature-oriented-go-modular-monolith.md).
