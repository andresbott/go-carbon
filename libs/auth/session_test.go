package auth

import (
	"bytes"
	"github.com/davecgh/go-spew/spew"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func TestSessionManagement(t *testing.T) {

	stores := []string{
		useCookieStore,
		useFsStore,
	}

	for _, storeType := range stores {
		t.Run(storeType, func(t *testing.T) {
			t.Run("access resource after login", func(t *testing.T) {
				svr, c := testServer(50*time.Millisecond, 200*time.Millisecond, 0, useCookieStore)
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
				svr, c := testServer(50*time.Millisecond, 500*time.Millisecond, 0, useCookieStore)
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
				svr, c := testServer(50*time.Millisecond, 2000*time.Millisecond, 0, useCookieStore)
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
				svr, c := testServer(500*time.Millisecond, 50*time.Millisecond, 0, useCookieStore)
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

			t.Run("401 forced session to not be updated", func(t *testing.T) {
				svr, c := testServer(50*time.Millisecond, 500*time.Millisecond, 5*time.Minute, useCookieStore)
				defer svr.Close()
				// perform login
				resp := doReq(c, svr.URL+"/login", t)
				want := http.StatusOK
				if resp.StatusCode != want {
					t.Errorf("[login request] got unexpected response code expected: %d, got: %d", want, resp.StatusCode)
				}

				// sleep a bit and trigger a session renew, this is not exercised
				time.Sleep(40 * time.Millisecond)
				// assert user is still logged in
				resp = doReq(c, svr.URL+"/something", t)
				want = http.StatusOK
				if resp.StatusCode != want {
					t.Errorf("[first request] got unexpected response code expected: %d, got: %d", want, resp.StatusCode)
				}

				// sleep another bit and check that session was not renewed
				time.Sleep(40 * time.Millisecond)
				// assert user is logged in
				resp = doReq(c, svr.URL+"/something", t)
				want = http.StatusUnauthorized
				if resp.StatusCode != want {
					t.Errorf("[second request] got unexpected response code expected: %d, got: %d", want, resp.StatusCode)
				}
			})
		})
	}

}

func TestAuthHandler(t *testing.T) {

	tcs := []struct {
		name     string
		password string
		expect   int
	}{
		{
			name:     "valid login",
			password: "admin",
			expect:   200,
		},
		{
			name:     "invalid login",
			password: "nope",
			expect:   401,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			svr, client := testServer(50*time.Millisecond, 500*time.Millisecond, 5*time.Minute, useCookieStore)
			defer svr.Close()

			// perform login
			var param = url.Values{}

			param.Set("Name", "admin")
			param.Set("Pw", tc.password)

			var payload = bytes.NewBufferString(param.Encode())
			request, err := http.NewRequest("POST", svr.URL+"/form-login", payload)

			if err != nil {
				t.Fatal(err)
			}
			request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			response, err := client.Do(request)
			if err != nil {
				t.Fatal(err)
			}

			want := tc.expect
			if response.StatusCode != want {
				t.Errorf("[login request] got unexpected response code expected: %d, got: %d", want, response.StatusCode)
			}
		})
	}
}

//func TestFsStore(t *testing.T) {
//	svr, client := testServer(50*time.Millisecond, 500*time.Millisecond, 5*time.Minute, useFsStore)
//	defer svr.Close()
//
//	resp := doReq(client, svr.URL+"/login", t)
//	want := http.StatusOK
//	if resp.StatusCode != want {
//		t.Errorf("[login request] got unexpected response code expected: %d, got: %d", want, resp.StatusCode)
//	}
//
//	client2 := getClient()
//	resp2 := doReq(client2, svr.URL+"/login", t)
//	if resp2.StatusCode != want {
//		t.Errorf("[login request] got unexpected response code expected: %d, got: %d", want, resp2.StatusCode)
//	}
//
//}

func TestProcessSessionData(t *testing.T) {

	tcs := []struct {
		name string
		in   SessionData
		want SessionData
	}{
		{
			name: "session valid",
			in: SessionData{IsAuthenticated: true,
				Expiration:  getTime("10m"),
				ForceReAuth: getTime("1m"),
			},
			want: SessionData{IsAuthenticated: true},
		},
		{
			name: "session expired",
			in: SessionData{IsAuthenticated: true,
				Expiration: getTime("-1s"),
			},
			want: SessionData{IsAuthenticated: false},
		},
		{
			name: "session NOT expired, but hard logout",
			in: SessionData{IsAuthenticated: true,
				Expiration:  getTime("10m"),
				ForceReAuth: getTime("-1s"),
			},
			want: SessionData{IsAuthenticated: false},
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

func doReq(client *http.Client, url string, t *testing.T) *http.Response {
	req, _ := http.NewRequest(http.MethodGet, url, nil)

	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	return resp
}

const useCookieStore = "cookie"
const useFsStore = "fs"

func testServer(SessionDur, MaxSessionDur, update time.Duration, storeType string) (*httptest.Server, *http.Client) {
	var store sessions.Store
	if storeType == useCookieStore {
		store, _ = CookieStore(securecookie.GenerateRandomKey(64), securecookie.GenerateRandomKey(32))
	}
	if storeType == useFsStore {
		store, _ = FsStore("", securecookie.GenerateRandomKey(64), securecookie.GenerateRandomKey(32))
	}

	auth, err := NewSessionMgr(SessionCfg{
		Store:         store,
		SessionDur:    SessionDur,
		MaxSessionDur: MaxSessionDur,
	})
	auth.minWriteSpace = update

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.RequestURI == "/login" {
			err = auth.Login(r, w, "tester")
			if err != nil {
				spew.Dump(err)
				http.Error(w, "internal error", http.StatusInternalServerError)
				return
			}
			http.Error(w, "ok", http.StatusOK)
		} else if r.RequestURI == "/form-login" {
			user := dummyUser{}
			handler := FormAuthHandler(auth, user)
			handler.ServeHTTP(w, r)
		} else {
			h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "ok", http.StatusOK)
			})
			handler := Middleware(auth, h)
			handler.ServeHTTP(w, r)
		}
	})

	svr := httptest.NewServer(handler)

	jar, _ := cookiejar.New(nil)
	c := svr.Client()
	c.Jar = jar

	return svr, c
}

func getClient() *http.Client {
	c := http.Client{}
	jar, _ := cookiejar.New(nil)
	c.Jar = jar
	return &c
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
