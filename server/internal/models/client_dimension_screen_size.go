package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type ClientScreenSize uint8

const (
	// Persisted analytics enum codes are hard-coded on purpose.
	// Do not reorder these values or switch to iota, because existing rows and migrations depend on them.
	ClientScreenSizeUnknown ClientScreenSize = 0
	ClientScreenSizeWatch   ClientScreenSize = 1
	ClientScreenSizeXS      ClientScreenSize = 2
	ClientScreenSizeSM      ClientScreenSize = 3
	ClientScreenSizeMD      ClientScreenSize = 4
	ClientScreenSizeLG      ClientScreenSize = 5
	ClientScreenSizeXL      ClientScreenSize = 6
)

func (s ClientScreenSize) String() string {
	switch s {
	case ClientScreenSizeWatch:
		return "watch"
	case ClientScreenSizeXS:
		return "xs"
	case ClientScreenSizeSM:
		return "sm"
	case ClientScreenSizeMD:
		return "md"
	case ClientScreenSizeLG:
		return "lg"
	case ClientScreenSizeXL:
		return "xl"
	default:
		return ""
	}
}

func (s ClientScreenSize) Value() (driver.Value, error) {
	return int64(s), nil
}

func (s *ClientScreenSize) Scan(src any) error {
	return scanClientEnumUint8((*uint8)(s), src)
}

func (s ClientScreenSize) MarshalJSON() ([]byte, error) {
	bytes, err := json.Marshal(s.String())
	if err != nil {
		return nil, fmt.Errorf("marshal client screen size: %w", err)
	}
	return bytes, nil
}

func ClientScreenSizeFromWidth(width int) ClientScreenSize {
	switch {
	case width > 0 && width < 320:
		return ClientScreenSizeWatch
	case width < 576:
		return ClientScreenSizeXS
	case width < 768:
		return ClientScreenSizeSM
	case width < 992:
		return ClientScreenSizeMD
	case width < 1200:
		return ClientScreenSizeLG
	default:
		return ClientScreenSizeXL
	}
}

func ClientScreenSizeFromLegacyLabel(value string) ClientScreenSize {
	switch normalizeClientDimensionLabel(value) {
	case "":
		return ClientScreenSizeUnknown
	case "watch":
		return ClientScreenSizeWatch
	case "xs":
		return ClientScreenSizeXS
	case "sm":
		return ClientScreenSizeSM
	case "md":
		return ClientScreenSizeMD
	case "lg":
		return ClientScreenSizeLG
	case "xl":
		return ClientScreenSizeXL
	default:
		width, ok := parseLegacyScreenWidth(value)
		if !ok {
			return ClientScreenSizeUnknown
		}
		return ClientScreenSizeFromWidth(width)
	}
}
