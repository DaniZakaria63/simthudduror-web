package httpx

import (
	"bytes"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDecodeJSONAcceptsValidObject(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"email":"a@example.com"}`))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	var payload struct {
		Email string `json:"email"`
	}
	rec := httptest.NewRecorder()
	if err := DecodeJSON(rec, req, &payload); err != nil {
		t.Fatalf("DecodeJSON: %v", err)
	}
	if payload.Email != "a@example.com" {
		t.Fatalf("unexpected payload: %+v", payload)
	}
}

func TestDecodeJSONRejectsBadInput(t *testing.T) {
	cases := []struct {
		name        string
		contentType string
		body        string
		status      int
	}{
		{"wrong content type", "text/plain", `{}`, http.StatusUnsupportedMediaType},
		{"malformed", "application/json", `{`, http.StatusBadRequest},
		{"unknown field", "application/json", `{"unknown":true,"email":"a"}`, http.StatusBadRequest},
		{"trailing data", "application/json", `{"email":"a"} {"email":"b"}`, http.StatusBadRequest},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tc.body))
			req.Header.Set("Content-Type", tc.contentType)
			var payload struct {
				Email string `json:"email"`
			}
			err := DecodeJSON(httptest.NewRecorder(), req, &payload)
			var problem *Problem
			if !errors.As(err, &problem) || problem.Status != tc.status {
				t.Fatalf("unexpected error: %#v", err)
			}
		})
	}
}

func TestDecodeJSONRejectsOversizedBody(t *testing.T) {
	body := `{"email":"` + strings.Repeat("a", MaxBodyBytes+10) + `"}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	var payload struct {
		Email string `json:"email"`
	}
	err := DecodeJSON(httptest.NewRecorder(), req, &payload)
	var problem *Problem
	if !errors.As(err, &problem) || problem.Status != http.StatusRequestEntityTooLarge {
		t.Fatalf("unexpected error: %#v", err)
	}
}

func TestHandleProblem(t *testing.T) {
	rec := httptest.NewRecorder()
	HandleProblem(rec, &Problem{Status: http.StatusForbidden, Code: "forbidden", Message: "no"})
	if rec.Code != http.StatusForbidden || !strings.Contains(rec.Body.String(), "forbidden") {
		t.Fatalf("unexpected response: %d %s", rec.Code, rec.Body.String())
	}
	rec = httptest.NewRecorder()
	HandleProblem(rec, errors.New("secret internal detail"))
	if rec.Code != http.StatusInternalServerError || strings.Contains(rec.Body.String(), "secret") {
		t.Fatalf("internal detail leaked: %s", rec.Body.String())
	}
}

func TestSameOrigin(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "https://maulid.id/api", nil)
	req.Host = "maulid.id"
	if !SameOrigin(req, nil, true) {
		t.Fatal("missing Origin should be allowed for non-browser clients")
	}
	req.Header.Set("Origin", "https://maulid.id")
	if !SameOrigin(req, nil, true) {
		t.Fatal("same-origin request rejected")
	}
	if SameOrigin(req, nil, false) == false {
		t.Fatal("development same-origin request rejected")
	}
	req.Header.Set("Origin", "https://evil.example")
	if SameOrigin(req, nil, true) {
		t.Fatal("cross-origin request accepted")
	}
	if !SameOrigin(req, []string{"https://evil.example/"}, true) {
		t.Fatal("trusted origin rejected")
	}
	req.Header.Set("Origin", "http://maulid.id")
	if SameOrigin(req, nil, true) {
		t.Fatal("insecure scheme accepted in a secure deployment")
	}
	req.Header.Set("Origin", "not a url")
	if SameOrigin(req, nil, false) {
		t.Fatal("malformed origin accepted")
	}
}

func TestClientIP(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "192.0.2.10:443"
	if got := ClientIP(req); got != "192.0.2.10" {
		t.Fatalf("ClientIP = %q", got)
	}
	req.RemoteAddr = "malformed"
	if got := ClientIP(req); got != "malformed" {
		t.Fatalf("ClientIP fallback = %q", got)
	}
}

func TestSecurityHeaders(t *testing.T) {
	handler := SecurityHeaders(true)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNoContent) }))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/v1/posts", nil))
	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d", rec.Code)
	}
	if rec.Header().Get("Content-Security-Policy") == "" || rec.Header().Get("Strict-Transport-Security") == "" {
		t.Fatal("security headers missing")
	}
	if rec.Header().Get("Cache-Control") != "no-store" {
		t.Fatal("API response should not be cached")
	}
}

func TestRecovererHidesPanic(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(bytes.NewBuffer(nil), nil))
	handler := Recoverer(logger)(http.HandlerFunc(func(http.ResponseWriter, *http.Request) { panic("boom") }))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))
	if rec.Code != http.StatusInternalServerError || strings.Contains(rec.Body.String(), "boom") {
		t.Fatalf("panic leaked: %d %s", rec.Code, rec.Body.String())
	}
}

func TestAccessLog(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, nil))
	handler := AccessLog(logger)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	}))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/healthz", nil))
	if !strings.Contains(buf.String(), "status=418") || !strings.Contains(buf.String(), "path=/healthz") {
		t.Fatalf("access log missing fields: %s", buf.String())
	}
}

func TestJSONWritesNoStore(t *testing.T) {
	rec := httptest.NewRecorder()
	JSON(rec, http.StatusOK, map[string]string{"ok": "yes"})
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), `"ok":"yes"`) {
		t.Fatalf("unexpected response: %d %s", rec.Code, rec.Body.String())
	}
	if rec.Header().Get("Cache-Control") != "no-store" {
		t.Fatal("missing no-store")
	}
}
