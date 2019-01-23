package protobuf

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethpb "gitlab.com/ConsenSys/client/fr/core-stack/core.git/protobuf/ethereum"
	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/types"
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

// LoadLog load a Log protobuffer to a Log object
func LoadLog(pb *ethpb.Log, l *types.Log) error {
	if ok := common.IsHexAddress(pb.GetAddress()); !ok {
		return fmt.Errorf("Invalid Address")
	}

	data, err := hexutil.Decode(pb.GetData())
	if err != nil {
		return err
	}

	l.Address.SetBytes(common.FromHex(pb.GetAddress()))

	l.Topics = l.Topics[0:0]
	for _, topic := range pb.GetTopics() {
		l.Topics = append(l.Topics, common.HexToHash(topic))
	}

	l.Data = data
	l.DecodedData = pb.GetDecodedData()
	l.BlockNumber = pb.GetBlockNumber()
	l.TxHash.SetBytes(common.FromHex(pb.GetTxHash()))
	l.TxIndex = uint(pb.GetTxIndex())
	l.BlockHash.SetBytes(common.FromHex(pb.GetBlockHash()))

	l.Index = uint(pb.GetIndex())
	l.Removed = pb.GetRemoved()

	return nil
}

// DumpLog dump a Log object to protobuffer
func DumpLog(l *types.Log, pb *ethpb.Log) {
	pb.Address = l.Address.Hex()
	pb.Topics = pb.Topics[0:0]
	for _, topic := range l.Topics {
		pb.Topics = append(pb.Topics, topic.Hex())
	}
	pb.Data = hexutil.Encode(l.Data)
	pb.DecodedData = l.DecodedData
	pb.BlockNumber = l.BlockNumber
	pb.TxHash = l.TxHash.Hex()
	pb.TxIndex = uint64(l.TxIndex)
	pb.BlockHash = l.BlockHash.Hex()
	pb.Index = uint64(l.Index)
	pb.Removed = l.Removed
}

// LoadReceipt load a Receipt protobuffer to a Receipt object
func LoadReceipt(pb *ethpb.Receipt, r *types.Receipt) error {
	s, err := hexutil.Decode(pb.GetPostState())
	if err != nil {
		return err
	}

	h, err := hexutil.Decode(pb.GetTxHash())
	if err != nil {
		return err
	}

	b, err := hexutil.Decode(pb.GetBloom())
	if err != nil {
		return err
	}

	if ok := common.IsHexAddress(pb.GetContractAddress()); !ok {
		return fmt.Errorf("Invalid Address")
	}

	logs := []*types.Log{}
	for _, log := range pb.GetLogs() {
		var l types.Log
		LoadLog(log, &l)
		logs = append(logs, &l)
	}

	r.Logs = logs
	r.ContractAddress.SetBytes(common.FromHex(pb.GetContractAddress()))
	r.PostState = s
	r.Status = pb.GetStatus()
	r.TxHash.SetBytes(h)
	r.Bloom.SetBytes(b)

	r.GasUsed = pb.GetGasUsed()
	r.CumulativeGasUsed = pb.GetCumulativeGasUsed()

	return nil
}

// DumpReceipt dump a Receipt object into a protobuffer
func DumpReceipt(r *types.Receipt, pb *ethpb.Receipt) error {
	pb.ContractAddress = r.ContractAddress.Hex()
	pb.Status = r.Status

	pb.GasUsed = r.GasUsed
	pb.CumulativeGasUsed = r.CumulativeGasUsed

	pb.TxHash = r.TxHash.Hex()
	pb.Bloom = common.ToHex(r.Bloom.Bytes())

	pb.PostState = hexutil.Encode(r.PostState)

	pb.Logs = pb.Logs[0:0]
	for _, log := range r.Logs {
		var pblog ethpb.Log
		DumpLog(log, &pblog)
		pb.Logs = append(pb.Logs, &pblog)
	}

	return nil
}
