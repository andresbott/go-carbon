package auth

import (
	"encoding/gob"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/sessions"
	"net/http"
	"time"
)

type SessionMgr struct {
	user  UserLogin
	store sessions.Store

	sessionDur    time.Duration
	maxSessionDur time.Duration

	logger func(action int, user string)
}

type SessionCfg struct {
	User  UserLogin
	Store sessions.Store

	SessionDur    time.Duration // normal session duration, can be renewed on subsequent requests
	MinWriteSpace time.Duration // time between the last session update
	MaxSessionDur time.Duration // force a logout after this time

	logger func(action int, user string)
}

func NewSessionAuth(cfg SessionCfg) (*SessionMgr, error) {
	gob.Register(SessionData{})
	if cfg.User == nil {
		return nil, fmt.Errorf("user login cannot be empty")
	}

	if cfg.logger == nil {
		cfg.logger = func(action int, user string) {}
	}
	if cfg.SessionDur == 0 {
		cfg.SessionDur = time.Hour * 1
	}
	if cfg.MaxSessionDur == 0 {
		cfg.MaxSessionDur = time.Hour * 24
	}

	c := SessionMgr{
		user:          cfg.User,
		sessionDur:    cfg.SessionDur,
		maxSessionDur: cfg.MaxSessionDur,
		store:         cfg.Store,
		logger:        cfg.logger,
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
	UserId   string // ID or username
	IsAuth   bool
	DeviceID string // hold information about the device
	// expiration of the session, e.g. 2 days, after a login is required, this value can be updated by "keep me logged in"
	Expiration  time.Time
	ForceReAuth time.Time // force re-auth, max time a session is valid, even if keep logged in is in place.
}

func (d *SessionData) process(extend time.Duration) {
	// check expiration
	if d.Expiration.Before(time.Now()) {
		d.IsAuth = false
	}
	// check hard expiration
	if d.ForceReAuth.Before(time.Now()) {
		d.IsAuth = false
	}
	// extend normal expiration
	if d.IsAuth && extend > 0 {
		d.Expiration = d.Expiration.Add(extend)
	}
}

const (
	SessionName    = "_c_auth"
	sessionDataKey = "data"
)

func (auth *SessionMgr) WriteSession(r *http.Request, w http.ResponseWriter, data SessionData) error {
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
func (auth *SessionMgr) Read(r *http.Request) (SessionData, error) {
	session, err := auth.store.Get(r, SessionName)
	if err != nil {
		return SessionData{}, err
	}
	key := session.Values[sessionDataKey]
	if key == nil {
		return SessionData{}, nil
	}
	authData := key.(SessionData)
	authData.process(auth.sessionDur)

	return authData, nil
}

// FormAuthHandler is a http handler that responds to login requests made from a Form
// it will check if the provided data is correct and write the login cookie into the response.
// todo better err handling
func (auth *SessionMgr) FormAuthHandler() http.Handler {
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
				UserId:      payload.Name,
				IsAuth:      true,
				Expiration:  time.Now().Add(auth.sessionDur),
				ForceReAuth: time.Now().Add(auth.maxSessionDur),
			}
			err = auth.WriteSession(r, w, authData)
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

func (auth *SessionMgr) Middleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// read the cookie

		// check data
		data, err := auth.Read(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		if data.IsAuth {
			auth.logger(ActionLoginOk, data.UserId)

			// todo only write if time is big enough to not overload session write
			err = auth.WriteSession(r, w, data)
			if err != nil {
				spew.Dump(err)
				http.Error(w, "internal error", http.StatusInternalServerError)
				return
			}

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
