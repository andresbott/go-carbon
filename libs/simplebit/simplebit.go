package simplebit

import (
	"fmt"
	"strings"
)

// Get returns true if the bit on position N (starting from right is 1)
func Get(b byte, n int) bool {
	return b&(1<<n) > 0
}

func Set(b byte, n int) byte {
	return b | (1 << n)
}

func Clear(b byte, n int) byte {
	return b &^ (1 << n)
}

func Flip(b byte, n int) byte {
	if Get(b, n) {
		return Clear(b, n)
	}
	return Set(b, n)
}

func StringByte(b byte) string {
	return fmt.Sprintf("%08b", b)
}
func StringBytes(b []byte) string {
	var buf strings.Builder
	for _, n := range b {
		buf.WriteString(fmt.Sprintf("%08b ", n))
	}
	return buf.String()
}
