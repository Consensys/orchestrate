package ethereum

import (
	"github.com/ConsenSys/orchestrate/pkg/errors"
)

type SignETHTransactionRequest struct {
	Namespace string `json:"namespace,omitempty" example:"tenant_id"`
	Nonce     uint64 `json:"nonce" example:"1"`
	To        string `json:"to,omitempty" validate:"isHex" example:"0x905B88EFf8Bda1543d4d6f4aA05afef143D27E18"`
	Amount    string `json:"amount,omitempty" validate:"isBig" example:"100000000000"`
	GasPrice  string `json:"gasPrice" validate:"required,isBig" example:"100000000000"`
	GasLimit  uint64 `json:"gasLimit" validate:"required" example:"21000"`
	Data      string `json:"data,omitempty" validate:"isHex" example:"0xfeaeee..."`
	ChainID   string `json:"chainID" validate:"required,isBig" example:"1 (mainnet)"`
}

type SignQuorumPrivateTransactionRequest struct {
	Namespace string `json:"namespace,omitempty" example:"tenant_id"`
	Nonce     uint64 `json:"nonce" example:"1"`
	To        string `json:"to,omitempty" validate:"isHex" example:"0x905B88EFf8Bda1543d4d6f4aA05afef143D27E18"`
	Amount    string `json:"amount,omitempty" validate:"isBig" example:"100000000000"`
	GasPrice  string `json:"gasPrice" validate:"required,isBig" example:"100000000000"`
	GasLimit  uint64 `json:"gasLimit" validate:"required" example:"21000"`
	Data      string `json:"data,omitempty" validate:"isHex" example:"0xfeaeee..."`
}

type SignEEATransactionRequest struct {
	Namespace      string   `json:"namespace,omitempty" example:"tenant_id"`
	Nonce          uint64   `json:"nonce" example:"1"`
	To             string   `json:"to,omitempty" validate:"isHex" example:"0x905B88EFf8Bda1543d4d6f4aA05afef143D27E18"`
	Data           string   `json:"data,omitempty" validate:"isHex" example:"0xfeaeee..."`
	ChainID        string   `json:"chainID" validate:"required,isBig" example:"1 (mainnet)"`
	PrivateFrom    string   `json:"privateFrom" validate:"required,base64,required_with=PrivateFor PrivacyGroupID" example:"A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo="`
	PrivateFor     []string `json:"privateFor,omitempty" validate:"omitempty,min=1,unique,dive,base64" example:"A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo=,B1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo="`
	PrivacyGroupID string   `json:"privacyGroupId,omitempty" validate:"omitempty,base64" example:"A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo="`
}

type SignTypedDataRequest struct {
	Namespace       string                 `json:"namespace,omitempty" example:"tenant_id"`
	DomainSeparator DomainSeparator        `json:"domainSeparator" validate:"required"`
	Types           map[string][]Type      `json:"types" validate:"required"`
	Message         map[string]interface{} `json:"message" validate:"required"`
	MessageType     string                 `json:"messageType" validate:"required" example:"Mail"`
}

type DomainSeparator struct {
	Name              string `json:"name" validate:"required" example:"MyDApp"`
	Version           string `json:"version" validate:"required" example:"v1.0.0"`
	ChainID           int64  `json:"chainID" validate:"required" example:"1"`
	VerifyingContract string `json:"verifyingContract,omitempty" validate:"omitempty,isHex" example:"0x905B88EFf8Bda1543d4d6f4aA05afef143D27E18"`
	Salt              string `json:"salt,omitempty" validate:"omitempty" example:"some-random-string"`
}

type Type struct {
	Name string `json:"name" validate:"required" example:"fieldName"`
	Type string `json:"type" validate:"required" example:"string"`
}

type VerifyTypedDataRequest struct {
	TypedData SignTypedDataRequest `json:"data" validate:"required"`
	Signature string               `json:"signature" validate:"required,isHex" example:"0x6019a3c8..."`
	Address   string               `json:"address" validate:"required,isHex" example:"0x905B88EFf8Bda1543d4d6f4aA05afef143D27E18"`
}

func (req *SignEEATransactionRequest) Validate() error {
	if len(req.PrivateFor) > 0 && req.PrivacyGroupID != "" {
		return errors.InvalidFormatError("privacyGroupId and privateFor fields are mutually exclusive")
	}

	return nil
}
