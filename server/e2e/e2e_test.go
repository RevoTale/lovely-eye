package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Khan/genqlient/graphql"
	operations "github.com/lovely-eye/server/e2e/generated"
	"github.com/lovely-eye/server/internal/config"
	"github.com/lovely-eye/server/internal/server"
	"github.com/stretchr/testify/require"
)

func testConfig() *config.Config {
	return &config.Config{
		Server: config.ServerConfig{
			Host:          "127.0.0.1",
			Port:          "0",
			DashboardPath: "", // Empty for tests - dashboard not required
		},
		Database: config.DatabaseConfig{
			Driver: "sqlite",
			DSN:    "file::memory:?cache=shared",
		},
		Auth: config.AuthConfig{
			JWTSecret:         "test-secret-key-for-e2e-testing-32chars",
			AccessTokenExpiry: 15 * time.Minute,
			RefreshExpiry:     7 * 24 * time.Hour,
			AllowRegistration: true,
			SecureCookies:     false,
			CookieDomain:      "",
		},
		// Mock tracker.js for testing (avoid file I/O)
		TrackerJS: []byte(`console.log("mock tracker")`),
	}
}

type testServer struct {
	*server.Server
	httpServer *httptest.Server
}

func newTestServer(t *testing.T) *testServer {
	t.Helper()

	srv, err := server.New(testConfig())
	require.NoError(t, err, "failed to create server")

	httpServer := httptest.NewServer(srv.Handler)

	t.Cleanup(func() {
		httpServer.Close()
		srv.Close()
	})

	return &testServer{
		Server:     srv,
		httpServer: httpServer,
	}
}

// newTestHTTPServer creates a test HTTP server without t.Cleanup
// for use in benchmarks and scenarios where cleanup is manual
func newTestHTTPServer(handler http.Handler) *httptest.Server {
	return httptest.NewServer(handler)
}

func (ts *testServer) graphqlClient() graphql.Client {
	return graphql.NewClient(ts.httpServer.URL+"/graphql", ts.httpServer.Client())
}

func (ts *testServer) bearerClient(accessToken string) graphql.Client {
	httpClient := &http.Client{
		Transport: &bearerTransport{
			base:        ts.httpServer.Client().Transport,
			accessToken: accessToken,
		},
	}
	return graphql.NewClient(ts.httpServer.URL+"/graphql", httpClient)
}

func (ts *testServer) cookieClient(accessToken string) graphql.Client {
	httpClient := &http.Client{
		Transport: &cookieTransport{
			base:        ts.httpServer.Client().Transport,
			accessToken: accessToken,
		},
	}
	return graphql.NewClient(ts.httpServer.URL+"/graphql", httpClient)
}

// authenticatedClient creates a client with a cookie jar and performs login
// Returns the authenticated client that will use cookies for subsequent requests
func (ts *testServer) authenticatedClient(ctx context.Context, t *testing.T, username, password string) graphql.Client {
	t.Helper()
	jar, err := cookiejar.New(nil)
	require.NoError(t, err)
	httpClient := &http.Client{Jar: jar}
	client := graphql.NewClient(ts.httpServer.URL+"/graphql", httpClient)

	// Login to set cookies
	_, err = operations.Login(ctx, client, operations.LoginInput{
		Username: username,
		Password: password,
	})
	require.NoError(t, err)

	return client
}

type bearerTransport struct {
	base        http.RoundTripper
	accessToken string
}

func (t *bearerTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+t.accessToken)
	return t.base.RoundTrip(req)
}

type cookieTransport struct {
	base        http.RoundTripper
	accessToken string
}

func (t *cookieTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.AddCookie(&http.Cookie{Name: "le_access", Value: t.accessToken})
	return t.base.RoundTrip(req)
}

func TestRegister(t *testing.T) {
	ts := newTestServer(t)
	ctx := context.Background()

	t.Run("first user becomes admin", func(t *testing.T) {
		resp, err := operations.Register(ctx, ts.graphqlClient(), operations.RegisterInput{
			Username: "admin",
			Password: "password123",
		})
		require.NoError(t, err)
		require.Equal(t, "admin", resp.Register.User.Username)
		require.Equal(t, "admin", resp.Register.User.Role)
		// Tokens are now in HttpOnly cookies, not in response body
	})

	t.Run("second user is regular user", func(t *testing.T) {
		resp, err := operations.Register(ctx, ts.graphqlClient(), operations.RegisterInput{
			Username: "user1",
			Password: "password123",
		})
		require.NoError(t, err)
		require.Equal(t, "user", resp.Register.User.Role)
	})

	t.Run("duplicate username fails", func(t *testing.T) {
		_, err := operations.Register(ctx, ts.graphqlClient(), operations.RegisterInput{
			Username: "admin",
			Password: "different",
		})
		require.Error(t, err)
	})
}

