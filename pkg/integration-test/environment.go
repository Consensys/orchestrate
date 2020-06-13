package integrationtest

import (
	"context"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
)

type TestSuiteEnv interface {
	Start(ctx context.Context) error
}

func StartEnvironment(env TestSuiteEnv) (gerr error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// We do not exit in Fatal errors
	log.StandardLogger().ExitFunc = func(code int) {
		gerr = fmt.Errorf("fatal error has been occurred")
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
