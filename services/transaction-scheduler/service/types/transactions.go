package types

import (
	"time"
)

type TransactionResponse struct {
	IdempotencyKey string                 `json:"idempotencyKey" validate:"required"`
	Params         map[string]interface{} `json:"params"`
	Schedule       *ScheduleResponse      `json:"schedule"`
	CreatedAt      time.Time              `json:"createdAt"`
}

type BaseTransactionRequest struct {
	IdempotencyKey string            `json:"idempotencyKey" validate:"required"`
	Labels         map[string]string `json:"labels,omitempty"`
}

type SendTransactionRequest struct {
	BaseTransactionRequest
	Params TransactionParams `json:"params" validate:"required"`
}

type TransferRequest struct {
	BaseTransactionRequest
	Params TransferParams `json:"params" validate:"required"`
}

type RawTransactionRequest struct {
	BaseTransactionRequest
	Method string               `json:"method" validate:"required"`
	Params RawTransactionParams `json:"params" validate:"required"`
}

type DeployContractRequest struct {
	BaseTransactionRequest
	Params DeployContractParams `json:"params" validate:"required"`
}

type PrivateTxRequest struct {
	BaseTransactionRequest
	Method string                   `json:"method" validate:"required"`
	Params PrivateTransactionParams `json:"params" validate:"required"`
}

/**
Transaction Request Param Types
*/

type BaseTransactionParams struct {
	Value    string `json:"value,omitempty" validate:"omitempty,isBig"`
	Gas      string `json:"gas,omitempty"`
	GasPrice string `json:"gasPrice,omitempty" validate:"omitempty,isBig"`
}

type TransactionParams struct {
	BaseTransactionParams
	From            string   `json:"from" validate:"required,eth_addr"`
	To              string   `json:"to" validate:"required,eth_addr"`
	MethodSignature string   `json:"methodSignature" validate:"required,isValidMethodSig"`
	Args            []string `json:"args,omitempty"`
}

type RawTransactionParams struct {
	Raw string `json:"raw" validate:"required,isHex"`
}

type TransferParams struct {
	BaseTransactionParams
	From string `json:"from" validate:"required,eth_addr"`
	To   string `json:"to" validate:"required,eth_addr"`
}

type DeployContractParams struct {
	BaseTransactionParams
	From         string   `json:"from" validate:"required,eth_addr"`
	ContractName string   `json:"contractName" validate:"required"`
	ContractTag  *string  `json:"contractTag,omitempty"`
	Args         []string `json:"args,omitempty"`
}

type PrivateTransactionParams struct {
	BaseTransactionParams
	From            string   `json:"from" validate:"required,eth_addr"`
	To              string   `json:"to" validate:"required,eth_addr"`
	MethodSignature string   `json:"methodSignature" validate:"required,isValidMethodSig"`
	Args            []string `json:"args,omitempty"`
	PrivateFrom     string   `json:"privateFrom" validate:"required,base64"`
	PrivateFor      []string `json:"privateFor" validate:"required,dive,base64"`
	PrivacyGroupID  *string  `json:"privacyGroupId,omitempty"`
}
