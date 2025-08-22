package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Role represents user roles in the system
type Role string

const (
	RoleAdmin   Role = "ADMIN"
	RoleManager Role = "MANAGER"
	RoleCashier Role = "CASHIER"
)

// User represents the main user model
type User struct {
	ID          uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Email       string     `json:"email" gorm:"uniqueIndex;not null"`
	Name        string     `json:"name" gorm:"not null"`
	Avatar      *string    `json:"avatar,omitempty"`
	Role        Role       `json:"role" gorm:"type:user_role;not null;default:'CASHIER'"`
	IsActive    bool       `json:"isActive" gorm:"not null;default:true"`
	LastLoginAt *time.Time `json:"lastLoginAt,omitempty"`
	CreatedAt   time.Time  `json:"createdAt" gorm:"not null;default:now()"`
	UpdatedAt   time.Time  `json:"updatedAt" gorm:"not null;default:now()"`

	// Relationships
	Accounts []Account `json:"accounts,omitempty" gorm:"foreignKey:UserID"`
	Sessions []Session `json:"sessions,omitempty" gorm:"foreignKey:UserID"`
	Password *Password `json:"password,omitempty" gorm:"foreignKey:UserID"`
}

// TableName specifies the table name for GORM
func (User) TableName() string {
	return "users"
}

