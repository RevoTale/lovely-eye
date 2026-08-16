package graph

import (
	"strconv"
	"strings"
	"time"

	"github.com/lovely-eye/server/internal/analytics"
	"github.com/lovely-eye/server/internal/auth"
	"github.com/lovely-eye/server/internal/event"
	"github.com/lovely-eye/server/internal/geoip"
	"github.com/lovely-eye/server/internal/graph/model"
)

func convertToGraphQLUser(user *auth.User) *model.User {
	return &model.User{
		ID:        strconv.FormatInt(user.ID, 10),
		Username:  user.Username,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	}
}

func parseDateRangeInput(input *model.DateRangeInput, maxRangeDays int) (time.Time, time.Time, error) {
	now := time.Now()
	defaultFrom := now.AddDate(0, 0, -30)
	defaultTo := now

	if input == nil {
		return defaultFrom, defaultTo, nil
	}

	from := defaultFrom
	to := defaultTo

	if input.From != nil {
		from = *input.From
	}
	if input.To != nil {
		to = *input.To
	}

	if err := validateDateRange(from, to, maxRangeDays); err != nil {
		return time.Time{}, time.Time{}, err
	}
	return from, to, nil
}

func validateDateRange(from, to time.Time, maxRangeDays int) error {
	if from.After(to) {
		return badUserInput("date range from must be before to")
	}
	if maxRangeDays > 0 && to.Sub(from) > time.Duration(maxRangeDays)*24*time.Hour {
		return badUserInputf("date range exceeds %d days", maxRangeDays)
	}
	return nil
}

func parseFilterInput(input *model.FilterInput, limits DashboardLimits) (analytics.Filter, error) {
	if input == nil {
		return analytics.Filter{}, nil
	}

	if err := validateStringFilters(limits, input.Referrer, input.Browser, input.Device, input.Os, input.Page, input.Country, input.EventName, input.EventPath, input.EventDefinitionID); err != nil {
		return analytics.Filter{}, err
	}
	if limits.MaxFilterValues > 0 && len(input.EventType) > limits.MaxFilterValues {
		return analytics.Filter{}, badUserInputf("filter eventType exceeds %d values", limits.MaxFilterValues)
	}
	eventDefinitionIDs, err := parseEventDefinitionIDs(input.EventDefinitionID)
	if err != nil {
		return analytics.Filter{}, err
	}

	referrers := make([]string, 0, len(input.Referrer))
	for _, referrer := range input.Referrer {
		if referrer == "(direct)" {
			referrer = ""
		}
		referrers = append(referrers, referrer)
	}

	return analytics.Filter{
		Referrer:           referrers,
		Browser:            input.Browser,
		Device:             input.Device,
		OS:                 input.Os,
		Page:               input.Page,
		Country:            input.Country,
		EventTypes:         parseEventTypes(input.EventType),
		EventName:          input.EventName,
		EventPath:          input.EventPath,
		EventDefinitionIDs: eventDefinitionIDs,
	}, nil
}

func validateStringFilters(limits DashboardLimits, groups ...[]string) error {
	for _, values := range groups {
		if limits.MaxFilterValues > 0 && len(values) > limits.MaxFilterValues {
			return badUserInputf("filter exceeds %d values", limits.MaxFilterValues)
		}
		if limits.MaxFilterStringLength <= 0 {
			continue
		}
		for _, value := range values {
			if len(value) > limits.MaxFilterStringLength {
				return badUserInputf("filter value exceeds %d bytes", limits.MaxFilterStringLength)
			}
		}
	}
	return nil
}

func isFilterEmpty(filter analytics.Filter) bool {
	return len(filter.Referrer) == 0 &&
		len(filter.Browser) == 0 &&
		len(filter.Device) == 0 &&
		len(filter.OS) == 0 &&
		len(filter.Page) == 0 &&
		len(filter.Country) == 0 &&
		len(filter.EventTypes) == 0 &&
		len(filter.EventName) == 0 &&
		len(filter.EventPath) == 0 &&
		len(filter.EventDefinitionIDs) == 0
}

