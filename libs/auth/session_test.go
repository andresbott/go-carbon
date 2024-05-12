package auth

import (
	"github.com/gorilla/securecookie"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSessionAuth(t *testing.T) {

	tcs := []struct {
		name string
		data SessionData

		expectedStatusCode int
	}{
		{
			name:               "expect 401 with auth false",
			data:               SessionData{IsAuth: false},
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			name:               "expect 200 with auth ok",
			data:               SessionData{IsAuth: true},
			expectedStatusCode: http.StatusOK,
		},
	}

	store, _ := CookieStore(securecookie.GenerateRandomKey(64), securecookie.GenerateRandomKey(32))
	auth, err := NewSessionAuth(SessionCfg{
		User:  dummyUser{},
		Store: store,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	handler := auth.Middleware(dummyHandler())

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {

			respRec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/bla", nil)
			err = auth.WriteSessionData(req, respRec, tc.data)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			handler.ServeHTTP(respRec, req)
			resp := respRec.Result()

			if resp.StatusCode != tc.expectedStatusCode {
				t.Errorf("got unexpected response code expected: %d, got: %d", tc.expectedStatusCode, resp.StatusCode)
			}
		})
	}
}
