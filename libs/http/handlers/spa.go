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

	// get the absolute path to prevent directory traversal
	//path, err := filepath.Abs(r.URL.Path)
	//if err != nil {
	//	http.Error(w, err.Error(), http.StatusBadRequest)
	//	return
	//}

	// strip the path prefix from the requests and prepend with the path of the static directory
	//path = filepath.Join(h.fsSubDir, strings.TrimPrefix(path, h.pathPrefix))

	//spew.Dump(path)
	//
	//fi, err := h.fs.Open(path)
	//_ = fi
	//spew.Dump(err)

	//if os.IsNotExist(err) {
	//	spew.Dump("HERE")
	//	spew.Dump(err)
	//	spew.Dump(fi)
	//	//fstat, _ := fi.Stat()
	//	//|| fstat.IsDir()
	//	// file does not exist, serve the index html file
	//	//index, err := h.StaticFS.ReadFile(filepath.Join(h.StaticPath, indexFile))
	//	//if err != nil {
	//	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//	//	return
	//	//}
	//	index := []byte("ddd")
	//
	//	spew.Dump(filepath.Join(h.FsSubDir, h.IndexFile))
	//	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	//	w.WriteHeader(http.StatusAccepted)
	//	if _, err := w.Write(index); err != nil {
	//		http.Error(w, err.Error(), http.StatusInternalServerError)
	//	}
	//	return
	//} else if err != nil {
	//	spew.Dump("err")
	//	// return 500 if it's not an error because the file is not found
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//	return
	//}

	// get the subdirectory of the static dir and simply serve the file

	//statics, err := fs.Sub(h.fs, h.fsSubDir)
	//if err != nil {
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//	return
	//}
	//
	//spew.Dump(statics)
	//
	//spew.Dump(r.URL)
	//spew.Dump(h.fs)

	// serve index on root
	if r.URL.Path == "/" {
		path := filepath.Join(r.URL.Path, index)
		isDir, err := checkDir(path, h.fs)
		if os.IsNotExist(err) {
			http.Error(w, "404 page not found", http.StatusNotFound)
		} else if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !isDir {
			http.FileServerFS(h.fs).ServeHTTP(w, r)
		}
		return
	}

	// disable directory listing
	if strings.HasSuffix(r.URL.Path, "/") {
		http.NotFound(w, r)
		return
	}

	// don't redirect an existing folder e.g. if r.URL.Path = /assets/ui
	// instead of redirecting to /assets/ui/ and then returning a 404
	// immediately check if it is a folder a return 404
	isDir, err := checkDir(r.URL.Path, h.fs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if isDir {
		http.Error(w, "404 page not found", http.StatusNotFound)
		return
	}

	http.FileServerFS(h.fs).ServeHTTP(w, r)
	////StripAndServeFS(h.StripPrefix, statics).ServeHTTP(w, r)
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
