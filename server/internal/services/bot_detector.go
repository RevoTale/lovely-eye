package services

import (
	"strings"
)

type BotDetector struct {
	botPatterns []string
}

func NewBotDetector() *BotDetector {
	return &BotDetector{
		botPatterns: []string{

			"Googlebot", "Bingbot", "Slurp", "DuckDuckBot", "Baiduspider", "YandexBot",
			"Sogou", "Exabot", "facebot", "ia_archiver",

			"AhrefsBot", "SemrushBot", "DotBot", "Applebot", "MJ12bot", "rogerbot",
			"LinkpadBot", "PingdomBot", "DataForSeoBot", "SeznamBot",

			"Twitterbot", "facebookexternalhit", "LinkedInBot", "Discordbot", "TelegramBot",
			"WhatsApp", "SkypeUriPreview", "Slackbot",

			"UptimeRobot", "StatusCake", "Pingdom", "GTmetrix", "Site24x7",

			"scrapy", "curl", "wget", "python-requests",
			"Apache-HttpClient", "axios", "node-fetch",

			"HeadlessChrome", "PhantomJS", "Selenium", "Playwright", "Puppeteer",

			"bot", "crawler", "spider", "scraper", "monitor",

			"Prerender", "rendertron",
		},
	}
}

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

func (bd *BotDetector) IsPrefetchRequest(purpose string) bool {
	purposeLower := strings.ToLower(purpose)
	return strings.Contains(purposeLower, "prefetch") ||
		strings.Contains(purposeLower, "prerender")
}
