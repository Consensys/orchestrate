package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/containous/traefik/v2/pkg/log"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	clientutils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/client-utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/types"
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

func NewHTTPClient(h *http.Client, c *Config) *HTTPClient {
	return &HTTPClient{
		client: h,
		config: c,
	}
}

func (c *HTTPClient) SendTransaction(ctx context.Context, txRequest *types.TransactionRequest) (*types.TransactionResponse, error) {
	reqURL := fmt.Sprintf("%v/transactions/send", c.config.URL)

	response, err := clientutils.PostRequest(ctx, c.client, reqURL, txRequest)
	if err != nil {
		errMessage := "error while sending transaction"
		log.FromContext(ctx).WithError(err).Error(errMessage)
		return nil, errors.ServiceConnectionError(errMessage).ExtendComponent(component)
	}
	defer clientutils.CloseResponse(response)

	switch response.StatusCode {
	case http.StatusAccepted:
		resp := &types.TransactionResponse{}
		if err := json.NewDecoder(response.Body).Decode(resp); err != nil {
			log.FromContext(ctx).WithError(err).Error(invalidRequestBody)
			return nil, errors.ServiceConnectionError(invalidRequestBody).ExtendComponent(component)
		}

		return resp, nil
	case http.StatusBadRequest:
		log.FromContext(ctx).Error(invalidRequestErrorMessage)
		return nil, errors.InvalidFormatError(invalidRequestErrorMessage)
	case http.StatusConflict:
		errMessage := "transaction already exists"
		log.FromContext(ctx).Error(errMessage)
		return nil, errors.ConflictedError(errMessage)
	case http.StatusUnprocessableEntity:
		errMessage := "unprocessable transaction"
		log.FromContext(ctx).Error(errMessage)
		return nil, errors.InvalidParameterError(errMessage)
	default:
		log.FromContext(ctx).Error(invalidStatus)
		return nil, errors.ServiceConnectionError(invalidStatus)
	}
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

	switch response.StatusCode {
	case http.StatusOK:
		resp := &types.ScheduleResponse{}
		if err := json.NewDecoder(response.Body).Decode(resp); err != nil {
			log.FromContext(ctx).WithError(err).Error(invalidRequestBody)
			return nil, errors.ServiceConnectionError(invalidRequestBody).ExtendComponent(component)
		}

		return resp, nil
	case http.StatusNotFound:
		errMessage := "schedule not found"
		log.FromContext(ctx).Error(errMessage)
		return nil, errors.NotFoundError(errMessage)
	default:
		log.FromContext(ctx).Error(invalidStatus)
		return nil, errors.ServiceConnectionError(invalidStatus)
	}
}
