package sender

import (
	"fmt"
	"strings"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/ethclient"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/chain-registry/proxy"
)

// If we need to send a transaction to Tessera enclave we first need to send a transaction data to Tessera
// to get a hash of this data. Then we need to replace data in a transaction object with a hash returned by
// Tessera enclave. We then need to sign the updated transaction
func TesseraPrivateTxSender(tesseraClient ethclient.QuorumTransactionSender) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		url, err := proxy.GetURL(txctx)
		if err != nil {
			return
		}

		txctx.Logger.WithField("envelope_id", txctx.Envelope.GetID()).
			WithField("job_uuid", txctx.Envelope.GetJobUUID()).
			Debugf("tessera handler starts")

		if txctx.Envelope.GetData() == "" {
			err := errors.DataError("cannot send transaction without data to Tessera").SetComponent(component)
			txctx.Logger.WithError(err).Errorf("failed to get transaction hash from Tessera")
			_ = txctx.AbortWithError(err)
			return
		}

		chainRegistryURL := strings.Replace(url, "/"+txctx.Envelope.GetChainUUID(), "", 1)
		proxyTessera := fmt.Sprintf("%s/tessera/%s", chainRegistryURL, txctx.Envelope.GetChainUUID())
		enclaveKey, err := tesseraClient.StoreRaw(
			txctx.Context(),
			proxyTessera,
			txctx.Envelope.MustGetDataBytes(),
			txctx.Envelope.GetPrivateFrom(),
		)

		if err != nil {
			e := txctx.AbortWithError(err).ExtendComponent(component)
			txctx.Logger.WithError(e).Errorf("failed to get transaction hash from Tessera")
			return
		}

		// Set TxHash to be the newly returned one instead of the computed one
		_ = txctx.Envelope.SetEnclaveKey(enclaveKey)

		txctx.Logger.Debugf("Sent transaction body to 'storeraw' endpoint to be signed: EnclaveKey %s", enclaveKey)
	}
}
