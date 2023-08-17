package client

import (
	"context"
	"fmt"
	"time"

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
	start := time.Now().Unix()
	if c.metricsRegistry != nil {
		cv, err := c.metricsRegistry.GetFloat64Counter(metrics.MetricsOpts{Name: "activate_hit_count", Description: "Number of times the activate called"})
		if err != nil {
			return err
		}
		cv.Add(context.Background(), 1)

		hs, err := c.metricsRegistry.GetFloat64Histogram(metrics.MetricsOpts{Name: "activate_response_time", Description: "Histogram for the activate response time"})
		if err != nil {
			return err
		}

		defer func() {
			end := time.Now().Unix()
			hs.Record(context.Background(), float64(end-start+1))
		}()
	}
	fmt.Println("Mock activate")
	return nil
}
