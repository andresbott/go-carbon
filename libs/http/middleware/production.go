package middleware

import (
	"github.com/rs/zerolog"
	"net/http"
	"time"
)

type ProdReady struct {
	hist Histogram
	log  *zerolog.Logger
}

func NewProd(log *zerolog.Logger) *ProdReady {
	h := NewHistogram("", nil, nil)

	return &ProdReady{
		hist: h,
		log:  log,
	}
}
func (m *ProdReady) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		timeStart := time.Now()
		respWriter := NewWriter(w)
		// serve the request
		next.ServeHTTP(respWriter, r)
		// get the duration
		timeDiff := time.Since(timeStart)
		// log the request
		log(m.log, r, respWriter.StatusCode(), timeDiff)
		// add prometheus metric
		observe(m.hist, r, respWriter.StatusCode(), timeDiff)
		return
	})
}
