package api

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/utils"
)

type SendTransactionRequest struct {
	ChainName string            `json:"chain" validate:"required" example:"myChain"`
	Labels    map[string]string `json:"labels,omitempty"`
	Params    TransactionParams `json:"params" validate:"required"`
}

// go validator does not support mutually exclusive parameters for now
// See more https://github.com/go-playground/validator/issues/608
type TransactionParams struct {
	Value           string                        `json:"value,omitempty" validate:"omitempty,isBig" example:"71500000 (wei)"`
	Gas             string                        `json:"gas,omitempty" example:"21000"`
	GasPrice        string                        `json:"gasPrice,omitempty" validate:"omitempty,isBig" example:"71500000 (wei)"`
	From            string                        `json:"from" validate:"omitempty,eth_addr" example:"0x1abae27a0cbfb02945720425d3b80c7e09728534"`
	To              string                        `json:"to" validate:"required,eth_addr" example:"0x1abae27a0cbfb02945720425d3b80c7e09728534"`
	MethodSignature string                        `json:"methodSignature" validate:"required,isValidMethodSig" example:"transfer(address,uint256)"`
	Args            []interface{}                 `json:"args,omitempty"`
	OneTimeKey      bool                          `json:"oneTimeKey,omitempty" example:"true"`
	GasPricePolicy  GasPriceParams                `json:"gasPricePolicy,omitempty"`
	Protocol        entities.PrivateTxManagerType `json:"protocol,omitempty" validate:"omitempty,isPrivateTxManagerType" example:"Tessera"`
	PrivateFrom     string                        `json:"privateFrom,omitempty" validate:"omitempty,base64,required_with=PrivateFor PrivacyGroupID" example:"A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo="`
	PrivateFor      []string                      `json:"privateFor,omitempty" validate:"omitempty,min=1,unique,dive,base64" example:"A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo=,B1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo="`
	PrivacyGroupID  string                        `json:"privacyGroupId,omitempty" validate:"omitempty,base64" example:"A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo="`
}

func (params *TransactionParams) Validate() error {
	if err := utils.GetValidator().Struct(params); err != nil {
		return err
	}

	if params.PrivateFrom != "" {
		return validatePrivateTxParams(params.Protocol, params.PrivacyGroupID, params.PrivateFor)
	}

	if params.From != "" && params.OneTimeKey {
		return errors.InvalidParameterError("fields 'from' and 'oneTimeKey' are mutually exclusive")
	}

	if params.From == "" && !params.OneTimeKey {
		return errors.InvalidParameterError("field 'from' is required")
	}

	return params.GasPricePolicy.RetryPolicy.Validate()
}
