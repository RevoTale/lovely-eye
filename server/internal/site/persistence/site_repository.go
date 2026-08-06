package persistence

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	sitefeature "github.com/lovely-eye/server/internal/site"
	"github.com/uptrace/bun"
)

type Repository struct {
	db *bun.DB
}

var _ sitefeature.Store = (*Repository)(nil)

func New(db *bun.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetByID(ctx context.Context, id int64) (*sitefeature.Site, error) {
	site := new(Site)
	err := r.db.NewSelect().
		Model(site).
		Where("id = ?", id).
		Relation("Domains", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Order("position ASC", "id ASC")
		}).
		Relation("BlockedIPs", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Order("ip ASC")
		}).
		Relation("BlockedCountries", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Order("country_code ASC")
		}).
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("failed to get site by id: %w", sitefeature.ErrSiteNotFound)
		}
		return nil, fmt.Errorf("failed to get site by id: %w", err)
	}
	return siteFromModel(site), nil
}

func (r *Repository) GetOwnerID(ctx context.Context, id int64) (int64, error) {
	var ownerID int64
	err := r.db.NewSelect().
		Model((*Site)(nil)).
		Column("user_id").
		Where("id = ?", id).
		Scan(ctx, &ownerID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, sitefeature.ErrSiteNotFound
		}
		return 0, fmt.Errorf("failed to get site owner: %w", err)
	}
	return ownerID, nil
}

func (r *Repository) GetByPublicKey(ctx context.Context, publicKey string) (*sitefeature.Site, error) {
	site := new(Site)
	err := r.db.NewSelect().
		Model(site).
		Where("public_key = ?", publicKey).
		Relation("Domains", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Order("position ASC", "id ASC")
		}).
		Relation("BlockedIPs", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Order("ip ASC")
		}).
		Relation("BlockedCountries", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Order("country_code ASC")
		}).
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("failed to get site by public key: %w", sitefeature.ErrSiteNotFound)
		}
		return nil, fmt.Errorf("failed to get site by public key: %w", err)
	}
	return siteFromModel(site), nil
}

func (r *Repository) GetByDomainForUser(
	ctx context.Context,
	userID int64,
	domain string,
) (*sitefeature.Site, error) {
	site := new(Site)
	err := r.db.NewSelect().
		Model(site).
		Join("JOIN site_domains AS sd ON sd.site_id = s.id").
		Where("s.user_id = ?", userID).
		Where("sd.domain = ?", domain).
		Order("s.id ASC").
		Limit(1).
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("failed to get site by domain for user: %w", sitefeature.ErrSiteNotFound)
		}
		return nil, fmt.Errorf("failed to get site by domain for user: %w", err)
	}
	return siteFromModel(site), nil
}

func (r *Repository) DomainExistsForUser(
	ctx context.Context,
	userID int64,
	domain string,
	excludedSiteID int64,
) (bool, error) {
	query := r.db.NewSelect().
		Model((*Domain)(nil)).
		Join("JOIN sites AS s ON s.id = sd.site_id").
		Where("s.user_id = ?", userID).
		Where("sd.domain = ?", domain)
	if excludedSiteID != 0 {
		query = query.Where("sd.site_id <> ?", excludedSiteID)
	}
	exists, err := query.Exists(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to check site domain: %w", err)
	}
	return exists, nil
}

func (r *Repository) GetByUserID(
	ctx context.Context,
	userID int64,
	limit,
	offset int,
) ([]*sitefeature.Site, error) {
	var sites []*Site
	q := r.db.NewSelect().
		Model(&sites).
		Where("user_id = ?", userID).
		Relation("Domains", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Order("position ASC", "id ASC")
		}).
		Relation("BlockedIPs", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Order("ip ASC")
		}).
		Relation("BlockedCountries", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Order("country_code ASC")
		})
	if limit > 0 {
		q = q.Limit(limit)
	}
	if offset > 0 {
		q = q.Offset(offset)
	}
	err := q.Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get sites by user id: %w", err)
	}
	result := make([]*sitefeature.Site, 0, len(sites))
	for _, site := range sites {
		result = append(result, siteFromModel(site))
	}
	return result, nil
}

func (r *Repository) AnyGeoIPRequirement(ctx context.Context) (bool, error) {
	var exists bool
	err := r.db.NewSelect().
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

func (r *Repository) Update(ctx context.Context, site *sitefeature.Site) error {
	row := siteModel(site)
	result, err := r.db.NewUpdate().Model(row).WherePK().Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to update site: %w", err)
	}
	if err := requireAffectedSite(result, "update site"); err != nil {
		return err
	}
	copySite(site, row)
	return nil
}

