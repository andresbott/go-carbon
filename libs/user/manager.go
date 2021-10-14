package user

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
	"time"
)

type Manager struct {
	// todo use interface instead of database
	db *gorm.DB
}

// NewManager creates an instance of user manager
func NewManager(db *gorm.DB) *Manager {

	// Migrate the schema
	// todo auto migration error handling
	db.AutoMigrate(&user{})

	return &Manager{
		db: db,
	}
}

type user struct {
	gorm.Model
	Name  string
	Email string
	Pw    string
	// last login
	// login location
}

type CreateUserOpts struct {
	Name  string
	Email string
	Pw    string
}

func (mng Manager) CreateUser(opts CreateUserOpts) error {

	if opts.Name == "" {
		return errors.New("name cannot be empty")
	}

	if opts.Email == "" {
		// todo add email structure verifications
		return errors.New("email cannot be empty")
	}

	if opts.Pw == "" {
		// todo pw length verification
		return errors.New("password cannot be empty")
	}

	hashPasswd, err := hashPw(opts.Pw, 15)
	if err != nil {
		return err
	}

	usr := user{
		Name:  opts.Name,
		Email: opts.Email,
		Pw:    hashPasswd,
	}

	mng.db.Create(&usr)
	return nil
}

// CheckLogin checks if the user provided password is correct for login
// if no error is returned login is successful
func (mng Manager) CheckLogin(email string, pw string) bool {

	var usr user

	result := mng.db.First(&usr, "email = ?", email)

	if result.RowsAffected == 0 {
		return false
	}

	return comparePw(pw, usr.Pw)
}

// GenJwtToken generates a signed jwt token
func GenJwtToken() (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"foo": "bar",
		"nbf": time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
	})

	hmacSampleSecret := []byte("secret")
	return token.SignedString(hmacSampleSecret)

}

// use jwt to create a session ?
// login to create a session directly
// - interface session storage ( db, memory etc)
