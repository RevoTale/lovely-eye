package analytics

import (
	analyticspersistence "github.com/lovely-eye/server/internal/analytics/persistence"
	eventpersistence "github.com/lovely-eye/server/internal/event/persistence"
)

func repositoryAnalyticsQuery(query Query) analyticspersistence.AnalyticsQuery {
	eventTypes := make([]analyticspersistence.EventType, 0, len(query.Filter.EventTypes))
	for _, eventType := range query.Filter.EventTypes {
		eventTypes = append(eventTypes, analyticspersistence.EventType(eventType))
	}
	return analyticspersistence.AnalyticsQuery{
		SiteID: query.SiteID,
		From:   query.From,
		To:     query.To,
		Limit:  query.Limit,
		Offset: query.Offset,
		Bucket: analyticspersistence.TimeBucket(query.Bucket),
		Filter: analyticspersistence.AnalyticsFilter{
			Referrer:           query.Filter.Referrer,
			Browser:            query.Filter.Browser,
			Device:             query.Filter.Device,
			OS:                 query.Filter.OS,
			Page:               query.Filter.Page,
			PagePathContains:   query.Filter.PagePathContains,
			Country:            query.Filter.Country,
			EventTypes:         eventTypes,
			EventName:          query.Filter.EventName,
			EventPath:          query.Filter.EventPath,
			EventDefinitionIDs: query.Filter.EventDefinitionIDs,
		},
	}
}

func pageStats(values []analyticspersistence.PageStats) []PageStats {
	result := make([]PageStats, 0, len(values))
	for _, value := range values {
		result = append(result, PageStats{
			Path: value.Path, Views: value.Views, Visitors: value.Visitors,
		})
	}
	return result
}

func referrerStats(values []analyticspersistence.ReferrerStats) []ReferrerStats {
	result := make([]ReferrerStats, 0, len(values))
	for _, value := range values {
		result = append(result, ReferrerStats{
			Referrer: value.Referrer, Visitors: value.Visitors,
		})
	}
	return result
}

func browserStats(values []analyticspersistence.BrowserStats) []BrowserStats {
	result := make([]BrowserStats, 0, len(values))
	for _, value := range values {
		result = append(result, BrowserStats{
			Browser: value.Browser.String(), Visitors: value.Visitors,
		})
	}
	return result
}

func deviceStats(values []analyticspersistence.DeviceStats) []DeviceStats {
	result := make([]DeviceStats, 0, len(values))
	for _, value := range values {
		result = append(result, DeviceStats{
			Device: value.Device.String(), Visitors: value.Visitors,
		})
	}
	return result
}

func operatingSystemStats(values []analyticspersistence.OperatingSystemStats) []OperatingSystemStats {
	result := make([]OperatingSystemStats, 0, len(values))
	for _, value := range values {
		result = append(result, OperatingSystemStats{
			OS: value.OS.String(), Visitors: value.Visitors,
		})
	}
	return result
}

func countryStats(values []analyticspersistence.CountryStats) []CountryStats {
	result := make([]CountryStats, 0, len(values))
	for _, value := range values {
		result = append(result, CountryStats{
			CountryCode: value.CountryCode, Visitors: value.Visitors,
		})
	}
	return result
}

func timeSeriesStats(values []analyticspersistence.DailyVisitorStats) []TimeSeriesStats {
	result := make([]TimeSeriesStats, 0, len(values))
	for _, value := range values {
		result = append(result, TimeSeriesStats{
			DateBucket: value.DateBucket,
			Visitors:   value.Visitors,
			PageViews:  value.PageViews,
			Sessions:   value.Sessions,
		})
	}
	return result
}

func activePageStats(values []analyticspersistence.ActivePageStats) []ActivePageStats {
	result := make([]ActivePageStats, 0, len(values))
	for _, value := range values {
		result = append(result, ActivePageStats{
			Path: value.Path, Visitors: value.Visitors,
		})
	}
	return result
}

func analyticsEvents(values []*analyticspersistence.Event) []*Event {
	result := make([]*Event, 0, len(values))
	for _, value := range values {
		result = append(result, analyticsEvent(value))
	}
	return result
}

func analyticsEvent(value *analyticspersistence.Event) *Event {
	if value == nil {
		return nil
	}
	result := &Event{
		ID:           value.ID,
		SessionID:    value.SessionID,
		Time:         value.Time,
		Hour:         value.Hour,
		Day:          value.Day,
		Path:         value.Path,
		DefinitionID: value.DefinitionID,
		Definition:   analyticsEventDefinition(value.Definition),
	}
	for _, data := range value.Data {
		if data == nil {
			continue
		}
		result.Data = append(result.Data, &EventData{
			ID:      data.ID,
			EventID: data.EventID,
			FieldID: data.FieldID,
			Value:   data.Value,
			Field:   analyticsEventField(data.Field),
		})
	}
	return result
}

func analyticsEventDefinition(value *eventpersistence.Definition) *EventDefinition {
	if value == nil {
		return nil
	}
	result := &EventDefinition{
		ID:        value.ID,
		SiteID:    value.SiteID,
		Name:      value.Name,
		CreatedAt: value.CreatedAt,
		UpdatedAt: value.UpdatedAt,
	}
	for _, field := range value.Fields {
		result.Fields = append(result.Fields, analyticsEventField(field))
	}
	return result
}

func analyticsEventField(value *eventpersistence.Field) *EventField {
	if value == nil {
		return nil
	}
	return &EventField{
		ID:                value.ID,
		EventDefinitionID: value.EventDefinitionID,
		Key:               value.Key,
		Type:              EventFieldType(value.Type),
		Required:          value.Required,
		MaxLength:         value.MaxLength,
		CreatedAt:         value.CreatedAt,
		UpdatedAt:         value.UpdatedAt,
	}
}
