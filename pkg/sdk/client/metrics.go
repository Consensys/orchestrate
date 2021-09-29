package client

import (
	"context"
	"fmt"
	"time"

	clientutils "github.com/consensys/orchestrate/pkg/toolkit/app/http/client-utils"
	promcli "github.com/consensys/orchestrate/pkg/toolkit/app/metrics/client"
	healthz "github.com/heptiolabs/healthcheck"
	dto "github.com/prometheus/client_model/go"
)

func (c *HTTPClient) Checker() healthz.Check {
	return healthz.HTTPGetCheck(fmt.Sprintf("%s/live", c.config.MetricsURL), time.Second)
}

func (c *HTTPClient) Prometheus(ctx context.Context) (map[string]*dto.MetricFamily, error) {
	resp, err := clientutils.GetRequest(ctx, c.client, fmt.Sprintf("%s/metrics", c.config.MetricsURL))
	if err != nil {
		return nil, err
	}

	mf, err := promcli.ParseResponse(resp)
	if err != nil {
		return nil, err
	}

	return mf, nil
}
