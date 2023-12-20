package user

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Manager is an opinionated user manager that stores the information on a gorm database
type Manager struct {
	db               *gorm.DB
	bcryptDifficulty int // exposed as parameter for make tests faster
}

type ManagerOpts struct {
	BcryptDifficulty int
}

// NewManager creates an instance of user manager
func NewManager(db *gorm.DB, opts ManagerOpts) (*Manager, error) {

	// Migrate the schema
	err := db.AutoMigrate(&userModel{})
	if err != nil {
		return nil, err
	}

	return &Manager{
		db:               db,
		bcryptDifficulty: opts.BcryptDifficulty, // set the cost of the difficulty
	}, nil
}

// userModel is the database representation of the user
type userModel struct {
	gorm.Model
	Email   string `gorm:"uniqueIndex"`
	Name    string
	Pw      string
	Enabled bool
	// last login
	// login location
}

type User struct {
	Name    string
	Email   string
	Pw      string
	Enabled bool
}

func (mng Manager) Create(id string, pw string) error {
	usr := User{
		Email: id,
		Pw:    pw,
	}
	return mng.CreateUser(usr)
}

func (mng Manager) CreateUser(usr User) error {

	if usr.Email == "" {
		// todo add email structure verifications
		return errors.New("email cannot be empty")
	}

	if usr.Pw == "" {
		// todo pw length and complexity verification
		return errors.New("password cannot be empty")
	}

	// generate bcrypt hashed password
	hashedPasswd, err := bcrypt.GenerateFromPassword([]byte(usr.Pw), mng.bcryptDifficulty)
	if err != nil {
		return err
	}

	usrModel := userModel{
		Name:    usr.Name,
		Email:   usr.Email,
		Pw:      string(hashedPasswd),
		Enabled: usr.Enabled,
	}

	mng.db.Create(&usrModel)
	return nil
}

// CheckLogin checks if the user provided password is correct for login
// if no error is returned login is successful
func (mng Manager) CheckLogin(user string, providedPass string) bool {

	var usr userModel

	result := mng.db.First(&usr, "email = ?", user)

	if result.RowsAffected == 0 {
		return false
	}

	err := bcrypt.CompareHashAndPassword([]byte(usr.Pw), []byte(providedPass))
	return err == nil

}

//// GenJwtToken generates a signed jwt token
//func GenJwtToken() (string, error) {
//
//	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
//		"foo": "bar",
//		"nbf": time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
//	})
//
//	hmacSampleSecret := []byte("secret")
//	return token.SignedString(hmacSampleSecret)
//
//}

// use jwt to create a session ?
// login to create a session directly
// - interface session storage ( db, memory etc)
