package file_test

import (
	"fmt"
	"git.andresbott.com/Golang/carbon/libs/file"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func createTmpDir() (string, error) {
	dir, err := ioutil.TempDir("", "file_test")
	if err != nil {
		return "", err
	}

	root := dir + "/root"
	os.MkdirAll(root, os.ModePerm)

	_, err = os.Create(root + "/image.jpg")
	if err != nil {
		return "", err
	}

	dir1 := root + "/dir1"
	os.MkdirAll(dir1, os.ModePerm)

	_, _ = os.Create(dir1 + "/image1.jpg")
	_, _ = os.Create(dir1 + "/image2.jpg")
	_, _ = os.Create(dir1 + "/image3.jpg")
	_, _ = os.Create(dir1 + "/text.pdf")

	mp3 := dir1 + "/mp3"
	os.MkdirAll(mp3, os.ModePerm)
	_, _ = os.Create(mp3 + "/song1.mp3")
	_, _ = os.Create(mp3 + "/song2.mp3")
	_, _ = os.Create(mp3 + "/song3.mp3")

	dir2 := root + "/dir2"
	os.MkdirAll(dir2, os.ModePerm)

	space := root + "/dir2 with spaces in name and special chars like ñ and ß."
	os.MkdirAll(space, os.ModePerm)
	_, _ = os.Create(space + "/image1.jpg")

	// folder structure for counting files and folders
	count := root + "/count"
	os.MkdirAll(count, os.ModePerm)
	os.MkdirAll(count+"/d1", os.ModePerm)
	os.MkdirAll(count+"/d2", os.ModePerm)
	os.MkdirAll(count+"/d3", os.ModePerm)
	os.MkdirAll(count+"/.hidden1", os.ModePerm)
	os.MkdirAll(count+"/.hidden2", os.ModePerm)
	_, _ = os.Create(count + "/file1.jpg")
	_, _ = os.Create(count + "/file2.jpg")
	_, _ = os.Create(count + "/file3.jpg")
	_, _ = os.Create(count + "/file4.jpg")
	_, _ = os.Create(count + "/audio1.mp3")
	_, _ = os.Create(count + "/.hidden1.jpg")
	_, _ = os.Create(count + "/.hidden2.jpg")
	_, _ = os.Create(count + "/.hidden3.jpg")

	// folder structure for counting files and folders
	nav := root + "/navigation e2e"
	os.MkdirAll(nav, os.ModePerm)
	os.MkdirAll(nav+"/d1", os.ModePerm)
	os.MkdirAll(nav+"/d2", os.ModePerm)
	os.MkdirAll(nav+"/d3", os.ModePerm)
	os.MkdirAll(nav+"/.hidden1", os.ModePerm)
	os.MkdirAll(nav+"/.hidden2", os.ModePerm)
	_, _ = os.Create(nav + "/file1.jpg")
	_, _ = os.Create(nav + "/file2.JPG")
	_, _ = os.Create(nav + "/file3.Jpg")
	_, _ = os.Create(nav + "/file4.mp3")
	_, _ = os.Create(nav + "/file5.txt")
	_, _ = os.Create(nav + "/file6.mp3")
	_, _ = os.Create(nav + "/file7")
	_, _ = os.Create(nav + "/file8.mp3")
	_, _ = os.Create(nav + "/file9.png")
	_, _ = os.Create(nav + "/file9.gif")
	_, _ = os.Create(nav + "/file9.avi")
	_, _ = os.Create(nav + "/.hidden1.jpg")
	_, _ = os.Create(nav + "/.hiddenFile")
	_, _ = os.Create(nav + "/.hidden3.jpg")
	_, _ = os.Create(nav + "/.hidden4.mp3")

	return dir, nil
}

func TestNewFileManager_Nav(t *testing.T) {

	tmpDir, err := createTmpDir()
	if err != nil {
		t.Error(err)
	}
	defer func() {
		err := os.RemoveAll(tmpDir)
		if err != nil {
			fmt.Printf("ERROR: unable to delete tmpDir: %s", err)
		}
	}()

	type test struct {
		name         string
		path         string
		chdir        string
		workingDir   string
		expected     string
		expecterdErr string
	}

	tcs := []test{
		{
			name:       "assert current working dir",
			path:       "",
			workingDir: tmpDir + "/root/dir1",
			expected:   tmpDir + "/root/dir1",
		},
		{
			name:     "assert relative sub dir navigation",
			path:     "root/dir1/mp3",
			expected: tmpDir + "/root/dir1/mp3",
		},
		{
			name:     "assert absolute sub dir navigation",
			path:     tmpDir + "/root/dir1/mp3",
			expected: tmpDir + "/root/dir1/mp3",
		},
		{
			name:     "assert access on chdir",
			path:     "/dir1/mp3",
			chdir:    tmpDir + "/root/",
			expected: "/dir1/mp3",
		},
		{
			// when using chroot, working directories need to be ignored, and cwd() is considered the root of the chroot
			name:       "assert relative path with wd",
			path:       "dir1/mp3",
			chdir:      tmpDir + "/root/",
			workingDir: tmpDir + "/root/dir1/",
			expected:   "/dir1/mp3",
		},
		{
			name:     "assert not exiting chroot",
			path:     "dir1/././../..",
			chdir:    tmpDir + "/root/",
			expected: "/",
		},
		{
			name:     "assert special chars",
			path:     "dir2 with spaces in name and special chars like ñ and ß.",
			chdir:    tmpDir + "/root/",
			expected: "/dir2 with spaces in name and special chars like ñ and ß.",
		},
		{
			name:         "passing a file",
			path:         tmpDir + "/root/dir1/mp3/song1.mp3",
			expected:     tmpDir + "/root/dir1/mp3",
			expecterdErr: "provided path is not a directory",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {

			if tc.workingDir != "" {
				err := os.Chdir(tc.workingDir)
				if err != nil {
					t.Fatal(err)
				}
			}
			defer os.Chdir(tmpDir)

			fm := file.Manager{ChRoot: tc.chdir}

			err := fm.Nav(tc.path)

			if tc.expecterdErr != "" {
				if err == nil {
					t.Fatalf("expeciting an error but got: <nil>")
				}
				if !strings.Contains(err.Error(), tc.expecterdErr) {
					t.Fatalf("error messages do not match, expected: %s bug got %s", tc.expecterdErr, err.Error())
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %s", err)
				}
				if fm.CurPath() != tc.expected {
					t.Errorf("path does not match, expected: %s bug got %s ", tc.expected, fm.CurPath())
				}
			}
		})
	}

}

func TestFileManager_Scan(t *testing.T) {
	tmpDir, err := createTmpDir()
	if err != nil {
		t.Error(err)
	}
	defer func() {
		err := os.RemoveAll(tmpDir)
		if err != nil {
			fmt.Printf("ERROR: unable to delete tmpDir: %s", err)
		}
	}()

	type test struct {
		name       string
		path       string
		chdir      string
		workingDir string
		expected   string
	}

	tcs := []test{
		{
			name:     "assert files in sub dir",
			path:     "dir1",
			chdir:    tmpDir + "/root",
			expected: "image1.jpg,image2.jpg,image3.jpg,text.pdf,",
		},
		{
			name:       "assert files in sub dir and chroot",
			path:       tmpDir + "/root/dir1/mp3",
			workingDir: tmpDir + "/root",
			expected:   "song1.mp3,song2.mp3,song3.mp3,",
		},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			fm := file.Manager{ChRoot: tc.chdir}

			fm.SetFilter(file.FilterCondition{})
			err := fm.Nav(tc.path)
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			got := ""
			for _, f := range fm.Files() {
				got = got + f.Name + ","
			}

			if got != tc.expected {
				t.Errorf("content does not match, expected: %s bug got %s ", tc.expected, got)
			}

		})
	}
}

