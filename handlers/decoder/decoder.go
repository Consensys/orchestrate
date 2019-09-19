package decoder

import (
	"fmt"
	"strings"

	ethAbi "github.com/ethereum/go-ethereum/accounts/abi"
	log "github.com/sirupsen/logrus"

	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/encoding/json"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine"
	contractregistry "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/services/contract-registry"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/ethereum.git/abi/decoder"
)

// Decoder creates a decode handler
func Decoder(r contractregistry.RegistryClient) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		txctx.Logger = txctx.Logger.WithFields(log.Fields{
			"chain.id":    txctx.Envelope.GetChain().ID().String(),
			"tx.hash":     txctx.Envelope.GetReceipt().GetTxHash().Hex(),
			"metadata.id": txctx.Envelope.GetMetadata().GetId(),
		})

		// For each log in receipt
		for _, l := range txctx.Envelope.GetReceipt().GetLogs() {
			if len(l.GetTopics()) == 0 {
				// This scenario is not supposed to append
				err := fmt.Errorf("invalid receipt (no topics in log)")
				txctx.Logger.WithError(err).Errorf("decoder: invalid receipt")
				_ = txctx.AbortWithError(err)
				return
			}

			// Retrieve event ABI from registry
			eventResp, err := r.GetEventsBySigHash(
				txctx.Context(),
				&contractregistry.GetEventsBySigHashRequest{
					SigHash: l.Topics[0].GetRaw(),
					AccountInstance: &common.AccountInstance{
						Chain:   txctx.Envelope.GetChain(),
						Account: nil,
					},
					IndexedInputCount: uint32(len(l.Topics) - 1),
				})
			if err != nil || (len(eventResp.GetEvent()) == 0 && len(eventResp.GetDefaultEvents()) == 0) {
				txctx.Logger.WithError(err).Errorf("decoder: could not retrieve event ABI")
				_ = txctx.AbortWithError(err)
				return
			}

			// Decode log
			var mapping map[string]string
			var event *ethAbi.Event
			if len(eventResp.GetEvent()) != 0 {
				err = json.Unmarshal(eventResp.GetEvent(), event)
				if err != nil {
					txctx.Logger.WithError(err).Errorf("decoder: could not unmarshal event")
					_ = txctx.AbortWithError(err)
					return
				}
				mapping, err = decoder.Decode(event, l)
			} else {
				for _, potentialEvent := range eventResp.GetDefaultEvents() {
					// Try to unmarshal
					err = json.Unmarshal(potentialEvent, event)
					if err != nil {
						// If it fails to unmarshal, try the next potential event
						continue
					}

					// Try to decode
					mapping, err = decoder.Decode(event, l)
					if err == nil {
						break
					}
				}
			}

			if err != nil {
				txctx.Logger.WithError(err).Errorf("decoder: could not decode log")
				_ = txctx.AbortWithError(err)
				return
			}

			// Set decoded data on log
			l.DecodedData = mapping
			l.Event = GetAbi(event)

			txctx.Logger.WithFields(log.Fields{
				"log": mapping,
			}).Debug("decoder: log decoded")
		}
	}
}

// GetAbi creates a string ABI (format EventName(argType1, argType2)) from an event
func GetAbi(e *ethAbi.Event) string {
	inputs := make([]string, len(e.Inputs))
	for i := range e.Inputs {
		inputs[i] = fmt.Sprintf("%v", e.Inputs[i].Type)
	}
	return fmt.Sprintf("%v(%v)", e.Name, strings.Join(inputs, ","))
}
