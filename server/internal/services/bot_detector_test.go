package services

import (
	"testing"

	"github.com/mileusna/useragent"
)

func TestUserAgentParsing(t *testing.T) {
	tests := []struct {
		name            string
		userAgent       string
		expectedBrowser string
		expectedOS      string
		expectedDevice  string
	}{
		{
			name:            "iPhone user agent",
			userAgent:       "Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) Safari/604.1",
			expectedBrowser: "Safari",
			expectedOS:      "iOS",
			expectedDevice:  "mobile",
		},
		{
			name:            "Android mobile user agent",
			userAgent:       "Mozilla/5.0 (Linux; Android 13; Pixel 7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0 Mobile Safari/537.36",
			expectedBrowser: "Chrome",
			expectedOS:      "Android",
			expectedDevice:  "mobile",
		},
		{
			name:            "Windows desktop user agent",
			userAgent:       "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/120.0",
			expectedBrowser: "Chrome",
			expectedOS:      "Windows",
			expectedDevice:  "desktop",
		},
		{
			name:            "macOS desktop user agent",
			userAgent:       "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) Safari/605.1.15",
			expectedBrowser: "Safari",
			expectedOS:      "macOS",
			expectedDevice:  "desktop",
		},
		{
			name:            "iPad tablet user agent",
			userAgent:       "Mozilla/5.0 (iPad; CPU OS 15_0 like Mac OS X) AppleWebKit/605.1.15",
			expectedBrowser: "Safari",
			expectedOS:      "iPadOS",
			expectedDevice:  "tablet",
		},
		{
			name:            "Android TV user agent",
			userAgent:       "Mozilla/5.0 (Linux; Android 12; BRAVIA 4K GB ATV3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
			expectedBrowser: "Chrome",
			expectedOS:      "Android",
			expectedDevice:  "smart-tv",
		},
		{
			name:            "Apple Watch user agent",
			userAgent:       "Mozilla/5.0 (Apple Watch; CPU WatchOS 10_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/10.0 Mobile/15E148 Safari/604.1",
			expectedBrowser: "Safari",
			expectedOS:      "watchOS",
			expectedDevice:  "watch",
		},
		{
			name:            "PlayStation user agent",
			userAgent:       "Mozilla/5.0 (PlayStation 5 3.20) AppleWebKit/605.1.15 (KHTML, like Gecko)",
			expectedBrowser: "PlayStation Browser",
			expectedOS:      "PlayStation OS",
			expectedDevice:  "console",
		},
		{
			name:            "Samsung Internet user agent",
			userAgent:       "Mozilla/5.0 (Linux; Android 14; SAMSUNG SM-S921B) AppleWebKit/537.36 (KHTML, like Gecko) SamsungBrowser/24.0 Chrome/117.0.0.0 Mobile Safari/537.36",
			expectedBrowser: "Samsung Internet",
			expectedOS:      "Android",
			expectedDevice:  "mobile",
		},
		{
			name:            "Edge on ChromeOS user agent",
			userAgent:       "Mozilla/5.0 (X11; CrOS x86_64 15474.84.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 Edg/120.0.0.0",
			expectedBrowser: "Edge",
			expectedOS:      "ChromeOS",
			expectedDevice:  "desktop",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ua := useragent.Parse(tt.userAgent)
			browser := normalizeBrowser(ua)
			os := normalizeOS(ua)
			device := categorizeDevice(ua)

			if browser != tt.expectedBrowser {
				t.Errorf("Expected browser %s, got %s for UA: %s", tt.expectedBrowser, browser, tt.userAgent)
			}

			if os != tt.expectedOS {
				t.Errorf("Expected OS %s, got %s for UA: %s", tt.expectedOS, os, tt.userAgent)
			}

			if device != tt.expectedDevice {
				t.Errorf("Expected device %s, got %s for UA: %s", tt.expectedDevice, device, tt.userAgent)
			}
		})
	}
}

func TestCategorizeScreenSize(t *testing.T) {
	tests := []struct {
		name     string
		width    int
		expected string
	}{
		{
			name:     "watch width",
			width:    280,
			expected: "watch",
		},
		{
			name:     "phone width",
			width:    390,
			expected: "xs",
		},
		{
			name:     "tablet width",
			width:    820,
			expected: "md",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := categorizeScreenSize(tt.width); got != tt.expected {
				t.Errorf("Expected screen size %s, got %s for width %d", tt.expected, got, tt.width)
			}
		})
	}
}

func TestBotDetection(t *testing.T) {
	bd := NewBotDetector()

	tests := []struct {
		name      string
		userAgent string
		isBot     bool
	}{
		{"Googlebot", "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)", true},
		{"Normal Chrome", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36", false},
		{"Empty UA", "", false},
		{"curl", "curl/7.68.0", true},
		{"wget", "Wget/1.20.3 (linux-gnu)", true},
		{"Normal Safari", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0 Safari/605.1.15", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := bd.IsBot(tt.userAgent)
			if result != tt.isBot {
				t.Errorf("Expected IsBot(%s) = %v, got %v", tt.userAgent, tt.isBot, result)
			}
		})
	}
}
