package models

import (
	"fmt"
	"strconv"
	"strings"
)

func normalizeClientDimensionLabel(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

func parseClientDimensionFilters[T comparable](values []string, parse func(string) (T, bool)) []T {
	if len(values) == 0 {
		return nil
	}

	seen := make(map[T]struct{}, len(values))
	parsed := make([]T, 0, len(values))
	for _, value := range values {
		enumValue, ok := parse(value)
		if !ok {
			continue
		}
		if _, exists := seen[enumValue]; exists {
			continue
		}
		seen[enumValue] = struct{}{}
		parsed = append(parsed, enumValue)
	}
	return parsed
}

func parseLegacyScreenWidth(value string) (int, bool) {
	widthPart := strings.TrimSpace(value)
	if widthPart == "" {
		return 0, false
	}
	if separator := strings.Index(widthPart, "x"); separator >= 0 {
		widthPart = widthPart[:separator]
	}
	width, err := strconv.Atoi(strings.TrimSpace(widthPart))
	if err != nil {
		return 0, false
	}
	return width, true
}

func scanClientEnumUint8(target *uint8, src any) error {
	switch value := src.(type) {
	case nil:
		*target = 0
		return nil
	case int64:
		return setScannedClientEnum(target, value)
	case int32:
		return setScannedClientEnum(target, int64(value))
	case int16:
		return setScannedClientEnum(target, int64(value))
	case int8:
		return setScannedClientEnum(target, int64(value))
	case int:
		return setScannedClientEnum(target, int64(value))
	case uint64:
		return setScannedClientEnumUnsigned(target, value)
	case uint32:
		return setScannedClientEnumUnsigned(target, uint64(value))
	case uint16:
		return setScannedClientEnumUnsigned(target, uint64(value))
	case uint8:
		*target = value
		return nil
	case []byte:
		return scanClientEnumUint8Bytes(target, value)
	case string:
		return scanClientEnumUint8Bytes(target, []byte(value))
	default:
		return fmt.Errorf("unsupported enum source type %T", src)
	}
}

func scanClientEnumUint8Bytes(target *uint8, value []byte) error {
	trimmed := strings.TrimSpace(string(value))
	if trimmed == "" {
		*target = 0
		return nil
	}
	parsed, err := strconv.ParseUint(trimmed, 10, 8)
	if err != nil {
		return fmt.Errorf("parse enum value %q: %w", trimmed, err)
	}
	*target = uint8(parsed)
	return nil
}

func setScannedClientEnum(target *uint8, value int64) error {
	if value < 0 || value > 255 {
		return fmt.Errorf("enum value %d out of uint8 range", value)
	}
	*target = uint8(value)
	return nil
}

func setScannedClientEnumUnsigned(target *uint8, value uint64) error {
	if value > 255 {
		return fmt.Errorf("enum value %d out of uint8 range", value)
	}
	*target = uint8(value)
	return nil
}
