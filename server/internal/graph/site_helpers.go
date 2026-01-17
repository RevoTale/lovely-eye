package graph

import (
	"strconv"

	"github.com/lovely-eye/server/internal/graph/model"
	"github.com/lovely-eye/server/internal/models"
)

func buildGraphQLSite(site *models.Site) *model.Site {
	return &model.Site{
		ID:           strconv.FormatInt(site.ID, 10),
		Domains:      siteDomains(site),
		Name:         site.Name,
		PublicKey:    site.PublicKey,
		TrackCountry: site.TrackCountry,
		CreatedAt:    site.CreatedAt,
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
