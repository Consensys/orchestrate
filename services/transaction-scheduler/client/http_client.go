package client

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	healthz "github.com/heptiolabs/healthcheck"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/entities"

	"github.com/cenkalti/backoff/v4"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/httputil"
	types "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/txscheduler"

	"github.com/containous/traefik/v2/pkg/log"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	clientutils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/client-utils"
)

type HTTPClient struct {
	client *http.Client
	config *Config
}

func NewHTTPClient(h *http.Client, c *Config) TransactionSchedulerClient {
	return &HTTPClient{
		client: h,
		config: c,
	}
}

func (c HTTPClient) Checker() healthz.Check {
	return healthz.HTTPGetCheck(fmt.Sprintf("%s/live", c.config.MetricsURL), time.Second)
}

func (c *HTTPClient) SendContractTransaction(ctx context.Context, txRequest *types.SendTransactionRequest) (*types.TransactionResponse, error) {
	reqURL := fmt.Sprintf("%v/transactions/send", c.config.URL)
	resp := &types.TransactionResponse{}

	err := callWithBackOff(ctx, c.config.backOff, func() error {
		response, err := clientutils.PostRequest(ctx, c.client, reqURL, txRequest)
		if err != nil {
			errMessage := "error while sending transaction"
			log.FromContext(ctx).WithError(err).Error(errMessage)
			return errors.ServiceConnectionError(errMessage).ExtendComponent(component)
		}

		defer clientutils.CloseResponse(response)
		return parseResponse(ctx, response, resp)
	})

	return resp, err
}

func (c *HTTPClient) SendDeployTransaction(ctx context.Context, txRequest *types.DeployContractRequest) (*types.TransactionResponse, error) {
	reqURL := fmt.Sprintf("%v/transactions/deploy-contract", c.config.URL)
	resp := &types.TransactionResponse{}

	err := callWithBackOff(ctx, c.config.backOff, func() error {
		response, err := clientutils.PostRequest(ctx, c.client, reqURL, txRequest)
		if err != nil {
			errMessage := "error while sending deploy contract transaction"
			log.FromContext(ctx).WithError(err).Error(errMessage)
			return errors.ServiceConnectionError(errMessage).ExtendComponent(component)
		}

		defer clientutils.CloseResponse(response)
		return parseResponse(ctx, response, resp)
	})

	return resp, err
}

func (c *HTTPClient) SendRawTransaction(ctx context.Context, txRequest *types.RawTransactionRequest) (*types.TransactionResponse, error) {
	reqURL := fmt.Sprintf("%v/transactions/send-raw", c.config.URL)
	resp := &types.TransactionResponse{}

	err := callWithBackOff(ctx, c.config.backOff, func() error {
		response, err := clientutils.PostRequest(ctx, c.client, reqURL, txRequest)
		if err != nil {
			errMessage := "error while sending raw transaction"
			log.FromContext(ctx).WithError(err).Error(errMessage)
			return errors.ServiceConnectionError(errMessage).ExtendComponent(component)
		}

		defer clientutils.CloseResponse(response)
		return parseResponse(ctx, response, resp)
	})

	return resp, err
}

func (c *HTTPClient) SendTransferTransaction(ctx context.Context, txRequest *types.TransferRequest) (*types.TransactionResponse, error) {
	reqURL := fmt.Sprintf("%v/transactions/transfer", c.config.URL)
	resp := &types.TransactionResponse{}

	err := callWithBackOff(ctx, c.config.backOff, func() error {
		response, err := clientutils.PostRequest(ctx, c.client, reqURL, txRequest)
		if err != nil {
			errMessage := "error while sending transfer transaction"
			log.FromContext(ctx).WithError(err).Error(errMessage)
			return errors.ServiceConnectionError(errMessage).ExtendComponent(component)
		}

		defer clientutils.CloseResponse(response)
		return parseResponse(ctx, response, resp)
	})

	return resp, err
}

func (c *HTTPClient) GetTxRequest(ctx context.Context, txRequestUUID string) (*types.TransactionResponse, error) {
	reqURL := fmt.Sprintf("%v/transactions/%v", c.config.URL, txRequestUUID)
	resp := &types.TransactionResponse{}

	err := callWithBackOff(ctx, c.config.backOff, func() error {
		response, err := clientutils.GetRequest(ctx, c.client, reqURL)
		if err != nil {
			errMessage := "error while getting transaction request"
			log.FromContext(ctx).WithError(err).Error(errMessage)
			return errors.ServiceConnectionError(errMessage).ExtendComponent(component)
		}

		defer clientutils.CloseResponse(response)
		return parseResponse(ctx, response, resp)
	})

	return resp, err
}

