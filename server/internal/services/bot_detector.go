package services

import (
	"strings"
)

// BotDetector provides bot detection capabilities
type BotDetector struct {
	botPatterns []string
}

// NewBotDetector creates a new bot detector
func NewBotDetector() *BotDetector {
	return &BotDetector{
		botPatterns: []string{
			// Search engine bots
			"Googlebot", "Bingbot", "Slurp", "DuckDuckBot", "Baiduspider", "YandexBot",
			"Sogou", "Exabot", "facebot", "ia_archiver",
			// SEO/monitoring bots
			"AhrefsBot", "SemrushBot", "DotBot", "Applebot", "MJ12bot", "rogerbot",
			"LinkpadBot", "PingdomBot", "DataForSeoBot", "SeznamBot",
			// Social media bots
			"Twitterbot", "facebookexternalhit", "LinkedInBot", "Discordbot", "TelegramBot",
			"WhatsApp", "SkypeUriPreview", "Slackbot",
			// Monitoring/uptime bots
			"UptimeRobot", "StatusCake", "Pingdom", "GTmetrix", "Site24x7",
			// Scrapers
			"scrapy", "curl", "wget", "python-requests",
			"Apache-HttpClient", "axios", "node-fetch",
			// Headless browsers (often used for scraping)
			"HeadlessChrome", "PhantomJS", "Selenium", "Playwright", "Puppeteer",
			// Other common bots
			"bot", "crawler", "spider", "scraper", "monitor",
			// Prerendering services
			"Prerender", "rendertron",
		},
	}
}

// IsBot checks if a user agent string is likely a bot
func (bd *BotDetector) IsBot(userAgent string) bool {
	if userAgent == "" {
		// Empty user agent - allow it (could be legitimate traffic, privacy tools, or tests)
		// Real analytics tools like GoatCounter and Umami don't block empty user agents
		return false
	}

	userAgentLower := strings.ToLower(userAgent)

	// Check against known bot patterns
	for _, pattern := range bd.botPatterns {
		if strings.Contains(userAgentLower, strings.ToLower(pattern)) {
			return true
		}
	}

	// Additional heuristics for suspicious patterns
	if strings.HasPrefix(userAgentLower, "mozilla/5.0 (compatible;") &&
		!strings.Contains(userAgentLower, "chrome") &&
		!strings.Contains(userAgentLower, "firefox") &&
		!strings.Contains(userAgentLower, "safari") {
		// Likely a bot pretending to be compatible
		return true
	}

	return false
}

// IsPrefetchRequest checks if a request is a prefetch/prerender
func (bd *BotDetector) IsPrefetchRequest(purpose string) bool {
	purposeLower := strings.ToLower(purpose)
	return strings.Contains(purposeLower, "prefetch") ||
		strings.Contains(purposeLower, "prerender")
}
