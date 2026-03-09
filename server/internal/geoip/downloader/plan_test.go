package downloader

import (
	"testing"
	"time"

	"github.com/lovely-eye/server/internal/geoip"
	"github.com/stretchr/testify/require"
)

func TestBuildDBIPCandidateURLs(t *testing.T) {
	t.Parallel()

	urls, err := buildDBIPCandidateURLs(
		"https://download.db-ip.com/free/dbip-country-lite.mmdb.gz",
		time.Date(2026, time.March, 8, 0, 0, 0, 0, time.UTC),
	)
	require.NoError(t, err)
	require.Equal(t, []string{
		"https://download.db-ip.com/free/dbip-country-lite-2026-03.mmdb.gz",
		"https://download.db-ip.com/free/dbip-country-lite-2026-02.mmdb.gz",
		"https://download.db-ip.com/free/dbip-country-lite-2026-01.mmdb.gz",
		"https://download.db-ip.com/free/dbip-country-lite.mmdb.gz",
	}, urls)
}

func TestBuildDBIPCandidateURLs_MonthlyURL(t *testing.T) {
	t.Parallel()

	urls, err := buildDBIPCandidateURLs(
		"https://download.db-ip.com/free/dbip-country-lite-2026-03.mmdb.gz",
		time.Date(2026, time.March, 8, 0, 0, 0, 0, time.UTC),
	)
	require.NoError(t, err)
	require.Equal(t, []string{
		"https://download.db-ip.com/free/dbip-country-lite-2026-03.mmdb.gz",
	}, urls)
}

func TestBuildDownloadPlan_MaxMind(t *testing.T) {
	t.Parallel()

	downloader := New(geoip.Config{
		DBPath:            "/tmp/GeoLite2-Country.mmdb",
		MaxMindLicenseKey: "license-key",
	})

	plan, err := downloader.BuildDownloadPlan()
	require.NoError(t, err)
	require.Equal(t, geoip.SourceMaxMind, plan.Source)
	require.Len(t, plan.CandidateURLs, 1)
	require.True(t, isTarGzDownload(plan.CandidateURLs[0]))
}
