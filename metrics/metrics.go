package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	otel_prom "go.opentelemetry.io/otel/exporters/prometheus"
	api "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric"
)

type MetricsOpts struct {
	Name        string
	Description string
}

type Registry interface {
	GetCounter(MetricsOpts) (api.Float64Counter, error)
}

type promRegistry struct {
	metrics  map[string]api.Float64Counter
	provider *metric.MeterProvider
}

func NewRegistry() (Registry, error) {
	prom, err := otel_prom.New()
	if err != nil {
		return nil, err
	}
	return &promRegistry{
		metrics:  make(map[string]api.Float64Counter),
		provider: metric.NewMeterProvider(metric.WithReader(prom)),
	}, nil
}

func (r *promRegistry) GetCounter(opts MetricsOpts) (api.Float64Counter, error) {
	val, found := r.metrics[opts.Name]
	if found {
		return val, nil
	}

	meter := r.provider.Meter("optimizely-go-sdk")

	return meter.Float64Counter(opts.Name, api.WithDescription(opts.Description))
}

func GetPrometheusHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		promhttp.Handler().ServeHTTP(w, r)
	}
}
