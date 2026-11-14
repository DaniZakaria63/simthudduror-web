package app

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"maulid.id/internal/account"
	"maulid.id/internal/platform/authn"
	"maulid.id/internal/platform/config"
)

type fixture struct {
	t      *testing.T
	app    *App
	server *httptest.Server
	client *http.Client
	cfg    config.Config
}

func newFixture(t *testing.T) *fixture {
	t.Helper()
	cfg := config.Config{
		Addr: ":0", DBPath: filepath.Join(t.TempDir(), "maulid.db"),
		Environment: config.EnvironmentDevelopment, SessionTTL: time.Hour,
		CookieSecure: false, RateLimitMax: 1000, RateLimitFor: time.Minute,
	}
	return newFixtureWithConfig(t, cfg)
}

func newFixtureWithConfig(t *testing.T, cfg config.Config) *fixture {
	t.Helper()
	application, err := New(cfg, slog.New(slog.NewTextHandler(io.Discard, nil)))
	if err != nil {
		t.Fatalf("app.New: %v", err)
	}
	server := httptest.NewServer(application.Handler)
	t.Cleanup(server.Close)
	return &fixture{t: t, app: application, server: server, client: server.Client(), cfg: cfg}
}

func (f *fixture) request(method, path string, body any, cookie *http.Cookie, csrf, bearer string) *http.Response {
	f.t.Helper()
	var reader io.Reader
	if body != nil {
		raw, err := json.Marshal(body)
		if err != nil {
			f.t.Fatalf("marshal: %v", err)
		}
		reader = bytes.NewReader(raw)
	}
	req, err := http.NewRequest(method, f.server.URL+path, reader)
	if err != nil {
		f.t.Fatalf("request: %v", err)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if cookie != nil {
		req.AddCookie(cookie)
	}
	if csrf != "" {
		req.Header.Set("X-CSRF-Token", csrf)
	}
	if bearer != "" {
		req.Header.Set("Authorization", "Bearer "+bearer)
	}
	response, err := f.client.Do(req)
	if err != nil {
		f.t.Fatalf("do: %v", err)
	}
	return response
}

type sessionPayload struct {
	Account account.View `json:"account"`
	CSRF    string       `json:"csrf_token"`
	Token   string       `json:"token"`
}

func (f *fixture) register(email, name string) (sessionPayload, *http.Cookie) {
	f.t.Helper()
	response := f.request(http.MethodPost, "/api/v1/auth/register", map[string]string{
		"email": email, "name": name, "password": "kata-sandi-aman",
	}, nil, "", "")
	if response.StatusCode != http.StatusCreated {
		f.t.Fatalf("register status = %d", response.StatusCode)
	}
	defer response.Body.Close()
	var payload sessionPayload
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		f.t.Fatalf("decode: %v", err)
	}
	for _, cookie := range response.Cookies() {
		if cookie.Name == authn.CookieName {
			return payload, cookie
		}
	}
	f.t.Fatal("session cookie missing")
	return payload, nil
}

func (f *fixture) login(email, name string) (sessionPayload, *http.Cookie) {
	f.t.Helper()
	response := f.request(http.MethodPost, "/api/v1/auth/login", map[string]string{
		"email": email, "password": "kata-sandi-aman",
	}, nil, "", "")
	if response.StatusCode != http.StatusOK {
		f.t.Fatalf("login status = %d", response.StatusCode)
	}
	defer response.Body.Close()
	var payload sessionPayload
	json.NewDecoder(response.Body).Decode(&payload)
	for _, cookie := range response.Cookies() {
		if cookie.Name == authn.CookieName {
			return payload, cookie
		}
	}
	f.t.Fatal("session cookie missing")
	return payload, nil
}

