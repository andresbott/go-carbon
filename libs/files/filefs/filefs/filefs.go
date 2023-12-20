package filefs

import (
	file "git.andresbott.com/Golang/carbon/libs/files/filefs"
	"os"
	"path/filepath"
	"strings"
)

type FileFs struct {
	root string
}

func New(root string) (*FileFs, error) {

	if root == "" {
		root = "/"
	}

	ffs := FileFs{
		root: root,
	}

	return &ffs, nil

}

// TODO:  filter by type, order, dirs first, deep?, limit, page?
func (f FileFs) List(path string) ([]file.FSEntry, error) {
	if path == "" {
		path = "/"
	}

	rel, _ := filepath.Rel(f.root, filepath.Join(f.root, path))
	if strings.HasPrefix(rel, "..") {
		// prevent escaping the jail
		path = "/"
	}

	abs, _ := filepath.Abs(filepath.Join(f.root, path))

	// read FS
	files, err := os.ReadDir(abs)
	if err != nil {
		// Todo explicit dir does not exists error
		return nil, err
	}

	// generate return
	ret := make([]file.FSEntry, len(files))
	for i := 0; i < len(files); i++ {
		ret[i] = files[i]
	}
	return ret, nil
}
