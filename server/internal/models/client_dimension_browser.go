package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

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
