package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/lovely-eye/server/internal/models"
	"github.com/uptrace/bun"
)

func (r *AnalyticsRepository) FindOrCreateClient(
	ctx context.Context,
	siteID int64,
	hash string,
	device models.ClientDevice,
	browser models.ClientBrowser,
	os models.ClientOS,
	screenSize models.ClientScreenSize,
	country string,
) (*models.Client, error) {
	// Try to find existing client by hash
	client := new(models.Client)
	err := r.db.NewSelect().
		Model(client).
		Where("site_id = ?", siteID).
		Where("hash = ?", hash).
		Limit(1).
		Scan(ctx)

	if err == nil {

		return client, nil
	}

	client = &models.Client{
		SiteID:     siteID,
		Hash:       hash,
		Device:     device,
		Browser:    browser,
		OS:         os,
		ScreenSize: screenSize,
		Country:    country,
	}

	_, err = r.db.NewInsert().Model(client).Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	return client, nil
}

func (r *AnalyticsRepository) FindClientByHashesTx(ctx context.Context, tx bun.IDB, siteID int64, todayHash, yesterdayHash string) (*models.Client, error) {
	client := new(models.Client)
	err := tx.NewSelect().
		Model(client).
		Where("site_id = ?", siteID).
		Where("(hash = ? OR hash = ?)", todayHash, yesterdayHash).
		OrderExpr("CASE WHEN hash = ? THEN 0 ELSE 1 END", todayHash).
		Order("id DESC").
		Limit(1).
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("find client by hashes: %w", err)
	}
	return client, nil
}

func (r *AnalyticsRepository) FindClientByHashTx(ctx context.Context, tx bun.IDB, siteID int64, hash string) (*models.Client, error) {
	client := new(models.Client)
	err := tx.NewSelect().
		Model(client).
		Where("site_id = ?", siteID).
		Where("hash = ?", hash).
		Limit(1).
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("find client by hash: %w", err)
	}
	return client, nil
}

func (r *AnalyticsRepository) CreateClientTx(ctx context.Context, tx bun.IDB, client *models.Client) error {
	_, err := tx.NewInsert().Model(client).Exec(ctx)
	if err != nil {
		return fmt.Errorf("create client: %w", err)
	}
	return nil
}

func (r *AnalyticsRepository) UpdateClientTx(ctx context.Context, tx bun.IDB, client *models.Client) error {
	_, err := tx.NewUpdate().Model(client).WherePK().Exec(ctx)
	if err != nil {
		return fmt.Errorf("update client: %w", err)
	}
	return nil
}

func (r *AnalyticsRepository) GetActiveSession(ctx context.Context, siteID, clientID int64, sinceUnix int64) (*models.Session, error) {
	session := new(models.Session)
	err := r.db.NewSelect().
		Model(session).
		Where("site_id = ?", siteID).
		Where("client_id = ?", clientID).
		Where("exit_time > ?", sinceUnix).
		Order("exit_time DESC").
		Limit(1).
		Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get active session: %w", err)
	}
	return session, nil
}

func (r *AnalyticsRepository) GetActiveSessionTx(ctx context.Context, tx bun.IDB, siteID, clientID int64, sinceUnix int64) (*models.Session, error) {
	session := new(models.Session)
	q := tx.NewSelect().
		Model(session).
		Where("site_id = ?", siteID).
		Where("client_id = ?", clientID).
		Where("exit_time > ?", sinceUnix).
		Order("exit_time DESC").
		Limit(1)
	if r.isPostgres() {
		q = q.For("UPDATE")
	}
	err := q.Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("get active session: %w", err)
	}
	return session, nil
}

func (r *AnalyticsRepository) isPostgres() bool {
	dialect := fmt.Sprint(r.db.Dialect().Name())
	return dialect == "pg" || dialect == "postgres" || dialect == "postgresql"
}

func (r *AnalyticsRepository) GetRecentPageViewEvent(ctx context.Context, sessionID int64, path string, since int64) (*models.Event, error) {
	event := new(models.Event)
	err := r.db.NewSelect().
		Model(event).
		Where("session_id = ?", sessionID).
		Where("path = ?", path).
		Where("definition_id IS NULL").
		Where("time > ?", since).
		Order("time DESC").
		Limit(1).
		Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent page view event: %w", err)
	}
	return event, nil
}

func (r *AnalyticsRepository) GetRecentPageViewEventTx(ctx context.Context, tx bun.IDB, sessionID int64, path string, since int64) (*models.Event, error) {
	event := new(models.Event)
	err := tx.NewSelect().
		Model(event).
		Where("session_id = ?", sessionID).
		Where("path = ?", path).
		Where("definition_id IS NULL").
		Where("time > ?", since).
		Order("time DESC").
		Limit(1).
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("get recent page view event: %w", err)
	}
	return event, nil
}

