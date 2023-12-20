package file_old

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

type Manager struct {
	files     []File
	dirs      []File
	fileIndex int // holds the current index of the file

	filesCounter int
	dirsCounter  int

	filterCache map[string]bool // see filterFile() for details
	once        sync.Once

	// opts
	realAbsPath string // actual absolute path to scan
	path        string // path representation in relation to the chroot
	ChRoot      string // limit all interactions to this directory
	filter      FilterCondition
}

type FilterCondition struct {
	ShowHiddenFiles bool
	ShowExtension   []string
	ShowType        []FileType // see file type ints as constants
}

func (fm *Manager) init() {
	fm.once.Do(func() {
		fm.filterCache = map[string]bool{}
		fm.path = "/"
		fm.realAbsPath = fm.ChRoot
	})
}

func (fm *Manager) Nav(path string) error {

	fm.init()

	var err error
	// get the absolute path
	absPath := filepath.Join(fm.ChRoot, path)
	if !filepath.IsAbs(absPath) {
		absPath, err = filepath.Abs(absPath)
		if err != nil {
			return fmt.Errorf("unable to get absolute path: %s", err)
		}
	}

	// chek if it exists
	file, err := os.Stat(absPath)
	if err != nil {
		return errors.New("Path does not exist: " + absPath)
	}

	if file.Mode().IsDir() {
		fm.realAbsPath = filepath.Clean(absPath)
		if fm.ChRoot != "" && fm.ChRoot != "/" {
			rel, err := filepath.Rel(fm.ChRoot, absPath)
			if err != nil {
				return fmt.Errorf("unable to calculate chroot: %s", err)
			}
			fm.path = filepath.Clean("/" + rel)
		} else {
			fm.path = filepath.Clean(absPath)
		}
	} else {
		return errors.New("provided path is not a directory")
	}

	return fm.Scan()
}

// CurPath returns the current absolute path of file manager, either to the root of the FS or to the chroot (if defined)
func (fm *Manager) CurPath() string {
	return fm.path
}

// Scan reads the contents of a dir and load it in memory for further work on
func (fm *Manager) Scan() error {

	fm.init()
	fm.ResetIndex()
	fm.reset()

	files, err := ioutil.ReadDir(fm.realAbsPath)
	if err != nil {
		return err
	}

	for _, file := range files {

		nf := File{
			Name:    file.Name(),
			Size:    file.Size(),
			Mode:    file.Mode(),
			ModTime: file.ModTime(),
			IsDir:   file.IsDir(),
		}

		if strings.HasPrefix(file.Name(), ".") {
			nf.IsHidden = true
		}
		// skip hidden
		if nf.IsHidden && !fm.filter.ShowHiddenFiles {
			continue
		}

		if file.IsDir() {
			fm.dirsCounter++
			nf.FileType = FOLDER
			fm.dirs = append(fm.dirs, nf)
		} else {
			// get type by extension
			ext := filepath.Ext(nf.Name)
			if ext != "" {
				ext = ext[1:]
			} else {
				ext = ""
			}
			ft := fileTypeFromExtension(ext)
			nf.FileType = ft.Type

			// filter files depending on configured filter
			if !fm.filterFile(ext, nf.FileType) {
				continue
			}
			fm.filesCounter++
			fm.files = append(fm.files, nf)
		}
	}
	return nil
}

// reset will reset all data that is modified by scan to the original state
func (fm *Manager) reset() {
	fm.filesCounter = 0
	fm.dirsCounter = 0
	fm.files = []File{}
	fm.dirs = []File{}
}

func (fm *Manager) SetFilter(f FilterCondition) {
	fm.filter = f
	fm.filterCache = map[string]bool{}
}

// return a slice of all scanned files and directories
func (fm *Manager) Files() []File {
	return fm.files
}

// return a slice of all scanned files and directories
func (fm *Manager) Dirs() []File {
	return fm.dirs
}

// return the next item if it exists
func (fm *Manager) Prev() *File {

	if fm.fileIndex >= 0 {
		fm.fileIndex--
		return fm.Current()
	}
	return nil
}

// return the next item if it exists
func (fm *Manager) Next() *File {
	if fm.fileIndex != (fm.filesCounter - 1) {
		fm.fileIndex++
		return fm.Current()
	}
	return nil
}

// return the next item if it exists
func (fm *Manager) Current() *File {
	if len(fm.files) == 0 {
		return nil
	}

	if fm.fileIndex < 0 || fm.fileIndex > len(fm.files) {
		return nil
	}
	return &fm.files[fm.fileIndex]
}

// puts the files index back to the beginning
func (fm *Manager) ResetIndex() {
	fm.fileIndex = 0
}

// puts the files index to the specified place
func (fm *Manager) SetIndex(i int) {
	// todo, this will not work with different filters of hidden and unhidden files
	if i >= fm.filesCounter {
		fm.fileIndex = fm.filesCounter - 1
	} else {
		fm.fileIndex = i
	}
}

// searches the file list for a name, and puts the index on that file
// if the file is not found, the index is put to 0
func (fm *Manager) SetIndexFile(fname string) {
	for i := range fm.Files() {
		if fname == fm.files[i].Name {
			fm.fileIndex = i
			return
		}
	}
}

// CountFiles returns the amount of files the current path, the boolean hidden = true, counts also hidden files
func (fm *Manager) CountFiles() int {
	return fm.filesCounter
}

// CountAllDirs returns the amount of directories the current path, the boolean hidden = true, counts also hidden files
func (fm *Manager) CountDirs() int {
	return fm.dirsCounter
}

// filterFile applies the extension and type filter to a file,
// return is true if the filters match the given type
// about cache: to avoid iterating over the combination of available extensions and fileTypes, this code uses a
// map of [string]bool to quicly check if a file is going to be added or discarded.
func (fm *Manager) filterFile(ext string, typ FileType) bool {

	cacheKey := strings.ToLower(ext) + strconv.Itoa(int(typ))

	v, found := fm.filterCache[cacheKey]
	if found {
		return v
	}

	if len(fm.filter.ShowExtension) == 0 && len(fm.filter.ShowType) == 0 {
		fm.filterCache[cacheKey] = true
		return true
	}

	ShowByExt := true
	showByType := true

	if len(fm.filter.ShowExtension) > 0 {
		ShowByExt = false
		showByType = false
		for i := range fm.filter.ShowExtension {
			if strings.ToLower(ext) == strings.ToLower(fm.filter.ShowExtension[i]) {
				ShowByExt = true
				continue
			}
		}
	}

	if len(fm.filter.ShowType) > 0 {
		if len(fm.filter.ShowExtension) == 0 {
			ShowByExt = false
		}

		showByType = false
		for i := range fm.filter.ShowType {
			if typ == fm.filter.ShowType[i] {
				showByType = true
				continue
			}
		}
	}

	if ShowByExt || showByType {
		fm.filterCache[cacheKey] = true
		return true
	}
	fm.filterCache[cacheKey] = false
	return false
}
