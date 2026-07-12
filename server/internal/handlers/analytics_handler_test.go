package handlers

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/lovely-eye/server/internal/database"
	"github.com/lovely-eye/server/internal/models"
	"github.com/lovely-eye/server/internal/repository"
	"github.com/lovely-eye/server/internal/services"
	"github.com/lovely-eye/server/pkg/clientip"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"

	_ "modernc.org/sqlite"
)

func TestAnalyticsHandlerCollectRequiresQuerySiteKey(t *testing.T) {
	handler := NewAnalyticsHandler(nil, nil, AnalyticsHandlerConfig{}, nil, nil)
	req := httptest.NewRequest(http.MethodPost, "/api/collect", strings.NewReader(`{"path":"/pricing"}`))
	rec := httptest.NewRecorder()

	handler.Collect(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestAnalyticsHandlerCollectRejectsOversizedBody(t *testing.T) {
	handler, site := newAnalyticsHandlerTestFixture(t, AnalyticsHandlerConfig{
		MaxBodyBytes:       8,
		MaxPropertiesBytes: 1024,
	}, nil)
	req := newAnalyticsCollectRequest(site.PublicKey, `{"path":"/pricing"}`)
	rec := httptest.NewRecorder()

	handler.Collect(rec, req)

	require.Equal(t, http.StatusRequestEntityTooLarge, rec.Code)
}

func TestAnalyticsHandlerCollectRejectsOversizedProperties(t *testing.T) {
	handler, site := newAnalyticsHandlerTestFixture(t, AnalyticsHandlerConfig{
		MaxBodyBytes:       4096,
		MaxPropertiesBytes: 10,
	}, nil)
	req := newAnalyticsCollectRequest(site.PublicKey, `{"path":"/pricing","properties":"{\"code\":\"PAYMENT_DECLINED\"}"}`)
	rec := httptest.NewRecorder()

	handler.Collect(rec, req)

	require.Equal(t, http.StatusRequestEntityTooLarge, rec.Code)
}

func TestAnalyticsHandlerCollectRateLimitsBySiteAndIP(t *testing.T) {
	limiter := NewCollectRateLimiter(true, 1, 1)
	handler, site := newAnalyticsHandlerTestFixture(t, AnalyticsHandlerConfig{
		MaxBodyBytes:       4096,
		MaxPropertiesBytes: 1024,
	}, limiter)

	req1 := newAnalyticsCollectRequest(site.PublicKey, `{"path":"/pricing"}`)
	rec1 := httptest.NewRecorder()
	handler.Collect(rec1, req1)
	require.Equal(t, http.StatusNoContent, rec1.Code)

	req2 := newAnalyticsCollectRequest(site.PublicKey, `{"path":"/pricing"}`)
	rec2 := httptest.NewRecorder()
	handler.Collect(rec2, req2)
	require.Equal(t, http.StatusTooManyRequests, rec2.Code)
}

func TestAnalyticsHandlerCollectRateLimitsRotatedInvalidSiteKeysByIP(t *testing.T) {
	limiter := NewCollectRateLimiter(true, 1, 1)
	handler, _ := newAnalyticsHandlerTestFixture(t, AnalyticsHandlerConfig{
		MaxBodyBytes:       4096,
		MaxPropertiesBytes: 1024,
	}, limiter)

	req1 := newAnalyticsCollectRequest("missing-site-key-a", `{"path":"/pricing"}`)
	rec1 := httptest.NewRecorder()
	handler.Collect(rec1, req1)
	require.Equal(t, http.StatusNoContent, rec1.Code)

	req2 := newAnalyticsCollectRequest("missing-site-key-b", `{"path":"/pricing"}`)
	rec2 := httptest.NewRecorder()
	handler.Collect(rec2, req2)
	require.Equal(t, http.StatusTooManyRequests, rec2.Code)
}

func TestAnalyticsHandlerCollectRejectsCustomEventWithoutPath(t *testing.T) {
	handler, site := newAnalyticsHandlerTestFixture(t, AnalyticsHandlerConfig{
		MaxBodyBytes:       4096,
		MaxPropertiesBytes: 1024,
	}, nil)
	req := newAnalyticsCollectRequest(site.PublicKey, `{"name":"checkout_failed"}`)
	rec := httptest.NewRecorder()

	handler.Collect(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestAnalyticsHandlerCollectRejectsLegacyDurationOnlyPayload(t *testing.T) {
	handler, site := newAnalyticsHandlerTestFixture(t, AnalyticsHandlerConfig{
		MaxBodyBytes:       4096,
		MaxPropertiesBytes: 1024,
	}, nil)
	req := newAnalyticsCollectRequest(site.PublicKey, `{"duration":60}`)
	rec := httptest.NewRecorder()

	handler.Collect(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestAnalyticsHandlerCollectHonorsForwardedIPFromTrustedRemote(t *testing.T) {
	fixture := newAnalyticsHandlerTestFixtureWithResolver(t, AnalyticsHandlerConfig{
		MaxBodyBytes:       4096,
		MaxPropertiesBytes: 1024,
	}, nil, clientip.MustNewResolver([]string{"10.0.0.0/8"}))
	insertAnalyticsHandlerBlockedIP(t, fixture.db, fixture.site.ID, "198.51.100.9")

	req := newAnalyticsCollectRequest(fixture.site.PublicKey, `{"path":"/pricing"}`)
	req.RemoteAddr = "10.1.2.3:12345"
	req.Header.Set("X-Forwarded-For", "198.51.100.9")
	rec := httptest.NewRecorder()

	fixture.handler.Collect(rec, req)

	require.Equal(t, http.StatusNoContent, rec.Code)
	require.Equal(t, 0, countAnalyticsHandlerPageViews(t, fixture.db, fixture.site.ID))
}

func TestAnalyticsHandlerCollectIgnoresForwardedIPFromUntrustedRemote(t *testing.T) {
	fixture := newAnalyticsHandlerTestFixtureWithResolver(t, AnalyticsHandlerConfig{
		MaxBodyBytes:       4096,
		MaxPropertiesBytes: 1024,
	}, nil, clientip.MustNewResolver([]string{"10.0.0.0/8"}))
	insertAnalyticsHandlerBlockedIP(t, fixture.db, fixture.site.ID, "198.51.100.9")

	req := newAnalyticsCollectRequest(fixture.site.PublicKey, `{"path":"/pricing"}`)
	req.RemoteAddr = "203.0.113.10:12345"
	req.Header.Set("X-Forwarded-For", "198.51.100.9")
	rec := httptest.NewRecorder()

	fixture.handler.Collect(rec, req)

	require.Equal(t, http.StatusNoContent, rec.Code)
	require.Equal(t, 1, countAnalyticsHandlerPageViews(t, fixture.db, fixture.site.ID))
}

type analyticsHandlerTestFixture struct {
	handler *AnalyticsHandler
	site    *models.Site
	db      *bun.DB
}

func newAnalyticsHandlerTestFixture(t *testing.T, handlerConfig AnalyticsHandlerConfig, limiter *CollectRateLimiter) (*AnalyticsHandler, *models.Site) {
	t.Helper()

	fixture := newAnalyticsHandlerTestFixtureWithResolver(t, handlerConfig, limiter, clientip.MustNewResolver(nil))
	return fixture.handler, fixture.site
}

func newAnalyticsHandlerTestFixtureWithResolver(
	t *testing.T,
	handlerConfig AnalyticsHandlerConfig,
	limiter *CollectRateLimiter,
	resolver *clientip.Resolver,
) *analyticsHandlerTestFixture {
	t.Helper()

	db := setupAnalyticsHandlerTestDB(t)
	ctx := context.Background()

	user := &models.User{
		Username:     "handler-user",
		PasswordHash: "hash",
		Role:         "admin",
	}
	_, err := db.NewInsert().Model(user).Exec(ctx)
	require.NoError(t, err)

	site := &models.Site{
		UserID:    user.ID,
		Name:      "Handler Site",
		PublicKey: "handler-site-key",
	}
	_, err = db.NewInsert().Model(site).Exec(ctx)
	require.NoError(t, err)

	domain := &models.SiteDomain{
		SiteID:   site.ID,
		Domain:   "handler.test",
		Position: 0,
	}
	_, err = db.NewInsert().Model(domain).Exec(ctx)
	require.NoError(t, err)

	siteRepo := repository.NewSiteRepository(db)
	analyticsRepo := repository.NewAnalyticsRepository(db)
	eventDefinitionRepo := repository.NewEventDefinitionRepository(db)
	siteService := services.NewSiteService(siteRepo)
	analyticsService := services.NewAnalyticsService(analyticsRepo, siteRepo, eventDefinitionRepo, nil, nil, strings.Repeat("a", 32))
	if resolver == nil {
		resolver = clientip.MustNewResolver(nil)
	}

	return &analyticsHandlerTestFixture{
		handler: NewAnalyticsHandler(analyticsService, siteService, handlerConfig, resolver, limiter),
		site:    site,
		db:      db,
	}
}

func setupAnalyticsHandlerTestDB(t *testing.T) *bun.DB {
	t.Helper()

	sqldb, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)

	db := bun.NewDB(sqldb, sqlitedialect.New())
	require.NoError(t, database.Migrate(context.Background(), db))

	t.Cleanup(func() {
		require.NoError(t, db.Close())
	})

	return db
}

func newAnalyticsCollectRequest(siteKey string, body string) *http.Request {
	req := httptest.NewRequest(http.MethodPost, "/api/collect?site_key="+siteKey, strings.NewReader(body))
	req.RemoteAddr = "203.0.113.10:12345"
	req.Header.Set("Origin", "https://handler.test")
	req.Header.Set("Content-Type", "text/plain;charset=UTF-8")
	return req
}

func insertAnalyticsHandlerBlockedIP(t *testing.T, db *bun.DB, siteID int64, ip string) {
	t.Helper()

	_, err := db.NewInsert().Model(&models.SiteBlockedIP{
		SiteID: siteID,
		IP:     ip,
	}).Exec(context.Background())
	require.NoError(t, err)
}

func countAnalyticsHandlerPageViews(t *testing.T, db *bun.DB, siteID int64) int {
	t.Helper()

	count, err := db.NewSelect().
		TableExpr("events e").
		Join("INNER JOIN sessions s ON e.session_id = s.id").
		Where("s.site_id = ?", siteID).
		Where("e.definition_id IS NULL").
		Count(context.Background())
	require.NoError(t, err)
	return count
}
