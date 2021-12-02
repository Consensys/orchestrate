package entities

import (
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
)

type ETHTransactionParams struct {
	From            *ethcommon.Address   `json:"from,omitempty"  example:"0x1abae27a0cbfb02945720425d3b80c7e09728534" swaggertype:"string"`
	To              *ethcommon.Address   `json:"to,omitempty" example:"0x1abae27a0cbfb02945720425d3b80c7e09728534" swaggertype:"string"`
	Value           *hexutil.Big         `json:"value,omitempty" validate:"omitempty" example:"0x59682f00" swaggertype:"string"`
	Gas             *uint64              `json:"gas,omitempty" example:"21000"`
	GasPrice        *hexutil.Big         `json:"gasPrice,omitempty" validate:"omitempty" example:"0x5208" swaggertype:"string"`
	GasFeeCap       *hexutil.Big         `json:"maxFeePerGas,omitempty" example:"0x4c4b40" swaggertype:"string"`
	GasTipCap       *hexutil.Big         `json:"maxPriorityFeePerGas,omitempty" example:"0x59682f00" swaggertype:"string"`
	AccessList      types.AccessList     `json:"accessList,omitempty" swaggertype:"array,object"`
	TransactionType string               `json:"transactionType,omitempty" example:"dynamic_fee" enums:"legacy,dynamic_fee"`
	MethodSignature string               `json:"methodSignature,omitempty" example:"transfer(address,uint256)"`
	Args            []interface{}        `json:"args,omitempty"`
	Raw             hexutil.Bytes        `json:"raw,omitempty" example:"0xfe378324abcde723"`
	ContractName    string               `json:"contractName,omitempty" example:"MyContract"`
	ContractTag     string               `json:"contractTag,omitempty" example:"v1.1.0"`
	Nonce           *uint64              `json:"nonce,omitempty" example:"1"`
	Protocol        PrivateTxManagerType `json:"protocol,omitempty" example:"Tessera"`
	PrivateFrom     string               `json:"privateFrom,omitempty" example:"A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo="`
	PrivateFor      []string             `json:"privateFor,omitempty" validate:"omitempty,min=1,unique,dive,base64" example:"[A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo=,B1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo=]"`
	MandatoryFor    []string             `json:"mandatoryFor,omitempty" validate:"omitempty,min=1,unique,dive,base64" example:"[A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo=,B1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo=]"`
	PrivacyGroupID  string               `json:"privacyGroupId,omitempty" example:"A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo="`
	PrivacyFlag     PrivacyFlag          `json:"privacyFlag,omitempty" validate:"omitempty,isPrivacyFlag" example:"0"`
}

type PrivateETHTransactionParams struct {
	PrivateFrom    string
	PrivateFor     []string
	PrivacyGroupID string
	PrivateTxType  PrivateTxType
}
