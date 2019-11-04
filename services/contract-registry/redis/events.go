package redis

import (
	"reflect"

	ethabi "github.com/ethereum/go-ethereum/accounts/abi"
	ethcommon "github.com/ethereum/go-ethereum/common"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/common"
)

const eventsPrefix = "EventsPrefix"

// EventsModel is a zero-byte object gathering events useful to interact with Events
type EventsModel struct{}

// Events returns an
var Events = &EventsModel{}

// Key serialize a lookup key for a set of events stored on redis
func (*EventsModel) Key(deployedByteCodeHash, eventID ethcommon.Hash, index uint) []byte {
	indexBytes := []byte{byte(index)}
	// Allocate memory to build the key
	res := make([]byte, 0,
		len(eventsPrefix)+len(deployedByteCodeHash)+len(eventID)+len(indexBytes))
	// Fill in the buffer and return it
	res = append(res, eventsPrefix...)
	res = append(res, deployedByteCodeHash[:]...)
	res = append(res, eventID[:]...)
	return append(res, indexBytes...)
}

// Get returns a serialized contract from its corresponding bytecode hash
func (e *EventsModel) Get(conn *Conn, deployedByteCodeHash, eventID ethcommon.Hash, index uint) (events [][]byte, ok bool, err error) {
	return conn.LRange(e.Key(deployedByteCodeHash, eventID, index))
}

// Push stores a new event in the contract-registry
func (e *EventsModel) Push(conn *Conn, deployedByteCodeHash, eventID ethcommon.Hash, index uint, eventBytes []byte) error {
	return conn.RPush(e.Key(deployedByteCodeHash, eventID, index), eventBytes)
}

// Find checks if an event is already registered for a given tuple (deployed bytecode hash, eventID, indexed count)
func (e *EventsModel) Find(registeredEvents [][]byte, event []byte) bool {
	for _, registeredMethod := range registeredEvents {
		if reflect.DeepEqual(registeredMethod, event) {
			return true
		}
	}

	return false
}

// Registers a batch of event from an Abi
func (e *EventsModel) Registers(conn *Conn,
	deployedByteCodeHash, defaultByteCodeHash ethcommon.Hash,
	events map[string]ethabi.Event,
	eventJSONs map[string][]byte,
) error {

	// Build a list of map keys, so that we always iterate on them with the same order
	// And precompute the events ids and indexedCounts, and concatenate them in the same order.
	eventKeys := make([]string, 0, len(events))
	eventIDs := make([]ethcommon.Hash, 0, len(events))
	indexedCounts := make([]uint, 0, len(events))

	for eventKey, event := range events {
		eventKeys = append(eventKeys, eventKey)
		eventIDs = append(eventIDs, event.ID())
		indexedCounts = append(indexedCounts, common.GetIndexedCount(event))
	}

	// Push all events to the new contract's bytecodehash
	for index, eventKey := range eventKeys {
		err := conn.SendRPush(
			e.Key(deployedByteCodeHash, eventIDs[index], indexedCounts[index]),
			eventJSONs[events[eventKey].Name])

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
	for index := range eventKeys {
		err = conn.ReceiveCheck()
		if err != nil {
			return err
		}

		err = conn.SendLRange(e.Key(defaultCodeHash, eventIDs[index], indexedCounts[index]))
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

	// Fetch events if they have already been registered
	for index, eventKey := range eventKeys {
		var registeredEvents [][]byte
		var ok bool
		registeredEvents, ok, err = conn.ReceiveByteSlices()
		if err != nil {
			return err
		}

		if !ok || !e.Find(registeredEvents, eventJSONs[events[eventKey].Name]) {
			notFoundCount++

			err = conn.SendRPush(
				e.Key(defaultCodeHash, eventIDs[index], indexedCounts[index]),
				eventJSONs[events[eventKey].Name],
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
