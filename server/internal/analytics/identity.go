package analytics

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/netip"
	"strings"
	"time"

	analyticspersistence "github.com/lovely-eye/server/internal/analytics/persistence"
	"github.com/lovely-eye/server/internal/event"
	"golang.org/x/crypto/hkdf"
)

const unknownVisitorIPPrefix = "unknown"

const defaultEventPropertyMaxLength = 500

// generateVisitorID creates the site-scoped daily UTC hash used by the
// UTC-day-skipped client rotation helper. Client reuse compares today's and
// yesterday's hash to preserve continuity across adjacent UTC days only.
func (s *Service) generateVisitorID(
	siteID int64,
	ip string,
	browser analyticspersistence.ClientBrowser,
	device analyticspersistence.ClientDevice,
	now time.Time,
) string {
	dateBucket := now.UTC().Format("2006-01-02")
	key := s.deriveVisitorIdentityKey(siteID, dateBucket)
	ipPrefix := truncateVisitorIPPrefix(ip)
	data := fmt.Sprintf("%d|%s|%s|%s", siteID, ipPrefix, browser.String(), device.String())

	mac := hmac.New(sha256.New, key)
	_, _ = mac.Write([]byte(data))

	sum := mac.Sum(nil)
	return hex.EncodeToString(sum[:16])
}

func (s *Service) deriveVisitorIdentityKey(siteID int64, dateBucket string) []byte {
	info := fmt.Appendf(nil, "analytics:%d:%s", siteID, dateBucket)
	reader := hkdf.New(sha256.New, s.identitySecret, nil, info)
	key := make([]byte, sha256.Size)
	_, _ = io.ReadFull(reader, key)
	return key
}

func truncateVisitorIPPrefix(ip string) string {
	addr, err := netip.ParseAddr(strings.TrimSpace(ip))
	if err != nil {
		return unknownVisitorIPPrefix
	}

	if addr.Is4In6() {
		addr = netip.AddrFrom4(addr.As4())
	}

	if addr.Is4() {
		return netip.PrefixFrom(addr, 24).Masked().String()
	}

	return netip.PrefixFrom(addr, 64).Masked().String()
}

func sanitizeEventProperties(propsJSON string, fields []*event.Field) (string, bool, error) {
	if len(fields) == 0 {
		return "", true, nil
	}

	allowed := make(map[string]*event.Field, len(fields))
	for _, field := range fields {
		allowed[field.Key] = field
	}

	var props map[string]interface{}
	if propsJSON != "" {
		if err := json.Unmarshal([]byte(propsJSON), &props); err != nil {
			return "", false, fmt.Errorf("unmarshal properties json: %w", err)
		}
	} else {
		props = map[string]interface{}{}
	}

	sanitized := make(map[string]interface{})
	for key, value := range props {
		field, ok := allowed[key]
		if !ok {
			continue
		}
		sanitizedValue, ok := sanitizeEventPropertyValue(field, value)
		if !ok {
			return "", false, nil
		}
		sanitized[key] = sanitizedValue
	}

	for _, field := range fields {
		if field.Required {
			if _, ok := sanitized[field.Key]; !ok {
				return "", false, nil
			}
		}
	}

	if len(sanitized) == 0 {
		return "", true, nil
	}

	bytes, err := json.Marshal(sanitized)
	if err != nil {
		return "", false, fmt.Errorf("marshal sanitized properties: %w", err)
	}
	return string(bytes), true, nil
}

func sanitizeEventPropertyValue(field *event.Field, value interface{}) (interface{}, bool) {
	switch field.Type {
	case event.FieldTypeString:
		return sanitizeStringEventProperty(value, field.MaxLength)
	case event.FieldTypeInt:
		return sanitizeIntEventProperty(value)
	case event.FieldTypeFloat:
		numberValue, ok := value.(float64)
		return numberValue, ok
	case event.FieldTypeBool:
		boolValue, ok := value.(bool)
		return boolValue, ok
	default:
		return nil, false
	}
}

func sanitizeStringEventProperty(value interface{}, maxLength int) (string, bool) {
	strValue, ok := value.(string)
	if !ok {
		return "", false
	}
	if maxLength <= 0 {
		maxLength = defaultEventPropertyMaxLength
	}
	if len(strValue) > maxLength {
		strValue = strValue[:maxLength]
	}
	return strValue, true
}

func sanitizeIntEventProperty(value interface{}) (int64, bool) {
	switch v := value.(type) {
	case float64:
		return int64(v), true
	case int:
		return int64(v), true
	case int64:
		return v, true
	default:
		return 0, false
	}
}
