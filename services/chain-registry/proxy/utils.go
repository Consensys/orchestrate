package proxy

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
)

type proxyURLCtxKey string

const ProxyURLKey proxyURLCtxKey = "proxy_url"

func With(ctx context.Context, proxyURL string) context.Context {
	return context.WithValue(ctx, ProxyURLKey, proxyURL)
}

func FromContext(ctx context.Context) string {
	proxyURL, _ := ctx.Value(ProxyURLKey).(string)
	return proxyURL
}

// GetURL is a util to retrieve the chain proxy url from the txctx context
func GetURL(txctx *engine.TxContext) (string, error) {
	chainProxyURL := FromContext(txctx.Context())
	if chainProxyURL == "" {
		err := txctx.AbortWithError(errors.InternalError("chain proxy url not found"))
		txctx.Logger.Error(err)
		return "", err
	}

	return chainProxyURL, nil
}
