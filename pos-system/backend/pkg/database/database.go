package database

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB *gorm.DB
)

// Connect establishes a connection to PostgreSQL using GORM
func Connect(databaseURL string) (*gorm.DB, error) {
	// Configure GORM logger
	gormLogger := logger.New(
		log.New(log.Writer(), "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,        // Enable color
		},
	)

	// Configure GORM
	config := &gorm.Config{
		Logger:                 gormLogger,
		SkipDefaultTransaction: false,
		PrepareStmt:            true,
	}

	// Connect to PostgreSQL
	db, err := gorm.Open(postgres.Open(databaseURL), config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	// Get underlying sql.DB to configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Configure connection pool
	sqlDB.SetMaxOpenConns(100)          // Maximum number of open connections
	sqlDB.SetMaxIdleConns(10)           // Maximum number of idle connections
	sqlDB.SetConnMaxLifetime(time.Hour) // Maximum connection lifetime

	// Test the connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping PostgreSQL: %w", err)
	}

	log.Println("Connected to PostgreSQL successfully")

	// Set global variable
	DB = db

	return db, nil
}

// Disconnect closes the database connection
func Disconnect() {
	if DB == nil {
		return
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Printf("Error getting underlying sql.DB: %v", err)
		return
	}

	err = sqlDB.Close()
	if err != nil {
		log.Printf("Error closing database connection: %v", err)
	} else {
		log.Println("Database connection closed successfully")
	}
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return DB
}

// Health checks the database connection health
func Health() error {
	if DB == nil {
		return fmt.Errorf("database connection is nil")
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	return nil
}

// Transaction executes a function within a database transaction
func Transaction(fn func(*gorm.DB) error) error {
	return DB.Transaction(fn)
}

// ExecuteRawSQL executes raw SQL queries (for running migration files)
func ExecuteRawSQL(sql string) error {
	if DB == nil {
		return fmt.Errorf("database connection is nil")
	}

	return DB.Exec(sql).Error
}

// BeginTransaction starts a new transaction
func BeginTransaction() *gorm.DB {
	return DB.Begin()
}

// IsConnectionValid checks if the database connection is valid
func IsConnectionValid() bool {
	return Health() == nil
}