func convertToGraphQLEvent(e *analytics.Event) *model.Event {

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
		Definition: convertAnalyticsEventDefinition(e.Definition),
		Properties: properties,
		CreatedAt:  createdAt,
	}
}

func convertToGraphQLEvents(events []*analytics.Event, total int) *model.EventsResult {
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
			Definition: convertAnalyticsEventDefinition(e.Definition),
			Properties: properties,
			CreatedAt:  createdAt,
		}
		result.Events = append(result.Events, event)
	}

	return result
}

func convertAnalyticsEventDefinition(definition *analytics.EventDefinition) *model.EventDefinition {
	if definition == nil {
		return nil
	}
	fields := make([]*model.EventDefinitionField, 0, len(definition.Fields))
	for _, field := range definition.Fields {
		if field == nil {
			continue
		}
		fieldType := model.EventFieldTypeString
		switch field.Type {
		case analytics.EventFieldTypeInt:
			fieldType = model.EventFieldTypeInt
		case analytics.EventFieldTypeBool:
			fieldType = model.EventFieldTypeBoolean
		}
		fields = append(fields, &model.EventDefinitionField{
			ID:        strconv.FormatInt(field.ID, 10),
			Key:       field.Key,
			Type:      fieldType,
			Required:  field.Required,
			MaxLength: field.MaxLength,
		})
	}
	return &model.EventDefinition{
		ID:        strconv.FormatInt(definition.ID, 10),
		Name:      definition.Name,
		Fields:    fields,
		CreatedAt: definition.CreatedAt,
		UpdatedAt: definition.UpdatedAt,
	}
}

func convertToGraphQLEventDefinition(def *event.Definition) *model.EventDefinition {
	if def == nil {
		return nil
	}
	fields := make([]*model.EventDefinitionField, 0, len(def.Fields))
	for _, field := range def.Fields {
		var fieldTypeStr string
		switch field.Type {
		case event.FieldTypeString:
			fieldTypeStr = "STRING"
		case event.FieldTypeInt:
			fieldTypeStr = "INT"
		case event.FieldTypeBool:
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

func convertToGraphQLEventDefinitions(definitions []*event.Definition) []*model.EventDefinition {
	result := make([]*model.EventDefinition, 0, len(definitions))
	for _, def := range definitions {
		result = append(result, convertToGraphQLEventDefinition(def))
	}
	return result
}

func convertToGraphQLGeoIPStatus(status geoip.Status) *model.GeoIPStatus {
	var source *string
	if status.Source != "" {
		value := status.Source.String()
		source = &value
	}
	var lastError *string
	if status.LastError != "" {
		lastError = &status.LastError
	}
	return &model.GeoIPStatus{
		State:     graphQLGeoIPState(status.State),
		DbPath:    status.DBPath,
		Source:    source,
		LastError: lastError,
		UpdatedAt: status.UpdatedAt,
	}
}

func graphQLGeoIPState(state geoip.State) model.GeoIPState {
	switch state {
	case geoip.StateMissing:
		return model.GeoIPStateMissing
	case geoip.StateDownloading:
		return model.GeoIPStateDownloading
	case geoip.StateReady:
		return model.GeoIPStateReady
	case geoip.StateError:
		return model.GeoIPStateError
	default:
		return model.GeoIPStateDisabled
	}
}

func newGraphQLCountry(code string, name string) *model.Country {
	graphQLCountry := &model.Country{Code: code}
	trimmedName := strings.TrimSpace(name)
	if trimmedName != "" {
		graphQLCountry.NameCache = &trimmedName
	}
	return graphQLCountry
}

func parseEventDefinitionIDs(values []string) ([]int64, error) {
	if len(values) == 0 {
		return nil, nil
	}
	ids := make([]int64, 0, len(values))
	for _, value := range values {
		id, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return nil, badUserInput("invalid event definition ID")
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func parseEventTypes(values []model.EventType) []analytics.EventType {
	if len(values) == 0 {
		return nil
	}
	types := make([]analytics.EventType, 0, len(values))
	for _, value := range values {
		switch value {
		case model.EventTypePageView:
			types = append(types, analytics.EventTypePageView)
		case model.EventTypePredefined:
			types = append(types, analytics.EventTypePredefined)
		}
	}
	return types
}