func (r *Repository) Delete(ctx context.Context, id int64) error {
	err := r.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if err := deleteSiteAnalytics(ctx, tx, id); err != nil {
			return err
		}
		if err := deleteSiteConfiguration(ctx, tx, id); err != nil {
			return err
		}

		result, err := tx.NewDelete().Model((*Site)(nil)).Where("id = ?", id).Exec(ctx)
		if err != nil {
			return fmt.Errorf("delete site: %w", err)
		}
		return requireAffectedSite(result, "delete site")
	})
	if err != nil {
		return fmt.Errorf("delete site transaction: %w", err)
	}
	return nil
}

func deleteSiteAnalytics(ctx context.Context, tx bun.Tx, siteID int64) error {
	if _, err := tx.NewDelete().
		Model((*ownedEventData)(nil)).
		Where("event_id IN (SELECT e.id FROM events AS e JOIN sessions AS s ON s.id = e.session_id WHERE s.site_id = ?)", siteID).
		Exec(ctx); err != nil {
		return fmt.Errorf("delete site event data: %w", err)
	}
	if _, err := tx.NewDelete().
		Model((*ownedEvent)(nil)).
		Where("session_id IN (SELECT id FROM sessions WHERE site_id = ?)", siteID).
		Exec(ctx); err != nil {
		return fmt.Errorf("delete site events: %w", err)
	}
	if _, err := tx.NewDelete().
		Model((*ownedEventDefinitionField)(nil)).
		Where("event_definition_id IN (SELECT id FROM event_definitions WHERE site_id = ?)", siteID).
		Exec(ctx); err != nil {
		return fmt.Errorf("delete site event definition fields: %w", err)
	}
	if _, err := tx.NewDelete().
		Model((*ownedEventDefinition)(nil)).
		Where("site_id = ?", siteID).
		Exec(ctx); err != nil {
		return fmt.Errorf("delete site event definitions: %w", err)
	}
	if _, err := tx.NewDelete().
		Model((*ownedSession)(nil)).
		Where("site_id = ?", siteID).
		Exec(ctx); err != nil {
		return fmt.Errorf("delete site sessions: %w", err)
	}
	if _, err := tx.NewDelete().
		Model((*ownedClient)(nil)).
		Where("site_id = ?", siteID).
		Exec(ctx); err != nil {
		return fmt.Errorf("delete site clients: %w", err)
	}
	return nil
}

func deleteSiteConfiguration(ctx context.Context, tx bun.Tx, siteID int64) error {
	if _, err := tx.NewDelete().
		Model((*BlockedCountry)(nil)).
		Where("site_id = ?", siteID).
		Exec(ctx); err != nil {
		return fmt.Errorf("delete site blocked countries: %w", err)
	}
	if _, err := tx.NewDelete().
		Model((*BlockedIP)(nil)).
		Where("site_id = ?", siteID).
		Exec(ctx); err != nil {
		return fmt.Errorf("delete site blocked IPs: %w", err)
	}
	if _, err := tx.NewDelete().
		Model((*Domain)(nil)).
		Where("site_id = ?", siteID).
		Exec(ctx); err != nil {
		return fmt.Errorf("delete site domains: %w", err)
	}
	return nil
}

