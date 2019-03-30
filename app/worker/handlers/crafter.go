package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/core/services"
	coreworker "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/core/worker"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/protos/ethereum"
	// "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/common/utils"
)

// Crafter creates a crafter handler
func Crafter(r services.ABIRegistry, c services.Crafter) coreworker.HandlerFunc {
	return func(ctx *coreworker.Context) {
		// Retrieve method identifier from trace
		if ctx.T.GetTx().GetTxData().GetData() != "" {
			// Transaction has already been crafted
			return
		}

		var method *abi.Method
		if ABI := ctx.T.GetCall().GetMethod().GetAbi(); len(ABI) > 0 {
			err := json.Unmarshal(ABI, method)
			if err != nil {
				ctx.AbortWithError(err)
				return
			}
		}

		if method == nil {
			methodID := ctx.T.GetCall().Short()
			if methodID == "" {
				// Nothing to craft
				return
			}

			m, err := r.GetMethodByID(methodID)
			if err != nil {
				ctx.Logger.WithError(err).Errorf("crafter: could not retrieve method ABI")
				ctx.AbortWithError(err)
				return
			}
			method = &m
			ctx.Logger = ctx.Logger.WithFields(log.Fields{
				"crafter.method": methodID,
			})
		}

		// Retrieve  args from trace
		args := ctx.T.GetCall().GetArgs()
		ctx.Logger = ctx.Logger.WithFields(log.Fields{
			"crafter.args": args,
		})
		log.Debugf("method: %v", method.Sig())
		log.Debugf("method: %v", method.String())
		log.Debugf("method: %v", hexutil.Encode(method.Id()))

		var (payload []byte; err error)

		if ctx.T.GetCall().GetMethod().GetName() == "constructor" {
			// Craft transaction payload
			payload, err = c.CraftConstructor(*method, args...)
			if err != nil {
				ctx.Logger.WithError(err).Errorf("crafter: could not craft constructor tx data payload")
				ctx.AbortWithError(err)
				return
			}

			// This is a deployment call
			contractName := ctx.T.GetCall().GetContract().Short()
			bytecode, err := r.GetBytecodeByID(
				fmt.Sprintf("constructor@%v", contractName),
			)
			if err != nil {
				ctx.Logger.WithError(err).Errorf("crafter: could not craft tx data payload")
			}
			if len(bytecode) == 0 {
				ctx.Logger.WithError(fmt.Errorf("Invalid empty bytecode")).Errorf("crafter: could not craft tx data payload")
			}
			payload = append(bytecode, payload...)
		} else {
			// Craft transaction payload
			payload, err = c.Craft(*method, args...)
			if err != nil {
				ctx.Logger.WithError(err).Errorf("crafter: could not craft tx data payload")
				ctx.AbortWithError(err)
				return
			}
		}

		ctx.Logger = ctx.Logger.WithFields(log.Fields{
			"tx.data": hexutil.Encode(payload),
		})

		// Update Trace
		if ctx.T.GetTx() == nil {
			ctx.T.Tx = &ethereum.Transaction{TxData: &ethereum.TxData{}}
		} else if ctx.T.GetTx().GetTxData() == nil {
			ctx.T.Tx.TxData = &ethereum.TxData{}
		}
		ctx.T.GetTx().GetTxData().SetData(payload)

		ctx.Logger.Debugf("crafter: tx data payload set")
	}
}
