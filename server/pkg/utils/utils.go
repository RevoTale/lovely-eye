package utils

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"regexp"
	"strings"
)

func GenerateRandomString(length int) (string, error) {
	bytes := make([]byte, length/2)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random string: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}

func GetClientIP(xForwardedFor, xRealIP, remoteAddr string) string {

	if xForwardedFor != "" {
		ips := strings.Split(xForwardedFor, ",")
		if len(ips) > 0 {
			ip := strings.TrimSpace(ips[0])
			if net.ParseIP(ip) != nil {
				return ip
			}
		}
	}

	if xRealIP != "" {
		if net.ParseIP(xRealIP) != nil {
			return xRealIP
		}
	}

	if remoteAddr != "" {
		host, _, err := net.SplitHostPort(remoteAddr)
		if err == nil {
			return host
		}
		return remoteAddr
	}

	return ""
}

func NormalizeURL(path string) string {
	if path == "" {
		return "/"
	}

	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	if len(path) > 1 && strings.HasSuffix(path, "/") {
		path = strings.TrimSuffix(path, "/")
	}

	return path
}

func TruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen]
}

var (
	ErrInvalidDomain = errors.New("invalid domain format")

	ErrInvalidSiteName = errors.New("invalid site name")

	ErrDomainTooLong = errors.New("domain name too long")

	ErrSiteNameTooLong = errors.New("site name must be between 1 and 100 characters")

	ErrInvalidIPAddress = errors.New("invalid IP address")

	ErrInvalidCountryCode = errors.New("invalid country code")
)

// Domain regex pattern for valid domain names
// Matches: example.com, sub.example.com, example-site.com
// Does not match: -example.com, example-.com, example..com, http://example.com
var domainRegex = regexp.MustCompile(`^[a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?(\.[a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?)*$`)
var countryCodeRegex = regexp.MustCompile(`^[A-Z]{2}$`)

func ValidateDomain(domain string) (string, error) {

	domain = strings.TrimSpace(domain)

	domain = strings.ToLower(domain)

	domain = strings.TrimPrefix(domain, "https://")
	domain = strings.TrimPrefix(domain, "http://")

	domain = strings.TrimPrefix(domain, "www.")

	if idx := strings.Index(domain, "/"); idx != -1 {
		domain = domain[:idx]
	}

	if domain == "" {
		return "", ErrInvalidDomain
	}

	if len(domain) > 253 {
		return "", ErrDomainTooLong
	}

	// Validate format using regex
	if !domainRegex.MatchString(domain) {
		return "", ErrInvalidDomain
	}

	return domain, nil
}

func ValidateSiteName(name string) (string, error) {

	name = strings.TrimSpace(name)

	if len(name) < 1 || len(name) > 100 {
		return "", ErrSiteNameTooLong
	}

	if name == "" {
		return "", ErrInvalidSiteName
	}

	return name, nil
}

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

func ValidateCountryCode(code string) (string, error) {
	code = strings.TrimSpace(strings.ToUpper(code))
	if !countryCodeRegex.MatchString(code) {
		return "", ErrInvalidCountryCode
	}
	return code, nil
}
