package site

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
)

var (
	ErrSiteNotFound            = errors.New("site not found")
	ErrSiteExists              = errors.New("site with this domain already exists")
	ErrNotAuthorized           = errors.New("not authorized to access this site")
	ErrTooManyBlockedIPs       = errors.New("blocked IP list exceeds 500 entries")
	ErrTooManyBlockedCountries = errors.New("blocked country list exceeds 250 entries")
)

type Store interface {
	GetByID(ctx context.Context, id int64) (*Site, error)
	GetOwnerID(ctx context.Context, id int64) (int64, error)
	GetByPublicKey(ctx context.Context, publicKey string) (*Site, error)
	GetByUserID(ctx context.Context, userID int64, limit, offset int) ([]*Site, error)
	AnyGeoIPRequirement(ctx context.Context) (bool, error)
	DomainExistsForUser(ctx context.Context, userID int64, domain string, excludedSiteID int64) (bool, error)
	CreateWithDomains(ctx context.Context, site *Site, domains []string) error
	Update(ctx context.Context, site *Site) error
	UpdateWithRelations(ctx context.Context, site *Site, domains, blockedIPs, blockedCountries []string) error
	Delete(ctx context.Context, id int64) error
}

type Site struct {
	ID               int64
	UserID           int64
	Name             string
	PublicKey        string
	TrackCountry     bool
	CreatedAt        time.Time
	UpdatedAt        time.Time
	Domains          []*Domain
	BlockedIPs       []*BlockedIP
	BlockedCountries []*BlockedCountry
}

type Domain struct {
	ID        int64
	SiteID    int64
	Domain    string
	Position  int
	CreatedAt time.Time
	UpdatedAt time.Time
}

type BlockedIP struct {
	ID        int64
	SiteID    int64
	IP        string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type BlockedCountry struct {
	ID          int64
	SiteID      int64
	CountryCode string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Service struct {
	store Store
}

func NewService(store Store) *Service {
	return &Service{store: store}
}

type CreateSiteInput struct {
	Domains []string
	Name    string
	UserID  int64
}

type UpdateSiteInput struct {
	Name             string
	TrackCountry     *bool
	Domains          []string
	BlockedIPs       []string
	BlockedCountries []string
}

func (s *Service) Create(ctx context.Context, input CreateSiteInput) (*Site, error) {
	normalizedDomains, err := normalizeDomains(input.Domains)
	if err != nil {
		return nil, err
	}

	validatedName, err := ValidateSiteName(input.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to validate site name: %w", err)
	}

	for _, domain := range normalizedDomains {
		exists, err := s.store.DomainExistsForUser(ctx, input.UserID, domain, 0)
		if err != nil {
			return nil, fmt.Errorf("failed to check domain availability: %w", err)
		}
		if exists {
			return nil, ErrSiteExists
		}
	}

	publicKey, err := generatePublicKey()
	if err != nil {
		return nil, err
	}

	site := &Site{
		UserID:    input.UserID,
		Name:      validatedName,
		PublicKey: publicKey,
	}

	if err := s.store.CreateWithDomains(ctx, site, normalizedDomains); err != nil {
		return nil, fmt.Errorf("failed to create site with domains: %w", err)
	}

	site.Domains = buildSiteDomains(site.ID, normalizedDomains)
	return site, nil
}

func (s *Service) GetByID(ctx context.Context, id, userID int64) (*Site, error) {
	return s.getAuthorizedSite(ctx, id, userID)
}

func (s *Service) GetByPublicKey(ctx context.Context, publicKey string) (*Site, error) {
	site, err := s.store.GetByPublicKey(ctx, publicKey)
	if err != nil {
		if errors.Is(err, ErrSiteNotFound) {
			return nil, ErrSiteNotFound
		}
		return nil, fmt.Errorf("failed to get site by public key: %w", err)
	}
	return site, nil
}

func (s *Service) GetUserSites(ctx context.Context, userID int64, limit, offset int) ([]*Site, error) {
	sites, err := s.store.GetByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get user sites: %w", err)
	}
	return sites, nil
}

func (s *Service) Update(ctx context.Context, id, userID int64, input UpdateSiteInput) (*Site, error) {
	validatedName, err := ValidateSiteName(input.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to validate site name: %w", err)
	}

	site, err := s.getAuthorizedSite(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	site.Name = validatedName
	if input.TrackCountry != nil {
		site.TrackCountry = *input.TrackCountry
	}

	var normalizedDomains []string
	if input.Domains != nil {
		normalizedDomains, err = normalizeDomains(input.Domains)
		if err != nil {
			return nil, err
		}

		for _, domain := range normalizedDomains {
			exists, err := s.store.DomainExistsForUser(ctx, userID, domain, site.ID)
			if err != nil {
				return nil, fmt.Errorf("failed to check domain availability: %w", err)
			}
			if exists {
				return nil, ErrSiteExists
			}
		}
	}

	var normalizedBlockedIPs []string
	if input.BlockedIPs != nil {
		normalizedBlockedIPs, err = normalizeBlockedIPs(input.BlockedIPs)
		if err != nil {
			return nil, err
		}
	}

	var normalizedBlockedCountries []string
	if input.BlockedCountries != nil {
		normalizedBlockedCountries, err = normalizeBlockedCountries(input.BlockedCountries)
		if err != nil {
			return nil, err
		}
	}

	if input.Domains == nil && input.BlockedIPs == nil && input.BlockedCountries == nil {
		if err := s.store.Update(ctx, site); err != nil {
			return nil, classifySiteWriteError("update site", err)
		}
		return site, nil
	}

	if err := s.store.UpdateWithRelations(ctx, site, normalizedDomains, normalizedBlockedIPs, normalizedBlockedCountries); err != nil {
		return nil, classifySiteWriteError("update site with relations", err)
	}

	if input.Domains != nil {
		site.Domains = buildSiteDomains(site.ID, normalizedDomains)
	}
	if input.BlockedIPs != nil {
		site.BlockedIPs = buildBlockedIPs(site.ID, normalizedBlockedIPs)
	}
	if input.BlockedCountries != nil {
		site.BlockedCountries = buildBlockedCountries(site.ID, normalizedBlockedCountries)
	}
	return site, nil
}

func (s *Service) Delete(ctx context.Context, id, userID int64) error {
	_, err := s.getAuthorizedSite(ctx, id, userID)
	if err != nil {
		return err
	}

	if err := s.store.Delete(ctx, id); err != nil {
		return classifySiteWriteError("delete site", err)
	}
	return nil
}

func (s *Service) RegeneratePublicKey(ctx context.Context, id, userID int64) (*Site, error) {
	site, err := s.getAuthorizedSite(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	publicKey, err := generatePublicKey()
	if err != nil {
		return nil, err
	}

	site.PublicKey = publicKey
	if err := s.store.Update(ctx, site); err != nil {
		return nil, classifySiteWriteError("update site public key", err)
	}

	return site, nil
}

func (s *Service) getAuthorizedSite(ctx context.Context, id, userID int64) (*Site, error) {
	site, err := s.store.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrSiteNotFound) {
			return nil, ErrSiteNotFound
		}
		return nil, fmt.Errorf("failed to get site by id: %w", err)
	}
	if site.UserID != userID {
		return nil, ErrNotAuthorized
	}
	return site, nil
}

// RequireOwnership verifies access without loading relations that the caller does not consume.
func (s *Service) RequireOwnership(ctx context.Context, id, userID int64) error {
	ownerID, err := s.store.GetOwnerID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrSiteNotFound) {
			return ErrSiteNotFound
		}
		return fmt.Errorf("failed to get site owner: %w", err)
	}
	if ownerID != userID {
		return ErrNotAuthorized
	}
	return nil
}

