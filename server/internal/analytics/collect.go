package analytics

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	analyticspersistence "github.com/lovely-eye/server/internal/analytics/persistence"
	"github.com/lovely-eye/server/internal/event"
	"github.com/lovely-eye/server/internal/site"
	"github.com/mileusna/useragent"
	"github.com/uptrace/bun"
)

type clientDimensions struct {
	device     analyticspersistence.ClientDevice
	browser    analyticspersistence.ClientBrowser
	os         analyticspersistence.ClientOS
	screenSize analyticspersistence.ClientScreenSize
}

type activeSessionLookup struct {
	session *analyticspersistence.Session
}

func (s *Service) CollectPageView(ctx context.Context, input CollectInput) error {
	site, accepted, err := s.acceptedAnalyticsSite(ctx, input.SiteKey, input.UserAgent, input.Origin, input.Referer, input.IP)
	if err != nil {
		return err
	}
	if !accepted {
		return nil
	}

	dimensions := parseClientDimensions(input.UserAgent)
	now := s.now()
	nowUnix := now.Unix()
	country := s.collectCountry(site, input.IP, !input.Exit)

	if err := s.analyticsRepo.RunInTx(ctx, func(ctx context.Context, tx bun.Tx) error {
		return s.collectPageViewTx(ctx, tx, site.ID, input, dimensions, country, now, nowUnix)
	}); err != nil {
		return fmt.Errorf("collect page view transaction: %w", err)
	}

	return nil
}

func (s *Service) acceptedAnalyticsSite(
	ctx context.Context,
	siteKey string,
	userAgent string,
	origin string,
	referer string,
	ip string,
) (*site.Site, bool, error) {
	if s.botDetector.IsBot(userAgent) {
		return nil, false, nil
	}
	site, err := s.siteRepo.GetByPublicKey(ctx, siteKey)
	if err != nil {
		return nil, false, fmt.Errorf("get site by public key: %w", err)
	}
	if !IsAllowedDomain(origin, referer, site.Domains) || s.isBlockedRequest(site, ip) {
		return site, false, nil
	}
	return site, true, nil
}

func parseClientDimensions(userAgent string) clientDimensions {
	ua := useragent.Parse(userAgent)
	return clientDimensions{
		device:     categorizeDevice(ua),
		browser:    normalizeBrowser(ua),
		os:         normalizeOS(ua),
		screenSize: analyticspersistence.ClientScreenSizeUnknown,
	}
}

func (s *Service) collectCountry(site *site.Site, ip string, enabled bool) Country {
	if enabled && site.TrackCountry && s.geoIPService != nil {
		return s.resolveCountryBestEffort(ip)
	}
	return UnknownCountry
}

func (s *Service) collectPageViewTx(
	ctx context.Context,
	tx bun.Tx,
	siteID int64,
	input CollectInput,
	dimensions clientDimensions,
	country Country,
	now time.Time,
	nowUnix int64,
) error {
	client, err := s.resolvePageViewClient(ctx, tx, siteID, input, dimensions, country, now)
	if errors.Is(err, errExitClientNotFound) {
		return nil
	}
	if err != nil {
		return err
	}

	activeSession, err := s.activeSessionTx(ctx, tx, siteID, client.ID, now, s.sessionLookupWindow(input.Exit))
	if err != nil {
		return err
	}
	session := activeSession.session
	session, insertEvent, err := s.applyPageViewSessionTx(ctx, tx, siteID, client.ID, input, session, now, nowUnix)
	if err != nil || !insertEvent {
		return err
	}
	return s.insertPageViewEventTx(ctx, tx, session.ID, input.Path, nowUnix)
}

func (s *Service) resolvePageViewClient(
	ctx context.Context,
	tx bun.Tx,
	siteID int64,
	input CollectInput,
	dimensions clientDimensions,
	country Country,
	now time.Time,
) (*analyticspersistence.Client, error) {
	if input.Exit {
		client, err := s.findClientForExit(ctx, tx, siteID, input.IP, dimensions.browser, dimensions.device, now)
		if err != nil {
			return nil, fmt.Errorf("find client for exit: %w", err)
		}
		return client, nil
	}
	client, err := s.resolveClientWithRotation(
		ctx,
		tx,
		siteID,
		input.IP,
		dimensions.device,
		dimensions.browser,
		dimensions.os,
		dimensions.screenSize,
		country.ISOCode,
		now,
	)
	if err != nil {
		return nil, fmt.Errorf("resolve client with rotation: %w", err)
	}
	return client, nil
}

