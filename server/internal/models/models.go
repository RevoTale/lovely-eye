package models

import (
	"time"

	"github.com/uptrace/bun"
)

// User represents an authenticated user who can view analytics
type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	ID           int64     `bun:"id,pk,autoincrement" json:"id"`
	Username     string    `bun:"username,unique,notnull" json:"username"`
	PasswordHash string    `bun:"password_hash,notnull" json:"-"`
	Role         string    `bun:"role,notnull,default:'user'" json:"role"`
	Email        string    `bun:"email" json:"email"`
	CreatedAt    time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt    time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp" json:"updated_at"`

	Sites []*Site `bun:"rel:has-many,join:id=user_id" json:"sites,omitempty"`
}

type Site struct {
	bun.BaseModel `bun:"table:sites,alias:s"`

	ID           int64     `bun:"id,pk,autoincrement" json:"id"`
	UserID       int64     `bun:"user_id,notnull" json:"user_id"`
	Name         string    `bun:"name,notnull" json:"name"`
	PublicKey    string    `bun:"public_key,unique,notnull" json:"public_key"`
	TrackCountry bool      `bun:"track_country,notnull,default:false" json:"track_country"`
	CreatedAt    time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt    time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp" json:"updated_at"`

	User             *User                 `bun:"rel:belongs-to,join:user_id=id" json:"user,omitempty"`
	Domains          []*SiteDomain         `bun:"rel:has-many,join:id=site_id" json:"domains,omitempty"`
	BlockedIPs       []*SiteBlockedIP      `bun:"rel:has-many,join:id=site_id" json:"blocked_ips,omitempty"`
	BlockedCountries []*SiteBlockedCountry `bun:"rel:has-many,join:id=site_id" json:"blocked_countries,omitempty"`
	Clients          []*Client             `bun:"rel:has-many,join:id=site_id" json:"clients,omitempty"`
	Sessions         []*Session            `bun:"rel:has-many,join:id=site_id" json:"sessions,omitempty"`
	EventDefinitions []*EventDefinition    `bun:"rel:has-many,join:id=site_id" json:"event_definitions,omitempty"`
}

type SiteDomain struct {
	bun.BaseModel `bun:"table:site_domains,alias:sd"`

	ID        int64     `bun:"id,pk,autoincrement" json:"id"`
	SiteID    int64     `bun:"site_id,notnull,unique:site_domains_site_id_domain,index:site_domains_site_id_position" json:"site_id"`
	Domain    string    `bun:"domain,notnull,unique:site_domains_site_id_domain" json:"domain"`
	Position  int       `bun:"position,notnull,default:0,index:site_domains_site_id_position" json:"position"`
	CreatedAt time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp" json:"updated_at"`

	Site *Site `bun:"rel:belongs-to,join:site_id=id" json:"site,omitempty"`
}

type SiteBlockedIP struct {
	bun.BaseModel `bun:"table:site_blocked_ips,alias:sbi"`

	ID        int64     `bun:"id,pk,autoincrement" json:"id"`
	SiteID    int64     `bun:"site_id,notnull" json:"site_id"`
	IP        string    `bun:"ip,notnull" json:"ip"`
	CreatedAt time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp" json:"updated_at"`

	Site *Site `bun:"rel:belongs-to,join:site_id=id" json:"site,omitempty"`
}

type SiteBlockedCountry struct {
	bun.BaseModel `bun:"table:site_blocked_countries,alias:sbc"`

	ID          int64     `bun:"id,pk,autoincrement" json:"id"`
	SiteID      int64     `bun:"site_id,notnull" json:"site_id"`
	CountryCode string    `bun:"country_code,notnull" json:"country_code"`
	CreatedAt   time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt   time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp" json:"updated_at"`

	Site *Site `bun:"rel:belongs-to,join:site_id=id" json:"site,omitempty"`
}

// Client represents a unique visitor, deduplicated by hash.
// Stores STABLE attributes (don't change between sessions).
// Same visitor = same hash = same row → keeps DB small.
type Client struct {
	bun.BaseModel `bun:"table:clients,alias:c"`

	ID         int64  `bun:"id,pk,autoincrement" json:"id"`
	SiteID     int64  `bun:"site_id,notnull" json:"site_id"`
	Hash       string `bun:"hash,notnull,type:varchar(64)" json:"hash"` // SHA-256 hex, from current architecture
	Country    string `bun:"country,type:varchar(2)" json:"country"`
	Device     string `bun:"device,type:varchar(10)" json:"device"`
	Browser    string `bun:"browser,type:varchar(32)" json:"browser"`
	OS         string `bun:"os,type:varchar(32)" json:"os"`
	ScreenSize string `bun:"screen_size,type:varchar(16)" json:"screen_size"`

	Site     *Site      `bun:"rel:belongs-to,join:site_id=id" json:"site,omitempty"`
	Sessions []*Session `bun:"rel:has-many,join:id=client_id" json:"sessions,omitempty"`
}