func (c *HTTPClient) GetSchedule(ctx context.Context, scheduleUUID string) (*types.ScheduleResponse, error) {
	reqURL := fmt.Sprintf("%v/schedules/%v", c.config.URL, scheduleUUID)
	resp := &types.ScheduleResponse{}

	err := callWithBackOff(ctx, c.config.backOff, func() error {
		response, err := clientutils.GetRequest(ctx, c.client, reqURL)
		if err != nil {
			errMessage := "error while getting schedule"
			log.FromContext(ctx).WithError(err).Error(errMessage)
			return errors.ServiceConnectionError(errMessage).ExtendComponent(component)
		}

		defer clientutils.CloseResponse(response)
		return parseResponse(ctx, response, resp)
	})

	return resp, err
}

func (c *HTTPClient) GetSchedules(ctx context.Context) ([]*types.ScheduleResponse, error) {
	reqURL := fmt.Sprintf("%v/schedules", c.config.URL)
	var resp []*types.ScheduleResponse

	err := callWithBackOff(ctx, c.config.backOff, func() error {
		response, err := clientutils.GetRequest(ctx, c.client, reqURL)
		if err != nil {
			errMessage := "error while getting schedules"
			log.FromContext(ctx).WithError(err).Error(errMessage)
			return errors.ServiceConnectionError(errMessage).ExtendComponent(component)
		}

		defer clientutils.CloseResponse(response)
		return parseResponse(ctx, response, &resp)
	})

	return resp, err
}

func (c *HTTPClient) CreateSchedule(ctx context.Context, request *types.CreateScheduleRequest) (*types.ScheduleResponse, error) {
	reqURL := fmt.Sprintf("%v/schedules", c.config.URL)
	resp := &types.ScheduleResponse{}

	err := callWithBackOff(ctx, c.config.backOff, func() error {
		response, err := clientutils.PostRequest(ctx, c.client, reqURL, request)
		if err != nil {
			errMessage := "error while creating schedule"
			log.FromContext(ctx).WithError(err).Error(errMessage)
			return errors.ServiceConnectionError(errMessage).ExtendComponent(component)
		}

		defer clientutils.CloseResponse(response)
		return parseResponse(ctx, response, resp)
	})

	return resp, err
}

func (c *HTTPClient) GetJob(ctx context.Context, jobUUID string) (*types.JobResponse, error) {
	reqURL := fmt.Sprintf("%v/jobs/%s", c.config.URL, jobUUID)
	resp := &types.JobResponse{}

	err := callWithBackOff(ctx, c.config.backOff, func() error {
		response, err := clientutils.GetRequest(ctx, c.client, reqURL)
		if err != nil {
			errMessage := "error while getting job"
			log.FromContext(ctx).WithError(err).Error(errMessage)
			return errors.ServiceConnectionError(errMessage).ExtendComponent(component)
		}

		defer clientutils.CloseResponse(response)
		return parseResponse(ctx, response, resp)
	})

	return resp, err
}

func (c *HTTPClient) GetJobs(ctx context.Context) ([]*types.JobResponse, error) {
	reqURL := fmt.Sprintf("%v/jobs", c.config.URL)
	var resp []*types.JobResponse

	err := callWithBackOff(ctx, c.config.backOff, func() error {
		response, err := clientutils.GetRequest(ctx, c.client, reqURL)
		if err != nil {
			errMessage := "error while getting jobs"
			log.FromContext(ctx).WithError(err).Error(errMessage)
			return errors.ServiceConnectionError(errMessage).ExtendComponent(component)
		}

		defer clientutils.CloseResponse(response)
		return parseResponse(ctx, response, &resp)
	})

	return resp, err
}