func (s *Service) activeSessionTx(
	ctx context.Context,
	tx bun.Tx,
	siteID int64,
	clientID int64,
	now time.Time,
	window time.Duration,
) (activeSessionLookup, error) {
	session, err := s.analyticsRepo.GetActiveSessionTx(ctx, tx, siteID, clientID, now.Add(-window).Unix())
	if errors.Is(err, sql.ErrNoRows) {
		return activeSessionLookup{}, nil
	}
	if err != nil {
		return activeSessionLookup{}, fmt.Errorf("get active session: %w", err)
	}
	return activeSessionLookup{session: session}, nil
}

func (s *Service) applyPageViewSessionTx(
	ctx context.Context,
	tx bun.Tx,
	siteID int64,
	clientID int64,
	input CollectInput,
	session *analyticspersistence.Session,
	now time.Time,
	nowUnix int64,
) (*analyticspersistence.Session, bool, error) {
	if session == nil {
		return s.createPageViewSessionTx(ctx, tx, siteID, clientID, input, nowUnix)
	}
	activeSinceUnix := now.Add(-activeSessionWindow).Unix()
	return s.updatePageViewSessionTx(ctx, tx, input, session, activeSinceUnix, nowUnix)
}

func (s *Service) createPageViewSessionTx(
	ctx context.Context,
	tx bun.Tx,
	siteID int64,
	clientID int64,
	input CollectInput,
	nowUnix int64,
) (*analyticspersistence.Session, bool, error) {
	if input.Exit {
		return nil, false, nil
	}
	session := &analyticspersistence.Session{
		SiteID:        siteID,
		ClientID:      clientID,
		EnterTime:     nowUnix,
		EnterHour:     nowUnix / 3600,
		EnterDay:      nowUnix / 86400,
		EnterPath:     input.Path,
		ExitTime:      nowUnix,
		ExitHour:      nowUnix / 3600,
		ExitDay:       nowUnix / 86400,
		ExitPath:      input.Path,
		Referrer:      input.Referrer,
		UTMSource:     input.UTMSource,
		UTMMedium:     input.UTMMedium,
		UTMCampaign:   input.UTMCampaign,
		Duration:      0,
		PageViewCount: 1,
	}
	if err := s.analyticsRepo.CreateSessionTx(ctx, tx, session); err != nil {
		return nil, false, fmt.Errorf("create session: %w", err)
	}
	return session, true, nil
}

func (s *Service) updatePageViewSessionTx(
	ctx context.Context,
	tx bun.Tx,
	input CollectInput,
	session *analyticspersistence.Session,
	activeSinceUnix int64,
	nowUnix int64,
) (*analyticspersistence.Session, bool, error) {
	if input.Exit && session.ExitTime <= activeSinceUnix && session.ExitPath != input.Path {
		return session, false, nil
	}
	if input.Exit && session.ExitPath == input.Path {
		updateSessionExit(session, input.Path, s.singlePageExitUnix(session, nowUnix))
		return session, false, s.updateSessionTx(ctx, tx, session, "update exit session")
	}
	recentEvent, err := s.analyticsRepo.GetRecentPageViewEventTx(ctx, tx, session.ID, input.Path, nowUnix-10)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, false, fmt.Errorf("get recent pageview event: %w", err)
	}
	if recentEvent != nil {
		return s.applyDuplicatePageViewTx(ctx, tx, input, session, nowUnix)
	}
	updateSessionExit(session, input.Path, nowUnix)
	session.PageViewCount++
	return session, true, s.updateSessionTx(ctx, tx, session, "update session")
}

func (s *Service) applyDuplicatePageViewTx(
	ctx context.Context,
	tx bun.Tx,
	input CollectInput,
	session *analyticspersistence.Session,
	nowUnix int64,
) (*analyticspersistence.Session, bool, error) {
	if input.Exit {
		updateSessionExit(session, input.Path, s.exitUnixForDuplicateEvent(session, input.Path, nowUnix))
		return session, false, s.updateSessionTx(ctx, tx, session, "update duplicate exit session")
	}
	return session, false, nil
}

func (s *Service) updateSessionTx(ctx context.Context, tx bun.Tx, session *analyticspersistence.Session, label string) error {
	if err := s.analyticsRepo.UpdateSessionTx(ctx, tx, session); err != nil {
		return fmt.Errorf("%s: %w", label, err)
	}
	return nil
}

