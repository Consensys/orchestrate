package utils

import (
	"context"
	"time"

	"github.com/ConsenSys/orchestrate/pkg/toolkit/app/log"
	"github.com/ConsenSys/orchestrate/pkg/toolkit/ethclient"
	"github.com/ConsenSys/orchestrate/pkg/utils"

	"github.com/cenkalti/backoff/v4"
)

func WaitForProxy(ctx context.Context, proxyHost, chainUUID string, ec ethclient.ChainSyncReader) error {
	chainProxyURL := utils.GetProxyURL(proxyHost, chainUUID)
	return backoff.RetryNotify(
		func() error {
			_, err2 := ec.Network(ctx, chainProxyURL)
			return err2
		},
		backoff.WithMaxRetries(backoff.NewConstantBackOff(time.Second), 5),
		func(err error, duration time.Duration) {
			log.FromContext(ctx).WithField("chain", chainUUID).WithError(err).Debug("chain proxy is still not ready")
		},
	)
}
