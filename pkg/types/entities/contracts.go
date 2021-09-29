package entities

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/consensys/orchestrate/pkg/errors"
	ethabi "github.com/consensys/orchestrate/pkg/go-ethereum/v1_9_12/accounts/abi"
)

const DefaultTagValue = "latest"

type Contract struct {
	Name             string
	Tag              string
	Registry         string
	ABI              string
	Bytecode         string
	DeployedBytecode string
	Constructor      Method
	Methods          []Method
	Events           []Event
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
