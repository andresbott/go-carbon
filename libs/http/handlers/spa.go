package handlers

import (
	_ "embed"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Spa is mainly inspired on the gorilla spa handler documentation
// https://github.com/gorilla/mux?tab=readme-ov-file#serving-single-page-applications
type Spa struct {
	staticPath string
	indexPath  string
}

// ServeHTTP inspects the URL path to locate a file within the static dir
// on the SPA handler. If a file is found, it will be served. If not, the
// file located at the index path on the SPA handler will be served. This
// is suitable behavior for serving an SPA (single page application).
func (h Spa) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Join internally call path.Clean to prevent directory traversal
	path := filepath.Join(h.staticPath, r.URL.Path)

	// check whether a file exists or is a directory at the given path
	fi, err := os.Stat(path)
	if os.IsNotExist(err) || fi.IsDir() {
		// file does not exist or path is a directory, serve index.html
		http.ServeFile(w, r, filepath.Join(h.staticPath, h.indexPath))
		return
	}

	if err != nil {
		// if we got an error (that wasn't that the file doesn't exist) stating the
		// file, return a 500 internal server error and stop
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// otherwise, use http.FileServer to serve the static file
	http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
}

// SpaHandler is a http handler capable of serving SPAs from an fs.FS ( tested are os.DirFS and embed.FS)
// configuration:
// FsSubDir allows to keep more files that only the SPA in an FS and serve the data from a sub dir
// Notice that the dir path needs to be relative and cannot be ./ or ../; empty string will be replaced by "."

func NewSpaHAndler(inputFs fs.FS, fsSubDir, pathPrefix string) (SpaHandler, error) {
	if inputFs == nil {
		return SpaHandler{}, fmt.Errorf("fs cannot be nil")
	}
	if fsSubDir != "" {
		newFs, err := fs.Sub(inputFs, fsSubDir)
		if err != nil {
			return SpaHandler{}, err
		}
		inputFs = newFs
	}

	s := SpaHandler{
		fs:         inputFs,
		fsSubDir:   fsSubDir,
		pathPrefix: pathPrefix,
	}
	return s, nil
}

func MustSpaHandler(fs fs.FS, fsSubDir, pathPrefix string) SpaHandler {
	h, err := NewSpaHAndler(fs, fsSubDir, pathPrefix)
	if err != nil {
		panic(err)
	}
	return h
}

type SpaHandler struct {
	fs         fs.FS
	fsSubDir   string // if the SPA is in a subdirectory e.g. static/dist
	pathPrefix string // if the SPA is served with a path prefix, e.g. "ui" in  http://my-app.com/ui/
	indexFile  string // main index file of the spa e.g. index.html
}

const index = "index.html"

func (h SpaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	reqPath := strings.TrimPrefix(r.URL.Path, h.pathPrefix)
	if reqPath == "" {
		reqPath = "/"
	}

	// serve index on root
	if reqPath == "/" {
		path := filepath.Join(reqPath, index)
		isDir, err := checkDir(path, h.fs)
		if os.IsNotExist(err) {
			http.Error(w, "404 page not found", http.StatusNotFound)
		} else if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !isDir {
			http.StripPrefix(h.pathPrefix, http.FileServerFS(h.fs)).ServeHTTP(w, r)
		}
		return
	}

	// disable directory listing
	if strings.HasSuffix(reqPath, "/") {
		http.NotFound(w, r)
		return
	}

	// don't redirect an existing folder e.g. if r.URL.Path = /assets/ui
	// instead of redirecting to /assets/ui/ and then returning a 404
	// immediately check if it is a folder a return 404
	isDir, err := checkDir(reqPath, h.fs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if isDir {
		http.Error(w, "404 page not found", http.StatusNotFound)
		return
	}
	http.StripPrefix(h.pathPrefix, http.FileServerFS(h.fs)).ServeHTTP(w, r)
}

func checkDir(path string, fs fs.FS) (bool, error) {
	path = strings.TrimPrefix(path, "/")
	f, err := fs.Open(path)
	if err != nil {
		return false, err
	}
	stat, err := f.Stat()
	if err != nil {
		return false, err
	}
	return stat.IsDir(), nil
}
