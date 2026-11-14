package account

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"maulid.id/internal/platform/authn"
)

func TestPasswordHashRoundTrip(t *testing.T) {
	hash, err := HashPassword("kata-sandi-yang-kuat")
	if err != nil {
		t.Fatalf("HashPassword: %v", err)
	}
	if !strings.HasPrefix(hash, "$argon2id$") {
		t.Fatalf("unexpected encoding: %s", hash)
	}
	if !VerifyPassword(hash, "kata-sandi-yang-kuat") {
		t.Fatal("correct password rejected")
	}
	if VerifyPassword(hash, "kata-sandi-yang-salah") {
		t.Fatal("wrong password accepted")
	}
	for _, malformed := range []string{"", "plain", "$argon2i$v=19$m=1,t=1,p=1$aa$bb", "$argon2id$v=19$m=0,t=2,p=1$aa$bb", "$argon2id$v=19$m=19456,t=2,p=1$bad$bad"} {
		if VerifyPassword(malformed, "x") {
			t.Fatalf("malformed hash accepted: %q", malformed)
		}
	}
}

func TestNewCredentialIsRandom(t *testing.T) {
	first, err := NewCredential()
	if err != nil {
		t.Fatalf("NewCredential: %v", err)
	}
	second, err := NewCredential()
	if err != nil {
		t.Fatalf("NewCredential: %v", err)
	}
	if len(first) < 40 || first == second {
		t.Fatalf("credentials are not sufficiently random: %q %q", first, second)
	}
}

func TestRegisterCreatesUnvalidatedMember(t *testing.T) {
	service, _ := newTestService(t)
	ctx := context.Background()
	result, err := service.Register(ctx, "  Rina@Example.com ", "Rina", "kata-sandi-aman")
	if err != nil {
		t.Fatalf("Register: %v", err)
	}
	if result.Account.Email != "rina@example.com" {
		t.Fatalf("email not normalised: %s", result.Account.Email)
	}
	if result.Account.Role != RoleMember || result.Account.ValidationStatus != ValidationUnvalidated {
		t.Fatalf("registration assigned privileges: %+v", result.Account)
	}
	if result.Token == "" || result.CSRFToken == "" || result.ExpiresAt.IsZero() {
		t.Fatalf("incomplete session: %+v", result)
	}

	principal, err := service.Authenticate(ctx, result.Token)
	if err != nil {
		t.Fatalf("Authenticate: %v", err)
	}
	if principal.AccountID != result.Account.ID || principal.Kind != "" {
		t.Fatalf("unexpected principal: %+v", principal)
	}
	if principal.SessionCSRFToken != result.CSRFToken {
		t.Fatal("session CSRF token not preserved")
	}

	if _, err := service.Register(ctx, "rina@example.com", "Rina", "kata-sandi-aman"); !errors.Is(err, ErrEmailTaken) {
		t.Fatalf("duplicate email error = %v", err)
	}
}

func TestRegisterRejectsInvalidInput(t *testing.T) {
	service, _ := newTestService(t)
	ctx := context.Background()
	cases := []struct {
		name, email, person, password string
	}{
		{"bad email", "nope", "Nama", "kata-sandi-aman"},
		{"long email", strings.Repeat("a", 250) + "@x.id", "Nama", "kata-sandi-aman"},
		{"short name", "a@example.com", "A", "kata-sandi-aman"},
		{"long name", "a@example.com", strings.Repeat("a", 81), "kata-sandi-aman"},
		{"control name", "a@example.com", "Na\x00ma", "kata-sandi-aman"},
		{"short password", "a@example.com", "Nama", "pendek"},
		{"long password", "a@example.com", "Nama", strings.Repeat("a", 129)},
		{"nul password", "a@example.com", "Nama", "kata-sandi\x00aman"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := service.Register(ctx, tc.email, tc.person, tc.password)
			var input *InputError
			if !errors.As(err, &input) {
				t.Fatalf("expected input error, got %v", err)
			}
		})
	}
}

