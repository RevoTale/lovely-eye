package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"testing"

	"github.com/Khan/genqlient/graphql"
	operations "github.com/lovely-eye/server/e2e/generated"
	"github.com/stretchr/testify/require"
)

func TestGraphQLPlaygroundIsNotExposed(t *testing.T) {
	ts := newTestServer(t)

	response, err := ts.httpServer.Client().Get(ts.httpServer.URL + "/graphql")
	require.NoError(t, err)
	defer response.Body.Close()
	require.Equal(t, http.StatusMethodNotAllowed, response.StatusCode)
}

func TestGraphQLRejectsOperationsAboveComplexityLimit(t *testing.T) {
	ts := newTestServer(t)

	var query strings.Builder
	query.WriteString(`{"query":"query {`)
	for i := range 301 {
		fmt.Fprintf(&query, " field%d: __typename", i)
	}
	query.WriteString(` }"}`)

	assertGraphQLError(
		t,
		postGraphQLRaw(t, ts.httpServer.Client(), ts.httpServer.URL+"/graphql", query.String()),
		"BAD_USER_INPUT",
		"operation has complexity 301, which exceeds the limit of 300",
	)
}

func TestGraphQLSemanticErrorCategories(t *testing.T) {
	ts := newTestServer(t)
	ctx := context.Background()

	assertGraphQLError(t, postGraphQLRaw(
		t,
		ts.httpServer.Client(),
		ts.httpServer.URL+"/graphql",
		`{"query":"mutation { createSite(input: { name: \"Unauthorized\", domains: [\"unauthorized.example\"] }) { id } }"}`,
	), "UNAUTHENTICATED", "unauthorized")

	_, err := operations.Register(ctx, ts.graphqlClient(), operations.RegisterInput{
		Username: "contract-admin",
		Password: "password123",
	})
	require.NoError(t, err)
	jar, err := cookiejar.New(nil)
	require.NoError(t, err)
	httpClient := &http.Client{Jar: jar}
	client := graphql.NewClient(ts.httpServer.URL+"/graphql", httpClient)
	_, err = operations.Login(ctx, client, operations.LoginInput{
		Username: "contract-admin",
		Password: "password123",
	})
	require.NoError(t, err)

	assertGraphQLError(t, postGraphQLRaw(
		t,
		httpClient,
		ts.httpServer.URL+"/graphql",
		`{"query":"mutation { createSite(input: { name: \"Invalid Input\", domains: [\"not a domain\"] }) { id } }"}`,
	), "BAD_USER_INPUT", "failed to validate domain")
	assertGraphQLError(t, postGraphQLRaw(
		t,
		httpClient,
		ts.httpServer.URL+"/graphql",
		`{"query":"query { site(id: \"999999\") { id } }"}`,
	), "NOT_FOUND", "site not found")
	created, err := operations.CreateSite(ctx, client, operations.CreateSiteInput{
		Domains: []string{"duplicate.example"},
		Name:    "Original",
	})
	require.NoError(t, err)
	assertGraphQLError(t, postGraphQLRaw(
		t,
		httpClient,
		ts.httpServer.URL+"/graphql",
		fmt.Sprintf(
			`{"query":"query { dashboard(siteId: \"%s\", filter: { eventDefinitionId: [\"invalid\"] }) { visitors } }"}`,
			created.CreateSite.Id,
		),
	), "BAD_USER_INPUT", "invalid event definition ID")
	assertGraphQLError(t, postGraphQLRaw(
		t,
		httpClient,
		ts.httpServer.URL+"/graphql",
		`{"query":"mutation { createSite(input: { name: \"Duplicate\", domains: [\"duplicate.example\"] }) { id } }"}`,
	), "CONFLICT", "site with this domain already exists")
}

func assertGraphQLError(t *testing.T, response rawGraphQLResponse, code, message string) {
	t.Helper()
	var payload struct {
		Errors []struct {
			Message    string `json:"message"`
			Extensions struct {
				Code string `json:"code"`
			} `json:"extensions"`
		} `json:"errors"`
	}
	require.NoError(t, json.Unmarshal([]byte(response.body), &payload))
	require.Len(t, payload.Errors, 1, response.body)
	require.Equal(t, code, payload.Errors[0].Extensions.Code)
	require.Contains(t, payload.Errors[0].Message, message)
}

func TestAuthResponsesUsePersistedCreationTimestamp(t *testing.T) {
	ts := newTestServer(t)
	ctx := context.Background()
	credentials := operations.RegisterInput{
		Username: "timestamp-admin",
		Password: "password123",
	}

	registered, err := operations.Register(ctx, ts.graphqlClient(), credentials)
	require.NoError(t, err)
	require.False(t, registered.Register.User.CreatedAt.IsZero())

	client := ts.authenticatedClient(ctx, t, credentials.Username, credentials.Password)
	loggedIn, err := operations.Login(ctx, client, operations.LoginInput(credentials))
	require.NoError(t, err)
	current, err := operations.Me(ctx, client)
	require.NoError(t, err)
	require.NotNil(t, current.Me)
	require.NotNil(t, current.Me.CreatedAt)
	require.Equal(t, registered.Register.User.CreatedAt, loggedIn.Login.User.CreatedAt)
	require.Equal(t, registered.Register.User.CreatedAt, *current.Me.CreatedAt)
}

