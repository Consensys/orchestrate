package infra

import (
	"fmt"

	listener "gitlab.com/ConsenSys/client/fr/core-stack/infra/ethereum.git/tx-listener"
	common "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/protos/common"
	ethereum "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/protos/ethereum"
	trace "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/protos/trace"
)

// ReceiptUnmarshaller assumes that input message is a go-ethereum receipt
type ReceiptUnmarshaller struct{}

// Unmarshal message expected to be a trace protobuffer
func (u *ReceiptUnmarshaller) Unmarshal(msg interface{}, t *trace.Trace) error {
	// Cast message into receipt
	receipt, ok := msg.(*listener.TxListenerReceipt)
	if !ok {
		return fmt.Errorf("Message does not match expected format")
	}

	// Set receipt
	t.Receipt = ethereum.FromGethReceipt(&receipt.Receipt).
		SetBlockHash(receipt.BlockHash).
		SetBlockNumber(uint64(receipt.BlockNumber)).
		SetTxIndex(receipt.TxIndex)
	t.Chain = (&common.Chain{}).SetID(receipt.ChainID)

	return nil
}