func TestFileManager_Filter(t *testing.T) {
	tmpDir, err := createTmpDir()
	if err != nil {
		t.Error(err)
	}
	defer func() {
		err := os.RemoveAll(tmpDir)
		if err != nil {
			fmt.Printf("ERROR: unable to delete tmpDir: %s", err)
		}
	}()

	os.Chdir(tmpDir + "/root")

	type test struct {
		name   string
		path   string
		filter file.FilterCondition

		expected string
	}

	tcs := []test{
		{
			name: "assert hidden files not visible",
			path: "navigation e2e",
			filter: file.FilterCondition{
				ShowHiddenFiles: false,
			},
			expected: "file1.jpg,file2.JPG,file3.Jpg,file4.mp3,file5.txt,file6.mp3,file7,file8.mp3,file9.avi,file9.gif,file9.png,",
		},
		{
			name: "assert hidden files visible",
			path: "navigation e2e",
			filter: file.FilterCondition{
				ShowHiddenFiles: true,
			},
			expected: ".hidden1.jpg,.hidden3.jpg,.hidden4.mp3,.hiddenFile,file1.jpg,file2.JPG,file3.Jpg,file4.mp3,file5.txt,file6.mp3,file7,file8.mp3,file9.avi,file9.gif,file9.png,",
		},
		{
			name: "assert only jpg extension",
			path: "navigation e2e",
			filter: file.FilterCondition{
				ShowHiddenFiles: true,
				ShowExtension: []string{
					"jpg",
				},
			},
			expected: ".hidden1.jpg,.hidden3.jpg,file1.jpg,file2.JPG,file3.Jpg,",
		},
		{
			name: "assert only mp3 extension AND videos",
			path: "navigation e2e",
			filter: file.FilterCondition{
				ShowHiddenFiles: false,
				ShowExtension: []string{
					"mp3",
				},
				ShowType: []file.FileType{
					file.IMAGE,
					file.VIDEO,
				},
			},
			expected: "file1.jpg,file2.JPG,file3.Jpg,file4.mp3,file6.mp3,file8.mp3,file9.avi,file9.gif,file9.png,",
		},
		{
			name: "assert only audio types",
			path: "navigation e2e",
			filter: file.FilterCondition{
				ShowHiddenFiles: false,
				ShowType: []file.FileType{
					file.AUDIO,
				},
			},
			expected: "file4.mp3,file6.mp3,file8.mp3,",
		},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {

			fm := file.Manager{}

			fm.SetFilter(tc.filter)
			err := fm.Nav(tc.path)

			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			got := ""
			for _, f := range fm.Files() {
				got = got + f.Name + ","
			}
			if got != tc.expected {
				t.Errorf("result does not match, expected: %s bug got %s ", tc.expected, got)
			}
		})
	}
}

