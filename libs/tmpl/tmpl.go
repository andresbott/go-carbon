package tmpl

import (
	"html/template"
	"io"
	"io/fs"
)

func newTmpl(in fs.FS, prefix string)(*tmpl,error)  {

	var err error
	// strip prefix if provided
	if prefix != ""{
		in,err = fs.Sub(in, prefix)
		if err !=nil{
			return nil,err
		}
	}
	return &tmpl{
		fs: in,
		templates: map[string]*template.Template{},
	},nil
}

type tmpl struct {
	fs fs.FS
	templates map[string]*template.Template
}

func (tp tmpl) Write(file string, out io.Writer) error  {

	// cache templates in an internal map
	// we assume we wont have size problems
	if _, ok := tp.templates[file]; !ok {
		t, err  := template.ParseFS(tp.fs,file)
		if err != nil{
			return err
		}
		tp.templates[file]=t
	}

	here => another atribute as interface to execute the template
	type bla struct {}

	//standard output to print merged data
	err := tp.templates[file].Execute(out, bla{})
	if err != nil{
		return err
	}
	return nil

}

func par()  {
	// Parsing the required html
	// file in same directory
	//t, err := template.ParseFiles("index.html")
}
