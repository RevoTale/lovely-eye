package e2e

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/cookiejar"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAuthFlowWithCookies tests the complete authentication flow
// that simulates what happens in the dashboard:
// 1. User logs in via GraphQL mutation
// 2. Server sets cookies
// 3. Page reloads - ME query runs with cookies
// 4. User should still be authenticated
func TestAuthFlowWithCookies(t *testing.T) {
	ts := newTestServer(t)

	// Create HTTP client with cookie jar (simulates browser)
	jar, err := cookiejar.New(nil)
	require.NoError(t, err)
	client := &http.Client{Jar: jar}

	// STEP 1: Register a user
	t.Log("STEP 1: Register user")
	registerMutation := `{
		"query": "mutation Register($input: RegisterInput!) { register(input: $input) { user { id username role } } }",
		"variables": {"input": {"username": "testuser", "password": "testpass123"}}
	}`

	registerResp := mustGraphQL(t, client, ts.httpServer.URL+"/graphql", registerMutation)
	t.Logf("Register response: %s", registerResp)

	// STEP 2: Login
	t.Log("\nSTEP 2: User logs in via mutation")
	loginMutation := `{
		"query": "mutation Login($input: LoginInput!) { login(input: $input) { user { id username role } } }",
		"variables": {"input": {"username": "testuser", "password": "testpass123"}}
	}`

	loginResp := mustGraphQL(t, client, ts.httpServer.URL+"/graphql", loginMutation)
	t.Logf("Login response: %s", loginResp)

	// Check cookies were set
	testURL := ts.httpServer.URL
	parsedURL, _ := http.NewRequest("GET", testURL, nil)
	cookies := jar.Cookies(parsedURL.URL)
	t.Logf("Cookies after login: %d cookies", len(cookies))
	for _, cookie := range cookies {
		maxLen := 20
		if len(cookie.Value) < maxLen {
			maxLen = len(cookie.Value)
		}
		t.Logf("  - %s: %s...", cookie.Name, cookie.Value[:maxLen])
	}

	require.GreaterOrEqual(t, len(cookies), 2, "Expected at least 2 cookies (access and refresh)")

	// NOTE: No CSRF tokens needed! Modern auth uses HttpOnly + Secure cookies with SameSite
	// See https://www.reddit.com/r/node/comments/1im7yj0/comment/mc0ylfd/

	// STEP 3: Simulate page reload - ME query with cookies
	t.Log("\nSTEP 3: Page reloads, ME query runs (cookies auto-included)")
	meQuery := `{"query": "query { me { id username role } }"}`

	// Create request with cookies (client will auto-add them)
	req, err := http.NewRequest(http.MethodPost, ts.httpServer.URL+"/graphql", bytes.NewBufferString(meQuery))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", ts.httpServer.URL)

	// Make request (cookies will be auto-attached by jar)
	meResp, err := client.Do(req)
	require.NoError(t, err)
	defer func ()  {
		err := meResp.Body.Close()
		if nil != err {
			slog.Error("body close err","error",err)
		}
	}()

	meBody, _ := io.ReadAll(meResp.Body)
	t.Logf("ME query status: %d", meResp.StatusCode)
	t.Logf("ME query response: %s", string(meBody))

	// Parse response
	var meResponse struct {
		Data struct {
			Me *struct {
				ID       string `json:"id"`
				Username string `json:"username"`
				Role     string `json:"role"`
			} `json:"me"`
		} `json:"data"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}
	require.NoError(t, json.Unmarshal(meBody, &meResponse))
	require.Empty(t, meResponse.Errors, "ME query should not return errors")
	require.NotNil(t, meResponse.Data.Me, "User should be authenticated after page reload")
	assert.Equal(t, "testuser", meResponse.Data.Me.Username)

	t.Log("\n✅ SUCCESS: User remains authenticated after page reload")

	// STEP 4: Test authenticated mutation works
	t.Log("\nSTEP 4: Test authenticated mutation works")
	createSiteMutation := `{
		"query": "mutation CreateSite($input: CreateSiteInput!) { createSite(input: $input) { id domains } }",
		"variables": {"input": {"domains": ["example.com"], "name": "Example Site"}}
	}`

	req2, err := http.NewRequest(http.MethodPost, ts.httpServer.URL+"/graphql", bytes.NewBufferString(createSiteMutation))
	require.NoError(t, err)
	req2.Header.Set("Content-Type", "application/json")

	resp2, err := client.Do(req2)
	require.NoError(t, err)
	body2, _ := io.ReadAll(resp2.Body)
	err = resp2.Body.Close()
	require.NoError(t,err)

	t.Logf("Mutation status: %d, body: %s", resp2.StatusCode, string(body2))
	assert.Equal(t, http.StatusOK, resp2.StatusCode, "Authenticated mutation should succeed")

	t.Log("\n✅ All authentication flow tests passed!")
}

func mustGraphQL(t *testing.T, client *http.Client, url, query string) string {
	t.Helper()

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBufferString(query))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer func ()  {
	err:=	resp.Body.Close()
	if nil != err {
		slog.Error("gql resp close failed","error",err)
	}
	}()

	body, _ := io.ReadAll(resp.Body)
	require.Equal(t, http.StatusOK, resp.StatusCode, "GraphQL request should succeed: %s", string(body))
	require.NotContains(t, string(body), `"errors"`, "GraphQL should not return errors")

	return string(body)
}
