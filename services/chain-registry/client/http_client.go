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
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/api"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/types"
)

type Config struct {
	URL string
}

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

func (c *HTTPClient) GetNodes(ctx context.Context) ([]*types.Node, error) {
	reqURL := fmt.Sprintf("%v/nodes", c.config.URL)
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
		return nil, errors.FromError(fmt.Errorf("could not get nodes %s - got %d", reqURL, r.StatusCode)).ExtendComponent(component)
	}

	var nodes []*types.Node
	if err := json.NewDecoder(r.Body).Decode(&nodes); err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	return nodes, nil
}

func (c *HTTPClient) GetNodeByID(ctx context.Context, nodeID string) (*types.Node, error) {
	reqURL := fmt.Sprintf("%v/nodes/%v", c.config.URL, nodeID)
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
		return nil, errors.FromError(fmt.Errorf("could not get node %s - got %d", reqURL, r.StatusCode)).ExtendComponent(component)
	}

	node := &types.Node{}
	if err := json.NewDecoder(r.Body).Decode(node); err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}
	return node, nil
}

func (c *HTTPClient) GetNodeByTenantAndNodeID(ctx context.Context, tenantID, nodeID string) (*types.Node, error) {
	reqURL := fmt.Sprintf("%s/%s/nodes/%s", c.config.URL, tenantID, nodeID)
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
		return nil, errors.FromError(fmt.Errorf("could not get node %s - got %d", reqURL, r.StatusCode)).ExtendComponent(component)
	}

	node := &types.Node{}
	if err := json.NewDecoder(r.Body).Decode(node); err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}
	return node, nil
}

func (c *HTTPClient) GetNodeByTenantAndNodeName(ctx context.Context, tenantID, nodeName string) (*types.Node, error) {
	baseURL, _ := url.Parse(c.config.URL)
	baseURL.Path = fmt.Sprintf("%s/nodes", tenantID)
	params := url.Values{}
	params.Add("name", nodeName)
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
		return nil, errors.FromError(fmt.Errorf("could not get node %s - got %d", reqURL, r.StatusCode)).ExtendComponent(component)
	}

	nodes := make([]*types.Node, 0)
	if err := json.NewDecoder(r.Body).Decode(&nodes); err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}
	if len(nodes) != 1 {
		return nil, errors.FromError(fmt.Errorf("did not expected to get many nodes with same for tenantID:%s and name:%s  from the chain registry - %s", tenantID, nodeName, reqURL)).ExtendComponent(component)
	}

	return nodes[0], nil
}

func (c *HTTPClient) UpdateBlockPosition(ctx context.Context, nodeID string, blockNumber int64) error {
	body := new(bytes.Buffer)
	_ = json.NewEncoder(body).Encode(&api.PatchBlockPositionRequest{
		BlockPosition: blockNumber,
	})

	reqURL := fmt.Sprintf("%v/nodes/%v/block-position", c.config.URL, nodeID)
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
