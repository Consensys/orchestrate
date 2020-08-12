package types

import (
	"time"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
)

type CreateJobRequest struct {
	ScheduleUUID string            `json:"scheduleUUID" validate:"required,uuid4" example:"b4374e6f-b28a-4bad-b4fe-bda36eaf849c"`
	ChainUUID    string            `json:"chainUUID" validate:"required,uuid4" example:"b4374e6f-b28a-4bad-b4fe-bda36eaf849c"`
	Type         string            `json:"type" validate:"required,isJobType" example:"eth://ethereum/transaction"`
	Labels       map[string]string `json:"labels,omitempty"`
	Annotations  *Annotations      `json:"annotations,omitempty"`
	Transaction  *ETHTransaction   `json:"transaction" validate:"required"`
}

type UpdateJobRequest struct {
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations *Annotations      `json:"annotations,omitempty"`
	Transaction *ETHTransaction   `json:"transaction,omitempty"`
	Status      string            `json:"status,omitempty" validate:"isJobStatus" example:"MINED"`
	Message     string            `json:"message,omitempty" example:"Update message"`
}

type JobResponse struct {
	UUID         string            `json:"uuid" example:"b4374e6f-b28a-4bad-b4fe-bda36eaf849c"`
	ChainUUID    string            `json:"chainUUID" example:"b4374e6f-b28a-4bad-b4fe-bda36eaf849c"`
	ScheduleUUID string            `json:"scheduleUUID" example:"b4374e6f-b28a-4bad-b4fe-bda36eaf849c"`
	Transaction  *ETHTransaction   `json:"transaction"`
	Logs         []*Log            `json:"logs"`
	Labels       map[string]string `json:"labels,omitempty"`
	Annotations  *Annotations      `json:"annotations,omitempty"`
	Status       string            `json:"status" example:"MINED"`
	Type         string            `json:"type" example:"eth://ethereum/transaction"`
	CreatedAt    time.Time         `json:"createdAt" example:"2020-07-09T12:35:42.115395Z"`
	UpdatedAt    time.Time         `json:"updatedAt" example:"2020-07-09T12:35:42.115395Z"`
}

type CreateScheduleRequest struct{}

type ScheduleResponse struct {
	UUID      string         `json:"uuid" example:"b4374e6f-b28a-4bad-b4fe-bda36eaf849c"`
	TenantID  string         `json:"tenantID" example:"tenant_id"`
	Jobs      []*JobResponse `json:"jobs"`
	CreatedAt time.Time      `json:"createdAt" example:"2020-07-09T12:35:42.115395Z"`
}

// go validator does not support mutually exclusive parameters for now
// See more https://github.com/go-playground/validator/issues/608
type TransactionParams struct {
	Value           string               `json:"value,omitempty" validate:"omitempty,isBig" example:"71500000 (wei)"`
	Gas             string               `json:"gas,omitempty" example:"21000"`
	GasPrice        string               `json:"gasPrice,omitempty" validate:"omitempty,isBig" example:"71500000 (wei)"`
	From            string               `json:"from" validate:"required_without=OneTimeKey,omitempty,eth_addr" example:"0x1abae27a0cbfb02945720425d3b80c7e09728534"`
	To              string               `json:"to" validate:"required,eth_addr" example:"0x1abae27a0cbfb02945720425d3b80c7e09728534"`
	MethodSignature string               `json:"methodSignature" validate:"required,isValidMethodSig" example:"transfer(address,uint256)"`
	Args            []interface{}        `json:"args,omitempty"`
	OneTimeKey      bool                 `json:"oneTimeKey,omitempty" example:"true"`
	Priority        string               `json:"priority,omitempty" validate:"isPriority" example:"very-high" `
	Retry           *GasPriceRetryParams `json:"retryPolicy,omitempty"`
	Protocol        string               `json:"protocol,omitempty" validate:"omitempty,isPrivateTxManagerType" example:"Tessera"`
	PrivateFrom     string               `json:"privateFrom,omitempty" validate:"omitempty,base64" example:"A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo="`
	PrivateFor      []string             `json:"privateFor,omitempty" validate:"omitempty,min=1,unique,dive,base64" example:"A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo=,B1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo="`
	PrivacyGroupID  string               `json:"privacyGroupId,omitempty" validate:"omitempty,base64" example:"A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo="`
}

func (params *TransactionParams) Validate() error {
	if err := utils.GetValidator().Struct(params); err != nil {
		return err
	}

	if params.Retry != nil {
		if err := params.Retry.Validate(); err != nil {
			return err
		}
	}

	if params.OneTimeKey && params.From != "" {
		return errors.InvalidParameterError("from account cannot be included when OneTimeKey is enabled")
	}

	if params.PrivateFrom != "" {
		return validatePrivateTxParams(params.Protocol, params.PrivacyGroupID, params.PrivateFor)
	}

	return nil
}

