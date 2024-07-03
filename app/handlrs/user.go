package handlrs

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

type userStatus struct {
	User     string `json:"user"`
	LoggedIn bool   `json:"logged-in"`
}

func UserStatusHandler(session *auth.SessionMgr) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// TODO remove this
		w.Header().Set("Access-Control-Allow-Origin", "*")

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

func UserLoginHandler(session *auth.SessionMgr, user auth.UserLogin) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload loginData
		// TODO move to another place
		w.Header().Set("Access-Control-Allow-Origin", "*")

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