func TestEndToEndValidationPublishingAndPersistence(t *testing.T) {
	fixture := newFixture(t)

	// Health endpoints are public.
	for _, path := range []string{"/healthz", "/readyz"} {
		response := fixture.request(http.MethodGet, path, nil, nil, "", "")
		if response.StatusCode != http.StatusOK {
			t.Fatalf("%s status = %d", path, response.StatusCode)
		}
		response.Body.Close()
	}

	// Anonymous visitors read an empty feed.
	response := fixture.request(http.MethodGet, "/api/v1/posts", nil, nil, "", "")
	if response.StatusCode != http.StatusOK {
		t.Fatalf("feed status = %d", response.StatusCode)
	}
	response.Body.Close()

	member, memberCookie := fixture.register("member@example.com", "Rina")
	response = fixture.request(http.MethodPost, "/api/v1/posts", map[string]string{
		"body": "Kabar Maulid dari komunitas kampung.",
	}, memberCookie, member.CSRF, "")
	if response.StatusCode != http.StatusForbidden {
		t.Fatalf("unvalidated publish status = %d", response.StatusCode)
	}
	response.Body.Close()

	// A staff operator validates the member.
	if _, err := fixture.app.Accounts.BootstrapStaff(context.Background(), "staff@example.com", "Staf", "kata-sandi-aman"); err != nil {
		t.Fatalf("bootstrap staff: %v", err)
	}
	staff, staffCookie := fixture.login("staff@example.com", "Staf")
	response = fixture.request(http.MethodPatch, "/api/v1/staff/accounts/"+member.Account.ID+"/validation",
		map[string]string{"validation_status": "validated"}, staffCookie, staff.CSRF, "")
	if response.StatusCode != http.StatusOK {
		t.Fatalf("validation update status = %d", response.StatusCode)
	}
	response.Body.Close()

	// The same member session can publish immediately; status comes from DB.
	response = fixture.request(http.MethodPost, "/api/v1/posts", map[string]string{
		"body": "Kabar Maulid dari komunitas kampung.",
	}, memberCookie, member.CSRF, "")
	if response.StatusCode != http.StatusCreated {
		t.Fatalf("validated publish status = %d body=%s", response.StatusCode, rawBody(response))
	}
	var created struct {
		Post struct {
			ID string `json:"id"`
		} `json:"post"`
	}
	json.NewDecoder(response.Body).Decode(&created)
	response.Body.Close()
	if created.Post.ID == "" {
		t.Fatal("created post id missing")
	}

	// Public feed contains the post but never the author's email.
	response = fixture.request(http.MethodGet, "/api/v1/posts", nil, nil, "", "")
	raw := rawBody(response)
	if !strings.Contains(raw, "Kabar Maulid") || strings.Contains(raw, "member@example.com") {
		t.Fatalf("public feed privacy/content failure: %s", raw)
	}

	// Any authenticated account may like, including an unvalidated one.
	other, otherCookie := fixture.register("other@example.com", "Bima")
	response = fixture.request(http.MethodPut, "/api/v1/posts/"+created.Post.ID+"/like", nil, otherCookie, other.CSRF, "")
	if response.StatusCode != http.StatusOK {
		t.Fatalf("unvalidated like status = %d", response.StatusCode)
	}
	response.Body.Close()

	// Revocation applies to the next publish call without a new login.
	response = fixture.request(http.MethodPatch, "/api/v1/staff/accounts/"+member.Account.ID+"/validation",
		map[string]string{"validation_status": "unvalidated"}, staffCookie, staff.CSRF, "")
	if response.StatusCode != http.StatusOK {
		t.Fatalf("revoke status = %d", response.StatusCode)
	}
	response.Body.Close()
	response = fixture.request(http.MethodPost, "/api/v1/posts", map[string]string{
		"body": "Kabar setelah validasi dicabut.",
	}, memberCookie, member.CSRF, "")
	if response.StatusCode != http.StatusForbidden {
		t.Fatalf("revoked publish status = %d", response.StatusCode)
	}
	response.Body.Close()

	// Staff-only data rejects the member.
	response = fixture.request(http.MethodGet, "/api/v1/staff/accounts", nil, memberCookie, "", "")
	if response.StatusCode != http.StatusForbidden {
		t.Fatalf("member staff list status = %d", response.StatusCode)
	}
	response.Body.Close()

	// Restart: recreate the app over the same file. Session, post and like survive.
	if err := fixture.app.Close(); err != nil {
		t.Fatalf("close app: %v", err)
	}
	fixture.server.Close()
	reopened := newFixtureWithConfig(t, fixture.cfg)
	response = reopened.request(http.MethodGet, "/api/v1/posts", nil, otherCookie, "", "")
	if response.StatusCode != http.StatusOK {
		t.Fatalf("feed after restart status = %d", response.StatusCode)
	}
	var feed struct {
		Posts []struct {
			ID        string `json:"id"`
			LikeCount int    `json:"like_count"`
			LikedByMe bool   `json:"liked_by_me"`
		} `json:"posts"`
	}
	json.NewDecoder(response.Body).Decode(&feed)
	response.Body.Close()
	if len(feed.Posts) != 1 || feed.Posts[0].LikeCount != 1 || !feed.Posts[0].LikedByMe {
		t.Fatalf("persistence failed: %+v", feed)
	}
	response = reopened.request(http.MethodGet, "/api/v1/auth/session", nil, memberCookie, "", "")
	if response.StatusCode != http.StatusOK {
		t.Fatalf("session did not survive restart: %d", response.StatusCode)
	}
	response.Body.Close()
}

func TestPublicAndAPIUnknownRoutes(t *testing.T) {
	fixture := newFixture(t)
	response := fixture.request(http.MethodGet, "/api/v1/unknown", nil, nil, "", "")
	if response.StatusCode != http.StatusNotFound {
		t.Fatalf("unknown API status = %d", response.StatusCode)
	}
	response.Body.Close()
	// The SPA fallback is served for browser routes (the frontend may be a
	// placeholder build in a clean checkout).
	response = fixture.request(http.MethodGet, "/akun", nil, nil, "", "")
	if response.StatusCode == http.StatusNotFound {
		t.Fatalf("browser route not handled")
	}
	response.Body.Close()
}

func rawBody(response *http.Response) string {
	defer response.Body.Close()
	raw, _ := io.ReadAll(response.Body)
	return string(raw)
}
