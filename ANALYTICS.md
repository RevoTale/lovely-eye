# Analytics Implementation

Privacy-focused web analytics for Lovely Eye.

## Visitor Identification

Anonymous hash computed from IP address, user agent, site key, and current date:
- Hash algorithm: SHA256
- Daily rotation: visitor ID changes every 24 hours
- No persistent storage or cookies
- Same visitor receives consistent ID throughout the day

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
- Filters double-clicks, script reloads, SPA navigation
- Ensures accurate page view metrics

## Query Parameters

- By default, query parameters are not included in tracked page paths
- UTM parameters are captured separately for campaign reporting
- Use `data-include-query="true"` on the tracker script to include full query strings

## IP Address Handling

Extracts real client IP from proxied requests:
- Parses X-Forwarded-For header (first IP)
- Falls back to X-Real-IP header
- Strips port from RemoteAddr
- IPs used only for hashing and geolocation, never stored

## Session Management

Tracks browsing sessions:
- 30-minute inactivity timeout
- Records entry page, exit page, duration
- Calculates bounce rate
- Captures device, browser, OS, screen size
- Optional country detection via GeoIP

## Privacy

- No cookies or persistent identifiers
- Visitor IDs rotate daily
- Site key prevents cross-site tracking
- IP addresses never stored in database
- Country-level geolocation only (no city data)

## Event Allowlist

- Custom events are recorded only if the event name is allowlisted for the site
- Event properties are filtered to the allowed keys and types
- Required fields must be present for the event to be stored
