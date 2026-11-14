// Package account owns accounts, credentials and browser/mobile sessions.
package account

import (
	"errors"
	"time"
)

// Role is the account's staff/member role. It is separate from validation.
type Role string

const (
	RoleMember Role = "member"
	RoleStaff  Role = "staff"
)

// ValidationStatus is the canonical publishing eligibility field.
type ValidationStatus string

const (
	ValidationUnvalidated ValidationStatus = "unvalidated"
	ValidationValidated   ValidationStatus = "validated"
)

// ParseValidationStatus accepts only the canonical values.
func ParseValidationStatus(raw string) (ValidationStatus, bool) {
	switch ValidationStatus(raw) {
	case ValidationUnvalidated:
		return ValidationUnvalidated, true
	case ValidationValidated:
		return ValidationValidated, true
	default:
		return "", false
	}
}

// Account is the persisted account record. PasswordHash is never serialised.
type Account struct {
	ID               string
	Email            string
	Name             string
	PasswordHash     string `json:"-"`
	Role             Role
	ValidationStatus ValidationStatus
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// Session is a server-side session record. The primary key is the SHA-256 of
// the opaque credential; the raw credential is never stored.
type Session struct {
	ID        string
	AccountID string
	CSRFToken string
	CreatedAt time.Time
	ExpiresAt time.Time
	RevokedAt *time.Time
}

// View is the account shape returned to its owner. Email is private account
// data and is only exposed to the account itself and to staff.
type View struct {
	ID               string           `json:"id"`
	Email            string           `json:"email"`
	Name             string           `json:"name"`
	Role             Role             `json:"role"`
	ValidationStatus ValidationStatus `json:"validation_status"`
	CreatedAt        time.Time        `json:"created_at"`
}

// ToView maps an account to its owner-facing representation.
func (a Account) ToView() View {
	return View{
		ID:               a.ID,
		Email:            a.Email,
		Name:             a.Name,
		Role:             a.Role,
		ValidationStatus: a.ValidationStatus,
		CreatedAt:        a.CreatedAt,
	}
}

var (
	// ErrEmailTaken is returned when registration would duplicate an email.
	ErrEmailTaken = errors.New("email already registered")
	// ErrInvalidCredentials covers unknown accounts and wrong passwords alike.
	ErrInvalidCredentials = errors.New("invalid credentials")
	// ErrAccountNotFound is returned by repository lookups.
	ErrAccountNotFound = errors.New("account not found")
)

// InputError is a safe validation failure suitable for a form response.
type InputError struct {
	Message string
}

func (e *InputError) Error() string { return e.Message }

func invalidInput(message string) error { return &InputError{Message: message} }
