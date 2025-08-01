package assets

import (
	"embed"
	"net/http"
	"path/filepath"
	"strings"
)

//go:embed *.svg
var EmbeddedSVGs embed.FS

//go:embed css js
var EmbeddedStatic embed.FS


// GetStaticFS returns the embedded static file system for use with http.FileServer
func GetStaticFS() http.FileSystem {
	return http.FS(EmbeddedStatic)
}


// ServeSVG serves embedded SVG files
func ServeSVG(w http.ResponseWriter, r *http.Request) {
	filename := strings.TrimPrefix(r.URL.Path, "/assets/")
	
	// Security check - only allow SVG files
	if !strings.HasSuffix(filename, ".svg") {
		http.NotFound(w, r)
		return
	}
	
	content, err := EmbeddedSVGs.ReadFile(filename)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	
	w.Header().Set("Content-Type", "image/svg+xml")
	w.Header().Set("Cache-Control", "public, max-age=3600")
	w.Write(content)
}

// ServeStatic serves embedded static files with proper content types
func ServeStatic(w http.ResponseWriter, r *http.Request) {
	filename := strings.TrimPrefix(r.URL.Path, "/static/")
	
	// Security check - prevent directory traversal
	if strings.Contains(filename, "..") {
		http.NotFound(w, r)
		return
	}
	
	// Read file from embedded filesystem
	content, err := EmbeddedStatic.ReadFile(filename)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	
	// Set appropriate content type based on file extension
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".css":
		w.Header().Set("Content-Type", "text/css")
	case ".js":
		w.Header().Set("Content-Type", "application/javascript")
	case ".json":
		w.Header().Set("Content-Type", "application/json")
	case ".txt":
		w.Header().Set("Content-Type", "text/plain")
	default:
		w.Header().Set("Content-Type", "application/octet-stream")
	}
	
	// Set cache headers for static assets
	w.Header().Set("Cache-Control", "public, max-age=3600")
	w.Write(content)
}