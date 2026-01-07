package graph

import (
	"time"

	"github.com/lovely-eye/server/internal/graph/model"
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

func convertToGraphQLStats(stats *services.DashboardStats) *model.DashboardStats {
	result := &model.DashboardStats{
		Visitors:    stats.Visitors,
		PageViews:   stats.PageViews,
		Sessions:    stats.Sessions,
		BounceRate:  stats.BounceRate,
		AvgDuration: stats.AvgDuration,
	}

	for _, p := range stats.TopPages {
		result.TopPages = append(result.TopPages, &model.PageStats{
			Path:     p.Path,
			Views:    p.Views,
			Visitors: p.Visitors,
		})
	}

	for _, ref := range stats.TopReferrers {
		result.TopReferrers = append(result.TopReferrers, &model.ReferrerStats{
			Referrer: ref.Referrer,
			Visitors: ref.Visitors,
		})
	}

	for _, b := range stats.Browsers {
		result.Browsers = append(result.Browsers, &model.BrowserStats{
			Browser:  b.Browser,
			Visitors: b.Visitors,
		})
	}

	for _, d := range stats.Devices {
		result.Devices = append(result.Devices, &model.DeviceStats{
			Device:   d.Device,
			Visitors: d.Visitors,
		})
	}

	for _, c := range stats.Countries {
		result.Countries = append(result.Countries, &model.CountryStats{
			Country:  c.Country,
			Visitors: c.Visitors,
		})
	}

	for _, d := range stats.DailyStats {
		result.DailyStats = append(result.DailyStats, &model.DailyStats{
			Date:      d.Date,
			Visitors:  d.Visitors,
			PageViews: d.PageViews,
			Sessions:  d.Sessions,
		})
	}

	return result
}
