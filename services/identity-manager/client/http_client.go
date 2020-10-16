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

func (c HTTPClient) CreateIdentity(ctx context.Context, req *types.CreateIdentityRequest) (*types.IdentityResponse, error) {
	reqURL := fmt.Sprintf("%v/identities", c.config.URL)
	resp := &types.IdentityResponse{}

	response, err := clientutils.PostRequest(ctx, c.client, reqURL, req)
	if err != nil {
		errMessage := "error while creating identity"
		log.FromContext(ctx).WithError(err).Error(errMessage)
		return nil, errors.ServiceConnectionError(errMessage).ExtendComponent(component)
	}

	defer clientutils.CloseResponse(response)
	if err := httputil.ParseResponse(ctx, response, resp); err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	return resp, nil
}

func (c HTTPClient) ImportIdentity(ctx context.Context, req *types.ImportIdentityRequest) (*types.IdentityResponse, error) {
	reqURL := fmt.Sprintf("%v/identities/import", c.config.URL)
	resp := &types.IdentityResponse{}

	response, err := clientutils.PostRequest(ctx, c.client, reqURL, req)
	if err != nil {
		errMessage := "error while importing identity"
		log.FromContext(ctx).WithError(err).Error(errMessage)
		return nil, errors.ServiceConnectionError(errMessage).ExtendComponent(component)
	}

	defer clientutils.CloseResponse(response)
	if err := httputil.ParseResponse(ctx, response, resp); err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	return resp, nil
}
