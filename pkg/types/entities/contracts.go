package entities

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	ethabi "github.com/ethereum/go-ethereum/accounts/abi"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
)

const DefaultTagValue = "latest"

type Contract struct {
	Name             string   `json:"name" example:"ERC20"`
	Tag              string   `json:"tag" example:"v1.0.0"`
	Registry         string   `json:"registry" example:"registry.consensys.net/orchestrate"`
	ABI              string   `json:"abi" example:"[{anonymous: false, inputs: [{indexed: false, name: account, type: address}, name: MinterAdded, type: event}]}]"`
	Bytecode         string   `json:"bytecode,omitempty" example:"0x6080604052348015600f57600080f..."`
	DeployedBytecode string   `json:"deployedBytecode,omitempty" example:"0x6080604052348015600f57600080f..."`
	Constructor      Method   `json:"constructor"`
	Methods          []Method `json:"methods"`
	Events           []Event  `json:"events"`
}

type Method struct {
	Signature string `json:"signature" example:"transfer(address,uint256)"`
	ABI       string `json:"abi" example:"[{anonymous: false, inputs: [{indexed: false, name: account, type: address}, name: MinterAdded, type: event}]}]"`
}

type Event struct {
	Signature string `json:"signature" example:"transfer(address,uint256)"`
	ABI       string `json:"abi" example:"[{anonymous: false, inputs: [{indexed: false, name: account, type: address}, name: MinterAdded, type: event}]}]"`
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

type Artifact struct {
	Abi              string
	Bytecode         string
	DeployedBytecode string
}

// Short returns a short string representation of contractId information
func (c *Contract) Short() string {
	if c.Name == "" {
		return ""
	}

	if c.Tag == "" {
		return fmt.Sprintf("%v[%v]", c.Name, DefaultTagValue)
	}

	return fmt.Sprintf("%v[%v]", c.Name, c.Tag)
}

// Long return a long string representation of contract information
func (c *Contract) Long() string {
	return fmt.Sprintf("%v:%v:%v", c.Short(), c.ABI, c.Bytecode)
}

// GetABICompacted returns a compacted version of the ABI
func (c *Contract) GetABICompacted() (string, error) {
	buffer := new(bytes.Buffer)
	if err := json.Compact(buffer, []byte(c.ABI)); err != nil {
		return "", err
	}
	return buffer.String(), nil
}

// CompactABI compact inplace the ABI
func (c *Contract) CompactABI() error {
	compactedABI, err := c.GetABICompacted()
	if err != nil {
		return err
	}
	c.ABI = compactedABI
	return nil
}

// IsConstructor indicate whether the method refers to a deployment
func (m *Method) IsConstructor() bool {
	return m.GetName() == "constructor"
}

// ToABI returns a Geth ABI object built from a contract ABI
func (c *Contract) ToABI() (*ethabi.ABI, error) {
	a := &ethabi.ABI{}

	if c.ABI != "" {
		err := a.UnmarshalJSON([]byte(c.ABI))
		if err != nil {
			return nil, errors.EncodingError(err.Error())
		}
	}

	return a, nil
}

// Short returns a short string representation of contract information
func (m *Method) GetName() string {
	return strings.Split(m.Signature, "(")[0]
}
