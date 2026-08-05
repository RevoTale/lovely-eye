package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/lovely-eye/server/internal/models"
	"github.com/lovely-eye/server/internal/repository"
	"github.com/stretchr/testify/require"
)

func TestSanitizeEventPropertiesTruncatesAndStrips(t *testing.T) {
	fields := []*models.EventDefinitionField{
		{
			Key:       "error",
			Type:      models.FieldTypeString,
			Required:  true,
			MaxLength: 5,
		},
	}

	props := `{"error":"toolong","extra":"drop"}`
	sanitized, ok, err := sanitizeEventProperties(props, fields)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatalf("expected ok=true")
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(sanitized), &parsed); err != nil {
		t.Fatalf("failed to unmarshal sanitized props: %v", err)
	}

	if parsed["error"] != "toolo" {
		t.Fatalf("expected truncated value, got %v", parsed["error"])
	}
	if _, ok := parsed["extra"]; ok {
		t.Fatalf("expected extra key to be stripped")
	}
}

func TestSanitizeEventPropertiesMissingRequired(t *testing.T) {
	fields := []*models.EventDefinitionField{
		{
			Key:      "code",
			Type:     models.FieldTypeString,
			Required: true,
		},
	}

	sanitized, ok, err := sanitizeEventProperties(`{}`, fields)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Fatalf("expected ok=false for missing required field, got %v", sanitized)
	}
}

func TestSanitizeEventPropertiesTypeMismatch(t *testing.T) {
	fields := []*models.EventDefinitionField{
		{
			Key:  "count",
			Type: models.FieldTypeFloat,
		},
	}

	sanitized, ok, err := sanitizeEventProperties(`{"count":"nope"}`, fields)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Fatalf("expected ok=false for type mismatch, got %v", sanitized)
	}
}

func TestSanitizeEventPropertiesBoolean(t *testing.T) {
	fields := []*models.EventDefinitionField{
		{
			Key:  "retry",
			Type: models.FieldTypeBool,
		},
	}

	sanitized, ok, err := sanitizeEventProperties(`{"retry":true}`, fields)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatalf("expected ok=true")
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(sanitized), &parsed); err != nil {
		t.Fatalf("failed to unmarshal sanitized props: %v", err)
	}
	if parsed["retry"] != true {
		t.Fatalf("expected boolean value, got %v", parsed["retry"])
	}
}

func TestCollectEventStoresNonStringProperties(t *testing.T) {
	ctx := context.Background()
	db := setupAnalyticsServiceTestDB(t)
	site := createAnalyticsIdentitySite(t, db)
	eventDefinitionRepo := repository.NewEventDefinitionRepository(db)
	_, err := eventDefinitionRepo.Upsert(ctx, site.ID, "checkout_failed", []*models.EventDefinitionField{
		{
			Key:       "count",
			Type:      models.FieldTypeInt,
			MaxLength: defaultEventMaxLength,
		},
		{
			Key:       "retry",
			Type:      models.FieldTypeBool,
			MaxLength: defaultEventMaxLength,
		},
	})
	require.NoError(t, err)

	service := NewAnalyticsService(
		repository.NewAnalyticsRepository(db),
		repository.NewSiteRepository(db),
		eventDefinitionRepo,
		nil,
		nil,
		testAnalyticsIdentitySecret,
	)

	err = service.CollectEvent(ctx, EventInput{
		SiteKey:    site.PublicKey,
		Name:       "checkout_failed",
		Path:       "/checkout",
		Properties: `{"count":42,"retry":true}`,
		UserAgent:  "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/120.0",
		IP:         "203.0.113.42",
		Origin:     "https://identity.test",
	})
	require.NoError(t, err)

	var rows []struct {
		Key   string
		Value string
	}
	err = db.NewSelect().
		TableExpr("event_data evd").
		Join("INNER JOIN event_definition_fields edf ON evd.field_id = edf.id").
		ColumnExpr("edf.key").
		ColumnExpr("evd.value").
		Order("edf.key ASC").
		Scan(ctx, &rows)
	require.NoError(t, err)
	require.Equal(t, []struct {
		Key   string
		Value string
	}{
		{Key: "count", Value: "42"},
		{Key: "retry", Value: "true"},
	}, rows)
}
