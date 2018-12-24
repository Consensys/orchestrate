package protobuf

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethpb "gitlab.com/ConsenSys/client/fr/core-stack/core/protobuf/ethereum"
	"gitlab.com/ConsenSys/client/fr/core-stack/core/types"
)

// LoadAddress load an hex string with 0x prefix to a go-ethereum Address object
func LoadAddress(hex string) (common.Address, error) {
	// Ensure address is a valid ethereum address
	var a common.Address
	if !common.IsHexAddress(hex) {
		return a, fmt.Errorf("Validation Error: %q is not a valid Ethereum address", hex)
	}

	// Set address value
	return common.HexToAddress(hex), nil
}

// LoadTx load a Transaction protobuffer to a Tx object
func LoadTx(pb *ethpb.Transaction, tx *types.Tx) error {
	txData := pb.GetTxData()

	tx.SetNonce(txData.GetNonce())

	a := common.HexToAddress(txData.GetTo())
	tx.SetTo(&a)

	v, err := hexutil.DecodeBig(txData.GetValue())
	if err != nil {
		return err
	}
	tx.SetValue(v)

	tx.SetGasLimit(txData.GetGas())

	p, err := hexutil.DecodeBig(txData.GetGasPrice())
	if err != nil {
		return err
	}
	tx.SetGasPrice(p)

	data, err := hexutil.Decode(txData.GetData())
	if err != nil {
		return err
	}
	tx.SetData(data)

	raw, err := hexutil.Decode(pb.GetRaw())
	if err != nil {
		return err
	}
	tx.SetRaw(raw)

	h := common.HexToHash(pb.GetHash())
	tx.SetHash(&h)

	return nil
}

// DumpTx dump Tx object to a transaction protobuffer
func DumpTx(tx *types.Tx, pb *ethpb.Transaction) {
	if pb.TxData == nil {
		pb.TxData = &ethpb.TxData{}
	}
	pb.TxData.Nonce = tx.Nonce()
	pb.TxData.To = tx.To().Hex()
	pb.TxData.Value = hexutil.EncodeBig(tx.Value())
	pb.TxData.Gas = tx.GasLimit()
	pb.TxData.GasPrice = hexutil.EncodeBig(tx.GasPrice())
	pb.TxData.Data = hexutil.Encode(tx.Data())

	pb.Raw = hexutil.Encode(tx.Raw())
	pb.Hash = tx.Hash().Hex()
}
