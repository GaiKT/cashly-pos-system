package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims represents the JWT claims structure
type Claims struct {
	UserID    string `json:"userId"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	Name      string `json:"name"`
	TokenType string `json:"tokenType"` // "access" or "refresh"
	jwt.RegisteredClaims
}

// JWTManager handles JWT token operations
type JWTManager struct {
	secretKey       string
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

// NewJWTManager creates a new JWT manager instance
func NewJWTManager(secretKey string, accessTTLHours, refreshTTLDays int) *JWTManager {
	return &JWTManager{
		secretKey:       secretKey,
		accessTokenTTL:  time.Duration(accessTTLHours) * time.Hour,
		refreshTokenTTL: time.Duration(refreshTTLDays) * 24 * time.Hour,
	}
}

// GenerateAccessToken generates a new access token
func (j *JWTManager) GenerateAccessToken(userID, email, role, name string) (string, error) {
	// Generate a unique ID for this token
	tokenID, err := generateUniqueID()
	if err != nil {
		return "", err
	}

	claims := Claims{
		UserID:    userID,
		Email:     email,
		Role:      role,
		Name:      name,
		TokenType: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        tokenID, // Unique ID for each token
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.accessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "pos-system",
			Subject:   userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

// GenerateRefreshToken generates a new refresh token
func (j *JWTManager) GenerateRefreshToken(userID, email string) (string, error) {
	// Generate a unique ID for this token
	tokenID, err := generateUniqueID()
	if err != nil {
		return "", err
	}

	claims := Claims{
		UserID:    userID,
		Email:     email,
		TokenType: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        tokenID, // Unique ID for each token
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.refreshTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "pos-system",
			Subject:   userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

// ValidateToken validates a JWT token and returns the claims
func (j *JWTManager) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid token signing method")
		}
		return []byte(j.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		// Check if token is expired
		if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
			return nil, errors.New("token has expired")
		}
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// ValidateAccessToken specifically validates access tokens
func (j *JWTManager) ValidateAccessToken(tokenString string) (*Claims, error) {
	claims, err := j.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	if claims.TokenType != "access" {
		return nil, errors.New("invalid token type, expected access token")
	}

	return claims, nil
}

// ValidateRefreshToken specifically validates refresh tokens
func (j *JWTManager) ValidateRefreshToken(tokenString string) (*Claims, error) {
	claims, err := j.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	if claims.TokenType != "refresh" {
		return nil, errors.New("invalid token type, expected refresh token")
	}

	return claims, nil
}

// ExtractTokenFromBearerString extracts the token from "Bearer <token>" format
func ExtractTokenFromBearerString(authHeader string) string {
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		return authHeader[7:]
	}
	return ""
}

// GetTokenClaims is a utility function to get claims without validation (for debugging)
func GetTokenClaims(tokenString string) (*Claims, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, &Claims{})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok {
		return claims, nil
	}

	return nil, errors.New("could not parse token claims")
}

// generateUniqueID generates a unique identifier for tokens
func generateUniqueID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
