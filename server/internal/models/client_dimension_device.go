package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
)

type ClientDevice uint8

const (
	// Persisted analytics enum codes are hard-coded on purpose.
	// Do not reorder these values or switch to iota, because existing rows and migrations depend on them.
	ClientDeviceUnknown ClientDevice = 0
	ClientDeviceConsole ClientDevice = 1
	ClientDeviceDesktop ClientDevice = 2
	ClientDeviceMobile  ClientDevice = 3
	ClientDeviceOther   ClientDevice = 4
	ClientDeviceSmartTV ClientDevice = 5
	ClientDeviceTablet  ClientDevice = 6
	ClientDeviceWatch   ClientDevice = 7
)

func (d ClientDevice) String() string {
	switch d {
	case ClientDeviceConsole:
		return "console"
	case ClientDeviceDesktop:
		return "desktop"
	case ClientDeviceMobile:
		return "mobile"
	case ClientDeviceOther:
		return "other"
	case ClientDeviceSmartTV:
		return "smart-tv"
	case ClientDeviceTablet:
		return "tablet"
	case ClientDeviceWatch:
		return "watch"
	default:
		return ""
	}
}

func (d ClientDevice) Value() (driver.Value, error) {
	return int64(d), nil
}

func (d *ClientDevice) Scan(src any) error {
	return scanClientEnumUint8((*uint8)(d), src)
}

func (d ClientDevice) MarshalJSON() ([]byte, error) {
	bytes, err := json.Marshal(d.String())
	if err != nil {
		return nil, fmt.Errorf("marshal client device: %w", err)
	}
	return bytes, nil
}

func ClientDeviceFromLabel(value string) (ClientDevice, bool) {
	switch normalizeClientDimensionLabel(value) {
	case "console":
		return ClientDeviceConsole, true
	case "desktop":
		return ClientDeviceDesktop, true
	case "mobile":
		return ClientDeviceMobile, true
	case "other":
		return ClientDeviceOther, true
	case "smart-tv", "smart tv", "smarttv":
		return ClientDeviceSmartTV, true
	case "tablet":
		return ClientDeviceTablet, true
	case "watch":
		return ClientDeviceWatch, true
	default:
		return ClientDeviceUnknown, false
	}
}

func ClientDeviceFromLegacyLabel(value string) ClientDevice {
	if device, ok := ClientDeviceFromLabel(value); ok {
		return device
	}
	if strings.TrimSpace(value) == "" {
		return ClientDeviceUnknown
	}
	return ClientDeviceOther
}

func ParseClientDeviceFilters(values []string) []ClientDevice {
	return parseClientDimensionFilters(values, ClientDeviceFromLabel)
}
