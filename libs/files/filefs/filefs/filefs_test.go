package filefs

import (
	file "git.andresbott.com/Golang/carbon/libs/files/filefs"
	"github.com/google/go-cmp/cmp"
	"path/filepath"
	"runtime"
	"testing"
)

var sampleDataRoot = []string{"D- a", "file1.txt"}

func TestList(t *testing.T) {
	tcs := []struct {
		name      string
		root      string
		path      string
		expect    []string
		expectErr string
	}{
		{
			name:   "default root to /",
			root:   "",
			path:   sampleDataDir(""),
			expect: sampleDataRoot,
		},
		{
			name:   "default jail to sampledata",
			root:   sampleDataDir(""),
			path:   "/",
			expect: sampleDataRoot,
		},
		{
			name:   "try to escape the jail to sampledata",
			root:   sampleDataDir(""),
			path:   "/../",
			expect: sampleDataRoot,
		},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			ffs, err := New(tc.root)
			if err != nil {
				t.Fatalf("unable to create file fs instance: %v", err)
			}

			got, err := ffs.List(tc.path)
			if err != nil {
				t.Fatalf("unable to list files: %v", err)
			}
			if diff := cmp.Diff(fsListToStr(got), tc.expect); diff != "" {
				t.Errorf("unexpected value (-got +want)\n%s", diff)
			}
		})
	}
}

// helper function that takes a list of FS.entry and returns a slice of the file names
func fsListToStr(in []file.FSEntry) []string {
	out := []string{}
	for _, f := range in {
		if f.IsDir() {
			out = append(out, "D- "+f.Name())
		} else {
			out = append(out, f.Name())
		}
	}
	return out
}

// helper function to get the absolute path of the directory containing the sample data
func sampleDataDir(path string) string {
	_, filename, _, _ := runtime.Caller(0)
	d := filepath.Dir(filename)
	d2 := filepath.Join(d, "/sampledata/"+path)
	return d2
}
