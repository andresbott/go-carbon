package auth_test

import (
	"fmt"
	"github.com/andresbott/go-carbon/libs/auth"
	"github.com/gorilla/securecookie"
	"net/http"
	"net/http/httptest"
	"time"
)

// nolint: govet
func ExampleBasicAuth() {
	protectedSite := dummyHandler()
	// create an instance of basic auth
	basicAuth := auth.Basic{
		User: dummyUser{ // pass a UserLogin
			user: "demo",
			pass: "demo",
		},
	}
	// use the middleware to protect the page
	protectedHandler := basicAuth.Middleware(protectedSite)

	// the client will make a request with credentials
	req := httptest.NewRequest(http.MethodGet, "/some/page", nil)
	req.SetBasicAuth("demo", "demo")

	// check the response
	respRec := httptest.NewRecorder()
	protectedHandler.ServeHTTP(respRec, req)
	resp := respRec.Result()
	fmt.Println(resp.StatusCode)

	// Output: 200

}

// nolint: govet
func ExampleSessionAuth() {

	protectedSite := dummyHandler()

	// create a session store:
	store, _ := auth.CookieStore(securecookie.GenerateRandomKey(64), securecookie.GenerateRandomKey(32))
	// create an instance of session auth
	sessionAuth, _ := auth.NewSessionMgr(auth.SessionCfg{
		Store:         store,
		SessionDur:    time.Hour,       // time the user is logged in
		MaxSessionDur: 24 * time.Hour,  // time after the user is forced to re-login anyway
		MinWriteSpace: 2 * time.Minute, // throttle write operations on the session
	})

	// make a call to the loging handler
	loginReq, _ := http.NewRequest(http.MethodGet, "", nil)
	loginRespRec := httptest.NewRecorder()
	_ = sessionAuth.Login(loginReq, loginRespRec, "demo")

	// the client will make a request with an authenticated session
	req := httptest.NewRequest(http.MethodGet, "/some/page", nil)
	// copy the session cookie from the login Request into the new request
	// normally the browser/client takes care of this
	loginResp := http.Response{Header: loginRespRec.Header()}
	req.Header.Set("Cookie", loginResp.Cookies()[0].String())

	// use the middleware to protect the page
	protectedHandler := sessionAuth.Middleware(protectedSite)

	// check the response
	respRec2 := httptest.NewRecorder()
	protectedHandler.ServeHTTP(respRec2, req)
	resp := respRec2.Result()
	fmt.Println(resp.StatusCode)

	// Output: 200
}
