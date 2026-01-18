package graph

import (
	"strconv"

	"github.com/lovely-eye/server/internal/graph/model"
	"github.com/lovely-eye/server/internal/models"
)

func buildGraphQLSite(site *models.Site) *model.Site {
	return &model.Site{
		ID:               strconv.FormatInt(site.ID, 10),
		Domains:          siteDomains(site),
		Name:             site.Name,
		PublicKey:        site.PublicKey,
		TrackCountry:     site.TrackCountry,
		BlockedIPs:       siteBlockedIPs(site),
		BlockedCountries: siteBlockedCountries(site),
		CreatedAt:        site.CreatedAt,
	}
}

func siteDomains(site *models.Site) []string {
	domains := make([]string, 0, len(site.Domains))
	for _, domain := range site.Domains {
		if domain != nil && domain.Domain != "" {
			domains = append(domains, domain.Domain)
		}
	}
	return domains
}

func siteBlockedIPs(site *models.Site) []string {
	ips := make([]string, 0, len(site.BlockedIPs))
	for _, entry := range site.BlockedIPs {
		if entry != nil && entry.IP != "" {
			ips = append(ips, entry.IP)
		}
	}
	return ips
}

func siteBlockedCountries(site *models.Site) []string {
	countries := make([]string, 0, len(site.BlockedCountries))
	for _, entry := range site.BlockedCountries {
		if entry != nil && entry.CountryCode != "" {
			countries = append(countries, entry.CountryCode)
		}
	}
	return countries
}
