package assets

import (
	"context"
	"fmt"

	"github.com/consensys/orchestrate/pkg/toolkit/app/log"
	"github.com/consensys/orchestrate/pkg/types/api"
	ethcommon "github.com/ethereum/go-ethereum/common"

	orchestrateclient "github.com/consensys/orchestrate/pkg/sdk/client"
)

var accountCtxKey ctxKey = "accounts"

func CreateNewAccount(ctx context.Context, client orchestrateclient.OrchestrateClient) (context.Context, error) {
	logger := log.FromContext(ctx)
	logger.Debug("registering new account...")
	resp, err := client.CreateAccount(ctx, &api.CreateAccountRequest{})
	if err != nil {
		errMsg := "failed to create account"
		logger.WithError(err).Error(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	logger.WithField("address", resp.Address).Info("account has been registered")
	return contextWithAccounts(ctx, append(ContextAccounts(ctx), resp.Address)), nil
}

func contextWithAccounts(ctx context.Context, accounts []ethcommon.Address) context.Context {
	return context.WithValue(ctx, accountCtxKey, accounts)
}

func ContextAccounts(ctx context.Context) []ethcommon.Address {
	if v, ok := ctx.Value(accountCtxKey).([]ethcommon.Address); ok {
		return v
	}
	return []ethcommon.Address{}
}
