package account

import (
	"database/sql"
	"path/filepath"
	"testing"
	"time"

	"maulid.id/internal/platform/database"
)

func newTestRepository(t *testing.T) (*SQLRepository, *sql.DB) {
	t.Helper()
	db, err := database.Open(filepath.Join(t.TempDir(), "maulid.db"))
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	return NewRepository(db), db
}

func newTestService(t *testing.T) (*Service, *SQLRepository) {
	t.Helper()
	repo, _ := newTestRepository(t)
	service := NewService(repo, time.Hour)
	service.now = func() time.Time { return time.Date(2026, 11, 14, 8, 0, 0, 0, time.UTC) }
	return service, repo
}