func TestGeoIPRefreshReturnsActionableFailureStatus(t *testing.T) {
	ts := newTestServer(t)
	ctx := context.Background()
	_, err := operations.Register(ctx, ts.graphqlClient(), operations.RegisterInput{
		Username: "geoip-admin",
		Password: "password123",
	})
	require.NoError(t, err)
	jar, err := cookiejar.New(nil)
	require.NoError(t, err)
	httpClient := &http.Client{Jar: jar}
	client := graphql.NewClient(ts.httpServer.URL+"/graphql", httpClient)
	_, err = operations.Login(ctx, client, operations.LoginInput{
		Username: "geoip-admin",
		Password: "password123",
	})
	require.NoError(t, err)
	site, err := operations.CreateSite(ctx, client, operations.CreateSiteInput{
		Domains: []string{"geoip.example"},
		Name:    "GeoIP Site",
	})
	require.NoError(t, err)

	updateBody, err := json.Marshal(map[string]any{
		"query": `mutation UpdateSite($id: ID!, $input: UpdateSiteInput!) {
			updateSite(id: $id, input: $input) { id }
		}`,
		"variables": map[string]any{
			"id":    site.CreateSite.Id,
			"input": map[string]any{"name": "GeoIP Site", "trackCountry": true},
		},
	})
	require.NoError(t, err)
	updated := postGraphQLRaw(t, httpClient, ts.httpServer.URL+"/graphql", string(updateBody))
	require.NotContains(t, updated.body, `"errors"`)

	refreshed := postGraphQLRaw(
		t,
		httpClient,
		ts.httpServer.URL+"/graphql",
		`{"query":"mutation { refreshGeoIPDatabase { state lastError } }"}`,
	)
	var response struct {
		Data struct {
			Refresh struct {
				State     string  `json:"state"`
				LastError *string `json:"lastError"`
			} `json:"refreshGeoIPDatabase"`
		} `json:"data"`
		Errors []json.RawMessage `json:"errors"`
	}
	require.NoError(t, json.Unmarshal([]byte(refreshed.body), &response))
	require.Empty(t, response.Errors)
	require.Equal(t, "MISSING", response.Data.Refresh.State)
	require.NotNil(t, response.Data.Refresh.LastError)
	require.Contains(t, *response.Data.Refresh.LastError, "not configured")
}

func TestMultiDomainCollectionAggregatesAtSiteLevel(t *testing.T) {
	ts := newTestServer(t)
	ctx := context.Background()

	_, err := operations.Register(ctx, ts.graphqlClient(), operations.RegisterInput{
		Username: "multi-domain-admin",
		Password: "password123",
	})
	require.NoError(t, err)
	client := ts.authenticatedClient(ctx, t, "multi-domain-admin", "password123")

	siteResponse, err := operations.CreateSite(ctx, client, operations.CreateSiteInput{
		Domains: []string{"primary.example", "secondary.example"},
		Name:    "Multi-domain Site",
	})
	require.NoError(t, err)
	require.Equal(t, []string{"primary.example", "secondary.example"}, siteResponse.CreateSite.Domains)

	for _, event := range []struct {
		origin string
		path   string
	}{
		{origin: "https://primary.example", path: "/from-primary"},
		{origin: "https://secondary.example", path: "/from-secondary"},
	} {
		postPageView(t, ts.httpServer.Client(), collectURL(ts.httpServer.URL, siteResponse.CreateSite.PublicKey), event.origin, event.path)
	}

	dashboard, err := operations.Dashboard(
		ctx,
		client,
		siteResponse.CreateSite.Id,
		nil,
		nil,
		defaultPaging,
		defaultPaging,
		defaultPaging,
		defaultPaging,
		defaultPaging,
		defaultPaging,
		nil,
		defaultPaging,
	)
	require.NoError(t, err)
	require.Equal(t, 2, dashboard.Dashboard.PageViews)
	require.Equal(t, 1, dashboard.Dashboard.Visitors)
}

func postPageView(t *testing.T, client *http.Client, endpoint, origin, path string) {
	t.Helper()
	body, err := json.Marshal(map[string]string{"path": path})
	require.NoError(t, err)
	request, err := http.NewRequestWithContext(context.Background(), http.MethodPost, endpoint, bytes.NewReader(body))
	require.NoError(t, err)
	request.Header.Set("Content-Type", "text/plain;charset=UTF-8")
	request.Header.Set("Origin", origin)
	request.Header.Set("User-Agent", "Mozilla/5.0")
	response, err := client.Do(request)
	require.NoError(t, err)
	require.NoError(t, response.Body.Close())
	require.Equal(t, http.StatusNoContent, response.StatusCode)
}
