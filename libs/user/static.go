package user

type StaticUsers struct {
	Users map[string]string
}

func (fu StaticUsers) AllowLogin(user string, hash string) bool {
	for u, p := range fu.Users {
		if user == u {
			if p == hash {
				return true
			} else {
				return false
			}
		}
	}
	return false
}
