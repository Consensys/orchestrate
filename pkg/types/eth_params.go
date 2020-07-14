package types

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
)

type ETHTransactionParams struct {
	From            string        `json:"from,omitempty" validate:"omitempty,eth_addr" example:"0x1abae27a0cbfb02945720425d3b80c7e09728534"`
	To              string        `json:"to,omitempty" validate:"omitempty,eth_addr" example:"0x1abae27a0cbfb02945720425d3b80c7e09728534"`
	Value           string        `json:"value,omitempty" validate:"omitempty,numeric" example:"71500000 (wei)"`
	GasPrice        string        `json:"gasPrice,omitempty" example:"71500000 (wei)"`
	Gas             string        `json:"gas,omitempty" example:"21000"`
	MethodSignature string        `json:"methodSignature,omitempty" example:"transfer(address,uint256)"`
	Args            []interface{} `json:"args,omitempty"`
	Raw             string        `json:"raw,omitempty" validate:"omitempty,isHex" example:"0xfe378324abcde723"`
	ContractName    string        `json:"contractName,omitempty" example:"MyContract"`
	ContractTag     string        `json:"contractTag,omitempty" example:"v1.1.0"`
	Nonce           string        `json:"nonce,omitempty" validate:"omitempty,numeric" example:"1"`
	PrivateTransactionParams
}

type PrivateTransactionParams struct {
	Protocol       string   `json:"protocol,omitempty" validate:"omitempty,isPrivateTxManagerType" example:"Tessera"`
	PrivateFrom    string   `json:"privateFrom,omitempty" validate:"omitempty,base64" example:"A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo="`
	PrivateFor     []string `json:"privateFor,omitempty" validate:"omitempty,min=1,unique,dive,base64" example:"A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo=,B1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo="`
	PrivacyGroupID string   `json:"privacyGroupId,omitempty" validate:"omitempty,base64" example:"A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo="`
}

func (tx *PrivateTransactionParams) Validate() error {
	if err := utils.GetValidator().Struct(tx); err != nil {
		return err
	}

	if tx.Protocol == "" {
		return nil
	}

	if tx.PrivateFrom == "" {
		return errors.InvalidParameterError("fields 'privateFrom' cannot be empty")
	}

	if len(tx.PrivateFor) == 0 && tx.PrivacyGroupID == "" {
		return errors.InvalidParameterError("fields 'privacyGroupId' and 'privateFor' cannot be both empty")
	}

	if len(tx.PrivateFor) > 0 && tx.PrivacyGroupID != "" {
		return errors.InvalidParameterError("fields 'privacyGroupId' and 'privateFor' are mutually exclusive")
	}

	return nil
}
