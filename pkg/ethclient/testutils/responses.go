package testutils

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"

	"github.com/ConsenSys/orchestrate/pkg/ethclient/utils"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

func MakeRespBody(result interface{}, errMsg string) io.ReadCloser {
	respMsg := &utils.JSONRpcMessage{}
	if result != nil {
		if b, ok := result.([]byte); ok {
			respMsg.Result = json.RawMessage(b)
		} else {
			jsonResult, _ := json.Marshal(result)
			respMsg.Result = json.RawMessage(jsonResult)
		}
	}

	if errMsg != "" {
		respMsg.Error = &utils.JSONError{
			Message: errMsg,
		}
	}
	b, _ := json.Marshal(respMsg)
	return ioutil.NopCloser(bytes.NewReader(b))
}

func NewReceiptResp(r *ethtypes.Receipt) *ReceiptResp {
	return &ReceiptResp{
		CumulativeGasUsed: hexutil.Uint64(r.CumulativeGasUsed),
		Bloom:             r.Bloom,
		Logs:              r.Logs,
		TxHash:            r.TxHash,
		GasUsed:           hexutil.Uint64(r.GasUsed),
	}
}

type ReceiptResp struct {
	CumulativeGasUsed hexutil.Uint64  `json:"cumulativeGasUsed" gencodec:"required"`
	Bloom             ethtypes.Bloom  `json:"logsBloom"         gencodec:"required"`
	Logs              []*ethtypes.Log `json:"logs"              gencodec:"required"`
	TxHash            ethcommon.Hash  `json:"transactionHash" gencodec:"required"`
	GasUsed           hexutil.Uint64  `json:"gasUsed" gencodec:"required"`
}
