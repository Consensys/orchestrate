package redis

import (
	"reflect"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/common"

	ethabi "github.com/ethereum/go-ethereum/accounts/abi"
	ethcommon "github.com/ethereum/go-ethereum/common"
)

const methodsPrefix = "MethodsPrefix"

// MethodsModel is a zero-byte object gathering methods useful to interact with Events
type MethodsModel struct{}

// Methods returns an
var Methods = &MethodsModel{}

// Key serialize a lookup key for a set of methods stored on redis
func (*MethodsModel) Key(deployedByteCodeHash ethcommon.Hash, selector [4]byte) []byte {
	// Allocate memory to build the key
	res := make([]byte, 0, len(methodsPrefix)+len(deployedByteCodeHash)+4)
	// Fill in the buffer and return it
	res = append(res, methodsPrefix...)
	res = append(res, deployedByteCodeHash[:]...)
	return append(res, selector[:4]...)
}

// Get returns a serialized contract from its corresponding bytecode hash
func (m *MethodsModel) Get(conn *Conn, deployedByteCodeHash ethcommon.Hash, selector [4]byte) (methods [][]byte, ok bool, err error) {
	return conn.LRange(m.Key(deployedByteCodeHash, selector))
}

// Push a new methods to the registered methods list that have the same selector
func (m *MethodsModel) Push(conn *Conn, deployedByteCodeHash ethcommon.Hash, selector [4]byte, methodBytes []byte) error {
	return conn.RPush(m.Key(deployedByteCodeHash, selector), methodBytes)
}

// Find checks if a method is already registered for a given tuple (deployed bytecode hash, selector)
func (m *MethodsModel) Find(registeredMethods [][]byte, methodBytes []byte) bool {
	for _, registeredMethod := range registeredMethods {
		if reflect.DeepEqual(registeredMethod, methodBytes) {
			return true
		}
	}

	return false
}

// Registers all the methods of an abi in the contract-registry
func (m *MethodsModel) Registers(
	conn *Conn,
	deployedByteCodeHash, defaultByteCodeHash ethcommon.Hash,
	methods map[string]ethabi.Method,
	methodJSONs map[string][]byte) error {

	// Build a list of map keys, so that we always iterate on them with the same order
	// And precompute the methods selectors, and concatenate them in the same order.
	methodKeys := make([]string, 0, len(methods))
	selectors := make([][4]byte, 0, len(methods))

	for methodKey, method := range methods {
		selector := common.SigHashToSelector(method.ID())
		methodKeys = append(methodKeys, methodKey)
		selectors = append(selectors, selector)
	}

	// Push all methods to the new contract bytecodehash
	for index, methodKey := range methodKeys {
		err := conn.SendRPush(
			m.Key(deployedByteCodeHash, selectors[index]),
			methodJSONs[methods[methodKey].Sig()])

		if err != nil {
			return err
		}
	}

	// Warning: Only the first error will be returned
	err := conn.Flush()
	if err != nil {
		return err
	}

	// Check the outcome of the redis request
	for _, selector := range selectors {
		err = conn.ReceiveCheck()
		if err != nil {
			return err
		}

		err = conn.SendLRange(m.Key(defaultCodeHash, selector))
		if err != nil {
			return err
		}
	}

	// Warning: Only the first error will be returned
	err = conn.Flush()
	if err != nil {
		return err
	}

	// Accumulate the results if Find
	notFoundCount := 0

	// Fetch methods if they have already been registered
	for index, methodKey := range methodKeys {
		var registeredMethod [][]byte
		var ok bool
		registeredMethod, ok, err = conn.ReceiveByteSlices()
		if err != nil {
			return err
		}

		if !ok || !m.Find(registeredMethod, methodJSONs[methods[methodKey].Sig()]) {
			notFoundCount++

			err = conn.SendRPush(
				m.Key(defaultCodeHash, selectors[index]),
				methodJSONs[methods[methodKey].Sig()],
			)

			if err != nil {
				return err
			}
		}
	}

	// Warning: Only the first error will be returned
	err = conn.Flush()
	if err != nil {
		return err
	}

	// Check that all new registrations methods as default have been successful
	for index := 0; index < notFoundCount; index++ {
		err := conn.ReceiveCheck()
		if err != nil {
			return err
		}
	}

	return nil
}
