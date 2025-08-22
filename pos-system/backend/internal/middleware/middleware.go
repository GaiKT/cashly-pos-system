package middleware

import (
	"github.com/pos-system/backend/internal/services"
)

// Middleware holds all middleware instances
type Middleware struct {
	Auth *AuthMiddleware
}

// NewMiddleware creates all middleware instances
func NewMiddleware(services *services.Services) *Middleware {
	return &Middleware{
		Auth: NewAuthMiddleware(services.Auth),
	}
}
