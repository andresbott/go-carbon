package auth

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/schema"
	"github.com/gorilla/securecookie"
	"net/http"
	"time"
)

const DefaultCookieAuth = "auth-login"

// Cookie auth manager uses an encrypted cookie to store the login status of the users
// this is less secure than using a serve side sessions management but more scalable
// we relly on gorilla securecookie for the crypto part
type Cookie struct {
	user         UserLogin
	cookieName   string
	expire       time.Duration
	redirect     string
	redirectCode int
	secCookie    *securecookie.SecureCookie
	logger       func(action int, user string)
}

type CookieCfg struct {
	User         UserLogin
	Redirect     string
	RedirectCode int

	CookieName string
	HashKey    []byte
	BlockKey   []byte

	Expire time.Duration
	logger func(action int, user string)
}

// NewCookieAuth returns a new instance of a cookie auth manager
// note the HashKey and BlockKey are passed directly into gorilla secure cookie
// https://github.com/gorilla/securecookie
func NewCookieAuth(cfg CookieCfg) (*Cookie, error) {

	if cfg.User == nil {
		return nil, fmt.Errorf("user login cannot be empty")
	}
	if cfg.CookieName == "" {
		cfg.CookieName = DefaultCookieAuth
	}
	if cfg.Expire == 0 {
		cfg.Expire = 6 * time.Hour
	}
	// TODO probably it would be good to add some validation for the crypto keys are correctly formatted etc
	c := Cookie{
		user:         cfg.User,
		cookieName:   cfg.CookieName,
		redirect:     cfg.Redirect,
		redirectCode: cfg.RedirectCode,
		logger:       cfg.logger,
		expire:       cfg.Expire,
		secCookie:    securecookie.New(cfg.HashKey, cfg.BlockKey),
	}
	return &c, nil
}

type CookieData struct {
	UserId string // ID or username
	IsAuth bool
}

// WriteAuthCookie writes the auth cookie into a http response
func (auth *Cookie) WriteAuthCookie(data CookieData, w http.ResponseWriter) error {

	encoded, err := auth.secCookie.Encode(auth.cookieName, data)
	if err != nil {
		return fmt.Errorf("unable to encode cokkie: %s", err.Error())
	}

	expire := time.Now().Add(auth.expire)
	spew.Dump(auth.expire)

	cookie := &http.Cookie{
		Name:     auth.cookieName,
		Value:    encoded,
		Path:     "/",
		Expires:  expire,
		Secure:   true,
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
	return nil

}

type LoginFormData struct {
	Name     string
	Pw       string
	Redirect string
}

// Set a Decoder instance as a package global, because it caches
// meta-data about structs, and an instance can be shared safely.
var formDecoder = schema.NewDecoder()

// FormAuthHandler is a http handler that responds to login requests made from a Form
// it will check if the provided data is correct and write the login cookie into the response.
// todo better err handling
func (auth *Cookie) FormAuthHandler() http.Handler {
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

			authData := CookieData{
				UserId: payload.Name,
				IsAuth: true,
			}
			err := auth.WriteAuthCookie(authData, w)
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

// Read will verify if user is logged in based on the secure cookie data
func (auth *Cookie) Read(r *http.Request) (CookieData, error) {
	value := CookieData{}
	if cookie, err := r.Cookie(auth.cookieName); err == nil {
		if err = auth.secCookie.Decode(auth.cookieName, cookie.Value, &value); err != nil {
			return value, err
		}
	}
	return value, nil
}

func (auth *Cookie) Middleware(next http.Handler) http.Handler {
	if auth.logger == nil {
		auth.logger = func(action int, user string) {}
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// check cookie
		data, err := auth.Read(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		if data.IsAuth {
			auth.logger(ActionLoginOk, data.UserId)
			next.ServeHTTP(w, r)
			return
		}

		if auth.redirect != "" {
			http.Redirect(w, r, auth.redirect, auth.redirectCode)
			return
		}
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
}