func TestFileManager_Count(t *testing.T) {
	tmpDir, err := createTmpDir()
	if err != nil {
		t.Error(err)
	}
	defer func() {
		err := os.RemoveAll(tmpDir)
		if err != nil {
			fmt.Printf("ERROR: unable to delete tmpDir: %s", err)
		}
	}()

	type test struct {
		name          string
		path          string
		chroot        string
		filter        file.FilterCondition
		expectedFiles int
		expectedDirs  int
	}

	tcs := []test{
		{
			name:          "count all",
			path:          "count",
			chroot:        tmpDir + "/root",
			filter:        file.FilterCondition{ShowHiddenFiles: true},
			expectedFiles: 8,
			expectedDirs:  5,
		},
		{
			name:          "count non-hidden",
			path:          "count",
			chroot:        tmpDir + "/root",
			filter:        file.FilterCondition{ShowHiddenFiles: false},
			expectedFiles: 5,
			expectedDirs:  3,
		},
		{
			name:          "count non hidden and filtered",
			path:          "count",
			chroot:        tmpDir + "/root",
			filter:        file.FilterCondition{ShowHiddenFiles: false, ShowType: []file.FileType{file.AUDIO}},
			expectedFiles: 1,
			expectedDirs:  3,
		},
		{
			name:          "count all jpg",
			path:          "count",
			chroot:        tmpDir + "/root",
			filter:        file.FilterCondition{ShowHiddenFiles: true, ShowExtension: []string{"jpg"}},
			expectedFiles: 7,
			expectedDirs:  5,
		},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {

			fm := file.Manager{ChRoot: tc.chroot}

			fm.SetFilter(tc.filter)
			err := fm.Nav(tc.path)

			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			nFiles := fm.CountFiles()
			if nFiles != tc.expectedFiles {
				t.Errorf("files count does not match, expected: %d but got %d ", tc.expectedFiles, nFiles)
			}

			nDirs := fm.CountDirs()
			if nDirs != tc.expectedDirs {
				t.Errorf("direcotries count does not match, expected: %d but got %d ", tc.expectedDirs, nDirs)
			}

		})
	}
}

