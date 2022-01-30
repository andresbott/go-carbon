package simpleTextHandler

import (
	"fmt"
	"net/http"
	"strings"
)

// Handler Is a simple handler that will print a list of navigation destination based on the map passed upon creation.
type Handler struct {
	Text  string
	Links []Link
}
type Link struct {
	Text string
	Url  string
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	w.Header().Add("Content-Type", " text/html")
	if r.Method == http.MethodGet {

		var s strings.Builder

		s.WriteString("GET: " + h.Text)

		if len(h.Links) > 0 {
			s.WriteString("<ul>")
			for _, link := range h.Links {
				s.WriteString(fmt.Sprintf("<li><a href=\"%s\">%s</a></li>", link.Url, link.Text))
			}
			s.WriteString("</ul>")
		}

		fmt.Fprint(w, s.String())

		return
	}

	if r.Method == http.MethodPost {
		fmt.Fprintf(w, "POST: %s", h.Text)
		return
	}

	http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
}
