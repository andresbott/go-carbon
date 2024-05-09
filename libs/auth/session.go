package auth

import (
	"github.com/gorilla/sessions"
	"net/http"
	"os"
)

type Session struct {
}

// for now to kee a reference
var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

func MyHandler(w http.ResponseWriter, r *http.Request) {
	// Get a session. We're ignoring the error resulted from decoding an
	// existing session: Get() always returns a session, even if empty.
	session, _ := store.Get(r, "session-name")
	// Set some session values.
	session.Values["foo"] = "bar"
	session.Values[42] = 43
	// Save it before we write to the response/return from the handler.
	err := session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
