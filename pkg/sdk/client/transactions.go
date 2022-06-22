package client

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/consensys/orchestrate/pkg/types/entities"

	"github.com/consensys/orchestrate/pkg/toolkit/app/http/httputil"
	types "github.com/consensys/orchestrate/pkg/types/api"

	clientutils "github.com/consensys/orchestrate/pkg/toolkit/app/http/client-utils"
)

func (c *HTTPClient) SendContractTransaction(ctx context.Context, txRequest *types.SendTransactionRequest) (*types.TransactionResponse, error) {
	reqURL := fmt.Sprintf("%v/transactions/send", c.config.URL)
	resp := &types.TransactionResponse{}

	err := callWithBackOff(ctx, c.config.backOff, func() error {
		response, err := clientutils.PostRequest(ctx, c.client, reqURL, txRequest)
		if err != nil {
			return err
		}

		defer clientutils.CloseResponse(response)
		return httputil.ParseResponse(ctx, response, resp)
	})

	return resp, err
}

func (c *HTTPClient) SendDeployTransaction(ctx context.Context, txRequest *types.DeployContractRequest) (*types.TransactionResponse, error) {
	reqURL := fmt.Sprintf("%v/transactions/deploy-contract", c.config.URL)
	resp := &types.TransactionResponse{}

	err := callWithBackOff(ctx, c.config.backOff, func() error {
		response, err := clientutils.PostRequest(ctx, c.client, reqURL, txRequest)
		if err != nil {
			return err
		}

		defer clientutils.CloseResponse(response)
		return httputil.ParseResponse(ctx, response, resp)
	})

	return resp, err
}

func (c *HTTPClient) SendRawTransaction(ctx context.Context, txRequest *types.RawTransactionRequest) (*types.TransactionResponse, error) {
	reqURL := fmt.Sprintf("%v/transactions/send-raw", c.config.URL)
	resp := &types.TransactionResponse{}

	err := callWithBackOff(ctx, c.config.backOff, func() error {
		response, err := clientutils.PostRequest(ctx, c.client, reqURL, txRequest)
		if err != nil {
			return err
		}

		defer clientutils.CloseResponse(response)
		return httputil.ParseResponse(ctx, response, resp)
	})

	return resp, err
}

func (c *HTTPClient) SendTransferTransaction(ctx context.Context, txRequest *types.TransferRequest) (*types.TransactionResponse, error) {
	reqURL := fmt.Sprintf("%v/transactions/transfer", c.config.URL)
	resp := &types.TransactionResponse{}

	err := callWithBackOff(ctx, c.config.backOff, func() error {
		response, err := clientutils.PostRequest(ctx, c.client, reqURL, txRequest)
		if err != nil {
			return err
		}

		defer clientutils.CloseResponse(response)
		return httputil.ParseResponse(ctx, response, resp)
	})

	return resp, err
}

func (c *HTTPClient) GetTxRequest(ctx context.Context, txRequestUUID string) (*types.TransactionResponse, error) {
	reqURL := fmt.Sprintf("%v/transactions/%v", c.config.URL, txRequestUUID)
	resp := &types.TransactionResponse{}

	err := callWithBackOff(ctx, c.config.backOff, func() error {
		response, err := clientutils.GetRequest(ctx, c.client, reqURL)
		if err != nil {
			return err
		}

		defer clientutils.CloseResponse(response)
		return httputil.ParseResponse(ctx, response, resp)
	})

	return resp, err
}

func (c *HTTPClient) CallOffTransaction(ctx context.Context, txRequestUUID string) (*types.TransactionResponse, error) {
	reqURL := fmt.Sprintf("%v/transactions/%v/call-off", c.config.URL, txRequestUUID)
	resp := &types.TransactionResponse{}

	err := callWithBackOff(ctx, c.config.backOff, func() error {
		response, err := clientutils.PutRequest(ctx, c.client, reqURL, nil)
		if err != nil {
			return err
		}

		defer clientutils.CloseResponse(response)
		return httputil.ParseResponse(ctx, response, resp)
	})

	return resp, err
}

func (c *HTTPClient) SpeedUpTransaction(ctx context.Context, txRequestUUID string, boost *float64) (*types.TransactionResponse, error) {
	reqURL := fmt.Sprintf("%v/transactions/%v/speed-up", c.config.URL, txRequestUUID)
	if boost != nil {
		reqURL += fmt.Sprintf("?boost=%f", *boost)
	}

	resp := &types.TransactionResponse{}

	err := callWithBackOff(ctx, c.config.backOff, func() error {
		response, err := clientutils.PutRequest(ctx, c.client, reqURL, nil)
		if err != nil {
			return err
		}

		defer clientutils.CloseResponse(response)
		return httputil.ParseResponse(ctx, response, resp)
	})

	return resp, err
}

func (c *HTTPClient) SearchTransaction(ctx context.Context, filters *entities.TransactionRequestFilters) (*types.TransactionSearchResponse, error) {
	reqURL := fmt.Sprintf("%v/transactions", c.config.URL)
	var resp *types.TransactionSearchResponse

	var qParams []string
	if len(filters.IdempotencyKeys) > 0 {
		qParams = append(qParams, "idempotency_keys="+strings.Join(filters.IdempotencyKeys, ","))
	}

	if filters.Pagination.Limit > 0 {
		qParams = append(qParams, "limit="+strconv.Itoa(filters.Pagination.Limit))
	}

	if filters.Pagination.Page > 0 {
		qParams = append(qParams, "page="+strconv.Itoa(filters.Pagination.Page))
	}

	if len(qParams) > 0 {
		reqURL = reqURL + "?" + strings.Join(qParams, "&")
	}

	response, err := clientutils.GetRequest(ctx, c.client, reqURL)
	if err != nil {
		return nil, err
	}

	defer clientutils.CloseResponse(response)
	if err := httputil.ParseResponse(ctx, response, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *HTTPClient) ResumeTransaction(ctx context.Context, txRequestUUID string) (*types.TransactionResponse, error) {
	reqURL := fmt.Sprintf("%v/transactions/%v/resume", c.config.URL, txRequestUUID)
	resp := &types.TransactionResponse{}

	err := callWithBackOff(ctx, c.config.backOff, func() error {
		response, err := clientutils.PutRequest(ctx, c.client, reqURL, nil)
		if err != nil {
			return err
		}

		defer clientutils.CloseResponse(response)
		return httputil.ParseResponse(ctx, response, resp)
	})

	return resp, err
}
