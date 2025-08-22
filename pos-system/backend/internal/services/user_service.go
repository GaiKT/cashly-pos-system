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
)

var (
	ErrUserProfileNotFound    = errors.New("user profile not found")
	ErrCannotUpdateOwnRole    = errors.New("cannot update your own role")
	ErrCannotDeactivateAdmin  = errors.New("cannot deactivate admin user")
	ErrCannotDeleteOwnAccount = errors.New("cannot delete your own account")
	ErrSuperAdminRequired     = errors.New("super admin permissions required")
	ErrEmailUpdateNotAllowed  = errors.New("email update not allowed for this account type")
)

// UserService handles user management operations
type UserService struct {
	userRepo     repository.UserRepository
	accountRepo  repository.AccountRepository
	sessionRepo  repository.SessionRepository
	passwordRepo repository.PasswordRepository
	auditRepo    repository.AuditLogRepository
	db           *gorm.DB
}

// NewUserService creates a new user management service
func NewUserService(
	userRepo repository.UserRepository,
	accountRepo repository.AccountRepository,
	sessionRepo repository.SessionRepository,
	passwordRepo repository.PasswordRepository,
	auditRepo repository.AuditLogRepository,
	db *gorm.DB,
) *UserService {
	return &UserService{
		userRepo:     userRepo,
		accountRepo:  accountRepo,
		sessionRepo:  sessionRepo,
		passwordRepo: passwordRepo,
		auditRepo:    auditRepo,
		db:           db,
	}
}

// GetUserProfile retrieves a user's profile information
func (s *UserService) GetUserProfile(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserProfileNotFound
		}
		return nil, fmt.Errorf("failed to get user profile: %w", err)
	}

	return user, nil
}

// UpdateUserProfile updates a user's profile information
func (s *UserService) UpdateUserProfile(ctx context.Context, userID uuid.UUID, req *models.UpdateProfileRequest) (*models.User, error) {
	// Get current user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserProfileNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Update fields if provided
	if req.Name != nil {
		user.Name = *req.Name
	}
	if req.Avatar != nil {
		user.Avatar = req.Avatar
	}

	// Update user
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user profile: %w", err)
	}

	// Log the update
	s.logUserAction(ctx, userID, "profile_updated", fmt.Sprintf("User %s updated their profile", user.Email))

	return user, nil
}

// ListUsers retrieves a paginated list of users (admin only)
func (s *UserService) ListUsers(ctx context.Context, requestorID uuid.UUID, filters map[string]interface{}, pagination *models.PaginationQuery) ([]models.User, int64, error) {
	// Verify requestor has admin permissions
	requestor, err := s.userRepo.GetByID(ctx, requestorID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get requestor: %w", err)
	}

	if requestor.Role != models.RoleAdmin {
		return nil, 0, ErrInsufficientRole
	}

	users, total, err := s.userRepo.List(ctx, filters, pagination)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list users: %w", err)
	}

	// Log the action
	s.logUserAction(ctx, requestorID, "users_listed", fmt.Sprintf("Admin %s listed users", requestor.Email))

	return users, total, nil
}

// CreateUser creates a new user account (admin only)
func (s *UserService) CreateUser(ctx context.Context, requestorID uuid.UUID, req *models.CreateUserRequest) (*models.User, error) {
	// Verify requestor has admin permissions
	requestor, err := s.userRepo.GetByID(ctx, requestorID)
	if err != nil {
		return nil, fmt.Errorf("failed to get requestor: %w", err)
	}

	if requestor.Role != models.RoleAdmin {
		return nil, ErrInsufficientRole
	}

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
		IsActive: true,               // Default to active
	}

	// Set role if provided
	if req.Role != nil {
		user.Role = *req.Role
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Create password if provided
	if req.Password != nil && *req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}

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
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Log the action
	s.logUserAction(ctx, requestorID, "user_created", fmt.Sprintf("Admin %s created user %s with role %s", requestor.Email, user.Email, user.Role))

	return user, nil
}

