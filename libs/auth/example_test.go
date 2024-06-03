package auth_test

import (
	"fmt"
	"git.andresbott.com/Golang/carbon/libs/auth"
	"net/http"
	"net/http/httptest"
)

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
