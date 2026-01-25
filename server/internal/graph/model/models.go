package model

import "time"

type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
	Sites     []*Site   `json:"sites,omitempty"`
}

type Site struct {
	ID               string    `json:"id"`
	Domains          []string  `json:"domains"`
	Name             string    `json:"name"`
	PublicKey        string    `json:"publicKey"`
	TrackCountry     bool      `json:"trackCountry"`
	BlockedIPs       []string  `json:"blockedIPs"`
	BlockedCountries []string  `json:"blockedCountries"`
	CreatedAt        time.Time `json:"createdAt"`
}

type AuthPayload struct {
	User *User `json:"user"`
	// Tokens are set as HttpOnly cookies, not returned in response
	// See: https://www.reddit.com/r/node/comments/1im7yj0/comment/mc0ylfd/
}

type TokenPayload struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type PageStats struct {
	Path     string `json:"path"`
	Views    int    `json:"views"`
	Visitors int    `json:"visitors"`
}

type ReferrerStats struct {
	Referrer string `json:"referrer"`
	Visitors int    `json:"visitors"`
}

type BrowserStats struct {
	Browser  string `json:"browser"`
	Visitors int    `json:"visitors"`
}

type DeviceStats struct {
	Device   string `json:"device"`
	Visitors int    `json:"visitors"`
}

type CountryStats struct {
	Country  string `json:"country"`
	Visitors int    `json:"visitors"`
}

type DailyStats struct {
	Date      time.Time `json:"date"`
	Visitors  int       `json:"visitors"`
	PageViews int       `json:"pageViews"`
	Sessions  int       `json:"sessions"`
}

type RealtimeStats struct {
	Visitors int   `json:"visitors"`
	SiteID   int64 `json:"-"`
}

type RegisterInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type CreateSiteInput struct {
	Domains []string `json:"domains"`
	Name    string   `json:"name"`
}

type UpdateSiteInput struct {
	Name             string   `json:"name"`
	TrackCountry     *bool    `json:"trackCountry,omitempty"`
	Domains          []string `json:"domains,omitempty"`
	BlockedIPs       []string `json:"blockedIPs,omitempty"`
	BlockedCountries []string `json:"blockedCountries,omitempty"`
}

type DateRangeInput struct {
	From *time.Time `json:"from,omitempty"`
	To   *time.Time `json:"to,omitempty"`
}

type Event struct {
	ID         string           `json:"id"`
	Name       string           `json:"name"`
	Path       string           `json:"path"`
	Definition *EventDefinition `json:"definition,omitempty"`
	Properties []*EventProperty `json:"properties"`
	CreatedAt  time.Time        `json:"createdAt"`
}

type EventProperty struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type EventsResult struct {
	Events []*Event `json:"events"`
	Total  int      `json:"total"`
}

type EventFieldType string

const (
	EventFieldTypeString  EventFieldType = "STRING"
	EventFieldTypeInt     EventFieldType = "INT"
	EventFieldTypeBoolean EventFieldType = "BOOLEAN"
)

type EventType string

const (
	EventTypePageView   EventType = "PAGE_VIEW"
	EventTypePredefined EventType = "PREDEFINED"
)

type EventDefinitionField struct {
	ID        string         `json:"id"`
	Key       string         `json:"key"`
	Type      EventFieldType `json:"type"`
	Required  bool           `json:"required"`
	MaxLength int            `json:"maxLength"`
}

type EventDefinition struct {
	ID        string                  `json:"id"`
	Name      string                  `json:"name"`
	Fields    []*EventDefinitionField `json:"fields"`
	CreatedAt time.Time               `json:"createdAt"`
	UpdatedAt time.Time               `json:"updatedAt"`
}

type EventDefinitionFieldInput struct {
	Key       string         `json:"key"`
	Type      EventFieldType `json:"type"`
	Required  bool           `json:"required"`
	MaxLength *int           `json:"maxLength,omitempty"`
}

type EventDefinitionInput struct {
	Name   string                       `json:"name"`
	Fields []*EventDefinitionFieldInput `json:"fields"`
}
