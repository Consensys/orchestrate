package api

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
)

type RegisterContractRequest struct {
	ABI              interface{} `json:"abi,omitempty" validate:"required"`
	Bytecode         string      `json:"bytecode,omitempty" validate:"omitempty,isHex" example:"0x6080604052348015600f57600080f..."`
	DeployedBytecode string      `json:"deployedBytecode,omitempty" validate:"omitempty,isHex" example:"0x6080604052348015600f57600080f..."`
	Name             string      `json:"name" validate:"required" example:"ERC20"`
	Tag              string      `json:"tag,omitempty" example:"v1.0.0"`
}

type ContractResponse struct {
	*entities.Contract
}

type GetContractEventsRequest struct {
	SigHash           string `json:"sig_hash" validate:"required,isHex" example:"0x6080604052348015600f57600080f..."`
	IndexedInputCount uint32 `json:"indexed_input_count" validate:"omitempty" example:"1"`
}

type GetContractEventsBySignHashResponse struct {
	Event         string   `json:"event" validate:"omitempty" example:"{anonymous:false,inputs:[{indexed:true,name:from,type:address},{indexed:true,name:to,type:address},{indexed:false,name:value,type:uint256}],name:Transfer,type:event}"`
	DefaultEvents []string `json:"defaultEvents" validate:"omitempty" example:"[{anonymous:false,inputs:[{indexed:true,name:from,type:address},{indexed:true,name:to,type:address},{indexed:false,name:value,type:uint256}],name:Transfer,type:event},..."`
}

type SetContractCodeHashRequest struct {
	CodeHash string `json:"code_hash" validate:"required,isHex" example:"0x6080604052348015600f57600080f..."`
}
