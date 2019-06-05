package crafter

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/ethereum"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/ethereum.git/abi/crafter"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/ethereum.git/abi/registry"
)

// Crafter creates a crafter handler
func Crafter(r registry.Registry, c crafter.Crafter) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		if txctx.Envelope.GetTx().GetTxData().GetData() != "" {
			// If transaction has already been crafted there is nothing to do
			return
		}

		// Try to read ABI in Envelope Call
		var method *abi.Method
		if ABI := txctx.Envelope.GetCall().GetMethod().GetAbi(); len(ABI) > 0 {
			// ABI was provided in the Envelope
			err := json.Unmarshal(ABI, method)
			if err != nil {
				_ = txctx.AbortWithError(err)
				return
			}
		} else if methodSig := txctx.Envelope.GetCall().GetMethod().GetSignature(); methodSig != "" {
			// Generate method ABI from signature
			m, err := crafter.SignatureToMethod(methodSig)

			if err != nil {
				txctx.Logger.WithError(err).Errorf("crafter: could not generate method ABI from signature")
				_ = txctx.AbortWithError(err)
				return
			}
			method = m
			txctx.Logger = txctx.Logger.WithFields(log.Fields{
				"crafter.method": methodSig,
			})
		} else {
			// Nothing to craft
			return
		}

		// Retrieve  Args from Envelope
		args := txctx.Envelope.GetCall().GetArgs()
		txctx.Logger = txctx.Logger.WithFields(log.Fields{
			"crafter.args": args,
		})

		log.WithFields(log.Fields{
			"method.sig":    method.Sig(),
			"method.string": method.String(),
			"method.id":     hexutil.Encode(method.Id()),
		}).Debugf("crafter: method call")

		var (
			bytecode []byte
			payload  []byte
			err      error
		)
		if txctx.Envelope.GetCall().GetMethod().IsConstructor() {
			// Transaction to be crafted is a Contract deployment
			// Retrieve Bytecode from registry
			bytecode, err = r.GetContractBytecode(txctx.Envelope.GetCall().GetContract())
			if err != nil {
				txctx.Logger.WithError(err).Errorf("crafter: could not retrieve contract bytecode")
				_ = txctx.AbortWithError(err)
				return
			}

			// Craft Deployment
			payload, err = c.CraftConstructor(bytecode, *method, args...)
			if err != nil {
				txctx.Logger.WithError(err).Errorf("crafter: could not craft tx payload")
				_ = txctx.AbortWithError(err)
				return
			}
		} else {
			// Transaction to be crafted is a contract call
			payload, err = c.CraftCall(*method, args...)
			if err != nil {
				txctx.Logger.WithError(err).Errorf("crafter: could not craft tx payload")
				_ = txctx.AbortWithError(err)
				return
			}
		}

		txctx.Logger = txctx.Logger.WithFields(log.Fields{
			"tx.data": utils.ShortString(hexutil.Encode(payload), 10),
		})

		// Attach transaction payload to Envelope
		if txctx.Envelope.GetTx() == nil {
			txctx.Envelope.Tx = &ethereum.Transaction{TxData: &ethereum.TxData{}}
		} else if txctx.Envelope.GetTx().GetTxData() == nil {
			txctx.Envelope.Tx.TxData = &ethereum.TxData{}
		}
		txctx.Envelope.GetTx().GetTxData().SetData(payload)

		txctx.Logger.Tracef("crafter: tx data payload set")
	}
}
