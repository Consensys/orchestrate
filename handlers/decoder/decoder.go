package decoder

import (
	"fmt"
	"strings"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"

	ethAbi "github.com/ethereum/go-ethereum/accounts/abi"
	log "github.com/sirupsen/logrus"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/ethereum/abi/decoder"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/encoding/json"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/common"
	svc "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/contract-registry"
)

// Decoder creates a decode handler
func Decoder(r svc.ContractRegistryClient) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		txctx.Logger = txctx.Logger.WithFields(log.Fields{
			"chainID": txctx.Builder.GetChainIDString(),
			"tx.hash": txctx.Builder.Receipt.GetTxHash(),
			"id":      txctx.Builder.GetID(),
		})

		// For each log in receipt
		for _, l := range txctx.Builder.Receipt.GetLogs() {
			if len(l.GetTopics()) == 0 {
				// This scenario is not supposed to happen
				err := errors.InternalError("invalid receipt (no topics in log)").ExtendComponent(component)
				txctx.Logger.WithError(err).Errorf("invalid receipt")
				_ = txctx.AbortWithError(err)
				return
			}

			// Retrieve event ABI from contract-registry
			eventResp, err := r.GetEventsBySigHash(
				txctx.Context(),
				&svc.GetEventsBySigHashRequest{
					SigHash: l.Topics[0],
					AccountInstance: &common.AccountInstance{
						ChainId: txctx.Builder.GetChainIDString(),
						Account: l.GetAddress(),
					},
					IndexedInputCount: uint32(len(l.Topics) - 1),
				},
			)
			if err != nil || (eventResp.GetEvent() == "" && len(eventResp.GetDefaultEvents()) == 0) {
				txctx.Logger.WithError(err).Tracef("%s: could not retrieve event ABI, txHash: %s sigHash: %s, ", component, l.GetTxHash(), l.GetTopics()[0])
				continue
			}

			// Decode log
			var mapping map[string]string
			event := &ethAbi.Event{}

			if eventResp.GetEvent() != "" {
				err = json.Unmarshal([]byte(eventResp.GetEvent()), event)
				if err != nil {
					txctx.Logger.WithError(err).Warnf("%s: could not unmarshal event ABI provided by the Contract Registry, txHash: %s sigHash: %s, ", component, l.GetTxHash(), l.GetTopics()[0])
					continue
				}
				mapping, err = decoder.Decode(event, l)
			} else {
				for _, potentialEvent := range eventResp.GetDefaultEvents() {
					// Try to unmarshal
					err = json.Unmarshal([]byte(potentialEvent), event)
					if err != nil {
						// If it fails to unmarshal, try the next potential event
						txctx.Logger.WithError(err).Tracef("%s: could not unmarshal potential event ABI, txHash: %s sigHash: %s, ", component, l.GetTxHash(), l.GetTopics()[0])
						continue
					}

					// Try to decode
					mapping, err = decoder.Decode(event, l)
					if err == nil {
						// As the decoding is successful, stop looping
						break
					}
				}
			}

			if err != nil {
				// As all potentialEvents fail to unmarshal, go to the next log
				txctx.Logger.WithError(err).Tracef("%s: could not unmarshal potential event ABI, txHash: %s sigHash: %s, ", component, l.GetTxHash(), l.GetTopics()[0])
				continue
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
