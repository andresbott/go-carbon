package limitedbuffer

import (
	"bytes"
	"fmt"
)

// Buffer is a wrapper around bytes.Buffer than only accepts up to certain size
// it implements the io.ReadWriter interface
type Buffer struct {
	Size int
	Buf  *bytes.Buffer
}

func (b Buffer) Write(p []byte) (n int, err error) {
	if len(p)+b.Buf.Len() >= b.Size {
		return len(p), fmt.Errorf("buf limit reached")
	} else {
		b.Buf.Write(p)
	}
	return len(p), nil
}

func (b Buffer) Read(p []byte) (n int, err error) {
	return b.Buf.Read(p)
}
