package handlers

import (
	"git.andresbott.com/Golang/carbon/libs/http/handlers"
	"git.andresbott.com/Golang/carbon/libs/prometheus"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
	"net/http"
)

type MyAppHandler struct {
	router *mux.Router
}

func (h *MyAppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}

// NewAppHandler generates the main url router handler to be used in the server
func NewAppHandler(l *zerolog.Logger, db *gorm.DB) *MyAppHandler {

	r := mux.NewRouter()

	//// add logging middleware
	//r.Use(func(handler http.Handler) http.Handler {
	//	return log.LoggingMiddleware(handler, l)
	//})

	promMiddle := prometheus.NewMiddleware(prometheus.Cfg{
		MetricPrefix: "myApp",
	})
	r.Use(func(handler http.Handler) http.Handler {
		return promMiddle.Handler(handler)
	})

	// root page
	// --------------------------
	rootPage := handlers.SimpleText{
		Text: "root page",
		Links: []handlers.Link{
			{
				Text: "Basic auth protected",
				Url:  "/basic",
			},
			{
				Text: "User handling",
				Url:  "/user",
			},
		},
	}

	r.Path("/").Handler(&rootPage)

	return &MyAppHandler{
		router: r,
	}
}