func (r *Repository) CreateWithDomains(
	ctx context.Context,
	site *sitefeature.Site,
	domains []string,
) error {
	row := siteModel(site)
	err := r.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.NewInsert().Model(row).Exec(ctx); err != nil {
			return fmt.Errorf("insert site: %w", err)
		}

		if len(domains) == 0 {
			return nil
		}

		siteDomains := make([]*Domain, 0, len(domains))
		for index, domain := range domains {
			siteDomains = append(siteDomains, &Domain{
				SiteID:   row.ID,
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
	copySite(site, row)
	return nil
}

func (r *Repository) UpdateWithRelations(
	ctx context.Context,
	site *sitefeature.Site,
	domains,
	blockedIPs,
	blockedCountries []string,
) error {
	row := siteModel(site)
	err := r.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		result, err := tx.NewUpdate().Model(row).WherePK().Exec(ctx)
		if err != nil {
			return fmt.Errorf("update site: %w", err)
		}
		if err := requireAffectedSite(result, "update site"); err != nil {
			return err
		}

		if domains != nil {
			if err := replaceSiteDomains(ctx, tx, row.ID, domains); err != nil {
				return err
			}
		}

		if blockedIPs != nil {
			if err := replaceBlockedIPs(ctx, tx, row.ID, blockedIPs); err != nil {
				return err
			}
		}

		if blockedCountries != nil {
			if err := replaceBlockedCountries(ctx, tx, row.ID, blockedCountries); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to update site with relations: %w", err)
	}
	copySite(site, row)
	return nil
}

func requireAffectedSite(result sql.Result, operation string) error {
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s rows affected: %w", operation, err)
	}
	if affected == 0 {
		return fmt.Errorf("%s: %w", operation, sitefeature.ErrSiteNotFound)
	}
	return nil
}

func siteFromModel(row *Site) *sitefeature.Site {
	site := &sitefeature.Site{
		ID:           row.ID,
		UserID:       row.UserID,
		Name:         row.Name,
		PublicKey:    row.PublicKey,
		TrackCountry: row.TrackCountry,
		CreatedAt:    row.CreatedAt,
		UpdatedAt:    row.UpdatedAt,
	}
	for _, domain := range row.Domains {
		if domain == nil {
			continue
		}
		site.Domains = append(site.Domains, &sitefeature.Domain{
			ID:        domain.ID,
			SiteID:    domain.SiteID,
			Domain:    domain.Domain,
			Position:  domain.Position,
			CreatedAt: domain.CreatedAt,
			UpdatedAt: domain.UpdatedAt,
		})
	}
	for _, blocked := range row.BlockedIPs {
		if blocked == nil {
			continue
		}
		site.BlockedIPs = append(site.BlockedIPs, &sitefeature.BlockedIP{
			ID:        blocked.ID,
			SiteID:    blocked.SiteID,
			IP:        blocked.IP,
			CreatedAt: blocked.CreatedAt,
			UpdatedAt: blocked.UpdatedAt,
		})
	}
	for _, blocked := range row.BlockedCountries {
		if blocked == nil {
			continue
		}
		site.BlockedCountries = append(site.BlockedCountries, &sitefeature.BlockedCountry{
			ID:          blocked.ID,
			SiteID:      blocked.SiteID,
			CountryCode: blocked.CountryCode,
			CreatedAt:   blocked.CreatedAt,
			UpdatedAt:   blocked.UpdatedAt,
		})
	}
	return site
}

func siteModel(site *sitefeature.Site) *Site {
	return &Site{
		ID:           site.ID,
		UserID:       site.UserID,
		Name:         site.Name,
		PublicKey:    site.PublicKey,
		TrackCountry: site.TrackCountry,
		CreatedAt:    site.CreatedAt,
		UpdatedAt:    site.UpdatedAt,
	}
}

func copySite(destination *sitefeature.Site, source *Site) {
	domains := destination.Domains
	blockedIPs := destination.BlockedIPs
	blockedCountries := destination.BlockedCountries
	*destination = *siteFromModel(source)
	destination.Domains = domains
	destination.BlockedIPs = blockedIPs
	destination.BlockedCountries = blockedCountries
}

func replaceSiteDomains(ctx context.Context, tx bun.Tx, siteID int64, domains []string) error {
	if _, err := tx.NewDelete().
		Model((*Domain)(nil)).
		Where("site_id = ?", siteID).
		Exec(ctx); err != nil {
		return fmt.Errorf("failed to delete site domains: %w", err)
	}

	if len(domains) == 0 {
		return nil
	}

	siteDomains := make([]*Domain, 0, len(domains))
	for index, domain := range domains {
		siteDomains = append(siteDomains, &Domain{
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
		Model((*BlockedIP)(nil)).
		Where("site_id = ?", siteID).
		Exec(ctx); err != nil {
		return fmt.Errorf("failed to delete blocked ips: %w", err)
	}

	if len(blockedIPs) == 0 {
		return nil
	}

	entries := make([]*BlockedIP, 0, len(blockedIPs))
	for _, ip := range blockedIPs {
		entries = append(entries, &BlockedIP{
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
		Model((*BlockedCountry)(nil)).
		Where("site_id = ?", siteID).
		Exec(ctx); err != nil {
		return fmt.Errorf("failed to delete blocked countries: %w", err)
	}

	if len(blockedCountries) == 0 {
		return nil
	}

	entries := make([]*BlockedCountry, 0, len(blockedCountries))
	for _, code := range blockedCountries {
		entries = append(entries, &BlockedCountry{
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
