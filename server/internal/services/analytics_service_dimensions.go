package services

import (
	"strings"

	"github.com/lovely-eye/server/internal/models"
	"github.com/mileusna/useragent"
)

func categorizeDevice(ua useragent.UserAgent) models.ClientDevice {
	uaString := strings.ToLower(ua.String)

	switch {
	case strings.Contains(uaString, "watchos"),
		strings.Contains(uaString, "watch os"),
		strings.Contains(uaString, "apple watch"),
		strings.Contains(uaString, "wear os"),
		strings.Contains(uaString, "smartwatch"),
		strings.Contains(uaString, "galaxy watch"):
		return models.ClientDeviceWatch
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
		return models.ClientDeviceSmartTV
	case strings.Contains(uaString, "playstation"),
		strings.Contains(uaString, "xbox"),
		strings.Contains(uaString, "nintendo switch"):
		return models.ClientDeviceConsole
	}

	if ua.Tablet {
		return models.ClientDeviceTablet
	}
	if ua.Mobile {
		return models.ClientDeviceMobile
	}
	if ua.Desktop {
		return models.ClientDeviceDesktop
	}

	if strings.Contains(uaString, "ipad") || strings.Contains(uaString, "tablet") {
		return models.ClientDeviceTablet
	}
	if strings.Contains(uaString, "iphone") || strings.Contains(uaString, "ipod") {
		return models.ClientDeviceMobile
	}
	if strings.Contains(uaString, "android") {
		if strings.Contains(uaString, "mobile") {
			return models.ClientDeviceMobile
		}
		return models.ClientDeviceTablet
	}
	if strings.Contains(uaString, "windows") ||
		strings.Contains(uaString, "macintosh") ||
		strings.Contains(uaString, "linux") ||
		strings.Contains(uaString, "cros") {
		return models.ClientDeviceDesktop
	}

	return models.ClientDeviceDesktop
}

func normalizeBrowser(ua useragent.UserAgent) models.ClientBrowser {
	uaString := strings.ToLower(ua.String)

	switch {
	case strings.Contains(uaString, "playstation"):
		return models.ClientBrowserPlayStation
	case strings.Contains(uaString, "xbox"):
		return models.ClientBrowserXbox
	case strings.Contains(uaString, "fb_iab"),
		strings.Contains(uaString, "fban"),
		strings.Contains(uaString, "fbav"):
		return models.ClientBrowserFacebookInApp
	case strings.Contains(uaString, "instagram"):
		return models.ClientBrowserInstagramInApp
	case strings.Contains(uaString, "edg/"),
		strings.Contains(uaString, "edgios"),
		ua.IsEdge():
		return models.ClientBrowserEdge
	case strings.Contains(uaString, "samsungbrowser"):
		return models.ClientBrowserSamsungInternet
	case strings.Contains(uaString, "opr/"),
		strings.Contains(uaString, "opera mini"),
		strings.Contains(uaString, "opera mobi"),
		ua.IsOpera(),
		ua.IsOperaMini():
		return models.ClientBrowserOpera
	case strings.Contains(uaString, "vivaldi"):
		return models.ClientBrowserVivaldi
	case strings.Contains(uaString, "yabrowser"),
		strings.Contains(uaString, "yowser"):
		return models.ClientBrowserYandex
	case strings.Contains(uaString, "duckduckgo"):
		return models.ClientBrowserDuckDuckGo
	case strings.Contains(uaString, "ucbrowser"),
		strings.Contains(uaString, "ucweb"):
		return models.ClientBrowserUCBrowser
	case strings.Contains(uaString, "miuibrowser"):
		return models.ClientBrowserMIUI
	case strings.Contains(uaString, "msie"),
		strings.Contains(uaString, "trident"),
		ua.IsInternetExplorer():
		return models.ClientBrowserInternetExplorer
	case strings.Contains(uaString, "wv"),
		strings.Contains(uaString, "webview"):
		return models.ClientBrowserAndroidWebView
	case strings.Contains(uaString, "crios"),
		strings.Contains(uaString, "chrome"),
		ua.IsChrome():
		return models.ClientBrowserChrome
	case strings.Contains(uaString, "fxios"),
		strings.Contains(uaString, "firefox"),
		ua.IsFirefox():
		return models.ClientBrowserFirefox
	case strings.Contains(uaString, "safari"),
		(strings.Contains(uaString, "applewebkit") &&
			(strings.Contains(uaString, "iphone") || strings.Contains(uaString, "ipad") || strings.Contains(uaString, "macintosh"))),
		ua.IsSafari():
		return models.ClientBrowserSafari
	}

	name := strings.TrimSpace(ua.Name)
	if name == "" {
		return models.ClientBrowserOther
	}

	if browser, ok := models.ClientBrowserFromLabel(name); ok {
		return browser
	}
	return models.ClientBrowserFromLegacyLabel(name)
}

func normalizeOS(ua useragent.UserAgent) models.ClientOS {
	uaString := strings.ToLower(ua.String)

	switch {
	case strings.Contains(uaString, "wear os"):
		return models.ClientOSWearOS
	case strings.Contains(uaString, "watchos"),
		strings.Contains(uaString, "watch os"),
		strings.Contains(uaString, "apple watch"):
		return models.ClientOSWatchOS
	case strings.Contains(uaString, "playstation"):
		return models.ClientOSPlayStation
	case strings.Contains(uaString, "xbox"):
		return models.ClientOSXbox
	case strings.Contains(uaString, "ipad"):
		return models.ClientOSIPadOS
	case strings.Contains(uaString, "iphone"),
		strings.Contains(uaString, "ipod"),
		ua.IsIOS():
		return models.ClientOSIOS
	case strings.Contains(uaString, "android"):
		return models.ClientOSAndroid
	case strings.Contains(uaString, "cros"),
		ua.IsChromeOS():
		return models.ClientOSChromeOS
	case strings.Contains(uaString, "windows"),
		ua.IsWindows():
		return models.ClientOSWindows
	case strings.Contains(uaString, "mac os"),
		strings.Contains(uaString, "macintosh"),
		ua.IsMacOS():
		return models.ClientOSMacOS
	case strings.Contains(uaString, "linux"),
		ua.IsLinux():
		return models.ClientOSLinux
	}

	if os, ok := models.ClientOSFromLabel(strings.TrimSpace(ua.OS)); ok {
		return os
	}
	return models.ClientOSFromLegacyLabel(strings.TrimSpace(ua.OS))
}

func categorizeScreenSize(width int) models.ClientScreenSize {
	return models.ClientScreenSizeFromWidth(width)
}
