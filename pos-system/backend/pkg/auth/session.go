package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"
)

// Session represents a user session
type Session struct {
	SessionToken string    `json:"sessionToken" gorm:"column:session_token;uniqueIndex"`
	UserID       string    `json:"userId" gorm:"column:user_id;index"`
	ExpiresAt    time.Time `json:"expiresAt" gorm:"column:expires_at;index"`
	CreatedAt    time.Time `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt    time.Time `json:"updatedAt" gorm:"column:updated_at"`
	IPAddress    string    `json:"ipAddress" gorm:"column:ip_address"`
	UserAgent    string    `json:"userAgent" gorm:"column:user_agent"`
	IsActive     bool      `json:"isActive" gorm:"column:is_active;index"`
}

// SessionManager handles session operations
type SessionManager struct {
	sessionTTL time.Duration
}

// NewSessionManager creates a new session manager
func NewSessionManager(sessionTTLHours int) *SessionManager {
	return &SessionManager{
		sessionTTL: time.Duration(sessionTTLHours) * time.Hour,
	}
}

// GenerateSession creates a new session
func (sm *SessionManager) GenerateSession(userID, ipAddress, userAgent string) (*Session, error) {
	sessionToken, err := generateSecureToken(32)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	session := &Session{
		SessionToken: sessionToken,
		UserID:       userID,
		ExpiresAt:    now.Add(sm.sessionTTL),
		CreatedAt:    now,
		UpdatedAt:    now,
		IPAddress:    ipAddress,
		UserAgent:    userAgent,
		IsActive:     true,
	}

	return session, nil
}

// IsSessionValid checks if a session is valid
func (sm *SessionManager) IsSessionValid(session *Session) bool {
	if session == nil {
		return false
	}

	return session.IsActive && time.Now().Before(session.ExpiresAt)
}

// RefreshSession extends the session expiration time
func (sm *SessionManager) RefreshSession(session *Session) *Session {
	session.ExpiresAt = time.Now().Add(sm.sessionTTL)
	session.UpdatedAt = time.Now()
	return session
}

// RevokeSession deactivates a session
func (sm *SessionManager) RevokeSession(session *Session) *Session {
	session.IsActive = false
	session.UpdatedAt = time.Now()
	return session
}

// generateSecureToken generates a cryptographically secure random token
func generateSecureToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// TokenValidationResult represents the result of token validation
type TokenValidationResult struct {
	Valid   bool     `json:"valid"`
	Claims  *Claims  `json:"claims,omitempty"`
	Session *Session `json:"session,omitempty"`
	Error   string   `json:"error,omitempty"`
}

// AuthManager combines JWT and session management
type AuthManager struct {
	jwtManager      *JWTManager
	sessionManager  *SessionManager
	passwordManager *PasswordManager
}

// NewAuthManager creates a comprehensive auth manager
func NewAuthManager(jwtSecret string, accessTTLHours, refreshTTLDays, sessionTTLHours, passwordSaltRounds int) *AuthManager {
	return &AuthManager{
		jwtManager:      NewJWTManager(jwtSecret, accessTTLHours, refreshTTLDays),
		sessionManager:  NewSessionManager(sessionTTLHours),
		passwordManager: NewPasswordManager(passwordSaltRounds),
	}
}

// GetJWTManager returns the JWT manager
func (am *AuthManager) GetJWTManager() *JWTManager {
	return am.jwtManager
}

// GetSessionManager returns the session manager
func (am *AuthManager) GetSessionManager() *SessionManager {
	return am.sessionManager
}

// GetPasswordManager returns the password manager
func (am *AuthManager) GetPasswordManager() *PasswordManager {
	return am.passwordManager
}

// ValidateAuthToken validates both JWT tokens and sessions
func (am *AuthManager) ValidateAuthToken(tokenString string) (*TokenValidationResult, error) {
	result := &TokenValidationResult{}

	// Try to validate as JWT token first
	claims, err := am.jwtManager.ValidateAccessToken(tokenString)
	if err != nil {
		result.Valid = false
		result.Error = err.Error()
		return result, err
	}

	result.Valid = true
	result.Claims = claims
	return result, nil
}

// CreateAuthTokens creates both access and refresh tokens
func (am *AuthManager) CreateAuthTokens(userID, email, role, name string) (accessToken, refreshToken string, err error) {
	accessToken, err = am.jwtManager.GenerateAccessToken(userID, email, role, name)
	if err != nil {
		return "", "", err
	}

	refreshToken, err = am.jwtManager.GenerateRefreshToken(userID, email)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// RefreshAuthTokens creates new tokens from a valid refresh token
func (am *AuthManager) RefreshAuthTokens(refreshTokenString string) (newAccessToken, newRefreshToken string, err error) {
	// Validate refresh token
	claims, err := am.jwtManager.ValidateRefreshToken(refreshTokenString)
	if err != nil {
		return "", "", errors.New("invalid refresh token")
	}

	// Generate new tokens
	newAccessToken, err = am.jwtManager.GenerateAccessToken(claims.UserID, claims.Email, claims.Role, claims.Name)
	if err != nil {
		return "", "", err
	}

	newRefreshToken, err = am.jwtManager.GenerateRefreshToken(claims.UserID, claims.Email)
	if err != nil {
		return "", "", err
	}

	return newAccessToken, newRefreshToken, nil
}
