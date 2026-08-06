package persistence

import (
	eventpersistence "github.com/lovely-eye/server/internal/event/persistence"
	"github.com/uptrace/bun"
)

// Client is a pseudonymous visitor row resolved by UTC-day-skipped rotation.
type Client struct {
	bun.BaseModel `bun:"table:clients,alias:c"`

	ID         int64            `bun:"id,pk,autoincrement"`
	SiteID     int64            `bun:"site_id,notnull,unique:clients_site_id_hash"`
	Hash       string           `bun:"hash,notnull,type:varchar(64),unique:clients_site_id_hash"`
	Country    string           `bun:"country,type:varchar(2)"`
	Device     ClientDevice     `bun:"device,notnull,default:0"`
	Browser    ClientBrowser    `bun:"browser,notnull,default:0"`
	OS         ClientOS         `bun:"os,notnull,default:0"`
	ScreenSize ClientScreenSize `bun:"screen_size,notnull,default:0"`

	Sessions []*Session `bun:"rel:has-many,join:id=client_id"`
}

type Session struct {
	bun.BaseModel `bun:"table:sessions,alias:s"`

	ID       int64 `bun:"id,pk,autoincrement"`
	SiteID   int64 `bun:"site_id,notnull"`
	ClientID int64 `bun:"client_id,notnull"`

	EnterTime int64  `bun:"enter_time,notnull"`
	EnterHour int64  `bun:"enter_hour,notnull"`
	EnterDay  int64  `bun:"enter_day,notnull"`
	EnterPath string `bun:"enter_path,notnull,type:varchar(2048)"`

	ExitTime int64  `bun:"exit_time,notnull"`
	ExitHour int64  `bun:"exit_hour,notnull"`
	ExitDay  int64  `bun:"exit_day,notnull"`
	ExitPath string `bun:"exit_path,notnull,type:varchar(2048)"`

	Referrer    string `bun:"referrer,type:varchar(2048)"`
	UTMSource   string `bun:"utm_source,type:varchar(128)"`
	UTMMedium   string `bun:"utm_medium,type:varchar(128)"`
	UTMCampaign string `bun:"utm_campaign,type:varchar(256)"`

	Duration      int `bun:"duration,notnull,default:0"`
	PageViewCount int `bun:"page_view_count,notnull,default:0"`

	Client *Client  `bun:"rel:belongs-to,join:client_id=id"`
	Events []*Event `bun:"rel:has-many,join:id=session_id"`
}

type Event struct {
	bun.BaseModel `bun:"table:events,alias:e"`

	ID        int64  `bun:"id,pk,autoincrement"`
	SessionID int64  `bun:"session_id,notnull"`
	Time      int64  `bun:"time,notnull"`
	Hour      int64  `bun:"hour,notnull"`
	Day       int64  `bun:"day,notnull"`
	Path      string `bun:"path,notnull,type:varchar(2048)"`

	DefinitionID *int64 `bun:"definition_id"`

	Session    *Session                     `bun:"rel:belongs-to,join:session_id=id"`
	Definition *eventpersistence.Definition `bun:"rel:belongs-to,join:definition_id=id"`
	Data       []*EventData                 `bun:"rel:has-many,join:id=event_id"`
}

type EventData struct {
	bun.BaseModel `bun:"table:event_data,alias:evd"`

	ID      int64  `bun:"id,pk,autoincrement"`
	EventID int64  `bun:"event_id,notnull"`
	FieldID int64  `bun:"field_id,notnull"`
	Value   string `bun:"value,notnull,type:varchar(1024)"`

	Event *Event                  `bun:"rel:belongs-to,join:event_id=id"`
	Field *eventpersistence.Field `bun:"rel:belongs-to,join:field_id=id"`
}
