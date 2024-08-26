package handlrs

import (
	"encoding/json"
	"github.com/andresbott/go-carbon/libs/auth"
	"net/http"
)

//type UserHandler struct {
//	user    auth.UserLogin
//	session *auth.SessionMgr
//}
//
//
//
//func (h *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
//
//	var payload loginData
//
//	err := json.NewDecoder(r.Body).Decode(&payload)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusBadRequest)
//		return
//	}
//
//	if h.user.AllowLogin(payload.User, payload.Pw) {
//		err = h.session.Login(r, w, payload.User)
//		if err != nil {
//			http.Error(w, "internal error", http.StatusInternalServerError)
//			return
//		}
//	} else {
//		http.Error(w, "Unauthorized", http.StatusUnauthorized)
//		return
//	}
//
//}

type loginData struct {
	User string `json:"user"`
	Pw   string `json:"password"`
}

type userStatus struct {
	User     string `json:"user"`
	LoggedIn bool   `json:"logged-in"`
}

func UserStatusHandler(session *auth.SessionMgr) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		data, err := session.Read(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		jsonData := userStatus{
			User:     data.UserId,
			LoggedIn: data.IsAuthenticated,
		}

		err = json.NewEncoder(w).Encode(jsonData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

	})
}

const anonymousUser = "anonymous"

type autDisabledStatus struct {
	AuthDisabled bool   `json:"auth-disabled,omitempty"`
	User         string `json:"user"`
}

func AuthDisabledHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		jsonData := autDisabledStatus{
			AuthDisabled: true,
			User:         anonymousUser,
		}
		err := json.NewEncoder(w).Encode(jsonData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

	})
}

func UserLogoutHandler(session *auth.SessionMgr) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		err := session.Logout(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		jsonData := userStatus{
			User:     "",
			LoggedIn: false,
		}

		err = json.NewEncoder(w).Encode(jsonData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

	})
}

func UserLoginHandler(session *auth.SessionMgr, user auth.UserLogin) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload loginData
		// TODO remove
		//time.Sleep(1 * time.Second)

		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if user.AllowLogin(payload.User, payload.Pw) {
			err = session.Login(r, w, payload.User)
			if err != nil {
				http.Error(w, "internal error", http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			// todo read user data...
			jsonData := userStatus{
				User:     payload.User,
				LoggedIn: true,
			}
			err = json.NewEncoder(w).Encode(jsonData)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		} else {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
	})
}
