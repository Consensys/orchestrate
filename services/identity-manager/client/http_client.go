package client

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/containous/traefik/v2/pkg/log"
	healthz "github.com/heptiolabs/healthcheck"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	clientutils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/client-utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/httputil"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/entities"
	types "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/identitymanager"
)

func NewHTTPClient(h *http.Client, c *Config) IdentityManagerClient {
	return &HTTPClient{
		client: h,
		config: c,
	}
}

type HTTPClient struct {
	client *http.Client
	config *Config
}

func (c HTTPClient) Checker() healthz.Check {
	return healthz.HTTPGetCheck(fmt.Sprintf("%s/live", c.config.MetricsURL), time.Second)
}

func (c HTTPClient) GetAccount(ctx context.Context, address string) (*types.AccountResponse, error) {
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

func (c HTTPClient) CreateAccount(ctx context.Context, req *types.CreateAccountRequest) (*types.AccountResponse, error) {
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

func (c HTTPClient) ImportAccount(ctx context.Context, req *types.ImportAccountRequest) (*types.AccountResponse, error) {
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

func (c HTTPClient) UpdateAccount(ctx context.Context, address string, req *types.UpdateAccountRequest) (*types.AccountResponse, error) {
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

func (c HTTPClient) SearchAccounts(ctx context.Context, filters *entities.AccountFilters) ([]*types.AccountResponse, error) {
	reqURL := fmt.Sprintf("%v/accounts", c.config.URL)
	resp := []*types.AccountResponse{}

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

func (c HTTPClient) SignPayload(ctx context.Context, address string, req *types.SignPayloadRequest) (string, error) {
	reqURL := fmt.Sprintf("%v/accounts/%s/sign", c.config.URL, address)

	response, err := clientutils.PostRequest(ctx, c.client, reqURL, req)
	if err != nil {
		errMessage := "error while signing payload"
		log.FromContext(ctx).WithError(err).Error(errMessage)
		return "", errors.ServiceConnectionError(errMessage).ExtendComponent(component)
	}

	defer clientutils.CloseResponse(response)
	signature, err := ioutil.ReadAll(response.Body)
	if err != nil {
		errMessage := "failed to decode response body"
		log.FromContext(ctx).WithError(err).Error(errMessage)
		return "", errors.ServiceConnectionError(errMessage).ExtendComponent(component)
	}

	return string(signature), nil
}
