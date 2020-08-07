package types

import (
	"time"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
)

type CreateJobRequest struct {
	ScheduleUUID string            `json:"scheduleUUID" validate:"required,uuid4" example:"b4374e6f-b28a-4bad-b4fe-bda36eaf849c"`
	ChainUUID    string            `json:"chainUUID" validate:"required,uuid4" example:"b4374e6f-b28a-4bad-b4fe-bda36eaf849c"`
	Type         string            `json:"type" validate:"required,isJobType" example:"eth://ethereum/transaction"` //  @TODO validate Type is valid
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
}

type CreateScheduleRequest struct{}

type ScheduleResponse struct {
	UUID      string         `json:"uuid" example:"b4374e6f-b28a-4bad-b4fe-bda36eaf849c"`
	TenantID  string         `json:"tenantID" example:"tenant_id"`
	Jobs      []*JobResponse `json:"jobs"`
	CreatedAt time.Time      `json:"createdAt" example:"2020-07-09T12:35:42.115395Z"`
}

type BaseTransactionRequest struct {
	ChainName string            `json:"chain" validate:"required" example:"myChain"`
	Labels    map[string]string `json:"labels,omitempty"`
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
	Retry           *GasPriceRetryParams `json:"retry,omitempty"`
	PrivateTransactionParams
}

func (t *TransactionParams) Validate() error {
	if t == nil {
		return nil
	}

	if err := utils.GetValidator().Struct(t); err != nil {
		return err
	}

	if err := t.PrivateTransactionParams.Validate(); err != nil {
		return err
	}

	if t.Retry != nil {
		if err := t.Retry.Validate(); err != nil {
			return err
		}
	}

	return nil
}

type SendTransactionRequest struct {
	BaseTransactionRequest
	Params TransactionParams `json:"params" validate:"required"`
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

func (t *TransferParams) Validate() error {
	if t == nil {
		return nil
	}

	if err := utils.GetValidator().Struct(t); err != nil {
		return err
	}

	if t.Retry != nil {
		if err := t.Retry.Validate(); err != nil {
			return err
		}
	}

	return nil
}

type TransferRequest struct {
	BaseTransactionRequest
	Params TransferParams `json:"params" validate:"required"`
}

type RawTransactionParams struct {
	Raw   string           `json:"raw" validate:"required,isHex" example:"0xfe378324abcde723..."`
	Retry *BaseRetryParams `json:"retry,omitempty"`
}

type BaseRetryParams struct {
	Interval string `json:"interval,omitempty" validate:"omitempty,isDuration" example:"2m"`
}

type RawTransactionRequest struct {
	BaseTransactionRequest
	Params RawTransactionParams `json:"params" validate:"required"`
}

type DeployContractParams struct {
	Value        string               `json:"value,omitempty" validate:"omitempty,isBig" example:"71500000 (wei)"`
	Gas          string               `json:"gas,omitempty" example:"21000"`
	GasPrice     string               `json:"gasPrice,omitempty" validate:"omitempty,isBig" example:"71500000 (wei)"`
	From         string               `json:"from" validate:"required_without=OneTimeKey,omitempty,eth_addr" example:"0x1abae27a0cbfb02945720425d3b80c7e09728534"`
	ContractName string               `json:"contractName" validate:"required" example:"MyContract"`
	ContractTag  string               `json:"contractTag,omitempty" example:"v1.1.0"`
	Args         []interface{}        `json:"args,omitempty"`
	OneTimeKey   bool                 `json:"oneTimeKey,omitempty" example:"true"`
	Priority     string               `json:"priority,omitempty" validate:"isPriority" example:"very-high" `
	Retry        *GasPriceRetryParams `json:"retry,omitempty"`
	PrivateTransactionParams
}

func (d *DeployContractParams) Validate() error {
	if d == nil {
		return nil
	}

	if err := utils.GetValidator().Struct(d); err != nil {
		return err
	}

	if err := d.PrivateTransactionParams.Validate(); err != nil {
		return err
	}

	if d.Retry != nil {
		if err := d.Retry.Validate(); err != nil {
			return err
		}
	}

	return nil
}

type DeployContractRequest struct {
	BaseTransactionRequest
	Params DeployContractParams `json:"params" validate:"required"`
}

type TransactionResponse struct {
	UUID           string                `json:"uuid" example:"b4374e6f-b28a-4bad-b4fe-bda36eaf849c"`
	IdempotencyKey string                `json:"idempotencyKey" example:"myIdempotencyKey"`
	ChainName      string                `json:"chain" example:"myChain"`
	Params         *ETHTransactionParams `json:"params"`
	Schedule       *ScheduleResponse     `json:"schedule"`
	CreatedAt      time.Time             `json:"createdAt" example:"2020-07-09T12:35:42.115395Z"`
}

type GasPriceRetryParams struct {
	BaseRetryParams
	GasPriceIncrementLevel string  `json:"gasPriceIncrementLevel,omitempty" validate:"omitempty,oneof=none very-low low medium high very-high" example:"medium"`
	GasPriceIncrement      float64 `json:"gasPriceIncrement,omitempty" validate:"omitempty,eqfield=GasPriceLimit|ltcsfield=GasPriceLimit" example:"1.05 for (5%)"`
	GasPriceLimit          float64 `json:"gasPriceLimit,omitempty" validate:"required_with=GasPriceIncrementLevel GasPriceIncrement,omitempty" example:"1.2 for (20%)"`
}

func (g *GasPriceRetryParams) Validate() error {
	if g == nil {
		return nil
	}

	if err := utils.GetValidator().Struct(g); err != nil {
		return err
	}

	if g.GasPriceIncrement > 0 && g.GasPriceIncrementLevel != "" {
		return errors.InvalidParameterError("fields 'gasPriceIncrement' and 'gasPriceIncrementLevel' are mutually exclusive")
	}

	if (g.Interval != "" || g.GasPriceLimit > 0) && (g.GasPriceIncrement == 0 && g.GasPriceIncrementLevel == "") {
		return errors.InvalidParameterError("fields 'gasPriceIncrement' and 'gasPriceIncrementLevel' cannot be both empty")
	}

	if (g.Interval != "" && g.GasPriceLimit == 0) || (g.Interval == "" && g.GasPriceLimit > 0) {
		return errors.InvalidParameterError("fields 'Interval' and 'GasPriceLimit' are both required")
	}

	return nil
}
