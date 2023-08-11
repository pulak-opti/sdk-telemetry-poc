package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Registry interface {
	GetCounter(prometheus.CounterOpts) (prometheus.Counter, error)
}

type promRegistry struct {
	metrics map[string]prometheus.Counter
}

func NewRegistry() Registry {
	return &promRegistry{
		metrics: make(map[string]prometheus.Counter),
	}
}

func (r *promRegistry) GetCounter(opts prometheus.CounterOpts) (prometheus.Counter, error) {
	val, found := r.metrics[opts.Name]
	if found {
		return val, nil
	}
	cv := prometheus.NewCounter(opts)
	if err := prometheus.Register(cv); err != nil {
		return nil, err
	}
	r.metrics[opts.Name] = cv
	return cv, nil
}
