// Package webui serves the embedded React build.
package webui

import (
	"bytes"
	"embed"
	"io"
	"io/fs"
	"net/http"
	"path"
	"strings"
	"time"
)

// dist is populated by `make build` from web/dist. The committed .keep file
// lets a clean checkout compile before the frontend build runs.
//
//go:embed all:dist
var dist embed.FS

// Embedded returns the production frontend filesystem.
func Embedded() fs.FS {
	sub, err := fs.Sub(dist, "dist")
	if err != nil {
		return nil
	}
	return sub
}

// Handler serves hashed assets and falls back to index.html for client routes.
func Handler(fsys fs.FS) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if fsys == nil {
			http.Error(w, "Frontend belum dibangun. Jalankan make build.", http.StatusServiceUnavailable)
			return
		}
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			w.Header().Set("Allow", "GET, HEAD")
			http.Error(w, "Metode tidak diizinkan.", http.StatusMethodNotAllowed)
			return
		}

		requestPath := path.Clean(r.URL.Path)
		if requestPath == "/" {
			serveIndex(w, r, fsys)
			return
		}
		name := strings.TrimPrefix(requestPath, "/")
		if file, err := fsys.Open(name); err == nil {
			defer file.Close()
			if info, statErr := file.Stat(); statErr == nil && !info.IsDir() {
				setAssetCache(w, requestPath)
				serveFile(w, r, name, info.ModTime(), file)
				return
			}
		}
		if path.Ext(requestPath) != "" {
			http.NotFound(w, r)
			return
		}
		serveIndex(w, r, fsys)
	})
}

func serveIndex(w http.ResponseWriter, r *http.Request, fsys fs.FS) {
	file, err := fsys.Open("index.html")
	if err != nil {
		http.Error(w, "Frontend belum dibangun. Jalankan make build.", http.StatusServiceUnavailable)
		return
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil {
		http.Error(w, "Frontend tidak dapat dibaca.", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Cache-Control", "no-cache")
	serveFile(w, r, "index.html", info.ModTime(), file)
}

func serveFile(w http.ResponseWriter, r *http.Request, name string, modTime time.Time, file fs.File) {
	if seeker, ok := file.(io.ReadSeeker); ok {
		http.ServeContent(w, r, name, modTime, seeker)
		return
	}
	content, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Berkas tidak dapat dibaca.", http.StatusInternalServerError)
		return
	}
	http.ServeContent(w, r, name, modTime, bytes.NewReader(content))
}

func setAssetCache(w http.ResponseWriter, requestPath string) {
	if strings.HasPrefix(requestPath, "/assets/") {
		w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		return
	}
	w.Header().Set("Cache-Control", "public, max-age=604800")
}
