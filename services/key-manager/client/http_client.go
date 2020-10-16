package client

import (
	"context"
	"fmt"
	"net/http"
	"time"

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

func (c HTTPClient) CreateETHAccount(ctx context.Context, req *types.CreateETHAccountRequest) (*types.ETHAccountResponse, error) {
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

func (c HTTPClient) ImportETHAccount(ctx context.Context, req *types.ImportETHAccountRequest) (*types.ETHAccountResponse, error) {
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
