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
	Role         string    `bun:"role,notnull,default:'user'" json:"role"` // admin, user
	Email        string    `bun:"email" json:"email"`                      // New field for testing
	CreatedAt    time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt    time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp" json:"updated_at"`

	Sites []*Site `bun:"rel:has-many,join:id=user_id" json:"sites,omitempty"`
}

// Site represents a website being tracked
type Site struct {
	bun.BaseModel `bun:"table:sites,alias:s"`

	ID           int64     `bun:"id,pk,autoincrement" json:"id"`
	UserID       int64     `bun:"user_id,notnull" json:"user_id"`
	Name         string    `bun:"name,notnull" json:"name"`
	PublicKey    string    `bun:"public_key,unique,notnull" json:"public_key"` // Used in tracking script
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

// SiteDomain represents an allowed domain for a site
type SiteDomain struct {
	bun.BaseModel `bun:"table:site_domains,alias:sd"`

	ID        int64     `bun:"id,pk,autoincrement" json:"id"`
	SiteID    int64     `bun:"site_id,notnull" json:"site_id"`
	Domain    string    `bun:"domain,unique,notnull" json:"domain"`
	Position  int       `bun:"position,notnull,default:0" json:"position"`
	CreatedAt time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp" json:"updated_at"`

	Site *Site `bun:"rel:belongs-to,join:site_id=id" json:"site,omitempty"`
}

// SiteBlockedIP represents a blocked IP for a site
type SiteBlockedIP struct {
	bun.BaseModel `bun:"table:site_blocked_ips,alias:sbi"`

	ID        int64     `bun:"id,pk,autoincrement" json:"id"`
	SiteID    int64     `bun:"site_id,notnull" json:"site_id"`
	IP        string    `bun:"ip,notnull" json:"ip"`
	CreatedAt time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp" json:"updated_at"`

	Site *Site `bun:"rel:belongs-to,join:site_id=id" json:"site,omitempty"`
}

// SiteBlockedCountry represents a blocked country for a site
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
	SiteID     int64  `bun:"site_id,notnull" json:"site_id"`                  // FK to Site
	Hash       string `bun:"hash,notnull,type:varchar(64)" json:"hash"`       // SHA-256 hex, from current architecture
	Country    string `bun:"country,type:varchar(2)" json:"country"`          // ISO 3166-1 alpha-2, only if site.track_country enabled
	Device     string `bun:"device,type:varchar(10)" json:"device"`           // "desktop", "mobile", "tablet" → Device Breakdown widget
	Browser    string `bun:"browser,type:varchar(32)" json:"browser"`         // "Chrome", "Firefox", etc. → Browser Stats widget
	OS         string `bun:"os,type:varchar(32)" json:"os"`                   // "Windows", "macOS", etc. → OS Stats widget
	ScreenSize string `bun:"screen_size,type:varchar(16)" json:"screen_size"` // "1920x1080" → Screen Size widget

	Site     *Site      `bun:"rel:belongs-to,join:site_id=id" json:"site,omitempty"`
	Sessions []*Session `bun:"rel:has-many,join:id=client_id" json:"sessions,omitempty"`
}

// Session represents a single visit by a Client.
// Stores PER-VISIT attributes (can change between sessions).
// Session timeout: 30min inactivity = new session.
type Session struct {
	bun.BaseModel `bun:"table:sessions,alias:s"`

	ID       int64 `bun:"id,pk,autoincrement" json:"id"`
	SiteID   int64 `bun:"site_id,notnull" json:"site_id"`   // denormalized from Client for query perf
	ClientID int64 `bun:"client_id,notnull" json:"client_id"` // FK to Client

	// Entry (set once when session starts)
	EnterTime int64  `bun:"enter_time,notnull" json:"enter_time"`                    // unix seconds
	EnterHour int64  `bun:"enter_hour,notnull" json:"enter_hour"`                    // unix / 3600 → hourly charts
	EnterDay  int64  `bun:"enter_day,notnull" json:"enter_day"`                      // unix / 86400 → daily charts
	EnterPath string `bun:"enter_path,notnull,type:varchar(2048)" json:"enter_path"` // landing page → Top Entry Pages

	// Exit (updated on each pageview/event)
	ExitTime int64  `bun:"exit_time,notnull" json:"exit_time"`                   // unix seconds
	ExitHour int64  `bun:"exit_hour,notnull" json:"exit_hour"`                   // unix / 3600
	ExitDay  int64  `bun:"exit_day,notnull" json:"exit_day"`                     // unix / 86400
	ExitPath string `bun:"exit_path,notnull,type:varchar(2048)" json:"exit_path"` // last page → Top Exit Pages

	// Attribution (per-session, not per-client)
	Referrer    string `bun:"referrer,type:varchar(2048)" json:"referrer"`       // external source → Top Referrers
	UTMSource   string `bun:"utm_source,type:varchar(128)" json:"utm_source"`    // utm_source param
	UTMMedium   string `bun:"utm_medium,type:varchar(128)" json:"utm_medium"`    // utm_medium param
	UTMCampaign string `bun:"utm_campaign,type:varchar(256)" json:"utm_campaign"` // utm_campaign param

	// Metrics (computed/updated during session)
	Duration      int `bun:"duration,notnull,default:0" json:"duration"`             // ExitTime - EnterTime (sec) → Avg Session Duration
	PageViewCount int `bun:"page_view_count,notnull,default:0" json:"page_view_count"` // incremented on pageview → Bounce Rate (count==1), Avg Pages/Session

	Site   *Site    `bun:"rel:belongs-to,join:site_id=id" json:"site,omitempty"`
	Client *Client  `bun:"rel:belongs-to,join:client_id=id" json:"client,omitempty"`
	Events []*Event `bun:"rel:has-many,join:id=session_id" json:"events,omitempty"`
}

