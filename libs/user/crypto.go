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

// nolint:nestif // Reason: nestif flagging false complexity
func checkPass(plainPass, hash string) (bool, error) {

	if strings.HasPrefix(hash, SHA1Prefix) {
		// sha1
		b64hash := strings.TrimPrefix(hash, SHA1Prefix)
		hashed, err := base64.StdEncoding.DecodeString(b64hash)
		if err != nil {
			return false, fmt.Errorf("malformed sha1 hash: %s", err.Error())
		}
		if len(hashed) != sha1.Size {
			return false, fmt.Errorf("malformed sha1 wrong length")
		}
		st := sha1.Sum([]byte(plainPass))
		if subtle.ConstantTimeCompare(st[:], hashed) == 1 {
			return true, nil
		}
		return false, nil
	} else if isbCryptString(hash) {
		// bcrypt
		err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(plainPass))
		if err != nil {
			return false, err
		}
		return true, nil
	} else {
		// plain
		if plainPass == hash {
			return true, nil
		} else {
			return false, nil
		}
	}
}

func isbCryptString(hash string) bool {
	if strings.HasPrefix(hash, BCryp1PRefix) {
		return true
	}
	if len(hash) >= 3 && slices.Contains(bCryptPrefix, hash[:4]) {
		return true
	}
	return false
}
