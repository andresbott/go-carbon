package tmpl

import (
	"bytes"
	"embed"
	"github.com/davecgh/go-spew/spew"
	"testing"
)

//go:embed sampledata/*.html
var content embed.FS

func TestName(t *testing.T) {

	tp, err := newTmpl(content, "sampledata")
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}

	buf := new(bytes.Buffer)
	err = tp.Write("file1.html", buf)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}

	spew.Dump(buf.String())
}
