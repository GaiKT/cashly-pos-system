package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

// PasswordManager handles password operations
type PasswordManager struct {
	saltRounds int
}

// NewPasswordManager creates a new password manager instance
func NewPasswordManager(saltRounds int) *PasswordManager {
	if saltRounds < 10 {
		saltRounds = 12 // Default to 12 rounds for security
	}
	return &PasswordManager{
		saltRounds: saltRounds,
	}
}

// HashPassword hashes a password using bcrypt
func (pm *PasswordManager) HashPassword(password string) (string, error) {
	if err := pm.ValidatePassword(password); err != nil {
		return "", err
	}

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), pm.saltRounds)
	if err != nil {
		return "", err
	}

	return string(hashedBytes), nil
}

// VerifyPassword verifies a password against a hash
func (pm *PasswordManager) VerifyPassword(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

// ValidatePassword validates password strength
func (pm *PasswordManager) ValidatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	if len(password) > 128 {
		return errors.New("password must be less than 128 characters long")
	}

	// Check for at least one uppercase letter
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	if !hasUpper {
		return errors.New("password must contain at least one uppercase letter")
	}

	// Check for at least one lowercase letter
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	if !hasLower {
		return errors.New("password must contain at least one lowercase letter")
	}

	// Check for at least one number
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	if !hasNumber {
		return errors.New("password must contain at least one number")
	}

	// Check for at least one special character
	hasSpecial := regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`).MatchString(password)
	if !hasSpecial {
		return errors.New("password must contain at least one special character")
	}

	return nil
}

// GenerateRandomPassword generates a cryptographically secure random password
func (pm *PasswordManager) GenerateRandomPassword(length int) (string, error) {
	if length < 8 {
		length = 12 // Default to 12 characters
	}

	// Character sets
	uppercase := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	lowercase := "abcdefghijklmnopqrstuvwxyz"
	numbers := "0123456789"
	special := "!@#$%^&*()_+-=[]{}|;:,.<>?"
	allChars := uppercase + lowercase + numbers + special

	password := make([]byte, length)

	// Ensure at least one character from each required set
	password[0] = uppercase[randInt(len(uppercase))]
	password[1] = lowercase[randInt(len(lowercase))]
	password[2] = numbers[randInt(len(numbers))]
	password[3] = special[randInt(len(special))]

	// Fill the rest with random characters
	for i := 4; i < length; i++ {
		password[i] = allChars[randInt(len(allChars))]
	}

	// Shuffle the password
	for i := len(password) - 1; i > 0; i-- {
		j := randInt(i + 1)
		password[i], password[j] = password[j], password[i]
	}

	return string(password), nil
}

// GenerateResetToken generates a secure token for password reset
func (pm *PasswordManager) GenerateResetToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// GenerateEmailVerificationToken generates a secure token for email verification
func (pm *PasswordManager) GenerateEmailVerificationToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// Helper function to generate cryptographically secure random integers
func randInt(max int) int {
	if max <= 0 {
		return 0
	}

	bytes := make([]byte, 4)
	rand.Read(bytes)

	// Convert bytes to int and ensure it's within range
	n := int(bytes[0])<<24 | int(bytes[1])<<16 | int(bytes[2])<<8 | int(bytes[3])
	if n < 0 {
		n = -n
	}
	return n % max
}

// ValidateEmail validates email format
func ValidateEmail(email string) error {
	if email == "" {
		return errors.New("email is required")
	}

	// Simple email validation regex
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return errors.New("invalid email format")
	}

	if len(email) > 254 {
		return errors.New("email is too long")
	}

	return nil
}
