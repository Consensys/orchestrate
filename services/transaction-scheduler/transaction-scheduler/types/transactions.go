package types

import "time"

type BaseTransactionRequest struct {
	IdempotencyKey string            `json:"idempotencyKey" validate:"required"`
	Labels         map[string]string `json:"labels,omitempty"`
}

type BaseTransactionParams struct {
	Value    string `json:"value,omitempty" validate:"omitempty,isBig"`
	Gas      string `json:"gas,omitempty"`
	GasPrice string `json:"gasPrice,omitempty" validate:"omitempty,isBig"`
}

type TransactionRequest struct {
	BaseTransactionRequest
	Params TransactionParams `json:"params" validate:"required"`
}
type TransactionParams struct {
	From            string             `json:"from" validate:"required,eth_addr"`
	To              string             `json:"to" validate:"required,eth_addr"`
	Value           string             `json:"value,omitempty" validate:"omitempty,isBig"`
	Gas             string             `json:"gas,omitempty"`
	GasPrice        string             `json:"gasPrice,omitempty" validate:"omitempty,isBig"`
	MethodSignature string             `json:"methodSignature" validate:"required,isValidMethodSig"`
	Args            *map[string]string `json:"args,omitempty"`
}

type RawTransactionRequest struct {
	BaseTransactionRequest
	Method string               `json:"method" validate:"required"`
	Params RawTransactionParams `json:"params" validate:"required"`
}
type RawTransactionParams struct {
	Raw string `json:"raw" validate:"required,isHex"`
}

type TransferRequest struct {
	BaseTransactionRequest
	Params TransferParams `json:"params" validate:"required"`
}
type TransferParams struct {
	BaseTransactionParams
	From string `json:"from" validate:"required,eth_addr"`
	To   string `json:"to" validate:"required,eth_addr"`
}

type DeployContractRequest struct {
	BaseTransactionRequest
	Params DeployContractParams `json:"params" validate:"required"`
}
type DeployContractParams struct {
	BaseTransactionParams
	From         string             `json:"from" validate:"required,eth_addr"`
	ContractName string             `json:"contractName" validate:"required"`
	ContractTag  *string            `json:"contractTag,omitempty"`
	Args         *map[string]string `json:"args,omitempty"`
}

type PrivateTransactionRequest struct {
	BaseTransactionRequest
	Method string                   `json:"method" validate:"required"`
	Params PrivateTransactionParams `json:"params" validate:"required"`
}
type PrivateTransactionParams struct {
	BaseTransactionParams
	From            string             `json:"from" validate:"required,eth_addr"`
	To              string             `json:"to" validate:"required,eth_addr"`
	MethodSignature string             `json:"methodSignature" validate:"required,isValidMethodSig"`
	Args            *map[string]string `json:"args,omitempty"`
	PrivateFrom     string             `json:"privateFrom" validate:"required,base64"`
	PrivateFor      []string           `json:"privateFor" validate:"required,dive,base64"`
	PrivacyGroupID  *string            `json:"privacyGroupId,omitempty"`
}

type TransactionResponse struct {
	IdempotencyKey string                 `json:"idempotencyKey" validate:"required"`
	Params         map[string]interface{} `json:"params"`
	Schedule       *ScheduleResponse      `json:"schedule"`
	CreatedAt      time.Time              `json:"createdAt"`
}
