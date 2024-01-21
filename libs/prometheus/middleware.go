package prometheus

import (
	"fmt"
	"git.andresbott.com/Golang/carbon/libs/http/extra/respstatus"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"strconv"
	"time"
)

type Cfg struct {
	MetricPrefix   string
	Buckets        []float64
	Registry       prometheus.Registerer
	GroupRespCodes bool
}

func NewMiddleware(cfg Cfg) *Middleware {
	// if null does it use the same as promauto?
	if cfg.Registry == nil {
		cfg.Registry = prometheus.DefaultRegisterer
	}

	if len(cfg.Buckets) == 0 {
		cfg.Buckets = prometheus.DefBuckets
	}

	if cfg.MetricPrefix == "" {
		panic("prometheus middleware metric prefix cannot be empty")
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

type Middleware struct {
	histogram      *prometheus.HistogramVec
	groupRespCodes bool
}

func (m *Middleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		respWriter := respstatus.NewWriter(w)
		timeStart := time.Now()
		next.ServeHTTP(respWriter, r)

		var statusCodeStr string
		if m.groupRespCodes {
			statusCodeStr = fmt.Sprintf("%dxx", respWriter.StatusCode()/100)
		} else {
			statusCodeStr = respWriter.StatusCodeStr()
		}

		isErrorStr := strconv.FormatBool(respstatus.IsStatusError(respWriter.StatusCode()))

		m.histogram.With(prometheus.Labels{
			"type":    r.Proto,
			"status":  statusCodeStr,
			"method":  r.Method,
			"addr":    r.URL.Path,
			"isError": isErrorStr,
		}).Observe(time.Since(timeStart).Seconds())
	})
}
