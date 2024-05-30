package spa

import (
	"embed"
	"github.com/davecgh/go-spew/spew"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

//go:embed files/ui/*
var UiFiles embed.FS

const (
	UIDir     = "files/ui"
	indexFile = "index.html"
)

type SpaHandler struct {
	StaticFS    embed.FS
	StaticPath  string
	IndexPath   string
	StripPrefix string
}

func (h SpaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// TODO add path calculation for all of the permutations

	// get the absolute path to prevent directory traversal
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		// if we failed to get the absolute path respond with a 400 bad request and stop
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// strip the path prefix from the requests and then
	// prepend the request path with the path to the static directory
	path = filepath.Join(h.StaticPath, strings.TrimPrefix(path, h.StripPrefix))

	spew.Dump(path)

	fi, err := h.StaticFS.Open(path)
	fstat, _ := fi.Stat()

	if os.IsNotExist(err) || fstat.IsDir() {
		// file does not exist, serve the index html file
		index, err := h.StaticFS.ReadFile(filepath.Join(h.StaticPath, indexFile))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		spew.Dump(filepath.Join(h.StaticPath, indexFile))
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusAccepted)
		if _, err := w.Write(index); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	} else if err != nil {
		spew.Dump("err")
		// return 500 if it's not an error because the file is not found
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// get the subdirectory of the static dir and simply serve the file
	statics, err := fs.Sub(h.StaticFS, h.StaticPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.FileServer(http.FS(statics)).ServeHTTP(w, r)
	//StripAndServeFS(h.StripPrefix, statics).ServeHTTP(w, r)
}
func StripAndServeFS(strip string, fs fs.FS) http.Handler {

	return http.StripPrefix(strip, http.FileServer(http.FS(fs)))
}