// Account represents OAuth provider accounts
type Account struct {
	ID                uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	UserID            uuid.UUID `json:"userId" gorm:"type:uuid;not null;index"`
	Type              string    `json:"type" gorm:"not null"`     // "oauth", "email"
	Provider          string    `json:"provider" gorm:"not null"` // "google", "facebook", "email"
	ProviderAccountID string    `json:"providerAccountId" gorm:"not null"`
	RefreshToken      *string   `json:"refreshToken,omitempty" gorm:"type:text"`
	AccessToken       *string   `json:"accessToken,omitempty" gorm:"type:text"`
	ExpiresAt         *int64    `json:"expiresAt,omitempty"`
	TokenType         *string   `json:"tokenType,omitempty"`
	Scope             *string   `json:"scope,omitempty" gorm:"type:text"`
	IDToken           *string   `json:"idToken,omitempty" gorm:"type:text"`
	SessionState      *string   `json:"sessionState,omitempty" gorm:"type:text"`
	CreatedAt         time.Time `json:"createdAt" gorm:"not null;default:now()"`
	UpdatedAt         time.Time `json:"updatedAt" gorm:"not null;default:now()"`

	// Relationships
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// TableName specifies the table name for GORM
func (Account) TableName() string {
	return "accounts"
}

// Session represents user sessions
type Session struct {
	ID           uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	SessionToken string    `json:"sessionToken" gorm:"uniqueIndex;not null"`
	UserID       uuid.UUID `json:"userId" gorm:"type:uuid;not null;index"`
	ExpiresAt    time.Time `json:"expiresAt" gorm:"not null;index"`
	IPAddress    *string   `json:"ipAddress,omitempty" gorm:"type:inet"`
	UserAgent    *string   `json:"userAgent,omitempty" gorm:"type:text"`
	IsActive     bool      `json:"isActive" gorm:"not null;default:true;index"`
	CreatedAt    time.Time `json:"createdAt" gorm:"not null;default:now()"`
	UpdatedAt    time.Time `json:"updatedAt" gorm:"not null;default:now()"`

	// Relationships
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// TableName specifies the table name for GORM
func (Session) TableName() string {
	return "sessions"
}

// Password represents password-based authentication
type Password struct {
	ID                     uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	UserID                 uuid.UUID  `json:"userId" gorm:"type:uuid;not null;uniqueIndex"`
	HashedPassword         string     `json:"-" gorm:"not null"`
	ResetToken             *string    `json:"-"`
	ResetTokenExpiresAt    *time.Time `json:"-"`
	EmailVerificationToken *string    `json:"-"`
	EmailVerified          bool       `json:"emailVerified" gorm:"not null;default:false"`
	EmailVerifiedAt        *time.Time `json:"emailVerifiedAt,omitempty"`
	LastPasswordChange     time.Time  `json:"lastPasswordChange" gorm:"not null;default:now()"`
	CreatedAt              time.Time  `json:"createdAt" gorm:"not null;default:now()"`
	UpdatedAt              time.Time  `json:"updatedAt" gorm:"not null;default:now()"`

	// Relationships
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// TableName specifies the table name for GORM
func (Password) TableName() string {
	return "passwords"
}

// UserWithRelations represents a user with all related data loaded
type UserWithRelations struct {
	User     `gorm:"embedded"`
	Accounts []Account `json:"accounts,omitempty"`
	Sessions []Session `json:"sessions,omitempty"`
	Password *Password `json:"password,omitempty"`
}

// CreateUserRequest represents the request to create a new user
type CreateUserRequest struct {
	Email    string  `json:"email" binding:"required,email"`
	Name     string  `json:"name" binding:"required,min=2,max=100"`
	Password *string `json:"password,omitempty" binding:"omitempty,min=8"`
	Role     *Role   `json:"role,omitempty"`
	Avatar   *string `json:"avatar,omitempty"`
}

// UpdateUserRequest represents the request to update user information
type UpdateUserRequest struct {
	Name     *string `json:"name,omitempty" binding:"omitempty,min=2,max=100"`
	Avatar   *string `json:"avatar,omitempty"`
	IsActive *bool   `json:"isActive,omitempty"`
}

// UpdateUserRoleRequest represents the request to update user role
type UpdateUserRoleRequest struct {
	Role Role `json:"role" binding:"required"`
}

// ChangePasswordRequest represents the request to change password
type ChangePasswordRequest struct {
	CurrentPassword string `json:"currentPassword" binding:"required"`
	NewPassword     string `json:"newPassword" binding:"required,min=8"`
}

// ResetPasswordRequest represents the request to reset password
type ResetPasswordRequest struct {
	Email       string `json:"email" binding:"required,email"`
	Token       string `json:"token,omitempty"`
	NewPassword string `json:"newPassword,omitempty"`
}

// ConfirmResetPasswordRequest represents the request to confirm password reset
type ConfirmResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required,min=8"`
}

// LoginRequest represents the login request
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents the login response
type LoginResponse struct {
	User         User   `json:"user"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int    `json:"expiresIn"`
}

// RefreshTokenRequest represents the token refresh request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

// RefreshTokenResponse represents the token refresh response
type RefreshTokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int    `json:"expiresIn"`
}

// RegisterRequest represents the user registration request
type RegisterRequest struct {
	Name     string `json:"name" binding:"required,min=2,max=100"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Role     *Role  `json:"role,omitempty"`
}

// AuthResponse represents the authentication response
type AuthResponse struct {
	User         User   `json:"user"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int    `json:"expiresIn"`
}

// UpdateProfileRequest represents the request to update user profile
type UpdateProfileRequest struct {
	Name   *string `json:"name,omitempty" binding:"omitempty,min=2,max=100"`
	Avatar *string `json:"avatar,omitempty"`
}

// Helper methods for User model

// IsAdmin checks if user has admin role
func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

// IsManager checks if user has manager role or higher
func (u *User) IsManager() bool {
	return u.Role == RoleAdmin || u.Role == RoleManager
}

// IsCashier checks if user has cashier role or higher
func (u *User) IsCashier() bool {
	return u.Role == RoleAdmin || u.Role == RoleManager || u.Role == RoleCashier
}

// CanAccessResource checks if user can access a resource based on required role
func (u *User) CanAccessResource(requiredRole Role) bool {
	switch requiredRole {
	case RoleAdmin:
		return u.IsAdmin()
	case RoleManager:
		return u.IsManager()
	case RoleCashier:
		return u.IsCashier()
	default:
		return false
	}
}

// ToPublic returns a user without sensitive information
func (u *User) ToPublic() User {
	publicUser := *u
	// Remove any sensitive fields if needed
	return publicUser
}

// ValidateRole checks if a role string is valid
func ValidateRole(role string) bool {
	switch Role(role) {
	case RoleAdmin, RoleManager, RoleCashier:
		return true
	default:
		return false
	}
}

// BeforeCreate hook for User model
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	return nil
}

// BeforeUpdate hook for User model
func (u *User) BeforeUpdate(tx *gorm.DB) error {
	u.UpdatedAt = time.Now()
	return nil
}

// BeforeCreate hook for Account model
func (a *Account) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	a.CreatedAt = time.Now()
	a.UpdatedAt = time.Now()
	return nil
}

// BeforeUpdate hook for Account model
func (a *Account) BeforeUpdate(tx *gorm.DB) error {
	a.UpdatedAt = time.Now()
	return nil
}

// BeforeCreate hook for Session model
func (s *Session) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	s.CreatedAt = time.Now()
	s.UpdatedAt = time.Now()
	return nil
}

// BeforeUpdate hook for Session model
func (s *Session) BeforeUpdate(tx *gorm.DB) error {
	s.UpdatedAt = time.Now()
	return nil
}

// BeforeCreate hook for Password model
func (p *Password) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
	return nil
}

// BeforeUpdate hook for Password model
func (p *Password) BeforeUpdate(tx *gorm.DB) error {
	p.UpdatedAt = time.Now()
	return nil
}

// UserStatistics represents user statistics for admin dashboard
type UserStatistics struct {
	TotalUsers    int `json:"totalUsers"`
	ActiveUsers   int `json:"activeUsers"`
	InactiveUsers int `json:"inactiveUsers"`
	AdminUsers    int `json:"adminUsers"`
	ManagerUsers  int `json:"managerUsers"`
	CashierUsers  int `json:"cashierUsers"`
}
