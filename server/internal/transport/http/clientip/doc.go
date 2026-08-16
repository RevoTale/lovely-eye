// Package clientip resolves client IP addresses from remote addresses and
// common reverse-proxy headers.
//
// RemoteAddr is authoritative unless it belongs to a configured trusted proxy
// CIDR. Only trusted remotes may supply X-Forwarded-For or X-Real-IP. For
// X-Forwarded-For chains, the resolver scans from right to left and chooses
// the last non-trusted hop, matching the common real_ip_recursive deployment
// model used by Nginx-style reverse proxies.
package clientip
