package analytics

import (
	"strings"

	analyticspersistence "github.com/lovely-eye/server/internal/analytics/persistence"
	"github.com/mileusna/useragent"
)

func categorizeDevice(ua useragent.UserAgent) analyticspersistence.ClientDevice {
	uaString := strings.ToLower(ua.String)

	switch {
	case strings.Contains(uaString, "watchos"),
		strings.Contains(uaString, "watch os"),
		strings.Contains(uaString, "apple watch"),
		strings.Contains(uaString, "wear os"),
		strings.Contains(uaString, "smartwatch"),
		strings.Contains(uaString, "galaxy watch"):
		return analyticspersistence.ClientDeviceWatch
	case strings.Contains(uaString, "smart-tv"),
		strings.Contains(uaString, "smarttv"),
		strings.Contains(uaString, "android tv"),
		strings.Contains(uaString, "bravia"),
		strings.Contains(uaString, "hbbtv"),
		strings.Contains(uaString, "googletv"),
		strings.Contains(uaString, "appletv"),
		strings.Contains(uaString, "crkey"),
		strings.Contains(uaString, "aft"),
		strings.Contains(uaString, "roku"),
		strings.Contains(uaString, "viera"),
		strings.Contains(uaString, "netcast"),
		strings.Contains(uaString, "tv;"):
		return analyticspersistence.ClientDeviceSmartTV
	case strings.Contains(uaString, "playstation"),
		strings.Contains(uaString, "xbox"),
		strings.Contains(uaString, "nintendo switch"):
		return analyticspersistence.ClientDeviceConsole
	}

	if ua.Tablet {
		return analyticspersistence.ClientDeviceTablet
	}
	if ua.Mobile {
		return analyticspersistence.ClientDeviceMobile
	}
	if ua.Desktop {
		return analyticspersistence.ClientDeviceDesktop
	}

	if strings.Contains(uaString, "ipad") || strings.Contains(uaString, "tablet") {
		return analyticspersistence.ClientDeviceTablet
	}
	if strings.Contains(uaString, "iphone") || strings.Contains(uaString, "ipod") {
		return analyticspersistence.ClientDeviceMobile
	}
	if strings.Contains(uaString, "android") {
		if strings.Contains(uaString, "mobile") {
			return analyticspersistence.ClientDeviceMobile
		}
		return analyticspersistence.ClientDeviceTablet
	}
	if strings.Contains(uaString, "windows") ||
		strings.Contains(uaString, "macintosh") ||
		strings.Contains(uaString, "linux") ||
		strings.Contains(uaString, "cros") {
		return analyticspersistence.ClientDeviceDesktop
	}

	return analyticspersistence.ClientDeviceDesktop
}

