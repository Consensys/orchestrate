package infra

import (
	"fmt"

	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/ethereum.git/tx-listener"
)

// ReceiptUnmarshaller assumes that input message is a go-ethereum receipt
type ReceiptUnmarshaller struct{}

// Unmarshal message expected to be a trace protobuffer
func (u *ReceiptUnmarshaller) Unmarshal(msg interface{}, t *types.Trace) error {
	// Cast message into receipt
	receipt, ok := msg.(*listener.TxListenerReceipt)
	if !ok {
		return fmt.Errorf("Message does not match expected format")
	}

	t.Chain().ID.Set(receipt.ChainID)

	// Load trace receipt from protobuffer
	t.Receipt().PostState = receipt.PostState
	t.Receipt().Status = receipt.Status
	t.Receipt().CumulativeGasUsed = receipt.CumulativeGasUsed
	t.Receipt().Bloom.SetBytes(receipt.Bloom.Bytes())
	for _, log := range receipt.Logs {
		t.Receipt().Logs = append(t.Receipt().Logs, &types.Log{Log: *log, DecodedData: map[string]string{}})
	}
	t.Receipt().TxHash.SetBytes(receipt.TxHash.Bytes())
	t.Receipt().ContractAddress.SetBytes(receipt.ContractAddress.Bytes())
	t.Receipt().GasUsed = receipt.GasUsed
	t.Receipt().BlockHash.SetBytes(receipt.BlockHash.Bytes())
	t.Receipt().BlockNumber = uint64(receipt.BlockNumber)
	t.Receipt().TxIndex = receipt.TxIndex

	return nil
}
