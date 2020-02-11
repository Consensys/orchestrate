package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

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
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	r, err := c.client.Do(req)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}
	defer func() {
		if deferErr := r.Body.Close(); err != nil {
			log.WithError(deferErr).Errorf("%s: could close body", component)
		}
	}()

	if r.StatusCode != http.StatusOK {
		return nil, errors.FromError(fmt.Errorf("could not get chains %s - got %d", reqURL, r.StatusCode)).ExtendComponent(component)
	}

	var chns []*types.Chain
	if err := json.NewDecoder(r.Body).Decode(&chns); err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	return chns, nil
}

func (c *HTTPClient) GetChainByUUID(ctx context.Context, chainUUID string) (*types.Chain, error) {
	reqURL := fmt.Sprintf("%v/chains/%v", c.config.URL, chainUUID)
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	r, err := c.client.Do(req)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}
	defer func() {
		if deferErr := r.Body.Close(); err != nil {
			log.WithError(deferErr).Errorf("%s: could close body", component)
		}
	}()

	if r.StatusCode != http.StatusOK {
		return nil, errors.FromError(fmt.Errorf("could not get chain %s - got %d", reqURL, r.StatusCode)).ExtendComponent(component)
	}

	chain := &types.Chain{}
	if err := json.NewDecoder(r.Body).Decode(chain); err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}
	return chain, nil
}

func (c *HTTPClient) GetChainByTenantAndUUID(ctx context.Context, tenantID, chainUUID string) (*types.Chain, error) {
	baseURL, _ := url.Parse(c.config.URL)
	baseURL.Path = fmt.Sprintf("%s/chains", tenantID)
	params := url.Values{}
	params.Add("uuid", chainUUID)
	baseURL.RawQuery = params.Encode()
	reqURL := baseURL.String()

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	r, err := c.client.Do(req)
	if err != nil {
		return nil, errors.FromError(fmt.Errorf("%v - url: %s", err, reqURL)).ExtendComponent(component)
	}
	defer func() {
		if deferErr := r.Body.Close(); err != nil {
			log.WithError(deferErr).Errorf("%s: could close body", component)
		}
	}()

	if r.StatusCode != http.StatusOK {
		return nil, errors.FromError(fmt.Errorf("could not get chain %s - got %d", reqURL, r.StatusCode)).ExtendComponent(component)
	}

	chns := make([]*types.Chain, 0)
	if err := json.NewDecoder(r.Body).Decode(&chns); err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}
	if len(chns) != 1 {
		return nil, errors.FromError(fmt.Errorf("did not expected to get many chains with same for tenantID:%s and uuid:%s  from the chain registry - %s", tenantID, chainUUID, reqURL)).ExtendComponent(component)
	}
	return chns[0], nil
}

func (c *HTTPClient) GetChainByTenantAndName(ctx context.Context, tenantID, chainName string) (*types.Chain, error) {
	reqURL := fmt.Sprintf("%s/%s/chains/%s", c.config.URL, tenantID, chainName)

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	r, err := c.client.Do(req)
	if err != nil {
		return nil, errors.FromError(fmt.Errorf("%v - url: %s", err, reqURL)).ExtendComponent(component)
	}
	defer func() {
		if deferErr := r.Body.Close(); err != nil {
			log.WithError(deferErr).Errorf("%s: could close body", component)
		}
	}()

	if r.StatusCode != http.StatusOK {
		return nil, errors.FromError(fmt.Errorf("could not get chain %s - got %d", reqURL, r.StatusCode)).ExtendComponent(component)
	}

	chain := &types.Chain{}
	if err := json.NewDecoder(r.Body).Decode(chain); err != nil {
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
	defer func() {
		if deferErr := r.Body.Close(); err != nil {
			log.WithError(deferErr).Errorf("%s: could close body", component)
		}
	}()

	if r.StatusCode != http.StatusOK {
		buf := new(bytes.Buffer)
		_, _ = buf.ReadFrom(r.Body)
		return errors.FromError(fmt.Errorf("could not update block position %s - got %d - body: %s", reqURL, r.StatusCode, buf.String())).ExtendComponent(component)
	}

	return nil
}
