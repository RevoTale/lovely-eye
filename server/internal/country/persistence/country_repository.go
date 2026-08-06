package persistence

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/lovely-eye/server/internal/country"
	"github.com/uptrace/bun"
)

type Repository struct {
	db *bun.DB
}

var _ country.Store = (*Repository)(nil)

func New(db *bun.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) UpsertCountries(ctx context.Context, countries []country.Info) error {
	if len(countries) == 0 {
		return nil
	}

	rows := make([]Country, 0, len(countries))
	for _, value := range countries {
		rows = append(rows, Country{Code: value.Code, Name: value.Name})
	}
	_, err := r.db.NewInsert().
		Model(&rows).
		On("CONFLICT (code) DO UPDATE").
		Set("name = EXCLUDED.name").
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("upsert countries: %w", err)
	}

	return nil
}

func (r *Repository) SearchCountries(
	ctx context.Context,
	search string,
	limit,
	offset int,
) ([]country.Info, error) {
	var countries []Country

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

	return countryInfo(countries), nil
}

func (r *Repository) GetCountriesByCodes(ctx context.Context, codes []string) ([]country.Info, error) {
	if len(codes) == 0 {
		return nil, nil
	}

	var countries []Country
	if err := r.db.NewSelect().
		Model(&countries).
		Where("co.code IN (?)", bun.List(codes)).
		Scan(ctx); err != nil {
		return nil, fmt.Errorf("get countries by code: %w", err)
	}

	return countryInfo(countries), nil
}

func (r *Repository) GetCountryByCode(ctx context.Context, code string) (*country.Info, error) {
	row := new(Country)
	err := r.db.NewSelect().
		Model(row).
		Where("co.code = ?", code).
		Limit(1).
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, country.ErrNotFound
		}
		return nil, fmt.Errorf("get country by code: %w", err)
	}

	return &country.Info{Code: row.Code, Name: row.Name}, nil
}

func countryInfo(rows []Country) []country.Info {
	result := make([]country.Info, 0, len(rows))
	for _, row := range rows {
		result = append(result, country.Info{Code: row.Code, Name: row.Name})
	}
	return result
}
