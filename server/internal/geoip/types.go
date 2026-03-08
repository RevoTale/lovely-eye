package geoip

import (
	"errors"
	"time"
)

const (
	StateDisabled    = "disabled"
	StateMissing     = "missing"
	StateDownloading = "downloading"
	StateReady       = "ready"
	StateError       = "error"
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
	State     string
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
