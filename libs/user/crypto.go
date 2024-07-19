package user

import (
	"crypto/sha1"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"slices"
	"strings"
)

const (
	SHA1Prefix   = "{SHA}"
	BCryp1PRefix = "$2$"
	BCryp2PRefix = "$2a$"
	BCryp3PRefix = "$2b$"
	BCryp4PRefix = "$2x$"
	BCryp5PRefix = "$2y$"
)

var bCryptPrefix = []string{
	BCryp2PRefix,
	BCryp3PRefix,
	BCryp4PRefix,
	BCryp5PRefix,
}

func checkPass(provided, stored string) (bool, error) {

	if strings.HasPrefix(stored, SHA1Prefix) {
		// sha1
		b64 := strings.TrimPrefix(stored, SHA1Prefix)
		hashed, err := base64.StdEncoding.DecodeString(b64)
		if err != nil {
			return false, fmt.Errorf("malformed sha1(%s): %s", stored, err.Error())
		}
		if len(hashed) != sha1.Size {
			return false, fmt.Errorf("malformed sha1(%s): wrong length", stored)
		}
		st := sha1.Sum([]byte(provided))
		if subtle.ConstantTimeCompare(st[:], hashed) == 1 {
			return true, nil
		}
	} else if strings.HasPrefix(stored, BCryp1PRefix) || (len(stored) >= 3 && slices.Contains(bCryptPrefix, stored[:4])) {
		// bcrypt
		err := bcrypt.CompareHashAndPassword([]byte(stored), []byte(provided))
		if err != nil {
			return false, err
		}
		return true, nil
	} else {
		// plain
		if provided == stored {
			return true, nil
		} else {
			return false, nil
		}
	}

	return false, nil
}
