package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/cookiejar"
	"testing"

	"github.com/Khan/genqlient/graphql"
	operations "github.com/lovely-eye/server/e2e/generated"
	"github.com/stretchr/testify/require"
)

func TestDashboardFiltering(t *testing.T) {
	ts := newTestServer(t)
	ctx := context.Background()
	defaultPaging := operations.PagingInput{Limit: 50, Offset: 0}

	// Create HTTP client with cookie jar (simulates browser)
	jar, err := cookiejar.New(nil)
	require.NoError(t, err)
	httpClient := &http.Client{Jar: jar}

	// GraphQL client with cookie support
	client := graphql.NewClient(ts.httpServer.URL+"/graphql", httpClient)

	// Register user
	_, err = operations.Register(ctx, client, operations.RegisterInput{
		Username: "admin",
		Password: "password",
	})
	require.NoError(t, err)

	// Login sets cookies in jar
	_, err = operations.Login(ctx, client, operations.LoginInput{
		Username: "admin",
		Password: "password",
	})
	require.NoError(t, err)

	siteResp, err := operations.CreateSite(ctx, client, operations.CreateSiteInput{
		Domains: []string{"filter-test.com"},
		Name:    "Filter Test Site",
	})
	require.NoError(t, err)

	siteKey := siteResp.CreateSite.PublicKey
	siteID := siteResp.CreateSite.Id

	// Collect test data with different referrers, devices, and pages
	testData := []struct {
		path        string
		referrer    string
		userAgent   string
		screenWidth int
	}{
		// Desktop from Google viewing /home
		{"/home", "https://google.com/search", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/120.0", 1920},
		{"/home", "https://google.com/search", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/120.0", 1920},

		// Mobile from Facebook viewing /about
		{"/about", "https://facebook.com", "Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) Safari/604.1", 375},
		{"/about", "https://facebook.com", "Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) Safari/604.1", 375},

		// Desktop direct traffic viewing /products
		{"/products", "", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) Safari/605.1.15", 1440},
		{"/products", "", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) Safari/605.1.15", 1440},

		// Mobile from Twitter viewing /home
		{"/home", "https://twitter.com", "Mozilla/5.0 (Android 13; Mobile) Chrome/119.0", 412},
	}

	// Collect all page views
	for i, data := range testData {
		payload := map[string]interface{}{
			"site_key":     siteKey,
			"path":         data.path,
			"title":        "Test Page",
			"referrer":     data.referrer,
			"screen_width": data.screenWidth,
		}
		body, _ := json.Marshal(payload)

		req, _ := http.NewRequest("POST", ts.httpServer.URL+"/api/collect", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Origin", "https://filter-test.com")
		// Append index to user agent to create different visitor IDs
		// This simulates different visitors to avoid deduplication
		req.Header.Set("User-Agent", data.userAgent+" TestVisitor/"+string(rune('A'+i)))

		resp, err := ts.httpServer.Client().Do(req)
		require.NoError(t, err)
		err = resp.Body.Close()
		if nil != err {
			require.NoError(t,err)
		}
		require.Equal(t, http.StatusNoContent, resp.StatusCode)
	}

	t.Run("no filter shows all data", func(t *testing.T) {
		resp, err := operations.Dashboard(ctx, client, siteID, nil, nil, defaultPaging, defaultPaging, defaultPaging, defaultPaging, defaultPaging, nil, nil)
		require.NoError(t, err)

		// Should have all page views
		require.Equal(t, 7, resp.Dashboard.PageViews)

		// Should have multiple referrers
		require.GreaterOrEqual(t, len(resp.Dashboard.TopReferrers.Items), 3)

		// Should have multiple pages
		require.GreaterOrEqual(t, len(resp.Dashboard.TopPages.Items), 3)
	})

	t.Run("filter by referrer", func(t *testing.T) {
		googleReferrer := "https://google.com/search"
		filter := &operations.FilterInput{
			Referrer: []string{googleReferrer},
		}

		resp, err := operations.Dashboard(ctx, client, siteID, nil, filter, defaultPaging, defaultPaging, defaultPaging, defaultPaging, defaultPaging, nil, nil)
		require.NoError(t, err)

		// Should only show Google traffic (2 page views)
		require.Equal(t, 2, resp.Dashboard.PageViews, "should only count Google referrer page views")

		// Should only show /home page (all Google traffic went to /home)
		require.Len(t, resp.Dashboard.TopPages.Items, 1, "should only show pages visited via Google")
		require.Equal(t, "/home", resp.Dashboard.TopPages.Items[0].Path)
		require.Equal(t, 2, resp.Dashboard.TopPages.Items[0].Views)

		// Should only show desktop device (Google traffic is desktop)
		require.Len(t, resp.Dashboard.Devices.Items, 1, "should only show devices from Google traffic")
		require.Equal(t, "desktop", resp.Dashboard.Devices.Items[0].Device)
	})

	t.Run("filter by device", func(t *testing.T) {
		mobileDevice := "mobile"
		filter := &operations.FilterInput{
			Device: []string{mobileDevice},
		}

		resp, err := operations.Dashboard(ctx, client, siteID, nil, filter, defaultPaging, defaultPaging, defaultPaging, defaultPaging, defaultPaging, nil, nil)
		require.NoError(t, err)

		// Should only show mobile traffic (3 page views: 2 from Facebook, 1 from Twitter)
		require.Equal(t, 3, resp.Dashboard.PageViews, "should only count mobile page views")

		// Should show pages visited on mobile
		require.GreaterOrEqual(t, len(resp.Dashboard.TopPages.Items), 1, "should show pages visited on mobile")

		// Find /about in top pages (Facebook mobile traffic)
		foundAbout := false
		for _, page := range resp.Dashboard.TopPages.Items {
			if page.Path == "/about" {
				foundAbout = true
				require.Equal(t, 2, page.Views, "/about should have 2 mobile views from Facebook")
			}
		}
		require.True(t, foundAbout, "should find /about page in mobile traffic")
	})

	t.Run("filter by page", func(t *testing.T) {
		homePage := "/home"
		filter := &operations.FilterInput{
			Page: []string{homePage},
		}

		resp, err := operations.Dashboard(ctx, client, siteID, nil, filter, defaultPaging, defaultPaging, defaultPaging, defaultPaging, defaultPaging, nil, nil)
		require.NoError(t, err)

		// Should show traffic to /home page (3 views: 2 Google desktop + 1 Twitter mobile)
		require.Equal(t, 3, resp.Dashboard.PageViews, "should only count /home page views")

		// Should show only /home in top pages
		require.Len(t, resp.Dashboard.TopPages.Items, 1, "should only show /home page")
		require.Equal(t, "/home", resp.Dashboard.TopPages.Items[0].Path)

		// Should show referrers that led to /home
		require.GreaterOrEqual(t, len(resp.Dashboard.TopReferrers.Items), 2, "should show Google and Twitter referrers")
	})

	t.Run("filter by multiple criteria - referrer and device", func(t *testing.T) {
		facebookReferrer := "https://facebook.com"
		mobileDevice := "mobile"
		filter := &operations.FilterInput{
			Referrer: []string{facebookReferrer},
			Device:   []string{mobileDevice},
		}

		resp, err := operations.Dashboard(ctx, client, siteID, nil, filter, defaultPaging, defaultPaging, defaultPaging, defaultPaging, defaultPaging, nil, nil)
		require.NoError(t, err)

		// Should show only Facebook mobile traffic (2 page views to /about)
		require.Equal(t, 2, resp.Dashboard.PageViews, "should only count Facebook mobile page views")

		// Should only show /about page
		require.Len(t, resp.Dashboard.TopPages.Items, 1, "should only show /about page")
		require.Equal(t, "/about", resp.Dashboard.TopPages.Items[0].Path)
		require.Equal(t, 2, resp.Dashboard.TopPages.Items[0].Views)
	})

	t.Run("filter by multiple criteria - page and device", func(t *testing.T) {
		homePage := "/home"
		desktopDevice := "desktop"
		filter := &operations.FilterInput{
			Page:   []string{homePage},
			Device: []string{desktopDevice},
		}

		resp, err := operations.Dashboard(ctx, client, siteID, nil, filter, defaultPaging, defaultPaging, defaultPaging, defaultPaging, defaultPaging, nil, nil)
		require.NoError(t, err)

		// Should show only desktop traffic to /home (2 Google views)
		require.Equal(t, 2, resp.Dashboard.PageViews, "should only count desktop views of /home")

		// Should show only Google referrer
		require.Len(t, resp.Dashboard.TopReferrers.Items, 1, "should only show Google referrer")
		require.Contains(t, resp.Dashboard.TopReferrers.Items[0].Referrer, "google")
	})

	t.Run("filter with non-existent values returns empty results", func(t *testing.T) {
		nonExistentReferrer := "https://nonexistent.com"
		filter := &operations.FilterInput{
			Referrer: []string{nonExistentReferrer},
		}

		resp, err := operations.Dashboard(ctx, client, siteID, nil, filter, defaultPaging, defaultPaging, defaultPaging, defaultPaging, defaultPaging, nil, nil)
		require.NoError(t, err)

		// Should have no data
		require.Equal(t, 0, resp.Dashboard.PageViews)
		require.Equal(t, 0, resp.Dashboard.Visitors)
		require.Equal(t, 0, resp.Dashboard.Sessions)
		require.Len(t, resp.Dashboard.TopPages.Items, 0)
		require.Len(t, resp.Dashboard.TopReferrers.Items, 0)
	})

	t.Run("direct traffic filter", func(t *testing.T) {
		directReferrer := "(direct)"
		filter := &operations.FilterInput{
			Referrer: []string{directReferrer},
		}

		resp, err := operations.Dashboard(ctx, client, siteID, nil, filter, defaultPaging, defaultPaging, defaultPaging, defaultPaging, defaultPaging, nil, nil)
		require.NoError(t, err)

		// Should show direct traffic (2 views to /products)
		require.Equal(t, 2, resp.Dashboard.PageViews, "should count direct traffic page views")

		// Should only show /products page
		require.Len(t, resp.Dashboard.TopPages.Items, 1, "should only show /products page")
		require.Equal(t, "/products", resp.Dashboard.TopPages.Items[0].Path)
	})
}
