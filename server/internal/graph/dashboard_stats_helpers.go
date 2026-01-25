package graph

import (
	"github.com/lovely-eye/server/internal/graph/model"
)

const (
	maxPageSize         = 100
	maxTimeSeriesPoints = 1000
	defaultDailyPoints  = 365
	defaultHourlyPoints = 168
	defaultEventsPage   = 50
)

func normalizePaging(paging model.PagingInput) (int, int) {
	limit := clampLimit(paging.Limit, maxPageSize)
	offset := max(paging.Offset, 0)
	return limit, offset
}

func clampLimit(value, max int) int {
	if value <= 0 {
		return 1
	}
	if value > max {
		return max
	}
	return value
}


func bucketValueOrDefault(value *model.TimeBucket) model.TimeBucket {
	if value == nil {
		return model.TimeBucketDaily
	}
	return *value
}
