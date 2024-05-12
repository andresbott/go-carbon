package auth

import (
	"encoding/gob"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/sessions"
	"net/http"
	"time"
)

type Session struct {
	user   UserLogin
	store  sessions.Store
	logger func(action int, user string)
}

type SessionCfg struct {
	User   UserLogin
	Store  sessions.Store
	logger func(action int, user string)
}

func NewSessionAuth(cfg SessionCfg) (*Session, error) {
	gob.Register(SessionData{})
	if cfg.User == nil {
		return nil, fmt.Errorf("user login cannot be empty")
	}

	if cfg.logger == nil {
		cfg.logger = func(action int, user string) {}
	}

	c := Session{
		user:   cfg.User,
		store:  cfg.Store,
		logger: cfg.logger,
	}
	return &c, nil
}

// CookieStore is a convenience function to generate a new secure cookiestore
// based on the securecookie.New doc:
//
// hashKey is required, used to authenticate values using HMAC. Create it using
// GenerateRandomKey(). It is recommended to use a key with 32 or 64 bytes.
//
// blockKey is optional, used to encrypt values. Create it using
// GenerateRandomKey(). The key length must correspond to the key size
// of the encryption algorithm. For AES, used by default, valid lengths are
// 16, 24, or 32 bytes to select AES-128, AES-192, or AES-256.
// The default encoder used for cookie serialization is encoding/gob.
func CookieStore(HashKey, BlockKey []byte) (*sessions.CookieStore, error) {
	hashL := len(HashKey)
	if hashL != 32 && hashL != 64 {
		return nil, fmt.Errorf("HashKey lenght should be 32 or 64 bytes")
	}
	blockKeyL := len(BlockKey)
	if blockKeyL != 16 && blockKeyL != 24 && blockKeyL != 32 {
		return nil, fmt.Errorf("blockKey lenght should be 16, 24 or 32 bytes")
	}
	return sessions.NewCookieStore(HashKey, BlockKey), nil
}

type SessionData struct {
	UserId string // ID or username
	IsAuth bool
	// expiration of the session, e.g. 2 days, after a login is required, this value can be updated by "keep me logged in"
	Expiration  time.Time
	ForceReAuth time.Time // force re-auth, max time a session is valid, even if keep logged in is in place.
}

const (
	SessionName    = "_c_auth"
	sessionDataKey = "data"
)

func (auth *Session) WriteSessionData(r *http.Request, w http.ResponseWriter, data SessionData) error {
	session, err := auth.store.Get(r, SessionName)
	if err != nil {
		return err
	}
	session.Values[sessionDataKey] = data
	err = session.Save(r, w)
	if err != nil {
		return err
	}
	return nil
}

// FormAuthHandler is a http handler that responds to login requests made from a Form
// it will check if the provided data is correct and write the login cookie into the response.
// todo better err handling
func (auth *Session) FormAuthHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			spew.Dump(err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		var payload LoginFormData

		// r.PostForm is a map of our POST form values
		err = formDecoder.Decode(&payload, r.PostForm)
		if err != nil {
			spew.Dump(err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		auth.logger(ActionLoginCheck, payload.Name)

		if auth.user.AllowLogin(payload.Name, payload.Pw) {
			auth.logger(ActionLoginOk, payload.Name)

			authData := SessionData{
				UserId: payload.Name,
				IsAuth: true,
			}
			err = auth.WriteSessionData(r, w, authData)
			if err != nil {
				spew.Dump(err)
				http.Error(w, "internal error", http.StatusInternalServerError)
				return
			}
		} else {
			auth.logger(ActionLoginFailed, payload.Name)
			http.Error(w, "401 unauthorized", http.StatusUnauthorized)
			return
		}

		// todo send path from get request
		http.Redirect(w, r, payload.Redirect, http.StatusSeeOther)

	})
}

func (auth *Session) Read(r *http.Request) (SessionData, error) {
	session, err := auth.store.Get(r, SessionName)
	if err != nil {
		return SessionData{}, err
	}
	key := session.Values[sessionDataKey]
	if key == nil {
		return SessionData{}, nil
	}
	authData := key.(SessionData)
	return authData, nil
	// TODO add additioanl login validation here, e.g. check that the login did not expire
	// TODO, see how we can as well extend login session validity, e.g. add X time before you need to re-auth

}

func (auth *Session) Middleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// read the cookie

		// check data
		data, err := auth.Read(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		if data.IsAuth {
			auth.logger(ActionLoginOk, data.UserId)
			next.ServeHTTP(w, r)
			return
		}

		//if auth.redirect != "" {
		//	http.Redirect(w, r, auth.redirect, auth.redirectCode)
		//	return
		//}
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
}

// for now to kee a reference
