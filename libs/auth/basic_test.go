package auth

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type dummyUser struct {
}

func (st dummyUser) AllowLogin(user string, hash string) bool {
	if user == "admin" && hash == "admin" {
		return true
	}
	return false
}

func TestBasicAuth(t *testing.T) {

	bauth := Basic{
		User: dummyUser{},
	}

	svr := httptest.NewServer(
		bauth.Middleware(func(writer http.ResponseWriter, request *http.Request) {
			fmt.Fprintln(writer, "protected")
		}),
	)
	defer svr.Close()

	t.Run("expect 401 without auth info", func(t *testing.T) {
		resp, err := http.Get(svr.URL)
		if err != nil {
			t.Error(err)
		}
		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("got unexpected response code expected: %d, got: %d", http.StatusUnauthorized, resp.StatusCode)
		}
	})

	t.Run("expect 401 on wrong auth credentials", func(t *testing.T) {

		client := http.Client{Timeout: 5 * time.Second}

		req, err := http.NewRequest(http.MethodGet, svr.URL, http.NoBody)
		if err != nil {
			t.Error(err)
		}

		req.SetBasicAuth("admin", "wrong")

		resp, err := client.Do(req)
		if err != nil {
			t.Error(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("got unexpected response code expected: %d, got: %d", http.StatusUnauthorized, resp.StatusCode)
		}
	})

	t.Run("expect 200 on correct credentials", func(t *testing.T) {

		client := http.Client{Timeout: 5 * time.Second}

		req, err := http.NewRequest(http.MethodGet, svr.URL, http.NoBody)
		if err != nil {
			t.Error(err)
		}

		req.SetBasicAuth("admin", "admin")

		resp, err := client.Do(req)
		if err != nil {
			t.Error(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("got unexpected response code expected: %d, got: %d", http.StatusOK, resp.StatusCode)
		}
	})

}
