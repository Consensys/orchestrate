package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/httputil"
	types "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/chainregistry"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/models"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	chainsctrl "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/service/controllers/chains"
	facuetsctrl "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/service/controllers/faucets"
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
		log.WithContext(ctx).WithError(err).WithField("chain_name", chainName).Error(notFoundErrorMessage)
		return nil, errors.NotFoundError(notFoundErrorMessage)
	case http.StatusBadRequest:
		log.WithContext(ctx).WithError(err).WithField("chain_name", chainName).Error(invalidRequestErrorMessage)
		return nil, errors.InvalidParameterError(invalidRequestErrorMessage)
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
		err := errors.DataError("maximum one element in PrivateTxManagers is allowed")
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
		err := errors.DataError("maximum one element in PrivateTxManagers is allowed")
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

func (c *HTTPClient) GetFaucetByUUID(ctx context.Context, faucetUUID string) (*models.Faucet, error) {
	reqURL := fmt.Sprintf("%v/faucets/%s", c.config.URL, faucetUUID)

	response, err := c.getRequest(ctx, reqURL)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case http.StatusOK:
		defer closeResponse(response)

		faucet := &models.Faucet{}
		if err = json.NewDecoder(response.Body).Decode(faucet); err != nil {
			return nil, errors.FromError(err).ExtendComponent(component)
		}

		return faucet, nil
	case http.StatusNotFound:
		log.WithContext(ctx).WithError(err).WithField("faucet_uuid", faucetUUID).Error(notFoundErrorMessage)
		return nil, errors.NotFoundError("faucet does not exist")
	case http.StatusBadRequest:
		log.WithContext(ctx).WithError(err).WithField("faucet_uuid", faucetUUID).Error(invalidRequestErrorMessage)
		return nil, errors.InvalidParameterError(invalidRequestErrorMessage)
	default:
		log.WithContext(ctx).WithError(err).WithField("faucet_uuid", faucetUUID).Error(failedToFetchErrorMessage)
		return nil, errors.ServiceConnectionError("failed to fetch faucet")
	}
}

func (c *HTTPClient) GetFaucetsByChainRule(ctx context.Context, chainRule string) ([]*models.Faucet, error) {
	reqURL := fmt.Sprintf("%v/faucets?chain_rule=%s", c.config.URL, chainRule)

	response, err := c.getRequest(ctx, reqURL)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case http.StatusOK:
		defer closeResponse(response)

		var faucetsResult []*models.Faucet
		if err = json.NewDecoder(response.Body).Decode(&faucetsResult); err != nil {
			return nil, errors.FromError(err).ExtendComponent(component)
		}

		return faucetsResult, nil
	default:
		errMessage := "failed to fetch faucets"
		log.WithContext(ctx).WithError(err).Error(errMessage)
		return nil, errors.ServiceConnectionError(errMessage)
	}
}

func (c *HTTPClient) GetFaucetCandidate(ctx context.Context, account ethcommon.Address, chainUUID string) (*types.Faucet, error) {
	reqURL := fmt.Sprintf("%v/faucets/candidate?chain_uuid=%s&account=%s", c.config.URL, chainUUID, account.Hex())

	// @TODO API-KEY
	response, err := c.getRequest(ctx, reqURL)
	if err != nil {
		return nil, err
	}
	defer closeResponse(response)

	switch response.StatusCode {
	case http.StatusOK:
		var faucetsResult *types.Faucet
		if err = json.NewDecoder(response.Body).Decode(&faucetsResult); err != nil {
			return nil, errors.FromError(err).ExtendComponent(component)
		}

		return faucetsResult, nil
	case http.StatusNotFound:
		return nil, nil
	default:
		errResp := httputil.ErrorResponse{}
		if err = json.NewDecoder(response.Body).Decode(&errResp); err != nil {
			return nil, errors.FromError(err).ExtendComponent(component)
		}
		errMessage := "failed to fetch faucet candidate"
		log.WithContext(ctx).WithError(fmt.Errorf(errResp.Message)).Error(errMessage)
		return nil, errors.ServiceConnectionError(errMessage)
	}
}

func (c *HTTPClient) DeleteFaucetByUUID(ctx context.Context, faucetUUID string) error {
	reqURL := fmt.Sprintf("%v/faucets/%s", c.config.URL, faucetUUID)

	response, err := c.deleteRequest(ctx, reqURL)
	if err != nil {
		return err
	}

	defer closeResponse(response)

	return nil
}

func (c *HTTPClient) RegisterFaucet(ctx context.Context, faucet *models.Faucet) (*models.Faucet, error) {
	reqURL := fmt.Sprintf("%v/faucets", c.config.URL)

	postReq := facuetsctrl.PostRequest{
		Name:            faucet.Name,
		Amount:          faucet.Amount,
		ChainRule:       faucet.ChainRule,
		MaxBalance:      faucet.MaxBalance,
		CreditorAccount: faucet.CreditorAccount,
		Cooldown:        faucet.Cooldown,
	}

	response, err := c.postRequest(ctx, reqURL, &postReq)

	if err != nil {
		return nil, err
	}

	defer closeResponse(response)

	resp := &models.Faucet{}
	if err := json.NewDecoder(response.Body).Decode(resp); err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	return resp, nil
}

func (c *HTTPClient) UpdateFaucetByUUID(ctx context.Context, uuid string, faucet *models.Faucet) (*models.Faucet, error) {
	reqURL := fmt.Sprintf("%v/faucets/%v", c.config.URL, uuid)

	postReq := facuetsctrl.PatchRequest{
		Name:            faucet.Name,
		Amount:          faucet.Amount,
		ChainRule:       faucet.ChainRule,
		MaxBalance:      faucet.MaxBalance,
		CreditorAccount: faucet.CreditorAccount,
		Cooldown:        faucet.Cooldown,
	}

	response, err := c.patchRequest(ctx, reqURL, &postReq)

	if err != nil {
		return nil, err
	}

	defer closeResponse(response)

	resp := &models.Faucet{}
	if err := json.NewDecoder(response.Body).Decode(resp); err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	return resp, nil
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
			return nil, errors.DataError(body.String())
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
			return nil, errors.DataError(body.String())
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
