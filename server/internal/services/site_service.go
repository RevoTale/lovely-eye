package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/lovely-eye/server/internal/models"
	"github.com/lovely-eye/server/internal/repository"
	"github.com/lovely-eye/server/pkg/utils"
)

var (
	ErrSiteNotFound            = errors.New("site not found")
	ErrSiteExists              = errors.New("site with this domain already exists")
	ErrNotAuthorized           = errors.New("not authorized to access this site")
	ErrTooManyBlockedIPs       = errors.New("blocked IP list exceeds 500 entries")
	ErrTooManyBlockedCountries = errors.New("blocked country list exceeds 250 entries")
)

type SiteService struct {
	siteRepo *repository.SiteRepository
}

func NewSiteService(siteRepo *repository.SiteRepository) *SiteService {
	return &SiteService{siteRepo: siteRepo}
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

func (s *SiteService) Create(ctx context.Context, input CreateSiteInput) (*models.Site, error) {
	normalizedDomains, err := normalizeDomains(input.Domains)
	if err != nil {
		return nil, err
	}

	validatedName, err := utils.ValidateSiteName(input.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to validate site name: %w", err)
	}

	for _, domain := range normalizedDomains {
		existing, _ := s.siteRepo.GetByDomainForUser(ctx, input.UserID, domain)
		if existing != nil {
			return nil, ErrSiteExists
		}
	}

	publicKey, err := generatePublicKey()
	if err != nil {
		return nil, err
	}

	site := &models.Site{
		UserID:    input.UserID,
		Name:      validatedName,
		PublicKey: publicKey,
	}

	if err := s.siteRepo.CreateWithDomains(ctx, site, normalizedDomains); err != nil {
		return nil, fmt.Errorf("failed to create site with domains: %w", err)
	}

	site.Domains = buildSiteDomains(site.ID, normalizedDomains)
	return site, nil
}

func (s *SiteService) GetByID(ctx context.Context, id, userID int64) (*models.Site, error) {
	site, err := s.siteRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrSiteNotFound
	}

	if site.UserID != userID {
		return nil, ErrNotAuthorized
	}

	return site, nil
}

func (s *SiteService) GetByPublicKey(ctx context.Context, publicKey string) (*models.Site, error) {
	site, err := s.siteRepo.GetByPublicKey(ctx, publicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get site by public key: %w", err)
	}
	return site, nil
}

func (s *SiteService) GetUserSites(ctx context.Context, userID int64, limit, offset int) ([]*models.Site, error) {
	sites, err := s.siteRepo.GetByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get user sites: %w", err)
	}
	return sites, nil
}

func (s *SiteService) Update(ctx context.Context, id, userID int64, input UpdateSiteInput) (*models.Site, error) {
	validatedName, err := utils.ValidateSiteName(input.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to validate site name: %w", err)
	}

	site, err := s.siteRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrSiteNotFound
	}

	if site.UserID != userID {
		return nil, ErrNotAuthorized
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
			existing, _ := s.siteRepo.GetByDomainForUser(ctx, userID, domain)
			if existing != nil && existing.ID != site.ID {
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
		if err := s.siteRepo.Update(ctx, site); err != nil {
			return nil, fmt.Errorf("failed to update site: %w", err)
		}
		return site, nil
	}

	if err := s.siteRepo.UpdateWithRelations(ctx, site, normalizedDomains, normalizedBlockedIPs, normalizedBlockedCountries); err != nil {
		return nil, fmt.Errorf("failed to update site with relations: %w", err)
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

func (s *SiteService) Delete(ctx context.Context, id, userID int64) error {
	site, err := s.siteRepo.GetByID(ctx, id)
	if err != nil {
		return ErrSiteNotFound
	}

	if site.UserID != userID {
		return ErrNotAuthorized
	}

	if err := s.siteRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete site: %w", err)
	}
	return nil
}

func (s *SiteService) RegeneratePublicKey(ctx context.Context, id, userID int64) (*models.Site, error) {
	site, err := s.siteRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrSiteNotFound
	}

	if site.UserID != userID {
		return nil, ErrNotAuthorized
	}

	publicKey, err := generatePublicKey()
	if err != nil {
		return nil, err
	}

	site.PublicKey = publicKey
	if err := s.siteRepo.Update(ctx, site); err != nil {
		return nil, fmt.Errorf("failed to update site: %w", err)
	}

	return site, nil
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
		normalizedDomain, err := utils.ValidateDomain(domain)
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
		return nil, utils.ErrInvalidDomain
	}

	return normalized, nil
}

func buildSiteDomains(siteID int64, domains []string) []*models.SiteDomain {
	result := make([]*models.SiteDomain, 0, len(domains))
	for index, domain := range domains {
		result = append(result, &models.SiteDomain{
			SiteID:   siteID,
			Domain:   domain,
			Position: index,
		})
	}
	return result
}

func buildBlockedIPs(siteID int64, ips []string) []*models.SiteBlockedIP {
	result := make([]*models.SiteBlockedIP, 0, len(ips))
	for _, ip := range ips {
		result = append(result, &models.SiteBlockedIP{
			SiteID: siteID,
			IP:     ip,
		})
	}
	return result
}

func buildBlockedCountries(siteID int64, codes []string) []*models.SiteBlockedCountry {
	result := make([]*models.SiteBlockedCountry, 0, len(codes))
	for _, code := range codes {
		result = append(result, &models.SiteBlockedCountry{
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
		ip, err := utils.ValidateIPAddress(value)
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
		code, err := utils.ValidateCountryCode(value)
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
