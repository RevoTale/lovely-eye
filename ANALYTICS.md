# Analytics Implementation

Implementation notes for Lovely Eye analytics.

## Public Collect Contract

`POST /api/collect?site_key=<public_key>` is the only collect shape. `site_key` in the JSON body is not part of the contract.

The tracker sends JSON with `Content-Type: text/plain;charset=UTF-8` so `navigator.sendBeacon` can queue small analytics payloads without a custom request setup.

Page view:

```json
{ "path": "/pricing" }
```

Exit ping:

```json
{ "path": "/pricing", "exit": true }
```

Initial attribution, when present, is sent only on the first page view:

```json
{ "path": "/pricing", "referrer": "https://google.com", "utm_source": "google" }
```

Custom event:

```json
{ "name": "checkout_failed", "path": "/checkout", "properties": "{\"code\":\"PAYMENT_DECLINED\"}" }
```

The client does not send `duration`, `screen_width`, `last_alive`, session IDs, client IDs, or page-state decisions. All client data is treated as untrusted hints.

## Tracker Lifecycle

- Normal page views send only the current path, plus first-touch attribution if present.
- SPA navigation hooks send a new page view when the path changes.
- Exit pings use `visibilitychange` when the document becomes hidden.
- `pagehide` is kept only as a fallback.
- `beforeunload` is intentionally not used.

This follows browser guidance for small analytics payloads:
- [MDN `sendBeacon`](https://developer.mozilla.org/en-US/docs/Web/API/Navigator/sendBeacon): intended for small analytics/diagnostic POSTs and avoids slowing navigation.
- [MDN `visibilitychange`](https://developer.mozilla.org/en-US/docs/Web/API/Document/visibilitychange_event): the hidden transition is the last reliably observable lifecycle point for many pages.
- [W3C Beacon](https://www.w3.org/TR/beacon/): defines asynchronous beacon delivery for analytics-style data.

## Visitor Identification

Server-generated visitor ID computed from minimized request signals:
- Hash algorithm: truncated HMAC-SHA-256
- Key derivation: site-scoped, daily key derived from a server secret
- Inputs: internal site ID, truncated IP prefix, browser family, device class
- UTC-day-skipped rotation: the server checks today's and yesterday's hash
- Adjacent-day reuse: if only yesterday matches, the same client row is rewritten to today's hash
- New client only after a full UTC day was skipped
- No client-side storage or cookies
- Same visitor receives a consistent ID within the day and across an adjacent UTC midnight
- Country is not part of the visitor ID
- The server secret helps reduce the impact of database-only leaks by making visitor IDs harder to recompute outside the app

## IP Address Handling

Client IP is resolved by `server/internal/transport/http/clientip`:
- `RemoteAddr` is authoritative unless it belongs to a configured trusted proxy CIDR.
- `X-Forwarded-For` and `X-Real-IP` are ignored from untrusted remotes.
- For trusted proxy chains, Lovely Eye scans `X-Forwarded-For` from right to left and selects the last non-trusted hop as the client.
- If every forwarded hop is trusted, the leftmost valid forwarded IP is used.
- IPs are truncated before hashing: IPv4 `/24`, IPv6 `/64`.
- IPs are used only for visitor identity, block checks, rate limiting, and optional country lookup. They are not stored.

The default `TRUSTED_PROXY_CIDRS` covers loopback, RFC1918 private IPv4 ranges, and IPv6 unique-local addresses for common Docker, Nginx, Traefik, and Kubernetes private-network deployments. Public CDN or edge proxy ranges must be configured explicitly. The hop-selection behavior matches the intent of Nginx `real_ip_recursive`: [Nginx realip module](https://nginx.org/en/docs/http/ngx_http_realip_module.html).

## Session Mechanics

Sessions are computed server-side:
- 30-minute inactivity timeout.
- Server receive time is used for event time, session entry time, exit time, and duration.
- A normal page view creates or extends the active session.
- A same-path exit ping updates `exit_time`, `exit_path`, and computed duration only.
- Same-path exit pings may close a single-page session for up to `ANALYTICS_MAX_SINGLE_PAGE_DURATION`, defaulting to `4h`; repeated exit pings cannot push that single-page duration past the cap.
- A different-path exit ping counts that path as a page view only while the session is still inside the normal 30-minute active window, then updates session exit fields.
- An exit ping without an active session is a no-op.
- A 10-second dedupe window suppresses repeated same-path page views before counters change.
- Bounce rate uses sessions with one page view.
- Average session duration uses positive server-computed durations, including single-page sessions closed by exit pings.

## Query Accuracy

- Top pages and active pages count distinct analytics clients, not sessions.
- Page-view time series are bucketed by event time, not session start time.
- Visitor and session time series are bucketed by session entry time.
- Dashboard overview errors are propagated instead of returning partial zero values.

## Limits And Hardening

Limits exist because analytics endpoints are public by design:
- `ANALYTICS_MAX_BODY_BYTES` defaults to `16384`; tracker payloads should be far smaller.
- `ANALYTICS_MAX_PROPERTIES_BYTES` defaults to `8192`; custom event properties are allowlisted and capped.
- `ANALYTICS_MAX_SINGLE_PAGE_DURATION` defaults to `4h`; this bounds long-lived open tabs while still allowing real long reads.
- `ANALYTICS_RATE_LIMIT_ENABLED` defaults to `true`.
- `ANALYTICS_RATE_LIMIT_PER_MINUTE` defaults to `120`; collect traffic is limited by client IP before site lookup and by site key plus client IP after validation.
- `ANALYTICS_RATE_LIMIT_BURST` defaults to `240`; it uses the same keying as the refill limit.
- `GRAPHQL_MAX_BODY_BYTES` defaults to `1048576`.
- `GRAPHQL_MAX_COMPLEXITY` defaults to `300`; gqlgen rejects more expensive operations before resolver execution.
- `DASHBOARD_MAX_DAILY_RANGE_DAYS` defaults to `730`.
- `DASHBOARD_MAX_HOURLY_RANGE_DAYS` defaults to `31`.
- `DASHBOARD_MAX_FILTER_VALUES` defaults to `100`.
- `DASHBOARD_MAX_FILTER_STRING_LENGTH` defaults to `2048`.

These controls follow OWASP API guidance to constrain resource consumption and validate request sizes at the boundary: [OWASP REST Security Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/REST_Security_Cheat_Sheet.html).

## Bot Filtering

Filters non-human traffic:
- Search engines: Googlebot, Bingbot
- Social media bots: facebookexternalhit, Twitterbot
- Monitoring tools: UptimeRobot, Pingdom
- Scrapers: curl, wget, python-requests
- Headless browsers: Puppeteer, Playwright

## Query Parameters

- By default, query parameters are not included in tracked page paths.
- Use `data-include-query="true"` on the tracker script to include full query strings.

## Privacy

- No client-side cookies or persistent identifiers.
- Visitor IDs use UTC-day-skipped rotation.
- Visitor IDs are derived server-side from minimized signals.
- Site-scoped keying prevents reuse across sites.
- Keyed visitor IDs reduce the value of database-only leaks.
- IP addresses are never stored in the database.
- Country-level geolocation only, no city data.

## Event Allowlist

- Custom events are recorded only if the event name is allowlisted for the site.
- Event properties are filtered to the allowed keys and types.
- Required fields must be present for the event to be stored.