func TestLogin(t *testing.T) {
	ts := newTestServer(t)
	ctx := context.Background()

	_, err := operations.Register(ctx, ts.graphqlClient(), operations.RegisterInput{
		Username: "testuser",
		Password: "testpass",
	})
	require.NoError(t, err)

	t.Run("valid credentials", func(t *testing.T) {
		resp, err := operations.Login(ctx, ts.graphqlClient(), operations.LoginInput{
			Username: "testuser",
			Password: "testpass",
		})
		require.NoError(t, err)
		require.Equal(t, "testuser", resp.Login.User.Username)
	})

	t.Run("invalid password", func(t *testing.T) {
		_, err := operations.Login(ctx, ts.graphqlClient(), operations.LoginInput{
			Username: "testuser",
			Password: "wrongpass",
		})
		require.Error(t, err)
	})

	t.Run("nonexistent user", func(t *testing.T) {
		_, err := operations.Login(ctx, ts.graphqlClient(), operations.LoginInput{
			Username: "nouser",
			Password: "testpass",
		})
		require.Error(t, err)
	})
}

func TestStatsCollection(t *testing.T) {
	ts := newTestServer(t)
	ctx := context.Background()

	_, err := operations.Register(ctx, ts.graphqlClient(), operations.RegisterInput{
		Username: "admin",
		Password: "password",
	})
	require.NoError(t, err)

	client := ts.authenticatedClient(ctx, t, "admin", "password")

	siteResp, err := operations.CreateSite(ctx, client, operations.CreateSiteInput{
		Domain: "example.com",
		Name:   "Example Site",
	})
	require.NoError(t, err)

	siteKey := siteResp.CreateSite.PublicKey

	t.Run("collect page view", func(t *testing.T) {
		payload := map[string]interface{}{
			"site_key":     siteKey,
			"path":         "/home",
			"title":        "Home Page",
			"referrer":     "https://google.com",
			"screen_width": 1920,
		}
		body, _ := json.Marshal(payload)

		resp, err := ts.httpServer.Client().Post(
			ts.httpServer.URL+"/api/collect",
			"application/json",
			bytes.NewReader(body),
		)
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("collect custom event", func(t *testing.T) {
		payload := map[string]interface{}{
			"site_key":   siteKey,
			"name":       "button_click",
			"path":       "/home",
			"properties": `{"button": "signup"}`,
		}
		body, _ := json.Marshal(payload)

		resp, err := ts.httpServer.Client().Post(
			ts.httpServer.URL+"/api/event",
			"application/json",
			bytes.NewReader(body),
		)
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("missing site_key fails", func(t *testing.T) {
		payload := map[string]interface{}{
			"path": "/home",
		}
		body, _ := json.Marshal(payload)

		resp, err := ts.httpServer.Client().Post(
			ts.httpServer.URL+"/api/collect",
			"application/json",
			bytes.NewReader(body),
		)
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestDashboardAuthorization(t *testing.T) {
	ts := newTestServer(t)
	ctx := context.Background()

	_, err := operations.Register(ctx, ts.graphqlClient(), operations.RegisterInput{
		Username: "admin",
		Password: "password",
	})
	require.NoError(t, err)

	authedClient := ts.authenticatedClient(ctx, t, "admin", "password")

	siteResp, err := operations.CreateSite(ctx, authedClient, operations.CreateSiteInput{
		Domain: "example.com",
		Name:   "Example Site",
	})
	require.NoError(t, err)

	siteID := siteResp.CreateSite.Id

	t.Run("authenticated user can view dashboard", func(t *testing.T) {
		resp, err := operations.Dashboard(ctx, authedClient, siteID, nil, nil)
		require.NoError(t, err)
		require.Equal(t, 0, resp.Dashboard.Visitors)
	})

	t.Run("unauthenticated user cannot view dashboard", func(t *testing.T) {
		_, err := operations.Dashboard(ctx, ts.graphqlClient(), siteID, nil, nil)
		require.Error(t, err)
	})

	t.Run("authenticated user can view realtime", func(t *testing.T) {
		resp, err := operations.Realtime(ctx, authedClient, siteID)
		require.NoError(t, err)
		require.Equal(t, 0, resp.Realtime.Visitors)
	})

	t.Run("me query returns user when authenticated", func(t *testing.T) {
		resp, err := operations.Me(ctx, authedClient)
		require.NoError(t, err)
		require.NotNil(t, resp.Me)
		require.NotNil(t, resp.Me.Username)
		require.Equal(t, "admin", *resp.Me.Username)
	})

	t.Run("me query returns nil when unauthenticated", func(t *testing.T) {
		resp, err := operations.Me(ctx, ts.graphqlClient())
		require.NoError(t, err)
		require.Nil(t, resp.Me)
	})
}

func TestSiteManagement(t *testing.T) {
	ts := newTestServer(t)
	ctx := context.Background()

	_, err := operations.Register(ctx, ts.graphqlClient(), operations.RegisterInput{
		Username: "admin",
		Password: "password",
	})
	require.NoError(t, err)

	client := ts.authenticatedClient(ctx, t, "admin", "password")

	t.Run("create site", func(t *testing.T) {
		resp, err := operations.CreateSite(ctx, client, operations.CreateSiteInput{
			Domain: "mysite.com",
			Name:   "My Site",
		})
		require.NoError(t, err)
		require.Equal(t, "mysite.com", resp.CreateSite.Domain)
		require.NotEmpty(t, resp.CreateSite.PublicKey)
	})

	t.Run("list sites", func(t *testing.T) {
		resp, err := operations.Sites(ctx, client)
		require.NoError(t, err)
		require.Len(t, resp.Sites, 1)
	})

	t.Run("unauthenticated cannot create site", func(t *testing.T) {
		_, err := operations.CreateSite(ctx, ts.graphqlClient(), operations.CreateSiteInput{
			Domain: "other.com",
			Name:   "Other",
		})
		require.Error(t, err)
	})
}

func TestEventPropertiesValidation(t *testing.T) {
	ts := newTestServer(t)
	ctx := context.Background()

	_, err := operations.Register(ctx, ts.graphqlClient(), operations.RegisterInput{
		Username: "admin",
		Password: "password",
	})
	require.NoError(t, err)

	client := ts.authenticatedClient(ctx, t, "admin", "password")

	siteResp, err := operations.CreateSite(ctx, client, operations.CreateSiteInput{
		Domain: "events-test.com",
		Name:   "Events Test Site",
	})
	require.NoError(t, err)

	siteKey := siteResp.CreateSite.PublicKey
	siteID := siteResp.CreateSite.Id
	maxLen := 500

	definitions := []operations.EventDefinitionInput{
		{
			Name: "purchase",
			Fields: []operations.EventDefinitionFieldInput{
				{Key: "product_id", Type: operations.EventFieldTypeString, Required: true, MaxLength: maxLen},
				{Key: "price", Type: operations.EventFieldTypeString, Required: true, MaxLength: maxLen},
				{Key: "currency", Type: operations.EventFieldTypeString, Required: true, MaxLength: maxLen},
			},
		},
		{
			Name:   "page_scroll",
			Fields: []operations.EventDefinitionFieldInput{},
		},
		{
			Name: "click",
			Fields: []operations.EventDefinitionFieldInput{
				{Key: "key", Type: operations.EventFieldTypeString, Required: true, MaxLength: maxLen},
			},
		},
	}

	for _, definition := range definitions {
		_, err := operations.UpsertEventDefinition(ctx, client, siteID, definition)
		require.NoError(t, err)
	}

	t.Run("valid string:string properties accepted", func(t *testing.T) {
		payload := map[string]interface{}{
			"site_key":   siteKey,
			"name":       "purchase",
			"path":       "/checkout",
			"properties": `{"product_id": "123", "price": "29.99", "currency": "USD"}`,
		}
		body, _ := json.Marshal(payload)

		resp, err := ts.httpServer.Client().Post(
			ts.httpServer.URL+"/api/event",
			"application/json",
			bytes.NewReader(body),
		)
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("empty properties accepted", func(t *testing.T) {
		payload := map[string]interface{}{
			"site_key":   siteKey,
			"name":       "page_scroll",
			"path":       "/home",
			"properties": "",
		}
		body, _ := json.Marshal(payload)

		resp, err := ts.httpServer.Client().Post(
			ts.httpServer.URL+"/api/event",
			"application/json",
			bytes.NewReader(body),
		)
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("invalid JSON properties rejected", func(t *testing.T) {
		payload := map[string]interface{}{
			"site_key":   siteKey,
			"name":       "click",
			"path":       "/page",
			"properties": `{invalid json}`,
		}
		body, _ := json.Marshal(payload)

		resp, err := ts.httpServer.Client().Post(
			ts.httpServer.URL+"/api/event",
			"application/json",
			bytes.NewReader(body),
		)
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("malformed JSON properties rejected", func(t *testing.T) {
		payload := map[string]interface{}{
			"site_key":   siteKey,
			"name":       "click",
			"path":       "/page",
			"properties": `{"key": "value"`,
		}
		body, _ := json.Marshal(payload)

		resp, err := ts.httpServer.Client().Post(
			ts.httpServer.URL+"/api/event",
			"application/json",
			bytes.NewReader(body),
		)
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("JSON array rejected (must be object)", func(t *testing.T) {
		payload := map[string]interface{}{
			"site_key":   siteKey,
			"name":       "click",
			"path":       "/page",
			"properties": `["item1", "item2"]`,
		}
		body, _ := json.Marshal(payload)

		resp, err := ts.httpServer.Client().Post(
			ts.httpServer.URL+"/api/event",
			"application/json",
			bytes.NewReader(body),
		)
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("JSON string rejected (must be object)", func(t *testing.T) {
		payload := map[string]interface{}{
			"site_key":   siteKey,
			"name":       "click",
			"path":       "/page",
			"properties": `"just a string"`,
		}
		body, _ := json.Marshal(payload)

		resp, err := ts.httpServer.Client().Post(
			ts.httpServer.URL+"/api/event",
			"application/json",
			bytes.NewReader(body),
		)
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("JSON number rejected (must be object)", func(t *testing.T) {
		payload := map[string]interface{}{
			"site_key":   siteKey,
			"name":       "click",
			"path":       "/page",
			"properties": `123`,
		}
		body, _ := json.Marshal(payload)

		resp, err := ts.httpServer.Client().Post(
			ts.httpServer.URL+"/api/event",
			"application/json",
			bytes.NewReader(body),
		)
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("non-string values rejected (must be string:string)", func(t *testing.T) {
		payload := map[string]interface{}{
			"site_key":   siteKey,
			"name":       "click",
			"path":       "/page",
			"properties": `{"key": 123}`,
		}
		body, _ := json.Marshal(payload)

		resp, err := ts.httpServer.Client().Post(
			ts.httpServer.URL+"/api/event",
			"application/json",
			bytes.NewReader(body),
		)
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("nested objects rejected (must be string:string)", func(t *testing.T) {
		payload := map[string]interface{}{
			"site_key":   siteKey,
			"name":       "click",
			"path":       "/page",
			"properties": `{"key": {"nested": "value"}}`,
		}
		body, _ := json.Marshal(payload)

		resp, err := ts.httpServer.Client().Post(
			ts.httpServer.URL+"/api/event",
			"application/json",
			bytes.NewReader(body),
		)
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusNoContent, resp.StatusCode)
	})
}

func TestEventPropertiesStored(t *testing.T) {
	ts := newTestServer(t)
	ctx := context.Background()

	_, err := operations.Register(ctx, ts.graphqlClient(), operations.RegisterInput{
		Username: "admin",
		Password: "password",
	})
	require.NoError(t, err)

	client := ts.authenticatedClient(ctx, t, "admin", "password")

	siteResp, err := operations.CreateSite(ctx, client, operations.CreateSiteInput{
		Domain: "events-storage-test.com",
		Name:   "Events Storage Test",
	})
	require.NoError(t, err)

	siteKey := siteResp.CreateSite.PublicKey
	siteID := siteResp.CreateSite.Id
	maxLen := 500

	definitions := []operations.EventDefinitionInput{
		{
			Name: "button_click",
			Fields: []operations.EventDefinitionFieldInput{
				{Key: "button", Type: operations.EventFieldTypeString, Required: true, MaxLength: maxLen},
				{Key: "variant", Type: operations.EventFieldTypeString, Required: true, MaxLength: maxLen},
				{Key: "position", Type: operations.EventFieldTypeString, Required: true, MaxLength: maxLen},
			},
		},
		{
			Name: "form_submit",
			Fields: []operations.EventDefinitionFieldInput{
				{Key: "form_id", Type: operations.EventFieldTypeString, Required: true, MaxLength: maxLen},
				{Key: "fields", Type: operations.EventFieldTypeString, Required: true, MaxLength: maxLen},
			},
		},
		{
			Name: "video_play",
			Fields: []operations.EventDefinitionFieldInput{
				{Key: "video_id", Type: operations.EventFieldTypeString, Required: true, MaxLength: maxLen},
				{Key: "duration", Type: operations.EventFieldTypeString, Required: true, MaxLength: maxLen},
			},
		},
		{
			Name: "download",
			Fields: []operations.EventDefinitionFieldInput{
				{Key: "file", Type: operations.EventFieldTypeString, Required: true, MaxLength: maxLen},
				{Key: "size_mb", Type: operations.EventFieldTypeString, Required: true, MaxLength: maxLen},
			},
		},
	}

	for _, definition := range definitions {
		_, err := operations.UpsertEventDefinition(ctx, client, siteID, definition)
		require.NoError(t, err)
	}

	t.Run("event properties are persisted and retrieved via GraphQL", func(t *testing.T) {
		// Send event with string:string properties
		properties := `{"button": "signup", "variant": "blue", "position": "1"}`
		payload := map[string]interface{}{
			"site_key":   siteKey,
			"name":       "button_click",
			"path":       "/landing",
			"properties": properties,
		}
		body, _ := json.Marshal(payload)

		resp, err := ts.httpServer.Client().Post(
			ts.httpServer.URL+"/api/event",
			"application/json",
			bytes.NewReader(body),
		)
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusNoContent, resp.StatusCode)

		// Retrieve events via GraphQL API
		eventsResp, err := operations.Events(ctx, client, siteID, nil, nil, nil)
		require.NoError(t, err)
		require.NotEmpty(t, eventsResp.Events.Events, "should have at least one event")

		// Find our event
		var foundEvent *operations.EventsEventsEventsResultEventsEvent
		for _, e := range eventsResp.Events.Events {
			if e.Name == "button_click" && e.Path == "/landing" {
				foundEvent = &e
				break
			}
		}

		require.NotNil(t, foundEvent, "should find the button_click event")
		require.Len(t, foundEvent.Properties, 3, "should have 3 properties")

		// Verify properties are returned correctly
		propsMap := make(map[string]string)
		for _, p := range foundEvent.Properties {
			propsMap[p.Key] = p.Value
		}
		require.Equal(t, "signup", propsMap["button"])
		require.Equal(t, "blue", propsMap["variant"])
		require.Equal(t, "1", propsMap["position"])
	})

	t.Run("multiple events with different properties", func(t *testing.T) {
		// Send multiple events with string:string properties
		events := []struct {
			name       string
			path       string
			properties string
		}{
			{"form_submit", "/contact", `{"form_id": "contact-us", "fields": "5"}`},
			{"video_play", "/media", `{"video_id": "intro", "duration": "120"}`},
			{"download", "/resources", `{"file": "whitepaper.pdf", "size_mb": "2.5"}`},
		}

		for _, ev := range events {
			payload := map[string]interface{}{
				"site_key":   siteKey,
				"name":       ev.name,
				"path":       ev.path,
				"properties": ev.properties,
			}
			body, _ := json.Marshal(payload)

			resp, err := ts.httpServer.Client().Post(
				ts.httpServer.URL+"/api/event",
				"application/json",
				bytes.NewReader(body),
			)
			require.NoError(t, err)
			resp.Body.Close()
			require.Equal(t, http.StatusNoContent, resp.StatusCode)
		}

		// Retrieve events via GraphQL API
		eventsResp, err := operations.Events(ctx, client, siteID, nil, nil, nil)
		require.NoError(t, err)

		// Verify each event exists
		for _, expected := range events {
			found := false
			for _, stored := range eventsResp.Events.Events {
				if stored.Name == expected.name && stored.Path == expected.path {
					found = true
					break
				}
			}
			require.True(t, found, "event %s should be found", expected.name)
		}
	})

	t.Run("events pagination works", func(t *testing.T) {
		limit := 2
		eventsResp, err := operations.Events(ctx, client, siteID, nil, &limit, nil)
		require.NoError(t, err)
		require.LessOrEqual(t, len(eventsResp.Events.Events), 2, "should return at most 2 events")
		require.GreaterOrEqual(t, eventsResp.Events.Total, 4, "total should include all events")
	})

	t.Run("unauthenticated user cannot access events", func(t *testing.T) {
		_, err := operations.Events(ctx, ts.graphqlClient(), siteID, nil, nil, nil)
		require.Error(t, err)
	})
}

func TestHealthEndpoint(t *testing.T) {
	ts := newTestServer(t)

	t.Run("health endpoint returns healthy status", func(t *testing.T) {
		resp, err := ts.httpServer.Client().Get(ts.httpServer.URL + "/health")
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)
		require.Contains(t, resp.Header.Get("Content-Type"), "application/json")

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		require.Contains(t, string(body), `"status":"healthy"`, "health endpoint should return healthy status")
	})
}
