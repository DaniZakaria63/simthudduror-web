package account

import (
	"context"
	"time"
)

// Repository is the persistence port owned by the account feature.
type Repository interface {
	CreateAccount(ctx context.Context, account Account) error
	CreateAccountAndSession(ctx context.Context, account Account, session Session) error
	CreateSession(ctx context.Context, session Session) error
	FindByEmail(ctx context.Context, email string) (Account, error)
	FindSessionAccount(ctx context.Context, sessionID string, now time.Time) (Account, Session, error)
	RevokeSession(ctx context.Context, sessionID string, revokedAt time.Time) error
	ValidationStatus(ctx context.Context, accountID string) (ValidationStatus, error)
}
