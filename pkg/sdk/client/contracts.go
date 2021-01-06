package client

import (
	"context"
	"fmt"

	"github.com/containous/traefik/v2/pkg/log"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	clientutils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/http/client-utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/http/httputil"
	types "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/api"
)

func (c *HTTPClient) RegisterContract(ctx context.Context, request *types.RegisterContractRequest) (*types.ContractResponse, error) {
	reqURL := fmt.Sprintf("%v/contracts", c.config.URL)
	resp := &types.ContractResponse{}

	err := callWithBackOff(ctx, c.config.backOff, func() error {
		response, err := clientutils.PostRequest(ctx, c.client, reqURL, request)
		if err != nil {
			errMessage := "error while registering contract"
			log.FromContext(ctx).WithError(err).Error(errMessage)
			return errors.ServiceConnectionError(errMessage).ExtendComponent(component)
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
			errMessage := "error while getting contract"
			log.FromContext(ctx).WithError(err).Error(errMessage)
			return errors.ServiceConnectionError(errMessage).ExtendComponent(component)
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
			errMessage := "error while getting contract catalog"
			log.FromContext(ctx).WithError(err).Error(errMessage)
			return errors.ServiceConnectionError(errMessage).ExtendComponent(component)
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
			errMessage := "error while getting contract tags"
			log.FromContext(ctx).WithError(err).Error(errMessage)
			return errors.ServiceConnectionError(errMessage).ExtendComponent(component)
		}

		defer clientutils.CloseResponse(response)
		return httputil.ParseResponse(ctx, response, &resp)
	})

	return resp, err
}

func (c *HTTPClient) DeregisterContract(_ context.Context, _, _ string) error {
	panic("method DeregisterContract is not implemented")
}

func (c *HTTPClient) SetContractAddressCodeHash(ctx context.Context, req *types.SetContractCodeHashRequest) error {
	reqURL := fmt.Sprintf("%v/contracts", c.config.URL)

	err := callWithBackOff(ctx, c.config.backOff, func() error {
		response, err := clientutils.PatchRequest(ctx, c.client, reqURL, req)
		if err != nil {
			errMessage := "error while setting contract code hash"
			log.FromContext(ctx).WithError(err).Error(errMessage)
			return errors.ServiceConnectionError(errMessage).ExtendComponent(component)
		}

		defer clientutils.CloseResponse(response)
		_, err = httputil.ParseStringResponse(ctx, response)
		return err
	})

	return err
}

func (c *HTTPClient) GetContractEventsBySigHash(ctx context.Context, address string, req *types.GetContractEventsBySignHashRequest) (*types.GetContractEventsBySignHashResponse, error) {
	reqURL := fmt.Sprintf("%v/contracts/%s/events?chain_id=%s&sig_hash=%s&indexed_input_count=%d", c.config.URL, address, req.ChainID, req.SigHash, req.IndexedInputCount)
	resp := &types.GetContractEventsBySignHashResponse{}

	err := callWithBackOff(ctx, c.config.backOff, func() error {
		response, err := clientutils.GetRequest(ctx, c.client, reqURL)
		if err != nil {
			errMessage := "error while getting contract events by sigHash"
			log.FromContext(ctx).WithError(err).Error(errMessage)
			return errors.ServiceConnectionError(errMessage).ExtendComponent(component)
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
			log.FromContext(ctx).WithError(err).Error(errMessage)
			return errors.ServiceConnectionError(errMessage).ExtendComponent(component)
		}

		defer clientutils.CloseResponse(response)
		return httputil.ParseResponse(ctx, response, &resp)
	})

	return resp, err
}
