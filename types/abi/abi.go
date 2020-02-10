//nolint:stylecheck // reason
package abi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	ethabi "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
)

var component = "types.abi"

// Short returns a short string representation of contractId information
func (c *ContractId) Short() string {
	if c.GetName() == "" {
		return ""
	}

	if c.GetTag() == "" {
		return c.GetName()
	}

	return fmt.Sprintf("%v[%v]", c.GetName(), c.GetTag())
}

func (c *Contract) GetName() string {
	if c.GetId() == nil {
		return ""
	}
	return c.GetId().GetName()
}

func (c *Contract) SetName(name string) {
	if c.Id == nil {
		c.Id = &ContractId{}
	}
	c.Id.Name = name
}

func (c *Contract) GetTag() string {
	if c.GetId() == nil {
		return ""
	}
	return c.GetId().GetTag()
}

func (c *Contract) SetTag(tag string) {
	if c.Id == nil {
		c.Id = &ContractId{}
	}
	c.Id.Tag = tag
}

// Short returns a short string representation of contract information
func (c *Contract) Short() string {
	if c.Id == nil {
		return ""
	}
	return c.GetId().Short()
}

// Long return a long string representation of contract information
func (c *Contract) Long() string {
	return fmt.Sprintf("%v:%v:%v", c.Short(), c.Abi, c.Bytecode)
}

var contractRegexp = `^(?P<contract>[a-zA-Z0-9]+)(?:\[(?P<tag>[0-9a-zA-Z-.]+)\])?(?::(?P<abi>\[.+\]))?(?::(?P<bytecode>0[xX][a-fA-F0-9]+))?(?::(?P<deployedBytecode>0[xX][a-fA-F0-9]+))?$`
var contractPattern = regexp.MustCompile(contractRegexp)

// StringToContract computes a Contract from is short representation
func StringToContract(s string) (*Contract, error) {
	parts := contractPattern.FindStringSubmatch(s)

	if len(parts) != 6 {
		return nil, errors.InvalidFormatError("invalid contract (expected format %s) %q", contractRegexp, s).SetComponent(component)
	}

	c := &Contract{
		Id: &ContractId{
			Name: parts[1],
			Tag:  parts[2],
		},
	}

	// Make sure bytecode is valid and set bytecode
	if parts[4] == "" {
		parts[4] = "0x"
	}
	_, err := hexutil.Decode(parts[4])
	if err != nil {
		return nil, errors.InvalidFormatError("invalid contract bytecode on %q", c.Short()).SetComponent(component)
	}
	c.Bytecode = parts[4]

	// Make sure deployedBytecode is valid and set deployedBytecode
	if parts[5] == "" {
		parts[5] = "0x"
	}
	_, err = hexutil.Decode(parts[5])
	if err != nil {
		return nil, errors.InvalidFormatError("invalid contract deployed bytecode on %q", c.Short()).SetComponent(component)
	}
	c.DeployedBytecode = parts[5]

	// Set ABI and make sure it is valid
	c.Abi = parts[3]
	_, err = c.ToABI()
	if err != nil {
		return nil, errors.InvalidFormatError("invalid contract ABI on %q", c.Short()).SetComponent(component)
	}

	return c, nil
}

// ToABI returns a Geth ABI object built from a contract ABI
func (c *Contract) ToABI() (*ethabi.ABI, error) {
	a := &ethabi.ABI{}

	if len(c.Abi) > 0 {
		err := a.UnmarshalJSON([]byte(c.Abi))
		if err != nil {
			return nil, errors.EncodingError(err.Error()).SetComponent(component)
		}
	}

	return a, nil
}

// GetABICompacted returns a compacted version of the ABI
func (c *Contract) GetABICompacted() (string, error) {
	buffer := new(bytes.Buffer)
	if err := json.Compact(buffer, []byte(c.Abi)); err != nil {
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
	c.Abi = compactedABI
	return nil
}

// IsConstructor indicate whether the method refers to a deployment
func (m *Method) IsConstructor() bool {
	return m.GetName() == "constructor"
}

// Short returns a short string representation of contract information
func (m *Method) GetName() string {
	return strings.Split(m.GetSignature(), "(")[0]
}
