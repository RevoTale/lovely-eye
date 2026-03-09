# Analytics Implementation

Implementation notes for Lovely Eye analytics.

## Visitor Identification

Server-generated visitor ID computed from minimized request signals:
- Hash algorithm: truncated HMAC-SHA-256
- Key derivation: site-scoped, daily key derived from a server secret
- Inputs: internal site ID, truncated IP prefix, browser family, device class
- Daily rotation: visitor ID changes every UTC day
- No client-side storage or cookies
- Same visitor receives consistent ID throughout the day
- Country is not part of the visitor ID
- The server secret helps reduce the impact of database-only leaks by making visitor IDs harder to recompute outside the app

## Bot Filtering

Filters non-human traffic:
- Search engines: Googlebot, Bingbot
- Social media bots: facebookexternalhit, Twitterbot
- Monitoring tools: UptimeRobot, Pingdom
- Scrapers: curl, wget, python-requests
- Headless browsers: Puppeteer, Playwright

## Page View Deduplication

Prevents duplicate counting:
- 10-second deduplication window per visitor per page
- Filters double-clicks, script reloads, same-path SPA updates within 10s
- Ensures accurate page view metrics

## Query Parameters

- By default, query parameters are not included in tracked page paths
- Use `data-include-query="true"` on the tracker script to include full query strings

## IP Address Handling

Extracts real client IP from proxied requests:
- Parses X-Forwarded-For header (first IP)
- Falls back to X-Real-IP header
- Strips port from RemoteAddr
- Truncates IP before hashing: IPv4 `/24`, IPv6 `/64`
- IPs used only for visitor identity and optional geolocation, never stored

## Session Management

Tracks browsing sessions:
- 30-minute inactivity timeout
- Records entry page, exit page, duration
- Calculates bounce rate
- Captures device, browser, OS, screen size
- Optional country detection via GeoIP

## Privacy

- No client-side cookies or persistent identifiers
- Visitor IDs rotate daily
- Visitor IDs are derived server-side from minimized signals
- Site-scoped keying prevents reuse across sites
- Keyed visitor IDs reduce the value of database-only leaks
- IP addresses never stored in database
- Country-level geolocation only (no city data)

## Event Allowlist

- Custom events are recorded only if the event name is allowlisted for the site
- Event properties are filtered to the allowed keys and types
- Required fields must be present for the event to be stored
