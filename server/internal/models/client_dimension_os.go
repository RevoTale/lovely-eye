package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type ClientOS uint8

const (
	// Persisted analytics enum codes are hard-coded on purpose.
	// Do not reorder these values or switch to iota, because existing rows and migrations depend on them.
	ClientOSUnknown     ClientOS = 0
	ClientOSOther       ClientOS = 1
	ClientOSAndroid     ClientOS = 2
	ClientOSChromeOS    ClientOS = 3
	ClientOSIOS         ClientOS = 4
	ClientOSIPadOS      ClientOS = 5
	ClientOSLinux       ClientOS = 6
	ClientOSMacOS       ClientOS = 7
	ClientOSPlayStation ClientOS = 8
	ClientOSWearOS      ClientOS = 9
	ClientOSWatchOS     ClientOS = 10
	ClientOSWindows     ClientOS = 11
	ClientOSXbox        ClientOS = 12
)

func (o ClientOS) String() string {
	switch o {
	case ClientOSOther:
		return "Other"
	case ClientOSAndroid:
		return "Android"
	case ClientOSChromeOS:
		return "ChromeOS"
	case ClientOSIOS:
		return "iOS"
	case ClientOSIPadOS:
		return "iPadOS"
	case ClientOSLinux:
		return "Linux"
	case ClientOSMacOS:
		return "macOS"
	case ClientOSPlayStation:
		return "PlayStation OS"
	case ClientOSWearOS:
		return "Wear OS"
	case ClientOSWatchOS:
		return "watchOS"
	case ClientOSWindows:
		return "Windows"
	case ClientOSXbox:
		return "Xbox OS"
	default:
		return ""
	}
}

func (o ClientOS) Value() (driver.Value, error) {
	return int64(o), nil
}

func (o *ClientOS) Scan(src any) error {
	return scanClientEnumUint8((*uint8)(o), src)
}

func (o ClientOS) MarshalJSON() ([]byte, error) {
	bytes, err := json.Marshal(o.String())
	if err != nil {
		return nil, fmt.Errorf("marshal client os: %w", err)
	}
	return bytes, nil
}

func ClientOSFromLabel(value string) (ClientOS, bool) {
	switch normalizeClientDimensionLabel(value) {
	case "other":
		return ClientOSOther, true
	case "android":
		return ClientOSAndroid, true
	case "chromeos", "chrome os":
		return ClientOSChromeOS, true
	case "ios":
		return ClientOSIOS, true
	case "ipados":
		return ClientOSIPadOS, true
	case "linux":
		return ClientOSLinux, true
	case "macos":
		return ClientOSMacOS, true
	case "playstation os":
		return ClientOSPlayStation, true
	case "wear os":
		return ClientOSWearOS, true
	case "watchos", "watch os":
		return ClientOSWatchOS, true
	case "windows":
		return ClientOSWindows, true
	case "xbox os":
		return ClientOSXbox, true
	default:
		return ClientOSUnknown, false
	}
}

func ClientOSFromLegacyLabel(value string) ClientOS {
	switch normalizeClientDimensionLabel(value) {
	case "":
		return ClientOSUnknown
	case "mac os x", "os x":
		return ClientOSMacOS
	default:
		if os, ok := ClientOSFromLabel(value); ok {
			return os
		}
		return ClientOSOther
	}
}

func ParseClientOSFilters(values []string) []ClientOS {
	return parseClientDimensionFilters(values, ClientOSFromLabel)
}
