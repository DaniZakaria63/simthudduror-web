package authn

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type fakeAuthenticator struct {
	principal Principal
	err       error
}

func (f fakeAuthenticator) Authenticate(context.Context, string) (Principal, error) {
	return f.principal, f.err
}

func TestRequireAuthenticatesCookieAndEnforcesCSRF(t *testing.T) {
	token := "browser-token"
	csrf := "csrf-value"
	sum := sha256.Sum256([]byte(csrf))
	auth := fakeAuthenticator{principal: Principal{AccountID: "acc", SessionCSRFHash: hex.EncodeToString(sum[:])}}
	handler := Require(auth)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		principal, ok := FromContext(r.Context())
		if !ok || principal.AccountID != "acc" {
			t.Fatal("principal missing")
		}
		w.WriteHeader(http.StatusNoContent)
	}))

	t.Run("safe method", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.AddCookie(&http.Cookie{Name: CookieName, Value: token})
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		if rec.Code != http.StatusNoContent {
			t.Fatalf("status = %d", rec.Code)
		}
	})

	t.Run("cookie write without CSRF", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		req.AddCookie(&http.Cookie{Name: CookieName, Value: token})
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		if rec.Code != http.StatusForbidden || !strings.Contains(rec.Body.String(), "csrf_invalid") {
			t.Fatalf("unexpected response: %d %s", rec.Code, rec.Body.String())
		}
	})

	t.Run("cookie write with CSRF", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		req.AddCookie(&http.Cookie{Name: CookieName, Value: token})
		req.Header.Set("X-CSRF-Token", csrf)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		if rec.Code != http.StatusNoContent {
			t.Fatalf("status = %d %s", rec.Code, rec.Body.String())
		}
	})

	t.Run("bearer write skips CSRF", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		req.Header.Set("Authorization", "Bearer mobile-token")
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		if rec.Code != http.StatusNoContent {
			t.Fatalf("status = %d %s", rec.Code, rec.Body.String())
		}
	})
}

func TestRequireRejectsMissingOrInvalidCredential(t *testing.T) {
	handler := Require(fakeAuthenticator{err: ErrInvalidCredential})(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("missing credential status = %d", rec.Code)
	}
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer bad")
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("invalid credential status = %d", rec.Code)
	}
}

func TestOptionalDoesNotRejectAnonymous(t *testing.T) {
	handler := Optional(fakeAuthenticator{err: ErrInvalidCredential})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := FromContext(r.Context()); ok {
			t.Fatal("anonymous request received a principal")
		}
		w.WriteHeader(http.StatusNoContent)
	}))

	t.Run("no credential", func(t *testing.T) {
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))
		if rec.Code != http.StatusNoContent || rec.Body.Len() != 0 {
			t.Fatalf("unexpected anonymous response: %d %q", rec.Code, rec.Body.String())
		}
	})

	t.Run("stale cookie", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.AddCookie(&http.Cookie{Name: CookieName, Value: "expired"})
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		if rec.Code != http.StatusNoContent || rec.Body.Len() != 0 {
			t.Fatalf("stale cookie polluted public response: %d %q", rec.Code, rec.Body.String())
		}
	})
}

func TestOptionalAttachesPrincipal(t *testing.T) {
	auth := fakeAuthenticator{principal: Principal{AccountID: "acc"}}
	handler := Optional(auth)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		principal, ok := FromContext(r.Context())
		if !ok || principal.AccountID != "acc" || principal.Kind != KindBearer {
			t.Fatalf("unexpected principal: %+v ok=%v", principal, ok)
		}
	}))
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer token")
	handler.ServeHTTP(httptest.NewRecorder(), req)
}

func TestPrincipalIsStaff(t *testing.T) {
	if !(Principal{Role: RoleStaff}).IsStaff() {
		t.Fatal("staff role not recognised")
	}
	if (Principal{Role: "member"}).IsStaff() {
		t.Fatal("member recognised as staff")
	}
}

func TestHashTokenIsStable(t *testing.T) {
	first := HashToken("token")
	if first != HashToken("token") || len(first) != 64 {
		t.Fatalf("unexpected hash: %s", first)
	}
}

func TestCredentialFromRequestIgnoresMalformedAuthorization(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Basic abc")
	req.AddCookie(&http.Cookie{Name: CookieName, Value: "cookie-token"})
	credential, kind := credentialFromRequest(req)
	if credential != "cookie-token" || kind != KindCookie {
		t.Fatalf("unexpected credential: %q %s", credential, kind)
	}
}

func TestRequireProblemEnvelope(t *testing.T) {
	handler := Require(fakeAuthenticator{err: errors.New("db down")})(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/", nil))
	if !strings.Contains(rec.Body.String(), "unauthenticated") {
		t.Fatalf("unexpected envelope: %s", rec.Body.String())
	}
}
