package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/uptrace/bun"
)

type AnalyticsRepository struct {
	db *bun.DB
}

type AnalyticsFilter struct {
	Referrer           []string
	Browser            []string
	Device             []string
	OS                 []string
	Page               []string
	Country            []string
	EventTypes         []EventType
	EventName          []string
	EventPath          []string
	EventDefinitionIDs []int64
}

type EventType string

const (
	EventTypePageView   EventType = "PAGE_VIEW"
	EventTypePredefined EventType = "PREDEFINED"
)

type AnalyticsQuery struct {
	SiteID int64
	From   time.Time
	To     time.Time
	Limit  int
	Offset int
	Bucket TimeBucket
	Filter AnalyticsFilter
}

func NewAnalyticsRepository(db *bun.DB) *AnalyticsRepository {
	return &AnalyticsRepository{db: db}
}

func (r *AnalyticsRepository) RunInTx(ctx context.Context, fn func(ctx context.Context, tx bun.Tx) error) error {
	if err := r.db.RunInTx(ctx, nil, fn); err != nil {
		return fmt.Errorf("run analytics transaction: %w", err)
	}
	return nil
}
