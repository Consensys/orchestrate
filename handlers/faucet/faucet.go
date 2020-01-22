package faucet

import (
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	authutils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/authentication/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/proxy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/faucet/faucet"
	faucettypes "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/faucet/types"
)

// Faucet creates a Faucet handler
func Faucet(multitenancy bool, fct faucet.Faucet) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		// Beneficiary
		beneficiary := txctx.Envelope.Sender()

		url, err := proxy.GetURL(txctx)
		if err != nil {
			return
		}

		// Create Faucet request
		req := &faucettypes.Request{
			ChainID:     txctx.Envelope.GetChain().ID(),
			NodeURL:     url,
			NodeID:      txctx.Envelope.GetChain().GetNodeId(),
			NodeName:    txctx.Envelope.GetChain().GetNodeName(),
			Beneficiary: beneficiary,
			Amount:      txctx.Envelope.GetTx().GetTxData().GetValueBig(),
		}

		if multitenancy {
			auth := authutils.AuthorizationFromContext(txctx.Context())
			if auth == "" {
				er := txctx.AbortWithError(errors.UnauthorizedError("missing Access Token")).ExtendComponent(component)
				txctx.Logger.WithError(er).Errorf("Token Not Found: could extract the Access Token from the envelop")
				return
			}
			req.AuthToken = auth
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
