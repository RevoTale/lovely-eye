package graph

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	analyticfeature "github.com/lovely-eye/server/internal/analytics"
	"github.com/lovely-eye/server/internal/auth"
	"github.com/lovely-eye/server/internal/event"
	"github.com/lovely-eye/server/internal/graph/model"
)

// UpsertEventDefinition is the resolver for the upsertEventDefinition field.
func (r *mutationResolver) UpsertEventDefinition(ctx context.Context, siteID string, input model.EventDefinitionInput) (*model.EventDefinition, error) {
	claims := auth.GetUserFromContext(ctx)
	if claims == nil {
		return nil, unauthenticated()
	}

	id, err := strconv.ParseInt(siteID, 10, 64)
	if err != nil {
		return nil, badUserInput("invalid site ID")
	}

	if err := r.SiteService.RequireOwnership(ctx, id, claims.UserID); err != nil {
		return nil, fmt.Errorf("failed to get site: %w", err)
	}

	fields := make([]event.FieldInput, 0, len(input.Fields))
	for _, field := range input.Fields {
		fields = append(fields, event.FieldInput{
			Key:       field.Key,
			Type:      strings.ToLower(string(field.Type)),
			Required:  field.Required,
			MaxLength: field.MaxLength,
		})
	}

	definition, err := r.EventDefService.Upsert(ctx, id, event.DefinitionInput{
		Name:   input.Name,
		Fields: fields,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to upsert event definition: %w", err)
	}

	results := convertToGraphQLEventDefinitions([]*event.Definition{definition})
	return results[0], nil
}

// DeleteEventDefinition is the resolver for the deleteEventDefinition field.
func (r *mutationResolver) DeleteEventDefinition(ctx context.Context, siteID string, name string) (bool, error) {
	claims := auth.GetUserFromContext(ctx)
	if claims == nil {
		return false, unauthenticated()
	}

	id, err := strconv.ParseInt(siteID, 10, 64)
	if err != nil {
		return false, badUserInput("invalid site ID")
	}

	if err := r.SiteService.RequireOwnership(ctx, id, claims.UserID); err != nil {
		return false, fmt.Errorf("failed to get site: %w", err)
	}

	if err := r.EventDefService.Delete(ctx, id, name); err != nil {
		return false, fmt.Errorf("failed to delete event definition: %w", err)
	}

	return true, nil
}

// Events is the resolver for the events field.
func (r *queryResolver) Events(ctx context.Context, siteID string, dateRange *model.DateRangeInput, filter *model.FilterInput, paging model.PagingInput) (*model.EventsResult, error) {
	claims := auth.GetUserFromContext(ctx)
	if claims == nil {
		return nil, unauthenticated()
	}

	id, err := strconv.ParseInt(siteID, 10, 64)
	if err != nil {
		return nil, badUserInput("invalid site ID")
	}

	if err := r.SiteService.RequireOwnership(ctx, id, claims.UserID); err != nil {
		return nil, fmt.Errorf("failed to get site: %w", err)
	}

	from, to, err := parseDateRangeInput(dateRange, r.DashboardLimits.MaxDailyRangeDays)
	if err != nil {
		return nil, err
	}

	lim, off := normalizePaging(paging)

	filterOpts, err := parseFilterInput(filter, r.DashboardLimits)
	if err != nil {
		return nil, err
	}

	var events []*analyticfeature.Event
	var total int
	if filter == nil || isFilterEmpty(filterOpts) {
		events, total, err = r.AnalyticsService.GetEventsWithTotal(ctx, id, from, to, lim, off)
	} else {
		events, total, err = r.AnalyticsService.GetEventsWithTotalAndFilter(ctx, analyticfeature.Query{
			SiteID: id,
			From:   from,
			To:     to,
			Limit:  lim,
			Offset: off,
			Filter: filterOpts,
		})
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get events: %w", err)
	}

	return convertToGraphQLEvents(events, total), nil
}

// EventCounts is the resolver for the eventCounts field.
func (r *queryResolver) EventCounts(ctx context.Context, siteID string, dateRange *model.DateRangeInput, filter *model.FilterInput, paging model.PagingInput) (*model.EventCountsResult, error) {
	claims := auth.GetUserFromContext(ctx)
	if claims == nil {
		return nil, unauthenticated()
	}

	id, err := strconv.ParseInt(siteID, 10, 64)
	if err != nil {
		return nil, badUserInput("invalid site ID")
	}

	if err := r.SiteService.RequireOwnership(ctx, id, claims.UserID); err != nil {
		return nil, fmt.Errorf("failed to get site: %w", err)
	}

	from, to, err := parseDateRangeInput(dateRange, r.DashboardLimits.MaxDailyRangeDays)
	if err != nil {
		return nil, err
	}

	limit, offset := normalizePaging(paging)

	filterOpts, err := parseFilterInput(filter, r.DashboardLimits)
	if err != nil {
		return nil, err
	}

	eventCounts, total, err := r.AnalyticsService.GetEventCounts(ctx, analyticfeature.Query{
		SiteID: id,
		From:   from,
		To:     to,
		Limit:  limit,
		Offset: offset,
		Filter: filterOpts,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get event counts: %w", err)
	}

	result := make([]*model.EventCount, 0, len(eventCounts))
	for _, ec := range eventCounts {
		result = append(result, &model.EventCount{
			Event: convertToGraphQLEvent(ec.Event),
			Count: ec.Count,
		})
	}

	return &model.EventCountsResult{Items: result, Total: total}, nil
}

// EventDefinitions is the resolver for the eventDefinitions field.
func (r *queryResolver) EventDefinitions(ctx context.Context, siteID string, paging model.PagingInput) ([]*model.EventDefinition, error) {
	claims := auth.GetUserFromContext(ctx)
	if claims == nil {
		return nil, unauthenticated()
	}

	id, err := strconv.ParseInt(siteID, 10, 64)
	if err != nil {
		return nil, badUserInput("invalid site ID")
	}

	if err := r.SiteService.RequireOwnership(ctx, id, claims.UserID); err != nil {
		return nil, fmt.Errorf("failed to get site: %w", err)
	}

	limit, offset := normalizePaging(paging)
	definitions, err := r.EventDefService.List(ctx, id, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list event definitions: %w", err)
	}

	return convertToGraphQLEventDefinitions(definitions), nil
}
