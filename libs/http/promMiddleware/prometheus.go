package promMiddleware

import (
	"fmt"
	"git.andresbott.com/Golang/carbon/libs/http/responseStatusCode"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"strconv"
	"time"
)

type Middleware struct {
	histogram      *prometheus.HistogramVec
	groupRespCodes bool
}

type Cfg struct {
	MetricPrefix   string
	Buckets        []float64
	Registry       prometheus.Registerer
	GroupRespCodes bool
}

func New(cfg Cfg) *Middleware {
	// if null does it use the same as promauto?
	if cfg.Registry == nil {
		cfg.Registry = prometheus.DefaultRegisterer
	}

	if len(cfg.Buckets) == 0 {
		cfg.Buckets = prometheus.DefBuckets
	}

	m := &Middleware{
		histogram: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: cfg.MetricPrefix,
			Subsystem: "http",
			Name:      "duration_seconds",
			Help:      "Duration of HTTP requests for different paths, methods, status codes",
			Buckets:   cfg.Buckets,
		},
			[]string{
				"type",
				"status",
				"method",
				"addr",
				"isError",
			},
		),
		groupRespCodes: cfg.GroupRespCodes,
	}

	cfg.Registry.MustRegister(
		m.histogram,
	)
	return m
}

func (m *Middleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		respWriter := responseStatusCode.New(w)
		timeStart := time.Now()
		next.ServeHTTP(respWriter, r)

		var statusCodeStr string
		if m.groupRespCodes {
			statusCodeStr = fmt.Sprintf("%dxx", respWriter.StatusCode()/100)
		} else {
			statusCodeStr = respWriter.StatusCodeStr()
		}

		isErrorStr := strconv.FormatBool(responseStatusCode.IsStatusError(respWriter.StatusCode()))

		m.histogram.With(prometheus.Labels{
			"type":    r.Proto,
			"status":  statusCodeStr,
			"method":  r.Method,
			"addr":    r.URL.Path,
			"isError": isErrorStr,
		}).Observe(time.Since(timeStart).Seconds())
	})
}
