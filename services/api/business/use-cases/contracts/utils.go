package contracts

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
)

// returns the count of indexed inputs in the event
func getIndexedCount(event *abi.Event) (indexedInputCount uint) {
	for i := range event.Inputs {
		if event.Inputs[i].Indexed {
			indexedInputCount++
		}
	}

	return indexedInputCount
}
