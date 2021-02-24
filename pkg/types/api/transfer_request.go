package api

import "github.com/ConsenSys/orchestrate/pkg/utils"

type TransferRequest struct {
	ChainName string            `json:"chain" validate:"required" example:"myChain"`
	Labels    map[string]string `json:"labels,omitempty"`
	Params    TransferParams    `json:"params" validate:"required"`
}

type TransferParams struct {
	Value          string         `json:"value" validate:"required,isBig" example:"71500000 (wei)"`
	Gas            string         `json:"gas,omitempty" example:"21000"`
	GasPrice       string         `json:"gasPrice,omitempty" validate:"omitempty,isBig" example:"71500000 (wei)"`
	From           string         `json:"from" validate:"required,eth_addr" example:"0x1abae27a0cbfb02945720425d3b80c7e09728534"`
	To             string         `json:"to" validate:"required,eth_addr" example:"0x1abae27a0cbfb02945720425d3b80c7e09728534"`
	GasPricePolicy GasPriceParams `json:"gasPricePolicy,omitempty"`
}

func (params *TransferParams) Validate() error {
	if err := utils.GetValidator().Struct(params); err != nil {
		return err
	}
	return params.GasPricePolicy.RetryPolicy.Validate()
}
