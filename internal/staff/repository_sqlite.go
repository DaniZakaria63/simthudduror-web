package staff

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"maulid.id/internal/account"
)

const persistedTimeLayout = "2006-01-02T15:04:05.000000000Z"

// SQLRepository is the SQLite adapter for staff listings and audit records.
type SQLRepository struct {
	db *sql.DB
}

// NewRepository returns the staff repository.
func NewRepository(db *sql.DB) *SQLRepository {
	return &SQLRepository{db: db}
}

func (r *SQLRepository) ListAccounts(ctx context.Context, filter ListFilter, limit int) ([]AccountSummary, error) {
	pattern := "%" + escapeLike(filter.Query) + "%"
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, email, name, role, validation_status, created_at
		FROM accounts
		WHERE (? = '' OR validation_status = ?)
		  AND (? = '' OR name LIKE ? ESCAPE '\' OR email LIKE ? ESCAPE '\')
		ORDER BY created_at DESC, id DESC
		LIMIT ?`,
		string(filter.ValidationStatus), string(filter.ValidationStatus),
		filter.Query, pattern, pattern, limit)
	if err != nil {
		return nil, fmt.Errorf("list accounts: %w", err)
	}
	defer rows.Close()
	accounts := make([]AccountSummary, 0)
	for rows.Next() {
		summary, err := scanSummary(rows)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, summary)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate accounts: %w", err)
	}
	return accounts, nil
}

func (r *SQLRepository) SetValidation(ctx context.Context, actorID, targetID string, status account.ValidationStatus, at time.Time) (AccountSummary, *AuditEntry, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return AccountSummary{}, nil, fmt.Errorf("begin validation update: %w", err)
	}
	defer tx.Rollback()

	row := tx.QueryRowContext(ctx, `
		SELECT id, email, name, role, validation_status, created_at
		FROM accounts WHERE id = ?`, targetID)
	summary, err := scanSummary(row)
	if errors.Is(err, sql.ErrNoRows) {
		return AccountSummary{}, nil, ErrAccountNotFound
	}
	if err != nil {
		return AccountSummary{}, nil, err
	}
	from := summary.ValidationStatus
	if from == status {
		if err := tx.Commit(); err != nil {
			return AccountSummary{}, nil, fmt.Errorf("commit unchanged validation: %w", err)
		}
		return summary, nil, nil
	}
	if _, err := tx.ExecContext(ctx, `UPDATE accounts SET validation_status = ?, updated_at = ? WHERE id = ?`,
		string(status), formatTime(at), targetID); err != nil {
		return AccountSummary{}, nil, fmt.Errorf("update validation status: %w", err)
	}
	summary.ValidationStatus = status
	var actorName string
	if err := tx.QueryRowContext(ctx, `SELECT name FROM accounts WHERE id = ?`, actorID).Scan(&actorName); err != nil {
		return AccountSummary{}, nil, fmt.Errorf("load audit actor: %w", err)
	}
	audit := AuditEntry{
		ID: newID(), ActorID: actorID, ActorName: actorName, TargetAccountID: targetID,
		TargetEmail: summary.Email,
		Action:      "validation_status_changed", FromStatus: string(from), ToStatus: string(status), CreatedAt: at,
	}
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO audit_logs (id, actor_id, target_account_id, action, from_status, to_status, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		audit.ID, actorID, targetID, audit.Action, audit.FromStatus, audit.ToStatus, formatTime(at)); err != nil {
		return AccountSummary{}, nil, fmt.Errorf("insert audit: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return AccountSummary{}, nil, fmt.Errorf("commit validation update: %w", err)
	}
	return summary, &audit, nil
}

func (r *SQLRepository) ListAudit(ctx context.Context, limit int) ([]AuditEntry, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT l.id, l.actor_id, actor.name, l.target_account_id, target.email,
		       l.action, l.from_status, l.to_status, l.created_at
		FROM audit_logs l
		JOIN accounts actor ON actor.id = l.actor_id
		JOIN accounts target ON target.id = l.target_account_id
		ORDER BY l.created_at DESC, l.id DESC
		LIMIT ?`, limit)
	if err != nil {
		return nil, fmt.Errorf("list audit: %w", err)
	}
	defer rows.Close()
	entries := make([]AuditEntry, 0)
	for rows.Next() {
		var entry AuditEntry
		var createdAt string
		if err := rows.Scan(&entry.ID, &entry.ActorID, &entry.ActorName, &entry.TargetAccountID,
			&entry.TargetEmail, &entry.Action, &entry.FromStatus, &entry.ToStatus, &createdAt); err != nil {
			return nil, fmt.Errorf("scan audit: %w", err)
		}
		parsed, err := time.Parse(persistedTimeLayout, createdAt)
		if err != nil {
			return nil, fmt.Errorf("parse audit timestamp: %w", err)
		}
		entry.CreatedAt = parsed
		entries = append(entries, entry)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate audit: %w", err)
	}
	return entries, nil
}

type scanner interface {
	Scan(dest ...any) error
}

func scanSummary(row scanner) (AccountSummary, error) {
	var summary AccountSummary
	var createdAt string
	err := row.Scan(&summary.ID, &summary.Email, &summary.Name, &summary.Role, &summary.ValidationStatus, &createdAt)
	if err != nil {
		return AccountSummary{}, err
	}
	parsed, err := time.Parse(persistedTimeLayout, createdAt)
	if err != nil {
		return AccountSummary{}, fmt.Errorf("parse account timestamp: %w", err)
	}
	summary.CreatedAt = parsed
	return summary, nil
}

func formatTime(t time.Time) string { return t.UTC().Format(persistedTimeLayout) }

func escapeLike(value string) string {
	return strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`).Replace(value)
}
