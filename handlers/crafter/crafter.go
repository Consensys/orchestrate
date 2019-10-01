package crafter

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/ethereum/abi/crafter"
	encoding "gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/encoding/json"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/engine"
	contractregistry "gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/services/contract-registry"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/types/ethereum"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/utils"
)

// Crafter creates a crafter handler
func Crafter(r contractregistry.RegistryClient, c crafter.Crafter) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		txctx.Logger = txctx.Logger.WithFields(log.Fields{
			"metadata.id": txctx.Envelope.GetMetadata().GetId(),
		})

		if txctx.Envelope.GetTx().GetTxData().GetData() != nil {
			// If transaction has already been crafted there is nothing to do
			return
		}

		// Try to read ABI in Envelope Call
		methodAbi, err := getMethodAbi(txctx)
		if err != nil || methodAbi == nil {
			return
		}

		txctx.Logger.WithFields(log.Fields{
			"method.sig":    methodAbi.Sig(),
			"method.string": methodAbi.String(),
			"method.id":     hexutil.Encode(methodAbi.Id()),
		}).Debugf("crafter: method call")

		payload, err := createTxPayload(txctx, methodAbi, r, c)
		if err != nil {
			return
		}

		txctx.Logger = txctx.Logger.WithFields(log.Fields{
			"tx.data": utils.ShortString(hexutil.Encode(payload), 10),
		})

		attachTxPayload(txctx, payload)

		txctx.Logger.Tracef("crafter: tx data payload set")
	}
}

func getMethodAbi(txctx *engine.TxContext) (*abi.Method, error) {
	var method *abi.Method
	if ABI := txctx.Envelope.GetArgs().GetCall().GetMethod().GetAbi(); len(ABI) > 0 {
		// ABI was provided in the Envelope
		err := encoding.Unmarshal(ABI, method)
		if err != nil {
			e := txctx.AbortWithError(err).ExtendComponent(component)
			txctx.Logger.WithError(e).Errorf("crafter: invalid ABI provided")
			return nil, err
		}
	} else if methodSig := txctx.Envelope.GetArgs().GetCall().GetMethod().GetSignature(); methodSig != "" {
		// Generate method ABI from signature
		m, err := crafter.SignatureToMethod(methodSig)
		if err != nil {
			e := txctx.AbortWithError(err).ExtendComponent(component)
			txctx.Logger.WithError(e).Errorf("crafter: could not generate method ABI from signature")
			return nil, e
		}
		method = m
		txctx.Logger = txctx.Logger.WithFields(log.Fields{
			"crafter.method": methodSig,
		})
	}
	// Nothing to craft

	return method, nil
}

func createTxPayload(txctx *engine.TxContext, methodAbi *abi.Method, r contractregistry.RegistryClient, c crafter.Crafter) ([]byte, error) {
	if txctx.Envelope.GetArgs().GetCall().GetMethod().IsConstructor() {
		return createContractDeploymentPayload(txctx, methodAbi, r, c)
	}

	return createTxCallPayload(txctx, methodAbi, c)
}

func createContractDeploymentPayload(txctx *engine.TxContext, methodAbi *abi.Method, r contractregistry.RegistryClient, c crafter.Crafter) ([]byte, error) {
	// Transaction to be crafted is a Contract deployment
	bytecodeResp, err := r.GetContractBytecode(
		txctx.Context(),
		&contractregistry.GetContractRequest{
			ContractId: txctx.Envelope.GetArgs().GetCall().GetContract().GetId(),
		},
	)
	if err != nil {
		e := txctx.AbortWithError(err).ExtendComponent(component)
		txctx.Logger.WithError(e).Errorf("crafter: could not retrieve contract bytecode")
		return nil, e
	}

	payload, err := c.CraftConstructor(bytecodeResp.GetBytecode(), *methodAbi, getTxArgs(txctx)...)
	if err != nil {
		e := txctx.AbortWithError(err).ExtendComponent(component)
		txctx.Logger.WithError(e).Errorf("crafter: could not craft tx payload")
		return nil, e
	}

	return payload, nil
}

func createTxCallPayload(txctx *engine.TxContext, methodAbi *abi.Method, c crafter.Crafter) ([]byte, error) {
	var payload, err = c.CraftCall(*methodAbi, getTxArgs(txctx)...)
	if err != nil {
		e := txctx.AbortWithError(err).ExtendComponent(component)
		txctx.Logger.WithError(e).Errorf("crafter: could not craft tx payload")
		return nil, e
	}

	return payload, nil
}

func getTxArgs(txctx *engine.TxContext) []string {
	args := txctx.Envelope.GetArgs().GetCall().GetArgs()
	txctx.Logger = txctx.Logger.WithFields(log.Fields{
		"crafter.args": args,
	})

	return args
}

func attachTxPayload(txctx *engine.TxContext, payload []byte) {
	attachTxDataIfMissing(txctx)
	setTxPayload(txctx, payload)
}

func attachTxDataIfMissing(txctx *engine.TxContext) {
	if txctx.Envelope.GetTx() == nil {
		txctx.Envelope.Tx = &ethereum.Transaction{TxData: &ethereum.TxData{}}
	} else if txctx.Envelope.GetTx().GetTxData() == nil {
		txctx.Envelope.Tx.TxData = &ethereum.TxData{}
	}
}

func setTxPayload(txctx *engine.TxContext, payload []byte) {
	txctx.Envelope.GetTx().GetTxData().SetData(payload)
}
