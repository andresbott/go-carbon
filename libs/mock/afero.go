package mock

import (
	"github.com/spf13/afero"
)

func AferoSample(dirs []string, files map[string]string) afero.Fs {
	fs := afero.NewMemMapFs()
	for _, dir := range dirs {
		err := fs.Mkdir(dir, 0755)
		if err != nil {
			panic(err)
		}
	}
	for fname, content := range files {
		// write the whole body at once
		err := afero.WriteFile(fs, fname, []byte(content), 0644)
		if err != nil {
			panic(err)
		}
	}
	return fs

}
func HttpFs(fs afero.Fs) *afero.HttpFs {
	return afero.NewHttpFs(fs)
}

func AferoSampleV1() afero.Fs {
	return AferoSample(DirsV1, filesV1)
}

var DirsV1 = []string{
	"/media/photos",
	"/media/video",
	"/media/music",
	"/text/plain",
	"/text/pdf",
	"/tree/a/a_b/a_b_a",
	"/tree/a/a_c/a_c_a",
}

var filesV1 = map[string]string{
	"/text/plain/file1.txt": "file1.txt",
	"/text/plain/file2.txt": "file2.txt",
}
