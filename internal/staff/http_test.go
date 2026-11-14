package staff

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"maulid.id/internal/account"
	"maulid.id/internal/platform/authn"
	"maulid.id/internal/platform/ratelimit"
)

type httpFixture struct {
	t      *testing.T
	server *httptest.Server
	client *http.Client
	staff  *account.Service
}

func newHTTPFixture(t *testing.T) *httpFixture {
	t.Helper()
	db := newStaffDB(t)
	accountRepo := account.NewRepository(db)
	accountService := account.NewService(accountRepo, time.Hour)
	staffService := NewService(NewRepository(db))
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	mux := http.NewServeMux()
	require := authn.Require(accountService)
	optional := authn.Optional(accountService)
	account.NewHandler(accountService, ratelimit.New(1000, time.Minute), logger, false, nil).RegisterRoutes(mux, require, optional)
	NewHandler(staffService, logger).RegisterRoutes(mux, require)
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)
	return &httpFixture{t: t, server: server, client: server.Client(), staff: accountService}
}

func (f *httpFixture) post(path string, body any, cookie *http.Cookie, csrf string) *http.Response {
	f.t.Helper()
	var reader io.Reader
	if body != nil {
		raw, _ := json.Marshal(body)
		reader = bytes.NewReader(raw)
	}
	method := http.MethodPost
	if path != "/api/v1/auth/register" && path != "/api/v1/auth/login" {
		method = http.MethodPatch
	}
	req, _ := http.NewRequest(method, f.server.URL+path, reader)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if cookie != nil {
		req.AddCookie(cookie)
	}
	if csrf != "" {
		req.Header.Set("X-CSRF-Token", csrf)
	}
	response, err := f.client.Do(req)
	if err != nil {
		f.t.Fatalf("request: %v", err)
	}
	return response
}

func (f *httpFixture) bootstrapAndLogin() (*http.Cookie, string) {
	f.t.Helper()
	if _, err := f.staff.BootstrapStaff(context.Background(), "staff@example.com", "Staf", "kata-sandi-aman"); err != nil {
		f.t.Fatalf("bootstrap: %v", err)
	}
	response := f.post("/api/v1/auth/login", map[string]string{
		"email": "staff@example.com", "password": "kata-sandi-aman",
	}, nil, "")
	if response.StatusCode != http.StatusOK {
		f.t.Fatalf("staff login status = %d", response.StatusCode)
	}
	defer response.Body.Close()
	var payload struct {
		CSRF string `json:"csrf_token"`
	}
	json.NewDecoder(response.Body).Decode(&payload)
	for _, cookie := range response.Cookies() {
		if cookie.Name == authn.CookieName {
			return cookie, payload.CSRF
		}
	}
	f.t.Fatal("staff cookie missing")
	return nil, ""
}

