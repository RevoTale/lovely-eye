package services

import (
	"testing"

	"github.com/mileusna/useragent"
)

func TestUserAgentParsing(t *testing.T) {
	tests := []struct {
		name           string
		userAgent      string
		expectedDevice string
	}{
		{
			name:           "iPhone user agent",
			userAgent:      "Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) Safari/604.1",
			expectedDevice: "mobile",
		},
		{
			name:           "Android mobile user agent",
			userAgent:      "Mozilla/5.0 (Android 13; Mobile) Chrome/119.0",
			expectedDevice: "mobile",
		},
		{
			name:           "Windows desktop user agent",
			userAgent:      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/120.0",
			expectedDevice: "desktop",
		},
		{
			name:           "macOS desktop user agent",
			userAgent:      "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) Safari/605.1.15",
			expectedDevice: "desktop",
		},
		{
			name:           "iPad tablet user agent",
			userAgent:      "Mozilla/5.0 (iPad; CPU OS 15_0 like Mac OS X) AppleWebKit/605.1.15",
			expectedDevice: "tablet",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ua := useragent.Parse(tt.userAgent)
			device := categorizeDevice(ua)

			if device != tt.expectedDevice {
				t.Errorf("Expected device %s, got %s for UA: %s", tt.expectedDevice, device, tt.userAgent)
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
		{"Empty UA", "", false}, // Allow empty user agents
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
