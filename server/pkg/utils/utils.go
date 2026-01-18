package utils

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net"
	"regexp"
	"strings"
)

// GenerateRandomString generates a random hex string of the specified length
func GenerateRandomString(length int) (string, error) {
	bytes := make([]byte, length/2)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// GetClientIP extracts the client IP from request headers
func GetClientIP(xForwardedFor, xRealIP, remoteAddr string) string {
	// Check X-Forwarded-For header first
	if xForwardedFor != "" {
		ips := strings.Split(xForwardedFor, ",")
		if len(ips) > 0 {
			ip := strings.TrimSpace(ips[0])
			if net.ParseIP(ip) != nil {
				return ip
			}
		}
	}

	// Check X-Real-IP header
	if xRealIP != "" {
		if net.ParseIP(xRealIP) != nil {
			return xRealIP
		}
	}

	// Fall back to RemoteAddr
	if remoteAddr != "" {
		host, _, err := net.SplitHostPort(remoteAddr)
		if err == nil {
			return host
		}
		return remoteAddr
	}

	return ""
}

// NormalizeURL removes trailing slashes and normalizes paths
func NormalizeURL(path string) string {
	if path == "" {
		return "/"
	}

	// Ensure path starts with /
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	// Remove trailing slash (except for root)
	if len(path) > 1 && strings.HasSuffix(path, "/") {
		path = strings.TrimSuffix(path, "/")
	}

	return path
}

// TruncateString truncates a string to the specified length
func TruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen]
}

var (
	// ErrInvalidDomain indicates the domain format is invalid
	ErrInvalidDomain = errors.New("invalid domain format")
	// ErrInvalidSiteName indicates the site name is invalid
	ErrInvalidSiteName = errors.New("invalid site name")
	// ErrDomainTooLong indicates the domain exceeds maximum length
	ErrDomainTooLong = errors.New("domain name too long")
	// ErrSiteNameTooLong indicates the site name exceeds maximum length
	ErrSiteNameTooLong = errors.New("site name must be between 1 and 100 characters")
	// ErrInvalidIPAddress indicates the IP address format is invalid
	ErrInvalidIPAddress = errors.New("invalid IP address")
	// ErrInvalidCountryCode indicates the country code format is invalid
	ErrInvalidCountryCode = errors.New("invalid country code")
)

// Domain regex pattern for valid domain names
// Matches: example.com, sub.example.com, example-site.com
// Does not match: -example.com, example-.com, example..com, http://example.com
var domainRegex = regexp.MustCompile(`^[a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?(\.[a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?)*$`)
var countryCodeRegex = regexp.MustCompile(`^[A-Z]{2}$`)

// ValidateDomain validates and normalizes a domain name
func ValidateDomain(domain string) (string, error) {
	// Trim whitespace
	domain = strings.TrimSpace(domain)

	// Convert to lowercase first
	domain = strings.ToLower(domain)

	// Remove protocol if present
	domain = strings.TrimPrefix(domain, "https://")
	domain = strings.TrimPrefix(domain, "http://")

	// Remove www. prefix
	domain = strings.TrimPrefix(domain, "www.")

	// Remove path and trailing slashes
	if idx := strings.Index(domain, "/"); idx != -1 {
		domain = domain[:idx]
	}

	// Check if empty after normalization
	if domain == "" {
		return "", ErrInvalidDomain
	}

	// Check maximum length (253 characters is DNS limit)
	if len(domain) > 253 {
		return "", ErrDomainTooLong
	}

	// Validate format using regex
	if !domainRegex.MatchString(domain) {
		return "", ErrInvalidDomain
	}

	return domain, nil
}

// ValidateSiteName validates a site name
func ValidateSiteName(name string) (string, error) {
	// Trim whitespace
	name = strings.TrimSpace(name)

	// Check length (1-100 characters)
	if len(name) < 1 || len(name) > 100 {
		return "", ErrSiteNameTooLong
	}

	// Site name should not be empty or just whitespace
	if name == "" {
		return "", ErrInvalidSiteName
	}

	return name, nil
}

// ValidateIPAddress validates and normalizes an IP address.
func ValidateIPAddress(ip string) (string, error) {
	ip = strings.TrimSpace(ip)
	if ip == "" {
		return "", ErrInvalidIPAddress
	}
	parsed := net.ParseIP(ip)
	if parsed == nil {
		return "", ErrInvalidIPAddress
	}
	return parsed.String(), nil
}

// ValidateCountryCode validates and normalizes a 2-letter ISO country code.
func ValidateCountryCode(code string) (string, error) {
	code = strings.TrimSpace(strings.ToUpper(code))
	if !countryCodeRegex.MatchString(code) {
		return "", ErrInvalidCountryCode
	}
	return code, nil
}
