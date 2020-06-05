package formatters

import (
	"net/http"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/service/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/entities"
)

func FormatJobResponse(job *entities.Job) *types.JobResponse {
	return &types.JobResponse{
		UUID:      job.UUID,
		ChainUUID: job.ChainUUID,
		Transaction: &entities.ETHTransaction{
			Hash:           job.Transaction.Hash,
			From:           job.Transaction.From,
			To:             job.Transaction.To,
			Nonce:          job.Transaction.Nonce,
			Value:          job.Transaction.Value,
			GasPrice:       job.Transaction.GasPrice,
			GasLimit:       job.Transaction.GasLimit,
			Data:           job.Transaction.Data,
			Raw:            job.Transaction.Raw,
			PrivateFrom:    job.Transaction.PrivateFrom,
			PrivateFor:     job.Transaction.PrivateFor,
			PrivacyGroupID: job.Transaction.PrivacyGroupID,
			CreatedAt:      job.Transaction.CreatedAt,
			UpdatedAt:      job.Transaction.UpdatedAt,
		},
		Labels:    job.Labels,
		Status:    job.Status,
		CreatedAt: job.CreatedAt,
	}
}

func FormatJobCreateRequest(request *types.CreateJobRequest) *entities.Job {
	return &entities.Job{
		Type:         request.Type,
		Labels:       request.Labels,
		ScheduleUUID: request.ScheduleUUID,
		Transaction:  request.Transaction,
	}
}

func FormatJobUpdateRequest(request *types.UpdateJobRequest) *entities.Job {
	return &entities.Job{
		Labels:      request.Labels,
		Transaction: request.Transaction,
		Status:      request.Status,
	}
}

func FormatJobFilterRequest(req *http.Request) (*entities.JobFilters, error) {
	filters := &entities.JobFilters{}

	qTxHashes := req.URL.Query().Get("tx_hashes")
	if qTxHashes != "" {
		for _, txHash := range strings.Split(qTxHashes, ",") {
			txHashT := strings.TrimSpace(txHash)
			if !utils.IsHash(txHashT) {
				err := errors.InvalidFormatError("invalid tx hash strings: %v", txHashT)
				return nil, err
			}

			filters.TxHashes = append(filters.TxHashes, common.HexToHash(txHashT))
		}
	}

	return filters, nil
}
