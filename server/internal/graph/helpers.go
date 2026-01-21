package graph

import (
	"strconv"
	"time"

	"github.com/lovely-eye/server/internal/graph/model"
	"github.com/lovely-eye/server/internal/models"
	"github.com/lovely-eye/server/internal/services"
)

func parseDateRangeInput(input *model.DateRangeInput) (time.Time, time.Time) {
	now := time.Now()
	defaultFrom := now.AddDate(0, 0, -30)
	defaultTo := now

	if input == nil {
		return defaultFrom, defaultTo
	}

	from := defaultFrom
	to := defaultTo

	if input.From != nil {
		from = *input.From
	}
	if input.To != nil {
		to = *input.To
	}

	return from, to
}

func parseFilterInput(input *model.FilterInput) (referrer []string, device []string, page []string, country []string) {
	if input == nil {
		return nil, nil, nil, nil
	}
	return input.Referrer, input.Device, input.Page, input.Country
}

func convertToGraphQLEvent(e *models.Event) *model.Event {
	// Convert unix timestamp to time.Time
	createdAt := time.Unix(e.Time, 0)

	// Convert EventData to EventProperty
	properties := make([]*model.EventProperty, 0, len(e.Data))
	for _, data := range e.Data {
		if data.Field != nil {
			properties = append(properties, &model.EventProperty{
				Key:   data.Field.Key,
				Value: data.Value,
			})
		}
	}

	return &model.Event{
		ID:         strconv.FormatInt(e.ID, 10),
		Name:       e.Name,
		Path:       e.Path,
		Properties: properties,
		CreatedAt:  createdAt,
	}
}


func convertToGraphQLEvents(events []*models.Event, total int) *model.EventsResult {
	result := &model.EventsResult{
		Events: make([]*model.Event, 0, len(events)),
		Total:  total,
	}

	for _, e := range events {
		// Convert unix timestamp to time.Time
		createdAt := time.Unix(e.Time, 0)

		// Convert EventData to EventProperty
		properties := make([]*model.EventProperty, 0, len(e.Data))
		for _, data := range e.Data {
			if data.Field != nil {
				properties = append(properties, &model.EventProperty{
					Key:   data.Field.Key,
					Value: data.Value,
				})
			}
		}

		event := &model.Event{
			ID:         strconv.FormatInt(e.ID, 10),
			Name:       e.Name,
			Path:       e.Path,
			Properties: properties,
			CreatedAt:  createdAt,
		}
		result.Events = append(result.Events, event)
	}

	return result
}

func convertToGraphQLEventDefinitions(definitions []*models.EventDefinition) []*model.EventDefinition {
	result := make([]*model.EventDefinition, 0, len(definitions))
	for _, def := range definitions {
		fields := make([]*model.EventDefinitionField, 0, len(def.Fields))
		for _, field := range def.Fields {
			// Convert FieldType enum to string
			var fieldTypeStr string
			switch field.Type {
			case models.FieldTypeString:
				fieldTypeStr = "STRING"
			case models.FieldTypeInt:
				fieldTypeStr = "INT"
			case models.FieldTypeFloat:
				fieldTypeStr = "FLOAT"
			case models.FieldTypeBool:
				fieldTypeStr = "BOOLEAN"
			default:
				fieldTypeStr = "STRING"
			}

			fields = append(fields, &model.EventDefinitionField{
				ID:        strconv.FormatInt(field.ID, 10),
				Key:       field.Key,
				Type:      model.EventFieldType(fieldTypeStr),
				Required:  field.Required,
				MaxLength: field.MaxLength,
			})
		}
		result = append(result, &model.EventDefinition{
			ID:        strconv.FormatInt(def.ID, 10),
			Name:      def.Name,
			Fields:    fields,
			CreatedAt: def.CreatedAt,
			UpdatedAt: def.UpdatedAt,
		})
	}
	return result
}

func convertToGraphQLGeoIPStatus(status services.GeoIPStatus) *model.GeoIPStatus {
	var source *string
	if status.Source != "" {
		source = &status.Source
	}
	var lastError *string
	if status.LastError != "" {
		lastError = &status.LastError
	}
	return &model.GeoIPStatus{
		State:     status.State,
		DbPath:    status.DBPath,
		Source:    source,
		LastError: lastError,
		UpdatedAt: status.UpdatedAt,
	}
}
