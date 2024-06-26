package handlers

const sessionLogin = "/session-login"
const sessionContent = "/session"

//func sessionAuthentication(r *mux.Router, demoUsers user.StaticUsers) error {
//	// session based auth
//	// --------------------------
//	hashKey := []byte("oach9iu2uavahcheephi4FahzaeNge8yeecie4jee9rah9ahrah6tithai7Oow5U")
//	blockKey := []byte("eeth3oon5eewifaogeibieShey5eiJ0E")
//
//	//cookieStore, err := auth.CookieStore(hashKey, blockKey)
//	cookieStore, err := auth.FsStore("", hashKey, blockKey)
//	if err != nil {
//		return err
//	}
//
//	sessionAuth, err := auth.NewSessionMgr(auth.SessionCfg{
//		Store: cookieStore,
//	})
//	if err != nil {
//		return err
//	}
//
//	pageHandler := handlers.SimpleText{
//		Text: "Page protected by session auth",
//		Links: []handlers.Link{
//			{Text: "back to root", Url: "../"},
//		},
//	}
//
//	ProtectedPage := sessionAuth.Middleware(&pageHandler)
//	r.Path(sessionContent).Handler(ProtectedPage)
//	// handle the post request
//	loginHandler := auth.FormAuthHandler(sessionAuth, demoUsers)
//	r.PathPrefix(sessionLogin).Methods(http.MethodPost).Handler(loginHandler)
//
//	// render the form
//	loginFormHandlr := handlers.TemplateWithRequest(server.loginForm, func(r *http.Request) map[string]interface{} {
//		payload := map[string]interface{}{}
//		payload["Path"] = sessionLogin
//		payload["Redirect"] = sessionContent
//		return payload
//	})
//	r.PathPrefix(sessionLogin).HandlerFunc(loginFormHandlr)
//	return nil
//}
