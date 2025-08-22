package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/pos-system/backend/pkg/config"
	"github.com/pos-system/backend/pkg/database"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize configuration
	cfg := config.New()

	// Initialize database connection
	_, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer database.Disconnect()

	// Setup Gin router
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Middleware setup
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// CORS configuration
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{cfg.CORSOrigin}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	corsConfig.AllowCredentials = true
	router.Use(cors.New(corsConfig))

	// Health check endpoint
	router.GET("/api/health", func(c *gin.Context) {
		dbStatus := "connected"
		if err := database.Health(); err != nil {
			dbStatus = "disconnected"
		}

		c.JSON(http.StatusOK, gin.H{
			"status":    "ok",
			"timestamp": time.Now().UTC(),
			"service":   "pos-backend",
			"version":   "1.0.0",
			"database":  dbStatus,
			"db_type":   "postgresql",
		})
	})

	// Basic API structure (handlers will be added in next tasks)
	api := router.Group("/api")
	{
		// Authentication routes (placeholder)
		auth := api.Group("/auth")
		{
			auth.GET("/status", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"message": "Authentication endpoints will be implemented in Task 4",
					"status":  "pending",
				})
			})
		}

		// Protected routes placeholder
		protected := api.Group("/")
		{
			protected.GET("/protected", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"message": "Protected endpoints will be implemented after authentication",
					"status":  "pending",
				})
			})
		}
	}

	// Setup HTTP server
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("🚀 Server starting on port %s", cfg.Port)
		log.Printf("🌍 Environment: %s", cfg.Environment)
		log.Printf("🔗 CORS Origin: %s", cfg.CORSOrigin)
		log.Printf("📊 Database: %s", cfg.DatabaseURL)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down server...")

	// The context is used to inform the server it has 30 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("✅ Server exited gracefully")
}
