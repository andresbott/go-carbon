package auth_test

import (
	"fmt"
	"github.com/andresbott/go-carbon/libs/auth"
	"net/http"
	"net/http/httptest"
	"testing"
	//"time"
)

type dummyUser struct {
	user string
	pass string
}

func (st dummyUser) AllowLogin(user string, hash string) bool {
	if st.user == "" {
		st.user = "admin"
	}
	if st.pass == "" {
		st.pass = "admin"
	}
	if user == st.user && hash == st.pass {
		return true
	}
	return false
}

func dummyHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//w.WriteHeader(statusCode)
		fmt.Fprint(w, "protected")
	})
}

func TestBasicAuthResponseCode(t *testing.T) {

	tcs := []struct {
		name               string
		request            func() *http.Request
		expectedStatusCode int
	}{
		{
			name: "expect 401 without auth info",
			request: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/bla", nil)
				return req
			},
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			name: "expect 401 on wrong auth credentials",
			request: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/bla", nil)
				req.SetBasicAuth("admin", "wrong")
				return req
			},
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			name: "expect 200 on correct credentials",
			request: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/bla", nil)
				req.SetBasicAuth("admin", "admin")
				return req
			},
			expectedStatusCode: http.StatusOK,
		},
	}

	bauth := auth.Basic{
		User: dummyUser{},
	}

	handler := bauth.Middleware(dummyHandler())

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {

			respRec := httptest.NewRecorder()
			handler.ServeHTTP(respRec, tc.request())
			resp := respRec.Result()

			if resp.StatusCode != tc.expectedStatusCode {
				t.Errorf("got unexpected response code expected: %d, got: %d", tc.expectedStatusCode, resp.StatusCode)
			}
		})
	}

}
