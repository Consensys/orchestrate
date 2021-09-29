package api

import "github.com/consensys/orchestrate/pkg/types/entities"

type RegisterContractRequest struct {
	ABI              interface{} `json:"abi,omitempty" validate:"required"`
	Bytecode         string      `json:"bytecode,omitempty" validate:"omitempty,isHex" example:"0x6080604052348015600f57600080f..."`
	DeployedBytecode string      `json:"deployedBytecode,omitempty" validate:"omitempty,isHex" example:"0x6080604052348015600f57600080f..."`
	Name             string      `json:"name" validate:"required" example:"ERC20"`
	Tag              string      `json:"tag,omitempty" example:"v1.0.0"`
}

type ContractResponse struct {
	Name             string            `json:"name" example:"ERC20"`
	Tag              string            `json:"tag" example:"v1.0.0"`
	Registry         string            `json:"registry" example:"registry.consensys.net/orchestrate"`
	ABI              string            `json:"abi" example:"[{anonymous: false, inputs: [{indexed: false, name: account, type: address}, name: MinterAdded, type: event}]}]"`
	Bytecode         string            `json:"bytecode,omitempty" example:"0x6080604052348015600f57600080f..."`
	DeployedBytecode string            `json:"deployedBytecode,omitempty" example:"0x6080604052348015600f57600080f..."`
	Constructor      entities.Method   `json:"constructor"`
	Methods          []entities.Method `json:"methods"`
	Events           []entities.Event  `json:"events"`
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
