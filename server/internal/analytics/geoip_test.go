package analytics

import (
	"errors"
	"testing"

	geoipcore "github.com/lovely-eye/server/internal/geoip"
	"github.com/stretchr/testify/require"
)

func TestService_ResolveCountryBestEffort_NoReader(t *testing.T) {
	t.Parallel()

	svc := &fakeGeoIPProvider{
		resolveErr: geoipcore.ErrNoDBReader,
	}

	analytics := NewService(nil, nil, nil, svc, nil, "test-analytics-identity-secret-32chars")
	country := analytics.resolveCountryBestEffort("8.8.8.8")

	require.Equal(t, UnknownCountry, country)
}

func TestService_ResolveCountryBestEffort_UnexpectedErrorFallsBack(t *testing.T) {
	t.Parallel()

	svc := &fakeGeoIPProvider{
		resolveErr: errors.New("broken reader"),
	}

	analytics := NewService(nil, nil, nil, svc, nil, "test-analytics-identity-secret-32chars")
	country := analytics.resolveCountryBestEffort("8.8.8.8")

	require.Equal(t, UnknownCountry, country)
}
