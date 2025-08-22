package services

import (
	"gorm.io/gorm"

	"github.com/pos-system/backend/internal/repository"
	"github.com/pos-system/backend/pkg/auth"
)

// Services holds all service instances
type Services struct {
	Auth *AuthService
	User *UserService
}

// NewServices creates all service instances
func NewServices(repos *repository.Repositories, jwtManager *auth.JWTManager) *Services {
	return &Services{
		Auth: NewAuthService(
			repos.User,
			repos.Account,
			repos.Session,
			repos.Password,
			jwtManager,
			repos.DB,
		),
		User: NewUserService(
			repos.User,
			repos.Account,
			repos.Session,
			repos.Password,
			repos.AuditLog,
			repos.DB,
		),
	}
}

// ServiceDependencies holds external dependencies needed by services
type ServiceDependencies struct {
	JWTManager *auth.JWTManager
	DB         *gorm.DB
}

// NewServiceDependencies creates service dependencies
func NewServiceDependencies(jwtManager *auth.JWTManager, db *gorm.DB) *ServiceDependencies {
	return &ServiceDependencies{
		JWTManager: jwtManager,
		DB:         db,
	}
}
