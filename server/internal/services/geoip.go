package services

import (
	"net"
	"sync"

	"github.com/oschwald/geoip2-golang"
)

// GeoIPService provides IP geolocation capabilities
type GeoIPService struct {
	reader *geoip2.Reader
	mu     sync.RWMutex
}

// NewGeoIPService creates a new GeoIP service
// dbPath should point to the GeoLite2-Country.mmdb file
// If the database file doesn't exist, the service will work but return "Unknown" for all lookups
func NewGeoIPService(dbPath string) (*GeoIPService, error) {
	if dbPath == "" {
		// No database path provided, return a service that always returns Unknown
		return &GeoIPService{}, nil
	}

	reader, err := geoip2.Open(dbPath)
	if err != nil {
		// If we can't open the database, log the error but continue
		// The service will return "Unknown" for all lookups
		return &GeoIPService{}, nil
	}

	return &GeoIPService{
		reader: reader,
	}, nil
}

// GetCountry returns the ISO country code for an IP address
func (g *GeoIPService) GetCountry(ipStr string) string {
	g.mu.RLock()
	defer g.mu.RUnlock()

	if g.reader == nil {
		return "Unknown"
	}

	ip := net.ParseIP(ipStr)
	if ip == nil {
		return "Unknown"
	}

	// Skip private/local IPs
	if isPrivateIP(ip) {
		return "Local"
	}

	record, err := g.reader.Country(ip)
	if err != nil {
		return "Unknown"
	}

	if record.Country.IsoCode == "" {
		return "Unknown"
	}

	return record.Country.IsoCode
}

// GetCountryName returns the country name for an IP address
func (g *GeoIPService) GetCountryName(ipStr string) string {
	g.mu.RLock()
	defer g.mu.RUnlock()

	if g.reader == nil {
		return "Unknown"
	}

	ip := net.ParseIP(ipStr)
	if ip == nil {
		return "Unknown"
	}

	if isPrivateIP(ip) {
		return "Local Network"
	}

	record, err := g.reader.Country(ip)
	if err != nil {
		return "Unknown"
	}

	if record.Country.Names["en"] == "" {
		return "Unknown"
	}

	return record.Country.Names["en"]
}

// Close closes the GeoIP database reader
func (g *GeoIPService) Close() error {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.reader != nil {
		return g.reader.Close()
	}
	return nil
}

// isPrivateIP checks if an IP is private/local
func isPrivateIP(ip net.IP) bool {
	if ip.IsLoopback() || ip.IsPrivate() {
		return true
	}
	return false
}
