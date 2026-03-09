package e2e

import (
	"context"
	"strconv"
	"testing"
	"time"

	operations "github.com/lovely-eye/server/e2e/generated"
	"github.com/lovely-eye/server/internal/models"
	"github.com/stretchr/testify/require"
)

func TestDashboardCountryNamesAreResolvedLazily(t *testing.T) {
	ts := newTestServer(t)
	ctx := context.Background()

	_, err := operations.Register(ctx, ts.graphqlClient(), operations.RegisterInput{
		Username: "admin",
		Password: "password",
	})
	require.NoError(t, err)

	client := ts.authenticatedClient(ctx, t, "admin", "password")

	siteResp, err := operations.CreateSite(ctx, client, operations.CreateSiteInput{
		Domains: []string{"example.com"},
		Name:    "Example Site",
	})
	require.NoError(t, err)

	siteID, err := strconv.ParseInt(siteResp.CreateSite.Id, 10, 64)
	require.NoError(t, err)

	_, err = ts.DB.NewInsert().Model(&models.Country{
		Code: "US",
		Name: "United States",
	}).Exec(ctx)
	require.NoError(t, err)

	now := time.Now().Unix()
	countryCodes := []string{"US", "", "ZZ"}
	for index, countryCode := range countryCodes {
		dbClient := &models.Client{
			SiteID:  siteID,
			Hash:    "hash-" + strconv.Itoa(index),
			Country: countryCode,
			Device:  "desktop",
			Browser: "Chrome",
			OS:      "Linux",
		}
		_, err = ts.DB.NewInsert().Model(dbClient).Exec(ctx)
		require.NoError(t, err)

		session := &models.Session{
			SiteID:        siteID,
			ClientID:      dbClient.ID,
			EnterTime:     now,
			EnterHour:     now / 3600,
			EnterDay:      now / 86400,
			EnterPath:     "/",
			ExitTime:      now,
			ExitHour:      now / 3600,
			ExitDay:       now / 86400,
			ExitPath:      "/",
			Duration:      1,
			PageViewCount: 1,
		}
		_, err = ts.DB.NewInsert().Model(session).Exec(ctx)
		require.NoError(t, err)
	}

	resp, err := operations.Dashboard(
		ctx,
		client,
		siteResp.CreateSite.Id,
		nil,
		nil,
		defaultPaging,
		defaultPaging,
		defaultPaging,
		defaultPaging,
		defaultPaging,
		nil,
		nil,
	)
	require.NoError(t, err)

	countryNames := make(map[string]string, len(resp.Dashboard.Countries.Items))
	for _, countryStat := range resp.Dashboard.Countries.Items {
		countryNames[countryStat.Country.Code] = countryStat.Country.Name
	}

	require.Equal(t, "United States", countryNames["US"])
	require.Equal(t, "Unknown", countryNames["-"])
	require.Equal(t, "ZZ", countryNames["ZZ"])
}