func (s *Service) insertPageViewEventTx(
	ctx context.Context,
	tx bun.Tx,
	sessionID int64,
	path string,
	nowUnix int64,
) error {
	event := &analyticspersistence.Event{
		SessionID:    sessionID,
		Time:         nowUnix,
		Hour:         nowUnix / 3600,
		Day:          nowUnix / 86400,
		Path:         path,
		DefinitionID: nil,
	}
	if err := s.analyticsRepo.CreateEventTx(ctx, tx, event); err != nil {
		return fmt.Errorf("create event: %w", err)
	}
	return nil
}

func (s *Service) singlePageExitUnix(session *analyticspersistence.Session, nowUnix int64) int64 {
	if session.PageViewCount > 1 {
		return nowUnix
	}

	maxExitUnix := session.EnterTime + s.maxSinglePageDurationSeconds()
	if nowUnix > maxExitUnix {
		return maxExitUnix
	}
	return nowUnix
}

func (s *Service) exitUnixForDuplicateEvent(session *analyticspersistence.Session, path string, nowUnix int64) int64 {
	if session.ExitPath == path {
		return s.singlePageExitUnix(session, nowUnix)
	}
	return nowUnix
}

func (s *Service) maxSinglePageDurationSeconds() int64 {
	duration := s.maxSinglePageDuration
	if duration <= 0 {
		duration = defaultMaxSinglePageExitDuration
	}

	seconds := int64(duration / time.Second)
	if seconds <= 0 {
		return int64(defaultMaxSinglePageExitDuration / time.Second)
	}
	return seconds
}

func updateSessionExit(session *analyticspersistence.Session, path string, nowUnix int64) {
	session.ExitTime = nowUnix
	session.ExitHour = nowUnix / 3600
	session.ExitDay = nowUnix / 86400
	if path != "" {
		session.ExitPath = path
	}
	session.Duration = int(nowUnix - session.EnterTime)
}

func (s *Service) CollectEvent(ctx context.Context, input EventInput) error {
	site, accepted, err := s.acceptedAnalyticsSite(ctx, input.SiteKey, input.UserAgent, input.Origin, input.Referer, input.IP)
	if err != nil {
		return err
	}
	if !accepted || s.eventDefinitionStore == nil {
		return nil
	}

	definition, sanitizedProps, ok, err := s.eventDefinitionForCollect(ctx, site.ID, input)
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}

	dimensions := parseClientDimensions(input.UserAgent)
	now := s.now()
	nowUnix := now.Unix()
	country := s.collectCountry(site, input.IP, true)

	if err := s.analyticsRepo.RunInTx(ctx, func(ctx context.Context, tx bun.Tx) error {
		return s.collectEventTx(ctx, tx, site.ID, input, dimensions, country, now, nowUnix, definition, sanitizedProps)
	}); err != nil {
		return fmt.Errorf("collect event transaction: %w", err)
	}

	return nil
}

func (s *Service) eventDefinitionForCollect(
	ctx context.Context,
	siteID int64,
	input EventInput,
) (*event.Definition, string, bool, error) {
	definition, err := s.eventDefinitionStore.GetByName(ctx, siteID, input.Name)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, "", false, nil
	}
	if err != nil {
		return nil, "", false, fmt.Errorf("get event definition by name: %w", err)
	}
	sanitizedProps, ok, err := sanitizeEventProperties(input.Properties, definition.Fields)
	if err != nil {
		return nil, "", false, fmt.Errorf("sanitize event properties: %w", err)
	}
	return definition, sanitizedProps, ok, nil
}

func (s *Service) collectEventTx(
	ctx context.Context,
	tx bun.Tx,
	siteID int64,
	input EventInput,
	dimensions clientDimensions,
	country Country,
	now time.Time,
	nowUnix int64,
	definition *event.Definition,
	sanitizedProps string,
) error {
	client, err := s.resolveClientWithRotation(
		ctx,
		tx,
		siteID,
		input.IP,
		dimensions.device,
		dimensions.browser,
		dimensions.os,
		dimensions.screenSize,
		country.ISOCode,
		now,
	)
	if err != nil {
		return fmt.Errorf("resolve client with rotation: %w", err)
	}

	session, err := s.eventSessionTx(ctx, tx, siteID, client.ID, input.Path, now, nowUnix)
	if err != nil {
		return err
	}
	event, err := s.insertCustomEventTx(ctx, tx, session.ID, input.Path, nowUnix, definition.ID)
	if err != nil || sanitizedProps == "" {
		return err
	}
	return s.insertCustomEventDataTx(ctx, tx, event.ID, sanitizedProps, definition.Fields)
}

