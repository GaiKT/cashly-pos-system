package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/pos-system/backend/internal/models"
	"github.com/pos-system/backend/internal/services"
)

// AuthMiddleware handles JWT authentication
type AuthMiddleware struct {
	authService *services.AuthService
}

// NewAuthMiddleware creates a new authentication middleware
func NewAuthMiddleware(authService *services.AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

// RequireAuth middleware validates JWT token and sets user context
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header required",
			})
			c.Abort()
			return
		}

		// Check if header starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid authorization header format",
			})
			c.Abort()
			return
		}

		// Extract token
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Token required",
			})
			c.Abort()
			return
		}

		// Validate token and get user
		user, err := m.authService.GetUserFromToken(c.Request.Context(), token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// Set user in context
		c.Set("user", user)
		c.Set("user_id", user.ID)
		c.Set("user_email", user.Email)
		c.Set("user_role", user.Role)

		c.Next()
	}
}

// OptionalAuth middleware validates JWT token if present but doesn't require it
func (m *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		// Check if header starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.Next()
			return
		}

		// Extract token
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			c.Next()
			return
		}

		// Validate token and get user
		user, err := m.authService.GetUserFromToken(c.Request.Context(), token)
		if err != nil {
			// Invalid token, but don't abort since auth is optional
			c.Next()
			return
		}

		// Set user in context
		c.Set("user", user)
		c.Set("user_id", user.ID)
		c.Set("user_email", user.Email)
		c.Set("user_role", user.Role)

		c.Next()
	}
}

// RequireRole middleware checks if authenticated user has required role
func (m *AuthMiddleware) RequireRole(requiredRole models.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user from context (should be set by RequireAuth middleware)
		userInterface, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authentication required",
			})
			c.Abort()
			return
		}

		user, ok := userInterface.(*models.User)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Invalid user context",
			})
			c.Abort()
			return
		}

		// Check role
		if err := m.authService.ValidateRole(user.Role, requiredRole); err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Insufficient permissions",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAdmin middleware shorthand for requiring admin role
func (m *AuthMiddleware) RequireAdmin() gin.HandlerFunc {
	return m.RequireRole(models.RoleAdmin)
}

// RequireManager middleware shorthand for requiring manager role or higher
func (m *AuthMiddleware) RequireManager() gin.HandlerFunc {
	return m.RequireRole(models.RoleManager)
}

// RequireActivation middleware checks if user account is active
func (m *AuthMiddleware) RequireActivation() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user from context (should be set by RequireAuth middleware)
		userInterface, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authentication required",
			})
			c.Abort()
			return
		}

		user, ok := userInterface.(*models.User)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Invalid user context",
			})
			c.Abort()
			return
		}

		// Check if user is active
		if !user.IsActive {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Account is not active",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// GetUserFromContext extracts the authenticated user from gin context
func GetUserFromContext(c *gin.Context) (*models.User, bool) {
	userInterface, exists := c.Get("user")
	if !exists {
		return nil, false
	}

	user, ok := userInterface.(*models.User)
	return user, ok
}

// GetUserIDFromContext extracts the authenticated user ID from gin context
func GetUserIDFromContext(c *gin.Context) (uuid.UUID, bool) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		return uuid.Nil, false
	}

	userID, ok := userIDInterface.(uuid.UUID)
	return userID, ok
}

// GetUserRoleFromContext extracts the authenticated user role from gin context
func GetUserRoleFromContext(c *gin.Context) (models.Role, bool) {
	roleInterface, exists := c.Get("user_role")
	if !exists {
		return "", false
	}

	role, ok := roleInterface.(models.Role)
	return role, ok
}

// Cors middleware for handling Cross-Origin Resource Sharing
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// You can configure allowed origins here
		allowedOrigins := []string{
			"http://localhost:3000",
			"http://localhost:3001",
			"http://localhost:5000",
			"http://127.0.0.1:3000",
		}

		isAllowed := false
		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin {
				isAllowed = true
				break
			}
		}

		if isAllowed {
			c.Header("Access-Control-Allow-Origin", origin)
		}

		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// RequestLogger middleware logs HTTP requests
func RequestLogger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format("2006/01/02 - 15:04:05"),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	})
}

// SecurityHeaders middleware adds security headers
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Content-Security-Policy", "default-src 'self'")
		c.Next()
	}
}

// RateLimiter middleware for basic rate limiting
// Note: In production, use a proper rate limiter like go-redis/redis_rate
func RateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		// This is a placeholder for rate limiting
		// In a real implementation, you'd use a proper rate limiter
		c.Next()
	}
}

// SessionUpdater middleware updates session information
func (m *AuthMiddleware) SessionUpdater() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Update session with IP and User-Agent if user is authenticated
		userInterface, exists := c.Get("user")
		if exists {
			if _, ok := userInterface.(*models.User); ok {
				// You could update session information here
				// For now, we'll just add the info to context
				c.Set("client_ip", c.ClientIP())
				c.Set("user_agent", c.Request.UserAgent())
			}
		}
		c.Next()
	}
}
