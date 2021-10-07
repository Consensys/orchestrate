// nolint
package abi

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type Method struct {
	Name    string
	RawName string
	Const   bool
	Inputs  abi.Arguments
	Outputs abi.Arguments
}

func (method *Method) Sig() string {
	types := make([]string, len(method.Inputs))

	for i, input := range method.Inputs {
		types[i] = input.Type.String()
	}
	return fmt.Sprintf("%v(%v)", method.RawName, strings.Join(types, ","))
}

func (method *Method) ID() []byte {
	return crypto.Keccak256([]byte(method.Sig()))[:4]
}

type Event struct {
	Name      string
	RawName   string
	Anonymous bool
	Inputs    abi.Arguments
}

func (e *Event) Sig() string {
	types := make([]string, len(e.Inputs))
	for i, input := range e.Inputs {
		types[i] = input.Type.String()
	}
	return fmt.Sprintf("%v(%v)", e.RawName, strings.Join(types, ","))
}

func (e *Event) ID() common.Hash {
	return common.BytesToHash(crypto.Keccak256([]byte(e.Sig())))
}
