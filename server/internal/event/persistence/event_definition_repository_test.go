package persistence

import (
	"context"
	"testing"

	"github.com/lovely-eye/server/internal/event"
	"github.com/stretchr/testify/require"
)

func TestRepository_DeleteByNameDeletesDefinitionFields(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	repo := New(db)
	site := createTestSite(t, db)
	ctx := context.Background()

	definition, err := repo.Upsert(ctx, site.ID, "signup_completed", []*event.Field{
		{
			Key:       "plan",
			Type:      event.FieldTypeString,
			Required:  true,
			MaxLength: 100,
		},
	})
	require.NoError(t, err)
	require.Len(t, definition.Fields, 1)

	require.NoError(t, repo.DeleteByName(ctx, site.ID, definition.Name))

	definitions, err := repo.GetBySite(ctx, site.ID, 100, 0)
	require.NoError(t, err)
	require.Empty(t, definitions)

	fieldCount, err := db.NewSelect().
		Model((*Field)(nil)).
		Where("event_definition_id = ?", definition.ID).
		Count(ctx)
	require.NoError(t, err)
	require.Zero(t, fieldCount)
}
