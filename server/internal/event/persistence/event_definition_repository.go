package persistence

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/lovely-eye/server/internal/event"
	"github.com/uptrace/bun"
)

type Repository struct {
	db *bun.DB
}

var _ event.Store = (*Repository)(nil)

func New(db *bun.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetBySite(
	ctx context.Context,
	siteID int64,
	limit,
	offset int,
) ([]*event.Definition, error) {
	var defs []*Definition
	q := r.db.NewSelect().
		Model(&defs).
		Where("site_id = ?", siteID).
		Relation("Fields").
		Order("name ASC")
	if limit > 0 {
		q = q.Limit(limit)
	}
	if offset > 0 {
		q = q.Offset(offset)
	}
	err := q.Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get event definitions by site: %w", err)
	}
	result := make([]*event.Definition, 0, len(defs))
	for _, definition := range defs {
		result = append(result, eventDefinition(definition))
	}
	return result, nil
}

func (r *Repository) GetByName(
	ctx context.Context,
	siteID int64,
	name string,
) (*event.Definition, error) {
	def := new(Definition)
	err := r.db.NewSelect().
		Model(def).
		Where("site_id = ?", siteID).
		Where("name = ?", name).
		Relation("Fields").
		Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get event definition by name: %w", err)
	}
	return eventDefinition(def), nil
}

func (r *Repository) Upsert(
	ctx context.Context,
	siteID int64,
	name string,
	fields []*event.Field,
) (*event.Definition, error) {
	var def *Definition
	fieldRows := make([]*Field, 0, len(fields))
	for _, field := range fields {
		fieldRows = append(fieldRows, eventFieldModel(field))
	}
	err := r.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		existing := new(Definition)
		err := tx.NewSelect().
			Model(existing).
			Where("site_id = ?", siteID).
			Where("name = ?", name).
			Scan(ctx)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return fmt.Errorf("select existing event definition: %w", err)
			}
			existing = nil
		}

		if existing == nil {
			newDef := &Definition{
				SiteID:    siteID,
				Name:      name,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			if _, err := tx.NewInsert().Model(newDef).Exec(ctx); err != nil {
				return fmt.Errorf("insert event definition: %w", err)
			}
			def = newDef
		} else {
			existing.UpdatedAt = time.Now()
			if _, err := tx.NewUpdate().Model(existing).Column("updated_at").WherePK().Exec(ctx); err != nil {
				return fmt.Errorf("failed to update event definition: %w", err)
			}
			def = existing
		}

		if _, err := tx.NewDelete().
			Model((*Field)(nil)).
			Where("event_definition_id = ?", def.ID).
			Exec(ctx); err != nil {
			return fmt.Errorf("failed to delete event definition fields: %w", err)
		}

		if len(fieldRows) == 0 {
			return nil
		}

		for _, field := range fieldRows {
			field.EventDefinitionID = def.ID
			field.CreatedAt = time.Now()
			field.UpdatedAt = time.Now()
		}

		if _, err := tx.NewInsert().Model(&fieldRows).Exec(ctx); err != nil {
			return fmt.Errorf("insert event definition fields: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to upsert event definition transaction: %w", err)
	}

	if def == nil {
		return nil, errors.New("failed to upsert event definition")
	}

	return r.GetByName(ctx, siteID, name)
}

func (r *Repository) DeleteByName(ctx context.Context, siteID int64, name string) error {
	err := r.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		def := new(Definition)
		if err := tx.NewSelect().
			Model(def).
			Where("site_id = ?", siteID).
			Where("name = ?", name).
			Scan(ctx); err != nil {
			return fmt.Errorf("select event definition: %w", err)
		}

		if _, err := tx.NewDelete().
			Model((*Field)(nil)).
			Where("event_definition_id = ?", def.ID).
			Exec(ctx); err != nil {
			return fmt.Errorf("delete event definition fields: %w", err)
		}

		if _, err := tx.NewDelete().
			Model((*Definition)(nil)).
			Where("id = ?", def.ID).
			Exec(ctx); err != nil {
			return fmt.Errorf("delete event definition: %w", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("delete event definition transaction: %w", err)
	}
	return nil
}

func eventDefinition(row *Definition) *event.Definition {
	if row == nil {
		return nil
	}
	definition := &event.Definition{
		ID:        row.ID,
		SiteID:    row.SiteID,
		Name:      row.Name,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}
	for _, field := range row.Fields {
		definition.Fields = append(definition.Fields, eventField(field))
	}
	return definition
}

func eventField(row *Field) *event.Field {
	if row == nil {
		return nil
	}
	return &event.Field{
		ID:                row.ID,
		EventDefinitionID: row.EventDefinitionID,
		Key:               row.Key,
		Type:              event.FieldType(row.Type),
		Required:          row.Required,
		MaxLength:         row.MaxLength,
		CreatedAt:         row.CreatedAt,
		UpdatedAt:         row.UpdatedAt,
	}
}

func eventFieldModel(field *event.Field) *Field {
	return &Field{
		ID:                field.ID,
		EventDefinitionID: field.EventDefinitionID,
		Key:               field.Key,
		Type:              FieldType(field.Type),
		Required:          field.Required,
		MaxLength:         field.MaxLength,
		CreatedAt:         field.CreatedAt,
		UpdatedAt:         field.UpdatedAt,
	}
}
