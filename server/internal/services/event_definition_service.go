package services

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/lovely-eye/server/internal/models"
	"github.com/lovely-eye/server/internal/repository"
)

const (
	defaultEventMaxLength = 500
	maxEventNameLength    = 100
	maxEventKeyLength     = 100
)

var (
	ErrInvalidEventName  = errors.New("invalid event name")
	ErrInvalidFieldKey   = errors.New("invalid field key")
	ErrInvalidFieldType  = errors.New("invalid field type")
	ErrInvalidFieldLimit = errors.New("invalid field max length")
)

type EventDefinitionService struct {
	repo *repository.EventDefinitionRepository
}

func NewEventDefinitionService(repo *repository.EventDefinitionRepository) *EventDefinitionService {
	return &EventDefinitionService{repo: repo}
}

type EventFieldInput struct {
	Key       string
	Type      string
	Required  bool
	MaxLength *int
}

type EventDefinitionInput struct {
	Name   string
	Fields []EventFieldInput
}

func (s *EventDefinitionService) List(ctx context.Context, siteID int64) ([]*models.EventDefinition, error) {
	defs, err := s.repo.GetBySite(ctx, siteID)
	if err != nil {
		return nil, fmt.Errorf("failed to list event definitions: %w", err)
	}
	return defs, nil
}

func (s *EventDefinitionService) Upsert(ctx context.Context, siteID int64, input EventDefinitionInput) (*models.EventDefinition, error) {
	name := strings.TrimSpace(input.Name)
	if name == "" || len(name) > maxEventNameLength {
		return nil, ErrInvalidEventName
	}

	fields := make([]*models.EventDefinitionField, 0, len(input.Fields))
	seen := make(map[string]struct{}, len(input.Fields))
	for _, field := range input.Fields {
		key := strings.TrimSpace(field.Key)
		if key == "" || len(key) > maxEventKeyLength {
			return nil, ErrInvalidFieldKey
		}
		if _, ok := seen[key]; ok {
			return nil, ErrInvalidFieldKey
		}
		seen[key] = struct{}{}

		fieldTypeStr := strings.ToLower(strings.TrimSpace(field.Type))
		if fieldTypeStr == "" {
			fieldTypeStr = "string"
		}

	var fieldType models.FieldType
	switch fieldTypeStr {
	case "string":
		fieldType = models.FieldTypeString
	case "int", "integer":
		fieldType = models.FieldTypeInt
	case "bool", "boolean":
		fieldType = models.FieldTypeBool
		default:
			return nil, ErrInvalidFieldType
		}

		maxLen := defaultEventMaxLength
		if field.MaxLength != nil {
			if *field.MaxLength <= 0 {
				return nil, ErrInvalidFieldLimit
			}
			maxLen = *field.MaxLength
		}

		fields = append(fields, &models.EventDefinitionField{
			Key:       key,
			Type:      fieldType,
			Required:  field.Required,
			MaxLength: maxLen,
		})
	}

	def, err := s.repo.Upsert(ctx, siteID, name, fields)
	if err != nil {
		return nil, fmt.Errorf("failed to upsert event definition: %w", err)
	}
	return def, nil
}

func (s *EventDefinitionService) Delete(ctx context.Context, siteID int64, name string) error {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return ErrInvalidEventName
	}
	if err := s.repo.DeleteByName(ctx, siteID, trimmed); err != nil {
		return fmt.Errorf("failed to delete event definition: %w", err)
	}
	return nil
}
