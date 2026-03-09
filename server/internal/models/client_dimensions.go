package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
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

type ClientBrowser uint8

const (
	// Persisted analytics enum codes are hard-coded on purpose.
	// Do not reorder these values or switch to iota, because existing rows and migrations depend on them.
	ClientBrowserUnknown          ClientBrowser = 0
	ClientBrowserOther            ClientBrowser = 1
	ClientBrowserAndroidWebView   ClientBrowser = 2
	ClientBrowserChrome           ClientBrowser = 3
	ClientBrowserDuckDuckGo       ClientBrowser = 4
	ClientBrowserEdge             ClientBrowser = 5
	ClientBrowserFacebookInApp    ClientBrowser = 6
	ClientBrowserFirefox          ClientBrowser = 7
	ClientBrowserInstagramInApp   ClientBrowser = 8
	ClientBrowserInternetExplorer ClientBrowser = 9
	ClientBrowserMIUI             ClientBrowser = 10
	ClientBrowserOpera            ClientBrowser = 11
	ClientBrowserPlayStation      ClientBrowser = 12
	ClientBrowserSafari           ClientBrowser = 13
	ClientBrowserSamsungInternet  ClientBrowser = 14
	ClientBrowserUCBrowser        ClientBrowser = 15
	ClientBrowserVivaldi          ClientBrowser = 16
	ClientBrowserXbox             ClientBrowser = 17
	ClientBrowserYandex           ClientBrowser = 18
)

func (b ClientBrowser) String() string {
	switch b {
	case ClientBrowserOther:
		return "Other"
	case ClientBrowserAndroidWebView:
		return "Android WebView"
	case ClientBrowserChrome:
		return "Chrome"
	case ClientBrowserDuckDuckGo:
		return "DuckDuckGo"
	case ClientBrowserEdge:
		return "Edge"
	case ClientBrowserFacebookInApp:
		return "Facebook In-App Browser"
	case ClientBrowserFirefox:
		return "Firefox"
	case ClientBrowserInstagramInApp:
		return "Instagram In-App Browser"
	case ClientBrowserInternetExplorer:
		return "Internet Explorer"
	case ClientBrowserMIUI:
		return "MIUI Browser"
	case ClientBrowserOpera:
		return "Opera"
	case ClientBrowserPlayStation:
		return "PlayStation Browser"
	case ClientBrowserSafari:
		return "Safari"
	case ClientBrowserSamsungInternet:
		return "Samsung Internet"
	case ClientBrowserUCBrowser:
		return "UC Browser"
	case ClientBrowserVivaldi:
		return "Vivaldi"
	case ClientBrowserXbox:
		return "Xbox Browser"
	case ClientBrowserYandex:
		return "Yandex Browser"
	default:
		return ""
	}
}

func (b ClientBrowser) Value() (driver.Value, error) {
	return int64(b), nil
}

func (b *ClientBrowser) Scan(src any) error {
	return scanClientEnumUint8((*uint8)(b), src)
}

func (b ClientBrowser) MarshalJSON() ([]byte, error) {
	bytes, err := json.Marshal(b.String())
	if err != nil {
		return nil, fmt.Errorf("marshal client browser: %w", err)
	}
	return bytes, nil
}

func ClientBrowserFromLabel(value string) (ClientBrowser, bool) {
	switch normalizeClientDimensionLabel(value) {
	case "other":
		return ClientBrowserOther, true
	case "android webview":
		return ClientBrowserAndroidWebView, true
	case "chrome":
		return ClientBrowserChrome, true
	case "duckduckgo":
		return ClientBrowserDuckDuckGo, true
	case "edge":
		return ClientBrowserEdge, true
	case "facebook in-app browser":
		return ClientBrowserFacebookInApp, true
	case "firefox":
		return ClientBrowserFirefox, true
	case "instagram in-app browser":
		return ClientBrowserInstagramInApp, true
	case "internet explorer":
		return ClientBrowserInternetExplorer, true
	case "miui browser":
		return ClientBrowserMIUI, true
	case "opera":
		return ClientBrowserOpera, true
	case "playstation browser":
		return ClientBrowserPlayStation, true
	case "safari":
		return ClientBrowserSafari, true
	case "samsung internet":
		return ClientBrowserSamsungInternet, true
	case "uc browser":
		return ClientBrowserUCBrowser, true
	case "vivaldi":
		return ClientBrowserVivaldi, true
	case "xbox browser":
		return ClientBrowserXbox, true
	case "yandex browser":
		return ClientBrowserYandex, true
	default:
		return ClientBrowserUnknown, false
	}
}

func ClientBrowserFromLegacyLabel(value string) ClientBrowser {
	switch normalizeClientDimensionLabel(value) {
	case "":
		return ClientBrowserUnknown
	case "chrome mobile ios", "chrome mobile":
		return ClientBrowserChrome
	case "mobile safari":
		return ClientBrowserSafari
	case "edge mobile":
		return ClientBrowserEdge
	case "firefox mobile":
		return ClientBrowserFirefox
	case "opera mobi":
		return ClientBrowserOpera
	default:
		if browser, ok := ClientBrowserFromLabel(value); ok {
			return browser
		}
		return ClientBrowserOther
	}
}

func ParseClientBrowserFilters(values []string) []ClientBrowser {
	return parseClientDimensionFilters(values, ClientBrowserFromLabel)
}

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
