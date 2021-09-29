package client

import (
	"context"
	"fmt"

	"github.com/consensys/orchestrate/pkg/errors"
	clientutils "github.com/consensys/orchestrate/pkg/toolkit/app/http/client-utils"
	"github.com/consensys/orchestrate/pkg/toolkit/app/http/httputil"
	types "github.com/consensys/orchestrate/pkg/types/api"
)

func (c *HTTPClient) RegisterContract(ctx context.Context, request *types.RegisterContractRequest) (*types.ContractResponse, error) {
	reqURL := fmt.Sprintf("%v/contracts", c.config.URL)
	resp := &types.ContractResponse{}

	err := callWithBackOff(ctx, c.config.backOff, func() error {
		response, err := clientutils.PostRequest(ctx, c.client, reqURL, request)
		if err != nil {
			return err
		}
		defer clientutils.CloseResponse(response)
		return httputil.ParseResponse(ctx, response, resp)
	})

	return resp, err
}

func (c *HTTPClient) GetContract(ctx context.Context, name, tag string) (*types.ContractResponse, error) {
	reqURL := fmt.Sprintf("%v/contracts/%s/%s", c.config.URL, name, tag)
	resp := &types.ContractResponse{}

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

func (c *HTTPClient) GetContractsCatalog(ctx context.Context) ([]string, error) {
	reqURL := fmt.Sprintf("%v/contracts", c.config.URL)
	var resp []string

	err := callWithBackOff(ctx, c.config.backOff, func() error {
		response, err := clientutils.GetRequest(ctx, c.client, reqURL)
		if err != nil {
			return err
		}

		defer clientutils.CloseResponse(response)
		return httputil.ParseResponse(ctx, response, &resp)
	})

	return resp, err
}

func (c *HTTPClient) GetContractTags(ctx context.Context, name string) ([]string, error) {
	reqURL := fmt.Sprintf("%v/contracts/%s", c.config.URL, name)
	var resp []string

	err := callWithBackOff(ctx, c.config.backOff, func() error {
		response, err := clientutils.GetRequest(ctx, c.client, reqURL)
		if err != nil {
			return err
		}

		defer clientutils.CloseResponse(response)
		return httputil.ParseResponse(ctx, response, &resp)
	})

	return resp, err
}

func (c *HTTPClient) DeregisterContract(_ context.Context, _, _ string) error {
	panic("method DeregisterContract is not implemented")
}

func (c *HTTPClient) SetContractAddressCodeHash(ctx context.Context, address, chainID string, req *types.SetContractCodeHashRequest) error {
	reqURL := fmt.Sprintf("%v/contracts/accounts/%s/%s", c.config.URL, chainID, address)

	err := callWithBackOff(ctx, c.config.backOff, func() error {
		response, err := clientutils.PostRequest(ctx, c.client, reqURL, req)
		if err != nil {
			return err
		}

		defer clientutils.CloseResponse(response)
		_, err = httputil.ParseStringResponse(ctx, response)
		return err
	})

	return err
}

func (c *HTTPClient) GetContractEvents(ctx context.Context, address, chainID string, req *types.GetContractEventsRequest) (*types.GetContractEventsBySignHashResponse, error) {
	reqURL := fmt.Sprintf("%v/contracts/accounts/%s/%s/events?&sig_hash=%s&indexed_input_count=%d", c.config.URL, chainID, address, req.SigHash, req.IndexedInputCount)
	resp := &types.GetContractEventsBySignHashResponse{}

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

func (c *HTTPClient) GetContractMethodSignatures(ctx context.Context, name, tag, method string) ([]string, error) {
	reqURL := fmt.Sprintf("%v/contracts/%s/%s/method-signatures", c.config.URL, name, tag)
	var resp []string

	if method != "" {
		reqURL = reqURL + "?method=" + method
	}

	err := callWithBackOff(ctx, c.config.backOff, func() error {
		response, err := clientutils.GetRequest(ctx, c.client, reqURL)
		if err != nil {
			errMessage := "error while getting contract method signatures"
			return errors.FromError(err).SetMessage(errMessage).AppendReason(err.Error()).ExtendComponent(component)
		}

		defer clientutils.CloseResponse(response)
		return httputil.ParseResponse(ctx, response, &resp)
	})

	return resp, err
}
