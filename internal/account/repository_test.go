package account

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	"maulid.id/internal/platform/authn"
)

func TestRepositoryAccountAndSessionLifecycle(t *testing.T) {
	repo, _ := newTestRepository(t)
	ctx := context.Background()
	now := time.Date(2026, 11, 14, 8, 0, 0, 123456789, time.UTC)
	acc := Account{
		ID: "acc-1", Email: "repo@example.com", Name: "Repo",
		PasswordHash: "hash", Role: RoleMember, ValidationStatus: ValidationUnvalidated,
		CreatedAt: now, UpdatedAt: now,
	}
	session := Session{
		ID: authn.HashToken("token"), AccountID: acc.ID, CSRFToken: "csrf",
		CreatedAt: now, ExpiresAt: now.Add(time.Hour),
	}
	if err := repo.CreateAccountAndSession(ctx, acc, session); err != nil {
		t.Fatalf("CreateAccountAndSession: %v", err)
	}
	if err := repo.CreateAccountAndSession(ctx, Account{
		ID: "acc-2", Email: "REPO@example.com", Name: "Other", PasswordHash: "hash",
		Role: RoleMember, ValidationStatus: ValidationUnvalidated, CreatedAt: now, UpdatedAt: now,
	}, Session{ID: "other", AccountID: "acc-2", CSRFToken: "csrf", CreatedAt: now, ExpiresAt: now}); !errors.Is(err, ErrEmailTaken) {
		t.Fatalf("duplicate error = %v", err)
	}

	found, err := repo.FindByEmail(ctx, "REPO@example.com")
	if err != nil || found.ID != "acc-1" || !found.CreatedAt.Equal(now) {
		t.Fatalf("FindByEmail = %+v, %v", found, err)
	}
	if _, err := repo.FindByEmail(ctx, "absent@example.com"); !errors.Is(err, ErrAccountNotFound) {
		t.Fatalf("missing account error = %v", err)
	}

	gotAcc, gotSession, err := repo.FindSessionAccount(ctx, session.ID, now.Add(30*time.Minute))
	if err != nil || gotAcc.ID != acc.ID || gotSession.CSRFToken != "csrf" {
		t.Fatalf("FindSessionAccount = %+v %+v %v", gotAcc, gotSession, err)
	}
	if _, _, err := repo.FindSessionAccount(ctx, session.ID, now.Add(2*time.Hour)); !errors.Is(err, ErrAccountNotFound) {
		t.Fatalf("expired session error = %v", err)
	}
	if err := repo.RevokeSession(ctx, session.ID, now); err != nil {
		t.Fatalf("RevokeSession: %v", err)
	}
	if _, _, err := repo.FindSessionAccount(ctx, session.ID, now); !errors.Is(err, ErrAccountNotFound) {
		t.Fatalf("revoked session error = %v", err)
	}

	status, err := repo.ValidationStatus(ctx, acc.ID)
	if err != nil || status != ValidationUnvalidated {
		t.Fatalf("ValidationStatus = %q, %v", status, err)
	}
	if _, err := repo.ValidationStatus(ctx, "absent"); !errors.Is(err, ErrAccountNotFound) {
		t.Fatalf("missing status error = %v", err)
	}
}

func TestRepositoryCreateAccountOnly(t *testing.T) {
	repo, _ := newTestRepository(t)
	now := time.Date(2026, 11, 14, 8, 0, 0, 0, time.UTC)
	acc := Account{ID: "staff", Email: "staff@example.com", Name: "Staff", PasswordHash: "hash",
		Role: RoleStaff, ValidationStatus: ValidationValidated, CreatedAt: now, UpdatedAt: now}
	if err := repo.CreateAccount(context.Background(), acc); err != nil {
		t.Fatalf("CreateAccount: %v", err)
	}
	if _, err := repo.FindByEmail(context.Background(), acc.Email); err != nil {
		t.Fatalf("created account not found: %v", err)
	}
}

func TestRepositoryRejectsMalformedTimestamp(t *testing.T) {
	repo, db := newTestRepository(t)
	_, err := db.Exec(`INSERT INTO accounts (id, email, name, password_hash, role, validation_status, created_at, updated_at)
		VALUES ('bad', 'bad@example.com', 'Bad', 'hash', 'member', 'unvalidated', 'not-a-time', 'not-a-time')`)
	if err != nil {
		t.Fatalf("insert fixture: %v", err)
	}
	if _, err := repo.FindByEmail(context.Background(), "bad@example.com"); err == nil {
		t.Fatal("malformed timestamp accepted")
	}
}

func TestPersistedTimesSortChronologicallyAtEveryPrecision(t *testing.T) {
	first := time.Date(2026, 11, 14, 8, 0, 0, 100000000, time.UTC)
	second := time.Date(2026, 11, 14, 8, 0, 0, 90000000, time.UTC)
	third := time.Date(2026, 11, 14, 8, 0, 1, 0, time.UTC)
	formattedFirst := formatTime(first)
	formattedSecond := formatTime(second)
	formattedThird := formatTime(third)
	if !(formattedSecond < formattedFirst && formattedFirst < formattedThird) {
		t.Fatalf("timestamps do not sort chronologically: %q %q %q", formattedSecond, formattedFirst, formattedThird)
	}
	for _, source := range []time.Time{first, second, third} {
		parsed, err := parseTime(formatTime(source))
		if err != nil || !parsed.Equal(source) {
			t.Fatalf("round trip failed: %v %v", parsed, err)
		}
	}
}

func TestAccountToViewHidesPassword(t *testing.T) {
	account := Account{ID: "id", Email: "a@example.com", Name: "A", Role: RoleMember,
		ValidationStatus: ValidationUnvalidated, PasswordHash: "secret"}
	view := account.ToView()
	if view.ID != "id" || view.ValidationStatus != ValidationUnvalidated {
		t.Fatalf("unexpected view: %+v", view)
	}
	serializedAccount, err := json.Marshal(account)
	if err != nil {
		t.Fatalf("marshal account: %v", err)
	}
	serializedView, err := json.Marshal(view)
	if err != nil {
		t.Fatalf("marshal view: %v", err)
	}
	if strings.Contains(string(serializedAccount), "secret") || strings.Contains(string(serializedView), "secret") {
		t.Fatalf("password hash was serialized: account=%s view=%s", serializedAccount, serializedView)
	}
}
