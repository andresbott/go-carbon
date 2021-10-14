package user

import "golang.org/x/crypto/bcrypt"

// hashPw uses the bcrypt function to generate a hash of a password
func hashPw(in string, cost int) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(in), cost)
	return string(bytes), err
}

// comparePw uses the bcrypt function to evaluate if the passed hash and password
// are the same
func comparePw(pw string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pw))
	return err == nil
}
