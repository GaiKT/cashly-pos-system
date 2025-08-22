package config

import (
	"os"
	"strconv"
)

// Config holds all configuration for the application
type Config struct {
	// Server configuration
	Port        string
	Environment string
	CORSOrigin  string

	// Database configuration
	DatabaseURL string

	// Authentication configuration
	JWTSecret          string
	JWTExpirationHours int
	PasswordSaltRounds int

	// OAuth configuration
	GoogleClientID     string
	GoogleClientSecret string
	FacebookAppID      string
	FacebookAppSecret  string

	// Email configuration
	EmailProvider  string
	SMTPHost       string
	SMTPPort       int
	SMTPUser       string
	SMTPPassword   string
	EmailFromName  string
	EmailFromEmail string

	// File upload configuration
	MaxFileSize  int64
	UploadPath   string
	AllowedTypes []string

	// Security configuration
	RateLimitRequests int
	RateLimitWindow   int
	SessionTimeout    int

	// Application settings
	CompanyName     string
	DefaultCurrency string
	TaxRate         float64
}

// New creates a new configuration instance with values from environment variables
func New() *Config {
	return &Config{
		// Server configuration
		Port:        getEnv("PORT", "8080"),
		Environment: getEnv("NODE_ENV", "development"),
		CORSOrigin:  getEnv("CORS_ORIGIN", "http://localhost:3000"),

		// Database configuration
		DatabaseURL: getEnv("DATABASE_URL", "postgres://admin:password123@localhost:5433/pos_db?sslmode=disable"),

		// Authentication configuration
		JWTSecret:          getEnv("JWT_SECRET", "your-super-secret-jwt-key-change-in-production"),
		JWTExpirationHours: getEnvAsInt("JWT_EXPIRATION_HOURS", 24),
		PasswordSaltRounds: getEnvAsInt("PASSWORD_SALT_ROUNDS", 12),

		// OAuth configuration
		GoogleClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
		FacebookAppID:      getEnv("FACEBOOK_APP_ID", ""),
		FacebookAppSecret:  getEnv("FACEBOOK_APP_SECRET", ""),

		// Email configuration
		EmailProvider:  getEnv("EMAIL_PROVIDER", "smtp"),
		SMTPHost:       getEnv("SMTP_HOST", "localhost"),
		SMTPPort:       getEnvAsInt("SMTP_PORT", 587),
		SMTPUser:       getEnv("SMTP_USER", ""),
		SMTPPassword:   getEnv("SMTP_PASSWORD", ""),
		EmailFromName:  getEnv("EMAIL_FROM_NAME", "POS System"),
		EmailFromEmail: getEnv("EMAIL_FROM_EMAIL", "noreply@possystem.com"),

		// File upload configuration
		MaxFileSize:  getEnvAsInt64("MAX_FILE_SIZE", 5*1024*1024), // 5MB default
		UploadPath:   getEnv("UPLOAD_PATH", "./uploads"),
		AllowedTypes: []string{"image/jpeg", "image/png", "image/gif", "image/webp"},

		// Security configuration
		RateLimitRequests: getEnvAsInt("RATE_LIMIT_REQUESTS", 100),
		RateLimitWindow:   getEnvAsInt("RATE_LIMIT_WINDOW", 900), // 15 minutes
		SessionTimeout:    getEnvAsInt("SESSION_TIMEOUT", 3600),  // 1 hour

		// Application settings
		CompanyName:     getEnv("COMPANY_NAME", "Your Store"),
		DefaultCurrency: getEnv("DEFAULT_CURRENCY", "USD"),
		TaxRate:         getEnvAsFloat64("TAX_RATE", 0.08), // 8% default
	}
}

// Helper functions to get environment variables with defaults
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsInt64(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsFloat64(key string, defaultValue float64) float64 {
	if value := os.Getenv(key); value != "" {
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			return floatValue
		}
	}
	return defaultValue
}

// IsProduction returns true if the environment is production
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}

// IsDevelopment returns true if the environment is development
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

// GetJWTExpirationDuration returns the JWT expiration as a time.Duration
func (c *Config) GetJWTExpirationDuration() int {
	return c.JWTExpirationHours
}
