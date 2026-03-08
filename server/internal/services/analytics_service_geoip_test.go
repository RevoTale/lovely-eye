package services

import (
	"errors"
	"testing"

	geoipcore "github.com/lovely-eye/server/internal/geoip"
	"github.com/stretchr/testify/require"
)

func TestAnalyticsService_ResolveCountryBestEffort_NoReader(t *testing.T) {
	t.Parallel()

	svc := NewGeoIPService(GeoIPConfig{DBPath: "/tmp/test.mmdb"})
	svc.lookup = &fakeGeoIPLookup{
		resolveErr: geoipcore.ErrNoDBReader,
	}

	analytics := NewAnalyticsService(nil, nil, nil, svc)
	country := analytics.resolveCountryBestEffort("8.8.8.8")

	require.Equal(t, UnknownCountry, country)
}

func TestAnalyticsService_ResolveCountryBestEffort_UnexpectedErrorFallsBack(t *testing.T) {
	t.Parallel()

	svc := NewGeoIPService(GeoIPConfig{DBPath: "/tmp/test.mmdb"})
	svc.lookup = &fakeGeoIPLookup{
		resolveErr: errors.New("broken reader"),
	}

	analytics := NewAnalyticsService(nil, nil, nil, svc)
	country := analytics.resolveCountryBestEffort("8.8.8.8")

	require.Equal(t, UnknownCountry, country)
}
