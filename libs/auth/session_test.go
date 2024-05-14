package auth

import (
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/gorilla/securecookie"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"
	"time"
)

func doReq(client *http.Client, url string, t *testing.T) *http.Response {
	req, _ := http.NewRequest(http.MethodGet, url, nil)

	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	return resp
}

func testServer(SessionDur, MaxSessionDur time.Duration) (*httptest.Server, *http.Client) {
	store, _ := CookieStore(securecookie.GenerateRandomKey(64), securecookie.GenerateRandomKey(32))
	auth, err := NewSessionAuth(SessionCfg{
		User:          dummyUser{},
		Store:         store,
		SessionDur:    SessionDur,
		MaxSessionDur: MaxSessionDur,
	})

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.RequestURI == "/login" {
			// todo use a commodity function here
			authData := SessionData{
				UserId:      "tester",
				IsAuth:      true,
				Expiration:  time.Now().Add(auth.sessionDur),
				ForceReAuth: time.Now().Add(auth.maxSessionDur),
			}
			err = auth.WriteSession(r, w, authData)
			if err != nil {
				http.Error(w, "internal error", http.StatusInternalServerError)
				return
			}
			http.Error(w, "ok", http.StatusOK)
		} else {
			h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "ok", http.StatusOK)
			})
			mid := auth.Middleware(h)
			mid.ServeHTTP(w, r)
		}
	})

	svr := httptest.NewServer(handler)

	jar, _ := cookiejar.New(nil)
	c := svr.Client()
	c.Jar = jar

	return svr, c
}

func TestSessionManagement(t *testing.T) {

	t.Run("access resource after login", func(t *testing.T) {
		svr, c := testServer(50*time.Millisecond, 200*time.Millisecond)
		defer svr.Close()

		// assert first request is not logged in
		resp := doReq(c, svr.URL+"/something", t)
		want := http.StatusUnauthorized
		if resp.StatusCode != want {
			t.Errorf("[first request] got unexpected response code expected: %d, got: %d", want, resp.StatusCode)
		}

		// perform login
		resp = doReq(c, svr.URL+"/login", t)
		want = http.StatusOK
		if resp.StatusCode != want {
			t.Errorf("[login request] got unexpected response code expected: %d, got: %d", want, resp.StatusCode)
		}
		// assert user is logged in
		resp = doReq(c, svr.URL+"/something", t)
		want = http.StatusOK
		if resp.StatusCode != want {
			t.Errorf("[login request] got unexpected response code expected: %d, got: %d", want, resp.StatusCode)
		}
	})

	t.Run("401 after session expired", func(t *testing.T) {
		svr, c := testServer(50*time.Millisecond, 500*time.Millisecond)
		defer svr.Close()
		// perform login
		resp := doReq(c, svr.URL+"/login", t)
		want := http.StatusOK
		if resp.StatusCode != want {
			t.Errorf("[login request] got unexpected response code expected: %d, got: %d", want, resp.StatusCode)
		}
		// assert user is logged in
		resp = doReq(c, svr.URL+"/something", t)
		want = http.StatusOK
		if resp.StatusCode != want {
			t.Errorf("[first request] got unexpected response code expected: %d, got: %d", want, resp.StatusCode)
		}
		// sleep longer than the 50ms expiry
		time.Sleep(100 * time.Millisecond)

		// assert user is logged in
		resp = doReq(c, svr.URL+"/something", t)
		want = http.StatusUnauthorized
		if resp.StatusCode != want {
			t.Errorf("[second request] got unexpected response code expected: %d, got: %d", want, resp.StatusCode)
		}
	})

	t.Run("renew session", func(t *testing.T) {
		svr, c := testServer(50*time.Millisecond, 2000*time.Millisecond)
		defer svr.Close()
		// perform login
		resp := doReq(c, svr.URL+"/login", t)
		want := http.StatusOK
		if resp.StatusCode != want {
			t.Errorf("[login request] got unexpected response code expected: %d, got: %d", want, resp.StatusCode)
		}
		// sleep a bit and renew the session
		time.Sleep(40 * time.Millisecond)
		// assert user is still logged in
		resp = doReq(c, svr.URL+"/something", t)
		want = http.StatusOK
		if resp.StatusCode != want {
			t.Errorf("[first request] got unexpected response code expected: %d, got: %d", want, resp.StatusCode)
		}
		// sleep another bit and renew the session
		time.Sleep(40 * time.Millisecond)

		// assert user is logged in
		resp = doReq(c, svr.URL+"/something", t)
		want = http.StatusOK
		if resp.StatusCode != want {
			t.Errorf("[second request] got unexpected response code expected: %d, got: %d", want, resp.StatusCode)
		}
	})

	t.Run("401 after max session duration", func(t *testing.T) {
		svr, c := testServer(500*time.Millisecond, 50*time.Millisecond)
		defer svr.Close()
		// perform login
		resp := doReq(c, svr.URL+"/login", t)
		want := http.StatusOK
		if resp.StatusCode != want {
			t.Errorf("[login request] got unexpected response code expected: %d, got: %d", want, resp.StatusCode)
		}
		// assert user is logged in
		resp = doReq(c, svr.URL+"/something", t)
		want = http.StatusOK
		if resp.StatusCode != want {
			t.Errorf("[first request] got unexpected response code expected: %d, got: %d", want, resp.StatusCode)
		}
		// sleep longer than the 50ms max session duration
		time.Sleep(60 * time.Millisecond)

		// assert user is logged in
		resp = doReq(c, svr.URL+"/something", t)
		want = http.StatusUnauthorized
		if resp.StatusCode != want {
			t.Errorf("[second request] got unexpected response code expected: %d, got: %d", want, resp.StatusCode)
		}
	})

}

func TestProcessSessionData(t *testing.T) {

	tcs := []struct {
		name string
		in   SessionData
		want SessionData
	}{
		{
			name: "session valid",
			in: SessionData{IsAuth: true,
				Expiration:  getTime("10m"),
				ForceReAuth: getTime("1m"),
			},
			want: SessionData{IsAuth: true},
		},
		{
			name: "session expired",
			in: SessionData{IsAuth: true,
				Expiration: getTime("-1s"),
			},
			want: SessionData{IsAuth: false},
		},
		{
			name: "session NOT expired, but hard logout",
			in: SessionData{IsAuth: true,
				Expiration:  getTime("10m"),
				ForceReAuth: getTime("-1s"),
			},
			want: SessionData{IsAuth: false},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {

			got := tc.in
			got.process(0)
			want := tc.want
			if diff := cmp.Diff(want, got, cmpopts.IgnoreFields(SessionData{}, "Expiration", "ForceReAuth")); diff != "" {
				t.Errorf("Content mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func getTime(add string) time.Time {
	if add == "" {
		add = "0s"
	}
	dur, err := time.ParseDuration(add)
	if err != nil {
		panic(err)
	}

	return time.Now().Add(dur)
}
