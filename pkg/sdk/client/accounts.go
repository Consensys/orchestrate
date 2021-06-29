package client

import (
	"context"
	"fmt"
	"strings"

	qkmtypes "github.com/ConsenSys/orchestrate/pkg/quorum-key-manager/types"
	"github.com/ConsenSys/orchestrate/pkg/types/api"

	clientutils "github.com/ConsenSys/orchestrate/pkg/toolkit/app/http/client-utils"
	"github.com/ConsenSys/orchestrate/pkg/toolkit/app/http/httputil"
	"github.com/ConsenSys/orchestrate/pkg/types/entities"
)

func (c *HTTPClient) GetAccount(ctx context.Context, address string) (*api.AccountResponse, error) {
	reqURL := fmt.Sprintf("%v/accounts/%s", c.config.URL, address)
	resp := &api.AccountResponse{}

	response, err := clientutils.GetRequest(ctx, c.client, reqURL)
	if err != nil {
		return nil, err
	}

	defer clientutils.CloseResponse(response)
	if err := httputil.ParseResponse(ctx, response, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *HTTPClient) CreateAccount(ctx context.Context, req *api.CreateAccountRequest) (*api.AccountResponse, error) {
	reqURL := fmt.Sprintf("%v/accounts", c.config.URL)
	resp := &api.AccountResponse{}

	response, err := clientutils.PostRequest(ctx, c.client, reqURL, req)
	if err != nil {
		return nil, err
	}

	defer clientutils.CloseResponse(response)
	if err := httputil.ParseResponse(ctx, response, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *HTTPClient) ImportAccount(ctx context.Context, req *api.ImportAccountRequest) (*api.AccountResponse, error) {
	reqURL := fmt.Sprintf("%v/accounts/import", c.config.URL)
	resp := &api.AccountResponse{}

	response, err := clientutils.PostRequest(ctx, c.client, reqURL, req)
	if err != nil {
		return nil, err
	}

	defer clientutils.CloseResponse(response)
	if err := httputil.ParseResponse(ctx, response, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *HTTPClient) UpdateAccount(ctx context.Context, address string, req *api.UpdateAccountRequest) (*api.AccountResponse, error) {
	reqURL := fmt.Sprintf("%v/accounts/%s", c.config.URL, address)
	resp := &api.AccountResponse{}

	response, err := clientutils.PatchRequest(ctx, c.client, reqURL, req)
	if err != nil {
		return nil, err
	}

	defer clientutils.CloseResponse(response)
	if err := httputil.ParseResponse(ctx, response, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *HTTPClient) SearchAccounts(ctx context.Context, filters *entities.AccountFilters) ([]*api.AccountResponse, error) {
	reqURL := fmt.Sprintf("%v/accounts", c.config.URL)
	var resp []*api.AccountResponse

	var qParams []string
	if len(filters.Aliases) > 0 {
		qParams = append(qParams, "aliases="+strings.Join(filters.Aliases, ","))
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

func (c *HTTPClient) SignPayload(ctx context.Context, address string, req *api.SignPayloadRequest) (string, error) {
	reqURL := fmt.Sprintf("%v/accounts/%s/sign", c.config.URL, address)

	response, err := clientutils.PostRequest(ctx, c.client, reqURL, req)
	if err != nil {
		return "", err
	}

	defer clientutils.CloseResponse(response)
	return httputil.ParseStringResponse(ctx, response)
}

func (c *HTTPClient) SignTypedData(ctx context.Context, address string, request *api.SignTypedDataRequest) (string, error) {
	reqURL := fmt.Sprintf("%v/accounts/%s/sign-typed-data", c.config.URL, address)

	response, err := clientutils.PostRequest(ctx, c.client, reqURL, request)
	if err != nil {
		return "", err
	}

	defer clientutils.CloseResponse(response)
	return httputil.ParseStringResponse(ctx, response)
}

func (c *HTTPClient) VerifySignature(ctx context.Context, request *qkmtypes.VerifyEth1SignatureRequest) error {
	reqURL := fmt.Sprintf("%v/accounts/verify-signature", c.config.URL)

	response, err := clientutils.PostRequest(ctx, c.client, reqURL, request)
	if err != nil {
		return err
	}

	defer clientutils.CloseResponse(response)
	return httputil.ParseEmptyBodyResponse(ctx, response)
}

func (c *HTTPClient) VerifyTypedDataSignature(ctx context.Context, request *qkmtypes.VerifyTypedDataRequest) error {
	reqURL := fmt.Sprintf("%v/accounts/verify-typed-data-signature", c.config.URL)

	response, err := clientutils.PostRequest(ctx, c.client, reqURL, request)
	if err != nil {
		return err
	}

	defer clientutils.CloseResponse(response)
	return httputil.ParseEmptyBodyResponse(ctx, response)
}
