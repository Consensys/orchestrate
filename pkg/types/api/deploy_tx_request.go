package api

import (
	"github.com/consensys/orchestrate/pkg/types/entities"
	"github.com/consensys/orchestrate/pkg/utils"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
)

type DeployContractRequest struct {
	ChainName string               `json:"chain" validate:"required" example:"myChain"`
	Labels    map[string]string    `json:"labels,omitempty"`
	Params    DeployContractParams `json:"params" validate:"required"`
}

type DeployContractParams struct {
	Value           *hexutil.Big                  `json:"value,omitempty" validate:"omitempty" example:"0x59682f00" swaggertype:"string"`
	Gas             *uint64                       `json:"gas,omitempty" example:"300000"`
	GasPrice        *hexutil.Big                  `json:"gasPrice,omitempty" validate:"omitempty" example:"0x5208" swaggertype:"string"`
	GasFeeCap       *hexutil.Big                  `json:"maxFeePerGas,omitempty" example:"0x4c4b40" swaggertype:"string"`
	GasTipCap       *hexutil.Big                  `json:"maxPriorityFeePerGas,omitempty" example:"0x59682f00" swaggertype:"string"`
	AccessList      types.AccessList              `json:"accessList,omitempty" swaggertype:"array,object"`
	TransactionType string                        `json:"transactionType,omitempty" validate:"omitempty,isTransactionType" example:"dynamic_fee" enums:"legacy,dynamic_fee"`
	From            *ethcommon.Address            `json:"from" validate:"omitempty" example:"0x1abae27a0cbfb02945720425d3b80c7e09728534" swaggertype:"string"`
	ContractName    string                        `json:"contractName" validate:"required" example:"MyContract"`
	ContractTag     string                        `json:"contractTag,omitempty" example:"v1.1.0"`
	Args            []interface{}                 `json:"args,omitempty"`
	OneTimeKey      bool                          `json:"oneTimeKey,omitempty" example:"true"`
	GasPricePolicy  GasPriceParams                `json:"gasPricePolicy,omitempty"`
	Protocol        entities.PrivateTxManagerType `json:"protocol,omitempty" validate:"omitempty,isPrivateTxManagerType" example:"Tessera"`
	PrivateFrom     string                        `json:"privateFrom,omitempty" validate:"omitempty,base64" example:"A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo="`
	PrivateFor      []string                      `json:"privateFor,omitempty" validate:"omitempty,min=1,unique,dive,base64" example:"[A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo=,B1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo=]"`
	MandatoryFor    []string                      `json:"mandatoryFor,omitempty" validate:"omitempty,min=1,unique,dive,base64" example:"[A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo=,B1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo=]"`
	PrivacyGroupID  string                        `json:"privacyGroupId,omitempty" validate:"omitempty,base64" example:"A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo="`
	PrivacyFlag     int                           `json:"privacyFlag,omitempty" validate:"omitempty,isPrivacyFlag" example:"0"`
}

func (params *DeployContractParams) Validate() error {
	if err := utils.GetValidator().Struct(params); err != nil {
		return err
	}

	if params.Protocol != "" || params.PrivateFrom != "" {
		return validatePrivateTxParams(params.Protocol, params.PrivateFrom, params.PrivacyGroupID, params.PrivateFor)
	}

	if err := validateTxFromParams(params.From, params.OneTimeKey); err != nil {
		return err
	}

	return params.GasPricePolicy.RetryPolicy.Validate()
}
