package web

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

//go:embed static/*
var staticFS embed.FS

// AddUIRoute sets up static file serving for the embedded React build with fallback.
func AddUIRoute(router chi.Router) {
	embeddedBuildFolder := newStaticFileSystem()
	fallback := newFallbackFileSystem(embeddedBuildFolder)

	// Serve static files and fallback to index.html for SPA routes
	router.Handle("/*", http.StripPrefix("/", http.FileServer(fallback)))
}

// ----------------------------------------------------------------------

// staticFileSystem wraps an http.FileSystem from embedded files.
type staticFileSystem struct {
	fs http.FileSystem
}

func newStaticFileSystem() *staticFileSystem {
	sub, err := fs.Sub(staticFS, "static")
	if err != nil {
		panic(err)
	}

	return &staticFileSystem{
		fs: http.FS(sub),
	}
}

func (s *staticFileSystem) Open(name string) (http.File, error) {
	buildpath := "static" + name

	// support for folders
	if strings.HasSuffix(name, "/") {
		_, err := staticFS.ReadDir(strings.TrimSuffix(buildpath, "/"))
		if err != nil {
			return nil, err
		}

		return s.fs.Open(name)
	}

	// support for files
	f, err := staticFS.Open(buildpath)
	if err != nil {
		return nil, err
	}

	_ = f.Close()

	return s.fs.Open(name)
}

// ----------------------------------------------------------------------

// fallbackFileSystem serves /index.html for any non-file path (SPA fallback).
type fallbackFileSystem struct {
	staticFS *staticFileSystem
}

func newFallbackFileSystem(fs *staticFileSystem) *fallbackFileSystem {
	return &fallbackFileSystem{staticFS: fs}
}

func (f *fallbackFileSystem) Open(path string) (http.File, error) {
	// Attempt to open the actual file
	file, err := f.staticFS.fs.Open(path)
	if err == nil {
		return file, nil
	}

	// Otherwise serve index.html
	return f.staticFS.fs.Open("index.html")
}

// root.Handle("/*", http.FileServer(http.FS(staticFS)))
