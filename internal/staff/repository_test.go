package staff

import (
	"context"
	"database/sql"
	"errors"
	"path/filepath"
	"testing"
	"time"

	"maulid.id/internal/account"
	"maulid.id/internal/platform/database"
)

func newStaffDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := database.Open(filepath.Join(t.TempDir(), "staff.db"))
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}

func seedUser(t *testing.T, db *sql.DB, id, email, name, role, status string, at time.Time) {
	t.Helper()
	_, err := db.Exec(`INSERT INTO accounts (id, email, name, password_hash, role, validation_status, created_at, updated_at)
		VALUES (?, ?, ?, 'hash', ?, ?, ?, ?)`, id, email, name, role, status, formatTime(at), formatTime(at))
	if err != nil {
		t.Fatalf("seed %s: %v", id, err)
	}
}

func TestRepositoryListAndSetValidationWithAudit(t *testing.T) {
	db := newStaffDB(t)
	repo := NewRepository(db)
	ctx := context.Background()
	base := time.Date(2026, 11, 14, 8, 0, 0, 0, time.UTC)
	seedUser(t, db, "staff", "staf@example.com", "Staf", "staff", "validated", base)
	seedUser(t, db, "member1", "satu@example.com", "Satu", "member", "unvalidated", base.Add(time.Hour))
	seedUser(t, db, "member2", "dua@example.com", "Dua_Spesial", "member", "validated", base.Add(2*time.Hour))

	all, err := repo.ListAccounts(ctx, ListFilter{}, 100)
	if err != nil || len(all) != 3 {
		t.Fatalf("list all = %d, %v", len(all), err)
	}
	unvalidated, err := repo.ListAccounts(ctx, ListFilter{ValidationStatus: account.ValidationUnvalidated}, 100)
	if err != nil || len(unvalidated) != 1 || unvalidated[0].ID != "member1" {
		t.Fatalf("filter status = %+v, %v", unvalidated, err)
	}
	escaped, err := repo.ListAccounts(ctx, ListFilter{Query: "Dua_"}, 100)
	if err != nil || len(escaped) != 1 || escaped[0].ID != "member2" {
		t.Fatalf("escape query = %+v, %v", escaped, err)
	}
	byEmail, err := repo.ListAccounts(ctx, ListFilter{Query: "satu@"}, 100)
	if err != nil || len(byEmail) != 1 || byEmail[0].ID != "member1" {
		t.Fatalf("email query = %+v, %v", byEmail, err)
	}

	updatedAt := base.Add(3 * time.Hour)
	updated, audit, err := repo.SetValidation(ctx, "staff", "member1", account.ValidationValidated, updatedAt)
	if err != nil || updated.ValidationStatus != account.ValidationValidated || audit == nil {
		t.Fatalf("SetValidation = %+v %+v %v", updated, audit, err)
	}
	if audit.FromStatus != "unvalidated" || audit.ToStatus != "validated" || audit.TargetEmail != "satu@example.com" {
		t.Fatalf("unexpected audit: %+v", audit)
	}
	if audit.ActorName != "Staf" {
		t.Fatalf("audit actor name = %q, want Staf", audit.ActorName)
	}

	same, audit, err := repo.SetValidation(ctx, "staff", "member1", account.ValidationValidated, updatedAt)
	if err != nil || same.ValidationStatus != account.ValidationValidated || audit != nil {
		t.Fatalf("unchanged update should not audit: %+v %+v %v", same, audit, err)
	}
	if _, _, err := repo.SetValidation(ctx, "staff", "missing", account.ValidationValidated, updatedAt); !errors.Is(err, ErrAccountNotFound) {
		t.Fatalf("missing target error = %v", err)
	}

	entries, err := repo.ListAudit(ctx, 10)
	if err != nil || len(entries) != 1 || entries[0].ActorName != "Staf" || entries[0].TargetEmail != "satu@example.com" {
		t.Fatalf("ListAudit = %+v, %v", entries, err)
	}

	status, err := NewRepository(db).ListAccounts(ctx, ListFilter{Query: "Satu"}, 10)
	if err != nil || status[0].ValidationStatus != account.ValidationValidated {
		t.Fatalf("validation not persisted: %+v, %v", status, err)
	}
}

func TestRepositoryRejectsMalformedTimestamps(t *testing.T) {
	db := newStaffDB(t)
	repo := NewRepository(db)
	ctx := context.Background()
	if _, err := db.Exec(`INSERT INTO accounts (id, email, name, password_hash, role, validation_status, created_at, updated_at)
		VALUES ('bad','bad@example.com','Bad','hash','member','unvalidated','nope','nope')`); err != nil {
		t.Fatalf("insert fixture: %v", err)
	}
	if _, err := repo.ListAccounts(ctx, ListFilter{}, 10); err == nil {
		t.Fatal("malformed timestamp accepted")
	}
}