// UpdateUser updates an existing user (admin only)
func (s *UserService) UpdateUser(ctx context.Context, requestorID uuid.UUID, targetUserID uuid.UUID, req *models.UpdateUserRequest) (*models.User, error) {
	// Verify requestor has admin permissions
	requestor, err := s.userRepo.GetByID(ctx, requestorID)
	if err != nil {
		return nil, fmt.Errorf("failed to get requestor: %w", err)
	}

	if requestor.Role != models.RoleAdmin {
		return nil, ErrInsufficientRole
	}

	// Get target user
	user, err := s.userRepo.GetByID(ctx, targetUserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserProfileNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Prevent updating own role (handled by separate method)
	// Role updates are handled by UpdateUserRole method

	// Update fields if provided
	changes := []string{}
	if req.Name != nil && *req.Name != user.Name {
		user.Name = *req.Name
		changes = append(changes, "name")
	}
	if req.IsActive != nil && *req.IsActive != user.IsActive {
		// Prevent deactivating admin user
		if user.Role == models.RoleAdmin && !*req.IsActive {
			return nil, ErrCannotDeactivateAdmin
		}
		user.IsActive = *req.IsActive
		changes = append(changes, "active_status")
	}

	// Update user if changes were made
	if len(changes) > 0 {
		if err := s.userRepo.Update(ctx, user); err != nil {
			return nil, fmt.Errorf("failed to update user: %w", err)
		}

		// If user was deactivated, revoke all their sessions
		if req.IsActive != nil && !*req.IsActive {
			s.sessionRepo.RevokeAllUserSessions(ctx, user.ID)
		}

		// Log the action
		s.logUserAction(ctx, requestorID, "user_updated", fmt.Sprintf("Admin %s updated user %s: %v", requestor.Email, user.Email, changes))
	}

	return user, nil
}

// UpdateUserRole updates a user's role (admin only)
func (s *UserService) UpdateUserRole(ctx context.Context, requestorID uuid.UUID, targetUserID uuid.UUID, req *models.UpdateUserRoleRequest) error {
	// Verify requestor has admin permissions
	requestor, err := s.userRepo.GetByID(ctx, requestorID)
	if err != nil {
		return fmt.Errorf("failed to get requestor: %w", err)
	}

	if requestor.Role != models.RoleAdmin {
		return ErrInsufficientRole
	}

	// Prevent updating own role
	if requestorID == targetUserID {
		return ErrCannotUpdateOwnRole
	}

	// Update user role
	if err := s.userRepo.UpdateRole(ctx, targetUserID, req.Role); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserProfileNotFound
		}
		return fmt.Errorf("failed to update user role: %w", err)
	}

	// Get updated user for logging
	user, _ := s.userRepo.GetByID(ctx, targetUserID)
	if user != nil {
		s.logUserAction(ctx, requestorID, "role_updated", fmt.Sprintf("Admin %s updated role for user %s to %s", requestor.Email, user.Email, req.Role))
	}

	return nil
}

// DeactivateUser deactivates a user account (admin only)
func (s *UserService) DeactivateUser(ctx context.Context, requestorID uuid.UUID, targetUserID uuid.UUID) error {
	// Verify requestor has admin permissions
	requestor, err := s.userRepo.GetByID(ctx, requestorID)
	if err != nil {
		return fmt.Errorf("failed to get requestor: %w", err)
	}

	if requestor.Role != models.RoleAdmin {
		return ErrInsufficientRole
	}

	// Get target user
	user, err := s.userRepo.GetByID(ctx, targetUserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserProfileNotFound
		}
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Prevent deactivating admin user
	if user.Role == models.RoleAdmin {
		return ErrCannotDeactivateAdmin
	}

	// Deactivate user
	if err := s.userRepo.SetActiveStatus(ctx, targetUserID, false); err != nil {
		return fmt.Errorf("failed to deactivate user: %w", err)
	}

	// Revoke all user sessions
	s.sessionRepo.RevokeAllUserSessions(ctx, targetUserID)

	// Log the action
	s.logUserAction(ctx, requestorID, "user_deactivated", fmt.Sprintf("Admin %s deactivated user %s", requestor.Email, user.Email))

	return nil
}

