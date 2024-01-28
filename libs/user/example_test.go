package user

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func ExampleUserManager() {

	// create a gorm DB
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	panicOnErr(err)

	// set some options
	opts := ManagerOpts{
		BcryptDifficulty: bcrypt.MinCost,
	}

	userMng, err := NewDbManager(db, opts)
	panicOnErr(err)

	// create a user
	err = userMng.CreateUser(User{
		Name:  "test",
		Email: "test@mail.com",
		Pw:    "1234",
	})
	panicOnErr(err)

	// check if the user can login
	isOK := userMng.AllowLogin("test@mail.com", "1234")
	fmt.Printf("user can login: %v", isOK)
	// Output: user can login: true
}

func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}
