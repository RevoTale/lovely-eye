package persistence

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/uptrace/bun"
)

func (r *Repository) FindClientByHashesTx(ctx context.Context, tx bun.IDB, siteID int64, todayHash, yesterdayHash string) (*Client, error) {
	client := new(Client)
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

func (r *Repository) FindClientByHashTx(ctx context.Context, tx bun.IDB, siteID int64, hash string) (*Client, error) {
	client := new(Client)
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

func (r *Repository) CreateClientTx(ctx context.Context, tx bun.IDB, client *Client) error {
	_, err := tx.NewInsert().Model(client).Exec(ctx)
	if err != nil {
		return fmt.Errorf("create client: %w", err)
	}
	return nil
}

func (r *Repository) UpdateClientTx(ctx context.Context, tx bun.IDB, client *Client) error {
	_, err := tx.NewUpdate().Model(client).WherePK().Exec(ctx)
	if err != nil {
		return fmt.Errorf("update client: %w", err)
	}
	return nil
}

func (r *Repository) GetActiveSessionTx(ctx context.Context, tx bun.IDB, siteID, clientID int64, sinceUnix int64) (*Session, error) {
	session := new(Session)
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

func (r *Repository) isPostgres() bool {
	dialect := fmt.Sprint(r.db.Dialect().Name())
	return dialect == "pg" || dialect == "postgres" || dialect == "postgresql"
}

func (r *Repository) GetRecentPageViewEventTx(ctx context.Context, tx bun.IDB, sessionID int64, path string, since int64) (*Event, error) {
	event := new(Event)
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

func (r *Repository) CreateSession(ctx context.Context, session *Session) error {
	_, err := r.db.NewInsert().Model(session).Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}
	return nil
}

func (r *Repository) CreateSessionTx(ctx context.Context, tx bun.IDB, session *Session) error {
	_, err := tx.NewInsert().Model(session).Exec(ctx)
	if err != nil {
		return fmt.Errorf("create session: %w", err)
	}
	return nil
}

func (r *Repository) UpdateSessionTx(ctx context.Context, tx bun.IDB, session *Session) error {
	_, err := tx.NewUpdate().Model(session).WherePK().Exec(ctx)
	if err != nil {
		return fmt.Errorf("update session: %w", err)
	}
	return nil
}

func (r *Repository) CreateEvent(ctx context.Context, event *Event) error {
	_, err := r.db.NewInsert().Model(event).Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to create event: %w", err)
	}
	return nil
}

func (r *Repository) CreateEventTx(ctx context.Context, tx bun.IDB, event *Event) error {
	_, err := tx.NewInsert().Model(event).Exec(ctx)
	if err != nil {
		return fmt.Errorf("create event: %w", err)
	}
	return nil
}

func (r *Repository) GetEvents(ctx context.Context, siteID int64, from, to time.Time, limit, offset int) ([]*Event, error) {
	var events []*Event
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

func (r *Repository) GetEventsWithFilter(ctx context.Context, query AnalyticsQuery) ([]*Event, error) {
	var events []*Event
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

func (r *Repository) CreateEventDataBatch(ctx context.Context, eventDataList []*EventData) error {
	if len(eventDataList) == 0 {
		return nil
	}
	_, err := r.db.NewInsert().Model(&eventDataList).Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to create event data batch: %w", err)
	}
	return nil
}

func (r *Repository) CreateEventDataBatchTx(ctx context.Context, tx bun.IDB, eventDataList []*EventData) error {
	if len(eventDataList) == 0 {
		return nil
	}
	_, err := tx.NewInsert().Model(&eventDataList).Exec(ctx)
	if err != nil {
		return fmt.Errorf("create event data batch: %w", err)
	}
	return nil
}
