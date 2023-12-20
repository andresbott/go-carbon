package view_test

// TODO: test for static file index, should not print index
// todo: test for content type headers

//
//import (
//	"bytes"
//	"embed"
//	"github.com/google/go-cmp/cmp"
//	"testing"
//)
//
////go:embed sampledata/*.html
//var content embed.FS
//
//func TestRead(t *testing.T) {
//
//	tp, err := newTmpl(content, "sampledata")
//	if err != nil {
//		t.Errorf("unexpected error %v", err)
//	}
//
//	t.Run("static content is loaded", func(t *testing.T) {
//		buf := new(bytes.Buffer)
//		err = tp.Read("static-file.html", buf)
//		if err != nil {
//			t.Errorf("unexpected error %v", err)
//		}
//		expected := "<html>file1</html>"
//
//		if diff := cmp.Diff(expected, buf.String()); diff != "" {
//			t.Errorf("(-got +want) %s", diff)
//		}
//	})
//
//	t.Run("error when file does not exist", func(t *testing.T) {
//		buf := new(bytes.Buffer)
//		err = tp.Read("no-no.html", buf)
//		if err == nil {
//			t.Fatalf("expecting an error but nore returned")
//		}
//
//		expectedErrMsg := "open no-no.html: file does not exist"
//		if err != nil && err.Error() != expectedErrMsg {
//			t.Fatalf("expecting error message: \"%s\", but got: \"%s\"", expectedErrMsg, err.Error())
//		}
//	})
//
//	//t.Run("test parsing a file content", func(t *testing.T) {
//	//	buf := new(bytes.Buffer)
//	//	err = tp.Execute("file1.html", nil,buf)
//	//	if err != nil {
//	//		t.Errorf("unexpected error %v", err)
//	//	}
//	//
//	//	spew.Dump(buf.String())
//	//})
//	//
//	//t.Run("test non existing file", func(t *testing.T) {
//	//	buf := new(bytes.Buffer)
//	//	err = tp.Execute("file1.html", nil,buf)
//	//	if err != nil {
//	//		t.Errorf("unexpected error %v", err)
//	//	}
//	//
//	//	spew.Dump(buf.String())
//	//})
//
//
//}
