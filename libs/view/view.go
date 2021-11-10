package view

import (
	"html/template"
	"io"
	"io/fs"
	"net/http"
)

type Cfg struct {
	Prefix    string // the prefix to the views, default to view
	TmplDir   string
	StaticDir string
}

type view struct {
	staticFs   fs.FS
	templateFs fs.FS
	templates  map[string]*template.Template
}

func NewView(in fs.FS, cfg Cfg) (*view, error) {

	v := view{
		templates: map[string]*template.Template{},
	}

	// view prefix
	prefix := "view"
	if cfg.Prefix != "" {
		prefix = cfg.Prefix
	}

	subDir, err := fs.Sub(in, prefix)
	if err != nil {
		return nil, err
	}

	// tmplDir
	tmpl := "tmpl"
	if cfg.TmplDir != "" {
		tmpl = cfg.TmplDir
	}
	v.templateFs, err = fs.Sub(subDir, tmpl)
	if err != nil {
		return nil, err
	}

	// staticDir
	static := "static"
	if cfg.StaticDir != "" {
		static = cfg.StaticDir
	}
	v.staticFs, err = fs.Sub(subDir, static)

	return &v, nil
}

// Execute Will parse a single template and execute it
// parsed templates are cached in a map to be reused
func (v view) Execute(file string, data interface{}, out io.Writer) error {

	// cache templates in an internal map
	// we assume we wont have size problems
	if _, ok := v.templates[file]; !ok {
		t, err := template.ParseFS(v.templateFs, file)
		if err != nil {
			return err
		}
		v.templates[file] = t
	}

	//standard output to print merged data
	err := v.templates[file].Execute(out, data)
	if err != nil {
		return err
	}
	return nil
}

//Read will read the contents of a static file and write it into a io.writer
func (v view) Read(file string, out io.Writer) error {
	fHandle, err := v.staticFs.Open(file)
	if err != nil {
		return err
	}

	b := make([]byte, 8)
	for {
		n, err := fHandle.Read(b)
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		_, err = out.Write(b[:n])
		if err != nil {
			return err
		}
	}
}

// StaticFsHandler returns a http handler responsible to serve the static content
func (v view) StaticFsHandler() http.Handler {
	return http.FileServer(http.FS(v.staticFs))
}
