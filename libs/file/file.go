package file

import (
	"os"
	"time"
)

type FileType int

const (
	FOLDER FileType = 0
	IMAGE  FileType = 1
	AUDIO  FileType = 2
	VIDEO  FileType = 3
	ZIP    FileType = 4
	TEXT   FileType = 5
	EXEC   FileType = 6
	PDF    FileType = 7
	FILE   FileType = 8
	WEB    FileType = 9
)

type File struct {
	Name     string
	Size     int64
	Mode     os.FileMode
	ModTime  time.Time
	IsDir    bool
	IsHidden bool
	FileType FileType
}
