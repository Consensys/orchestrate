package api

import (
	"github.com/consensys/orchestrate/pkg/utils"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
)

type TransferRequest struct {
	ChainName string            `json:"chain" validate:"required" example:"myChain"`
	Labels    map[string]string `json:"labels,omitempty"`
	Params    TransferParams    `json:"params" validate:"required"`
}

type TransferParams struct {
	Value           *hexutil.Big      `json:"value" validate:"required" example:"0x59682f00" swaggertype:"string"`
	Gas             *uint64           `json:"gas,omitempty" example:"21000"`
	GasPrice        *hexutil.Big      `json:"gasPrice,omitempty" example:"0x5208" swaggertype:"string"`
	GasFeeCap       *hexutil.Big      `json:"maxFeePerGas,omitempty" example:"0x4c4b40" swaggertype:"string"`
	GasTipCap       *hexutil.Big      `json:"maxPriorityFeePerGas,omitempty" example:"0x59682f00" swaggertype:"string"`
	AccessList      types.AccessList  `json:"accessList,omitempty" swaggertype:"array,object"`
	TransactionType string            `json:"transactionType,omitempty" validate:"omitempty,isTransactionType" example:"dynamic_fee" enums:"legacy,dynamic_fee"`
	From            ethcommon.Address `json:"from" validate:"required" example:"0x1abae27a0cbfb02945720425d3b80c7e09728534" swaggertype:"string"`
	To              ethcommon.Address `json:"to" validate:"required" example:"0x1abae27a0cbfb02945720425d3b80c7e09728534" swaggertype:"string"`
	GasPricePolicy  GasPriceParams    `json:"gasPricePolicy,omitempty"`
}

func (params *TransferParams) Validate() error {
	if err := utils.GetValidator().Struct(params); err != nil {
		return err
	}
	return params.GasPricePolicy.RetryPolicy.Validate()
}
