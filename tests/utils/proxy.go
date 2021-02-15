package utils

import (
	"context"
	"time"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/ethclient"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/log"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/utils"

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
