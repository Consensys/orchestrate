package api

import (
	"github.com/consensys/orchestrate/pkg/utils"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type TransferRequest struct {
	ChainName string            `json:"chain" validate:"required" example:"myChain"`
	Labels    map[string]string `json:"labels,omitempty"`
	Params    TransferParams    `json:"params" validate:"required"`
}

type TransferParams struct {
	Value           string            `json:"value" validate:"required,isBig" example:"71500000 (wei)"`
	Gas             string            `json:"gas,omitempty" example:"21000"`
	GasPrice        string            `json:"gasPrice,omitempty" validate:"omitempty,isBig" example:"71500000 (wei)"`
	GasFeeCap       string            `json:"maxFeePerGas,omitempty" example:"71500000 (wei)"`
	GasTipCap       string            `json:"maxPriorityFeePerGas,omitempty" example:"71500000 (wei)"`
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