func TestLoginAndLogout(t *testing.T) {
	service, _ := newTestService(t)
	ctx := context.Background()
	registered, err := service.Register(ctx, "sari@example.com", "Sari", "kata-sandi-aman")
	if err != nil {
		t.Fatalf("Register: %v", err)
	}

	if _, err := service.Login(ctx, "sari@example.com", "salah-sekali"); !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("wrong password error = %v", err)
	}
	if _, err := service.Login(ctx, "missing@example.com", "salah-sekali"); !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("unknown account error = %v", err)
	}

	login, err := service.Login(ctx, "SARI@example.com", "kata-sandi-aman")
	if err != nil {
		t.Fatalf("Login: %v", err)
	}
	if login.Token == registered.Token {
		t.Fatal("login reused an old credential")
	}
	if _, err := service.Authenticate(ctx, login.Token); err != nil {
		t.Fatalf("authenticate new session: %v", err)
	}

	if err := service.Logout(ctx, login.Token); err != nil {
		t.Fatalf("Logout: %v", err)
	}
	if _, err := service.Authenticate(ctx, login.Token); !errors.Is(err, authn.ErrInvalidCredential) {
		t.Fatalf("revoked session accepted: %v", err)
	}
	if err := service.Logout(ctx, login.Token); err != nil {
		t.Fatalf("repeated logout should be idempotent: %v", err)
	}
	if err := service.Logout(ctx, ""); err != nil {
		t.Fatalf("empty logout should be a no-op: %v", err)
	}
	if _, err := service.Authenticate(ctx, ""); !errors.Is(err, authn.ErrInvalidCredential) {
		t.Fatalf("empty credential accepted: %v", err)
	}
}

func TestSessionExpiry(t *testing.T) {
	service, _ := newTestService(t)
	service.sessionTTL = 30 * time.Minute
	ctx := context.Background()
	result, err := service.Register(ctx, "exp@example.com", "Exp", "kata-sandi-aman")
	if err != nil {
		t.Fatalf("Register: %v", err)
	}
	service.now = func() time.Time { return time.Date(2026, 11, 14, 8, 31, 0, 0, time.UTC) }
	if _, err := service.Authenticate(ctx, result.Token); !errors.Is(err, authn.ErrInvalidCredential) {
		t.Fatalf("expired session accepted: %v", err)
	}
}

func TestBootstrapStaff(t *testing.T) {
	service, _ := newTestService(t)
	ctx := context.Background()
	acc, err := service.BootstrapStaff(ctx, "staf@example.com", "Staf", "kata-sandi-aman")
	if err != nil {
		t.Fatalf("BootstrapStaff: %v", err)
	}
	if acc.Role != RoleStaff || acc.ValidationStatus != ValidationValidated {
		t.Fatalf("unexpected bootstrap account: %+v", acc)
	}
	if _, err := service.Login(ctx, "staf@example.com", "kata-sandi-aman"); err != nil {
		t.Fatalf("bootstrap login: %v", err)
	}
	if _, err := service.BootstrapStaff(ctx, "staf@example.com", "Staf", "kata-sandi-aman"); !errors.Is(err, ErrEmailTaken) {
		t.Fatalf("duplicate bootstrap error = %v", err)
	}
	if _, err := service.BootstrapStaff(ctx, "bad", "Staf", "kata-sandi-aman"); err == nil {
		t.Fatal("invalid bootstrap input accepted")
	}
}

func TestValidationVocabulary(t *testing.T) {
	if status, ok := ParseValidationStatus("validated"); !ok || status != ValidationValidated {
		t.Fatal("validated not parsed")
	}
	if status, ok := ParseValidationStatus("unvalidated"); !ok || status != ValidationUnvalidated {
		t.Fatal("unvalidated not parsed")
	}
	for _, raw := range []string{"pending", "rejected", "", "Validated"} {
		if _, ok := ParseValidationStatus(raw); ok {
			t.Fatalf("non-canonical status accepted: %q", raw)
		}
	}
}
