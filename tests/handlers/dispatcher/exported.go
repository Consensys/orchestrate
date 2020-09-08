package dispatcher

import (
	"context"
	"fmt"
	"sync"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/tests/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/tests/utils/chanregistry"
)

var (
	handler    engine.HandlerFunc
	initOnce   = &sync.Once{}
	keyOfFuncs []KeyOfFunc
)

func SetKeyOfFuncs(keyOfs ...KeyOfFunc) {
	keyOfFuncs = keyOfs
}

// Init initialize Dispatcher Handler
func Init(ctx context.Context) {
	initOnce.Do(func() {
		if handler != nil {
			return
		}

		// Initialize Channel registry
		chanregistry.Init(ctx)

		handler = Dispatcher(chanregistry.GlobalChanRegistry(), keyOfFuncs...)

		log.Infof("dispatcher: handler ready")
	})
}

// SetGlobalHandler sets global Cucumber Handler
func SetGlobalHandler(h engine.HandlerFunc) {
	handler = h
}

// GlobalHandler returns global Cucumber handler
func GlobalHandler() engine.HandlerFunc {
	return handler
}

func LongKeyOf(topics map[string]string) KeyOfFunc {
	return func(txctx *engine.TxContext) (string, error) {
		topic, ok := topics[txctx.In.Entrypoint()]
		if !ok {
			return "", fmt.Errorf("unknown message entrypoint")
		}

		return utils.LongKeyOf(topic, txctx.Envelope.GetID()), nil
	}
}

func LabelKey(topics map[string]string) KeyOfFunc {
	return func(txctx *engine.TxContext) (string, error) {
		topic, ok := topics[txctx.In.Entrypoint()]
		if !ok {
			return "", fmt.Errorf("unknown message entrypoint")
		}

		id := txctx.Envelope.GetContextLabelsValue("id")
		if id == "" {
			return "", fmt.Errorf("message has no id in context labels")
		}
		return utils.LongKeyOf(topic, id), nil
	}
}

func ShortKeyOf(topics map[string]string) KeyOfFunc {
	return func(txtcx *engine.TxContext) (string, error) {
		topic, ok := topics[txtcx.In.Entrypoint()]
		if !ok {
			return "", fmt.Errorf("unknown message entrypoint")
		}

		scenario := txtcx.Envelope.GetContextLabelsValue("scenario.id")
		if scenario == "" {
			return "", fmt.Errorf("message has no test scenario")
		}

		return utils.ShortKeyOf(topic, scenario), nil
	}
}
