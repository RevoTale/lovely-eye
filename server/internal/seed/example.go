// Package seed loads reusable development data into a Lovely Eye database.
package seed

import (
	"context"
	cryptorand "crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	mathrand "math/rand"
	"strings"
	"time"

	analyticspersistence "github.com/lovely-eye/server/internal/analytics/persistence"
	authpersistence "github.com/lovely-eye/server/internal/auth/persistence"
	"github.com/lovely-eye/server/internal/event"
	eventpersistence "github.com/lovely-eye/server/internal/event/persistence"
	"github.com/lovely-eye/server/internal/platform/config"
	"github.com/lovely-eye/server/internal/platform/database"
	sitefeature "github.com/lovely-eye/server/internal/site"
	sitepersistence "github.com/lovely-eye/server/internal/site/persistence"
	"github.com/uptrace/bun"
	"golang.org/x/crypto/bcrypt"
)

const (
	defaultUsers     = 80
	minSessions      = 1
	maxSessions      = 3
	maxDays          = 14
	recentSessions   = 12
	sessionBaseMins  = 4
	sessionExtraMins = 10
)

var (
	pathsMarketing = []string{"/", "/pricing", "/features", "/blog/launch"}
	pathsDocs      = []string{"/", "/docs", "/docs/setup", "/docs/api"}
	pathsProduct   = []string{"/", "/app", "/app/dashboard", "/settings", "/billing"}
	pathsUpgrade   = []string{"/", "/pricing", "/checkout", "/billing"}
	referrers      = []string{"https://google.com", "https://news.ycombinator.com", "https://github.com"}
)

type eventSeed struct {
	name  string
	path  string
	props map[string]string
}

type behaviorPattern struct {
	paths    []string
	referrer string
	events   []eventSeed
}

// Result summarizes one example-data load.
type Result struct {
	SiteName         string
	PublicKey        string
	CreatedSite      bool
	Clients          int
	Sessions         int
	PageViews        int
	PredefinedEvents int
}

// Run opens and migrates the configured database before loading example data.
func Run(ctx context.Context, cfg config.Config) (_ Result, err error) {
	db, err := database.New(ctx, cfg.Database)
	if err != nil {
		return Result{}, fmt.Errorf("connect to database: %w", err)
	}
	defer func() {
		err = errors.Join(err, database.Close(db))
	}()

	if err := database.Migrate(ctx, db); err != nil {
		return Result{}, fmt.Errorf("migrate database: %w", err)
	}
	return Load(ctx, db, cfg.Auth)
}

// Load adds example data to an already-open database.
func Load(ctx context.Context, db *bun.DB, authCfg config.AuthConfig) (Result, error) {
	owner, err := ensureSeedOwner(ctx, db, authCfg)
	if err != nil {
		return Result{}, fmt.Errorf("resolve seed user: %w", err)
	}

	site, created, err := ensureLocalhostSite(ctx, db, owner.ID)
	if err != nil {
		return Result{}, fmt.Errorf("ensure localhost site: %w", err)
	}

	defs, err := ensureEventDefinitions(ctx, db, site.ID)
	if err != nil {
		return Result{}, fmt.Errorf("ensure event definitions: %w", err)
	}

	counts, err := seedData(ctx, db, site.ID, defs)
	if err != nil {
		return Result{}, fmt.Errorf("seed analytics data: %w", err)
	}

	return Result{
		SiteName:         site.Name,
		PublicKey:        site.PublicKey,
		CreatedSite:      created,
		Clients:          counts.clients,
		Sessions:         counts.sessions,
		PageViews:        counts.pageViews,
		PredefinedEvents: counts.predefinedEvents,
	}, nil
}

type seedCounts struct {
	clients          int
	sessions         int
	pageViews        int
	predefinedEvents int
}

func firstUser(ctx context.Context, db *bun.DB) (*authpersistence.User, error) {
	user := new(authpersistence.User)
	if err := db.NewSelect().
		Model(user).
		Order("id ASC").
		Limit(1).
		Scan(ctx); err != nil {
		return nil, fmt.Errorf("scan first user: %w", err)
	}
	return user, nil
}

