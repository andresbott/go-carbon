package nonce

import (
	"strconv"
	"testing"
	"time"
)

func TestGet(t *testing.T) {

	cfg := Cfg{
		Duration: 100 * time.Millisecond,
		Store:    nil,
		Len:      10,
	}
	nonce := New(cfg)

	a := nonce.Get()
	if a == "" {
		t.Error("expect nonce, got empty string")
	}
	if len(a) != 10 {
		t.Error("nonce length does not match expected: 10")
	}
}

func TestValidate(t *testing.T) {

	cfg := Cfg{
		Duration: 2 * time.Millisecond,
		Store:    nil,
		Len:      10,
	}
	nonce := New(cfg)
	a := nonce.Get()

	res := nonce.Validate(a)
	if res != true {
		t.Error("expect nonce to be valid but got invalid")
	}
	// wait for expirations
	time.Sleep(6 * time.Millisecond)
	res = nonce.Validate(a)
	if res != false {
		t.Error("expect nonce to be invalid but got valid")
	}

}

func TestMemStore(t *testing.T) {

	s := memStore{
		expiration: 5 * time.Millisecond,
		db:         make(map[string]time.Time),
	}
	t.Run("validate key", func(t *testing.T) {
		s.Save("a")
		s.Save("b")
		s.Save("c")
		valid1 := s.Valid("b")
		if valid1 != true {
			t.Error("expected the string to be valid, got false")
		}
	})

	t.Run("validate expiration", func(t *testing.T) {
		s.Save("1a")
		s.Save("1b")
		s.Save("1c")

		valid := s.Valid("1b")
		if valid != true {
			t.Error("expected the string to be valid, got false")
		}

		time.Sleep(6 * time.Millisecond)

		valid = s.Valid("1b")
		if valid != false {
			t.Error("expected the string to be invalid, got true")
		}
	})

	t.Run("testMemPerformance", func(t *testing.T) {
		timeStart := time.Now()
		s.expiration = 1 * time.Millisecond

		c := 0
		lookups := 0
		for i := 1; i <= 3000; i++ {
			check := ""
			for j := 1; j <= 2000; j++ {
				check = strconv.Itoa(c)
				s.Save(check)
				c++
			}
			//n := check + "- "+strconv.Itoa(len(s.db))
			//spew.Dump(n)
			lookups++
			s.Valid(check)
		}

		timeEnd := time.Now()
		timeDiff := timeEnd.Sub(timeStart)
		t.Log("insertions", c, "lookups", lookups, "duration", timeDiff.String())

		//fmt.Println(timeDiff.String())
	})

}
