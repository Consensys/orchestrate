package crafter

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/ethclient"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/ethereum/abi"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/chain-registry/proxy"
	svc "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/contract-registry/proto"
)

// Crafter creates a crafter handler
func Crafter(r svc.ContractRegistryClient, crafter abi.Crafter, ec ethclient.EEAChainStateReader) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		txctx.Logger.
			WithField("envelope_id", txctx.Envelope.GetID()).
			WithField("job_uuid", txctx.Envelope.GetJobUUID()).
			Debugf("crafter handler starts")

		if txctx.Envelope.IsEeaSendMarkingTransaction() {
			url, err := proxy.GetURL(txctx)
			if err != nil {
				e := txctx.AbortWithError(err).ExtendComponent(component)
				txctx.Logger.WithError(e).Error("crafter: could not fetch envelope proxy url")
				return
			}
			privPContractAddr, err := ec.EEAPrivPrecompiledContractAddr(txctx.Context(), url)
			if err != nil {
				e := txctx.AbortWithError(err).ExtendComponent(component)
				txctx.Logger.WithError(e).Error("crafter: could not fetch eea precompiled contract address")
				return
			}

			_ = txctx.Envelope.SetTo(privPContractAddr)
		}

		if txctx.Envelope.GetData() != "" || txctx.Envelope.GetMethodSignature() == "" {
			// If transaction has already been crafted there is nothing to do
			return
		}

		txctx.Logger = txctx.Logger.WithFields(log.Fields{
			"method": txctx.Envelope.GetMethodSignature(),
			"args":   txctx.Envelope.GetArgs(),
		})

		var data []byte
		if txctx.Envelope.IsContractCreation() {
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
		txctx.Logger.Tracef("crafter: tx data payload set with %s", txctx.Envelope.GetData())
	}
}
