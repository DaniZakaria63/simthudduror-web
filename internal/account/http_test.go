package account

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"maulid.id/internal/platform/authn"
	"maulid.id/internal/platform/ratelimit"
)

type testServer struct {
	t      *testing.T
	server *httptest.Server
	client *http.Client
}

func newAccountServer(t *testing.T, max int) (*testServer, *Service) {
	t.Helper()
	service, _ := newTestService(t)
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	handler := NewHandler(service, ratelimit.New(max, time.Minute), logger, false, []string{"https://trusted.example"})
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux, authn.Require(service), authn.Optional(service))
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)
	return &testServer{t: t, server: server, client: server.Client()}, service
}

func (s *testServer) request(method, path string, body any, headers map[string]string, cookie *http.Cookie) *http.Response {
	s.t.Helper()
	var reader io.Reader
	if body != nil {
		raw, err := json.Marshal(body)
		if err != nil {
			s.t.Fatalf("marshal body: %v", err)
		}
		reader = bytes.NewReader(raw)
	}
	req, err := http.NewRequest(method, s.server.URL+path, reader)
	if err != nil {
		s.t.Fatalf("new request: %v", err)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	if cookie != nil {
		req.AddCookie(cookie)
	}
	response, err := s.client.Do(req)
	if err != nil {
		s.t.Fatalf("do request: %v", err)
	}
	return response
}

func decodeBody(t *testing.T, response *http.Response) map[string]any {
	t.Helper()
	defer response.Body.Close()
	var payload map[string]any
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	return payload
}

func sessionCookie(t *testing.T, response *http.Response) *http.Cookie {
	t.Helper()
	for _, cookie := range response.Cookies() {
		if cookie.Name == authn.CookieName && cookie.Value != "" {
			return cookie
		}
	}
	t.Fatal("session cookie not set")
	return nil
}

func TestHTTPRegisterSessionLogout(t *testing.T) {
	server, _ := newAccountServer(t, 100)
	response := server.request(http.MethodPost, "/api/v1/auth/register", map[string]string{
		"email": "web@example.com", "name": "Web", "password": "kata-sandi-aman",
	}, nil, nil)
	if response.StatusCode != http.StatusCreated {
		t.Fatalf("register status = %d", response.StatusCode)
	}
	payload := decodeBody(t, response)
	csrf, _ := payload["csrf_token"].(string)
	if csrf == "" {
		t.Fatal("csrf token missing")
	}
	if _, leaked := payload["token"]; leaked {
		t.Fatal("browser registration returned a bearer token")
	}
	cookie := sessionCookie(t, response)
	if !cookie.HttpOnly || cookie.SameSite != http.SameSiteLaxMode {
		t.Fatalf("unsafe cookie: %+v", cookie)
	}
	accountPayload := payload["account"].(map[string]any)
	if accountPayload["validation_status"] != "unvalidated" || accountPayload["role"] != "member" {
		t.Fatalf("unexpected account: %+v", accountPayload)
	}

	response = server.request(http.MethodGet, "/api/v1/auth/session", nil, nil, cookie)
	if response.StatusCode != http.StatusOK {
		t.Fatalf("session status = %d", response.StatusCode)
	}
	sessionPayload := decodeBody(t, response)
	if got := sessionPayload["csrf_token"]; got != csrf {
		t.Fatalf("session csrf = %v, want %s", got, csrf)
	}
	if sessionPayload["account"].(map[string]any)["created_at"] != accountPayload["created_at"] {
		t.Fatalf("session view lost created_at: register=%v session=%v", accountPayload["created_at"], sessionPayload["account"])
	}

	response = server.request(http.MethodPost, "/api/v1/auth/logout", nil, nil, cookie)
	if response.StatusCode != http.StatusForbidden {
		t.Fatalf("logout without CSRF status = %d", response.StatusCode)
	}
	response = server.request(http.MethodPost, "/api/v1/auth/logout", nil, map[string]string{"X-CSRF-Token": csrf}, cookie)
	if response.StatusCode != http.StatusNoContent {
		t.Fatalf("logout status = %d", response.StatusCode)
	}
	response = server.request(http.MethodGet, "/api/v1/auth/session", nil, nil, cookie)
	if response.StatusCode != http.StatusUnauthorized {
		t.Fatalf("session after logout status = %d", response.StatusCode)
	}
}

func TestHTTPLoginGenericFailureAndMobileToken(t *testing.T) {
	server, _ := newAccountServer(t, 100)
	register := server.request(http.MethodPost, "/api/v1/auth/register", map[string]string{
		"email": "mobile@example.com", "name": "Mobile", "password": "kata-sandi-aman",
	}, map[string]string{"X-Maulid-Client": "mobile"}, nil)
	payload := decodeBody(t, register)
	token, _ := payload["token"].(string)
	if token == "" {
		t.Fatal("mobile registration did not return a token")
	}

	response := server.request(http.MethodPost, "/api/v1/auth/login", map[string]string{
		"email": "mobile@example.com", "password": "salah-total",
	}, map[string]string{"X-Maulid-Client": "mobile"}, nil)
	body := decodeBody(t, response)
	if response.StatusCode != http.StatusUnauthorized || !strings.Contains(body["error"].(map[string]any)["message"].(string), "Email atau kata sandi salah") {
		t.Fatalf("unexpected login failure: %d %+v", response.StatusCode, body)
	}

	response = server.request(http.MethodPost, "/api/v1/auth/login", map[string]string{
		"email": "mobile@example.com", "password": "kata-sandi-aman",
	}, map[string]string{"X-Maulid-Client": "mobile"}, nil)
	payload = decodeBody(t, response)
	if payload["token"] == "" {
		t.Fatal("mobile login missing bearer token")
	}
	bearer := payload["token"].(string)
	response = server.request(http.MethodGet, "/api/v1/auth/session", nil, map[string]string{"Authorization": "Bearer " + bearer}, nil)
	if response.StatusCode != http.StatusOK {
		t.Fatalf("bearer session status = %d", response.StatusCode)
	}
}

func TestHTTPAnonymousAndDuplicateRegistration(t *testing.T) {
	server, _ := newAccountServer(t, 100)
	response := server.request(http.MethodGet, "/api/v1/auth/session", nil, nil, nil)
	if response.StatusCode != http.StatusUnauthorized {
		t.Fatalf("anonymous session status = %d", response.StatusCode)
	}
	response.Body.Close()

	for i := 0; i < 2; i++ {
		response = server.request(http.MethodPost, "/api/v1/auth/register", map[string]string{
			"email": "dup@example.com", "name": "Dup", "password": "kata-sandi-aman",
		}, nil, nil)
		if i == 0 && response.StatusCode != http.StatusCreated {
			t.Fatalf("first registration status = %d", response.StatusCode)
		}
		if i == 1 && response.StatusCode != http.StatusConflict {
			t.Fatalf("duplicate registration status = %d", response.StatusCode)
		}
		response.Body.Close()
	}
}

func TestHTTPBadJSONAndOversizedBody(t *testing.T) {
	server, _ := newAccountServer(t, 100)
	response := server.request(http.MethodPost, "/api/v1/auth/login", map[string]string{
		"email": "a@example.com", "role": "staff",
	}, nil, nil)
	if response.StatusCode != http.StatusBadRequest {
		t.Fatalf("unknown privilege field status = %d", response.StatusCode)
	}
	response.Body.Close()

	oversized := `{"email":"` + strings.Repeat("a", 40<<10) + `@example.com"}`
	req, _ := http.NewRequest(http.MethodPost, server.server.URL+"/api/v1/auth/login", strings.NewReader(oversized))
	req.Header.Set("Content-Type", "application/json")
	response, err := server.client.Do(req)
	if err != nil {
		t.Fatalf("oversized request: %v", err)
	}
	if response.StatusCode != http.StatusRequestEntityTooLarge {
		t.Fatalf("oversized status = %d", response.StatusCode)
	}
	response.Body.Close()

	req, _ = http.NewRequest(http.MethodPost, server.server.URL+"/api/v1/auth/login", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "text/plain")
	response, err = server.client.Do(req)
	if err != nil {
		t.Fatalf("wrong content type request: %v", err)
	}
	if response.StatusCode != http.StatusUnsupportedMediaType {
		t.Fatalf("wrong content type status = %d", response.StatusCode)
	}
	response.Body.Close()
}

func TestHTTPOriginAndRateLimit(t *testing.T) {
	server, _ := newAccountServer(t, 2)
	response := server.request(http.MethodPost, "/api/v1/auth/login", map[string]string{
		"email": "a@example.com", "password": "kata-sandi-aman",
	}, map[string]string{"Origin": "https://evil.example"}, nil)
	if response.StatusCode != http.StatusForbidden {
		t.Fatalf("cross-origin status = %d", response.StatusCode)
	}
	response.Body.Close()

	for i := 0; i < 2; i++ {
		response = server.request(http.MethodPost, "/api/v1/auth/login", map[string]string{
			"email": "a@example.com", "password": "kata-sandi-aman",
		}, map[string]string{"Origin": "https://trusted.example"}, nil)
		if response.StatusCode != http.StatusUnauthorized {
			t.Fatalf("login attempt %d status = %d", i, response.StatusCode)
		}
		response.Body.Close()
	}
	response = server.request(http.MethodPost, "/api/v1/auth/login", map[string]string{
		"email": "a@example.com", "password": "kata-sandi-aman",
	}, nil, nil)
	if response.StatusCode != http.StatusTooManyRequests || response.Header.Get("Retry-After") == "" {
		t.Fatalf("rate limit status = %d retry=%q", response.StatusCode, response.Header.Get("Retry-After"))
	}
	response.Body.Close()
}

func TestHTTPServiceFailureIsGeneric(t *testing.T) {
	service, repo := newTestService(t)
	_ = repo
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	handler := NewHandler(service, ratelimit.New(100, time.Minute), logger, false, nil)
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux, authn.Require(service), authn.Optional(service))
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", strings.NewReader(`{"email":"bad"}`))
	req.Header.Set("Content-Type", "application/json")
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d", rec.Code)
	}
}
