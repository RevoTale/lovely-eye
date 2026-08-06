package event

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
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

type FieldType int8

const (
	FieldTypeString FieldType = 0
	FieldTypeInt    FieldType = 1
	FieldTypeFloat  FieldType = 2
	FieldTypeBool   FieldType = 3
)

type Definition struct {
	ID        int64
	SiteID    int64
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	Fields    []*Field
}

type Field struct {
	ID                int64
	EventDefinitionID int64
	Key               string
	Type              FieldType
	Required          bool
	MaxLength         int
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type Store interface {
	GetBySite(ctx context.Context, siteID int64, limit, offset int) ([]*Definition, error)
	GetByName(ctx context.Context, siteID int64, name string) (*Definition, error)
	Upsert(ctx context.Context, siteID int64, name string, fields []*Field) (*Definition, error)
	DeleteByName(ctx context.Context, siteID int64, name string) error
}

type Service struct {
	store Store
}

func NewService(store Store) *Service {
	return &Service{store: store}
}

type FieldInput struct {
	Key       string
	Type      string
	Required  bool
	MaxLength *int
}

type DefinitionInput struct {
	Name   string
	Fields []FieldInput
}

func (s *Service) List(ctx context.Context, siteID int64, limit, offset int) ([]*Definition, error) {
	defs, err := s.store.GetBySite(ctx, siteID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list event definitions: %w", err)
	}
	return defs, nil
}

func (s *Service) Upsert(ctx context.Context, siteID int64, input DefinitionInput) (*Definition, error) {
	name := strings.TrimSpace(input.Name)
	if name == "" || len(name) > maxEventNameLength {
		return nil, ErrInvalidEventName
	}

	fields := make([]*Field, 0, len(input.Fields))
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

		var fieldType FieldType
		switch fieldTypeStr {
		case "string":
			fieldType = FieldTypeString
		case "int", "integer":
			fieldType = FieldTypeInt
		case "bool", "boolean":
			fieldType = FieldTypeBool
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

		fields = append(fields, &Field{
			Key:       key,
			Type:      fieldType,
			Required:  field.Required,
			MaxLength: maxLen,
		})
	}

	def, err := s.store.Upsert(ctx, siteID, name, fields)
	if err != nil {
		return nil, fmt.Errorf("failed to upsert event definition: %w", err)
	}
	return def, nil
}

func (s *Service) Delete(ctx context.Context, siteID int64, name string) error {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return ErrInvalidEventName
	}
	if err := s.store.DeleteByName(ctx, siteID, trimmed); err != nil {
		return fmt.Errorf("failed to delete event definition: %w", err)
	}
	return nil
}
