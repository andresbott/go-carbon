package view_test

import (
	"bytes"
	"embed"
	"fmt"
	"git.andresbott.com/Golang/carbon/libs/view"
	"net/http"
	"net/http/httptest"
)

//go:embed sampledata/*
var sampleContent embed.FS

func Example() {

	// we create a new instance of view ¡
	v, _ := view.NewView(sampleContent, view.Cfg{
		Prefix:    "sampledata",
		TmplDir:   "tmpl",
		StaticDir: "static",
	})

	// io.wirter to write the content into
	buf := new(bytes.Buffer)

	// data to fill out the template
	type data struct {
		Name string
	}
	// execute the template called file2.html in the tmpl folder
	_ = v.Execute("file2.html", data{Name: "Batman"}, buf)
	// print the string
	fmt.Println(buf.String())

	// printing a static file
	staticFileBuf := new(bytes.Buffer)
	// execute the static file called static-file.html in the tmpl folder
	err := v.Read("static-file.html", staticFileBuf)
	if err != nil {
		fmt.Printf("%s\n", err)
	}
	// print the string
	fmt.Println(staticFileBuf.String())

	// Output:
	// <html>Batman</html>
	// <html>file1</html>

}

func ExampleServe() {

	// we create a new instance of view ¡
	v, _ := view.NewView(sampleContent, view.Cfg{
		Prefix:    "sampledata",
		TmplDir:   "tmpl",
		StaticDir: "static",
	})

	// get an http handler to server the static content
	staticHandler := v.StaticFsHandler()

	req := httptest.NewRequest(http.MethodGet, "/static-file.html", nil)
	res := httptest.NewRecorder()

	staticHandler.ServeHTTP(res, req)

	// print the response of a request
	fmt.Println(res.Body)

	// Output:
	// <html>file1</html>

}
