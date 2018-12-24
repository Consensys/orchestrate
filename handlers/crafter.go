package handlers

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"gitlab.com/ConsenSys/client/fr/core-stack/core/ethereum"
	"gitlab.com/ConsenSys/client/fr/core-stack/core/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/core/infra"
)

// ABIGetter is an interface to retrieve ABIs
type ABIGetter interface {
	// Must return
	GetMethodByID(ID string) (*abi.Method, error)
}

// DummyABIGetter always return the same ABI method (useful for testing purpose)
type DummyABIGetter struct {
	abi *abi.Method
}

// NewDummyABIGetter creates a new DummyABIgetter
func NewDummyABIGetter(abi *abi.Method) *DummyABIGetter {
	return &DummyABIGetter{abi}
}

// GetMethodByID return ABI
func (g *DummyABIGetter) GetMethodByID(ID string) (*abi.Method, error) {
	return g.abi, nil
}

// Crafter creates a crafter handler
func Crafter(g ABIGetter) infra.HandlerFunc {
	return func(ctx *infra.Context) {
		// Retrieve method identifier and args from trace
		methodID, args := ctx.T.Call().MethodID, ctx.T.Call().Args

		// Retrieve method ABI object
		method, err := g.GetMethodByID(methodID)
		if err != nil {
			e := types.Error{
				Err:  err,
				Type: 0, // TODO: add an error type ErrorTypeABIGet
			}
			// Abort execution
			ctx.AbortWithError(e)
			return
		}

		// Craft transaction payload
		payload, err := ethereum.CraftPayload(method, args)
		if err != nil {
			e := types.Error{
				Err:  err,
				Type: 0, // TODO: add an error type ErrorTypeCraft
			}
			// Abort execution
			ctx.AbortWithError(e)
			return
		}

		// Update Trace
		ctx.T.Tx().SetData(payload)
	}
}
