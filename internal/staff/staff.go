// Package staff owns the protected account-validation backoffice operations
// and their audit trail.
package staff

import (
	"context"
	"errors"
	"time"

	"maulid.id/internal/account"
)

// AccountSummary is the staff-facing account projection. Email is visible only
// to authenticated staff.
type AccountSummary struct {
	ID               string                   `json:"id"`
	Email            string                   `json:"email"`
	Name             string                   `json:"name"`
	Role             account.Role             `json:"role"`
	ValidationStatus account.ValidationStatus `json:"validation_status"`
	CreatedAt        time.Time                `json:"created_at"`
}

// AuditEntry records one administrative change.
type AuditEntry struct {
	ID              string    `json:"id"`
	ActorID         string    `json:"actor_id"`
	ActorName       string    `json:"actor_name"`
	TargetAccountID string    `json:"target_account_id"`
	TargetEmail     string    `json:"target_email"`
	Action          string    `json:"action"`
	FromStatus      string    `json:"from_status"`
	ToStatus        string    `json:"to_status"`
	CreatedAt       time.Time `json:"created_at"`
}

// ListFilter narrows the account list.
type ListFilter struct {
	ValidationStatus account.ValidationStatus
	Query            string
}

// Store is the staff feature's persistence port.
type Store interface {
	ListAccounts(ctx context.Context, filter ListFilter, limit int) ([]AccountSummary, error)
	SetValidation(ctx context.Context, actorID, targetID string, status account.ValidationStatus, at time.Time) (AccountSummary, *AuditEntry, error)
	ListAudit(ctx context.Context, limit int) ([]AuditEntry, error)
}

var (
	// ErrForbidden is returned for non-staff principals.
	ErrForbidden = errors.New("staff access required")
	// ErrAccountNotFound is returned for unknown target accounts.
	ErrAccountNotFound = errors.New("account not found")
)

// InputError is a safe validation failure for the backoffice form.
type InputError struct {
	Message string
}

func (e *InputError) Error() string { return e.Message }

// Service enforces staff authorization itself, independent of the HTTP layer.
type Service struct {
	store Store
	now   func() time.Time
}

// NewService constructs the staff service.
func NewService(store Store) *Service {
	return &Service{store: store, now: time.Now}
}

// List returns accounts visible to staff.
func (s *Service) List(ctx context.Context, principalRole string, filter ListFilter) ([]AccountSummary, error) {
	if principalRole != string(account.RoleStaff) {
		return nil, ErrForbidden
	}
	if filter.ValidationStatus != "" {
		status, ok := account.ParseValidationStatus(string(filter.ValidationStatus))
		if !ok {
			return nil, &InputError{Message: "Status validasi tidak dikenal."}
		}
		filter.ValidationStatus = status
	}
	if len([]rune(filter.Query)) > 100 {
		return nil, &InputError{Message: "Pencarian maksimal 100 karakter."}
	}
	return s.store.ListAccounts(ctx, filter, 100)
}

// SetValidation changes an account's validation status and writes an audit
// entry in the same transaction.
func (s *Service) SetValidation(ctx context.Context, principalRole, actorID, targetID, rawStatus string) (AccountSummary, *AuditEntry, error) {
	if principalRole != string(account.RoleStaff) {
		return AccountSummary{}, nil, ErrForbidden
	}
	status, ok := account.ParseValidationStatus(rawStatus)
	if !ok {
		return AccountSummary{}, nil, &InputError{Message: "Status validasi tidak dikenal."}
	}
	if targetID == "" {
		return AccountSummary{}, nil, ErrAccountNotFound
	}
	return s.store.SetValidation(ctx, actorID, targetID, status, s.now().UTC())
}

// Audit returns recent administrative changes.
func (s *Service) Audit(ctx context.Context, principalRole string, limit int) ([]AuditEntry, error) {
	if principalRole != string(account.RoleStaff) {
		return nil, ErrForbidden
	}
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	return s.store.ListAudit(ctx, limit)
}
