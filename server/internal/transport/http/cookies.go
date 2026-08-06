package http

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strings"
	"time"

	"github.com/lovely-eye/server/internal/auth"
)

const (
	accessTokenCookiePrefix  = "le_access"
	refreshTokenCookiePrefix = "le_refresh"
)

type CookieConfig struct {
	AccessTokenExpiry time.Duration
	RefreshExpiry     time.Duration
	Secure            bool
	Domain            string
	BasePath          string
}

type CookieManager struct {
	secure            bool
	domain            string
	path              string
	accessCookieName  string
	refreshCookieName string
	accessExpiry      time.Duration
	refreshExpiry     time.Duration
}

func NewCookieManager(cfg CookieConfig) *CookieManager {
	path := normalizeCookiePath(cfg.BasePath)
	return &CookieManager{
		secure:            cfg.Secure,
		domain:            cfg.Domain,
		path:              path,
		accessCookieName:  scopedCookieName(accessTokenCookiePrefix, path),
		refreshCookieName: scopedCookieName(refreshTokenCookiePrefix, path),
		accessExpiry:      cfg.AccessTokenExpiry,
		refreshExpiry:     cfg.RefreshExpiry,
	}
}

func (m *CookieManager) SetAuthCookies(w http.ResponseWriter, tokens *auth.Tokens) {
	m.setCookie(w, m.accessCookieName, tokens.AccessToken, int(m.accessExpiry.Seconds()))
	m.setCookie(w, m.refreshCookieName, tokens.RefreshToken, int(m.refreshExpiry.Seconds()))
}

func (m *CookieManager) ClearAuthCookies(w http.ResponseWriter) {
	m.setCookie(w, m.accessCookieName, "", -1)
	m.setCookie(w, m.refreshCookieName, "", -1)
}

func (m *CookieManager) TokensFromRequest(r *http.Request) (accessToken, refreshToken string) {
	if cookie, err := r.Cookie(m.accessCookieName); err == nil {
		accessToken = cookie.Value
	}
	if cookie, err := r.Cookie(m.refreshCookieName); err == nil {
		refreshToken = cookie.Value
	}
	return accessToken, refreshToken
}

func (m *CookieManager) setCookie(w http.ResponseWriter, name, value string, maxAge int) {
	sameSite := http.SameSiteLaxMode
	if m.secure {
		sameSite = http.SameSiteStrictMode
	}
	// #nosec G124 -- Secure is configurable for local HTTP/test; production defaults it to true.
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     m.path,
		Domain:   m.domain,
		MaxAge:   maxAge,
		HttpOnly: true,
		Secure:   m.secure,
		SameSite: sameSite,
	})
}

func normalizeCookiePath(basePath string) string {
	trimmed := strings.Trim(basePath, "/")
	if trimmed == "" {
		return "/"
	}
	return "/" + trimmed
}

func scopedCookieName(prefix, cookiePath string) string {
	// Hashing keeps arbitrary nested paths within the conservative cookie-name character set.
	sum := sha256.Sum256([]byte(cookiePath))
	return prefix + "_" + hex.EncodeToString(sum[:4])
}
