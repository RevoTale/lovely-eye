package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/lovely-eye/server/internal/models"
	"github.com/uptrace/bun"
)

type CountryRepository struct {
	db *bun.DB
}

func NewCountryRepository(db *bun.DB) *CountryRepository {
	return &CountryRepository{db: db}
}

func (r *CountryRepository) UpsertCountries(ctx context.Context, countries []models.Country) error {
	if len(countries) == 0 {
		return nil
	}

	_, err := r.db.NewInsert().
		Model(&countries).
		On("CONFLICT (code) DO UPDATE").
		Set("name = EXCLUDED.name").
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("upsert countries: %w", err)
	}

	return nil
}

func (r *CountryRepository) SearchCountries(ctx context.Context, search string, limit, offset int) ([]models.Country, error) {
	var countries []models.Country

	q := r.db.NewSelect().
		Model(&countries).
		Order("co.name ASC", "co.code ASC")

	trimmedSearch := strings.ToLower(strings.TrimSpace(search))
	if trimmedSearch != "" {
		searchPattern := "%" + trimmedSearch + "%"
		q = q.Where("LOWER(co.code) LIKE ? OR LOWER(co.name) LIKE ?", searchPattern, searchPattern)
	}

	if limit > 0 {
		q = q.Limit(limit)
	}
	if offset > 0 {
		q = q.Offset(offset)
	}

	if err := q.Scan(ctx); err != nil {
		return nil, fmt.Errorf("search countries: %w", err)
	}

	return countries, nil
}

func (r *CountryRepository) GetCountriesByCodes(ctx context.Context, codes []string) ([]models.Country, error) {
	if len(codes) == 0 {
		return nil, nil
	}

	var countries []models.Country
	if err := r.db.NewSelect().
		Model(&countries).
		Where("co.code IN (?)", bun.List(codes)).
		Scan(ctx); err != nil {
		return nil, fmt.Errorf("get countries by code: %w", err)
	}

	return countries, nil
}

func (r *CountryRepository) GetCountryByCode(ctx context.Context, code string) (*models.Country, error) {
	country := new(models.Country)
	err := r.db.NewSelect().
		Model(country).
		Where("co.code = ?", code).
		Limit(1).
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get country by code: %w", err)
	}

	return country, nil
}
