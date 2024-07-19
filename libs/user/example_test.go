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

func ExampleStaticUsers() {

	// load Users from a file
	userMng, err := FromFile("testdata/Users.yaml")
	panicOnErr(err)

	// manually add a user
	u := User{Name: "u2", Pw: "u2", Enabled: true}
	userMng.Users = append(userMng.Users, u)

	// check if the user demo (from file) can login
	isOK := userMng.AllowLogin("demo", "demo")
	fmt.Printf("user demo can login: %v\n", isOK)

	// check if the user u2 can login
	isOK = userMng.AllowLogin("u2", "u2")
	fmt.Printf("user u2 can login: %v", isOK)

	// Output:
	// user demo can login: true
	// user u2 can login: true
}

func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}
