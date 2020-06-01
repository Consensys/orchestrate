package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/containous/traefik/v2/pkg/log"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	clientutils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/client-utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/service/types"
)

const (
	invalidRequestErrorMessage = "invalid request payload"
	invalidStatus              = "unhandled invalid response status"
	invalidRequestBody         = "failed to decode request body"
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

func (c *HTTPClient) SendTransaction(ctx context.Context, chainUUID string, txRequest *types.SendTransactionRequest) (*types.TransactionResponse, error) {
	reqURL := fmt.Sprintf("%v/transactions/%s/send", c.config.URL, chainUUID)

	response, err := clientutils.PostRequest(ctx, c.client, reqURL, txRequest)
	if err != nil {
		errMessage := "error while sending transaction"
		log.FromContext(ctx).WithError(err).Error(errMessage)
		return nil, errors.ServiceConnectionError(errMessage).ExtendComponent(component)
	}
	defer clientutils.CloseResponse(response)

	resp := &types.TransactionResponse{}
	err = parseResponse(ctx, response, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *HTTPClient) GetSchedule(ctx context.Context, scheduleUUID string) (*types.ScheduleResponse, error) {
	reqURL := fmt.Sprintf("%v/schedules/%v", c.config.URL, scheduleUUID)

	response, err := clientutils.GetRequest(ctx, c.client, reqURL)
	if err != nil {
		errMessage := "error while getting schedule"
		log.FromContext(ctx).WithError(err).Error(errMessage)
		return nil, errors.ServiceConnectionError(errMessage).ExtendComponent(component)
	}
	defer clientutils.CloseResponse(response)

	resp := &types.ScheduleResponse{}
	err = parseResponse(ctx, response, resp)
	return resp, err
}

func (c *HTTPClient) GetSchedules(ctx context.Context) ([]*types.ScheduleResponse, error) {
	reqURL := fmt.Sprintf("%v/schedules", c.config.URL)

	response, err := clientutils.GetRequest(ctx, c.client, reqURL)
	if err != nil {
		errMessage := "error while getting schedules"
		log.FromContext(ctx).WithError(err).Error(errMessage)
		return nil, errors.ServiceConnectionError(errMessage).ExtendComponent(component)
	}
	defer clientutils.CloseResponse(response)

	resp := []*types.ScheduleResponse{}
	err = parseResponse(ctx, response, &resp)
	return resp, err
}

func (c *HTTPClient) CreateSchedule(ctx context.Context, request *types.CreateScheduleRequest) (*types.ScheduleResponse, error) {
	reqURL := fmt.Sprintf("%v/schedules", c.config.URL)

	response, err := clientutils.PostRequest(ctx, c.client, reqURL, request)
	if err != nil {
		errMessage := "error while creating schedule"
		log.FromContext(ctx).WithError(err).Error(errMessage)
		return nil, errors.ServiceConnectionError(errMessage).ExtendComponent(component)
	}
	defer clientutils.CloseResponse(response)

	resp := &types.ScheduleResponse{}
	err = parseResponse(ctx, response, resp)
	return resp, err
}

func (c *HTTPClient) GetJob(ctx context.Context, jobUUID string) (*types.JobResponse, error) {
	reqURL := fmt.Sprintf("%v/jobs/%s", c.config.URL, jobUUID)

	response, err := clientutils.GetRequest(ctx, c.client, reqURL)
	if err != nil {
		errMessage := "error while getting job"
		log.FromContext(ctx).WithError(err).Error(errMessage)
		return nil, errors.ServiceConnectionError(errMessage).ExtendComponent(component)
	}
	defer clientutils.CloseResponse(response)

	resp := &types.JobResponse{}
	err = parseResponse(ctx, response, resp)
	return resp, err
}

func (c *HTTPClient) GetJobs(ctx context.Context) ([]*types.JobResponse, error) {
	reqURL := fmt.Sprintf("%v/jobs", c.config.URL)

	response, err := clientutils.GetRequest(ctx, c.client, reqURL)
	if err != nil {
		errMessage := "error while getting jobs"
		log.FromContext(ctx).WithError(err).Error(errMessage)
		return nil, errors.ServiceConnectionError(errMessage).ExtendComponent(component)
	}
	defer clientutils.CloseResponse(response)

	resp := []*types.JobResponse{}
	err = parseResponse(ctx, response, &resp)
	return resp, err
}

func (c *HTTPClient) SearchJob(ctx context.Context, txHashes []string) ([]*types.JobResponse, error) {
	reqURL := fmt.Sprintf("%v/jobs", c.config.URL)

	if len(txHashes) > 0 {
		reqURL = fmt.Sprintf("%s?tx_hashes=%s", reqURL, strings.Join(txHashes, ","))
	}

	response, err := clientutils.GetRequest(ctx, c.client, reqURL)
	if err != nil {
		errMessage := "error while searching jobs"
		log.FromContext(ctx).WithError(err).Error(errMessage)
		return nil, errors.ServiceConnectionError(errMessage).ExtendComponent(component)
	}
	defer clientutils.CloseResponse(response)

	resp := []*types.JobResponse{}
	err = parseResponse(ctx, response, &resp)
	return resp, err
}

func (c *HTTPClient) CreateJob(ctx context.Context, request *types.CreateJobRequest) (*types.JobResponse, error) {
	reqURL := fmt.Sprintf("%v/jobs", c.config.URL)

	response, err := clientutils.PostRequest(ctx, c.client, reqURL, request)
	if err != nil {
		errMessage := "error while creating job"
		log.FromContext(ctx).WithError(err).Error(errMessage)
		return nil, errors.ServiceConnectionError(errMessage).ExtendComponent(component)
	}
	defer clientutils.CloseResponse(response)

	resp := &types.JobResponse{}
	err = parseResponse(ctx, response, resp)
	return resp, err
}

func (c *HTTPClient) UpdateJob(ctx context.Context, jobUUID string, request *types.UpdateJobRequest) (*types.JobResponse, error) {
	reqURL := fmt.Sprintf("%v/jobs/%s", c.config.URL, jobUUID)

	response, err := clientutils.PatchRequest(ctx, c.client, reqURL, request)
	if err != nil {
		errMessage := "error while updating job"
		log.FromContext(ctx).WithError(err).Error(errMessage)
		return nil, errors.ServiceConnectionError(errMessage).ExtendComponent(component)
	}
	defer clientutils.CloseResponse(response)

	resp := &types.JobResponse{}
	err = parseResponse(ctx, response, resp)
	return resp, err
}

func (c *HTTPClient) StartJob(ctx context.Context, jobUUID string) error {
	reqURL := fmt.Sprintf("%v/jobs/%s", c.config.URL, jobUUID)

	response, err := clientutils.PutRequest(ctx, c.client, reqURL, nil)
	if err != nil {
		errMessage := "error while starting job"
		log.FromContext(ctx).WithError(err).Error(errMessage)
		return errors.ServiceConnectionError(errMessage).ExtendComponent(component)
	}
	defer clientutils.CloseResponse(response)

	resp := &types.JobResponse{}
	return parseResponse(ctx, response, resp)
}

func parseResponse(ctx context.Context, response *http.Response, resp interface{}) error {
	switch response.StatusCode {
	case http.StatusAccepted, http.StatusOK:
		if resp == nil {
			return nil
		}

		if err := json.NewDecoder(response.Body).Decode(resp); err != nil {
			log.FromContext(ctx).WithError(err).Error(invalidRequestBody)
			return errors.ServiceConnectionError(invalidRequestBody).ExtendComponent(component)
		}

		return nil
	case http.StatusBadRequest:
		log.FromContext(ctx).Error(invalidRequestErrorMessage)
		return errors.InvalidFormatError(invalidRequestErrorMessage)
	case http.StatusConflict:
		errMessage := "entity already exists"
		log.FromContext(ctx).Error(errMessage)
		return errors.ConflictedError(errMessage)
	case http.StatusUnprocessableEntity:
		errMessage := "unprocessable entity"
		log.FromContext(ctx).Error(errMessage)
		return errors.InvalidParameterError(errMessage)
	case http.StatusNotFound:
		errMessage := "entity not found"
		log.FromContext(ctx).Error(errMessage)
		return errors.NotFoundError(errMessage)
	default:
		log.FromContext(ctx).Error(invalidStatus)
		return errors.ServiceConnectionError(invalidStatus)
	}
}
