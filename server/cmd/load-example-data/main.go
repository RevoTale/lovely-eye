package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	mathrand "math/rand"
	"time"

	"github.com/lovely-eye/server/internal/config"
	"github.com/lovely-eye/server/internal/database"
	"github.com/lovely-eye/server/internal/models"
	"github.com/lovely-eye/server/internal/repository"
	"github.com/lovely-eye/server/internal/services"
	"github.com/lovely-eye/server/pkg/utils"
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

func main() {
	cfg := config.Load()

	db, err := database.New(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer func() {
		if err := database.Close(db); err != nil {
			log.Printf("DB close error: %v", err)
		}
	}()

	ctx := context.Background()
	if err := database.Migrate(ctx, db); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	owner, err := ensureSeedOwner(ctx, db, cfg.Auth)
	if err != nil {
		log.Fatalf("Failed to resolve seed user: %v", err)
	}

	site, created, err := ensureLocalhostSite(ctx, db, owner.ID)
	if err != nil {
		log.Fatalf("Failed to ensure localhost site: %v", err)
	}

	defs, err := ensureEventDefinitions(ctx, db, site.ID)
	if err != nil {
		log.Fatalf("Failed to ensure event definitions: %v", err)
	}

	counts, err := seedData(ctx, db, site.ID, defs)
	if err != nil {
		log.Fatalf("Failed to seed analytics data: %v", err)
	}

	if created {
		log.Printf("Created site %q for localhost with public key %s", site.Name, site.PublicKey)
	}

	log.Printf(
		"Seeded %d clients, %d sessions, %d page views, %d predefined events",
		counts.clients,
		counts.sessions,
		counts.pageViews,
		counts.predefinedEvents,
	)
}

type seedCounts struct {
	clients         int
	sessions        int
	pageViews       int
	predefinedEvents int
}

func firstUser(ctx context.Context, db *bun.DB) (*models.User, error) {
	user := new(models.User)
	if err := db.NewSelect().
		Model(user).
		Order("id ASC").
		Limit(1).
		Scan(ctx); err != nil {
		return nil, fmt.Errorf("scan first user: %w", err)
	}
	return user, nil
}

func ensureSeedOwner(ctx context.Context, db *bun.DB, authCfg config.AuthConfig) (*models.User, error) {
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

	userRepo := repository.NewUserRepository(db)
	newUser := &models.User{
		Username:     username,
		PasswordHash: string(hash),
		Role:         "admin",
	}
	if err := userRepo.Create(ctx, newUser); err != nil {
		return nil, fmt.Errorf("create seed user: %w", err)
	}
	return newUser, nil
}

func ensureLocalhostSite(ctx context.Context, db *bun.DB, userID int64) (*models.Site, bool, error) {
	siteRepo := repository.NewSiteRepository(db)
	siteService := services.NewSiteService(siteRepo)

	site, err := siteRepo.GetByDomainForUser(ctx, userID, "localhost")
	created := false
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, false, fmt.Errorf("get site by domain for user: %w", err)
		}
		site, err = siteService.Create(ctx, services.CreateSiteInput{
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

func ensureEventDefinitions(ctx context.Context, db *bun.DB, siteID int64) ([]*models.EventDefinition, error) {
	eventRepo := repository.NewEventDefinitionRepository(db)
	eventService := services.NewEventDefinitionService(eventRepo)

	definitions := []services.EventDefinitionInput{
		{
			Name: "signup",
			Fields: []services.EventFieldInput{
				{Key: "plan", Type: "string", Required: true},
				{Key: "referrer", Type: "string"},
			},
		},
		{
			Name: "purchase",
			Fields: []services.EventFieldInput{
				{Key: "amount", Type: "int", Required: true},
				{Key: "currency", Type: "string", Required: true},
				{Key: "plan", Type: "string"},
			},
		},
		{
			Name: "newsletter_subscribe",
			Fields: []services.EventFieldInput{
				{Key: "source", Type: "string"},
			},
		},
		{
			Name: "video_play",
			Fields: []services.EventFieldInput{
				{Key: "video", Type: "string", Required: true},
				{Key: "seconds", Type: "int"},
			},
		},
		{
			Name: "cta_click",
			Fields: []services.EventFieldInput{
				{Key: "cta", Type: "string", Required: true},
				{Key: "page", Type: "string"},
			},
		},
		{
			Name: "file_download",
			Fields: []services.EventFieldInput{
				{Key: "file", Type: "string", Required: true},
				{Key: "success", Type: "bool"},
			},
		},
	}

	results := make([]*models.EventDefinition, 0, len(definitions))
	for _, def := range definitions {
		created, err := eventService.Upsert(ctx, siteID, def)
		if err != nil {
			return nil, fmt.Errorf("upsert event definition %q: %w", def.Name, err)
		}
		results = append(results, created)
	}
	return results, nil
}

func seedData(ctx context.Context, db *bun.DB, siteID int64, defs []*models.EventDefinition) (seedCounts, error) {
	counts := seedCounts{}
	analyticsRepo := repository.NewAnalyticsRepository(db)
	//nolint:gosec // Deterministic randomness isn't required for seed data.
	rng := mathrand.New(mathrand.NewSource(time.Now().UnixNano()))

	defMap := make(map[string]*models.EventDefinition, len(defs))
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

	for i := 0; i < defaultUsers; i++ {
		hash, err := utils.GenerateRandomString(64)
		if err != nil {
			return counts, fmt.Errorf("generate client hash: %w", err)
		}

		client := &models.Client{
			SiteID:     siteID,
			Hash:       hash,
			Country:    pickString(rng, []string{"US", "GB", "DE", "FR", "CA", "NL"}),
			Device:     pickString(rng, []string{"desktop", "mobile", "tablet"}),
			Browser:    pickString(rng, []string{"Chrome", "Safari", "Firefox", "Edge"}),
			OS:         pickString(rng, []string{"Windows", "macOS", "Linux", "iOS", "Android"}),
			ScreenSize: pickString(rng, []string{"1920x1080", "1366x768", "390x844", "1440x900"}),
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
	pageViews       int
	predefinedEvents int
}

func applyPattern(
	ctx context.Context,
	analyticsRepo *repository.AnalyticsRepository,
	rng *mathrand.Rand,
	clientID, siteID int64,
	start time.Time,
	pattern behaviorPattern,
	defMap map[string]*models.EventDefinition,
) (patternCounts, error) {
	counts := patternCounts{}

	paths := normalizePaths(pattern.paths)
	durationMinutes := randRange(rng, sessionBaseMins, sessionBaseMins+sessionExtraMins)
	duration := durationMinutes*60 + randRange(rng, 0, 120)
	enterUnix := start.Unix()
	exitUnix := enterUnix + int64(duration)

	session := &models.Session{
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
		event := &models.Event{
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
		event := &models.Event{
			SessionID:    session.ID,
			Time:         eventTime,
			Hour:         eventTime / 3600,
			Day:          eventTime / 86400,
			Path:         utils.NormalizeURL(seed.path),
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
		paths = append(paths, utils.NormalizeURL(path))
	}
	if len(paths) == 0 {
		return []string{"/"}
	}
	return paths
}

func buildEventData(def *models.EventDefinition, props map[string]string) []*models.EventData {
	if def == nil {
		return nil
	}
	data := make([]*models.EventData, 0, len(def.Fields))
	for _, field := range def.Fields {
		value := props[field.Key]
		if value == "" {
			value = fallbackEventValue(field)
		}
		value = utils.TruncateString(value, field.MaxLength)
		data = append(data, &models.EventData{
			FieldID: field.ID,
			Value:   value,
		})
	}
	return data
}

func fallbackEventValue(field *models.EventDefinitionField) string {
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
		case models.FieldTypeInt:
			return "1"
		case models.FieldTypeBool:
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

func randRange(rng *mathrand.Rand, min, max int) int {
	if max <= min {
		return min
	}
	return rng.Intn(max-min+1) + min
}

func pickString(rng *mathrand.Rand, values []string) string {
	if len(values) == 0 {
		return ""
	}
	return values[rng.Intn(len(values))]
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