func classifySiteWriteError(operation string, err error) error {
	if errors.Is(err, ErrSiteNotFound) {
		return ErrSiteNotFound
	}
	return fmt.Errorf("failed to %s: %w", operation, err)
}

func generatePublicKey() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}

func normalizeDomains(domains []string) ([]string, error) {
	normalized := make([]string, 0, len(domains))
	seen := make(map[string]struct{}, len(domains))
	for _, domain := range domains {
		normalizedDomain, err := ValidateDomain(domain)
		if err != nil {
			return nil, fmt.Errorf("failed to validate domain: %w", err)
		}
		if _, ok := seen[normalizedDomain]; ok {
			continue
		}
		seen[normalizedDomain] = struct{}{}
		normalized = append(normalized, normalizedDomain)
	}

	if len(normalized) == 0 {
		return nil, ErrInvalidDomain
	}

	return normalized, nil
}

func buildSiteDomains(siteID int64, domains []string) []*Domain {
	result := make([]*Domain, 0, len(domains))
	for index, domain := range domains {
		result = append(result, &Domain{
			SiteID:   siteID,
			Domain:   domain,
			Position: index,
		})
	}
	return result
}

func buildBlockedIPs(siteID int64, ips []string) []*BlockedIP {
	result := make([]*BlockedIP, 0, len(ips))
	for _, ip := range ips {
		result = append(result, &BlockedIP{
			SiteID: siteID,
			IP:     ip,
		})
	}
	return result
}

func buildBlockedCountries(siteID int64, codes []string) []*BlockedCountry {
	result := make([]*BlockedCountry, 0, len(codes))
	for _, code := range codes {
		result = append(result, &BlockedCountry{
			SiteID:      siteID,
			CountryCode: code,
		})
	}
	return result
}

func normalizeBlockedIPs(ips []string) ([]string, error) {
	normalized := make([]string, 0, len(ips))
	seen := make(map[string]struct{}, len(ips))
	for _, value := range ips {
		ip, err := ValidateIPAddress(value)
		if err != nil {
			return nil, fmt.Errorf("failed to validate IP address: %w", err)
		}
		if _, ok := seen[ip]; ok {
			continue
		}
		seen[ip] = struct{}{}
		normalized = append(normalized, ip)
	}

	if len(normalized) > 500 {
		return nil, ErrTooManyBlockedIPs
	}
	return normalized, nil
}

func normalizeBlockedCountries(countries []string) ([]string, error) {
	normalized := make([]string, 0, len(countries))
	seen := make(map[string]struct{}, len(countries))
	for _, value := range countries {
		code, err := ValidateCountryCode(value)
		if err != nil {
			return nil, fmt.Errorf("failed to validate country code: %w", err)
		}
		if _, ok := seen[code]; ok {
			continue
		}
		seen[code] = struct{}{}
		normalized = append(normalized, code)
	}

	if len(normalized) > 250 {
		return nil, ErrTooManyBlockedCountries
	}
	return normalized, nil
}