func (s *Service) eventSessionTx(
	ctx context.Context,
	tx bun.Tx,
	siteID int64,
	clientID int64,
	path string,
	now time.Time,
	nowUnix int64,
) (*analyticspersistence.Session, error) {
	activeSession, err := s.activeSessionTx(ctx, tx, siteID, clientID, now, activeSessionWindow)
	if err != nil {
		return nil, err
	}
	session := activeSession.session
	if session == nil {
		return s.createEventSessionTx(ctx, tx, siteID, clientID, path, nowUnix)
	}
	updateSessionExit(session, path, nowUnix)
	if err := s.analyticsRepo.UpdateSessionTx(ctx, tx, session); err != nil {
		return nil, fmt.Errorf("update session: %w", err)
	}
	return session, nil
}

func (s *Service) createEventSessionTx(
	ctx context.Context,
	tx bun.Tx,
	siteID int64,
	clientID int64,
	path string,
	nowUnix int64,
) (*analyticspersistence.Session, error) {
	entryPath := path
	if entryPath == "" {
		entryPath = "/"
	}
	session := &analyticspersistence.Session{
		SiteID:        siteID,
		ClientID:      clientID,
		EnterTime:     nowUnix,
		EnterHour:     nowUnix / 3600,
		EnterDay:      nowUnix / 86400,
		EnterPath:     entryPath,
		ExitTime:      nowUnix,
		ExitHour:      nowUnix / 3600,
		ExitDay:       nowUnix / 86400,
		ExitPath:      entryPath,
		Referrer:      "",
		UTMSource:     "",
		UTMMedium:     "",
		UTMCampaign:   "",
		Duration:      0,
		PageViewCount: 0,
	}
	if err := s.analyticsRepo.CreateSessionTx(ctx, tx, session); err != nil {
		return nil, fmt.Errorf("create session: %w", err)
	}
	return session, nil
}

func (s *Service) insertCustomEventTx(
	ctx context.Context,
	tx bun.Tx,
	sessionID int64,
	path string,
	nowUnix int64,
	definitionID int64,
) (*analyticspersistence.Event, error) {
	event := &analyticspersistence.Event{
		SessionID:    sessionID,
		Time:         nowUnix,
		Hour:         nowUnix / 3600,
		Day:          nowUnix / 86400,
		Path:         path,
		DefinitionID: &definitionID,
	}
	if err := s.analyticsRepo.CreateEventTx(ctx, tx, event); err != nil {
		return nil, fmt.Errorf("create event: %w", err)
	}
	return event, nil
}

func (s *Service) insertCustomEventDataTx(
	ctx context.Context,
	tx bun.Tx,
	eventID int64,
	sanitizedProps string,
	fields []*event.Field,
) error {
	var propsMap map[string]interface{}
	if err := json.Unmarshal([]byte(sanitizedProps), &propsMap); err != nil {
		return fmt.Errorf("unmarshal sanitized properties: %w", err)
	}
	eventDataList := eventDataFromProperties(eventID, propsMap, fields)
	if err := s.analyticsRepo.CreateEventDataBatchTx(ctx, tx, eventDataList); err != nil {
		return fmt.Errorf("create event data batch: %w", err)
	}
	return nil
}

func eventDataFromProperties(
	eventID int64,
	propsMap map[string]interface{},
	fields []*event.Field,
) []*analyticspersistence.EventData {
	fieldMap := make(map[string]int64, len(fields))
	for _, field := range fields {
		fieldMap[field.Key] = field.ID
	}
	eventDataList := make([]*analyticspersistence.EventData, 0, len(propsMap))
	for key, value := range propsMap {
		fieldID, exists := fieldMap[key]
		if exists {
			eventDataList = append(eventDataList, &analyticspersistence.EventData{
				EventID: eventID,
				FieldID: fieldID,
				Value:   eventDataStorageValue(value),
			})
		}
	}
	return eventDataList
}

func eventDataStorageValue(value interface{}) string {
	switch typedValue := value.(type) {
	case string:
		return typedValue
	case bool:
		return strconv.FormatBool(typedValue)
	case float64:
		return strconv.FormatFloat(typedValue, 'f', -1, 64)
	default:
		return fmt.Sprint(typedValue)
	}
}