// ActivateUser activates a user account (admin only)
func (s *UserService) ActivateUser(ctx context.Context, requestorID uuid.UUID, targetUserID uuid.UUID) error {
	// Verify requestor has admin permissions
	requestor, err := s.userRepo.GetByID(ctx, requestorID)
	if err != nil {
		return fmt.Errorf("failed to get requestor: %w", err)
	}

	if requestor.Role != models.RoleAdmin {
		return ErrInsufficientRole
	}

	// Get target user
	user, err := s.userRepo.GetByID(ctx, targetUserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserProfileNotFound
		}
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Activate user
	if err := s.userRepo.SetActiveStatus(ctx, targetUserID, true); err != nil {
		return fmt.Errorf("failed to activate user: %w", err)
	}

	// Log the action
	s.logUserAction(ctx, requestorID, "user_activated", fmt.Sprintf("Admin %s activated user %s", requestor.Email, user.Email))

	return nil
}

// DeleteUser soft deletes a user account (admin only)
func (s *UserService) DeleteUser(ctx context.Context, requestorID uuid.UUID, targetUserID uuid.UUID) error {
	// Verify requestor has admin permissions
	requestor, err := s.userRepo.GetByID(ctx, requestorID)
	if err != nil {
		return fmt.Errorf("failed to get requestor: %w", err)
	}

	if requestor.Role != models.RoleAdmin {
		return ErrInsufficientRole
	}

	// Prevent deleting own account
	if requestorID == targetUserID {
		return ErrCannotDeleteOwnAccount
	}

	// Get target user
	user, err := s.userRepo.GetByID(ctx, targetUserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserProfileNotFound
		}
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Prevent deleting admin user
	if user.Role == models.RoleAdmin {
		return ErrCannotDeactivateAdmin
	}

	// Start transaction
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Revoke all user sessions
	if err := s.sessionRepo.RevokeAllUserSessions(ctx, targetUserID); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to revoke user sessions: %w", err)
	}

	// Soft delete user
	if err := s.userRepo.Delete(ctx, targetUserID); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete user: %w", err)
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Log the action
	s.logUserAction(ctx, requestorID, "user_deleted", fmt.Sprintf("Admin %s deleted user %s", requestor.Email, user.Email))

	return nil
}

// GetUserSessions retrieves active sessions for a user
func (s *UserService) GetUserSessions(ctx context.Context, requestorID uuid.UUID, targetUserID uuid.UUID) ([]models.Session, error) {
	// Users can view their own sessions, admins can view any user's sessions
	if requestorID != targetUserID {
		requestor, err := s.userRepo.GetByID(ctx, requestorID)
		if err != nil {
			return nil, fmt.Errorf("failed to get requestor: %w", err)
		}

		if requestor.Role != models.RoleAdmin {
			return nil, ErrInsufficientRole
		}
	}

	sessions, err := s.sessionRepo.GetByUserID(ctx, targetUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user sessions: %w", err)
	}

	return sessions, nil
}

// RevokeUserSession revokes a specific user session
func (s *UserService) RevokeUserSession(ctx context.Context, requestorID uuid.UUID, sessionID uuid.UUID) error {
	// Only admins can revoke specific sessions by ID
	// Regular users should use RevokeAllUserSessions for their own sessions
	requestor, err := s.userRepo.GetByID(ctx, requestorID)
	if err != nil {
		return fmt.Errorf("failed to get requestor: %w", err)
	}

	if requestor.Role != models.RoleAdmin {
		return ErrInsufficientRole
	}

	// Revoke session
	if err := s.sessionRepo.Delete(ctx, sessionID); err != nil {
		return fmt.Errorf("failed to revoke session: %w", err)
	}

	// Log the action
	s.logUserAction(ctx, requestorID, "session_revoked", fmt.Sprintf("Admin revoked session %s", sessionID))

	return nil
}

// RevokeAllUserSessions revokes all sessions for a user
func (s *UserService) RevokeAllUserSessions(ctx context.Context, requestorID uuid.UUID, targetUserID uuid.UUID) error {
	// Users can revoke their own sessions, admins can revoke any user's sessions
	if requestorID != targetUserID {
		requestor, err := s.userRepo.GetByID(ctx, requestorID)
		if err != nil {
			return fmt.Errorf("failed to get requestor: %w", err)
		}

		if requestor.Role != models.RoleAdmin {
			return ErrInsufficientRole
		}
	}

	// Revoke all sessions
	if err := s.sessionRepo.RevokeAllUserSessions(ctx, targetUserID); err != nil {
		return fmt.Errorf("failed to revoke all user sessions: %w", err)
	}

	// Log the action
	s.logUserAction(ctx, requestorID, "all_sessions_revoked", fmt.Sprintf("All sessions revoked for user %s", targetUserID))

	return nil
}

