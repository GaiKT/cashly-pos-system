package auth

import (
	"testing"
	"time"
)

func TestJWTManager(t *testing.T) {
	// Create JWT manager
	jwtManager := NewJWTManager("test-secret-key", 1, 7) // 1 hour access, 7 days refresh

	// Test data
	userID := "user123"
	email := "test@example.com"
	role := "MANAGER"
	name := "Test User"

	// Test access token generation
	accessToken, err := jwtManager.GenerateAccessToken(userID, email, role, name)
	if err != nil {
		t.Fatalf("Failed to generate access token: %v", err)
	}

	if accessToken == "" {
		t.Fatal("Access token is empty")
	}

	// Test access token validation
	claims, err := jwtManager.ValidateAccessToken(accessToken)
	if err != nil {
		t.Fatalf("Failed to validate access token: %v", err)
	}

	// Verify claims
	if claims.UserID != userID {
		t.Errorf("Expected UserID %s, got %s", userID, claims.UserID)
	}
	if claims.Email != email {
		t.Errorf("Expected Email %s, got %s", email, claims.Email)
	}
	if claims.Role != role {
		t.Errorf("Expected Role %s, got %s", role, claims.Role)
	}
	if claims.Name != name {
		t.Errorf("Expected Name %s, got %s", name, claims.Name)
	}
	if claims.TokenType != "access" {
		t.Errorf("Expected TokenType 'access', got %s", claims.TokenType)
	}

	// Test refresh token generation
	refreshToken, err := jwtManager.GenerateRefreshToken(userID, email)
	if err != nil {
		t.Fatalf("Failed to generate refresh token: %v", err)
	}

	if refreshToken == "" {
		t.Fatal("Refresh token is empty")
	}

	// Test refresh token validation
	refreshClaims, err := jwtManager.ValidateRefreshToken(refreshToken)
	if err != nil {
		t.Fatalf("Failed to validate refresh token: %v", err)
	}

	if refreshClaims.UserID != userID {
		t.Errorf("Expected UserID %s, got %s", userID, refreshClaims.UserID)
	}
	if refreshClaims.TokenType != "refresh" {
		t.Errorf("Expected TokenType 'refresh', got %s", refreshClaims.TokenType)
	}
}

func TestPasswordManager(t *testing.T) {
	// Create password manager
	passwordManager := NewPasswordManager(12)

	// Test password validation
	validPassword := "TestPassword123!"
	err := passwordManager.ValidatePassword(validPassword)
	if err != nil {
		t.Fatalf("Valid password failed validation: %v", err)
	}

	// Test invalid passwords
	invalidPasswords := []string{
		"short",           // Too short
		"nouppercase123!", // No uppercase
		"NOLOWERCASE123!", // No lowercase
		"NoNumbers!",      // No numbers
		"NoSpecial123",    // No special characters
	}

	for _, password := range invalidPasswords {
		err := passwordManager.ValidatePassword(password)
		if err == nil {
			t.Errorf("Invalid password %s passed validation", password)
		}
	}

	// Test password hashing
	hashedPassword, err := passwordManager.HashPassword(validPassword)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	if hashedPassword == "" {
		t.Fatal("Hashed password is empty")
	}

	// Test password verification
	err = passwordManager.VerifyPassword(validPassword, hashedPassword)
	if err != nil {
		t.Fatalf("Failed to verify correct password: %v", err)
	}

	// Test wrong password verification
	err = passwordManager.VerifyPassword("WrongPassword123!", hashedPassword)
	if err == nil {
		t.Fatal("Wrong password verification should fail")
	}

	// Test random password generation
	randomPassword, err := passwordManager.GenerateRandomPassword(16)
	if err != nil {
		t.Fatalf("Failed to generate random password: %v", err)
	}

	if len(randomPassword) != 16 {
		t.Errorf("Expected password length 16, got %d", len(randomPassword))
	}

	// Test that generated password is valid
	err = passwordManager.ValidatePassword(randomPassword)
	if err != nil {
		t.Fatalf("Generated password is invalid: %v", err)
	}
}

func TestSessionManager(t *testing.T) {
	// Create session manager
	sessionManager := NewSessionManager(24) // 24 hours

	userID := "user123"
	ipAddress := "192.168.1.1"
	userAgent := "Mozilla/5.0 Test"

	// Test session generation
	session, err := sessionManager.GenerateSession(userID, ipAddress, userAgent)
	if err != nil {
		t.Fatalf("Failed to generate session: %v", err)
	}

	if session.SessionToken == "" {
		t.Fatal("Session token is empty")
	}
	if session.UserID != userID {
		t.Errorf("Expected UserID %s, got %s", userID, session.UserID)
	}
	if !session.IsActive {
		t.Error("New session should be active")
	}

	// Test session validation
	if !sessionManager.IsSessionValid(session) {
		t.Error("New session should be valid")
	}

	// Test session refresh
	originalExpiry := session.ExpiresAt
	time.Sleep(1 * time.Millisecond) // Ensure time difference
	refreshedSession := sessionManager.RefreshSession(session)

	if !refreshedSession.ExpiresAt.After(originalExpiry) {
		t.Error("Refreshed session expiry should be later than original")
	}

	// Test session revocation
	revokedSession := sessionManager.RevokeSession(session)
	if revokedSession.IsActive {
		t.Error("Revoked session should not be active")
	}
	if sessionManager.IsSessionValid(revokedSession) {
		t.Error("Revoked session should not be valid")
	}
}

func TestAuthManager(t *testing.T) {
	// Create auth manager
	authManager := NewAuthManager("test-secret", 1, 7, 24, 12)

	userID := "user123"
	email := "test@example.com"
	role := "MANAGER"
	name := "Test User"

	// Test token creation
	accessToken, refreshToken, err := authManager.CreateAuthTokens(userID, email, role, name)
	if err != nil {
		t.Fatalf("Failed to create auth tokens: %v", err)
	}

	if accessToken == "" || refreshToken == "" {
		t.Fatal("Auth tokens should not be empty")
	}

	// Test token validation
	result, err := authManager.ValidateAuthToken(accessToken)
	if err != nil {
		t.Fatalf("Failed to validate auth token: %v", err)
	}

	if !result.Valid {
		t.Error("Auth token should be valid")
	}
	if result.Claims.UserID != userID {
		t.Errorf("Expected UserID %s, got %s", userID, result.Claims.UserID)
	}

	// Test token refresh
	time.Sleep(1 * time.Millisecond) // Ensure different timestamps
	newAccessToken, newRefreshToken, err := authManager.RefreshAuthTokens(refreshToken)
	if err != nil {
		t.Fatalf("Failed to refresh auth tokens: %v", err)
	}

	if newAccessToken == "" || newRefreshToken == "" {
		t.Fatal("Refreshed tokens should not be empty")
	}

	// Verify new tokens are different
	if newAccessToken == accessToken {
		t.Error("New access token should be different from old one")
	}
	if newRefreshToken == refreshToken {
		t.Error("New refresh token should be different from old one")
	}
}
