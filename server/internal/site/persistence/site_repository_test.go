package persistence

import (
	"context"
	"errors"
	"testing"
	"time"

	analyticspersistence "github.com/lovely-eye/server/internal/analytics/persistence"
	"github.com/lovely-eye/server/internal/event"
	eventpersistence "github.com/lovely-eye/server/internal/event/persistence"
	sitefeature "github.com/lovely-eye/server/internal/site"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
)

func TestRepository_DeleteRemovesOwnedData(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	siteRepo := New(db)
	eventDefinitionRepo := eventpersistence.New(db)
	site := createTestSite(t, db)
	ctx := context.Background()

	eventTime := time.Now()
	clientID := createTestClient(t, db, site.ID, "delete-site-client", "desktop", "Chrome", "Linux")
	sessionID := insertSessionWithPath(t, db, site.ID, clientID, "/delete", eventTime, 10, 1)
	definition, err := eventDefinitionRepo.Upsert(ctx, site.ID, "deleted_event", []*event.Field{
		{Key: "reason", Type: event.FieldTypeString, MaxLength: 100},
	})
	require.NoError(t, err)

	event := &analyticspersistence.Event{
		SessionID:    sessionID,
		Time:         eventTime.Unix(),
		Hour:         eventTime.Unix() / 3600,
		Day:          eventTime.Unix() / 86400,
		Path:         "/delete",
		DefinitionID: &definition.ID,
	}
	_, err = db.NewInsert().Model(event).Exec(ctx)
	require.NoError(t, err)
	_, err = db.NewInsert().Model(&analyticspersistence.EventData{
		EventID: event.ID,
		FieldID: definition.Fields[0].ID,
		Value:   "test",
	}).Exec(ctx)
	require.NoError(t, err)
	_, err = db.NewInsert().Model(&BlockedIP{SiteID: site.ID, IP: "203.0.113.10"}).Exec(ctx)
	require.NoError(t, err)
	_, err = db.NewInsert().Model(&BlockedCountry{SiteID: site.ID, CountryCode: "US"}).Exec(ctx)
	require.NoError(t, err)

	require.NoError(t, siteRepo.Delete(ctx, site.ID))

	_, err = siteRepo.GetByID(ctx, site.ID)
	require.True(t, errors.Is(err, sitefeature.ErrSiteNotFound))
	requireModelTableEmpty(t, db, (*analyticspersistence.EventData)(nil))
	requireModelTableEmpty(t, db, (*analyticspersistence.Event)(nil))
	requireModelTableEmpty(t, db, (*eventpersistence.Field)(nil))
	requireModelTableEmpty(t, db, (*eventpersistence.Definition)(nil))
	requireModelTableEmpty(t, db, (*analyticspersistence.Session)(nil))
	requireModelTableEmpty(t, db, (*analyticspersistence.Client)(nil))
	requireModelTableEmpty(t, db, (*BlockedIP)(nil))
	requireModelTableEmpty(t, db, (*BlockedCountry)(nil))
	requireModelTableEmpty(t, db, (*Domain)(nil))
}

func TestRepository_GetByPublicKeyReturnsAnalyticsRelations(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	repository := New(db)
	site := createTestSite(t, db)
	ctx := context.Background()

	_, err := db.NewInsert().Model(&Domain{SiteID: site.ID, Domain: "alias.example.com", Position: 1}).Exec(ctx)
	require.NoError(t, err)
	_, err = db.NewInsert().Model(&BlockedIP{SiteID: site.ID, IP: "203.0.113.10"}).Exec(ctx)
	require.NoError(t, err)
	_, err = db.NewInsert().Model(&BlockedCountry{SiteID: site.ID, CountryCode: "US"}).Exec(ctx)
	require.NoError(t, err)

	loaded, err := repository.GetByPublicKey(ctx, site.PublicKey)
	require.NoError(t, err)
	require.Equal(t, site.ID, loaded.ID)
	require.Equal(t, []string{"example.com", "alias.example.com"}, []string{
		loaded.Domains[0].Domain,
		loaded.Domains[1].Domain,
	})
	require.Equal(t, "203.0.113.10", loaded.BlockedIPs[0].IP)
	require.Equal(t, "US", loaded.BlockedCountries[0].CountryCode)
}

func requireModelTableEmpty(t *testing.T, db *bun.DB, model any) {
	t.Helper()
	count, err := db.NewSelect().Model(model).Count(context.Background())
	require.NoError(t, err)
	require.Zero(t, count)
}