type SendTransactionRequest struct {
	ChainName string            `json:"chain" validate:"required" example:"myChain"`
	Labels    map[string]string `json:"labels,omitempty"`
	Params    TransactionParams `json:"params" validate:"required"`
}

type TransferParams struct {
	Value    string               `json:"value" validate:"required,isBig" example:"71500000 (wei)"`
	Gas      string               `json:"gas,omitempty" example:"21000"`
	GasPrice string               `json:"gasPrice,omitempty" validate:"omitempty,isBig" example:"71500000 (wei)"`
	From     string               `json:"from" validate:"required,eth_addr" example:"0x1abae27a0cbfb02945720425d3b80c7e09728534"`
	To       string               `json:"to" validate:"required,eth_addr" example:"0x1abae27a0cbfb02945720425d3b80c7e09728534"`
	Priority string               `json:"priority,omitempty" validate:"isPriority" example:"very-high"`
	Retry    *GasPriceRetryParams `json:"retry,omitempty"`
}

func (params *TransferParams) Validate() error {
	if err := utils.GetValidator().Struct(params); err != nil {
		return err
	}

	if params.Retry != nil {
		if err := params.Retry.Validate(); err != nil {
			return err
		}
	}

	return nil
}

type TransferRequest struct {
	ChainName string            `json:"chain" validate:"required" example:"myChain"`
	Labels    map[string]string `json:"labels,omitempty"`
	Params    TransferParams    `json:"params" validate:"required"`
}

type RawTransactionParams struct {
	Raw   string               `json:"raw" validate:"required,isHex" example:"0xfe378324abcde723..."`
	Retry *IntervalRetryParams `json:"retryPolicy,omitempty"`
}

type RawTransactionRequest struct {
	ChainName string               `json:"chain" validate:"required" example:"myChain"`
	Labels    map[string]string    `json:"labels,omitempty"`
	Params    RawTransactionParams `json:"params" validate:"required"`
}

type DeployContractParams struct {
	Value          string               `json:"value,omitempty" validate:"omitempty,isBig" example:"71500000 (wei)"`
	Gas            string               `json:"gas,omitempty" example:"21000"`
	GasPrice       string               `json:"gasPrice,omitempty" validate:"omitempty,isBig" example:"71500000 (wei)"`
	From           string               `json:"from" validate:"required_without=OneTimeKey,omitempty,eth_addr" example:"0x1abae27a0cbfb02945720425d3b80c7e09728534"`
	ContractName   string               `json:"contractName" validate:"required" example:"MyContract"`
	ContractTag    string               `json:"contractTag,omitempty" example:"v1.1.0"`
	Args           []interface{}        `json:"args,omitempty"`
	OneTimeKey     bool                 `json:"oneTimeKey,omitempty" example:"true"`
	Priority       string               `json:"priority,omitempty" validate:"isPriority" example:"very-high" `
	Retry          *GasPriceRetryParams `json:"retryPolicy,omitempty"`
	Protocol       string               `json:"protocol,omitempty" validate:"omitempty,isPrivateTxManagerType" example:"Tessera"`
	PrivateFrom    string               `json:"privateFrom,omitempty" validate:"omitempty,base64" example:"A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo="`
	PrivateFor     []string             `json:"privateFor,omitempty" validate:"omitempty,min=1,unique,dive,base64" example:"A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo=,B1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo="`
	PrivacyGroupID string               `json:"privacyGroupId,omitempty" validate:"omitempty,base64" example:"A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo="`
}

func (params *DeployContractParams) Validate() error {
	if err := utils.GetValidator().Struct(params); err != nil {
		return err
	}

	if params.Retry != nil {
		if err := params.Retry.Validate(); err != nil {
			return err
		}
	}

	if params.OneTimeKey && params.From != "" {
		return errors.InvalidParameterError("transaction from account and OneTimeKey annotation are mutually exclusive")
	}

	if params.PrivateFrom != "" {
		return validatePrivateTxParams(params.Protocol, params.PrivacyGroupID, params.PrivateFor)
	}

	return nil
}

type DeployContractRequest struct {
	ChainName string               `json:"chain" validate:"required" example:"myChain"`
	Labels    map[string]string    `json:"labels,omitempty"`
	Params    DeployContractParams `json:"params" validate:"required"`
}

type TransactionResponse struct {
	UUID           string                `json:"uuid" example:"b4374e6f-b28a-4bad-b4fe-bda36eaf849c"`
	IdempotencyKey string                `json:"idempotencyKey" example:"myIdempotencyKey"`
	ChainName      string                `json:"chain" example:"myChain"`
	Params         *ETHTransactionParams `json:"params"`
	Schedule       *ScheduleResponse     `json:"schedule"`
	CreatedAt      time.Time             `json:"createdAt" example:"2020-07-09T12:35:42.115395Z"`
}
