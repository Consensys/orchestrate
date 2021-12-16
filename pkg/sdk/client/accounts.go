package client

import (
	"context"
	"fmt"
	"strings"

	utilstypes "github.com/consensys/quorum-key-manager/src/utils/api/types"

	"github.com/consensys/orchestrate/pkg/types/api"
	qkmtypes "github.com/consensys/quorum-key-manager/src/stores/api/types"

	clientutils "github.com/consensys/orchestrate/pkg/toolkit/app/http/client-utils"
	"github.com/consensys/orchestrate/pkg/toolkit/app/http/httputil"
	"github.com/consensys/orchestrate/pkg/types/entities"
	ethcommon "github.com/ethereum/go-ethereum/common"
)

func (c *HTTPClient) GetAccount(ctx context.Context, address ethcommon.Address) (*api.AccountResponse, error) {
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

func (c *HTTPClient) UpdateAccount(ctx context.Context, address ethcommon.Address, req *api.UpdateAccountRequest) (*api.AccountResponse, error) {
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

func (c *HTTPClient) SignMessage(ctx context.Context, address ethcommon.Address, req *qkmtypes.SignMessageRequest) (string, error) {
	reqURL := fmt.Sprintf("%v/accounts/%s/sign-message", c.config.URL, address)

	response, err := clientutils.PostRequest(ctx, c.client, reqURL, req)
	if err != nil {
		return "", err
	}

	defer clientutils.CloseResponse(response)
	return httputil.ParseStringResponse(ctx, response)
}

func (c *HTTPClient) SignTypedData(ctx context.Context, address ethcommon.Address, request *qkmtypes.SignTypedDataRequest) (string, error) {
	reqURL := fmt.Sprintf("%v/accounts/%s/sign-typed-data", c.config.URL, address)

	response, err := clientutils.PostRequest(ctx, c.client, reqURL, request)
	if err != nil {
		return "", err
	}

	defer clientutils.CloseResponse(response)
	return httputil.ParseStringResponse(ctx, response)
}

func (c *HTTPClient) VerifyMessageSignature(ctx context.Context, request *utilstypes.VerifyRequest) error {
	reqURL := fmt.Sprintf("%v/accounts/verify-message", c.config.URL)

	response, err := clientutils.PostRequest(ctx, c.client, reqURL, request)
	if err != nil {
		return err
	}

	defer clientutils.CloseResponse(response)
	return httputil.ParseEmptyBodyResponse(ctx, response)
}

func (c *HTTPClient) VerifyTypedDataSignature(ctx context.Context, request *utilstypes.VerifyTypedDataRequest) error {
	reqURL := fmt.Sprintf("%v/accounts/verify-typed-data", c.config.URL)

	response, err := clientutils.PostRequest(ctx, c.client, reqURL, request)
	if err != nil {
		return err
	}

	defer clientutils.CloseResponse(response)
	return httputil.ParseEmptyBodyResponse(ctx, response)
}
