package user

import (
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"testing"
	"time"
)

const dbFile = "test.db"

func setup() *Manager {

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Silent,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	db, err := gorm.Open(sqlite.Open(dbFile), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic("failed to connect database")
	}

	mng := NewManager(db)

	return mng

}

func clean() {
	defer func() {
		err := os.Remove(dbFile)
		if err != nil {
			panic(err)
		}
	}()
}

func TestCreateUser(t *testing.T) {

	mng := setup()
	defer clean()

	err := mng.CreateUser(CreateUserOpts{
		Name:  "test",
		Email: "test@mail.com",
		Pw:    "1234",
	})

	if err != nil {
		t.Fatalf("unexpected error %s", err)
	}

	// Reads
	var got user
	mng.db.First(&got, 1)

	want := user{
		Name:  "test",
		Email: "test@mail.com",
	}

	if diff := cmp.Diff(want, got, cmpopts.IgnoreFields(user{}, "Model", "Pw")); diff != "" {
		t.Errorf("Content mismatch (-want +got):\n%s", diff)
	}
}

func TestLogin(t *testing.T) {

	mng := setup()
	defer clean()

	mng.CreateUser(CreateUserOpts{
		Name:  "test",
		Email: "test@mail.com",
		Pw:    "1234",
	})

	t.Run("assert correct login", func(t *testing.T) {
		got := mng.CheckLogin("test@mail.com", "1234")
		if got != true {
			t.Errorf("expecting login failure")
		}
	})

	t.Run("assert wrong password login", func(t *testing.T) {
		got := mng.CheckLogin("test@mail.com", "12345")
		if got != false {
			t.Errorf("expecting login failure")
		}
	})

	t.Run("assert wrong user name", func(t *testing.T) {
		got := mng.CheckLogin("test_@mail.com", "1234")
		if got != false {
			t.Errorf("expecting login failure")
		}
	})

}
