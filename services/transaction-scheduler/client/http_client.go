package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v4"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/httputil"
	types "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/tx-scheduler"

	"github.com/containous/traefik/v2/pkg/log"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	clientutils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/client-utils"
)

const (
	invalidResponseBody = "failed to decode response body"
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

func (c *HTTPClient) SearchJob(ctx context.Context, txHashes []string, chainUUID, status string) ([]*types.JobResponse, error) {
	reqURL := fmt.Sprintf("%v/jobs", c.config.URL)
	var resp []*types.JobResponse

	var qParams []string
	if len(txHashes) > 0 {
		qParams = append(qParams, "tx_hashes="+strings.Join(txHashes, ","))
	}

	if chainUUID != "" {
		qParams = append(qParams, "chain_uuid="+chainUUID)
	}

	if status != "" {
		qParams = append(qParams, "status="+status)
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
	resp := &types.JobResponse{}

	return callWithBackOff(ctx, c.config.backOff, func() error {
		response, err := clientutils.PutRequest(ctx, c.client, reqURL, nil)
		if err != nil {
			errMessage := "error while starting job"
			log.FromContext(ctx).WithError(err).Error(errMessage)
			return errors.ServiceConnectionError(errMessage).ExtendComponent(component)
		}

		defer clientutils.CloseResponse(response)
		return parseResponse(ctx, response, resp)
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
	if response.StatusCode == http.StatusAccepted || response.StatusCode == http.StatusOK {
		if resp == nil {
			return nil
		}

		if err := json.NewDecoder(response.Body).Decode(resp); err != nil {
			log.FromContext(ctx).WithError(err).Error(invalidResponseBody)
			return errors.ServiceConnectionError(invalidResponseBody).ExtendComponent(component)
		}

		return nil
	}

	errResp := httputil.ErrorResponse{}
	if err := json.NewDecoder(response.Body).Decode(&errResp); err != nil {
		log.FromContext(ctx).WithError(err).Error(invalidResponseBody)
		return errors.ServiceConnectionError(invalidResponseBody).ExtendComponent(component)
	}

	log.FromContext(ctx).Error(errResp.Message)
	return errors.Errorf(errResp.Code, errResp.Message)
}
