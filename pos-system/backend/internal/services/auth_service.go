package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/pos-system/backend/internal/models"
	"github.com/pos-system/backend/internal/repository"
	"github.com/pos-system/backend/pkg/auth"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserNotFound       = errors.New("user not found")
	ErrUserNotActive      = errors.New("user account is not active")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidToken       = errors.New("invalid token")
	ErrTokenExpired       = errors.New("token expired")
	ErrInsufficientRole   = errors.New("insufficient role permissions")
)

// AuthService handles authentication operations
type AuthService struct {
	userRepo     repository.UserRepository
	accountRepo  repository.AccountRepository
	sessionRepo  repository.SessionRepository
	passwordRepo repository.PasswordRepository
	jwtManager   *auth.JWTManager
	db           *gorm.DB
}

// NewAuthService creates a new authentication service
func NewAuthService(
	userRepo repository.UserRepository,
	accountRepo repository.AccountRepository,
	sessionRepo repository.SessionRepository,
	passwordRepo repository.PasswordRepository,
	jwtManager *auth.JWTManager,
	db *gorm.DB,
) *AuthService {
	return &AuthService{
		userRepo:     userRepo,
		accountRepo:  accountRepo,
		sessionRepo:  sessionRepo,
		passwordRepo: passwordRepo,
		jwtManager:   jwtManager,
		db:           db,
	}
}

