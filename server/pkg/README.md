# Independent packages

App-agnostic code only. Packages that grow enough should be extracted and published as independent libraries.

See each package's `doc.go` for package-level documentation.

Current package boundaries:
- `clientip` - Resolves a client IP from `RemoteAddr`, `X-Forwarded-For`, and `X-Real-IP` only when the remote address is in a trusted proxy CIDR.
- `random` - Random value helpers.
- `textutil` - Text manipulation helpers.
- `urlpath` - URL path normalization helpers.
- `validation` - App-agnostic input validators and validation errors.
