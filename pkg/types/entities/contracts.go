package entities

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

const DefaultTagValue = "latest"

type Contract struct {
	Name             string
	Tag              string
	Registry         string
	ABI              abi.ABI
	RawABI           string
	Bytecode         hexutil.Bytes
	DeployedBytecode hexutil.Bytes
	Constructor      ABIComponent
	Methods          []ABIComponent
	Events           []ABIComponent
}

type ABIComponent struct {
	Signature string `json:"signature" example:"transfer(address,uint256)"`
	ABI       string `json:"abi,omitempty" example:"[{anonymous: false, inputs: [{indexed: false, name: account, type: address}, name: MinterAdded, type: event}]}]"`
}

type Arguments struct {
	Name    string
	Type    string
	Indexed bool
}

type RawABI struct {
	Type      string
	Name      string
	Constant  bool
	Anonymous bool
	Inputs    []Arguments
	Outputs   []Arguments
}

func (c *Contract) String() string {
	tag := DefaultTagValue
	if c.Tag != "" {
		tag = c.Tag
	}

	return fmt.Sprintf("%v[%v]", c.Name, tag)
}
