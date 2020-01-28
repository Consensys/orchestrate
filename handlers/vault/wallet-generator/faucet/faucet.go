package faucet

import (
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/faucet/faucet"
	faucettypes "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/faucet/types"
)

// Faucet creates a Faucet handler
func Faucet(fct faucet.Faucet) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {

		// Skip if the chainId is unset
		if txctx.Envelope.GetChain() == nil || txctx.Envelope.GetChain().GetBigChainID() == nil {
			txctx.Logger.Debugf("faucet: skipping faucet request because no chainID")
			return
		}

		// Skip if no sender has been set
		if txctx.Envelope.Sender().Hex() == "0x0000000000000000000000000000000000000000" {
			txctx.Logger.Debugf("faucet: skipping faucet request because no sender address has been set")
			return
		}

		// Create Faucet request
		req := &faucettypes.Request{
			ChainID:     txctx.Envelope.GetChain().GetBigChainID(),
			Beneficiary: txctx.Envelope.Sender(),
			Amount:      txctx.Envelope.GetTx().GetTxData().GetValueBig(),
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
