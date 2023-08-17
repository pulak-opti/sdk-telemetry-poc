package main

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/pulak-opti/sdk-telemetry-poc/client"
	"github.com/pulak-opti/sdk-telemetry-poc/metrics"
)

var metricsReg metrics.Registry

func activate(w http.ResponseWriter, r *http.Request) {
	optiClient := client.NewOptiClient(metricsReg)
	if err := optiClient.Activate(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte("OK"))
}

func main() {
	r := chi.NewRouter()
	var err error
	metricsReg, err = metrics.NewRegistry()
	if err != nil {
		panic(err)
	}
	r.Get("/metrics", metrics.GetPrometheusHandler().ServeHTTP)
	r.Get("/activate", activate)

	if err := http.ListenAndServe(":8080", r); err != nil {
		slog.Error(err.Error())
	}
}