func ensureSeedOwner(ctx context.Context, db *bun.DB, authCfg config.AuthConfig) (*authpersistence.User, error) {
	user, err := firstUser(ctx, db)
	if err == nil {
		return user, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	username := authCfg.InitialAdminUsername
	password := authCfg.InitialAdminPassword
	if username == "" {
		username = "demo-admin"
	}
	if password == "" {
		password = "demo-password"
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash seed password: %w", err)
	}

	userRepo := authpersistence.New(db)
	newUser := &authpersistence.User{
		Username:     username,
		PasswordHash: string(hash),
		Role:         "admin",
	}
	if err := userRepo.Create(ctx, newUser); err != nil {
		return nil, fmt.Errorf("create seed user: %w", err)
	}
	return newUser, nil
}

func ensureLocalhostSite(ctx context.Context, db *bun.DB, userID int64) (*sitefeature.Site, bool, error) {
	siteRepo := sitepersistence.New(db)
	siteService := sitefeature.NewService(siteRepo)

	site, err := siteRepo.GetByDomainForUser(ctx, userID, "localhost")
	created := false
	if err != nil {
		if !errors.Is(err, sitefeature.ErrSiteNotFound) {
			return nil, false, fmt.Errorf("get site by domain for user: %w", err)
		}
		site, err = siteService.Create(ctx, sitefeature.CreateSiteInput{
			Domains: []string{"localhost"},
			Name:    "Localhost",
			UserID:  userID,
		})
		if err != nil {
			return nil, false, fmt.Errorf("create localhost site: %w", err)
		}
		created = true
	}

	return site, created, nil
}

func ensureEventDefinitions(ctx context.Context, db *bun.DB, siteID int64) ([]*event.Definition, error) {
	eventRepo := eventpersistence.New(db)
	eventService := event.NewService(eventRepo)

	definitions := []event.DefinitionInput{
		{
			Name: "signup",
			Fields: []event.FieldInput{
				{Key: "plan", Type: "string", Required: true},
				{Key: "referrer", Type: "string"},
			},
		},
		{
			Name: "purchase",
			Fields: []event.FieldInput{
				{Key: "amount", Type: "int", Required: true},
				{Key: "currency", Type: "string", Required: true},
				{Key: "plan", Type: "string"},
			},
		},
		{
			Name: "newsletter_subscribe",
			Fields: []event.FieldInput{
				{Key: "source", Type: "string"},
			},
		},
		{
			Name: "video_play",
			Fields: []event.FieldInput{
				{Key: "video", Type: "string", Required: true},
				{Key: "seconds", Type: "int"},
			},
		},
		{
			Name: "cta_click",
			Fields: []event.FieldInput{
				{Key: "cta", Type: "string", Required: true},
				{Key: "page", Type: "string"},
			},
		},
		{
			Name: "file_download",
			Fields: []event.FieldInput{
				{Key: "file", Type: "string", Required: true},
				{Key: "success", Type: "bool"},
			},
		},
	}

	results := make([]*event.Definition, 0, len(definitions))
	for _, def := range definitions {
		created, err := eventService.Upsert(ctx, siteID, def)
		if err != nil {
			return nil, fmt.Errorf("upsert event definition %q: %w", def.Name, err)
		}
		results = append(results, created)
	}
	return results, nil
}

func seedData(ctx context.Context, db *bun.DB, siteID int64, defs []*event.Definition) (seedCounts, error) {
	counts := seedCounts{}
	analyticsRepo := analyticspersistence.New(db)
	//nolint:gosec // Deterministic randomness isn't required for seed data.
	rng := mathrand.New(mathrand.NewSource(time.Now().UnixNano()))

	defMap := make(map[string]*event.Definition, len(defs))
	for _, def := range defs {
		defMap[def.Name] = def
	}

	patterns := []behaviorPattern{
		{
			paths:    pathsMarketing,
			referrer: pickString(rng, referrers),
			events: []eventSeed{
				{name: "signup", path: "/pricing", props: map[string]string{"plan": "pro", "referrer": "pricing"}},
				{name: "cta_click", path: "/pricing", props: map[string]string{"cta": "Start Trial", "page": "/pricing"}},
			},
		},
		{
			paths:    pathsDocs,
			referrer: pickString(rng, referrers),
			events: []eventSeed{
				{name: "video_play", path: "/docs", props: map[string]string{"video": "setup", "seconds": "120"}},
				{name: "newsletter_subscribe", path: "/blog/launch", props: map[string]string{"source": "docs"}},
			},
		},
		{
			paths:    pathsProduct,
			referrer: "",
			events: []eventSeed{
				{name: "file_download", path: "/app", props: map[string]string{"file": "report.pdf", "success": "true"}},
			},
		},
		{
			paths:    pathsUpgrade,
			referrer: pickString(rng, referrers),
			events: []eventSeed{
				{name: "purchase", path: "/billing", props: map[string]string{"amount": "49", "currency": "USD", "plan": "business"}},
			},
		},
	}

	now := time.Now()
	recentRemaining := recentSessions

	for range defaultUsers {
		hash, err := generateRandomString(64)
		if err != nil {
			return counts, fmt.Errorf("generate client hash: %w", err)
		}

		client := &analyticspersistence.Client{
			SiteID:  siteID,
			Hash:    hash,
			Country: pickString(rng, []string{"US", "GB", "DE", "FR", "CA", "NL"}),
			Device: pickEnum(rng, []analyticspersistence.ClientDevice{
				analyticspersistence.ClientDeviceDesktop,
				analyticspersistence.ClientDeviceMobile,
				analyticspersistence.ClientDeviceTablet,
				analyticspersistence.ClientDeviceSmartTV,
				analyticspersistence.ClientDeviceConsole,
				analyticspersistence.ClientDeviceWatch,
			}),
			Browser: pickEnum(rng, []analyticspersistence.ClientBrowser{
				analyticspersistence.ClientBrowserChrome,
				analyticspersistence.ClientBrowserSafari,
				analyticspersistence.ClientBrowserFirefox,
				analyticspersistence.ClientBrowserEdge,
				analyticspersistence.ClientBrowserOpera,
				analyticspersistence.ClientBrowserSamsungInternet,
				analyticspersistence.ClientBrowserAndroidWebView,
			}),
			OS: pickEnum(rng, []analyticspersistence.ClientOS{
				analyticspersistence.ClientOSWindows,
				analyticspersistence.ClientOSMacOS,
				analyticspersistence.ClientOSLinux,
				analyticspersistence.ClientOSChromeOS,
				analyticspersistence.ClientOSIOS,
				analyticspersistence.ClientOSIPadOS,
				analyticspersistence.ClientOSAndroid,
				analyticspersistence.ClientOSWatchOS,
				analyticspersistence.ClientOSWearOS,
			}),
			ScreenSize: pickEnum(rng, []analyticspersistence.ClientScreenSize{
				analyticspersistence.ClientScreenSizeWatch,
				analyticspersistence.ClientScreenSizeXS,
				analyticspersistence.ClientScreenSizeSM,
				analyticspersistence.ClientScreenSizeLG,
				analyticspersistence.ClientScreenSizeXL,
			}),
		}

		if _, err := db.NewInsert().Model(client).Exec(ctx); err != nil {
			return counts, fmt.Errorf("insert client: %w", err)
		}
		counts.clients++

		sessionCount := randRange(rng, minSessions, maxSessions)
		for j := 0; j < sessionCount; j++ {
			pattern := patterns[rng.Intn(len(patterns))]
			start := randomStart(rng, now, &recentRemaining)
			created, err := applyPattern(ctx, analyticsRepo, rng, client.ID, siteID, start, pattern, defMap)
			if err != nil {
				return counts, fmt.Errorf("apply pattern: %w", err)
			}
			counts.sessions++
			counts.pageViews += created.pageViews
			counts.predefinedEvents += created.predefinedEvents
		}
	}

	return counts, nil
}

type patternCounts struct {
	pageViews        int
	predefinedEvents int
}

func applyPattern(
	ctx context.Context,
	analyticsRepo *analyticspersistence.Repository,
	rng *mathrand.Rand,
	clientID, siteID int64,
	start time.Time,
	pattern behaviorPattern,
	defMap map[string]*event.Definition,
) (patternCounts, error) {
	counts := patternCounts{}

	paths := normalizePaths(pattern.paths)
	durationMinutes := randRange(rng, sessionBaseMins, sessionBaseMins+sessionExtraMins)
	duration := durationMinutes*60 + randRange(rng, 0, 120)
	enterUnix := start.Unix()
	exitUnix := enterUnix + int64(duration)

	session := &analyticspersistence.Session{
		SiteID:        siteID,
		ClientID:      clientID,
		EnterTime:     enterUnix,
		EnterHour:     enterUnix / 3600,
		EnterDay:      enterUnix / 86400,
		EnterPath:     paths[0],
		ExitTime:      exitUnix,
		ExitHour:      exitUnix / 3600,
		ExitDay:       exitUnix / 86400,
		ExitPath:      paths[len(paths)-1],
		Referrer:      pattern.referrer,
		UTMSource:     "",
		UTMMedium:     "",
		UTMCampaign:   "",
		Duration:      duration,
		PageViewCount: len(paths),
	}
	if err := analyticsRepo.CreateSession(ctx, session); err != nil {
		return counts, fmt.Errorf("create session: %w", err)
	}

	interval := duration / maxInt(len(paths), 1)
	for idx, path := range paths {
		eventTime := enterUnix + int64(interval*idx)
		event := &analyticspersistence.Event{
			SessionID: session.ID,
			Time:      eventTime,
			Hour:      eventTime / 3600,
			Day:       eventTime / 86400,
			Path:      path,
		}
		if err := analyticsRepo.CreateEvent(ctx, event); err != nil {
			return counts, fmt.Errorf("create page view event: %w", err)
		}
		counts.pageViews++
	}

	for _, seed := range pattern.events {
		def := defMap[seed.name]
		if def == nil {
			continue
		}
		eventTime := enterUnix + int64(interval/2)
		defID := def.ID
		event := &analyticspersistence.Event{
			SessionID:    session.ID,
			Time:         eventTime,
			Hour:         eventTime / 3600,
			Day:          eventTime / 86400,
			Path:         normalizeURL(seed.path),
			DefinitionID: &defID,
		}
		if err := analyticsRepo.CreateEvent(ctx, event); err != nil {
			return counts, fmt.Errorf("create predefined event: %w", err)
		}

		data := buildEventData(def, seed.props)
		for _, entry := range data {
			entry.EventID = event.ID
		}
		if err := analyticsRepo.CreateEventDataBatch(ctx, data); err != nil {
			return counts, fmt.Errorf("create event data batch: %w", err)
		}
		counts.predefinedEvents++
	}

	return counts, nil
}

func normalizePaths(input []string) []string {
	paths := make([]string, 0, len(input))
	for _, path := range input {
		paths = append(paths, normalizeURL(path))
	}
	if len(paths) == 0 {
		return []string{"/"}
	}
	return paths
}

func buildEventData(def *event.Definition, props map[string]string) []*analyticspersistence.EventData {
	if def == nil {
		return nil
	}
	data := make([]*analyticspersistence.EventData, 0, len(def.Fields))
	for _, field := range def.Fields {
		value := props[field.Key]
		if value == "" {
			value = fallbackEventValue(field)
		}
		value = truncateString(value, field.MaxLength)
		data = append(data, &analyticspersistence.EventData{
			FieldID: field.ID,
			Value:   value,
		})
	}
	return data
}

func fallbackEventValue(field *event.Field) string {
	switch field.Key {
	case "plan":
		return "pro"
	case "referrer":
		return "pricing"
	case "amount":
		return "49"
	case "currency":
		return "USD"
	case "source":
		return "blog"
	case "video":
		return "demo"
	case "seconds":
		return "90"
	case "cta":
		return "Start Trial"
	case "page":
		return "/pricing"
	case "file":
		return "report.pdf"
	case "success":
		return "true"
	default:
		switch field.Type {
		case event.FieldTypeInt:
			return "1"
		case event.FieldTypeBool:
			return "false"
		default:
			return "value"
		}
	}
}

func randomStart(rng *mathrand.Rand, now time.Time, recentRemaining *int) time.Time {
	if recentRemaining != nil && *recentRemaining > 0 {
		*recentRemaining--
		return now.Add(-time.Duration(randRange(rng, 120, 900)) * time.Second)
	}
	seconds := randRange(rng, 60, maxDays*24*3600)
	return now.Add(-time.Duration(seconds) * time.Second)
}

func randRange(rng *mathrand.Rand, minimum, maximum int) int {
	if maximum <= minimum {
		return minimum
	}
	return rng.Intn(maximum-minimum+1) + minimum
}

func pickString(rng *mathrand.Rand, values []string) string {
	if len(values) == 0 {
		return ""
	}
	return values[rng.Intn(len(values))]
}

func pickEnum[T any](rng *mathrand.Rand, values []T) T {
	var zero T
	if len(values) == 0 {
		return zero
	}
	return values[rng.Intn(len(values))]
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func generateRandomString(length int) (string, error) {
	bytes := make([]byte, length/2)
	if _, err := cryptorand.Read(bytes); err != nil {
		return "", fmt.Errorf("generate random string: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}

func normalizeURL(path string) string {
	if path == "" {
		return "/"
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	if len(path) > 1 {
		path = strings.TrimSuffix(path, "/")
	}
	return path
}

func truncateString(value string, maxLength int) string {
	if len(value) <= maxLength {
		return value
	}
	return value[:maxLength]
}