func (c *HTTPClient) SearchJob(ctx context.Context, filters *entities.JobFilters) ([]*types.JobResponse, error) {
	reqURL := fmt.Sprintf("%v/jobs", c.config.URL)
	var resp []*types.JobResponse

	var qParams []string
	if len(filters.TxHashes) > 0 {
		qParams = append(qParams, "tx_hashes="+strings.Join(filters.TxHashes, ","))
	}

	if filters.ChainUUID != "" {
		qParams = append(qParams, "chain_uuid="+filters.ChainUUID)
	}

	if filters.Status != "" {
		qParams = append(qParams, "status="+filters.Status)
	}

	if !filters.UpdatedAfter.IsZero() {
		qParams = append(qParams, "updated_after="+filters.UpdatedAfter.Format(time.RFC3339))
	}

	if filters.OnlyParents {
		qParams = append(qParams, "only_parents=true")
	}

	if filters.ParentJobUUID != "" {
		qParams = append(qParams, "parent_job_uuid="+filters.ParentJobUUID)
	}

	if len(qParams) > 0 {
		reqURL = reqURL + "?" + strings.Join(qParams, "&")
	}

	err := callWithBackOff(ctx, c.config.backOff, func() error {
		response, err := clientutils.GetRequest(ctx, c.client, reqURL)
		if err != nil {
			errMessage := "error while searching jobs"
			log.FromContext(ctx).WithError(err).Error(errMessage)
			return errors.ServiceConnectionError(errMessage).ExtendComponent(component)
		}
		defer clientutils.CloseResponse(response)
		return parseResponse(ctx, response, &resp)
	})

	return resp, err
}

func (c *HTTPClient) CreateJob(ctx context.Context, request *types.CreateJobRequest) (*types.JobResponse, error) {
	reqURL := fmt.Sprintf("%v/jobs", c.config.URL)
	resp := &types.JobResponse{}

	err := callWithBackOff(ctx, c.config.backOff, func() error {
		response, err := clientutils.PostRequest(ctx, c.client, reqURL, request)
		if err != nil {
			errMessage := "error while creating job"
			log.FromContext(ctx).WithError(err).Error(errMessage)
			return errors.ServiceConnectionError(errMessage).ExtendComponent(component)
		}
		defer clientutils.CloseResponse(response)
		return parseResponse(ctx, response, resp)
	})

	return resp, err
}

func (c *HTTPClient) UpdateJob(ctx context.Context, jobUUID string, request *types.UpdateJobRequest) (*types.JobResponse, error) {
	reqURL := fmt.Sprintf("%v/jobs/%s", c.config.URL, jobUUID)
	resp := &types.JobResponse{}

	err := callWithBackOff(ctx, c.config.backOff, func() error {
		response, err := clientutils.PatchRequest(ctx, c.client, reqURL, request)
		if err != nil {
			errMessage := "error while updating job"
			log.FromContext(ctx).WithError(err).Error(errMessage)
			return errors.ServiceConnectionError(errMessage).ExtendComponent(component)
		}

		defer clientutils.CloseResponse(response)
		return parseResponse(ctx, response, resp)
	})

	return resp, err
}

func (c *HTTPClient) StartJob(ctx context.Context, jobUUID string) error {
	reqURL := fmt.Sprintf("%v/jobs/%s/start", c.config.URL, jobUUID)

	return callWithBackOff(ctx, c.config.backOff, func() error {
		response, err := clientutils.PutRequest(ctx, c.client, reqURL, nil)
		if err != nil {
			errMessage := "error while starting job"
			log.FromContext(ctx).WithError(err).Error(errMessage)
			return errors.ServiceConnectionError(errMessage).ExtendComponent(component)
		}

		defer clientutils.CloseResponse(response)
		return nil
	})
}

func (c *HTTPClient) ResendJobTx(ctx context.Context, jobUUID string) error {
	reqURL := fmt.Sprintf("%v/jobs/%s/resend", c.config.URL, jobUUID)

	return callWithBackOff(ctx, c.config.backOff, func() error {
		response, err := clientutils.PutRequest(ctx, c.client, reqURL, nil)
		if err != nil {
			errMessage := "error while resending job tx"
			log.FromContext(ctx).WithError(err).Error(errMessage)
			return errors.ServiceConnectionError(errMessage).ExtendComponent(component)
		}

		defer clientutils.CloseResponse(response)
		return nil
	})
}

func callWithBackOff(ctx context.Context, backOff backoff.BackOff, requestCall func() error) error {
	return backoff.RetryNotify(
		func() error {
			err := requestCall()
			// If not errors, it does not retry
			if err == nil {
				return nil
			}

			// Retry on following errors
			if errors.IsInvalidStateError(err) || errors.IsServiceConnectionError(err) {
				return err
			}

			// Otherwise, stop retrying
			return backoff.Permanent(err)
		}, backoff.WithContext(backOff, ctx),
		func(e error, duration time.Duration) {
			log.FromContext(ctx).
				WithError(e).
				Warnf("transaction-scheduler: http call failed, retrying in %v...", duration)
		},
	)
}

func parseResponse(ctx context.Context, response *http.Response, resp interface{}) error {
	err := httputil.ParseResponse(ctx, response, resp)
	if err != nil {
		return errors.FromError(err).ExtendComponent(component)
	}

	return nil
}
