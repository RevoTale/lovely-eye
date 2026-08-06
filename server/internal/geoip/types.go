package geoip

import (
	"errors"
	"time"
)

type State string

const (
	StateDisabled    State = "disabled"
	StateMissing     State = "missing"
	StateDownloading State = "downloading"
	StateReady       State = "ready"
	StateError       State = "error"
)

type Source string

const (
	SourceUnknown     Source = ""
	SourceFile        Source = "file"
	SourceDownloadURL Source = "download-url"
	SourceDBIP        Source = "dbip"
	SourceMaxMind     Source = "maxmind"
)

func (s Source) String() string {
	return string(s)
}

type Config struct {
	DBPath            string
	DownloadURL       string
	MaxMindLicenseKey string
}

type Status struct {
	State     State
	DBPath    string
	Source    Source
	LastError string
	UpdatedAt *time.Time
}

type ListedCountry struct {
	Code string
	Name string
}

type Country struct {
	Name    string
	ISOCode string
}

var ErrNoDBReader = errors.New("no IP reader")

var UnknownCountry = Country{
	Name:    "Unknown",
	ISOCode: "-",
}

var LocalNetworkCountry = Country{
	Name:    "Local Network",
	ISOCode: "-",
}
