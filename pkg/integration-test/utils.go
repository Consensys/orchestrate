package integrationtest

import (
	"context"
	"net/http"
	"time"

	"github.com/containous/traefik/v2/pkg/log"
)

func WaitForServiceLive(ctx context.Context, url, name string, timeout time.Duration) {
	logger := log.FromContext(ctx)
	rctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	for {
		req, _ := http.NewRequest("GET", url, nil)
		req = req.WithContext(rctx)

		resp, err := http.DefaultClient.Do(req)
		if err == nil {
			if resp != nil && resp.StatusCode == 200 {
				logger.Infof("service %s is live", name)
				return
			}

			logger.WithField("status", resp.StatusCode).Warnf("cannot reach %s service", name)
		}

		if rctx.Err() != nil {
			logger.WithError(rctx.Err()).Warnf("cannot reach %s service", name)
			return
		}

		logger.Debugf("waiting for 1 s for service %s to start...", name)
		time.Sleep(time.Second)
	}
}
