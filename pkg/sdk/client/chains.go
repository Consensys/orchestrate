package client

import (
	"context"
	"fmt"
	"strings"

	"github.com/containous/traefik/v2/pkg/log"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	clientutils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/http/client-utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/http/httputil"
	types "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/api"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
)

func (c *HTTPClient) GetChain(ctx context.Context, uuid string) (*types.ChainResponse, error) {
	reqURL := fmt.Sprintf("%v/chains/%s", c.config.URL, uuid)
	resp := &types.ChainResponse{}

	err := callWithBackOff(ctx, c.config.backOff, func() error {
		response, err := clientutils.GetRequest(ctx, c.client, reqURL)
		if err != nil {
			errMessage := "error while getting chain"
			log.FromContext(ctx).WithError(err).Error(errMessage)
			return errors.ServiceConnectionError(errMessage).ExtendComponent(component)
		}

		defer clientutils.CloseResponse(response)
		return httputil.ParseResponse(ctx, response, resp)
	})

	return resp, err
}

func (c *HTTPClient) SearchChains(ctx context.Context, filters *entities.ChainFilters) ([]*types.ChainResponse, error) {
	reqURL := fmt.Sprintf("%v/chains", c.config.URL)
	var resp []*types.ChainResponse

	var qParams []string
	if len(filters.Names) > 0 {
		qParams = append(qParams, "names="+strings.Join(filters.Names, ","))
	}

	if len(qParams) > 0 {
		reqURL = reqURL + "?" + strings.Join(qParams, "&")
	}

	err := callWithBackOff(ctx, c.config.backOff, func() error {
		response, err := clientutils.GetRequest(ctx, c.client, reqURL)
		if err != nil {
			errMessage := "error while searching chains"
			log.FromContext(ctx).WithError(err).Error(errMessage)
			return errors.ServiceConnectionError(errMessage).ExtendComponent(component)
		}
		defer clientutils.CloseResponse(response)
		return httputil.ParseResponse(ctx, response, &resp)
	})

	return resp, err
}

func (c *HTTPClient) RegisterChain(ctx context.Context, request *types.RegisterChainRequest) (*types.ChainResponse, error) {
	reqURL := fmt.Sprintf("%v/chains", c.config.URL)
	resp := &types.ChainResponse{}

	err := callWithBackOff(ctx, c.config.backOff, func() error {
		response, err := clientutils.PostRequest(ctx, c.client, reqURL, request)
		if err != nil {
			errMessage := "error while registering chain"
			log.FromContext(ctx).WithError(err).Error(errMessage)
			return errors.ServiceConnectionError(errMessage).ExtendComponent(component)
		}
		defer clientutils.CloseResponse(response)
		return httputil.ParseResponse(ctx, response, resp)
	})

	return resp, err
}

func (c *HTTPClient) UpdateChain(ctx context.Context, uuid string, request *types.UpdateChainRequest) (*types.ChainResponse, error) {
	reqURL := fmt.Sprintf("%v/chains/%v", c.config.URL, uuid)
	resp := &types.ChainResponse{}

	err := callWithBackOff(ctx, c.config.backOff, func() error {
		response, err := clientutils.PatchRequest(ctx, c.client, reqURL, request)
		if err != nil {
			errMessage := "error while updating chain"
			log.FromContext(ctx).WithError(err).Error(errMessage)
			return errors.ServiceConnectionError(errMessage).ExtendComponent(component)
		}

		defer clientutils.CloseResponse(response)
		return httputil.ParseResponse(ctx, response, resp)
	})

	return resp, err
}

func (c *HTTPClient) DeleteChain(ctx context.Context, uuid string) error {
	reqURL := fmt.Sprintf("%v/chains/%v", c.config.URL, uuid)

	response, err := clientutils.DeleteRequest(ctx, c.client, reqURL)
	if err != nil {
		errMessage := "error while deleting chain"
		log.FromContext(ctx).WithError(err).Error(errMessage)
		return errors.ServiceConnectionError(errMessage).ExtendComponent(component)
	}

	defer clientutils.CloseResponse(response)
	return httputil.ParseEmptyBodyResponse(ctx, response)
}
