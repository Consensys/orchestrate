package protobuf

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethpb "gitlab.com/ConsenSys/client/fr/core-stack/core/protobuf/ethereum"
	"gitlab.com/ConsenSys/client/fr/core-stack/core/types"
)

// LoadAddress load an hex string with 0x prefix to a go-ethereum Address object
func LoadAddress(hex string, a *common.Address) error {
	// Ensure address is a valid ethereum address
	if !common.IsHexAddress(hex) {
		return fmt.Errorf("Validation Error: %q is not a valid Ethereum address", hex)
	}

	// Set address value
	*a = common.HexToAddress(hex)

	return nil
}

// DumpAddress dump go-ethereum Address object to an hex string with 0x prefix
func DumpAddress(a *common.Address, hex *string) {
	if a != nil {
		// Set string hex Value
		*hex = a.Hex()
	}
}

// LoadQuantity load an hex string with 0x prefix to a big.Int
func LoadQuantity(hex string, q *big.Int) error {
	i, err := hexutil.DecodeBig(hex)
	if err != nil {
		return err
	}

	// Set quantity value
	*q = *i

	return nil
}

// DumpQuantity dump a big.Int to an hex string with 0x prefix
func DumpQuantity(q *big.Int, hex *string) {
	if q != nil {
		// Set string hex Value
		*hex = hexutil.EncodeBig(q)
	}
}

// LoadData load an hex string with 0x prefix to a []byte
func LoadData(hex string, data *[]byte) error {
	b, err := hexutil.Decode(hex)
	if err != nil {
		return err
	}

	// Set data value
	*data = b

	return nil
}

// DumpData dump a []byte to an hex string with 0x prefix
func DumpData(data []byte, hex *string) {
	*hex = hexutil.Encode(data)
}

// LoadHash load an hex string with 0x prefix to a go-ethereum Hash object
func LoadHash(hex string, h *common.Hash) error {
	*h = common.HexToHash(hex)
	return nil
}

// DumpHash dump a go-ethereum Hash object to an hex string with 0x prefix
func DumpHash(hash common.Hash, hex *string) {
	*hex = hash.Hex()
}

// LoadTxData load a TxData protobuffer to a TxData object
func LoadTxData(pb *ethpb.TxData, txData *types.TxData) {
	if pb == nil {
		pb = &ethpb.TxData{}
	}

	txData.Nonce = pb.GetNonce()

	if txData.To == nil {
		var a common.Address
		txData.To = &a
	}
	LoadAddress(pb.To, txData.To)

	if txData.Value == nil {
		var v big.Int
		txData.Value = &v
	}
	LoadQuantity(pb.GetValue(), txData.Value)

	txData.GasLimit = pb.GetGas()

	if txData.GasPrice == nil {
		var l big.Int
		txData.GasPrice = &l
	}
	LoadQuantity(pb.GetGasPrice(), txData.GasPrice)

	LoadData(pb.GetData(), &txData.Data)
}

// DumpTxData load a TxData protobuffer to a TxData object
func DumpTxData(txData *types.TxData, pb *ethpb.TxData) {
	pb.Nonce = txData.GetNonce()
	DumpAddress(txData.GetTo(), &pb.To)
	DumpQuantity(txData.GetValue(), &pb.Value)
	pb.Gas = txData.GetGasLimit()
	DumpQuantity(txData.GetGasPrice(), &pb.GasPrice)
	DumpData(txData.GetData(), &pb.Data)
}

// LoadTransaction load a Transaction protobuffer to a Transaction object
func LoadTransaction(pb *ethpb.Transaction, tx *types.Transaction) {
	if pb == nil {
		pb = &ethpb.Transaction{}
	}

	if tx.TxData == nil {
		var data types.TxData
		tx.TxData = &data
	}
	LoadTxData(pb.TxData, tx.TxData)

	LoadData(pb.GetRaw(), &tx.Raw)

	if tx.Hash == nil {
		var h common.Hash
		tx.Hash = &h
	}
	LoadHash(pb.GetHash(), tx.Hash)

	if tx.From == nil {
		var a common.Address
		tx.From = &a
	}
	LoadAddress(pb.From, tx.From)
}

// DumpTransaction dump Transaction object to a transaction protobuffer
func DumpTransaction(tx *types.Transaction, pb *ethpb.Transaction) {
	if pb.TxData == nil {
		var data ethpb.TxData
		pb.TxData = &data
	}
	DumpTxData(tx.GetTxData(), pb.TxData)
	DumpData(tx.GetRaw(), &pb.Raw)
	DumpHash(*tx.GetHash(), &pb.Hash)
	DumpAddress(tx.GetFrom(), &pb.From)
}
