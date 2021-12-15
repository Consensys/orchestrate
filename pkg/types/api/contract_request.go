package api

import (
	"github.com/consensys/orchestrate/pkg/types/entities"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

type RegisterContractRequest struct {
	ABI              interface{}   `json:"abi,omitempty" validate:"required"`
	Bytecode         hexutil.Bytes `json:"bytecode,omitempty" validate:"omitempty" example:"0x6080604052348015600f57600080f" swaggertype:"string"`
	DeployedBytecode hexutil.Bytes `json:"deployedBytecode,omitempty" validate:"omitempty" example:"0x6080604052348015600f57600080f" swaggertype:"string"`
	Name             string        `json:"name" validate:"required" example:"ERC20"`
	Tag              string        `json:"tag,omitempty" example:"v1.0.0"`
}

type ContractResponse struct {
	Name             string            `json:"name" example:"ERC20"`
	Tag              string            `json:"tag" example:"v1.0.0"`
	Registry         string            `json:"registry" example:"registry.consensys.net/orchestrate"`
	ABI              string            `json:"abi" example:"[{anonymous: false, inputs: [{indexed: false, name: account, type: address}, name: MinterAdded, type: event}]}]"`
	Bytecode         hexutil.Bytes     `json:"bytecode,omitempty" example:"0x6080604052348015600f57600080f..." swaggertype:"string"`
	DeployedBytecode hexutil.Bytes     `json:"deployedBytecode,omitempty" example:"0x6080604052348015600f57600080f..." swaggertype:"string"`
	Constructor      entities.Method   `json:"constructor"`
	Methods          []entities.Method `json:"methods"`
	Events           []entities.Event  `json:"events"`
}

type GetContractEventsRequest struct {
	SigHash           hexutil.Bytes `json:"sig_hash" validate:"required" example:"0x6080604052348015600f57600080f" swaggertype:"string"`
	IndexedInputCount uint32        `json:"indexed_input_count" validate:"omitempty" example:"1"`
}

type GetContractEventsBySignHashResponse struct {
	Event         string   `json:"event" validate:"omitempty" example:"{anonymous:false,inputs:[{indexed:true,name:from,type:address},{indexed:true,name:to,type:address},{indexed:false,name:value,type:uint256}],name:Transfer,type:event}"`
	DefaultEvents []string `json:"defaultEvents" validate:"omitempty" example:"[{anonymous:false,inputs:[{indexed:true,name:from,type:address},{indexed:true,name:to,type:address},{indexed:false,name:value,type:uint256}],name:Transfer,type:event},..."`
}

type SetContractCodeHashRequest struct {
	CodeHash hexutil.Bytes `json:"code_hash" validate:"required" example:"0x6080604052348015600f57600080f" swaggertype:"string"`
}

type SearchContractRequest struct {
	CodeHash hexutil.Bytes      `json:"code_hash" validate:"required" example:"0x6080604052348015600f57600080f" swaggertype:"string"`
	Address  *ethcommon.Address `json:"address" validate:"required" example:"0x1abae27a0cbfb02945720425d3b80c7e09728534" swaggertype:"string"`
}
