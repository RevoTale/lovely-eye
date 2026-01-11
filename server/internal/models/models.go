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

	ID        int64     `bun:"id,pk,autoincrement" json:"id"`
	UserID    int64     `bun:"user_id,notnull" json:"user_id"`
	Domain    string    `bun:"domain,unique,notnull" json:"domain"`
	Name      string    `bun:"name,notnull" json:"name"`
	PublicKey string    `bun:"public_key,unique,notnull" json:"public_key"` // Used in tracking script
	CreatedAt time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp" json:"updated_at"`

	User       *User        `bun:"rel:belongs-to,join:user_id=id" json:"user,omitempty"`
	PageViews  []*PageView  `bun:"rel:has-many,join:id=site_id" json:"page_views,omitempty"`
	Events     []*Event     `bun:"rel:has-many,join:id=site_id" json:"events,omitempty"`
	Sessions   []*Session   `bun:"rel:has-many,join:id=site_id" json:"sessions,omitempty"`
}

// Session represents a visitor session
type Session struct {
	bun.BaseModel `bun:"table:sessions,alias:sess"`

	ID          int64     `bun:"id,pk,autoincrement" json:"id"`
	SiteID      int64     `bun:"site_id,notnull" json:"site_id"`
	VisitorID   string    `bun:"visitor_id,notnull" json:"visitor_id"` // Anonymous hash
	StartedAt   time.Time `bun:"started_at,notnull" json:"started_at"`
	LastSeenAt  time.Time `bun:"last_seen_at,notnull" json:"last_seen_at"`
	EntryPage   string    `bun:"entry_page" json:"entry_page"`
	ExitPage    string    `bun:"exit_page" json:"exit_page"`
	Referrer    string    `bun:"referrer" json:"referrer"`
	UTMSource   string    `bun:"utm_source" json:"utm_source"`
	UTMMedium   string    `bun:"utm_medium" json:"utm_medium"`
	UTMCampaign string    `bun:"utm_campaign" json:"utm_campaign"`
	Country     string    `bun:"country" json:"country"`
	City        string    `bun:"city" json:"city"`
	Device      string    `bun:"device" json:"device"`     // desktop, mobile, tablet
	Browser     string    `bun:"browser" json:"browser"`
	OS          string    `bun:"os" json:"os"`
	ScreenSize  string    `bun:"screen_size" json:"screen_size"`
	PageViews   int       `bun:"page_views,default:0" json:"page_views"`
	Duration    int       `bun:"duration,default:0" json:"duration"` // seconds
	IsBounce    bool      `bun:"is_bounce,default:true" json:"is_bounce"`

	Site *Site `bun:"rel:belongs-to,join:site_id=id" json:"site,omitempty"`
}

// PageView represents a single page view
type PageView struct {
	bun.BaseModel `bun:"table:page_views,alias:pv"`

	ID        int64     `bun:"id,pk,autoincrement" json:"id"`
	SiteID    int64     `bun:"site_id,notnull" json:"site_id"`
	SessionID int64     `bun:"session_id" json:"session_id"`
	VisitorID string    `bun:"visitor_id,notnull" json:"visitor_id"`
	Path      string    `bun:"path,notnull" json:"path"`
	Title     string    `bun:"title" json:"title"`
	Referrer  string    `bun:"referrer" json:"referrer"`
	Duration  int       `bun:"duration,default:0" json:"duration"` // time on page in seconds
	CreatedAt time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp" json:"created_at"`

	Site    *Site    `bun:"rel:belongs-to,join:site_id=id" json:"site,omitempty"`
	Session *Session `bun:"rel:belongs-to,join:session_id=id" json:"session,omitempty"`
}

// Event represents a custom event (button click, form submission, etc.)
type Event struct {
	bun.BaseModel `bun:"table:events,alias:e"`

	ID         int64     `bun:"id,pk,autoincrement" json:"id"`
	SiteID     int64     `bun:"site_id,notnull" json:"site_id"`
	SessionID  int64     `bun:"session_id" json:"session_id"`
	VisitorID  string    `bun:"visitor_id,notnull" json:"visitor_id"`
	Name       string    `bun:"name,notnull" json:"name"`
	Path       string    `bun:"path" json:"path"`
	Properties string    `bun:"properties" json:"properties"` // JSON string
	CreatedAt  time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp" json:"created_at"`

	Site    *Site    `bun:"rel:belongs-to,join:site_id=id" json:"site,omitempty"`
	Session *Session `bun:"rel:belongs-to,join:session_id=id" json:"session,omitempty"`
}

// DailyStats represents aggregated daily statistics for a site
type DailyStats struct {
	bun.BaseModel `bun:"table:daily_stats,alias:ds"`

	ID              int64     `bun:"id,pk,autoincrement" json:"id"`
	SiteID          int64     `bun:"site_id,notnull" json:"site_id"`
	Date            time.Time `bun:"date,notnull" json:"date"`
	Visitors        int       `bun:"visitors,default:0" json:"visitors"`
	PageViews       int       `bun:"page_views,default:0" json:"page_views"`
	Sessions        int       `bun:"sessions,default:0" json:"sessions"`
	BounceRate      float64   `bun:"bounce_rate,default:0" json:"bounce_rate"`
	AvgDuration     float64   `bun:"avg_duration,default:0" json:"avg_duration"`

	Site *Site `bun:"rel:belongs-to,join:site_id=id" json:"site,omitempty"`
}
