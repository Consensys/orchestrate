package faucet

import (
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/faucet.git/faucet"
	faucettypes "gitlab.com/ConsenSys/client/fr/core-stack/service/faucet.git/types"
)

// Faucet creates a Faucet handler
func Faucet(fct faucet.Faucet) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		// Beneficiary
		beneficiary := txctx.Envelope.Sender()

		// Create Faucet request
		req := &faucettypes.Request{
			ChainID:     txctx.Envelope.GetChain().ID(),
			Beneficiary: beneficiary,
			Amount:      txctx.Envelope.GetTx().GetTxData().GetValueBig(),
		}

		// Credit
		amount, approved, err := fct.Credit(txctx.Context(), req)
		if err != nil {
			e := txctx.Error(err).ExtendComponent(component)
			txctx.Logger.WithError(e).Errorf("faucet: credit error")
			return
		}

		if !approved {
			txctx.Logger.Debugf("faucet: credit not approved")
			return
		}

		txctx.Logger = txctx.Logger.WithFields(log.Fields{
			"faucet.amount": amount.Text(10),
		})
		txctx.Logger.Debugf("faucet: credit approved")
	}
}
