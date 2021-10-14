package file_test

import (
	"fmt"
	"git.andresbott.com/Golang/carbon/libs/file"
	"io/ioutil"
	"os"
	"strings"
)

func ExampleFileManager() {
	// WARNING: this example ignores errors for simplicity, you should handle them correctly
	tmpDir, err := prepareExample()
	if err != nil {
		fmt.Errorf("error while preparing example: %s", err)
		return
	}
	defer func() {
		err := os.RemoveAll(tmpDir)
		if err != nil {
			fmt.Printf("ERROR: unable to delete tmpDir: %s", err)
		}
	}()

	// create a new file manager
	fm := file.Manager{
		ChRoot: tmpDir + "/root/", // limit all interaction to that folder and bellow, similar to linux chroot
	}

	// scan root directory
	_ = fm.Scan()
	rootDirs := fm.Dirs()
	fmt.Print("directory \"root\" contains: ")
	for _, d := range rootDirs {
		fmt.Printf("\"%s\"", d.Name)
	}
	fmt.Print("\n")

	// navigate to (and scan) dir1
	_ = fm.Nav(fm.CurPath() + "/dir1")
	// count the files in the directory (it contains 11 non-hidden files)
	fileCount := fm.CountFiles()
	fmt.Printf("directory \"root/dir1\" contains %d files\n", fileCount)

	// change the filter to get only audio files and include hidden files
	fm.SetFilter(file.FilterCondition{
		ShowHiddenFiles: true,
		ShowExtension:   nil,
		ShowType:        []file.FileType{file.AUDIO},
	})
	// scan again to apply filter
	_ = fm.Scan()
	// count the files in the directory
	fileCount = fm.CountFiles()
	fmt.Printf("directory \"root/dir1\" contains %d audio files: ", fileCount)

	// print the file names found
	var files []string
	for _, d := range fm.Files() {
		files = append(files, fmt.Sprintf("\"%s\"", d.Name))
	}
	fmt.Println(strings.Join(files, ", "))

	// Output:
	// directory "root" contains: "dir1"
	// directory "root/dir1" contains 11 files
	// directory "root/dir1" contains 4 audio files: ".hidden4.mp3", "file4.mp3", "file6.mp3", "file8.midi"
}

func ExampleNavigation() {
	// WARNING: this example ignores errors for simplicity, you should handle them correctly
	tmpDir, err := prepareExample()
	if err != nil {
		fmt.Errorf("error while preparing example: %s", err)
		return
	}
	defer func() {
		err := os.RemoveAll(tmpDir)
		if err != nil {
			fmt.Printf("ERROR: unable to delete tmpDir: %s", err)
		}
	}()

	// create a new file manager
	fm := file.Manager{
		ChRoot: tmpDir + "/root/", // limit all interaction to that folder and bellow, similar to linux chroot
	}

	fm.SetFilter(file.FilterCondition{
		ShowHiddenFiles: true,
		ShowExtension:   []string{"jpg"},
		ShowType:        nil,
	})
	_ = fm.Nav(fm.CurPath() + "/dir1")
	// files in the selection: .hidden1.jpg, .hidden3.jpg, file1.jpg, file2.JPG, file3.Jpg

	var f *file.File
	// print the first file in the selection:
	f = fm.Current()
	fmt.Println(f.Name)

	// print the next file
	f = fm.Next()
	fmt.Println(f.Name)

	// see that current has been updated
	f = fm.Current()
	fmt.Println(f.Name)

	// move index to the last position
	fm.SetIndex(4)
	f = fm.Current()
	fmt.Println(f.Name)

	// go one back
	f = fm.Prev()
	fmt.Println(f.Name)

	// go twice over the last item
	fm.Next()
	fm.Next()
	fm.Next()

	f = fm.Current()
	fmt.Println(f.Name)

	// Output:
	// .hidden1.jpg
	// .hidden3.jpg
	// .hidden3.jpg
	// file3.Jpg
	// file2.JPG
	// file3.Jpg

}

// prepareExample creates a tmpdir with fake data to be used in the example test
func prepareExample() (string, error) {
	dir, err := ioutil.TempDir("", "efile_example")
	if err != nil {
		return "", err
	}
	root := dir + "/root" // tmp/root
	os.MkdirAll(root, os.ModePerm)
	// folder structure for counting files and folders
	nav := root + "/dir1" // tmp/root/dir1
	os.MkdirAll(nav, os.ModePerm)
	os.MkdirAll(nav+"/d1", os.ModePerm) // tmp/root/dir1/d1
	_, _ = os.Create(nav + "/d1" + "/img1.jpg")
	_, _ = os.Create(nav + "/d1" + "/img2.jpg")
	_, _ = os.Create(nav + "/d1" + "/img3.jpg")
	_, _ = os.Create(nav + "/d1" + "/img4.jpg")
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
	_, _ = os.Create(nav + "/file8.midi")
	_, _ = os.Create(nav + "/file9.png")
	_, _ = os.Create(nav + "/file9.gif")
	_, _ = os.Create(nav + "/file9.avi")
	_, _ = os.Create(nav + "/.hidden1.jpg")
	_, _ = os.Create(nav + "/.hiddenFile")
	_, _ = os.Create(nav + "/.hidden3.jpg")
	_, _ = os.Create(nav + "/.hidden4.mp3")

	return dir, nil
}
