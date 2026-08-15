package collect

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/lovely-eye/server/internal/analytics"
	analyticspersistence "github.com/lovely-eye/server/internal/analytics/persistence"
	authpersistence "github.com/lovely-eye/server/internal/auth/persistence"
	eventpersistence "github.com/lovely-eye/server/internal/event/persistence"
	"github.com/lovely-eye/server/internal/platform/database"
	sitefeature "github.com/lovely-eye/server/internal/site"
	sitepersistence "github.com/lovely-eye/server/internal/site/persistence"
	"github.com/lovely-eye/server/internal/transport/http/clientip"
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
	limiter := NewRateLimiter(true, 1, 1)
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
	limiter := NewRateLimiter(true, 1, 1)
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

func TestAnalyticsHandlerCollectRejectsTrailingJSON(t *testing.T) {
	handler, site := newAnalyticsHandlerTestFixture(t, AnalyticsHandlerConfig{
		MaxBodyBytes:       4096,
		MaxPropertiesBytes: 1024,
	}, nil)
	req := newAnalyticsCollectRequest(site.PublicKey, `{"path":"/pricing"}{"path":"/second"}`)
	rec := httptest.NewRecorder()

	handler.Collect(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestAnalyticsHandlerCollectRejectsValuesBeyondPersistenceLimits(t *testing.T) {
	handler, site := newAnalyticsHandlerTestFixture(t, AnalyticsHandlerConfig{
		MaxBodyBytes:       8192,
		MaxPropertiesBytes: 1024,
	}, nil)
	tests := []string{
		`{"path":"` + strings.Repeat("p", 2049) + `"}`,
		`{"path":"/","referrer":"` + strings.Repeat("r", 2049) + `"}`,
		`{"path":"/","utm_source":"` + strings.Repeat("s", 129) + `"}`,
		`{"path":"/","utm_medium":"` + strings.Repeat("m", 129) + `"}`,
		`{"path":"/","utm_campaign":"` + strings.Repeat("c", 257) + `"}`,
	}
	for _, body := range tests {
		recorder := httptest.NewRecorder()
		handler.Collect(recorder, newAnalyticsCollectRequest(site.PublicKey, body))
		require.Equal(t, http.StatusBadRequest, recorder.Code)
	}
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

func TestAnalyticsHandlerCollectLoadsSiteOnce(t *testing.T) {
	fixture := newAnalyticsHandlerTestFixtureWithResolver(t, AnalyticsHandlerConfig{
		MaxBodyBytes:       4096,
		MaxPropertiesBytes: 1024,
	}, nil, clientip.MustNewResolver(nil))
	counter := new(publicKeyQueryCounter)
	fixture.db.AddQueryHook(counter)

	recorder := httptest.NewRecorder()
	fixture.handler.Collect(recorder, newAnalyticsCollectRequest(fixture.site.PublicKey, `{"path":"/pricing"}`))

	require.Equal(t, http.StatusNoContent, recorder.Code)
	require.Equal(t, 1, counter.queries)
}

func TestAnalyticsHandlerCollectAllocationBudget(t *testing.T) {
	handler, site := newAnalyticsHandlerTestFixture(t, AnalyticsHandlerConfig{
		MaxBodyBytes:       4096,
		MaxPropertiesBytes: 1024,
	}, nil)

	allocations := testing.AllocsPerRun(20, func() {
		recorder := httptest.NewRecorder()
		handler.Collect(recorder, newAnalyticsCollectRequest(site.PublicKey, `{"path":"/pricing"}`))
		if recorder.Code != http.StatusNoContent {
			t.Fatalf("collect status = %d, want %d", recorder.Code, http.StatusNoContent)
		}
	})

	require.LessOrEqual(t, allocations, 625.0)
}

type publicKeyQueryCounter struct {
	queries int
}

func (c *publicKeyQueryCounter) BeforeQuery(ctx context.Context, event *bun.QueryEvent) context.Context {
	if strings.Contains(event.Query, "public_key") {
		c.queries++
	}
	return ctx
}

func (*publicKeyQueryCounter) AfterQuery(context.Context, *bun.QueryEvent) {}

type analyticsHandlerTestFixture struct {
	handler *AnalyticsHandler
	site    *sitepersistence.Site
	db      *bun.DB
}

func newAnalyticsHandlerTestFixture(t testing.TB, handlerConfig AnalyticsHandlerConfig, limiter *RateLimiter) (*AnalyticsHandler, *sitepersistence.Site) {
	t.Helper()

	fixture := newAnalyticsHandlerTestFixtureWithResolver(t, handlerConfig, limiter, clientip.MustNewResolver(nil))
	return fixture.handler, fixture.site
}

func newAnalyticsHandlerTestFixtureWithResolver(
	t testing.TB,
	handlerConfig AnalyticsHandlerConfig,
	limiter *RateLimiter,
	resolver *clientip.Resolver,
) *analyticsHandlerTestFixture {
	t.Helper()

	db := setupAnalyticsHandlerTestDB(t)
	ctx := context.Background()

	user := &authpersistence.User{
		Username:     "handler-user",
		PasswordHash: "hash",
		Role:         "admin",
	}
	_, err := db.NewInsert().Model(user).Exec(ctx)
	require.NoError(t, err)

	site := &sitepersistence.Site{
		UserID:    user.ID,
		Name:      "Handler Site",
		PublicKey: "handler-site-key",
	}
	_, err = db.NewInsert().Model(site).Exec(ctx)
	require.NoError(t, err)

	domain := &sitepersistence.Domain{
		SiteID:   site.ID,
		Domain:   "handler.test",
		Position: 0,
	}
	_, err = db.NewInsert().Model(domain).Exec(ctx)
	require.NoError(t, err)

	siteRepo := sitepersistence.New(db)
	analyticsRepo := analyticspersistence.New(db)
	eventDefinitionRepo := eventpersistence.New(db)
	siteService := sitefeature.NewService(siteRepo)
	analyticsService := analytics.NewService(analyticsRepo, siteRepo, eventDefinitionRepo, nil, nil, strings.Repeat("a", 32))
	if resolver == nil {
		resolver = clientip.MustNewResolver(nil)
	}

	return &analyticsHandlerTestFixture{
		handler: NewAnalyticsHandler(analyticsService, siteService, handlerConfig, resolver, limiter),
		site:    site,
		db:      db,
	}
}

func setupAnalyticsHandlerTestDB(t testing.TB) *bun.DB {
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

	_, err := db.NewInsert().Model(&sitepersistence.BlockedIP{
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
