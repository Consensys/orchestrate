package crafter

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/ethereum/abi/crafter"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
	svc "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/contract-registry"
)

// Crafter creates a crafter handler
func Crafter(r svc.ContractRegistryClient, c crafter.Crafter) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		txctx.Logger = txctx.Logger.WithFields(log.Fields{
			"id": txctx.Builder.ID,
		})

		if txctx.Builder.GetData() != "" {
			// If transaction has already been crafted there is nothing to do
			return
		}

		// Try to read ABI in Builder Call
		methodAbi, err := getMethodAbi(txctx)
		if err != nil || methodAbi == nil {
			return
		}

		txctx.Logger.WithFields(log.Fields{
			"method.sig":    methodAbi.Sig(),
			"method.string": methodAbi.String(),
			"method.id":     hexutil.Encode(methodAbi.ID()),
		}).Debugf("crafter: method call")

		payload, err := createTxPayload(txctx, methodAbi, r, c)
		if err != nil {
			return
		}

		txctx.Logger = txctx.Logger.WithFields(log.Fields{
			"data": utils.ShortString(hexutil.Encode(payload), 10),
		})

		_ = txctx.Builder.SetData(payload)

		txctx.Logger.Tracef("crafter: tx data payload set")
	}
}

func getMethodAbi(txctx *engine.TxContext) (*abi.Method, error) {
	if txctx.Builder.MethodSignature == "" {
		return nil, errors.DataError("No method signature provided")
	}

	// Generate method ABI from signature
	method, err := crafter.SignatureToMethod(txctx.Builder.MethodSignature)
	if err != nil {
		e := txctx.AbortWithError(err).ExtendComponent(component)
		txctx.Logger.WithError(e).Errorf("crafter: could not generate method ABI from signature")
		return nil, e
	}
	txctx.Logger = txctx.Logger.WithFields(log.Fields{
		"crafter.method": txctx.Builder.GetMethodSignature(),
	})

	return method, nil
}

func createTxPayload(txctx *engine.TxContext, methodAbi *abi.Method, r svc.ContractRegistryClient, c crafter.Crafter) ([]byte, error) {
	if txctx.Builder.IsConstructor() {
		return createContractDeploymentPayload(txctx, methodAbi, r, c)
	}

	return createTxCallPayload(txctx, methodAbi, c)
}

func createContractDeploymentPayload(txctx *engine.TxContext, methodAbi *abi.Method, r svc.ContractRegistryClient, c crafter.Crafter) ([]byte, error) {
	// Transaction to be crafted is a Contract deployment
	bytecodeResp, err := r.GetContractBytecode(
		txctx.Context(),
		&svc.GetContractRequest{
			ContractId: txctx.Builder.GetContractID(),
		},
	)
	if err != nil {
		e := txctx.AbortWithError(err).ExtendComponent(component)
		txctx.Logger.WithError(e).Errorf("crafter: could not retrieve contract bytecode")
		return nil, e
	}

	payload, err := c.CraftConstructor([]byte(bytecodeResp.GetBytecode()), methodAbi, getTxArgs(txctx)...)
	if err != nil {
		e := txctx.AbortWithError(err).ExtendComponent(component)
		txctx.Logger.WithError(e).Errorf("crafter: could not craft tx payload")
		return nil, e
	}

	return payload, nil
}

func createTxCallPayload(txctx *engine.TxContext, methodAbi *abi.Method, c crafter.Crafter) ([]byte, error) {
	var payload, err = c.CraftCall(methodAbi, getTxArgs(txctx)...)
	if err != nil {
		e := txctx.AbortWithError(err).ExtendComponent(component)
		txctx.Logger.WithError(e).Errorf("crafter: could not craft tx payload")
		return nil, e
	}

	return payload, nil
}

func getTxArgs(txctx *engine.TxContext) []string {
	args := txctx.Builder.GetArgs()
	txctx.Logger = txctx.Logger.WithFields(log.Fields{
		"crafter.args": args,
	})

	return args
}