func TestFileManager_e2e_Navigation(t *testing.T) {
	tmpDir, err := createTmpDir()
	if err != nil {
		t.Error(err)
	}
	defer func() {
		err := os.RemoveAll(tmpDir)
		if err != nil {
			fmt.Printf("ERROR: unable to delete tmpDir: %s", err)
		}
	}()

	fm := file.Manager{ChRoot: tmpDir + "/root"}

	err = fm.Nav("/navigation e2e")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	t.Run("get first item", func(t *testing.T) {
		item0 := fm.Current()
		expect := "file1.jpg"
		if item0.Name != expect {
			t.Errorf("result does not match, expected: %s but got %s ", expect, item0.Name)
		}
	})

	t.Run("change index", func(t *testing.T) {
		fm.SetIndex(8)
		item := fm.Next()
		expect := "file9.gif"
		if item.Name != expect {
			t.Errorf("result does not match, expected: %s but got %s ", expect, item.Name)
		}
	})

	t.Run("reset index", func(t *testing.T) {
		fm.ResetIndex()
		item := fm.Current()
		expect := "file1.jpg"
		if item.Name != expect {
			t.Errorf("result does not match, expected: %s but got %s ", expect, item.Name)
		}
	})

	t.Run("change index out of boundaries", func(t *testing.T) {
		fm.SetIndex(100)
		item := fm.Current()
		expect := "file9.png"
		if item.Name != expect {
			t.Errorf("result does not match, expected: %s but got %s ", expect, item.Name)
		}
	})

	t.Run("get previous", func(t *testing.T) {
		fm.SetIndex(2)
		item := fm.Prev()
		expect := "file2.JPG"
		if item.Name != expect {
			t.Errorf("result does not match, expected: %s but got %s ", expect, item.Name)
		}

		item = fm.Prev()
		expect = "file1.jpg"
		if item.Name != expect {
			t.Errorf("result does not match, expected: %s but got %s ", expect, item.Name)
		}
	})

	t.Run("previous expects nil", func(t *testing.T) {
		item := fm.Prev()
		if item != nil {
			t.Errorf("result does not match, expected: <nil> but got %s", item.Name)
		}
	})

	t.Run("set filename index and Prev", func(t *testing.T) {
		fm.SetIndexFile("file4.mp3")
		item := fm.Prev()
		expect := "file3.Jpg"
		if item.Name != expect {
			t.Errorf("result does not match, expected: %s but got %s ", expect, item.Name)
		}
	})

	t.Run("set filename index and Current", func(t *testing.T) {
		fm.SetIndexFile("file4.mp3")
		item := fm.Current()
		expect := "file4.mp3"
		if item.Name != expect {
			t.Errorf("result does not match, expected: %s but got %s ", expect, item.Name)
		}
	})

	t.Run("set filename index and Next", func(t *testing.T) {
		fm.SetIndexFile("file4.mp3")
		item := fm.Next()
		expect := "file5.txt"
		if item.Name != expect {
			t.Errorf("result does not match, expected: %s but got %s ", expect, item.Name)
		}
	})
}
