package client

import (
	"context"
	"fmt"

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
		cv, err := c.metricsRegistry.GetCounter(metrics.MetricsOpts{Name: "activate_hit_count", Description: "Number of times the activate called"})
		if err != nil {
			return err
		}
		cv.Add(context.Background(), 1)
	}
	fmt.Println("Mock activate")
	return nil
}
