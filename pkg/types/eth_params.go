package types

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
)

type ETHTransactionParams struct {
	From            string   `json:"from,omitempty" validate:"omitempty,eth_addr"`
	To              string   `json:"to,omitempty" validate:"omitempty,eth_addr"`
	Value           string   `json:"value,omitempty" validate:"omitempty,numeric"`
	GasPrice        string   `json:"gasPrice,omitempty"`
	GasLimit        string   `json:"gas,omitempty"`
	MethodSignature string   `json:"methodSignature,omitempty"`
	Args            []string `json:"args,omitempty"`
	Raw             string   `json:"raw,omitempty" validate:"omitempty,isHex"`
	ContractName    string   `json:"contractName,omitempty"`
	ContractTag     string   `json:"contractTag,omitempty"`
	Nonce           string   `json:"nonce,omitempty" validate:"omitempty,numeric"`
	PrivateTransactionParams
}

type PrivateTransactionParams struct {
	Protocol       string   `json:"protocol,omitempty" validate:"omitempty,isPrivateTxManagerType"`
	PrivateFrom    string   `json:"privateFrom,omitempty" validate:"omitempty,base64"`
	PrivateFor     []string `json:"privateFor,omitempty" validate:"omitempty,min=1,unique,dive,base64"`
	PrivacyGroupID string   `json:"privacyGroupId,omitempty" validate:"omitempty,base64"`
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
