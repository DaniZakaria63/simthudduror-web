package account

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

// SQLRepository is the SQLite adapter for the account feature.
type SQLRepository struct {
	db *sql.DB
}

// NewRepository returns a repository backed by the given database.
func NewRepository(db *sql.DB) *SQLRepository {
	return &SQLRepository{db: db}
}

func (r *SQLRepository) CreateAccount(ctx context.Context, account Account) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO accounts (id, email, name, password_hash, role, validation_status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		account.ID, account.Email, account.Name, account.PasswordHash, string(account.Role),
		string(account.ValidationStatus), formatTime(account.CreatedAt), formatTime(account.UpdatedAt))
	if err != nil {
		if isUniqueEmail(err) {
			return ErrEmailTaken
		}
		return fmt.Errorf("insert account: %w", err)
	}
	return nil
}

func (r *SQLRepository) CreateAccountAndSession(ctx context.Context, account Account, session Session) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin registration: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `
		INSERT INTO accounts (id, email, name, password_hash, role, validation_status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		account.ID, account.Email, account.Name, account.PasswordHash, string(account.Role),
		string(account.ValidationStatus), formatTime(account.CreatedAt), formatTime(account.UpdatedAt))
	if err != nil {
		if isUniqueEmail(err) {
			return ErrEmailTaken
		}
		return fmt.Errorf("insert account: %w", err)
	}
	if err := insertSession(ctx, tx, session); err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit registration: %w", err)
	}
	return nil
}

func (r *SQLRepository) CreateSession(ctx context.Context, session Session) error {
	if err := insertSession(ctx, r.db, session); err != nil {
		return err
	}
	return nil
}

type execer interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

func insertSession(ctx context.Context, db execer, session Session) error {
	_, err := db.ExecContext(ctx, `
		INSERT INTO sessions (id, account_id, csrf_token, created_at, expires_at)
		VALUES (?, ?, ?, ?, ?)`,
		session.ID, session.AccountID, session.CSRFToken, formatTime(session.CreatedAt), formatTime(session.ExpiresAt))
	if err != nil {
		return fmt.Errorf("insert session: %w", err)
	}
	return nil
}

func (r *SQLRepository) FindByEmail(ctx context.Context, email string) (Account, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, email, name, password_hash, role, validation_status, created_at, updated_at
		FROM accounts WHERE email = ? COLLATE NOCASE`, email)
	return scanAccount(row)
}

func (r *SQLRepository) FindSessionAccount(ctx context.Context, sessionID string, now time.Time) (Account, Session, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT a.id, a.email, a.name, a.password_hash, a.role, a.validation_status, a.created_at, a.updated_at,
		       s.id, s.account_id, s.csrf_token, s.created_at, s.expires_at
		FROM sessions s
		JOIN accounts a ON a.id = s.account_id
		WHERE s.id = ? AND s.revoked_at IS NULL AND s.expires_at > ?`,
		sessionID, formatTime(now))
	var account Account
	var session Session
	var createdAt, updatedAt, sessionCreated, expiresAt string
	err := row.Scan(&account.ID, &account.Email, &account.Name, &account.PasswordHash,
		&account.Role, &account.ValidationStatus, &createdAt, &updatedAt,
		&session.ID, &session.AccountID, &session.CSRFToken, &sessionCreated, &expiresAt)
	if errors.Is(err, sql.ErrNoRows) {
		return Account{}, Session{}, ErrAccountNotFound
	}
	if err != nil {
		return Account{}, Session{}, fmt.Errorf("query session: %w", err)
	}
	account.CreatedAt, err = parseTime(createdAt)
	if err != nil {
		return Account{}, Session{}, err
	}
	account.UpdatedAt, err = parseTime(updatedAt)
	if err != nil {
		return Account{}, Session{}, err
	}
	session.CreatedAt, err = parseTime(sessionCreated)
	if err != nil {
		return Account{}, Session{}, err
	}
	session.ExpiresAt, err = parseTime(expiresAt)
	if err != nil {
		return Account{}, Session{}, err
	}
	return account, session, nil
}

func (r *SQLRepository) RevokeSession(ctx context.Context, sessionID string, revokedAt time.Time) error {
	_, err := r.db.ExecContext(ctx, `UPDATE sessions SET revoked_at = ? WHERE id = ? AND revoked_at IS NULL`,
		formatTime(revokedAt), sessionID)
	if err != nil {
		return fmt.Errorf("revoke session: %w", err)
	}
	return nil
}

func (r *SQLRepository) ValidationStatus(ctx context.Context, accountID string) (ValidationStatus, error) {
	var status ValidationStatus
	err := r.db.QueryRowContext(ctx, `SELECT validation_status FROM accounts WHERE id = ?`, accountID).Scan(&status)
	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrAccountNotFound
	}
	if err != nil {
		return "", fmt.Errorf("query validation status: %w", err)
	}
	return status, nil
}

type scanner interface {
	Scan(dest ...any) error
}

func scanAccount(row scanner) (Account, error) {
	var account Account
	var createdAt, updatedAt string
	err := row.Scan(&account.ID, &account.Email, &account.Name, &account.PasswordHash,
		&account.Role, &account.ValidationStatus, &createdAt, &updatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return Account{}, ErrAccountNotFound
	}
	if err != nil {
		return Account{}, fmt.Errorf("scan account: %w", err)
	}
	account.CreatedAt, err = parseTime(createdAt)
	if err != nil {
		return Account{}, err
	}
	account.UpdatedAt, err = parseTime(updatedAt)
	if err != nil {
		return Account{}, err
	}
	return account, nil
}

// persistedTimeLayout is fixed-width UTC so SQLite TEXT comparisons (session
// expiry, post ordering) are chronologically correct at every precision.
const persistedTimeLayout = "2006-01-02T15:04:05.000000000Z"

func formatTime(t time.Time) string {
	return t.UTC().Format(persistedTimeLayout)
}

func parseTime(raw string) (time.Time, error) {
	parsed, err := time.Parse(persistedTimeLayout, raw)
	if err != nil {
		return time.Time{}, fmt.Errorf("parse persisted timestamp: %w", err)
	}
	return parsed, nil
}

func isUniqueEmail(err error) bool {
	return strings.Contains(strings.ToLower(err.Error()), "unique constraint failed: accounts.email")
}
