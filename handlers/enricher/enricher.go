package enricher

import (
	"github.com/ethereum/go-ethereum/crypto"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/ethereum/ethclient"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/proxy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/common"
	svc "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/contract-registry"
)

// Enricher is a Middleware engine.HandlerFunc
func Enricher(r svc.ContractRegistryClient, ec ethclient.ChainStateReader) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		if txctx.Builder.GetReceipt().GetContractAddress() != "" {
			url, err := proxy.GetURL(txctx)
			if err != nil {
				return
			}

			code, err := ec.CodeAt(txctx.Context(),
				url,
				txctx.Builder.GetReceipt().GetContractAddr(),
				nil)
			if err != nil {
				_ = txctx.AbortWithError(errors.InternalError(
					"could not read account code for chain %s and account %s",
					url,
					txctx.Builder.GetReceipt().GetContractAddr().Hex(),
				)).SetComponent(component)
				return
			}

			_, err = r.SetAccountCodeHash(txctx.Context(),
				&svc.SetAccountCodeHashRequest{
					AccountInstance: &common.AccountInstance{},
					CodeHash:        crypto.Keccak256Hash(code).String(),
				},
			)
			if err != nil {
				_ = txctx.AbortWithError(errors.InternalError("invalid input message format")).
					SetComponent(component)
				return
			}
			txctx.Logger.Debugf("%s successfully SetAccountCodeHash in Contract Registry for chain %s and account %s with codehash",
				txctx.Builder.GetChainIDString(),
				txctx.Builder.GetReceipt().GetContractAddress(),
				crypto.Keccak256Hash(code))
		}
	}
}
