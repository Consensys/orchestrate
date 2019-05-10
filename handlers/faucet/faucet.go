package faucet

import (
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/faucet.git/faucet"
	faucettypes "gitlab.com/ConsenSys/client/fr/core-stack/infra/faucet.git/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine"
)

// Faucet creates a Faucet handler
func Faucet(fct faucet.Faucet) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		// Beneficiary
		beneficiary, _ := txctx.Envelope.GetSender().Address()

		// Create Faucet request
		req := &faucettypes.Request{
			ChainID:     txctx.Envelope.GetChain().ID(),
			Beneficiary: beneficiary,
			Amount:      txctx.Envelope.GetTx().GetTxData().ValueBig(),
		}

		// Credit
		amount, approved, err := fct.Credit(txctx.Context(), req)
		if err != nil {
			// TODO: handle error
			txctx.Logger.WithError(err).Errorf("faucet: credit error")
			_ = txctx.Error(err)
		} else {
			if !approved {
				txctx.Logger.Debugf("faucet: credit not approved")
			} else {
				txctx.Logger = txctx.Logger.WithFields(log.Fields{
					"faucet.amount": amount.Text(10),
				})
				txctx.Logger.Debugf("faucet: credit approved")
			}
		}
	}
}
