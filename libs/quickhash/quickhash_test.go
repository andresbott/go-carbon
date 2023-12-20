package quickhash

import (
	"bufio"
	"encoding/base64"
	"github.com/davecgh/go-spew/spew"
	"github.com/google/go-cmp/cmp"
	"github.com/spf13/afero"
	"testing"
)

var _ = spew.Dump // keep the dependency in the test
func TestQuick(t *testing.T) {

	tcs := []struct {
		name   string
		in     string
		expect string
	}{
		{
			name:   "hello world",
			in:     fileHelloWorld,
			expect: "C97qD6t+h4IfB07Hem8KTg==",
		},
		{
			name:   "hello world: expect same result",
			in:     fileHelloWorld,
			expect: "C97qD6t+h4IfB07Hem8KTg==",
		},
		{
			name:   "5MB emptyFile",
			in:     fileEmpty5MBFile,
			expect: "gIDAAjiWd8m1uXqfv2fgRw==",
		},
		{
			name:   "10MB emptyFile",
			in:     fileEmpty5MBFile,
			expect: "gIDAAjiWd8m1uXqfv2fgRw==",
		},
	}

	fs := populateFs()
	qh := NewQuick()
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {

			f, _ := fs.Open(tc.in)
			hash, err := qh.SumFile(f)
			if err != nil {
				t.Errorf("unexpected err %v", err)
			}
			gotB64 := base64.StdEncoding.EncodeToString(hash[:])
			if diff := cmp.Diff(gotB64, tc.expect); diff != "" {
				t.Errorf("unexpected value (-got +want)\n%s", diff)
			}

		})
	}

}

const fileHelloWorld = "hello-world.txt"
const fileEmpty5MBFile = "empty5MB"
const fileEmpty10MBFile = "empty10MB"

func populateFs() afero.Fs {
	memFs := afero.NewMemMapFs()

	f, _ := memFs.Create(fileHelloWorld)
	_, _ = f.Write([]byte("hello World"))

	err := emptyFile(memFs, fileEmpty5MBFile, 5)
	if err != nil {
		panic(err)
	}

	err = emptyFile(memFs, fileEmpty10MBFile, 10)
	if err != nil {
		panic(err)
	}

	return memFs
}

func emptyFile(fs afero.Fs, fname string, sizeMeg int) error {
	// https://go.dev/play/p/sbQ0u35PGx1
	f, err := fs.Create(fname)
	if err != nil {
		return err
	}
	defer f.Close()

	buf := make([]byte, 4096)
	//buf[len(buf)-1] = '\n'
	w := bufio.NewWriterSize(f, len(buf))

	var size int64
	size = 1024 * 1024 * int64(sizeMeg)

	for i := int64(0); i < size; i += int64(len(buf)) {
		_, err := w.Write(buf)
		if err != nil {
			return err
		}
	}
	err = w.Flush()
	if err != nil {
		return err
	}
	err = f.Sync()
	if err != nil {
		return err
	}

	f.Close()
	return nil
}
