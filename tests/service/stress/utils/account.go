package utils

import (
	"context"
	"time"

	"github.com/Shopify/sarama"
	"github.com/containous/traefik/v2/pkg/log"
	"github.com/gofrs/uuid"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/tx"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/tests/utils/chanregistry"
)

const NAccounts = 10

type ctxKey string

var accountCtxKey ctxKey = "accounts"

func CreateNewAccount(ctx context.Context, chanReg *chanregistry.ChanRegistry, producer sarama.SyncProducer) (string, error) {
	log.FromContext(ctx).Debugf("Registering new account...")
	evlp := tx.NewEnvelope()
	msgID := uuid.Must(uuid.NewV4()).String()
	_ = evlp.SetID(msgID)

	t := NewEnvelopeTracker(chanReg, evlp, msgID)
	err := SendEnvelope(producer, "account.generator", evlp, msgID)
	if err != nil {
		return "", nil
	}

	err = t.Load("account.generated", time.Second*10)
	if err != nil {
		return "", nil
	}

	if t.Current.MustGetFromAddress().String() == "" {
		return "", errors.DataError("account was not generated")
	}

	log.FromContext(ctx).Infof("New account generated: %s", t.Current.MustGetFromAddress().String())
	return t.Current.MustGetFromAddress().String(), nil
}

func ContextWithAccounts(ctx context.Context, accounts []string) context.Context {
	return context.WithValue(ctx, accountCtxKey, accounts)
}

func ContextAccounts(ctx context.Context) []string {
	v, ok := ctx.Value(accountCtxKey).([]string)
	if !ok {
		return []string{}
	}
	return v
}