// Register creates a new user account with email/password
func (s *AuthService) Register(ctx context.Context, req *models.RegisterRequest) (*models.AuthResponse, error) {
	// Check if user already exists
	existingUser, _ := s.userRepo.GetByEmail(ctx, req.Email)
	if existingUser != nil {
		return nil, ErrEmailAlreadyExists
	}

	// Start transaction
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create user
	user := &models.User{
		ID:       uuid.New(),
		Email:    req.Email,
		Name:     req.Name,
		Role:     models.RoleCashier, // Default role
		IsActive: true,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create password record
	password := &models.Password{
		ID:             uuid.New(),
		UserID:         user.ID,
		HashedPassword: string(hashedPassword),
	}

	if err := s.passwordRepo.Create(ctx, password); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create password: %w", err)
	}

	// Create email account record
	account := &models.Account{
		ID:                uuid.New(),
		UserID:            user.ID,
		Type:              "email",
		Provider:          "email",
		ProviderAccountID: req.Email,
	}

	if err := s.accountRepo.Create(ctx, account); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create account: %w", err)
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Generate tokens
	accessToken, err := s.jwtManager.GenerateAccessToken(user.ID.String(), user.Email, string(user.Role), user.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken(user.ID.String(), user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Create session
	userAgent := ""
	ipAddress := ""
	session := &models.Session{
		ID:           uuid.New(),
		UserID:       user.ID,
		SessionToken: refreshToken,
		UserAgent:    &userAgent,                          // Will be set by middleware
		IPAddress:    &ipAddress,                          // Will be set by middleware
		ExpiresAt:    time.Now().Add(24 * time.Hour * 30), // 30 days
	}

	if err := s.sessionRepo.Create(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return &models.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    3600, // 1 hour
		User:         *user,
	}, nil
}

// Login authenticates a user with email/password
func (s *AuthService) Login(ctx context.Context, req *models.LoginRequest) (*models.AuthResponse, error) {
	// Get user by email
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Check if user is active
	if !user.IsActive {
		return nil, ErrUserNotActive
	}

	// Get password record
	password, err := s.passwordRepo.GetByUserID(ctx, user.ID)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(password.HashedPassword), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	// Update last login
	if err := s.userRepo.UpdateLastLogin(ctx, user.ID); err != nil {
		// Log error but don't fail login
		fmt.Printf("Failed to update last login for user %s: %v\n", user.ID, err)
	}

	// Generate tokens
	accessToken, err := s.jwtManager.GenerateAccessToken(user.ID.String(), user.Email, string(user.Role), user.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken(user.ID.String(), user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Create or update session
	userAgent := ""
	ipAddress := ""
	session := &models.Session{
		ID:           uuid.New(),
		UserID:       user.ID,
		SessionToken: refreshToken,
		UserAgent:    &userAgent,                          // Will be set by middleware
		IPAddress:    &ipAddress,                          // Will be set by middleware
		ExpiresAt:    time.Now().Add(24 * time.Hour * 30), // 30 days
	}

	if err := s.sessionRepo.Create(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return &models.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    3600, // 1 hour
		User:         *user,
	}, nil
}

// RefreshToken generates new tokens using a refresh token
func (s *AuthService) RefreshToken(ctx context.Context, req *models.RefreshTokenRequest) (*models.RefreshTokenResponse, error) {
	// Validate refresh token
	claims, err := s.jwtManager.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, ErrInvalidToken
	}

	// Get user
	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, ErrInvalidToken
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	// Check if user is still active
	if !user.IsActive {
		return nil, ErrUserNotActive
	}

	// Check if session exists and is valid
	session, err := s.sessionRepo.GetByToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, ErrInvalidToken
	}

	if session.ExpiresAt.Before(time.Now()) {
		// Clean up expired session
		s.sessionRepo.Delete(ctx, session.ID)
		return nil, ErrTokenExpired
	}

	// Generate new access token
	accessToken, err := s.jwtManager.GenerateAccessToken(user.ID.String(), user.Email, string(user.Role), user.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	return &models.RefreshTokenResponse{
		AccessToken: accessToken,
		ExpiresIn:   3600, // 1 hour
	}, nil
}

// Logout invalidates a user's session
func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	session, err := s.sessionRepo.GetByToken(ctx, refreshToken)
	if err != nil {
		// Session not found, consider it already logged out
		return nil
	}

	return s.sessionRepo.Delete(ctx, session.ID)
}

// GetUserFromToken extracts and validates user information from an access token
func (s *AuthService) GetUserFromToken(ctx context.Context, token string) (*models.User, error) {
	claims, err := s.jwtManager.ValidateAccessToken(token)
	if err != nil {
		return nil, ErrInvalidToken
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, ErrInvalidToken
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if !user.IsActive {
		return nil, ErrUserNotActive
	}

	return user, nil
}

// ValidateRole checks if a user has the required role
func (s *AuthService) ValidateRole(userRole models.Role, requiredRole models.Role) error {
	roleHierarchy := map[models.Role]int{
		models.RoleCashier: 1,
		models.RoleManager: 2,
		models.RoleAdmin:   3,
	}

	userLevel := roleHierarchy[userRole]
	requiredLevel := roleHierarchy[requiredRole]

	if userLevel < requiredLevel {
		return ErrInsufficientRole
	}

	return nil
}

// ChangePassword allows a user to change their password
func (s *AuthService) ChangePassword(ctx context.Context, userID uuid.UUID, req *models.ChangePasswordRequest) error {
	// Get current password
	currentPassword, err := s.passwordRepo.GetByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get current password: %w", err)
	}

	// Verify current password
	if err := bcrypt.CompareHashAndPassword([]byte(currentPassword.HashedPassword), []byte(req.CurrentPassword)); err != nil {
		return ErrInvalidCredentials
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash new password: %w", err)
	}

	// Update password
	currentPassword.HashedPassword = string(hashedPassword)
	if err := s.passwordRepo.Update(ctx, currentPassword); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	// Invalidate all existing sessions for the user
	sessions, err := s.sessionRepo.GetByUserID(ctx, userID)
	if err == nil {
		for _, session := range sessions {
			s.sessionRepo.Delete(ctx, session.ID)
		}
	}

	return nil
}

// ResetPassword initiates a password reset process
func (s *AuthService) ResetPassword(ctx context.Context, req *models.ResetPasswordRequest) error {
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		// Don't reveal if email exists or not
		return nil
	}

	if !user.IsActive {
		return nil
	}

	// Generate reset token (in a real implementation, this would be sent via email)
	resetToken := uuid.New().String()

	// Store reset token (this would typically be stored in a separate table)
	// For now, we'll just log it (in production, send via email)
	fmt.Printf("Password reset token for %s: %s\n", user.Email, resetToken)

	return nil
}

// ConfirmResetPassword completes the password reset process
func (s *AuthService) ConfirmResetPassword(ctx context.Context, req *models.ConfirmResetPasswordRequest) error {
	// In a real implementation, validate the reset token and extract user info
	// For now, we'll return an error since we need the email or user ID
	return fmt.Errorf("password reset confirmation not implemented - requires token validation")
}
