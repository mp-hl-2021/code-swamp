package prom

import (
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"time"

	"net/http"
)

var (
	totalHttpRequests = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_total_requests",
		Help: "Handled HTTP requests",
	}, []string{"code", "method"})
	inflightHttpRequests = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "http_inflight_requests",
		Help: "HTTP requests currently inflight",
	})
	durationHttpRequests = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "http_handlers_duration_seconds",
		Help: "code-swamp HTTP requests duration in seconds",
		//Buckets:     nil,
	}, []string{"path"})
)

func Measurer() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return promhttp.InstrumentHandlerCounter(totalHttpRequests, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			now := time.Now()
			inflightHttpRequests.Inc()
			next.ServeHTTP(w, r)
			inflightHttpRequests.Dec()
			route := mux.CurrentRoute(r)
			path, _ := route.GetPathTemplate()

			durationHttpRequests.WithLabelValues(path).Observe(time.Since(now).Seconds())
			// todo: add success label
		}))
	}
}