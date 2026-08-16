package http

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/lovely-eye/server/internal/auth"
	"github.com/stretchr/testify/require"
)

func TestAuthCookiesAreDeterministicallyScopedByBasePath(t *testing.T) {
	first := newCookieTestManager("/tools/first/")
	second := newCookieTestManager("/tools/second")
	tokens := &auth.Tokens{AccessToken: "access", RefreshToken: "refresh"}

	firstRecorder := httptest.NewRecorder()
	first.SetAuthCookies(firstRecorder, tokens)
	firstCookies := firstRecorder.Result().Cookies()
	require.Len(t, firstCookies, 2)
	for _, cookie := range firstCookies {
		require.Equal(t, "/tools/first", cookie.Path)
		require.True(t, cookie.HttpOnly)
	}

	secondRecorder := httptest.NewRecorder()
	second.SetAuthCookies(secondRecorder, tokens)
	secondCookies := secondRecorder.Result().Cookies()
	require.Len(t, secondCookies, 2)
	require.NotEqual(t, firstCookies[0].Name, secondCookies[0].Name)
	require.NotEqual(t, firstCookies[1].Name, secondCookies[1].Name)

	request := httptest.NewRequest("POST", "/tools/first/graphql", nil)
	for _, cookie := range firstCookies {
		request.AddCookie(cookie)
	}
	access, refresh := first.TokensFromRequest(request)
	require.Equal(t, tokens.AccessToken, access)
	require.Equal(t, tokens.RefreshToken, refresh)

	clearRecorder := httptest.NewRecorder()
	first.ClearAuthCookies(clearRecorder)
	clearedCookies := clearRecorder.Result().Cookies()
	require.Len(t, clearedCookies, 2)
	for index, cookie := range clearedCookies {
		require.Equal(t, firstCookies[index].Name, cookie.Name)
		require.Equal(t, "/tools/first", cookie.Path)
		require.Equal(t, -1, cookie.MaxAge)
	}
}

func newCookieTestManager(basePath string) *CookieManager {
	return NewCookieManager(CookieConfig{
		AccessTokenExpiry: 15 * time.Minute,
		RefreshExpiry:     7 * 24 * time.Hour,
		BasePath:          basePath,
	})
}
