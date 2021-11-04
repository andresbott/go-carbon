package server

import (
	"fmt"
	"git.andresbott.com/Golang/carbon/libs/log"
	"net/http"
	"strings"
)

type textHandler struct {
	Text   string
	Links  map[string]string
	Logger log.LeveledLogger
}

func (h *textHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	h.Logger.Infof("Request", "method", r.Method, "url", r.RequestURI)

	w.Header().Add("Content-Type", " text/html")
	if r.Method == http.MethodGet {

		var s strings.Builder

		s.WriteString("GET: " + h.Text)

		if len(h.Links) > 0 {
			s.WriteString("<ul>")
			for name, url := range h.Links {
				s.WriteString(fmt.Sprintf("<li><a href=\"%s\">%s</a></li>", url, name))
			}
			s.WriteString("</ul>")
		}

		fmt.Fprintf(w, s.String())

		return
	}

	if r.Method == http.MethodPost {
		fmt.Fprintf(w, "POST: %s", h.Text)
		return
	}

	http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
}
