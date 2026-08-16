package analytics

import "time"

type EventType string

const (
	EventTypePageView   EventType = "PAGE_VIEW"
	EventTypePredefined EventType = "PREDEFINED"
)

type TimeBucket string

const (
	TimeBucketDaily  TimeBucket = "daily"
	TimeBucketHourly TimeBucket = "hourly"
)

type Filter struct {
	Referrer           []string
	Browser            []string
	Device             []string
	OS                 []string
	Page               []string
	Country            []string
	EventTypes         []EventType
	EventName          []string
	EventPath          []string
	EventDefinitionIDs []int64
}

type Query struct {
	SiteID int64
	From   time.Time
	To     time.Time
	Limit  int
	Offset int
	Bucket TimeBucket
	Filter Filter
}

type Overview struct {
	Visitors    int
	PageViews   int
	Sessions    int
	BounceRate  float64
	AvgDuration float64
}

type PageStats struct {
	Path     string
	Views    int
	Visitors int
}

type ReferrerStats struct {
	Referrer string
	Visitors int
}

type BrowserStats struct {
	Browser  string
	Visitors int
}

type DeviceStats struct {
	Device   string
	Visitors int
}

type OperatingSystemStats struct {
	OS       string
	Visitors int
}

type CountryStats struct {
	CountryCode string
	Visitors    int
}

type TimeSeriesStats struct {
	DateBucket int64
	Visitors   int
	PageViews  int
	Sessions   int
}

type ActivePageStats struct {
	Path     string
	Visitors int
}

type Event struct {
	ID           int64
	SessionID    int64
	Time         int64
	Hour         int64
	Day          int64
	Path         string
	DefinitionID *int64
	Definition   *EventDefinition
	Data         []*EventData
}

type EventDefinition struct {
	ID        int64
	SiteID    int64
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	Fields    []*EventField
}

type EventFieldType int8

const (
	EventFieldTypeString EventFieldType = iota
	EventFieldTypeInt
	EventFieldTypeFloat
	EventFieldTypeBool
)

type EventField struct {
	ID                int64
	EventDefinitionID int64
	Key               string
	Type              EventFieldType
	Required          bool
	MaxLength         int
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type EventData struct {
	ID      int64
	EventID int64
	FieldID int64
	Value   string
	Field   *EventField
}

type EventCount struct {
	Event *Event
	Count int
}
