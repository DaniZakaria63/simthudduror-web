package post

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"
	"time"

	"maulid.id/internal/platform/database"
)

func newTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := database.Open(filepath.Join(t.TempDir(), "posts.db"))
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}

func seedAccount(t *testing.T, db *sql.DB, id, name string) {
	t.Helper()
	now := formatTime(time.Date(2026, 11, 14, 8, 0, 0, 0, time.UTC))
	_, err := db.Exec(`INSERT INTO accounts (id, email, name, password_hash, role, validation_status, created_at, updated_at)
		VALUES (?, ?, ?, 'hash', 'member', 'unvalidated', ?, ?)`, id, id+"@example.com", name, now, now)
	if err != nil {
		t.Fatalf("seed account: %v", err)
	}
}

type fakeAuthorizer struct {
	statuses map[string]string
	err      error
}

func (f fakeAuthorizer) ValidationStatus(_ context.Context, accountID string) (string, error) {
	if f.err != nil {
		return "", f.err
	}
	return f.statuses[accountID], nil
}

func newTestService(t *testing.T) (*Service, *SQLRepository, *sql.DB) {
	t.Helper()
	db := newTestDB(t)
	repo := NewRepository(db)
	service := NewService(repo, fakeAuthorizer{statuses: map[string]string{}})
	service.now = func() time.Time { return time.Date(2026, 11, 14, 9, 0, 0, 123, time.UTC) }
	return service, repo, db
}