// GetUserAccounts retrieves OAuth accounts for a user
func (s *UserService) GetUserAccounts(ctx context.Context, requestorID uuid.UUID, targetUserID uuid.UUID) ([]models.Account, error) {
	// Users can view their own accounts, admins can view any user's accounts
	if requestorID != targetUserID {
		requestor, err := s.userRepo.GetByID(ctx, requestorID)
		if err != nil {
			return nil, fmt.Errorf("failed to get requestor: %w", err)
		}

		if requestor.Role != models.RoleAdmin {
			return nil, ErrInsufficientRole
		}
	}

	accounts, err := s.accountRepo.GetByUserID(ctx, targetUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user accounts: %w", err)
	}

	return accounts, nil
}

// GetUserStatistics retrieves user activity statistics (admin only)
func (s *UserService) GetUserStatistics(ctx context.Context, requestorID uuid.UUID) (*models.UserStatistics, error) {
	// Verify requestor has admin permissions
	requestor, err := s.userRepo.GetByID(ctx, requestorID)
	if err != nil {
		return nil, fmt.Errorf("failed to get requestor: %w", err)
	}

	if requestor.Role != models.RoleAdmin {
		return nil, ErrInsufficientRole
	}

	// Get user counts by role
	filters := map[string]interface{}{"is_active": true}
	allUsers, total, err := s.userRepo.List(ctx, filters, &models.PaginationQuery{Limit: 1000})
	if err != nil {
		return nil, fmt.Errorf("failed to get user statistics: %w", err)
	}

	stats := &models.UserStatistics{
		TotalUsers:    int(total),
		ActiveUsers:   0,
		InactiveUsers: 0,
		AdminUsers:    0,
		ManagerUsers:  0,
		CashierUsers:  0,
	}

	for _, user := range allUsers {
		if user.IsActive {
			stats.ActiveUsers++
		} else {
			stats.InactiveUsers++
		}

		switch user.Role {
		case models.RoleAdmin:
			stats.AdminUsers++
		case models.RoleManager:
			stats.ManagerUsers++
		case models.RoleCashier:
			stats.CashierUsers++
		}
	}

	// Log the action
	s.logUserAction(ctx, requestorID, "statistics_viewed", fmt.Sprintf("Admin %s viewed user statistics", requestor.Email))

	return stats, nil
}

// logUserAction logs user management actions for audit trail
func (s *UserService) logUserAction(ctx context.Context, userID uuid.UUID, actionType, description string) {
	if s.auditRepo == nil {
		return // Audit logging is optional
	}

	// Get user info for audit log
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return // Don't fail if we can't get user info
	}

	// Map string actions to AuditLogAction constants
	var action models.AuditLogAction
	switch actionType {
	case "profile_updated":
		action = models.AuditActionUpdateUser
	case "users_listed":
		action = models.AuditActionLogin // Use existing constant for now
	case "user_created":
		action = models.AuditActionCreateUser
	case "user_updated":
		action = models.AuditActionUpdateUser
	case "role_updated":
		action = models.AuditActionUpdateUser
	case "user_deactivated":
		action = models.AuditActionUpdateUser
	case "user_activated":
		action = models.AuditActionUpdateUser
	case "user_deleted":
		action = models.AuditActionDeleteUser
	case "session_revoked":
		action = models.AuditActionLogout
	case "all_sessions_revoked":
		action = models.AuditActionLogout
	case "statistics_viewed":
		action = models.AuditActionSystemConfig
	default:
		action = models.AuditActionSystemConfig
	}

	auditLog := &models.AuditLog{
		ID:        uuid.New(),
		UserID:    userID,
		UserName:  user.Name,
		UserRole:  user.Role,
		Action:    action,
		Resource:  "user_management",
		IPAddress: "", // Will be set by middleware
		UserAgent: "", // Will be set by middleware
		Timestamp: time.Now(),
	}

	// Log in background, don't fail the main operation if logging fails
	go func() {
		if err := s.auditRepo.Create(context.Background(), auditLog); err != nil {
			fmt.Printf("Failed to log audit action: %v\n", err)
		}
	}()
}
