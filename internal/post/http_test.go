package post

import (
	"bytes"
	"context"
	"database/sql"
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

type validationAdapter struct {
	repo account.Repository
}

func (a validationAdapter) ValidationStatus(ctx context.Context, accountID string) (string, error) {
	status, err := a.repo.ValidationStatus(ctx, accountID)
	return string(status), err
}

type apiServer struct {
	t      *testing.T
	server *httptest.Server
	client *http.Client
	db     *sql.DB
}

func newAPIServer(t *testing.T) *apiServer {
	t.Helper()
	db := newTestDB(t)
	accountRepo := account.NewRepository(db)
	accountService := account.NewService(accountRepo, time.Hour)
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	postRepo := NewRepository(db)
	postService := NewService(postRepo, validationAdapter{repo: accountRepo})
	postService.now = func() time.Time { return time.Date(2026, 11, 14, 9, 0, 0, 0, time.UTC) }

	mux := http.NewServeMux()
	require := authn.Require(accountService)
	optional := authn.Optional(accountService)
	account.NewHandler(accountService, ratelimit.New(1000, time.Minute), logger, false, nil).RegisterRoutes(mux, require, optional)
	NewHandler(postService, logger).RegisterRoutes(mux, require, optional)
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)
	return &apiServer{t: t, server: server, client: server.Client(), db: db}
}

func (s *apiServer) do(method, path string, body any, cookie *http.Cookie, csrf string, bearer string) *http.Response {
	s.t.Helper()
	var reader io.Reader
	if body != nil {
		raw, err := json.Marshal(body)
		if err != nil {
			s.t.Fatalf("marshal: %v", err)
		}
		reader = bytes.NewReader(raw)
	}
	req, err := http.NewRequest(method, s.server.URL+path, reader)
	if err != nil {
		s.t.Fatalf("request: %v", err)
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
	response, err := s.client.Do(req)
	if err != nil {
		s.t.Fatalf("do: %v", err)
	}
	return response
}

func (s *apiServer) register(email, name, password string) (*http.Cookie, string, string) {
	s.t.Helper()
	response := s.do(http.MethodPost, "/api/v1/auth/register", map[string]string{
		"email": email, "name": name, "password": password,
	}, nil, "", "")
	if response.StatusCode != http.StatusCreated {
		s.t.Fatalf("register status = %d", response.StatusCode)
	}
	defer response.Body.Close()
	var payload struct {
		CSRF    string       `json:"csrf_token"`
		Account account.View `json:"account"`
	}
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		s.t.Fatalf("decode register: %v", err)
	}
	for _, cookie := range response.Cookies() {
		if cookie.Name == authn.CookieName {
			return cookie, payload.CSRF, payload.Account.ID
		}
	}
	s.t.Fatal("missing session cookie")
	return nil, "", ""
}

func (s *apiServer) validate(accountID string) {
	s.t.Helper()
	if _, err := s.db.Exec(`UPDATE accounts SET validation_status = 'validated' WHERE id = ?`, accountID); err != nil {
		s.t.Fatalf("validate account: %v", err)
	}
}

func TestHTTPPostsAuthorizationAndInteraction(t *testing.T) {
	server := newAPIServer(t)
	cookie, csrf, accountID := server.register("member@example.com", "Member", "kata-sandi-aman")

	response := server.do(http.MethodGet, "/api/v1/posts", nil, nil, "", "")
	if response.StatusCode != http.StatusOK {
		t.Fatalf("anonymous list status = %d", response.StatusCode)
	}
	if body := readJSON(t, response); len(body["posts"].([]any)) != 0 {
		t.Fatal("new feed should be empty")
	}

	response = server.do(http.MethodPost, "/api/v1/posts", map[string]string{"body": "Kabar komunitas yang valid."}, cookie, "", "")
	if response.StatusCode != http.StatusForbidden {
		t.Fatalf("cookie write without CSRF status = %d", response.StatusCode)
	}
	response.Body.Close()

	response = server.do(http.MethodPost, "/api/v1/posts", map[string]string{"body": "Kabar komunitas yang valid."}, cookie, csrf, "")
	if response.StatusCode != http.StatusForbidden {
		t.Fatalf("unvalidated create status = %d", response.StatusCode)
	}
	if body := readJSON(t, response); body["error"].(map[string]any)["code"] != "validation_required" {
		t.Fatalf("unexpected validation error: %+v", body)
	}

	server.validate(accountID)
	response = server.do(http.MethodPost, "/api/v1/posts", map[string]string{"body": "Kabar komunitas yang valid."}, cookie, csrf, "")
	if response.StatusCode != http.StatusCreated {
		t.Fatalf("validated create status = %d body=%s", response.StatusCode, readRaw(response))
	}
	created := readJSON(t, response)["post"].(map[string]any)
	postID := created["id"].(string)

	response = server.do(http.MethodPost, "/api/v1/posts", map[string]any{
		"body": "Kabar dengan klaim palsu.", "role": "staff", "validation_status": "validated",
	}, cookie, csrf, "")
	if response.StatusCode != http.StatusBadRequest {
		t.Fatalf("forged privilege fields status = %d", response.StatusCode)
	}
	response.Body.Close()

	response = server.do(http.MethodGet, "/api/v1/posts?q=komunitas", nil, cookie, "", "")
	body := readJSON(t, response)
	if len(body["posts"].([]any)) != 1 {
		t.Fatalf("filter failed: %+v", body)
	}

	response = server.do(http.MethodPut, "/api/v1/posts/"+postID+"/like", nil, cookie, csrf, "")
	if response.StatusCode != http.StatusOK {
		t.Fatalf("like status = %d", response.StatusCode)
	}
	likeBody := readJSON(t, response)["post"].(map[string]any)
	if likeBody["like_count"].(float64) != 1 || likeBody["liked_by_me"] != true {
		t.Fatalf("unexpected like: %+v", likeBody)
	}
	response = server.do(http.MethodPut, "/api/v1/posts/"+postID+"/like", nil, cookie, csrf, "")
	if body := readJSON(t, response); body["post"].(map[string]any)["like_count"].(float64) != 1 {
		t.Fatalf("like not idempotent: %+v", body)
	}
	response = server.do(http.MethodDelete, "/api/v1/posts/"+postID+"/like", nil, cookie, csrf, "")
	if body := readJSON(t, response); body["post"].(map[string]any)["like_count"].(float64) != 0 {
		t.Fatalf("unlike failed: %+v", body)
	}

	response = server.do(http.MethodPut, "/api/v1/posts/missing/like", nil, cookie, csrf, "")
	if response.StatusCode != http.StatusNotFound {
		t.Fatalf("missing post like status = %d", response.StatusCode)
	}
	response.Body.Close()

	response = server.do(http.MethodPut, "/api/v1/posts/"+postID+"/like", nil, nil, "", "")
	if response.StatusCode != http.StatusUnauthorized {
		t.Fatalf("anonymous like status = %d", response.StatusCode)
	}
	response.Body.Close()
}

