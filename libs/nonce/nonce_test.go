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

}

func BenchmarkMemoryStore(b *testing.B) {

	insert := func(n int, b *testing.B) {
		s := memStore{
			//expiration: 5 * time.Millisecond,
			db: make(map[string]time.Time),
		}
		for i := 0; i < b.N; i++ {
			val := 0
			check := ""
			for j := 1; j <= n; j++ {
				check = strconv.Itoa(val)
				s.Save(check)
				val++
			}
		}
	}

	b.Run("insert", func(b *testing.B) {
		b.Run("1000", func(b *testing.B) {

			insert(1000, b)
		})
		b.Run("10K", func(b *testing.B) {
			insert(10*1000, b)
		})
		b.Run("100k", func(b *testing.B) {
			insert(100*1000, b)
		})
		b.Run("500k", func(b *testing.B) {
			insert(500*1000, b)
		})
	})

	insertRead := func(n int, b *testing.B) {
		s := memStore{
			db: make(map[string]time.Time),
		}

		for i := 0; i < b.N; i++ {
			val := 0
			check := ""
			for j := 1; j <= n; j++ {
				check = strconv.Itoa(val)
				s.Save(check)
				s.Valid(check)
				val++
			}
		}
	}

	b.Run("insertAndRead", func(b *testing.B) {
		b.Run("1000", func(b *testing.B) {
			insertRead(1000, b)
		})
		b.Run("10K", func(b *testing.B) {
			insertRead(10*1000, b)
		})
		b.Run("100k", func(b *testing.B) {
			insertRead(100*1000, b)
		})
		b.Run("500k", func(b *testing.B) {
			insertRead(500*1000, b)
		})
	})

}
