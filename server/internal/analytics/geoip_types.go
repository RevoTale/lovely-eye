package analytics

import "github.com/lovely-eye/server/internal/geoip"

const (
	geoIPStateDisabled = geoip.StateDisabled
	geoIPStateReady    = geoip.StateReady
	geoIPStateError    = geoip.StateError
)

type GeoIPStatus = geoip.Status
type GeoIPCountry = geoip.ListedCountry
type Country = geoip.Country

var ErrNoDBReader = geoip.ErrNoDBReader
var UnknownCountry = geoip.UnknownCountry
var LocalNetworkCountry = geoip.LocalNetworkCountry
