package crafter

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/ethereum/abi"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
	svc "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/proto"
)

// Crafter creates a crafter handler
func Crafter(r svc.ContractRegistryClient, crafter abi.Crafter) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		if txctx.Envelope.GetData() != "" || txctx.Envelope.GetMethodSignature() == "" {
			// If transaction has already been crafted there is nothing to do
			return
		}

		txctx.Logger = txctx.Logger.WithFields(log.Fields{
			"method": txctx.Envelope.GetMethodSignature(),
			"args":   txctx.Envelope.GetArgs(),
		})

		var data []byte
		if txctx.Envelope.IsConstructor() {
			// Load contract from contract registry
			resp, err := r.GetContractBytecode(
				txctx.Context(),
				&svc.GetContractRequest{
					ContractId: txctx.Envelope.GetContractID(),
				},
			)
			if err != nil {
				e := txctx.AbortWithError(err).ExtendComponent(component)
				txctx.Logger.WithError(e).Errorf("crafter: could not retrieve contract bytecode")
				return
			}

			// Craft transaction payload
			bytecode, err := hexutil.Decode(resp.GetBytecode())
			if err != nil {
				e := txctx.AbortWithError(err).ExtendComponent(component)
				txctx.Logger.WithError(e).Errorf("crafter: invalid contract bytecode")
				return
			}

			data, err = crafter.CraftConstructor(
				bytecode,
				txctx.Envelope.GetMethodSignature(),
				txctx.Envelope.GetArgs()...,
			)
			if err != nil {
				e := txctx.AbortWithError(err).ExtendComponent(component)
				txctx.Logger.WithError(e).Errorf("crafter: could not craft contract deployment payload")
				return
			}
		} else {
			// Craft transaction payload
			var err error
			data, err = crafter.CraftCall(
				txctx.Envelope.GetMethodSignature(),
				txctx.Envelope.GetArgs()...,
			)
			if err != nil {
				e := txctx.AbortWithError(err).ExtendComponent(component)
				txctx.Logger.WithError(e).Errorf("crafter: could not craft transaction payload")
				return
			}
		}

		txctx.Logger = txctx.Logger.WithFields(log.Fields{
			"data": utils.ShortString(hexutil.Encode(data), 10),
		})

		_ = txctx.Envelope.SetData(data)

		txctx.Logger.Tracef("crafter: tx data payload set")
	}
}
