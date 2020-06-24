package types

import (
	"time"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types"
)

type TransactionResponse struct {
	UUID           string                      `json:"uuid"`
	IdempotencyKey string                      `json:"idempotencyKey"`
	Params         *types.ETHTransactionParams `json:"params"`
	ChainName      string                      `json:"chain"`
	Schedule       *ScheduleResponse           `json:"schedule"`
	CreatedAt      time.Time                   `json:"createdAt"`
}

type BaseTransactionRequest struct {
	ChainName string            `json:"chain" validate:"required"`
	Labels    map[string]string `json:"labels,omitempty"`
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
	Params RawTransactionParams `json:"params" validate:"required"`
}

type DeployContractRequest struct {
	BaseTransactionRequest
	Params DeployContractParams `json:"params" validate:"required"`
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
	types.PrivateTransactionParams
}

type RawTransactionParams struct {
	Raw string `json:"raw" validate:"required,isHex"`
}

type TransferParams struct {
	BaseTransactionParams
	From  string `json:"from" validate:"required,eth_addr"`
	To    string `json:"to" validate:"required,eth_addr"`
	Value string `json:"value" validate:"required,isBig"`
}

type DeployContractParams struct {
	BaseTransactionParams
	From         string   `json:"from" validate:"required,eth_addr"`
	ContractName string   `json:"contractName" validate:"required"`
	ContractTag  string   `json:"contractTag,omitempty"`
	Args         []string `json:"args,omitempty"`
	types.PrivateTransactionParams
}