// EventType enum for unified Event table
type EventType int8

const (
	EventTypePageview EventType = 0 // page view
	EventTypeCustom   EventType = 1 // custom event (button click, form submit, etc.)
)

// Event represents unified table for pageviews and custom events.
// Pageview: Type=0, Name=page title, DefinitionID=null
// Custom:   Type=1, Name=event name, DefinitionID=FK
// Unified enables: single query for all activity, consistent time bucketing.
type Event struct {
	bun.BaseModel `bun:"table:events,alias:e"`

	ID        int64 `bun:"id,pk,autoincrement" json:"id"`
	SessionID int64 `bun:"session_id,notnull" json:"session_id"` // FK to Session

	// Time buckets (same strategy as Session)
	Time int64 `bun:"time,notnull" json:"time"` // unix seconds
	Hour int64 `bun:"hour,notnull" json:"hour"` // unix / 3600 → hourly charts
	Day  int64 `bun:"day,notnull" json:"day"`   // unix / 86400 → daily charts

	// Data
	Path string    `bun:"path,notnull,type:varchar(2048)" json:"path"` // page URL → Top Pages
	Name string    `bun:"name,notnull,type:varchar(256)" json:"name"`  // page title (pageview) or event name (custom)
	Type EventType `bun:"type,notnull" json:"type"`                    // 0=pageview, 1=custom (int8 for storage efficiency)

	DefinitionID *int64 `bun:"definition_id" json:"definition_id"` // FK to EventDefinition (null for pageviews)

	Session    *Session         `bun:"rel:belongs-to,join:session_id=id" json:"session,omitempty"`
	Definition *EventDefinition `bun:"rel:belongs-to,join:definition_id=id" json:"definition,omitempty"`
	Data       []*EventData     `bun:"rel:has-many,join:id=event_id" json:"data,omitempty"`
}

// EventDefinition represents an allowlisted event name for a site
type EventDefinition struct {
	bun.BaseModel `bun:"table:event_definitions,alias:ed"`

	ID        int64     `bun:"id,pk,autoincrement" json:"id"`
	SiteID    int64     `bun:"site_id,notnull" json:"site_id"`
	Name      string    `bun:"name,notnull" json:"name"`
	CreatedAt time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp" json:"updated_at"`

	Site   *Site                   `bun:"rel:belongs-to,join:site_id=id" json:"site,omitempty"`
	Fields []*EventDefinitionField `bun:"rel:has-many,join:id=event_definition_id" json:"fields,omitempty"`
}

// FieldType enum for event property types
type FieldType int8

const (
	FieldTypeString FieldType = 0
	FieldTypeInt    FieldType = 1
	FieldTypeFloat  FieldType = 2
	FieldTypeBool   FieldType = 3
)

// EventDefinitionField describes allowed properties for an event.
// Property definition for custom events.
// Specifies schema for EventData values: key name, type, constraints.
// Enables typed, validated properties instead of arbitrary JSON.
type EventDefinitionField struct {
	bun.BaseModel `bun:"table:event_definition_fields,alias:edf"`

	ID                int64     `bun:"id,pk,autoincrement" json:"id"`
	EventDefinitionID int64     `bun:"event_definition_id,notnull" json:"event_definition_id"`
	Key               string    `bun:"key,notnull,type:varchar(64)" json:"key"`      // property name
	Type              FieldType `bun:"type,notnull" json:"type"`                     // 0=string, 1=int, 2=float, 3=bool (int8)
	Required          bool      `bun:"required,notnull,default:false" json:"required"`
	MaxLength         int       `bun:"max_length,notnull,default:500" json:"max_length"` // for string type, 0 = no limit
	CreatedAt         time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt         time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp" json:"updated_at"`

	EventDefinition *EventDefinition `bun:"rel:belongs-to,join:event_definition_id=id" json:"event_definition,omitempty"`
}

// EventData represents property values for custom events.
// Normalized storage: one row per property (not JSON blob).
// Enables: SQL queries on properties, type safety, smaller storage.
type EventData struct {
	bun.BaseModel `bun:"table:event_data,alias:evd"`

	ID      int64  `bun:"id,pk,autoincrement" json:"id"`
	EventID int64  `bun:"event_id,notnull" json:"event_id"`                  // FK to Event
	FieldID int64  `bun:"field_id,notnull" json:"field_id"`                  // FK to EventDefinitionField
	Value   string `bun:"value,notnull,type:varchar(1024)" json:"value"` // stored as string, typed by FieldType

	Event *Event                `bun:"rel:belongs-to,join:event_id=id" json:"event,omitempty"`
	Field *EventDefinitionField `bun:"rel:belongs-to,join:field_id=id" json:"field,omitempty"`
}
