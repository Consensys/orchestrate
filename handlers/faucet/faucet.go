package faucet

import (
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/proxy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/faucet/faucet"
	faucettypes "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/faucet/types"
)

// Faucet creates a Faucet handler
func Faucet(fct faucet.Faucet) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		url, err := proxy.GetURL(txctx)
		if err != nil {
			return
		}

		if txctx.Builder.GetChainID() == nil || txctx.Builder.GetValue() == nil {
			return
		}

		// Create Faucet request
		req := &faucettypes.Request{
			ChainID:     txctx.Builder.GetChainID(),
			ChainURL:    url,
			ChainUUID:   txctx.Builder.GetChainUUID(),
			ChainName:   txctx.Builder.GetChainName(),
			Beneficiary: txctx.Builder.MustGetFromAddress(),
			Amount:      txctx.Builder.GetValue(),
		}

		// Credit
		amount, err := fct.Credit(txctx.Context(), req)
		if err != nil {
			switch {
			case errors.IsFaucetSelfCreditWarning(err):
				return
			case errors.IsFaucetNotConfiguredWarning(err):
				e := errors.FromError(err).ExtendComponent(component)
				txctx.Logger.WithError(e).Debug("faucet: not configured")
				return
			case errors.IsWarning(err):
				e := errors.FromError(err).ExtendComponent(component)
				txctx.Logger.WithError(e).Debugf("faucet: credit refused")
				return
			default:
				e := errors.FromError(err).ExtendComponent(component)
				txctx.Logger.WithError(e).Error("faucet: credit error")
				return
			}
		}

		txctx.Logger.WithFields(log.Fields{
			"faucet.amount": amount.Text(10),
		}).Infof("faucet: credit approved")
	}
}
