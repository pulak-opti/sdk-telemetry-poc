package client

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/pulak-opti/sdk-telemetry-poc/metrics"
)

type OptiClient struct {
	metricsRegistry metrics.Registry
}

func NewOptiClient(reg metrics.Registry) *OptiClient {
	return &OptiClient{
		metricsRegistry: reg,
	}
}

func (c *OptiClient) Activate() error {
	if c.metricsRegistry != nil {
		cv, err := c.metricsRegistry.GetCounter(prometheus.CounterOpts{Name: "activate_hit_count"})
		if err != nil {
			return err
		}
		cv.Add(1)
	}
	fmt.Println("Mock activate")
	return nil
}
