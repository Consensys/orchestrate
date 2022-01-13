package entities

import (
	"encoding/json"
	"time"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
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
	Hash            *ethcommon.Hash    `json:"hash,omitempty" validate:"omitempty" example:"0xd41551c714c8ec769d2edad9adc250ae955d263da161bf59142b7500eea6715e" swaggertype:"string"`
	From            *ethcommon.Address `json:"from,omitempty" validate:"omitempty" example:"0x1abae27a0cbfb02945720425d3b80c7e09728534" swaggertype:"string"`
	To              *ethcommon.Address `json:"to,omitempty" validate:"omitempty" example:"0x1abae27a0cbfb02945720425d3b80c7e09728534" swaggertype:"string"`
	Nonce           *uint64            `json:"nonce,omitempty" validate:"omitempty" example:"1"`
	Value           *hexutil.Big       `json:"value,omitempty" validate:"omitempty" example:"0x59682f00" swaggertype:"string"`
	Gas             *uint64            `json:"gas,omitempty" example:"21000"`
	GasPrice        *hexutil.Big       `json:"gasPrice,omitempty" validate:"omitempty" example:"0x5208" swaggertype:"string"`
	GasFeeCap       *hexutil.Big       `json:"maxFeePerGas,omitempty" example:"0x4c4b40" swaggertype:"string"`
	GasTipCap       *hexutil.Big       `json:"maxPriorityFeePerGas,omitempty" example:"0x59682f00" swaggertype:"string"`
	AccessList      types.AccessList   `json:"accessList,omitempty" swaggertype:"array,object"`
	TransactionType TransactionType    `json:"transactionType,omitempty" example:"dynamic_fee" enums:"legacy,dynamic_fee"`
	Data            hexutil.Bytes      `json:"data,omitempty" validate:"omitempty" example:"0xfe378324abcde723" swaggertype:"string"`
	Raw             hexutil.Bytes      `json:"raw,omitempty" validate:"omitempty" example:"0xfe378324abcde723" swaggertype:"string"`
	PrivateFrom     string             `json:"privateFrom,omitempty" validate:"omitempty,base64" example:"A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo="`
	PrivateFor      []string           `json:"privateFor,omitempty" validate:"omitempty,min=1,unique,dive,base64" example:"[A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo=,B1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo=]"`
	MandatoryFor    []string           `json:"mandatoryFor,omitempty" validate:"omitempty,min=1,unique,dive,base64" example:"[A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo=,B1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo=]"`
	PrivacyGroupID  string             `json:"privacyGroupId,omitempty" validate:"omitempty,base64" example:"A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo="`
	PrivacyFlag     PrivacyFlag        `json:"privacyFlag,omitempty" validate:"omitempty,isPrivacyFlag" example:"1"`
	EnclaveKey      hexutil.Bytes      `json:"enclaveKey,omitempty" example:"0xd41551c714c8ec769d2edad9adc250ae955d263da161bf59142b7500eea6715eadc250ae955d263da161bf59142b7500eea6715e" swaggertype:"string"`
	CreatedAt       time.Time          `json:"createdAt,omitempty" example:"2020-07-09T12:35:42.115395Z"`
	UpdatedAt       time.Time          `json:"updatedAt,omitempty" example:"2020-07-09T12:35:42.115395Z"`
}

type ethTransactionJSON struct {
	Hash            string            `json:"hash,omitempty"`
	From            string            `json:"from,omitempty"`
	To              string            `json:"to,omitempty"`
	Nonce           *uint64           `json:"nonce,omitempty"`
	Value           string            `json:"value,omitempty"`
	Gas             *uint64           `json:"gas,omitempty"`
	GasPrice        string            `json:"gasPrice,omitempty"`
	GasFeeCap       string            `json:"maxFeePerGas,omitempty"`
	GasTipCap       string            `json:"maxPriorityFeePerGas,omitempty"`
	AccessList      []accessTupleJSON `json:"accessList,omitempty"`
	TransactionType string            `json:"transactionType,omitempty"`
	Data            string            `json:"data,omitempty"`
	Raw             string            `json:"raw,omitempty"`
	PrivateFrom     string            `json:"privateFrom,omitempty"`
	PrivateFor      []string          `json:"privateFor,omitempty"`
	MandatoryFor    []string          `json:"mandatoryFor,omitempty"`
	PrivacyGroupID  string            `json:"privacyGroupId,omitempty"`
	PrivacyFlag     int               `json:"privacyFlag,omitempty"`
	EnclaveKey      string            `json:"enclaveKey,omitempty"`
	CreatedAt       time.Time         `json:"createdAt,omitempty"`
	UpdatedAt       time.Time         `json:"updatedAt,omitempty"`
}

type accessTupleJSON struct {
	Address     string   `json:"address"       `
	StorageKeys []string `json:"storageKeys"   `
}

func (tx *ETHTransaction) MarshalJSON() ([]byte, error) {
	r := &ethTransactionJSON{
		Nonce:           tx.Nonce,
		TransactionType: string(tx.TransactionType),
		Gas:             tx.Gas,
		PrivateFrom:     tx.PrivateFrom,
		PrivateFor:      tx.PrivateFor,
		MandatoryFor:    tx.MandatoryFor,
		PrivacyGroupID:  tx.PrivacyGroupID,
		PrivacyFlag:     int(tx.PrivacyFlag),
		CreatedAt:       tx.CreatedAt,
		UpdatedAt:       tx.UpdatedAt,
	}

	if tx.Hash != nil {
		r.Hash = tx.Hash.String()
	}

	if tx.From != nil {
		r.From = tx.From.String()
	}

	if tx.To != nil {
		r.To = tx.To.String()
	}

	if tx.Value != nil {
		r.Value = tx.Value.String()
	}

	if tx.GasPrice != nil {
		r.GasPrice = tx.GasPrice.String()
	}

	if tx.GasFeeCap != nil {
		r.GasFeeCap = tx.GasFeeCap.String()
	}

	if tx.GasTipCap != nil {
		r.GasTipCap = tx.GasTipCap.String()
	}

	if len(tx.Data) != 0 {
		r.Data = tx.Data.String()
	}

	if len(tx.Raw) != 0 {
		r.Raw = tx.Raw.String()
	}

	if len(tx.EnclaveKey) != 0 {
		r.EnclaveKey = tx.EnclaveKey.String()
	}

	r.AccessList = []accessTupleJSON{}
	for _, accessListItem := range tx.AccessList {
		elem := accessTupleJSON{
			Address:     accessListItem.Address.String(),
			StorageKeys: []string{},
		}
		for _, storeKey := range accessListItem.StorageKeys {
			elem.StorageKeys = append(elem.StorageKeys, storeKey.String())
		}
		r.AccessList = append(r.AccessList, elem)
	}

	return json.Marshal(r)
}
