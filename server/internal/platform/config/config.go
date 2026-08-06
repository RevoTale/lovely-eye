package config

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	defaultIPDBURL       string = "https://download.db-ip.com/free/dbip-country-lite.mmdb.gz"
	defaultIPDBLocalPath string = "/data/GeoLite2-Country.mmdb"
)

const defaultDBDSN string = "file:data/lovely_eye.db?cache=shared&mode=rwc"
const defaultTrustedProxyCIDRs string = "127.0.0.1/32,::1/128,10.0.0.0/8,172.16.0.0/12,192.168.0.0/16,fc00::/7"

const (
	DBDriverPG     DBDriver = "postgres"
	DBDriverSQLite DBDriver = "sqlite"
)

type DBDriver = string

type Config struct {
	Server    ServerConfig
	Database  DatabaseConfig
	Auth      AuthConfig
	Analytics AnalyticsConfig
	GraphQL   GraphQLConfig
	Dashboard DashboardConfig
	GeoIP     GeoIPConfig
	LogLevel  slog.Level // Log level: DEBUG(-4), INFO(0), WARN(4), ERROR(8) - default: WARN
	TrackerJS []byte     // Optional: for testing, to avoid loading from file
}

type ServerConfig struct {
	Host          string
	Port          string
	BasePath      string
	DashboardPath string
}

type DatabaseConfig struct {
	Driver         DBDriver
	DSN            string
	MaxConns       int
	MinConns       int
	ConnectTimeout time.Duration
}

type AuthConfig struct {
	JWTSecret            string
	AccessTokenExpiry    time.Duration
	RefreshExpiry        time.Duration
	AllowRegistration    bool
	SecureCookies        bool
	CookieDomain         string
	InitialAdminUsername string
	InitialAdminPassword string // password for initial admin (optional)
	RateLimitEnabled     bool
	RateLimitAttempts    int
	RateLimitWindow      time.Duration
}

type AnalyticsConfig struct {
	IdentitySecret        string
	MaxBodyBytes          int64
	MaxPropertiesBytes    int
	MaxSinglePageDuration time.Duration
	RateLimitEnabled      bool
	RateLimitPerMinute    int
	RateLimitBurst        int
	TrustedProxyCIDRs     []string
}

type GraphQLConfig struct {
	MaxBodyBytes  int64
	MaxComplexity int
}

type DashboardConfig struct {
	MaxDailyRangeDays     int
	MaxHourlyRangeDays    int
	MaxFilterValues       int
	MaxFilterStringLength int
}

type GeoIPConfig struct {
	DBPath            string
	DownloadURL       string
	MaxMindLicenseKey string
}

