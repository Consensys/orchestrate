package abi

import (
	"fmt"
	"regexp"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

// Short returns a string representation of the method
func (c *Contract) Short() string {
	if c.GetName() == "" {
		return ""
	}

	if c.GetTag() == "" {
		return c.GetName()
	}

	return fmt.Sprintf("%v[%v]", c.GetName(), c.GetTag())
}

var contractRegexp = `^(?P<contract>[a-zA-Z0-9]+)(\[(?P<tag>[0-9a-zA-Z-.]+)\])?$`
var contractPattern = regexp.MustCompile(contractRegexp)

// FromShortContract computes a Contract from is short representation
func FromShortContract(s string) (*Contract, error) {
	parts := contractPattern.FindStringSubmatch(s)

	if len(parts) != 4 {
		return nil, fmt.Errorf("%v is invalid short method (expected format %q)", s, contractRegexp)
	}

	name, tag := parts[1], parts[3]

	return &Contract{
		Name: name,
		Tag:  tag,
	}, nil
}

// ToABI returns a Geth ABI object built from a contract ABI
func (c *Contract) ToABI() (*abi.ABI, error) {
	abi := &abi.ABI{}
	err := abi.UnmarshalJSON(c.Abi)
	if err != nil {
		return nil, err
	}
	return abi, nil
}

// IsDeploy indicate wether the method refers to a deployment
func (m *Method) IsDeploy() bool {
	return m.GetName() == "constructor"
}