func TestHTTPStaffBackoffice(t *testing.T) {
	fixture := newHTTPFixture(t)
	staffCookie, csrf := fixture.bootstrapAndLogin()

	// Create a member through the public flow.
	response := fixture.post("/api/v1/auth/register", map[string]string{
		"email": "member@example.com", "name": "Member", "password": "kata-sandi-aman",
	}, nil, "")
	if response.StatusCode != http.StatusCreated {
		t.Fatalf("register status = %d", response.StatusCode)
	}
	var registered struct {
		Account account.View `json:"account"`
		CSRF    string       `json:"csrf_token"`
	}
	json.NewDecoder(response.Body).Decode(&registered)
	var memberCookie *http.Cookie
	for _, cookie := range response.Cookies() {
		if cookie.Name == authn.CookieName {
			memberCookie = cookie
		}
	}
	response.Body.Close()

	response = fixture.post("/api/v1/auth/login", map[string]string{
		"email": "member@example.com", "password": "kata-sandi-aman",
	}, nil, "")
	response.Body.Close()

	// Anonymous and member sessions cannot read staff data.
	req, _ := http.NewRequest(http.MethodGet, fixture.server.URL+"/api/v1/staff/accounts", nil)
	response, _ = fixture.client.Do(req)
	if response.StatusCode != http.StatusUnauthorized {
		t.Fatalf("anonymous staff list status = %d", response.StatusCode)
	}
	response.Body.Close()
	req, _ = http.NewRequest(http.MethodGet, fixture.server.URL+"/api/v1/staff/accounts", nil)
	req.AddCookie(memberCookie)
	response, _ = fixture.client.Do(req)
	if response.StatusCode != http.StatusForbidden {
		t.Fatalf("member staff list status = %d", response.StatusCode)
	}
	response.Body.Close()

	// Staff can list accounts, including email (staff-only projection).
	req, _ = http.NewRequest(http.MethodGet, fixture.server.URL+"/api/v1/staff/accounts", nil)
	req.AddCookie(staffCookie)
	response, _ = fixture.client.Do(req)
	if response.StatusCode != http.StatusOK {
		t.Fatalf("staff list status = %d", response.StatusCode)
	}
	var list struct {
		Accounts []AccountSummary `json:"accounts"`
	}
	json.NewDecoder(response.Body).Decode(&list)
	response.Body.Close()
	if len(list.Accounts) != 2 {
		t.Fatalf("staff list = %+v", list)
	}

	// CSRF is required, then the update and audit succeed.
	response = fixture.post("/api/v1/staff/accounts/"+registered.Account.ID+"/validation", map[string]string{"validation_status": "validated"}, staffCookie, "")
	if response.StatusCode != http.StatusForbidden {
		t.Fatalf("staff update without CSRF status = %d", response.StatusCode)
	}
	response.Body.Close()
	response = fixture.post("/api/v1/staff/accounts/"+registered.Account.ID+"/validation", map[string]string{"validation_status": "validated"}, staffCookie, csrf)
	if response.StatusCode != http.StatusOK {
		t.Fatalf("staff update status = %d", response.StatusCode)
	}
	var update struct {
		Account AccountSummary `json:"account"`
		Audit   *AuditEntry    `json:"audit"`
	}
	json.NewDecoder(response.Body).Decode(&update)
	response.Body.Close()
	if update.Account.ValidationStatus != account.ValidationValidated || update.Audit == nil {
		t.Fatalf("unexpected update: %+v", update)
	}
	if update.Audit.ActorName != "Staf" {
		t.Fatalf("update audit actor name = %q, want Staf", update.Audit.ActorName)
	}

	response = fixture.post("/api/v1/staff/accounts/"+registered.Account.ID+"/validation", map[string]string{"validation_status": "pending"}, staffCookie, csrf)
	if response.StatusCode != http.StatusBadRequest {
		t.Fatalf("invalid status update = %d", response.StatusCode)
	}
	response.Body.Close()
	response = fixture.post("/api/v1/staff/accounts/missing/validation", map[string]string{"validation_status": "validated"}, staffCookie, csrf)
	if response.StatusCode != http.StatusNotFound {
		t.Fatalf("missing account update = %d", response.StatusCode)
	}
	response.Body.Close()

	req, _ = http.NewRequest(http.MethodGet, fixture.server.URL+"/api/v1/staff/audit", nil)
	req.AddCookie(staffCookie)
	response, _ = fixture.client.Do(req)
	if response.StatusCode != http.StatusOK {
		t.Fatalf("audit status = %d", response.StatusCode)
	}
	var audit struct {
		Audit []AuditEntry `json:"audit"`
	}
	json.NewDecoder(response.Body).Decode(&audit)
	response.Body.Close()
	if len(audit.Audit) != 1 || audit.Audit[0].ActorName != "Staf" {
		t.Fatalf("unexpected audit: %+v", audit)
	}

	// The member's next authenticated request reflects the new status without a
	// new login, because the session endpoint reloads trusted persistence.
	req, _ = http.NewRequest(http.MethodGet, fixture.server.URL+"/api/v1/auth/session", nil)
	req.AddCookie(memberCookie)
	response, _ = fixture.client.Do(req)
	var session struct {
		Account account.View `json:"account"`
	}
	json.NewDecoder(response.Body).Decode(&session)
	response.Body.Close()
	if session.Account.ValidationStatus != account.ValidationValidated {
		t.Fatalf("stale validation in session: %+v", session)
	}
}
