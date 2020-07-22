package scheduler

import (
	"context"
	"math/big"
	"reflect"

	utils2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/multitenancy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/auth/utils"
	clientutils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/client-utils"
	types2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/client"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/faucet/types"
)

// Faucet allows to credit by sending messages to a Kafka topic
type Faucet struct {
	txSchedulerClient client.TransactionSchedulerClient
}

// NewFaucet creates a New Faucet that can send message to a Kafka Topic
func NewFaucet(txSchedulerClient client.TransactionSchedulerClient) *Faucet {
	return &Faucet{
		txSchedulerClient: txSchedulerClient,
	}
}

// Credit process a Faucet credit request
func (f *Faucet) Credit(ctx context.Context, r *types.Request) (*big.Int, error) {
	// Elect final faucet
	if len(r.FaucetsCandidates) == 0 {
		return nil, errors.FaucetWarning("no faucet request").ExtendComponent(component)
	}

	// Select a first faucet candidate for comparison
	faucet := r.FaucetsCandidates[electFaucet(r.FaucetsCandidates)]

	if authToken := utils.AuthorizationFromContext(ctx); authToken != "" {
		ctx = context.WithValue(ctx, clientutils.RequestHeaderKey, map[string]string{
			multitenancy.AuthorizationMetadata: authToken,
		})
	}
	if apiKey := utils.APIKeyFromContext(ctx); apiKey != "" {
		ctx = context.WithValue(ctx, clientutils.RequestHeaderKey, map[string]string{
			utils.APIKeyHeader: apiKey,
		})
	}

	// If we have a chainName, we are funding an account generated in the signer
	if r.ChainName != "" {
		transferRequest := &types2.TransferRequest{
			BaseTransactionRequest: types2.BaseTransactionRequest{
				ChainName: r.ChainName,
				Labels: map[string]string{
					"id":            r.ChildTxID,
					"parentJobUUID": r.ParentTxID,
				}},
			Params: types2.TransferParams{
				Value: faucet.Amount.String(),
				From:  faucet.Creditor.String(),
				To:    r.Beneficiary.String(),
			},
		}
		txResponse, err := f.txSchedulerClient.SendTransferTransaction(ctx, transferRequest)
		if err != nil {
			errMessage := "failed to transfer funds from faucet account. Failed to transfer ETH"
			log.WithError(err).Error(errMessage)
			return nil, errors.ServiceConnectionError(errMessage)
		}

		log.WithField("tx_uuid", txResponse.UUID).Info("faucet: transfer transaction sent successfully for account generation")
	} else {
		// Create new job
		jobRequest := &types2.CreateJobRequest{
			ScheduleUUID: r.ScheduleUUID,
			ChainUUID:    r.ChainUUID,
			Type:         utils2.EthereumTransaction,
			Labels: map[string]string{
				"parentJobUUID": r.ParentTxID,
			},
			Transaction: &types2.ETHTransaction{
				From:  faucet.Creditor.String(),
				To:    r.Beneficiary.String(),
				Value: faucet.Amount.String(),
			},
		}

		jobResponse, err := f.txSchedulerClient.CreateJob(ctx, jobRequest)
		if err != nil {
			errMessage := "failed to transfer funds from faucet account. Failed to create job"
			log.WithError(err).Error(errMessage)
			return nil, errors.ServiceConnectionError(errMessage)
		}
		err = f.txSchedulerClient.StartJob(ctx, jobResponse.UUID)
		if err != nil {
			errMessage := "failed to transfer funds from faucet account. Failed to start job"
			log.WithError(err).Error(errMessage)
			return nil, errors.ServiceConnectionError(errMessage)
		}

		log.WithField("job_uuid", jobResponse.UUID).Info("faucet: transfer transaction sent successfully")
	}

	return faucet.Amount, nil
}

// electFaucet is currently selecting the remaining faucet candidates with the highest amount
func electFaucet(faucetsCandidates map[string]types.Faucet) string {
	// Select a first faucet candidate for comparison
	electedFaucet := reflect.ValueOf(faucetsCandidates).MapKeys()[0].String()
	for key, candidate := range faucetsCandidates {
		if candidate.Amount.Cmp(faucetsCandidates[electedFaucet].Amount) > 0 {
			electedFaucet = key
		}
	}
	return electedFaucet
}