func normalizeBrowser(ua useragent.UserAgent) analyticspersistence.ClientBrowser {
	uaString := strings.ToLower(ua.String)

	switch {
	case strings.Contains(uaString, "playstation"):
		return analyticspersistence.ClientBrowserPlayStation
	case strings.Contains(uaString, "xbox"):
		return analyticspersistence.ClientBrowserXbox
	case strings.Contains(uaString, "fb_iab"),
		strings.Contains(uaString, "fban"),
		strings.Contains(uaString, "fbav"):
		return analyticspersistence.ClientBrowserFacebookInApp
	case strings.Contains(uaString, "instagram"):
		return analyticspersistence.ClientBrowserInstagramInApp
	case strings.Contains(uaString, "edg/"),
		strings.Contains(uaString, "edgios"),
		ua.IsEdge():
		return analyticspersistence.ClientBrowserEdge
	case strings.Contains(uaString, "samsungbrowser"):
		return analyticspersistence.ClientBrowserSamsungInternet
	case strings.Contains(uaString, "opr/"),
		strings.Contains(uaString, "opera mini"),
		strings.Contains(uaString, "opera mobi"),
		ua.IsOpera(),
		ua.IsOperaMini():
		return analyticspersistence.ClientBrowserOpera
	case strings.Contains(uaString, "vivaldi"):
		return analyticspersistence.ClientBrowserVivaldi
	case strings.Contains(uaString, "yabrowser"),
		strings.Contains(uaString, "yowser"):
		return analyticspersistence.ClientBrowserYandex
	case strings.Contains(uaString, "duckduckgo"):
		return analyticspersistence.ClientBrowserDuckDuckGo
	case strings.Contains(uaString, "ucbrowser"),
		strings.Contains(uaString, "ucweb"):
		return analyticspersistence.ClientBrowserUCBrowser
	case strings.Contains(uaString, "miuibrowser"):
		return analyticspersistence.ClientBrowserMIUI
	case strings.Contains(uaString, "msie"),
		strings.Contains(uaString, "trident"),
		ua.IsInternetExplorer():
		return analyticspersistence.ClientBrowserInternetExplorer
	case strings.Contains(uaString, "wv"),
		strings.Contains(uaString, "webview"):
		return analyticspersistence.ClientBrowserAndroidWebView
	case strings.Contains(uaString, "crios"),
		strings.Contains(uaString, "chrome"),
		ua.IsChrome():
		return analyticspersistence.ClientBrowserChrome
	case strings.Contains(uaString, "fxios"),
		strings.Contains(uaString, "firefox"),
		ua.IsFirefox():
		return analyticspersistence.ClientBrowserFirefox
	case strings.Contains(uaString, "safari"),
		(strings.Contains(uaString, "applewebkit") &&
			(strings.Contains(uaString, "iphone") || strings.Contains(uaString, "ipad") || strings.Contains(uaString, "macintosh"))),
		ua.IsSafari():
		return analyticspersistence.ClientBrowserSafari
	}

	name := strings.TrimSpace(ua.Name)
	if name == "" {
		return analyticspersistence.ClientBrowserOther
	}

	if browser, ok := analyticspersistence.ClientBrowserFromLabel(name); ok {
		return browser
	}
	return analyticspersistence.ClientBrowserFromLegacyLabel(name)
}

func normalizeOS(ua useragent.UserAgent) analyticspersistence.ClientOS {
	uaString := strings.ToLower(ua.String)

	switch {
	case strings.Contains(uaString, "wear os"):
		return analyticspersistence.ClientOSWearOS
	case strings.Contains(uaString, "watchos"),
		strings.Contains(uaString, "watch os"),
		strings.Contains(uaString, "apple watch"):
		return analyticspersistence.ClientOSWatchOS
	case strings.Contains(uaString, "playstation"):
		return analyticspersistence.ClientOSPlayStation
	case strings.Contains(uaString, "xbox"):
		return analyticspersistence.ClientOSXbox
	case strings.Contains(uaString, "ipad"):
		return analyticspersistence.ClientOSIPadOS
	case strings.Contains(uaString, "iphone"),
		strings.Contains(uaString, "ipod"),
		ua.IsIOS():
		return analyticspersistence.ClientOSIOS
	case strings.Contains(uaString, "android"):
		return analyticspersistence.ClientOSAndroid
	case strings.Contains(uaString, "cros"),
		ua.IsChromeOS():
		return analyticspersistence.ClientOSChromeOS
	case strings.Contains(uaString, "windows"),
		ua.IsWindows():
		return analyticspersistence.ClientOSWindows
	case strings.Contains(uaString, "mac os"),
		strings.Contains(uaString, "macintosh"),
		ua.IsMacOS():
		return analyticspersistence.ClientOSMacOS
	case strings.Contains(uaString, "linux"),
		ua.IsLinux():
		return analyticspersistence.ClientOSLinux
	}

	if os, ok := analyticspersistence.ClientOSFromLabel(strings.TrimSpace(ua.OS)); ok {
		return os
	}
	return analyticspersistence.ClientOSFromLegacyLabel(strings.TrimSpace(ua.OS))
}

func categorizeScreenSize(width int) analyticspersistence.ClientScreenSize {
	return analyticspersistence.ClientScreenSizeFromWidth(width)
}
