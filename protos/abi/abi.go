package abi

import (
	"fmt"
	"regexp"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

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

var contractRegexp = `^(?P<contract>[a-zA-Z0-9]+)(\[(?P<tag>[0-9a-zA-Z-.]+)\])?(:(?P<abi>\[.+\]))?(:(?P<bytecode>0[xX][a-fA-F0-9]+))?$`
var contractPattern = regexp.MustCompile(contractRegexp)

// StringToContract computes a Contract from is short representation
func StringToContract(s string) (*Contract, error) {
	parts := contractPattern.FindStringSubmatch(s)

	if len(parts) != 8 {
		return nil, fmt.Errorf("String format invalid (expected format %q): %q", contractRegexp, s)
	}

	c := &Contract{
		Name: parts[1],
		Tag:  parts[3],
	}

	// Make sure bytecode is valid and set bytecode
	if parts[7] == "" {
		parts[7] = "0x"
	}
	bytecode, err := hexutil.Decode(parts[7])
	if err != nil {
		return nil, fmt.Errorf("Contract %q bytecode is invalid", c.Short())
	}
	c.Bytecode = bytecode

	// Set ABI and make sure it is valid
	c.Abi = []byte(parts[5])
	_, err = c.ToABI()
	if err != nil {
		return nil, fmt.Errorf("Contract %q ABI is invalid", c.Short())
	}

	return c, nil
}

// ToABI returns a Geth ABI object built from a contract ABI
func (c *Contract) ToABI() (*abi.ABI, error) {
	abi := &abi.ABI{}

	if len(c.Abi) > 0 {
		err := abi.UnmarshalJSON(c.Abi)
		if err != nil {
			return nil, err
		}
	}

	return abi, nil
}

// IsDeploy indicate wether the method refers to a deployment
func (m *Method) IsDeploy() bool {
	return m.GetName() == "constructor"
}
