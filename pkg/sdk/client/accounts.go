package client

import (
	"context"
	"fmt"
	"strings"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/keymanager"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/keymanager/ethereum"

	"github.com/containous/traefik/v2/pkg/log"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	clientutils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/http/client-utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/http/httputil"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
	types "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/identitymanager"
)

func (c *HTTPClient) GetAccount(ctx context.Context, address string) (*types.AccountResponse, error) {
	reqURL := fmt.Sprintf("%v/accounts/%s", c.config.URL, address)
	resp := &types.AccountResponse{}

	response, err := clientutils.GetRequest(ctx, c.client, reqURL)
	if err != nil {
		errMessage := "error while getting account"
		log.FromContext(ctx).WithError(err).Error(errMessage)
		return nil, errors.ServiceConnectionError(errMessage).ExtendComponent(component)
	}

	defer clientutils.CloseResponse(response)
	if err := httputil.ParseResponse(ctx, response, resp); err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	return resp, nil
}

func (c *HTTPClient) CreateAccount(ctx context.Context, req *types.CreateAccountRequest) (*types.AccountResponse, error) {
	reqURL := fmt.Sprintf("%v/accounts", c.config.URL)
	resp := &types.AccountResponse{}

	response, err := clientutils.PostRequest(ctx, c.client, reqURL, req)
	if err != nil {
		errMessage := "error while creating account"
		log.FromContext(ctx).WithError(err).Error(errMessage)
		return nil, errors.ServiceConnectionError(errMessage).ExtendComponent(component)
	}

	defer clientutils.CloseResponse(response)
	if err := httputil.ParseResponse(ctx, response, resp); err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	return resp, nil
}

func (c *HTTPClient) ImportAccount(ctx context.Context, req *types.ImportAccountRequest) (*types.AccountResponse, error) {
	reqURL := fmt.Sprintf("%v/accounts/import", c.config.URL)
	resp := &types.AccountResponse{}

	response, err := clientutils.PostRequest(ctx, c.client, reqURL, req)
	if err != nil {
		errMessage := "error while importing account"
		log.FromContext(ctx).WithError(err).Error(errMessage)
		return nil, errors.ServiceConnectionError(errMessage).ExtendComponent(component)
	}

	defer clientutils.CloseResponse(response)
	if err := httputil.ParseResponse(ctx, response, resp); err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	return resp, nil
}

func (c *HTTPClient) UpdateAccount(ctx context.Context, address string, req *types.UpdateAccountRequest) (*types.AccountResponse, error) {
	reqURL := fmt.Sprintf("%v/accounts/%s", c.config.URL, address)
	resp := &types.AccountResponse{}

	response, err := clientutils.PatchRequest(ctx, c.client, reqURL, req)
	if err != nil {
		errMessage := "error while updating account"
		log.FromContext(ctx).WithError(err).Error(errMessage)
		return nil, errors.ServiceConnectionError(errMessage).ExtendComponent(component)
	}

	defer clientutils.CloseResponse(response)
	if err := httputil.ParseResponse(ctx, response, resp); err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	return resp, nil
}

func (c *HTTPClient) SearchAccounts(ctx context.Context, filters *entities.AccountFilters) ([]*types.AccountResponse, error) {
	reqURL := fmt.Sprintf("%v/accounts", c.config.URL)
	var resp []*types.AccountResponse

	var qParams []string
	if len(filters.Aliases) > 0 {
		qParams = append(qParams, "aliases="+strings.Join(filters.Aliases, ","))
	}

	if len(qParams) > 0 {
		reqURL = reqURL + "?" + strings.Join(qParams, "&")
	}

	response, err := clientutils.GetRequest(ctx, c.client, reqURL)
	if err != nil {
		errMessage := "error while searching accounts"
		log.FromContext(ctx).WithError(err).Error(errMessage)
		return nil, errors.ServiceConnectionError(errMessage).ExtendComponent(component)
	}

	defer clientutils.CloseResponse(response)
	if err := httputil.ParseResponse(ctx, response, &resp); err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	return resp, nil
}

func (c *HTTPClient) SignPayload(ctx context.Context, address string, req *types.SignPayloadRequest) (string, error) {
	reqURL := fmt.Sprintf("%v/accounts/%s/sign", c.config.URL, address)

	response, err := clientutils.PostRequest(ctx, c.client, reqURL, req)
	if err != nil {
		errMessage := "error while signing payload"
		log.FromContext(ctx).WithError(err).Error(errMessage)
		return "", errors.ServiceConnectionError(errMessage).ExtendComponent(component)
	}

	defer clientutils.CloseResponse(response)
	return httputil.ParseStringResponse(ctx, response)
}

func (c *HTTPClient) SignTypedData(ctx context.Context, address string, request *types.SignTypedDataRequest) (string, error) {
	reqURL := fmt.Sprintf("%v/accounts/%s/sign-typed-data", c.config.URL, address)

	response, err := clientutils.PostRequest(ctx, c.client, reqURL, request)
	if err != nil {
		errMessage := "error while signing typed data"
		log.FromContext(ctx).WithError(err).Error(errMessage)
		return "", errors.ServiceConnectionError(errMessage).ExtendComponent(component)
	}

	defer clientutils.CloseResponse(response)
	return httputil.ParseStringResponse(ctx, response)
}

func (c *HTTPClient) VerifySignature(ctx context.Context, request *keymanager.VerifyPayloadRequest) error {
	reqURL := fmt.Sprintf("%v/accounts/verify-signature", c.config.URL)

	response, err := clientutils.PostRequest(ctx, c.client, reqURL, request)
	if err != nil {
		errMessage := "error while verifying signature"
		log.FromContext(ctx).WithError(err).Error(errMessage)
		return errors.ServiceConnectionError(errMessage).ExtendComponent(component)
	}

	defer clientutils.CloseResponse(response)
	return httputil.ParseEmptyBodyResponse(ctx, response)
}

func (c *HTTPClient) VerifyTypedDataSignature(ctx context.Context, request *ethereum.VerifyTypedDataRequest) error {
	reqURL := fmt.Sprintf("%v/accounts/verify-typed-data-signature", c.config.URL)

	response, err := clientutils.PostRequest(ctx, c.client, reqURL, request)
	if err != nil {
		errMessage := "error while verifying typed data signature"
		log.FromContext(ctx).WithError(err).Error(errMessage)
		return errors.ServiceConnectionError(errMessage).ExtendComponent(component)
	}

	defer clientutils.CloseResponse(response)
	return httputil.ParseEmptyBodyResponse(ctx, response)
}
