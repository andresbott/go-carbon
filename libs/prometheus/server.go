package prometheus

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

const DefaultAddr = ":9090"

// NewServer creates a new sever instance that can be started individually
func NewServer(Addr string) *http.Server {

	if Addr == "" {
		Addr = DefaultAddr
	}

	return &http.Server{
		Addr:    Addr,
		Handler: promhttp.Handler(),
	}
}
