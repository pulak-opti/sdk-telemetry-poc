package main

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var cv prometheus.Counter

func init() {
	cv = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "activate_http_count",
	})
	prometheus.MustRegister(cv)
}

func activate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func addMetrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cv.Add(1)
		next.ServeHTTP(w, r)
	})
}

func main() {
	r := chi.NewRouter()
	r.Get("/metrics", promhttp.Handler().ServeHTTP)
	r.With(addMetrics).Get("/activate", activate)

	if err := http.ListenAndServe(":8080", r); err != nil {
		slog.Error(err.Error())
	}
}