func TestHTTPMobileBearerFlow(t *testing.T) {
	server := newAPIServer(t)
	cookie, csrf, accountID := server.register("mobilepost@example.com", "Mobile", "kata-sandi-aman")
	server.validate(accountID)
	_ = cookie
	_ = csrf

	response := server.do(http.MethodPost, "/api/v1/auth/login", map[string]string{
		"email": "mobilepost@example.com", "password": "kata-sandi-aman",
	}, nil, "", "")
	if response.StatusCode != http.StatusOK {
		t.Fatalf("login status = %d", response.StatusCode)
	}
	response.Body.Close()

	request, _ := http.NewRequest(http.MethodPost, server.server.URL+"/api/v1/auth/login", bytes.NewReader([]byte(
		`{"email":"mobilepost@example.com","password":"kata-sandi-aman"}`)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Maulid-Client", "mobile")
	response, err := server.client.Do(request)
	if err != nil {
		t.Fatalf("mobile login: %v", err)
	}
	var login struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(response.Body).Decode(&login); err != nil || login.Token == "" {
		t.Fatalf("mobile token missing: %v", err)
	}
	response.Body.Close()

	response = server.do(http.MethodPost, "/api/v1/posts", map[string]string{"body": "Kabar dari aplikasi mobile."}, nil, "", login.Token)
	if response.StatusCode != http.StatusCreated {
		t.Fatalf("bearer create status = %d body=%s", response.StatusCode, readRaw(response))
	}
	created := readJSON(t, response)["post"].(map[string]any)
	postID := created["id"].(string)

	// Bearer writes need no CSRF token.
	response = server.do(http.MethodPut, "/api/v1/posts/"+postID+"/like", nil, nil, "", login.Token)
	if response.StatusCode != http.StatusOK {
		t.Fatalf("bearer like status = %d", response.StatusCode)
	}
	response.Body.Close()
}

func TestHTTPBadLimitAndOversizedJSON(t *testing.T) {
	server := newAPIServer(t)
	response := server.do(http.MethodGet, "/api/v1/posts?limit=0", nil, nil, "", "")
	if response.StatusCode != http.StatusBadRequest {
		t.Fatalf("invalid limit status = %d", response.StatusCode)
	}
	response.Body.Close()

	cookie, csrf, _ := server.register("big@example.com", "Big", "kata-sandi-aman")
	big := bytes.Repeat([]byte("a"), 40<<10)
	req, _ := http.NewRequest(http.MethodPost, server.server.URL+"/api/v1/posts",
		bytes.NewReader(append(append([]byte(`{"body":"`), big...), []byte(`"}`)...)))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(cookie)
	req.Header.Set("X-CSRF-Token", csrf)
	response, err := server.client.Do(req)
	if err != nil {
		t.Fatalf("oversized post: %v", err)
	}
	if response.StatusCode != http.StatusRequestEntityTooLarge {
		t.Fatalf("oversized post status = %d", response.StatusCode)
	}
	response.Body.Close()
}

func readJSON(t *testing.T, response *http.Response) map[string]any {
	t.Helper()
	defer response.Body.Close()
	var payload map[string]any
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	return payload
}

func readRaw(response *http.Response) string {
	defer response.Body.Close()
	raw, _ := io.ReadAll(response.Body)
	return string(raw)
}
