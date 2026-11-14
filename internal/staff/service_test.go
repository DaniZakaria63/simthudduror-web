package staff

import (
	"context"
	"errors"
	"testing"
	"time"

	"maulid.id/internal/account"
)

type fakeStore struct {
	accounts  []AccountSummary
	audit     []AuditEntry
	err       error
	gotActor  string
	gotTarget string
	gotStatus account.ValidationStatus
}

func (f *fakeStore) ListAccounts(context.Context, ListFilter, int) ([]AccountSummary, error) {
	return f.accounts, f.err
}

func (f *fakeStore) SetValidation(_ context.Context, actorID, targetID string, status account.ValidationStatus, _ time.Time) (AccountSummary, *AuditEntry, error) {
	f.gotActor, f.gotTarget, f.gotStatus = actorID, targetID, status
	if f.err != nil {
		return AccountSummary{}, nil, f.err
	}
	return f.accounts[0], &f.audit[0], nil
}

func (f *fakeStore) ListAudit(context.Context, int) ([]AuditEntry, error) {
	return f.audit, f.err
}

func TestServiceEnforcesStaffRole(t *testing.T) {
	service := NewService(&fakeStore{})
	ctx := context.Background()
	if _, err := service.List(ctx, "member", ListFilter{}); !errors.Is(err, ErrForbidden) {
		t.Fatalf("member list error = %v", err)
	}
	if _, _, err := service.SetValidation(ctx, "member", "actor", "target", "validated"); !errors.Is(err, ErrForbidden) {
		t.Fatalf("member set error = %v", err)
	}
	if _, err := service.Audit(ctx, "member", 10); !errors.Is(err, ErrForbidden) {
		t.Fatalf("member audit error = %v", err)
	}
}

func TestServiceValidationAndAudit(t *testing.T) {
	store := &fakeStore{
		accounts: []AccountSummary{{ID: "target", ValidationStatus: account.ValidationValidated}},
		audit:    []AuditEntry{{ID: "audit"}},
	}
	service := NewService(store)
	ctx := context.Background()
	updated, audit, err := service.SetValidation(ctx, "staff", "actor", "target", "validated")
	if err != nil || updated.ID != "target" || audit == nil {
		t.Fatalf("SetValidation = %+v %+v %v", updated, audit, err)
	}
	if store.gotActor != "actor" || store.gotTarget != "target" || store.gotStatus != account.ValidationValidated {
		t.Fatalf("store received wrong request: %+v", store)
	}
	if _, _, err := service.SetValidation(ctx, "staff", "actor", "target", "pending"); err == nil {
		t.Fatal("non-canonical status accepted")
	}
	if _, _, err := service.SetValidation(ctx, "staff", "actor", "", "validated"); !errors.Is(err, ErrAccountNotFound) {
		t.Fatalf("empty target error = %v", err)
	}
	if _, err := service.List(ctx, "staff", ListFilter{ValidationStatus: "pending"}); err == nil {
		t.Fatal("invalid list filter accepted")
	}
	if _, err := service.List(ctx, "staff", ListFilter{Query: string(make([]rune, 101))}); err == nil {
		t.Fatal("oversized search accepted")
	}
	if _, err := service.Audit(ctx, "staff", 0); err != nil {
		t.Fatalf("audit: %v", err)
	}
}
