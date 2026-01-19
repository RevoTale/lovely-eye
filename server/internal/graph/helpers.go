package graph

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
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


func convertToGraphQLEvents(events []*models.Event, total int) *model.EventsResult {
	result := &model.EventsResult{
		Events: make([]*model.Event, 0, len(events)),
		Total:  total,
	}

	for _, e := range events {
		event := &model.Event{
			ID:         strconv.FormatInt(e.ID, 10),
			Name:       e.Name,
			Path:       e.Path,
			Properties: parseEventProperties(e.Properties),
			CreatedAt:  e.CreatedAt,
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
			fields = append(fields, &model.EventDefinitionField{
				ID:        strconv.FormatInt(field.ID, 10),
				Key:       field.Key,
				Type:      model.EventFieldType(strings.ToUpper(field.Type)),
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

func parseEventProperties(propsJSON string) []*model.EventProperty {
	if propsJSON == "" {
		return []*model.EventProperty{}
	}

	var props map[string]interface{}
	if err := json.Unmarshal([]byte(propsJSON), &props); err != nil {
		return []*model.EventProperty{}
	}

	result := make([]*model.EventProperty, 0, len(props))
	for k, v := range props {
		result = append(result, &model.EventProperty{
			Key:   k,
			Value: fmt.Sprint(v),
		})
	}

	return result
}
