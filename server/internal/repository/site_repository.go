package repository

import (
	"context"

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
	return err
}

func (r *SiteRepository) GetByID(ctx context.Context, id int64) (*models.Site, error) {
	site := new(models.Site)
	err := r.db.NewSelect().Model(site).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return site, nil
}

func (r *SiteRepository) GetByPublicKey(ctx context.Context, publicKey string) (*models.Site, error) {
	site := new(models.Site)
	err := r.db.NewSelect().Model(site).Where("public_key = ?", publicKey).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return site, nil
}

func (r *SiteRepository) GetByDomain(ctx context.Context, domain string) (*models.Site, error) {
	site := new(models.Site)
	err := r.db.NewSelect().Model(site).Where("domain = ?", domain).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return site, nil
}

func (r *SiteRepository) GetByUserID(ctx context.Context, userID int64) ([]*models.Site, error) {
	var sites []*models.Site
	err := r.db.NewSelect().Model(&sites).Where("user_id = ?", userID).Scan(ctx)
	return sites, err
}

func (r *SiteRepository) Update(ctx context.Context, site *models.Site) error {
	_, err := r.db.NewUpdate().Model(site).WherePK().Exec(ctx)
	return err
}

func (r *SiteRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.NewDelete().Model((*models.Site)(nil)).Where("id = ?", id).Exec(ctx)
	return err
}
