package client

import (
	"context"
	"fmt"
	"time"

	"github.com/containous/traefik/v2/pkg/log"
	healthz "github.com/heptiolabs/healthcheck"
	dto "github.com/prometheus/client_model/go"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	clientutils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/http/client-utils"
	promcli "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/metrics/client"
)

func (c *HTTPClient) Checker() healthz.Check {
	return healthz.HTTPGetCheck(fmt.Sprintf("%s/live", c.config.MetricsURL), time.Second)
}

func (c *HTTPClient) Prometheus(ctx context.Context) (map[string]*dto.MetricFamily, error) {
	resp, err := clientutils.GetRequest(ctx, c.client, fmt.Sprintf("%s/metrics", c.config.MetricsURL))
	if err != nil {
		errMessage := "error while getting prometheus metrics"
		log.FromContext(ctx).WithError(err).Error(errMessage)
		return nil, errors.ServiceConnectionError(errMessage).ExtendComponent(component)
	}

	mf, err := promcli.ParseResponse(resp)
	if err != nil {
		errMessage := "error while parsing prometheus metric response"
		log.FromContext(ctx).WithError(err).Error(errMessage)
		return nil, errors.ServiceConnectionError(errMessage).ExtendComponent(component)
	}

	return mf, nil
}
