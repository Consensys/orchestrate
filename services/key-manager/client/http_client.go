package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/keymanager"

	"github.com/containous/traefik/v2/pkg/log"
	healthz "github.com/heptiolabs/healthcheck"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	clientutils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/client-utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/httputil"
	types "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/keymanager/ethereum"
)

func NewHTTPClient(h *http.Client, c *Config) KeyManagerClient {
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

func (c HTTPClient) ETHCreateAccount(ctx context.Context, req *types.CreateETHAccountRequest) (*types.ETHAccountResponse, error) {
	reqURL := fmt.Sprintf("%v/ethereum/accounts", c.config.URL)
	resp := &types.ETHAccountResponse{}

	response, err := clientutils.PostRequest(ctx, c.client, reqURL, req)
	if err != nil {
		errMessage := "error while creating ethereum account"
		log.FromContext(ctx).WithError(err).Error(errMessage)
		return nil, errors.ServiceConnectionError(errMessage).ExtendComponent(component)
	}

	defer clientutils.CloseResponse(response)
	if err := httputil.ParseResponse(ctx, response, resp); err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	return resp, nil
}

func (c HTTPClient) ETHImportAccount(ctx context.Context, req *types.ImportETHAccountRequest) (*types.ETHAccountResponse, error) {
	reqURL := fmt.Sprintf("%v/ethereum/accounts/import", c.config.URL)
	resp := &types.ETHAccountResponse{}

	response, err := clientutils.PostRequest(ctx, c.client, reqURL, req)
	if err != nil {
		errMessage := "error while importing ethereum account"
		log.FromContext(ctx).WithError(err).Error(errMessage)
		return nil, errors.ServiceConnectionError(errMessage).ExtendComponent(component)
	}

	defer clientutils.CloseResponse(response)
	if err := httputil.ParseResponse(ctx, response, resp); err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	return resp, nil
}

func (c HTTPClient) ETHSign(ctx context.Context, address string, req *keymanager.PayloadRequest) (string, error) {
	reqURL := fmt.Sprintf("%v/ethereum/accounts/%v/sign", c.config.URL, address)

	response, err := clientutils.PostRequest(ctx, c.client, reqURL, req)
	if err != nil {
		errMessage := "error while signing data with Ethereum account"
		log.FromContext(ctx).WithError(err).Error(errMessage)
		return "", errors.ServiceConnectionError(errMessage).ExtendComponent(component)
	}

	defer clientutils.CloseResponse(response)

	if response.StatusCode != http.StatusOK {
		errResp := httputil.ErrorResponse{}
		if err = json.NewDecoder(response.Body).Decode(&errResp); err != nil {
			errMessage := "failed to decode error response body"
			log.FromContext(ctx).WithError(err).Error(errMessage)
			return "", errors.ServiceConnectionError(errMessage)
		}

		return "", errors.Errorf(errResp.Code, errResp.Message)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		errMessage := "failed to decode response body"
		log.FromContext(ctx).WithError(err).Error(errMessage)
		return "", errors.ServiceConnectionError(errMessage)
	}

	return string(responseData), nil
}
