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
		methodAbi, err := getMethodAbi(txctx)
		if err != nil || methodAbi == nil {
			return
		}

		log.WithFields(log.Fields{
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
	if ABI := txctx.Envelope.GetCall().GetMethod().GetAbi(); len(ABI) > 0 {
		// ABI was provided in the Envelope
		err := json.Unmarshal(ABI, method)
		if err != nil {
			_ = txctx.AbortWithError(err)
			return nil, err
		}
	} else if methodSig := txctx.Envelope.GetCall().GetMethod().GetSignature(); methodSig != "" {
		// Generate method ABI from signature
		m, err := crafter.SignatureToMethod(methodSig)

		if err != nil {
			txctx.Logger.WithError(err).Errorf("crafter: could not generate method ABI from signature")
			_ = txctx.AbortWithError(err)
			return nil, err
		}
		method = m
		txctx.Logger = txctx.Logger.WithFields(log.Fields{
			"crafter.method": methodSig,
		})
	}
	// Nothing to craft

	return method, nil
}

func createTxPayload(txctx *engine.TxContext, methodAbi *abi.Method, r registry.Registry, c crafter.Crafter) ([]byte, error) {
	if txctx.Envelope.GetCall().GetMethod().IsConstructor() {
		return createContractDeploymentPayload(txctx, methodAbi, r, c)
	}

	return createTxCallPayload(txctx, methodAbi, c)
}

func createContractDeploymentPayload(txctx *engine.TxContext, methodAbi *abi.Method, r registry.Registry, c crafter.Crafter) ([]byte, error) {
	var (
		bytecode []byte
		payload  []byte
		err      error
	)
	// Transaction to be crafted is a Contract deployment
	bytecode, err = r.GetContractBytecode(txctx.Envelope.GetCall().GetContract())
	if err != nil {
		txctx.Logger.WithError(err).Errorf("crafter: could not retrieve contract bytecode")
		_ = txctx.AbortWithError(err)
		return nil, err
	}

	payload, err = c.CraftConstructor(bytecode, *methodAbi, getTxArgs(txctx)...)
	if err != nil {
		txctx.Logger.WithError(err).Errorf("crafter: could not craft tx payload")
		_ = txctx.AbortWithError(err)
		return nil, err
	}

	return payload, nil
}

func createTxCallPayload(txctx *engine.TxContext, methodAbi *abi.Method, c crafter.Crafter) ([]byte, error) {
	var payload, err = c.CraftCall(*methodAbi, getTxArgs(txctx)...)
	if err != nil {
		txctx.Logger.WithError(err).Errorf("crafter: could not craft tx payload")
		_ = txctx.AbortWithError(err)
		return nil, err
	}

	return payload, nil
}

func getTxArgs(txctx *engine.TxContext) []string {
	args := txctx.Envelope.GetCall().GetArgs()
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
