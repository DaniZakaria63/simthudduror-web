package webui

import (
	"io/fs"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/fstest"
)

func testFS() fstest.MapFS {
	return fstest.MapFS{
		"index.html":        &fstest.MapFile{Data: []byte("<html>maulid</html>")},
		"assets/app-123.js": &fstest.MapFile{Data: []byte("console.log('ok')")},
		"images/photo.webp": &fstest.MapFile{Data: []byte("image")},
	}
}

func TestHandlerServesIndexAndClientRoutes(t *testing.T) {
	handler := Handler(testFS())
	for _, route := range []string{"/", "/akun", "/backoffice"} {
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, route, nil))
		if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), "maulid") {
			t.Fatalf("route %s = %d %q", route, rec.Code, rec.Body.String())
		}
		if rec.Header().Get("Cache-Control") != "no-cache" {
			t.Fatalf("index cache header = %q", rec.Header().Get("Cache-Control"))
		}
	}
}

func TestHandlerServesAssetsWithImmutableCache(t *testing.T) {
	handler := Handler(testFS())
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/assets/app-123.js", nil))
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), "console.log") {
		t.Fatalf("asset = %d %q", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Header().Get("Cache-Control"), "immutable") {
		t.Fatalf("asset cache header = %q", rec.Header().Get("Cache-Control"))
	}

	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/images/photo.webp", nil))
	if rec.Code != http.StatusOK || rec.Header().Get("Cache-Control") != "public, max-age=604800" {
		t.Fatalf("image = %d cache=%q", rec.Code, rec.Header().Get("Cache-Control"))
	}
}

func TestHandlerRejectsMissingAssetAndMethod(t *testing.T) {
	handler := Handler(testFS())
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/assets/missing.js", nil))
	if rec.Code != http.StatusNotFound {
		t.Fatalf("missing asset status = %d", rec.Code)
	}

	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/", nil))
	if rec.Code != http.StatusMethodNotAllowed || rec.Header().Get("Allow") == "" {
		t.Fatalf("method status = %d", rec.Code)
	}
}

func TestHandlerReportsMissingBuild(t *testing.T) {
	handler := Handler(nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))
	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("missing build status = %d", rec.Code)
	}

	handler = Handler(fstest.MapFS{})
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))
	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("empty build status = %d", rec.Code)
	}
}

func TestEmbeddedReturnsFS(t *testing.T) {
	var fsys fs.FS = Embedded()
	if fsys == nil {
		t.Fatal("embedded filesystem is nil")
	}
}
