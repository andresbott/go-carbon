package auth

import (
	"fmt"
	"net/http"
)

// CtxUserInfo extracts the user information from a request context
func CtxUserInfo(r *http.Request) (UserData, error) {

	var d UserData
	ctx := r.Context()

	val := ctx.Value(UserIdKey)
	userId, ok := val.(string)
	if !ok {
		return d, fmt.Errorf("unable to cast userId to string")
	}

	val = ctx.Value(UserIsLoggedInKey)
	isLoggedIn, ok := val.(bool)
	if !ok {
		return d, fmt.Errorf("unable to cast isLoggedIn to boolean")
	}

	d = UserData{
		UserId:          userId,
		IsAuthenticated: isLoggedIn,
	}
	return d, nil

}
