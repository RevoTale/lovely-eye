package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/lovely-eye/server/internal/models"
	"github.com/uptrace/bun"
)

type EventDefinitionRepository struct {
	db *bun.DB
}

func NewEventDefinitionRepository(db *bun.DB) *EventDefinitionRepository {
	return &EventDefinitionRepository{db: db}
}

func (r *EventDefinitionRepository) GetBySite(ctx context.Context, siteID int64) ([]*models.EventDefinition, error) {
	var defs []*models.EventDefinition
	err := r.db.NewSelect().
		Model(&defs).
		Where("site_id = ?", siteID).
		Relation("Fields").
		Order("name ASC").
		Scan(ctx)
	return defs, err
}

func (r *EventDefinitionRepository) GetByName(ctx context.Context, siteID int64, name string) (*models.EventDefinition, error) {
	def := new(models.EventDefinition)
	err := r.db.NewSelect().
		Model(def).
		Where("site_id = ?", siteID).
		Where("name = ?", name).
		Relation("Fields").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return def, nil
}

func (r *EventDefinitionRepository) Upsert(ctx context.Context, siteID int64, name string, fields []*models.EventDefinitionField) (*models.EventDefinition, error) {
	var def *models.EventDefinition
	err := r.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		existing := new(models.EventDefinition)
		err := tx.NewSelect().
			Model(existing).
			Where("site_id = ?", siteID).
			Where("name = ?", name).
			Scan(ctx)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return err
			}
			existing = nil
		}

		if existing == nil {
			newDef := &models.EventDefinition{
				SiteID:    siteID,
				Name:      name,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			if _, err := tx.NewInsert().Model(newDef).Exec(ctx); err != nil {
				return err
			}
			def = newDef
		} else {
			existing.UpdatedAt = time.Now()
			if _, err := tx.NewUpdate().Model(existing).Column("updated_at").WherePK().Exec(ctx); err != nil {
				return err
			}
			def = existing
		}

		if _, err := tx.NewDelete().
			Model((*models.EventDefinitionField)(nil)).
			Where("event_definition_id = ?", def.ID).
			Exec(ctx); err != nil {
			return err
		}

		if len(fields) == 0 {
			return nil
		}

		for _, field := range fields {
			field.EventDefinitionID = def.ID
			field.CreatedAt = time.Now()
			field.UpdatedAt = time.Now()
		}

		if _, err := tx.NewInsert().Model(&fields).Exec(ctx); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	if def == nil {
		return nil, errors.New("failed to upsert event definition")
	}

	return r.GetByName(ctx, siteID, name)
}

func (r *EventDefinitionRepository) DeleteByName(ctx context.Context, siteID int64, name string) error {
	def := new(models.EventDefinition)
	err := r.db.NewSelect().
		Model(def).
		Where("site_id = ?", siteID).
		Where("name = ?", name).
		Scan(ctx)
	if err != nil {
		return err
	}

	_, err = r.db.NewDelete().
		Model((*models.EventDefinition)(nil)).
		Where("id = ?", def.ID).
		Exec(ctx)
	return err
}
