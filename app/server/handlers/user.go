package handlers

import (
	"encoding/json"
	"git.andresbott.com/Golang/carbon/libs/auth"
	"net/http"
)

type loginData struct {
	User string `json:"user"`
	Pw   string `json:"password"`
}

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

func UserLoginHandler(session *auth.SessionMgr, user auth.UserLogin) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload loginData

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
		} else {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
	})
}
