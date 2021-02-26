package client

import (
	"context"
	"fmt"

	"github.com/ConsenSys/orchestrate/pkg/toolkit/app/http/httputil"
	types "github.com/ConsenSys/orchestrate/pkg/types/api"

	clientutils "github.com/ConsenSys/orchestrate/pkg/toolkit/app/http/client-utils"
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
