package handlers

import (
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/services"
	"github.com/ethereum/go-ethereum/common"
)

// HandlerSignature implements the interface TxSigner
func HandlerSignature(s *services.TxSigner) types.HandlerFunc {

	return func(ctx *types.Context) {

		if len(ctx.T.Tx().Raw()) > 0 {
			return
		}

	}

}