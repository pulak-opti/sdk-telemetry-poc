package metrics

import (
	"net/http"
	"sync"

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
	GetFloat64Counter(MetricsOpts) (api.Float64Counter, error)
	GetFloat64Histogram(MetricsOpts) (api.Float64Histogram, error)
}

func GetPrometheusHandler() http.Handler {
	return promhttp.Handler()
}

type promRegistry struct {
	counters  map[string]api.Float64Counter
	histogram map[string]api.Float64Histogram
	provider  *metric.MeterProvider
	rLock     *sync.RWMutex
}

func NewRegistry() (Registry, error) {
	prom, err := otel_prom.New()
	if err != nil {
		return nil, err
	}
	return &promRegistry{
		counters:  make(map[string]api.Float64Counter),
		histogram: make(map[string]api.Float64Histogram),
		provider:  metric.NewMeterProvider(metric.WithReader(prom)),
		rLock:     &sync.RWMutex{},
	}, nil
}

func (r *promRegistry) GetFloat64Counter(opts MetricsOpts) (api.Float64Counter, error) {
	r.rLock.RLock()
	defer r.rLock.RUnlock()

	val, found := r.counters[opts.Name]
	if found {
		return val, nil
	}

	meter := r.provider.Meter("optimizely-go-sdk")

	counter, err := meter.Float64Counter(opts.Name, api.WithDescription(opts.Description))
	if err != nil {
		return nil, err
	}
	r.counters[opts.Name] = counter
	return counter, nil
}

func (r *promRegistry) GetFloat64Histogram(opts MetricsOpts) (api.Float64Histogram, error) {
	r.rLock.RLock()
	defer r.rLock.RUnlock()

	val, found := r.histogram[opts.Name]
	if found {
		return val, nil
	}

	meter := r.provider.Meter("optimizely-go-sdk")

	histogram, err := meter.Float64Histogram(opts.Name)
	if err != nil {
		return nil, err
	}

	r.histogram[opts.Name] = histogram
	return histogram, nil
}
