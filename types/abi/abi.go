package abi

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"
	errors "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/errors"
)

var component = "types.abi"

func (c *Contract) GetName() string {
	return c.GetId().GetName()
}

func (c *Contract) SetName(name string) {
	c.Id.Name = name
}

func (c *Contract) GetTag() string {
	return c.GetId().GetTag()
}

func (c *Contract) SetTag(tag string) {
	c.Id.Tag = tag
}

// Short returns a short string representation of contract information
func (c *Contract) Short() string {
	if c.GetName() == "" {
		return ""
	}

	if c.GetTag() == "" {
		return c.GetName()
	}

	return fmt.Sprintf("%v[%v]", c.GetName(), c.GetTag())
}

// Long return a long string representation of contract information
func (c *Contract) Long() string {
	return fmt.Sprintf("%v:%v:%v", c.Short(), string(c.Abi), string(c.Bytecode))
}

var contractRegexp = `^(?P<contract>[a-zA-Z0-9]+)(?:\[(?P<tag>[0-9a-zA-Z-.]+)\])?(?::(?P<abi>\[.+\]))?(?::(?P<bytecode>0[xX][a-fA-F0-9]+))?(?::(?P<deployedBytecode>0[xX][a-fA-F0-9]+))?$`
var contractPattern = regexp.MustCompile(contractRegexp)

// StringToContract computes a Contract from is short representation
func StringToContract(s string) (*Contract, error) {
	parts := contractPattern.FindStringSubmatch(s)

	if len(parts) != 6 {
		return nil, errors.InvalidFormatError("invalid contract (expected format %q) %q", contractRegexp, s).SetComponent(component)
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
	bytecode, err := hexutil.Decode(parts[4])
	if err != nil {
		return nil, errors.InvalidFormatError("invalid contract bytecode on %q", c.Short()).SetComponent(component)
	}
	c.Bytecode = bytecode

	// Make sure deployedBytecode is valid and set deployedBytecode
	if parts[5] == "" {
		parts[5] = "0x"
	}
	deployedBytecode, err := hexutil.Decode(parts[5])
	if err != nil {
		return nil, errors.InvalidFormatError("invalid contract deployed bytecode on %q", c.Short()).SetComponent(component)
	}
	c.DeployedBytecode = deployedBytecode

	// Set ABI and make sure it is valid
	c.Abi = []byte(parts[3])
	_, err = c.ToABI()
	if err != nil {
		return nil, errors.InvalidFormatError("invalid contract ABI on %q", c.Short()).SetComponent(component)
	}

	return c, nil
}

// ToABI returns a Geth ABI object built from a contract ABI
func (c *Contract) ToABI() (*abi.ABI, error) {
	a := &abi.ABI{}

	if len(c.Abi) > 0 {
		err := a.UnmarshalJSON(c.Abi)
		if err != nil {
			return nil, errors.EncodingError(err.Error()).SetComponent(component)
		}
	}

	return a, nil
}

// IsConstructor indicate whether the method refers to a deployment
func (m *Method) IsConstructor() bool {
	return m.GetName() == "constructor"
}

// Short returns a short string representation of contract information
func (m *Method) GetName() string {
	return strings.Split(m.GetSignature(), "(")[0]
}