func Load() (Config, error) {
	reader := newEnvReader()
	authSecret, err := getJWTSecret()
	if err != nil {
		return Config{}, err
	}
	identitySecret, err := getAnalyticsIdentitySecret(authSecret)
	if err != nil {
		return Config{}, err
	}
	basePath := getEnv("BASE_PATH", "/")
	downloadURL := getEnv("GEOIP_DOWNLOAD_URL", "")
	maxMindKey := getEnv("GEOIP_MAXMIND_LICENSE_KEY", "")
	initialAdminUsername := getEnv("INITIAL_ADMIN_USERNAME", "")
	initialAdminPassword := getEnv("INITIAL_ADMIN_PASSWORD", "")
	allowRegistration := defaultAllowRegistration(initialAdminUsername, initialAdminPassword)
	if explicitAllowRegistration, ok := reader.OptionalBool("ALLOW_REGISTRATION"); ok {
		allowRegistration = explicitAllowRegistration
	}
	if downloadURL == "" && maxMindKey == "" {
		downloadURL = defaultIPDBURL
	}

	if basePath != "/" {
		basePath = "/" + strings.Trim(basePath, "/")
	}
	cfg := Config{
		Server: ServerConfig{
			Host:          getEnv("SERVER_HOST", "0.0.0.0"),
			Port:          getEnv("SERVER_PORT", "8080"),
			BasePath:      basePath,
			DashboardPath: getEnv("DASHBOARD_PATH", "dashboard"),
		},
		Database: DatabaseConfig{
			Driver:         getEnv("DB_DRIVER", DBDriverSQLite),
			DSN:            getEnv("DB_DSN", defaultDBDSN),
			MaxConns:       reader.Int("DB_MAX_CONNS", 10),
			MinConns:       reader.Int("DB_MIN_CONNS", 1),
			ConnectTimeout: reader.Duration("DB_CONNECT_TIMEOUT", 7*time.Second),
		},
		Auth: AuthConfig{
			JWTSecret:            authSecret,
			AccessTokenExpiry:    time.Duration(reader.Int("JWT_ACCESS_EXPIRY_MINUTES", 15)) * time.Minute,
			RefreshExpiry:        time.Duration(reader.Int("JWT_REFRESH_DAYS", 7)) * 24 * time.Hour,
			AllowRegistration:    allowRegistration,
			SecureCookies:        reader.Bool("SECURE_COOKIES", true),
			CookieDomain:         getEnv("COOKIE_DOMAIN", ""),
			InitialAdminUsername: initialAdminUsername,
			InitialAdminPassword: initialAdminPassword,
			RateLimitEnabled:     reader.Bool("AUTH_RATE_LIMIT_ENABLED", true),
			RateLimitAttempts:    reader.Int("AUTH_RATE_LIMIT_ATTEMPTS", 10),
			RateLimitWindow:      reader.Duration("AUTH_RATE_LIMIT_WINDOW", 15*time.Minute),
		},
		Analytics: AnalyticsConfig{
			IdentitySecret:        identitySecret,
			MaxBodyBytes:          int64(reader.Int("ANALYTICS_MAX_BODY_BYTES", 16*1024)),
			MaxPropertiesBytes:    reader.Int("ANALYTICS_MAX_PROPERTIES_BYTES", 8*1024),
			MaxSinglePageDuration: reader.Duration("ANALYTICS_MAX_SINGLE_PAGE_DURATION", 4*time.Hour),
			RateLimitEnabled:      reader.Bool("ANALYTICS_RATE_LIMIT_ENABLED", true),
			RateLimitPerMinute:    reader.Int("ANALYTICS_RATE_LIMIT_PER_MINUTE", 120),
			RateLimitBurst:        reader.Int("ANALYTICS_RATE_LIMIT_BURST", 240),
			TrustedProxyCIDRs:     getEnvCSV("TRUSTED_PROXY_CIDRS", defaultTrustedProxyCIDRs),
		},
		GraphQL: GraphQLConfig{
			MaxBodyBytes:  int64(reader.Int("GRAPHQL_MAX_BODY_BYTES", 1024*1024)),
			MaxComplexity: reader.Int("GRAPHQL_MAX_COMPLEXITY", 300),
		},
		Dashboard: DashboardConfig{
			MaxDailyRangeDays:     reader.Int("DASHBOARD_MAX_DAILY_RANGE_DAYS", 730),
			MaxHourlyRangeDays:    reader.Int("DASHBOARD_MAX_HOURLY_RANGE_DAYS", 31),
			MaxFilterValues:       reader.Int("DASHBOARD_MAX_FILTER_VALUES", 100),
			MaxFilterStringLength: reader.Int("DASHBOARD_MAX_FILTER_STRING_LENGTH", 2048),
		},
		GeoIP: GeoIPConfig{
			DBPath:            getEnv("GEOIP_DB_PATH", defaultIPDBLocalPath),
			DownloadURL:       downloadURL,
			MaxMindLicenseKey: maxMindKey,
		},
		LogLevel: reader.LogLevel("LOG_LEVEL", slog.LevelWarn),
	}
	if reader.err != nil {
		return Config{}, fmt.Errorf("parse environment: %w", reader.err)
	}
	if err := cfg.validate(); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func getJWTSecret() (string, error) {
	if secret := os.Getenv("JWT_SECRET"); secret != "" {
		if len(secret) < 32 {
			return "", errors.New("JWT_SECRET must be at least 32 characters")
		}
		return secret, nil
	}
	slog.Warn("JWT_SECRET is not set; generated development tokens will not survive restart")
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("generate JWT secret: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}

func getAnalyticsIdentitySecret(fallback string) (string, error) {
	if secret := os.Getenv("ANALYTICS_IDENTITY_SECRET"); secret != "" {
		if len(secret) < 32 {
			return "", errors.New("ANALYTICS_IDENTITY_SECRET must be at least 32 characters")
		}
		return secret, nil
	}
	return fallback, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvCSV(key, defaultValue string) []string {
	raw := getEnv(key, defaultValue)
	parts := strings.Split(raw, ",")
	values := make([]string, 0, len(parts))
	for _, part := range parts {
		value := strings.TrimSpace(part)
		if value != "" {
			values = append(values, value)
		}
	}
	return values
}

func defaultAllowRegistration(initialAdminUsername, initialAdminPassword string) bool {
	return initialAdminUsername == "" || initialAdminPassword == ""
}

type envReader struct {
	err error
}

func newEnvReader() *envReader {
	return &envReader{}
}

func (r *envReader) Int(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		r.addError(fmt.Errorf("%s=%q must be an integer: %w", key, value, err))
		return defaultValue
	}
	return parsed
}

func (r *envReader) Bool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	parsed, err := parseBool(value)
	if err != nil {
		r.addError(fmt.Errorf("%s: %w", key, err))
		return defaultValue
	}
	return parsed
}

func (r *envReader) OptionalBool(key string) (bool, bool) {
	value, ok := os.LookupEnv(key)
	if !ok || value == "" {
		return false, false
	}
	parsed, err := parseBool(value)
	if err != nil {
		r.addError(fmt.Errorf("%s: %w", key, err))
		return false, true
	}
	return parsed, true
}

func parseBool(value string) (bool, error) {
	switch strings.ToLower(value) {
	case "true", "1", "yes":
		return true, nil
	case "false", "0", "no":
		return false, nil
	default:
		return false, fmt.Errorf("value %q must be one of true, false, 1, 0, yes, or no", value)
	}
}

func (r *envReader) LogLevel(key string, defaultValue slog.Level) slog.Level {
	value := strings.ToLower(os.Getenv(key))
	if value == "" {
		return defaultValue
	}
	switch value {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		r.addError(fmt.Errorf("%s=%q must be debug, info, warn, warning, or error", key, value))
		return defaultValue
	}
}

func (r *envReader) Duration(key string, defaultValue time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	duration, err := time.ParseDuration(value)
	if err != nil {
		r.addError(fmt.Errorf("%s=%q must be a duration: %w", key, value, err))
		return defaultValue
	}
	return duration
}

func (r *envReader) addError(err error) {
	r.err = errors.Join(r.err, err)
}

func (cfg Config) validate() error {
	var err error
	requirePositive := func(name string, value int64) {
		if value <= 0 {
			err = errors.Join(err, fmt.Errorf("%s must be greater than zero", name))
		}
	}
	requirePositive("DB_MAX_CONNS", int64(cfg.Database.MaxConns))
	if cfg.Database.MinConns < 0 {
		err = errors.Join(err, errors.New("DB_MIN_CONNS must not be negative"))
	}
	if cfg.Database.MinConns > cfg.Database.MaxConns {
		err = errors.Join(err, errors.New("DB_MIN_CONNS must not exceed DB_MAX_CONNS"))
	}
	requirePositive("DB_CONNECT_TIMEOUT", int64(cfg.Database.ConnectTimeout))
	requirePositive("JWT_ACCESS_EXPIRY_MINUTES", int64(cfg.Auth.AccessTokenExpiry))
	requirePositive("JWT_REFRESH_DAYS", int64(cfg.Auth.RefreshExpiry))
	requirePositive("AUTH_RATE_LIMIT_ATTEMPTS", int64(cfg.Auth.RateLimitAttempts))
	requirePositive("AUTH_RATE_LIMIT_WINDOW", int64(cfg.Auth.RateLimitWindow))
	requirePositive("ANALYTICS_MAX_BODY_BYTES", cfg.Analytics.MaxBodyBytes)
	requirePositive("ANALYTICS_MAX_PROPERTIES_BYTES", int64(cfg.Analytics.MaxPropertiesBytes))
	requirePositive("ANALYTICS_MAX_SINGLE_PAGE_DURATION", int64(cfg.Analytics.MaxSinglePageDuration))
	requirePositive("ANALYTICS_RATE_LIMIT_PER_MINUTE", int64(cfg.Analytics.RateLimitPerMinute))
	requirePositive("ANALYTICS_RATE_LIMIT_BURST", int64(cfg.Analytics.RateLimitBurst))
	requirePositive("GRAPHQL_MAX_BODY_BYTES", cfg.GraphQL.MaxBodyBytes)
	requirePositive("GRAPHQL_MAX_COMPLEXITY", int64(cfg.GraphQL.MaxComplexity))
	requirePositive("DASHBOARD_MAX_DAILY_RANGE_DAYS", int64(cfg.Dashboard.MaxDailyRangeDays))
	requirePositive("DASHBOARD_MAX_HOURLY_RANGE_DAYS", int64(cfg.Dashboard.MaxHourlyRangeDays))
	requirePositive("DASHBOARD_MAX_FILTER_VALUES", int64(cfg.Dashboard.MaxFilterValues))
	requirePositive("DASHBOARD_MAX_FILTER_STRING_LENGTH", int64(cfg.Dashboard.MaxFilterStringLength))
	if err != nil {
		return fmt.Errorf("validate configuration: %w", err)
	}
	return nil
}
