package model

import (
	"time"

	"github.com/lovely-eye/server/internal/analytics"
)

type DashboardStats struct {
	Visitors    int
	PageViews   int
	Sessions    int
	BounceRate  float64
	AvgDuration float64

	SiteID int64
	From   time.Time
	To     time.Time
	Filter analytics.Filter
}
