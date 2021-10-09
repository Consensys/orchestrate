package tcp

import (
	"context"
	"io"

	"github.com/traefik/traefik/v2/pkg/log"
)

type Shutdownable interface {
	Shutdown(ctx context.Context) error
}

func Shutdown(ctx context.Context, s Shutdownable) error {
	logger := log.FromContext(ctx)

	err := s.Shutdown(ctx)
	if err == nil {
		return nil
	}

	if ctx.Err() == context.DeadlineExceeded {
		logger.WithError(err).Debugf("shutdown exceeded deadline")
		return ctx.Err()
	}

	logger.WithError(err).Error("failed to shut down")

	return err
}

func Close(s io.Closer) error {
	err := s.Close()
	if err != nil {
		log.WithoutContext().WithError(err).Errorf("failed to close")
		return err
	}
	return nil
}