type Session struct {
	bun.BaseModel `bun:"table:sessions,alias:s"`

	ID       int64 `bun:"id,pk,autoincrement" json:"id"`
	SiteID   int64 `bun:"site_id,notnull" json:"site_id"`
	ClientID int64 `bun:"client_id,notnull" json:"client_id"`

	EnterTime int64  `bun:"enter_time,notnull" json:"enter_time"`
	EnterHour int64  `bun:"enter_hour,notnull" json:"enter_hour"`
	EnterDay  int64  `bun:"enter_day,notnull" json:"enter_day"`
	EnterPath string `bun:"enter_path,notnull,type:varchar(2048)" json:"enter_path"`

	ExitTime int64  `bun:"exit_time,notnull" json:"exit_time"`
	ExitHour int64  `bun:"exit_hour,notnull" json:"exit_hour"`
	ExitDay  int64  `bun:"exit_day,notnull" json:"exit_day"`
	ExitPath string `bun:"exit_path,notnull,type:varchar(2048)" json:"exit_path"`

	Referrer    string `bun:"referrer,type:varchar(2048)" json:"referrer"`
	UTMSource   string `bun:"utm_source,type:varchar(128)" json:"utm_source"`
	UTMMedium   string `bun:"utm_medium,type:varchar(128)" json:"utm_medium"`
	UTMCampaign string `bun:"utm_campaign,type:varchar(256)" json:"utm_campaign"`

	Duration      int `bun:"duration,notnull,default:0" json:"duration"`
	PageViewCount int `bun:"page_view_count,notnull,default:0" json:"page_view_count"`

	Site   *Site    `bun:"rel:belongs-to,join:site_id=id" json:"site,omitempty"`
	Client *Client  `bun:"rel:belongs-to,join:client_id=id" json:"client,omitempty"`
	Events []*Event `bun:"rel:has-many,join:id=session_id" json:"events,omitempty"`
}

type Event struct {
	bun.BaseModel `bun:"table:events,alias:e"`

	ID        int64 `bun:"id,pk,autoincrement" json:"id"`
	SessionID int64 `bun:"session_id,notnull" json:"session_id"`

	Time int64 `bun:"time,notnull" json:"time"`
	Hour int64 `bun:"hour,notnull" json:"hour"`
	Day  int64 `bun:"day,notnull" json:"day"`

	Path string `bun:"path,notnull,type:varchar(2048)" json:"path"`

	DefinitionID *int64 `bun:"definition_id" json:"definition_id"`

	Session    *Session         `bun:"rel:belongs-to,join:session_id=id" json:"session,omitempty"`
	Definition *EventDefinition `bun:"rel:belongs-to,join:definition_id=id" json:"definition,omitempty"`
	Data       []*EventData     `bun:"rel:has-many,join:id=event_id" json:"data,omitempty"`
}

type EventDefinition struct {
	bun.BaseModel `bun:"table:event_definitions,alias:ed"`

	ID        int64     `bun:"id,pk,autoincrement" json:"id"`
	SiteID    int64     `bun:"site_id,notnull,unique:event_definitions_site_id_name" json:"site_id"`
	Name      string    `bun:"name,notnull,unique:event_definitions_site_id_name" json:"name"`
	CreatedAt time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp" json:"updated_at"`

	Site   *Site                   `bun:"rel:belongs-to,join:site_id=id" json:"site,omitempty"`
	Fields []*EventDefinitionField `bun:"rel:has-many,join:id=event_definition_id" json:"fields,omitempty"`
}

type FieldType int8

const (
	FieldTypeString FieldType = 0
	FieldTypeInt    FieldType = 1
	FieldTypeFloat  FieldType = 2
	FieldTypeBool   FieldType = 3
)

type EventDefinitionField struct {
	bun.BaseModel `bun:"table:event_definition_fields,alias:edf"`

	ID                int64     `bun:"id,pk,autoincrement" json:"id"`
	EventDefinitionID int64     `bun:"event_definition_id,notnull" json:"event_definition_id"`
	Key               string    `bun:"key,notnull,type:varchar(64)" json:"key"`
	Type              FieldType `bun:"type,notnull" json:"type"`
	Required          bool      `bun:"required,notnull,default:false" json:"required"`
	MaxLength         int       `bun:"max_length,notnull,default:500" json:"max_length"`
	CreatedAt         time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt         time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp" json:"updated_at"`

	EventDefinition *EventDefinition `bun:"rel:belongs-to,join:event_definition_id=id" json:"event_definition,omitempty"`
}

type EventData struct {
	bun.BaseModel `bun:"table:event_data,alias:evd"`

	ID      int64  `bun:"id,pk,autoincrement" json:"id"`
	EventID int64  `bun:"event_id,notnull" json:"event_id"`
	FieldID int64  `bun:"field_id,notnull" json:"field_id"`
	Value   string `bun:"value,notnull,type:varchar(1024)" json:"value"`

	Event *Event                `bun:"rel:belongs-to,join:event_id=id" json:"event,omitempty"`
	Field *EventDefinitionField `bun:"rel:belongs-to,join:field_id=id" json:"field,omitempty"`
}
