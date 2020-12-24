package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	healthz "github.com/heptiolabs/healthcheck"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/chain-registry/store/models"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	chainsctrl "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/chain-registry/service/controllers/chains"
)

const (
	notFoundErrorMessage       = "chain does not exist"
	invalidRequestErrorMessage = "invalid request"
	failedToFetchErrorMessage  = "failed to fetch chain"
)

type HTTPClient struct {
	client *http.Client
	config *Config
}

func NewHTTPClient(h *http.Client, c *Config) *HTTPClient {
	return &HTTPClient{
		client: h,
		config: c,
	}
}

func (c HTTPClient) Checker() healthz.Check {
	return healthz.HTTPGetCheck(fmt.Sprintf("%s/live", c.config.MetricsURL), time.Second)
}

func (c *HTTPClient) GetChains(ctx context.Context) ([]*models.Chain, error) {
	reqURL := fmt.Sprintf("%v/chains", c.config.URL)

	response, err := c.getRequest(ctx, reqURL)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case http.StatusOK:
		defer closeResponse(response)

		chainsResult := []*models.Chain{}
		if err = json.NewDecoder(response.Body).Decode(&chainsResult); err != nil {
			return nil, errors.FromError(err).ExtendComponent(component)
		}

		return chainsResult, nil
	default:
		errMessage := "failed to fetch chains"
		log.WithContext(ctx).WithError(err).Error(errMessage)
		return nil, errors.ServiceConnectionError(errMessage)
	}
}

func (c *HTTPClient) GetChainByName(ctx context.Context, chainName string) (*models.Chain, error) {
	logger := log.WithContext(ctx).WithField("chain_name", chainName)
	reqURL := fmt.Sprintf("%v/chains?name=%s", c.config.URL, chainName)

	response, err := c.getRequest(ctx, reqURL)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case http.StatusOK:
		defer closeResponse(response)

		chainsResult := []*models.Chain{}
		if err = json.NewDecoder(response.Body).Decode(&chainsResult); err != nil {
			return nil, errors.FromError(err).ExtendComponent(component)
		}

		if len(chainsResult) == 0 {
			return nil, errors.NotFoundError("no chain found with name %s", chainName).ExtendComponent(component)
		}

		return chainsResult[0], nil
	case http.StatusNotFound:
		logger.WithError(err).Error(notFoundErrorMessage)
		return nil, errors.NotFoundError(notFoundErrorMessage)
	case http.StatusBadRequest:
		logger.WithError(err).Error(invalidRequestErrorMessage)
		return nil, errors.InvalidFormatError(invalidRequestErrorMessage)
	default:
		log.WithContext(ctx).WithError(err).WithField("chain_name", chainName).Error(failedToFetchErrorMessage)
		return nil, errors.ServiceConnectionError(failedToFetchErrorMessage)
	}
}

func (c *HTTPClient) GetChainByUUID(ctx context.Context, chainUUID string) (*models.Chain, error) {
	reqURL := fmt.Sprintf("%v/chains/%s", c.config.URL, chainUUID)

	response, err := c.getRequest(ctx, reqURL)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case http.StatusOK:
		defer closeResponse(response)

		chain := &models.Chain{}
		if err = json.NewDecoder(response.Body).Decode(chain); err != nil {
			return nil, errors.ServiceConnectionError(failedToFetchErrorMessage)
		}

		return chain, nil
	case http.StatusNotFound:
		return nil, errors.NotFoundError(notFoundErrorMessage).ExtendComponent(component)
	case http.StatusBadRequest:
		return nil, errors.InvalidFormatError(invalidRequestErrorMessage).ExtendComponent(component)
	default:
		return nil, errors.ServiceConnectionError(failedToFetchErrorMessage).ExtendComponent(component)
	}
}

func (c *HTTPClient) DeleteChainByUUID(ctx context.Context, chainUUID string) error {
	reqURL := fmt.Sprintf("%v/chains/%s", c.config.URL, chainUUID)

	response, err := c.deleteRequest(ctx, reqURL)
	if err != nil {
		return err
	}

	defer closeResponse(response)

	return nil
}

