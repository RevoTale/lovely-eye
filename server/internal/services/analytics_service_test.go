package services

import (
	"encoding/json"
	"testing"

	"github.com/lovely-eye/server/internal/models"
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
