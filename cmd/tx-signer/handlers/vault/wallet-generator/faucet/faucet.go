package faucet

import (
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/services/faucet/faucet"
	faucettypes "gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/services/faucet/types"
)

// Faucet creates a Faucet handler
func Faucet(fct faucet.Faucet) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {

		// Skip if the chainId is unset
		if txctx.Envelope.GetChain() == nil || txctx.Envelope.GetChain().ID() == nil {
			txctx.Logger.Debugf("faucet: skipping faucet request because no chainI")
			return
		}

		// Skip if no sender has been set
		if txctx.Envelope.Sender().Hex() == "0x0000000000000000000000000000000000000000" {
			txctx.Logger.Debugf("faucet: skipping faucet request because no sender address has been set")
			return
		}

		// Create Faucet request
		req := &faucettypes.Request{
			ChainID:     txctx.Envelope.GetChain().ID(),
			Beneficiary: txctx.Envelope.Sender(),
			Amount:      txctx.Envelope.GetTx().GetTxData().GetValueBig(),
		}

		// Credit
		amount, approved, err := fct.Credit(txctx.Context(), req)
		if err != nil {
			e := txctx.Error(err).ExtendComponent(component)
			txctx.Logger.WithError(e).Errorf("faucet: credit error")
			return
		}

		if approved {
			txctx.Logger.WithFields(log.Fields{
				"faucet.amount": amount.Text(10),
			}).Debugf("faucet: credit approved")
		} else {
			txctx.Logger.WithFields(log.Fields{
				"faucet.amount": amount.Text(10),
			}).Debugf("faucet: credit not approved")
		}
	}
}
