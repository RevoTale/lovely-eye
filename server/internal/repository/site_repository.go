package repository

import (
	"context"
	"fmt"

	"github.com/lovely-eye/server/internal/models"
	"github.com/uptrace/bun"
)

type SiteRepository struct {
	db *bun.DB
}

func NewSiteRepository(db *bun.DB) *SiteRepository {
	return &SiteRepository{db: db}
}

func (r *SiteRepository) Create(ctx context.Context, site *models.Site) error {
	_, err := r.db.NewInsert().Model(site).Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to create site: %w", err)
	}
	return nil
}

func (r *SiteRepository) GetByID(ctx context.Context, id int64) (*models.Site, error) {
	site := new(models.Site)
	err := r.db.NewSelect().
		Model(site).
		Where("id = ?", id).
		Relation("Domains", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Order("position ASC")
		}).
		Relation("BlockedIPs", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Order("ip ASC")
		}).
		Relation("BlockedCountries", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Order("country_code ASC")
		}).
		Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get site by id: %w", err)
	}
	return site, nil
}

func (r *SiteRepository) GetByPublicKey(ctx context.Context, publicKey string) (*models.Site, error) {
	site := new(models.Site)
	err := r.db.NewSelect().
		Model(site).
		Where("public_key = ?", publicKey).
		Relation("Domains", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Order("position ASC")
		}).
		Relation("BlockedIPs", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Order("ip ASC")
		}).
		Relation("BlockedCountries", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Order("country_code ASC")
		}).
		Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get site by public key: %w", err)
	}
	return site, nil
}

func (r *SiteRepository) GetByDomain(ctx context.Context, domain string) (*models.Site, error) {
	site := new(models.Site)
	err := r.db.NewSelect().
		Model(site).
		Join("JOIN site_domains AS sd ON sd.site_id = s.id").
		Where("sd.domain = ?", domain).
		Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get site by domain: %w", err)
	}
	return site, nil
}

func (r *SiteRepository) GetByUserID(ctx context.Context, userID int64) ([]*models.Site, error) {
	var sites []*models.Site
	err := r.db.NewSelect().
		Model(&sites).
		Where("user_id = ?", userID).
		Relation("Domains", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Order("position ASC")
		}).
		Relation("BlockedIPs", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Order("ip ASC")
		}).
		Relation("BlockedCountries", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Order("country_code ASC")
		}).
		Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get sites by user id: %w", err)
	}
	return sites, nil
}

func (r *SiteRepository) AnyGeoIPRequirement(ctx context.Context) (bool, error) {
	var exists bool
	err := r.db.NewSelect().
		Model((*models.Site)(nil)).
		ColumnExpr(`EXISTS (
			SELECT 1 FROM sites WHERE track_country = true
			UNION
			SELECT 1 FROM site_blocked_countries
		)`).
		Scan(ctx, &exists)
	if err != nil {
		return false, fmt.Errorf("failed to check geoip requirement: %w", err)
	}
	return exists, nil
}

func (r *SiteRepository) Update(ctx context.Context, site *models.Site) error {
	_, err := r.db.NewUpdate().Model(site).WherePK().Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to update site: %w", err)
	}
	return nil
}

func (r *SiteRepository) UpdateWithDomains(ctx context.Context, site *models.Site, domains []string) error {
	return r.UpdateWithRelations(ctx, site, domains, nil, nil)
}

func (r *SiteRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.NewDelete().Model((*models.Site)(nil)).Where("id = ?", id).Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete site: %w", err)
	}
	return nil
}

func (r *SiteRepository) CreateWithDomains(ctx context.Context, site *models.Site, domains []string) error {
	err := r.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.NewInsert().Model(site).Exec(ctx); err != nil {
			return fmt.Errorf("insert site: %w", err)
		}

		if len(domains) == 0 {
			return nil
		}

		siteDomains := make([]*models.SiteDomain, 0, len(domains))
		for index, domain := range domains {
			siteDomains = append(siteDomains, &models.SiteDomain{
				SiteID:   site.ID,
				Domain:   domain,
				Position: index,
			})
		}

		if _, err := tx.NewInsert().Model(&siteDomains).Exec(ctx); err != nil {
			return fmt.Errorf("insert site domains: %w", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to create site with domains: %w", err)
	}
	return nil
}

func (r *SiteRepository) UpdateWithRelations(ctx context.Context, site *models.Site, domains []string, blockedIPs []string, blockedCountries []string) error {
	err := r.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.NewUpdate().Model(site).WherePK().Exec(ctx); err != nil {
			return fmt.Errorf("update site: %w", err)
		}

		if domains != nil {
			if err := replaceSiteDomains(ctx, tx, site.ID, domains); err != nil {
				return err
			}
		}

		if blockedIPs != nil {
			if err := replaceBlockedIPs(ctx, tx, site.ID, blockedIPs); err != nil {
				return err
			}
		}

		if blockedCountries != nil {
			if err := replaceBlockedCountries(ctx, tx, site.ID, blockedCountries); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to update site with relations: %w", err)
	}
	return nil
}

func replaceSiteDomains(ctx context.Context, tx bun.Tx, siteID int64, domains []string) error {
	if _, err := tx.NewDelete().
		Model((*models.SiteDomain)(nil)).
		Where("site_id = ?", siteID).
		Exec(ctx); err != nil {
		return fmt.Errorf("failed to delete site domains: %w", err)
	}

	if len(domains) == 0 {
		return nil
	}

	siteDomains := make([]*models.SiteDomain, 0, len(domains))
	for index, domain := range domains {
		siteDomains = append(siteDomains, &models.SiteDomain{
			SiteID:   siteID,
			Domain:   domain,
			Position: index,
		})
	}

	_, err := tx.NewInsert().Model(&siteDomains).Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to insert site domains: %w", err)
	}
	return nil
}

func replaceBlockedIPs(ctx context.Context, tx bun.Tx, siteID int64, blockedIPs []string) error {
	if _, err := tx.NewDelete().
		Model((*models.SiteBlockedIP)(nil)).
		Where("site_id = ?", siteID).
		Exec(ctx); err != nil {
		return fmt.Errorf("failed to delete blocked ips: %w", err)
	}

	if len(blockedIPs) == 0 {
		return nil
	}

	entries := make([]*models.SiteBlockedIP, 0, len(blockedIPs))
	for _, ip := range blockedIPs {
		entries = append(entries, &models.SiteBlockedIP{
			SiteID: siteID,
			IP:     ip,
		})
	}

	_, err := tx.NewInsert().Model(&entries).Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to insert blocked ips: %w", err)
	}
	return nil
}

func replaceBlockedCountries(ctx context.Context, tx bun.Tx, siteID int64, blockedCountries []string) error {
	if _, err := tx.NewDelete().
		Model((*models.SiteBlockedCountry)(nil)).
		Where("site_id = ?", siteID).
		Exec(ctx); err != nil {
		return fmt.Errorf("failed to delete blocked countries: %w", err)
	}

	if len(blockedCountries) == 0 {
		return nil
	}

	entries := make([]*models.SiteBlockedCountry, 0, len(blockedCountries))
	for _, code := range blockedCountries {
		entries = append(entries, &models.SiteBlockedCountry{
			SiteID:      siteID,
			CountryCode: code,
		})
	}

	_, err := tx.NewInsert().Model(&entries).Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to insert blocked countries: %w", err)
	}
	return nil
}
