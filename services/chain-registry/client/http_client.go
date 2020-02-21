package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/api/chains"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/types"
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

func (c *HTTPClient) GetChains(ctx context.Context) ([]*types.Chain, error) {
	reqURL := fmt.Sprintf("%v/chains", c.config.URL)

	response, err := c.getRequest(ctx, reqURL)
	if err != nil {
		return nil, err
	}
	defer closeResponse(response)

	chainsResult := []*types.Chain{}
	if err := json.NewDecoder(response.Body).Decode(&chainsResult); err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	return chainsResult, nil
}

func (c *HTTPClient) GetChainByName(ctx context.Context, chainName string) (*types.Chain, error) {
	reqURL := fmt.Sprintf("%v/chains?name=%s", c.config.URL, chainName)

	response, err := c.getRequest(ctx, reqURL)
	if err != nil {
		return nil, err
	}
	defer closeResponse(response)

	chainsResult := []*types.Chain{}
	if err := json.NewDecoder(response.Body).Decode(&chainsResult); err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	if len(chainsResult) == 0 {
		return nil, errors.FromError(fmt.Errorf("no chain found with name %s", chainName)).ExtendComponent(component)
	}

	return chainsResult[0], nil
}

func (c *HTTPClient) GetChainByUUID(ctx context.Context, chainUUID string) (*types.Chain, error) {
	reqURL := fmt.Sprintf("%v/chains/%s", c.config.URL, chainUUID)

	response, err := c.getRequest(ctx, reqURL)
	if err != nil {
		return nil, err
	}
	defer closeResponse(response)

	chain := &types.Chain{}
	if err := json.NewDecoder(response.Body).Decode(chain); err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	return chain, nil
}

func (c *HTTPClient) UpdateBlockPosition(ctx context.Context, chainUUID string, blockNumber int64) error {
	reqURL := fmt.Sprintf("%v/chains/%v", c.config.URL, chainUUID)
	body := new(bytes.Buffer)
	_ = json.NewEncoder(body).Encode(&chains.PatchRequest{
		Listener: &chains.Listener{BlockPosition: &blockNumber},
	})

	req, _ := http.NewRequestWithContext(ctx, http.MethodPatch, reqURL, body)
	r, err := c.client.Do(req)
	if err != nil {
		return errors.FromError(err).ExtendComponent(component)
	}
	defer closeResponse(r)

	if r.StatusCode != http.StatusOK {
		buf := new(bytes.Buffer)
		_, _ = buf.ReadFrom(r.Body)
		return errors.FromError(fmt.Errorf("could not update block position %s - got %d - body: %s", reqURL, r.StatusCode, buf.String())).ExtendComponent(component)
	}

	return nil
}

func (c *HTTPClient) getRequest(ctx context.Context, reqURL string) (*http.Response, error) {
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	r, err := c.client.Do(req)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	if r.StatusCode != http.StatusOK {
		return nil, errors.FromError(fmt.Errorf("get request: %s failed with error %d", reqURL, r.StatusCode)).ExtendComponent(component)
	}

	return r, nil
}

func closeResponse(response *http.Response) {
	if deferErr := response.Body.Close(); deferErr != nil {
		log.WithError(deferErr).Errorf("%s: could not close body", component)
	}
}
