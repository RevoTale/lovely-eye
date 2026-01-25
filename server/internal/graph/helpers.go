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

func parseFilterInput(input *model.FilterInput) services.DashboardFilter {
	if input == nil {
		return services.DashboardFilter{}
	}
	return services.DashboardFilter{
		Referrer:           input.Referrer,
		Device:             input.Device,
		Page:               input.Page,
		Country:            input.Country,
		EventName:          input.EventName,
		EventPath:          input.EventPath,
		EventDefinitionIDs: parseEventDefinitionIDs(input.EventDefinitionID),
	}
}

func isFilterEmpty(filter services.DashboardFilter) bool {
	return len(filter.Referrer) == 0 &&
		len(filter.Device) == 0 &&
		len(filter.Page) == 0 &&
		len(filter.Country) == 0 &&
		len(filter.EventName) == 0 &&
		len(filter.EventPath) == 0 &&
		len(filter.EventDefinitionIDs) == 0
}

func convertToGraphQLEvent(e *models.Event) *model.Event {

	createdAt := time.Unix(e.Time, 0)

	name := e.Path
	if e.Definition != nil {
		name = e.Definition.Name
	}

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
		Name:       name,
		Path:       e.Path,
		Definition: convertToGraphQLEventDefinition(e.Definition),
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

		createdAt := time.Unix(e.Time, 0)

		name := e.Path
		if e.Definition != nil {
			name = e.Definition.Name
		}

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
			Name:       name,
			Path:       e.Path,
			Definition: convertToGraphQLEventDefinition(e.Definition),
			Properties: properties,
			CreatedAt:  createdAt,
		}
		result.Events = append(result.Events, event)
	}

	return result
}

func convertToGraphQLEventDefinition(def *models.EventDefinition) *model.EventDefinition {
	if def == nil {
		return nil
	}
	fields := make([]*model.EventDefinitionField, 0, len(def.Fields))
	for _, field := range def.Fields {
		var fieldTypeStr string
		switch field.Type {
		case models.FieldTypeString:
			fieldTypeStr = "STRING"
		case models.FieldTypeInt:
			fieldTypeStr = "INT"
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
	return &model.EventDefinition{
		ID:        strconv.FormatInt(def.ID, 10),
		Name:      def.Name,
		Fields:    fields,
		CreatedAt: def.CreatedAt,
		UpdatedAt: def.UpdatedAt,
	}
}

func convertToGraphQLEventDefinitions(definitions []*models.EventDefinition) []*model.EventDefinition {
	result := make([]*model.EventDefinition, 0, len(definitions))
	for _, def := range definitions {
		result = append(result, convertToGraphQLEventDefinition(def))
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

func parseEventDefinitionIDs(values []string) []int64 {
	if len(values) == 0 {
		return nil
	}
	ids := make([]int64, 0, len(values))
	for _, value := range values {
		id, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			continue
		}
		ids = append(ids, id)
	}
	return ids
}
