package formatters

import (
	"net/http"
	"strings"

	pkgtypes "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/service/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/entities"
)

func FormatJobResponse(job *pkgtypes.Job) *types.JobResponse {
	return &types.JobResponse{
		UUID:      job.UUID,
		ChainUUID: job.ChainUUID,
		Transaction: &pkgtypes.ETHTransaction{
			Hash:           job.Transaction.Hash,
			From:           job.Transaction.From,
			To:             job.Transaction.To,
			Nonce:          job.Transaction.Nonce,
			Value:          job.Transaction.Value,
			GasPrice:       job.Transaction.GasPrice,
			Gas:            job.Transaction.Gas,
			Data:           job.Transaction.Data,
			Raw:            job.Transaction.Raw,
			PrivateFrom:    job.Transaction.PrivateFrom,
			PrivateFor:     job.Transaction.PrivateFor,
			PrivacyGroupID: job.Transaction.PrivacyGroupID,
			CreatedAt:      job.Transaction.CreatedAt,
			UpdatedAt:      job.Transaction.UpdatedAt,
		},
		Logs:        job.Logs,
		Labels:      job.Labels,
		Annotations: job.Annotations,
		Type:        job.Type,
		Status:      job.GetStatus(),
		CreatedAt:   job.CreatedAt,
	}
}

func FormatJobCreateRequest(request *types.CreateJobRequest) *pkgtypes.Job {
	return &pkgtypes.Job{
		Type:         request.Type,
		Labels:       request.Labels,
		Annotations:  request.Annotations,
		ScheduleUUID: request.ScheduleUUID,
		ChainUUID:    request.ChainUUID,
		Transaction:  request.Transaction,
	}
}

func FormatJobUpdateRequest(request *types.UpdateJobRequest) *pkgtypes.Job {
	return &pkgtypes.Job{
		Labels:      request.Labels,
		Annotations: request.Annotations,
		Transaction: request.Transaction,
	}
}

func FormatJobFilterRequest(req *http.Request) (*entities.JobFilters, error) {
	filters := &entities.JobFilters{}

	qTxHashes := req.URL.Query().Get("tx_hashes")
	if qTxHashes != "" {
		filters.TxHashes = strings.Split(qTxHashes, ",")
	}

	qChainUUID := req.URL.Query().Get("chain_uuid")
	if qChainUUID != "" {
		filters.ChainUUID = qChainUUID
	}

	if err := utils.GetValidator().Struct(filters); err != nil {
		return nil, err
	}

	return filters, nil
}
