package types

import "time"

type BaseTransactionRequest struct {
	IdempotencyKey string             `json:"idempotencyKey" validate:"required"`
	ChainID        string             `json:"chainID" validate:"required,uuid4"`
	Labels         *map[string]string `json:"labels,omitempty"`
}

type TransactionRequest struct {
	BaseTransactionRequest
	Params TransactionParams `json:"params" validate:"required"`
}
type TransactionParams struct {
	From            string             `json:"from" validate:"required,eth_addr"`
	To              string             `json:"to" validate:"required,eth_addr"`
	Value           *string            `json:"value,omitempty" validate:"omitempty,isBig"`
	Gas             *uint32            `json:"gas,omitempty"`
	GasPrice        *string            `json:"gasPrice,omitempty" validate:"omitempty,isBig"`
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
	From     string  `json:"from" validate:"required,eth_addr"`
	To       string  `json:"to" validate:"required,eth_addr"`
	Value    string  `json:"value" validate:"isBig"`
	Gas      *uint32 `json:"gas,omitempty"`
	GasPrice *string `json:"gasPrice,omitempty" validate:"omitempty,isBig"`
}

type DeployContractRequest struct {
	BaseTransactionRequest
	Params DeployContractParams `json:"params" validate:"required"`
}
type DeployContractParams struct {
	From         string             `json:"from" validate:"required,eth_addr"`
	Value        *string            `json:"value,omitempty" validate:"omitempty,isBig"`
	Gas          *uint32            `json:"gas,omitempty"`
	GasPrice     *string            `json:"gasPrice,omitempty" validate:"omitempty,isBig"`
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
	From            string             `json:"from" validate:"required,eth_addr"`
	To              string             `json:"to" validate:"required,eth_addr"`
	Gas             *uint32            `json:"gas,omitempty"`
	GasPrice        *string            `json:"gasPrice,omitempty" validate:"omitempty,isBig"`
	MethodSignature string             `json:"methodSignature" validate:"required,isValidMethodSig"`
	Args            *map[string]string `json:"args,omitempty"`
	PrivateFrom     string             `json:"privateFrom" validate:"required,base64"`
	PrivateFor      []string           `json:"privateFor" validate:"required,dive,base64"`
	PrivacyGroupID  *string            `json:"privacyGroupId,omitempty"`
}

type TransactionResponse struct {
	IdempotencyKey string                 `json:"idempotencyKey" validate:"required"`
	ChainID        string                 `json:"chainID" validate:"required,uuid4"`
	Labels         *map[string]string     `json:"labels,omitempty"`
	Method         string                 `json:"method"`
	Params         map[string]interface{} `json:"params"`
	Schedule       ScheduleResponse       `json:"schedule"`
	CreatedAt      time.Time              `json:"createdAt"`
}
