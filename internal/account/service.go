package account

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/mail"
	"strings"
	"sync"
	"time"
	"unicode"
	"unicode/utf8"

	"maulid.id/internal/platform/authn"
)

// AuthResult carries a freshly created session. The raw token is returned to
// the caller once and never persisted.
type AuthResult struct {
	Account   Account
	Token     string
	CSRFToken string
	ExpiresAt time.Time
}

// Service implements registration, login, session lookup, logout and the
// operator-only staff bootstrap.
type Service struct {
	repo       Repository
	sessionTTL time.Duration
	now        func() time.Time
}

// NewService constructs the account service.
func NewService(repo Repository, sessionTTL time.Duration) *Service {
	return &Service{repo: repo, sessionTTL: sessionTTL, now: time.Now}
}

var (
	dummyOnce sync.Once
	dummyHash string
)

func comparisonHash() string {
	dummyOnce.Do(func() {
		hash, err := HashPassword("not-a-real-password-placeholder")
		if err == nil {
			dummyHash = hash
		}
	})
	return dummyHash
}

// Register creates an unvalidated member account and a session.
func (s *Service) Register(ctx context.Context, email, name, password string) (AuthResult, error) {
	email, err := normalizeEmail(email)
	if err != nil {
		return AuthResult{}, err
	}
	name, err = normalizeName(name)
	if err != nil {
		return AuthResult{}, err
	}
	if err := validatePassword(password); err != nil {
		return AuthResult{}, err
	}
	hash, err := HashPassword(password)
	if err != nil {
		return AuthResult{}, err
	}
	now := s.now().UTC()
	acc := Account{
		ID:               newID(),
		Email:            email,
		Name:             name,
		PasswordHash:     hash,
		Role:             RoleMember,
		ValidationStatus: ValidationUnvalidated,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
	return s.createSessionFor(ctx, acc, true)
}

// Login verifies credentials and creates a new session. Unknown accounts are
// compared against a dummy hash so both failure paths cost the same.
func (s *Service) Login(ctx context.Context, email, password string) (AuthResult, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	acc, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, ErrAccountNotFound) {
			VerifyPassword(comparisonHash(), password)
			return AuthResult{}, ErrInvalidCredentials
		}
		return AuthResult{}, err
	}
	if !VerifyPassword(acc.PasswordHash, password) {
		return AuthResult{}, ErrInvalidCredentials
	}
	return s.createSessionFor(ctx, acc, false)
}

func (s *Service) createSessionFor(ctx context.Context, acc Account, persistAccount bool) (AuthResult, error) {
	token, err := NewCredential()
	if err != nil {
		return AuthResult{}, err
	}
	csrfToken, err := NewCredential()
	if err != nil {
		return AuthResult{}, err
	}
	now := s.now().UTC()
	session := Session{
		ID:        authn.HashToken(token),
		AccountID: acc.ID,
		CSRFToken: csrfToken,
		CreatedAt: now,
		ExpiresAt: now.Add(s.sessionTTL),
	}
	if persistAccount {
		if err := s.repo.CreateAccountAndSession(ctx, acc, session); err != nil {
			return AuthResult{}, err
		}
	} else if err := s.repo.CreateSession(ctx, session); err != nil {
		return AuthResult{}, err
	}
	return AuthResult{Account: acc, Token: token, CSRFToken: csrfToken, ExpiresAt: session.ExpiresAt}, nil
}

// Authenticate implements authn.Authenticator and reloads the account row, so
// role and validation changes apply on the next request.
func (s *Service) Authenticate(ctx context.Context, credential string) (authn.Principal, error) {
	if strings.TrimSpace(credential) == "" {
		return authn.Principal{}, authn.ErrInvalidCredential
	}
	acc, session, err := s.repo.FindSessionAccount(ctx, authn.HashToken(credential), s.now().UTC())
	if err != nil {
		if errors.Is(err, ErrAccountNotFound) {
			return authn.Principal{}, authn.ErrInvalidCredential
		}
		return authn.Principal{}, err
	}
	return authn.Principal{
		AccountID:        acc.ID,
		AccountCreatedAt: acc.CreatedAt,
		Email:            acc.Email,
		Name:             acc.Name,
		Role:             string(acc.Role),
		ValidationStatus: string(acc.ValidationStatus),
		SessionID:        session.ID,
		SessionCSRFHash:  authn.HashToken(session.CSRFToken),
		SessionCSRFToken: session.CSRFToken,
		SessionExpiresAt: session.ExpiresAt,
	}, nil
}

// Logout revokes the presented credential. Unknown credentials are ignored so
// logout stays idempotent.
func (s *Service) Logout(ctx context.Context, credential string) error {
	if strings.TrimSpace(credential) == "" {
		return nil
	}
	return s.repo.RevokeSession(ctx, authn.HashToken(credential), s.now().UTC())
}

// BootstrapStaff is the operator-only path for creating a validated staff
// account. It is never reachable over HTTP.
func (s *Service) BootstrapStaff(ctx context.Context, email, name, password string) (Account, error) {
	email, err := normalizeEmail(email)
	if err != nil {
		return Account{}, err
	}
	name, err = normalizeName(name)
	if err != nil {
		return Account{}, err
	}
	if err := validatePassword(password); err != nil {
		return Account{}, err
	}
	hash, err := HashPassword(password)
	if err != nil {
		return Account{}, err
	}
	now := s.now().UTC()
	acc := Account{
		ID:               newID(),
		Email:            email,
		Name:             name,
		PasswordHash:     hash,
		Role:             RoleStaff,
		ValidationStatus: ValidationValidated,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
	if err := s.repo.CreateAccount(ctx, acc); err != nil {
		return Account{}, err
	}
	return acc, nil
}

func normalizeEmail(raw string) (string, error) {
	email := strings.ToLower(strings.TrimSpace(raw))
	if email == "" || len(email) > 254 || strings.ContainsAny(email, " \t\r\n") {
		return "", invalidInput("Alamat email tidak valid.")
	}
	parsed, err := mail.ParseAddress(email)
	if err != nil || parsed.Address != email || parsed.Name != "" {
		return "", invalidInput("Alamat email tidak valid.")
	}
	return email, nil
}

func normalizeName(raw string) (string, error) {
	name := strings.TrimSpace(raw)
	if utf8.RuneCountInString(name) < 2 || utf8.RuneCountInString(name) > 80 {
		return "", invalidInput("Nama harus 2 sampai 80 karakter.")
	}
	if strings.ContainsFunc(name, unicode.IsControl) {
		return "", invalidInput("Nama mengandung karakter yang tidak diizinkan.")
	}
	return name, nil
}

func validatePassword(password string) error {
	if len(password) < 10 || len(password) > 128 {
		return invalidInput("Kata sandi harus 10 sampai 128 karakter.")
	}
	if strings.ContainsRune(password, '\x00') {
		return invalidInput("Kata sandi mengandung karakter yang tidak diizinkan.")
	}
	return nil
}

func newID() string {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		panic(fmt.Sprintf("generate id: %v", err))
	}
	return hex.EncodeToString(buf)
}
