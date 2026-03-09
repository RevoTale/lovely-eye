# Privacy Policy (Template)

Lovely Eye is self-hosted analytics. The site owner is the data controller. This template describes default behavior and can be adapted for your deployment.

## Data We Collect

- Page path (query parameters excluded by default)
- Referrer URL
- UTM source, medium, and campaign
- Device, browser, OS, screen size
- Session timing (entry/exit page, duration, bounce)
- Country (only if enabled)
- Custom events and allowlisted metadata fields

## Visitor Identifiers

We derive a keyed visitor identifier on the server from site ID, truncated IP prefix, browser family, and device class. The hash is computed per UTC day, but the server reuses the same client across `today` and `yesterday`; if only yesterday matches, that row is rewritten to today's hash. A new client is created only after a UTC day was skipped, so the identifier is still short-lived and not persistent. This keyed approach helps reduce the impact of database-only leaks because the stored analytics rows do not include enough information to recompute the identifier on their own.

## IP Addresses Under GDPR

IP addresses can be personal data. GDPR does not ban storing them, but it requires a lawful basis, minimized retention, security, and justification.
Lovely Eye uses IPs transiently for visitor identity and country lookup and does not store them by default. For identity, the IP is truncated before hashing.

## Data Retention

Retention is controlled by the site owner. Lovely Eye does not delete analytics data automatically.

## Cookies and Local Storage

Lovely Eye does not use cookies or local storage for analytics tracking by default.

## Contact

If you have questions about this policy, contact the site owner.
