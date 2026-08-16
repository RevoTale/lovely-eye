package analytics

import (
	"errors"
	"log/slog"
	"net"
	"net/url"
	"strings"

	"github.com/lovely-eye/server/internal/site"
)

func IsAllowedDomain(origin, referer string, domains []*site.Domain) bool {
	host := hostFromHeader(origin)
	if host == "" {
		host = hostFromHeader(referer)
	}
	if host == "" {
		return false
	}

	for _, domain := range domains {
		if domain != nil && domain.Domain == host {
			return true
		}
	}

	return false
}

func hostFromHeader(raw string) string {
	if raw == "" {
		return ""
	}

	parsed, err := url.Parse(raw)
	if err != nil {
		return ""
	}

	host := parsed.Hostname()
	if host == "" {
		return ""
	}

	normalized, err := site.ValidateDomain(host)
	if err != nil {
		return ""
	}

	return normalized
}

func (s *Service) isBlockedRequest(site *site.Site, ip string) bool {
	if site == nil {
		return false
	}
	if ip == "" {
		return false
	}

	if s.isIPBlocked(site.BlockedIPs, ip) {
		return true
	}

	return s.isCountryBlocked(site.BlockedCountries, ip)
}

func (s *Service) isIPBlocked(blocked []*site.BlockedIP, ip string) bool {
	if len(blocked) == 0 {
		return false
	}
	parsed := net.ParseIP(strings.TrimSpace(ip))
	if parsed == nil {
		return false
	}
	normalized := parsed.String()
	for _, entry := range blocked {
		if entry != nil && entry.IP == normalized {
			return true
		}
	}
	return false
}

func (s *Service) isCountryBlocked(blocked []*site.BlockedCountry, ip string) bool {
	if len(blocked) == 0 || s.geoIPService == nil {
		return false
	}
	c := s.resolveCountryBestEffort(ip)

	if c == (Country{}) || c == UnknownCountry || c == LocalNetworkCountry {
		return false
	}
	for _, entry := range blocked {
		if entry != nil && strings.EqualFold(entry.CountryCode, c.ISOCode) {
			return true
		}
	}
	return false
}

func (s *Service) resolveCountryBestEffort(ip string) Country {
	if s.geoIPService == nil {
		return UnknownCountry
	}

	country, err := s.geoIPService.ResolveCountry(ip)
	if err == nil {
		return country
	}

	if errors.Is(err, ErrNoDBReader) {
		return UnknownCountry
	}

	slog.Error("country resolve failed", "error", err)
	return UnknownCountry
}
