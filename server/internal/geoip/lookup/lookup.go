package lookup

import (
	"fmt"
	"log/slog"
	"net/netip"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/lovely-eye/server/internal/geoip"
	"github.com/oschwald/geoip2-golang/v2"
	"github.com/oschwald/maxminddb-golang/v2"
)

type Service struct {
	dbPath string

	reader *geoip2.Reader
	mu     sync.RWMutex
}

func New(dbPath string) *Service {
	return &Service{dbPath: dbPath}
}

func (s *Service) DBPath() string {
	return s.dbPath
}

func (s *Service) FileExists() bool {
	info, err := os.Stat(s.dbPath)
	return err == nil && !info.IsDir()
}

func (s *Service) UpdatedAt() *time.Time {
	info, err := os.Stat(s.dbPath)
	if err != nil || info.IsDir() {
		return nil
	}
	modTime := info.ModTime()
	return &modTime
}

func (s *Service) Load() error {
	reader, err := geoip2.Open(s.dbPath)
	if err != nil {
		return fmt.Errorf("open GeoIP database: %w", err)
	}

	s.mu.Lock()
	if s.reader != nil {
		_ = s.reader.Close()
	}
	s.reader = reader
	s.mu.Unlock()

	return nil
}

func (s *Service) HasReader() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.reader != nil
}

// ListCountries scans the whole MMDB to build a unique country list.
// Do not use it on hot paths or for frequent polling.
func (s *Service) ListCountries(search string) ([]geoip.ListedCountry, error) {
	reader, err := maxminddb.Open(s.dbPath)
	if err != nil {
		return nil, fmt.Errorf("open GeoIP database: %w", err)
	}
	defer func() {
		if err := reader.Close(); err != nil {
			slog.Error("failed to close maxminddb reader", "error", err)
		}
	}()

	type countryRecord struct {
		Country struct {
			ISOCode string            `maxminddb:"iso_code"`
			Names   map[string]string `maxminddb:"names"`
		} `maxminddb:"country"`
		RegisteredCountry struct {
			ISOCode string            `maxminddb:"iso_code"`
			Names   map[string]string `maxminddb:"names"`
		} `maxminddb:"registered_country"`
	}

	seen := make(map[string]string)
	for result := range reader.Networks() {
		var record countryRecord
		if err := result.Decode(&record); err != nil {
			return nil, fmt.Errorf("read GeoIP network: %w", err)
		}
		if !result.Found() {
			continue
		}

		code := record.Country.ISOCode
		name := record.Country.Names["en"]
		if code == "" {
			code = record.RegisteredCountry.ISOCode
			name = record.RegisteredCountry.Names["en"]
		}
		if code == "" {
			continue
		}
		if name == "" {
			name = code
		}
		if _, ok := seen[code]; ok {
			continue
		}
		seen[code] = name
	}

	countries := make([]geoip.ListedCountry, 0, len(seen))
	for code, name := range seen {
		countries = append(countries, geoip.ListedCountry{
			Code: code,
			Name: name,
		})
	}
	sort.Slice(countries, func(i, j int) bool {
		if countries[i].Name == countries[j].Name {
			return countries[i].Code < countries[j].Code
		}
		return countries[i].Name < countries[j].Name
	})

	return filterCountries(countries, search), nil
}

func (s *Service) ResolveCountry(ipStr string) (geoip.Country, error) {
	ip, err := netip.ParseAddr(ipStr)
	if err != nil {
		return geoip.UnknownCountry, fmt.Errorf("failed parse country IP: %s", err.Error())
	}

	if !ip.IsGlobalUnicast() || ip.IsPrivate() {
		return geoip.LocalNetworkCountry, nil
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.reader == nil {
		return geoip.UnknownCountry, geoip.ErrNoDBReader
	}

	record, err := s.reader.Country(ip)
	if err != nil {
		return geoip.UnknownCountry, fmt.Errorf("failed to get country: %s", err.Error())
	}

	name := record.Country.Names.English
	code := record.Country.ISOCode
	if code == "" {
		name = record.RegisteredCountry.Names.English
		code = record.RegisteredCountry.ISOCode
	}
	if name == "" {
		name = code
	}
	if code == "" {
		return geoip.UnknownCountry, nil
	}

	return geoip.Country{Name: name, ISOCode: code}, nil
}

func (s *Service) Close() error {
	s.mu.Lock()
	reader := s.reader
	s.reader = nil
	s.mu.Unlock()

	if reader != nil {
		if err := reader.Close(); err != nil {
			return fmt.Errorf("failed to close geoip reader: %w", err)
		}
	}

	return nil
}

func filterCountries(countries []geoip.ListedCountry, search string) []geoip.ListedCountry {
	query := strings.TrimSpace(strings.ToLower(search))
	if query == "" {
		result := make([]geoip.ListedCountry, len(countries))
		copy(result, countries)
		return result
	}

	filtered := make([]geoip.ListedCountry, 0, len(countries))
	for _, country := range countries {
		code := strings.ToLower(country.Code)
		name := strings.ToLower(country.Name)
		if strings.Contains(code, query) || strings.Contains(name, query) {
			filtered = append(filtered, country)
		}
	}

	return filtered
}
