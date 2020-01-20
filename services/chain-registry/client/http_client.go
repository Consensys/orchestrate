package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

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

func (c *HTTPClient) GetNodeByID(ctx context.Context, nodeID string) (*types.Node, error) {
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%v/nodes/%v", c.config.URL, nodeID), nil)
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
		return nil, errors.FromError(fmt.Errorf("could not get node - got %d", r.StatusCode)).ExtendComponent(component)
	}

	node := &types.Node{}
	if err := json.NewDecoder(r.Body).Decode(node); err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}
	return node, nil
}

func (c *HTTPClient) GetNodes(ctx context.Context) ([]*types.Node, error) {
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%v/nodes", c.config.URL), nil)
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
		return nil, errors.FromError(fmt.Errorf("could not get nodes - got %d", r.StatusCode)).ExtendComponent(component)
	}

	var nodes []*types.Node
	if err := json.NewDecoder(r.Body).Decode(&nodes); err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	return nodes, nil
}

func (c *HTTPClient) UpdateBlockPosition(ctx context.Context, nodeID string, blockNumber int64) error {
	body := new(bytes.Buffer)
	_ = json.NewEncoder(body).Encode(&api.PatchBlockPositionRequest{
		BlockPosition: blockNumber,
	})

	req, _ := http.NewRequestWithContext(ctx, http.MethodPatch, fmt.Sprintf("%v/nodes/%v/block-position", c.config.URL, nodeID), body)
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
		return errors.FromError(fmt.Errorf("could not update block position - got %d - body: %s", r.StatusCode, buf.String())).ExtendComponent(component)
	}

	return nil
}