func (r *AnalyticsRepository) CreateSession(ctx context.Context, session *models.Session) error {
	_, err := r.db.NewInsert().Model(session).Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}
	return nil
}

func (r *AnalyticsRepository) CreateSessionTx(ctx context.Context, tx bun.IDB, session *models.Session) error {
	_, err := tx.NewInsert().Model(session).Exec(ctx)
	if err != nil {
		return fmt.Errorf("create session: %w", err)
	}
	return nil
}

func (r *AnalyticsRepository) GetSession(ctx context.Context, id int64) (*models.Session, error) {
	session := new(models.Session)
	err := r.db.NewSelect().Model(session).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}
	return session, nil
}

func (r *AnalyticsRepository) GetSessionTx(ctx context.Context, tx bun.IDB, id int64) (*models.Session, error) {
	session := new(models.Session)
	err := tx.NewSelect().Model(session).Where("id = ?", id).Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("get session: %w", err)
	}
	return session, nil
}

func (r *AnalyticsRepository) UpdateSession(ctx context.Context, session *models.Session) error {
	_, err := r.db.NewUpdate().Model(session).WherePK().Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to update session: %w", err)
	}
	return nil
}

func (r *AnalyticsRepository) UpdateSessionTx(ctx context.Context, tx bun.IDB, session *models.Session) error {
	_, err := tx.NewUpdate().Model(session).WherePK().Exec(ctx)
	if err != nil {
		return fmt.Errorf("update session: %w", err)
	}
	return nil
}

func (r *AnalyticsRepository) CreateEvent(ctx context.Context, event *models.Event) error {
	_, err := r.db.NewInsert().Model(event).Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to create event: %w", err)
	}
	return nil
}

func (r *AnalyticsRepository) CreateEventTx(ctx context.Context, tx bun.IDB, event *models.Event) error {
	_, err := tx.NewInsert().Model(event).Exec(ctx)
	if err != nil {
		return fmt.Errorf("create event: %w", err)
	}
	return nil
}

func (r *AnalyticsRepository) GetEvents(ctx context.Context, siteID int64, from, to time.Time, limit, offset int) ([]*models.Event, error) {
	var events []*models.Event
	fromUnix := from.Unix()
	toUnix := to.Unix()
	err := r.db.NewSelect().
		Model(&events).
		Relation("Data.Field").
		Relation("Definition.Fields").
		Join("INNER JOIN sessions s ON e.session_id = s.id").
		Where("s.site_id = ?", siteID).
		Where("e.time >= ?", fromUnix).
		Where("e.time <= ?", toUnix).
		Order("e.time DESC").
		Limit(limit).
		Offset(offset).
		Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get events: %w", err)
	}
	return events, nil
}

func (r *AnalyticsRepository) GetEventsWithFilter(ctx context.Context, query AnalyticsQuery) ([]*models.Event, error) {
	var events []*models.Event
	fromUnix := query.From.Unix()
	toUnix := query.To.Unix()
	q := r.db.NewSelect().
		Model(&events).
		Relation("Data.Field").
		Relation("Definition.Fields").
		Join("INNER JOIN sessions s ON e.session_id = s.id").
		Where("s.site_id = ?", query.SiteID).
		Where("e.time >= ?", fromUnix).
		Where("e.time <= ?", toUnix)
	q = applyEventFilters(q, query.Filter)
	q = applyEventNamePathFilters(q, query.Filter)
	err := q.Order("e.time DESC").
		Limit(query.Limit).
		Offset(query.Offset).
		Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get events with filter: %w", err)
	}
	return events, nil
}

func (r *AnalyticsRepository) CreateEventData(ctx context.Context, eventData *models.EventData) error {
	_, err := r.db.NewInsert().Model(eventData).Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to create event data: %w", err)
	}
	return nil
}

func (r *AnalyticsRepository) CreateEventDataBatch(ctx context.Context, eventDataList []*models.EventData) error {
	if len(eventDataList) == 0 {
		return nil
	}
	_, err := r.db.NewInsert().Model(&eventDataList).Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to create event data batch: %w", err)
	}
	return nil
}

func (r *AnalyticsRepository) CreateEventDataBatchTx(ctx context.Context, tx bun.IDB, eventDataList []*models.EventData) error {
	if len(eventDataList) == 0 {
		return nil
	}
	_, err := tx.NewInsert().Model(&eventDataList).Exec(ctx)
	if err != nil {
		return fmt.Errorf("create event data batch: %w", err)
	}
	return nil
}
