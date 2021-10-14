package user

import (
	"testing"
)

func TestComparePw(t *testing.T) {

	t.Run("assert correct password", func(t *testing.T) {
		pw := "123456789"

		hash1, _ := hashPw(pw, 4)
		got := comparePw(pw, hash1)
		if got == false {
			t.Errorf("password does not match")
		}

		hash2, _ := hashPw(pw, 4)
		got = comparePw(pw, hash2)
		if got == false {
			t.Errorf("password does not match")
		}
	})

	t.Run("assert wrong password", func(t *testing.T) {
		pw := "123456789"
		hash1, _ := hashPw(pw, 4)

		got := comparePw("no match", hash1)

		if got == true {
			t.Errorf("password verification returns true")
		}
	})
}