func (c *HTTPClient) RegisterChain(ctx context.Context, chain *models.Chain) (*models.Chain, error) {
	reqURL := fmt.Sprintf("%v/chains", c.config.URL)

	var fromBlock *string = nil
	if chain.ListenerStartingBlock != nil {
		fromBlock = &(&struct{ x string }{strconv.FormatUint(*chain.ListenerStartingBlock, 10)}).x
	}
	postReq := chainsctrl.PostRequest{
		Name: chain.Name,
		URLs: chain.URLs,
		Listener: &chainsctrl.ListenerPostRequest{
			BackOffDuration:   chain.ListenerBackOffDuration,
			FromBlock:         fromBlock,
			Depth:             chain.ListenerDepth,
			ExternalTxEnabled: chain.ListenerExternalTxEnabled,
		},
	}

	if len(chain.PrivateTxManagers) > 1 {
		err := errors.InvalidFormatError("maximum one element in PrivateTxManagers is allowed")
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	if len(chain.PrivateTxManagers) > 0 {
		postReq.PrivateTxManager = &chainsctrl.PrivateTxManagerRequest{
			URL:  chain.PrivateTxManagers[0].URL,
			Type: chain.PrivateTxManagers[0].Type,
		}
	}

	response, err := c.postRequest(ctx, reqURL, &postReq)
	if err != nil {
		return nil, err
	}

	defer closeResponse(response)

	respChain := &models.Chain{}
	if err := json.NewDecoder(response.Body).Decode(respChain); err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	return respChain, nil
}

func (c *HTTPClient) UpdateBlockPosition(ctx context.Context, chainUUID string, blockNumber uint64) error {
	reqURL := fmt.Sprintf("%v/chains/%v", c.config.URL, chainUUID)

	response, err := c.patchRequest(ctx, reqURL, &chainsctrl.PatchRequest{
		Listener: &chainsctrl.ListenerPatchRequest{CurrentBlock: &blockNumber},
	})
	if err != nil {
		return errors.FromError(err).ExtendComponent(component)
	}
	defer closeResponse(response)

	return nil
}

func (c *HTTPClient) UpdateChainByUUID(ctx context.Context, chainUUID string, chain *models.Chain) error {
	reqURL := fmt.Sprintf("%v/chains/%v", c.config.URL, chainUUID)

	patchReq := chainsctrl.PatchRequest{
		Name: chain.Name,
		URLs: chain.URLs,
		Listener: &chainsctrl.ListenerPatchRequest{
			BackOffDuration:   chain.ListenerBackOffDuration,
			Depth:             chain.ListenerDepth,
			CurrentBlock:      chain.ListenerCurrentBlock,
			ExternalTxEnabled: chain.ListenerExternalTxEnabled,
		},
	}

	if len(chain.PrivateTxManagers) > 1 {
		err := errors.InvalidFormatError("maximum one element in PrivateTxManagers is allowed")
		return errors.FromError(err).ExtendComponent(component)
	}

	if len(chain.PrivateTxManagers) > 0 {
		patchReq.PrivateTxManager = &chainsctrl.PrivateTxManagerRequest{
			URL:  chain.PrivateTxManagers[0].URL,
			Type: chain.PrivateTxManagers[0].Type,
		}
	}

	response, err := c.patchRequest(ctx, reqURL, &patchReq)
	if err != nil {
		return errors.FromError(err).ExtendComponent(component)
	}
	defer closeResponse(response)

	return nil
}

func (c *HTTPClient) getRequest(ctx context.Context, reqURL string) (*http.Response, error) {
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	r, err := c.client.Do(req)
	if err != nil {
		errMessage := "failed to execute get request"
		log.WithContext(ctx).WithError(err).Error(errMessage)
		return nil, errors.HTTPConnectionError(errMessage)
	}

	return r, nil
}

func (c *HTTPClient) deleteRequest(ctx context.Context, reqURL string) (*http.Response, error) {
	req, _ := http.NewRequestWithContext(ctx, http.MethodDelete, reqURL, nil)
	r, err := c.client.Do(req)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	if r.StatusCode != http.StatusNoContent {
		return nil, errors.FromError(fmt.Errorf("DELETE request: %s failed with error %d", reqURL, r.StatusCode)).ExtendComponent(component)
	}

	return r, nil
}

func (c *HTTPClient) postRequest(ctx context.Context, reqURL string, postRequest interface{}) (*http.Response, error) {
	body := new(bytes.Buffer)
	_ = json.NewEncoder(body).Encode(postRequest)

	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, body)
	r, err := c.client.Do(req)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	if r.StatusCode != http.StatusOK {
		if r.StatusCode == http.StatusBadRequest {
			return nil, errors.InvalidFormatError(body.String())
		}

		return nil, errors.FromError(fmt.Errorf("POST request: %s failed with error %d", reqURL, r.StatusCode)).ExtendComponent(component)
	}

	return r, nil
}

func (c *HTTPClient) patchRequest(ctx context.Context, reqURL string, patchRequest interface{}) (*http.Response, error) {
	body := new(bytes.Buffer)
	_ = json.NewEncoder(body).Encode(patchRequest)

	req, _ := http.NewRequestWithContext(ctx, http.MethodPatch, reqURL, body)
	r, err := c.client.Do(req)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	if r.StatusCode != http.StatusOK {
		if r.StatusCode == http.StatusBadRequest {
			return nil, errors.InvalidFormatError(body.String())
		}

		return nil, errors.FromError(fmt.Errorf("PATH request: %s failed with error %d", reqURL, r.StatusCode)).ExtendComponent(component)
	}

	return r, nil
}

func closeResponse(response *http.Response) {
	if deferErr := response.Body.Close(); deferErr != nil {
		log.WithError(deferErr).Errorf("%s: could not close body", component)
	}
}
