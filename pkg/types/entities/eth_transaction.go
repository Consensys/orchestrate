package entities

import (
	"time"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type TransactionType string

const (
	LegacyTxType     TransactionType = "legacy"
	DynamicFeeTxType TransactionType = "dynamic_fee"
)

type PrivacyFlag int

const (
	PrivacyFlagSP  PrivacyFlag = iota
	PrivacyFlagPP  PrivacyFlag = 1
	PrivacyFlagMPP PrivacyFlag = 2
	PrivacyFlagPSV PrivacyFlag = 3
)

type ETHTransaction struct {
	Hash            string           `json:"hash,omitempty" validate:"omitempty,isHex" example:"0xd41551c714c8ec769d2edad9adc250ae955d263da161bf59142b7500eea6715e"`
	From            string           `json:"from,omitempty" validate:"omitempty,eth_addr" example:"0x1abae27a0cbfb02945720425d3b80c7e09728534"`
	To              string           `json:"to,omitempty" validate:"omitempty,eth_addr" example:"0x1abae27a0cbfb02945720425d3b80c7e09728534"`
	Nonce           string           `json:"nonce,omitempty" validate:"omitempty,isBig" example:"1"`
	Value           string           `json:"value,omitempty" validate:"omitempty,isBig" example:"71500000 (wei)"`
	GasPrice        string           `json:"gasPrice,omitempty" validate:"omitempty,isBig" example:"71500000 (wei)"`
	GasFeeCap       string           `json:"maxFeePerGas,omitempty" example:"71500000 (wei)"`
	GasTipCap       string           `json:"maxPriorityFeePerGas,omitempty" example:"71500000 (wei)"`
	AccessList      types.AccessList `json:"accessList,omitempty" swaggertype:"array,object"`
	TransactionType TransactionType  `json:"transactionType,omitempty" example:"dynamic_fee" enums:"legacy,dynamic_fee"`
	Gas             string           `json:"gas,omitempty" example:"21000"`
	Data            string           `json:"data,omitempty" validate:"omitempty,isHex" example:"0xfe378324abcde723"`
	Raw             string           `json:"raw,omitempty" validate:"omitempty,isHex" example:"0xfe378324abcde723"`
	PrivateFrom     string           `json:"privateFrom,omitempty" validate:"omitempty,base64" example:"A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo="`
	PrivateFor      []string         `json:"privateFor,omitempty" validate:"omitempty,min=1,unique,dive,base64" example:"[A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo=,B1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo=]"`
	MandatoryFor    []string         `json:"mandatoryFor,omitempty" validate:"omitempty,min=1,unique,dive,base64" example:"[A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo=,B1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo=]"`
	PrivacyGroupID  string           `json:"privacyGroupId,omitempty" validate:"omitempty,base64" example:"A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo="`
	PrivacyFlag     PrivacyFlag      `json:"privacyFlag,omitempty" validate:"omitempty,isPrivacyFlag" example:"A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo="`
	EnclaveKey      string           `json:"enclaveKey,omitempty" example:"0xd41551c714c8ec769d2edad9adc250ae955d263da161bf59142b7500eea6715eadc250ae955d263da161bf59142b7500eea6715e"`
	CreatedAt       time.Time        `json:"createdAt,omitempty" example:"2020-07-09T12:35:42.115395Z"`
	UpdatedAt       time.Time        `json:"updatedAt,omitempty" example:"2020-07-09T12:35:42.115395Z"`
}

func (t *ETHTransaction) GetHash() ethcommon.Hash {
	return ethcommon.HexToHash(t.Hash)
}
