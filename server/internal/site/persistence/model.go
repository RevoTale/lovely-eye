package persistence

import (
	"time"

	"github.com/uptrace/bun"
)

type Site struct {
	bun.BaseModel `bun:"table:sites,alias:s"`

	ID           int64     `bun:"id,pk,autoincrement"`
	UserID       int64     `bun:"user_id,notnull"`
	Name         string    `bun:"name,notnull"`
	PublicKey    string    `bun:"public_key,unique,notnull"`
	TrackCountry bool      `bun:"track_country,notnull,default:false"`
	CreatedAt    time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt    time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp"`

	Domains          []*Domain         `bun:"rel:has-many,join:id=site_id"`
	BlockedIPs       []*BlockedIP      `bun:"rel:has-many,join:id=site_id"`
	BlockedCountries []*BlockedCountry `bun:"rel:has-many,join:id=site_id"`
}

type Domain struct {
	bun.BaseModel `bun:"table:site_domains,alias:sd"`

	ID        int64     `bun:"id,pk,autoincrement"`
	SiteID    int64     `bun:"site_id,notnull,unique:site_domains_site_id_domain"`
	Domain    string    `bun:"domain,notnull,unique:site_domains_site_id_domain"`
	Position  int       `bun:"position,notnull,default:0"`
	CreatedAt time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
}

type BlockedIP struct {
	bun.BaseModel `bun:"table:site_blocked_ips,alias:sbi"`

	ID        int64     `bun:"id,pk,autoincrement"`
	SiteID    int64     `bun:"site_id,notnull"`
	IP        string    `bun:"ip,notnull"`
	CreatedAt time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
}

type BlockedCountry struct {
	bun.BaseModel `bun:"table:site_blocked_countries,alias:sbc"`

	ID          int64     `bun:"id,pk,autoincrement"`
	SiteID      int64     `bun:"site_id,notnull"`
	CountryCode string    `bun:"country_code,notnull"`
	CreatedAt   time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt   time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
}
