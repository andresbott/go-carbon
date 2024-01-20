package prometheus_test

import (
	"fmt"
	prommidleware "git.andresbott.com/Golang/carbon/libs/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"io"

	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func testHandler(statusCode int) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
		fmt.Fprint(w, "ok")
	})
}

func TestMiddleware(t *testing.T) {
	tcs := []struct {
		name          string
		requests      func(h http.Handler)
		cfg           prommidleware.Cfg
		statusCode    int
		expectedLines []string
	}{
		{
			name: "simple test",
			requests: func(h http.Handler) {
				r := httptest.NewRequest("GET", "/bla", nil)
				r2 := httptest.NewRequest("POST", "/ble/bli", nil)
				rec := httptest.NewRecorder()
				h.ServeHTTP(rec, r)
				h.ServeHTTP(rec, r2)
			},
			statusCode: 200,
			expectedLines: []string{
				`http_duration_seconds_bucket{addr="/bla",isError="false",method="GET",status="200",type="HTTP/1.1",le="0.005"} 1`,
				`http_duration_seconds_bucket{addr="/bla",isError="false",method="GET",status="200",type="HTTP/1.1",le="0.01"} 1`,
				`http_duration_seconds_bucket{addr="/ble/bli",isError="false",method="POST",status="200",type="HTTP/1.1",le="0.01"} 1`,
				`http_duration_seconds_bucket{addr="/ble/bli",isError="false",method="POST",status="200",type="HTTP/1.1",le="0.25"} 1`,
			},
		},

		{
			name: "status codes are grouped",
			requests: func(h http.Handler) {
				r := httptest.NewRequest("GET", "/bla", nil)
				rec := httptest.NewRecorder()
				h.ServeHTTP(rec, r)
			},
			cfg: prommidleware.Cfg{
				GroupRespCodes: true,
			},
			statusCode: 200,
			expectedLines: []string{
				`http_duration_seconds_bucket{addr="/bla",isError="false",method="GET",status="2xx",type="HTTP/1.1",le="0.005"} 1`,
				`http_duration_seconds_bucket{addr="/bla",isError="false",method="GET",status="2xx",type="HTTP/1.1",le="0.01"} 1`,
			},
		},
		{
			name: "metric prefix",
			requests: func(h http.Handler) {
				r := httptest.NewRequest("GET", "/bla", nil)
				rec := httptest.NewRecorder()
				h.ServeHTTP(rec, r)
			},
			cfg: prommidleware.Cfg{
				MetricPrefix: "ehmm",
			},
			statusCode: 200,
			expectedLines: []string{
				`ehmm_http_duration_seconds_bucket{addr="/bla",isError="false",method="GET",status="200",type="HTTP/1.1",le="0.005"} 1`,
				`ehmm_http_duration_seconds_bucket{addr="/bla",isError="false",method="GET",status="200",type="HTTP/1.1",le="0.01"} 1`,
			},
		},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {

			reg := prometheus.NewRegistry()
			cfg := tc.cfg
			cfg.Registry = reg

			mid := prommidleware.NewMiddleware(cfg)

			h := mid.Handler(testHandler(tc.statusCode))

			tc.requests(h)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/metrics", nil)

			promhttp.HandlerFor(reg, promhttp.HandlerOpts{}).ServeHTTP(rec, req)
			resp := rec.Result()

			body, _ := io.ReadAll(resp.Body)

			respBody := string(body)

			//fmt.Print(respBody)
			for _, line := range tc.expectedLines {
				if !strings.Contains(respBody, line) {
					t.Errorf("response does not contains expected line: %s", line)
				}
			}

		})
	}
}
