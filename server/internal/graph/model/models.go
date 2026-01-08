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
	ID        string    `json:"id"`
	Domain    string    `json:"domain"`
	Name      string    `json:"name"`
	PublicKey string    `json:"publicKey"`
	CreatedAt time.Time `json:"createdAt"`
}

type AuthPayload struct {
	User         *User  `json:"user"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type TokenPayload struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type DashboardStats struct {
	Visitors     int              `json:"visitors"`
	PageViews    int              `json:"pageViews"`
	Sessions     int              `json:"sessions"`
	BounceRate   float64          `json:"bounceRate"`
	AvgDuration  float64          `json:"avgDuration"`
	TopPages     []*PageStats     `json:"topPages"`
	TopReferrers []*ReferrerStats `json:"topReferrers"`
	Browsers     []*BrowserStats  `json:"browsers"`
	Devices      []*DeviceStats   `json:"devices"`
	Countries    []*CountryStats  `json:"countries"`
	DailyStats   []*DailyStats    `json:"dailyStats"`
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
	Visitors int `json:"visitors"`
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
	Domain string `json:"domain"`
	Name   string `json:"name"`
}

type UpdateSiteInput struct {
	Name string `json:"name"`
}

type DateRangeInput struct {
	From *time.Time `json:"from,omitempty"`
	To   *time.Time `json:"to,omitempty"`
}

type Event struct {
	ID         string           `json:"id"`
	Name       string           `json:"name"`
	Path       string           `json:"path"`
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
