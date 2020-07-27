package faucet

import (
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/client"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/proxy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/faucet/faucet"
	faucettypes "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/faucet/types"
)

// Faucet creates a Faucet handler
func Faucet(fct faucet.Faucet, faucetClient client.FaucetClient) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		txctx.Logger.WithField("envelope_id", txctx.Envelope.GetID()).Debugf("faucet handler starts")

		if txctx.Envelope.GetChainUUID() == "" || txctx.Envelope.GetContextLabelsValue("faucet.parentTxID") != "" ||
			txctx.Envelope.GetFromString() == "" {
			return
		}

		faucets, err := faucetClient.GetFaucetsByChainRule(txctx.Context(), txctx.Envelope.GetChainUUID())
		if err != nil {
			_ = txctx.Error(errors.FaucetWarning("could not get faucets for chain rule '%s' - got %v", txctx.Envelope.GetChainUUID(), err)).ExtendComponent(component)
			return
		}
		if len(faucets) == 0 {
			return
		}

		url, err := proxy.GetURL(txctx)
		if err != nil {
			_ = txctx.Error(errors.FaucetWarning("could not get chain url - got %v", err)).ExtendComponent(component)
			return
		}

		req := &faucettypes.Request{
			ScheduleUUID:      txctx.Envelope.GetContextLabelsValue("scheduleUUID"),
			ParentTxID:        txctx.Envelope.GetID(),
			ChildTxID:         txctx.Envelope.GetContextLabelsValue("faucetChildTxID"),
			ChainID:           txctx.Envelope.GetChainID(),
			ChainURL:          url,
			ChainName:         txctx.Envelope.GetChainName(),
			ChainUUID:         txctx.Envelope.GetChainUUID(),
			Beneficiary:       txctx.Envelope.MustGetFromAddress(),
			FaucetsCandidates: faucettypes.NewFaucetsCandidates(faucets),
		}

		// Credit
		amount, err := fct.Credit(txctx.Context(), req)
		if err != nil {
			switch {
			case errors.IsFaucetSelfCreditWarning(err):
				return
			case errors.IsFaucetWarning(err):
				e := errors.FromError(err).ExtendComponent(component)
				txctx.Logger.WithError(e).Debug("faucet: credit refused")
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
