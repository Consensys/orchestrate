package integrationtest

import (
	"context"
	"fmt"
	"os"

	"github.com/ConsenSys/orchestrate/pkg/utils"
	log "github.com/sirupsen/logrus"
)

type TestSuiteEnv interface {
	Start(ctx context.Context) error
}

func StartEnvironment(ctx context.Context, env TestSuiteEnv) (gerr error) {
	ctx, cancel := context.WithCancel(ctx)

	// We do not exit in Fatal errors
	log.StandardLogger().ExitFunc = func(code int) {
		gerr = fmt.Errorf("fatal error has occurred")
		cancel()
	}
	defer func() {
		log.StandardLogger().ExitFunc = nil
	}()

	sig := utils.NewSignalListener(func(signal os.Signal) {
		gerr = fmt.Errorf("interrupt signal has been sent")
		cancel()
	})
	defer sig.Close()

	err := env.Start(ctx)
	if err != nil {
		if gerr == nil {
			return err
		}
	}

	return
}
