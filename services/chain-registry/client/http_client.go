package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/api"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/types"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
)

type Config struct {
	URL string
}

type HTTPClient struct {
	client http.Client

	config Config
}

func NewHTTPClient(h http.Client, c Config) *HTTPClient {
	return &HTTPClient{
		client: h,
		config: c,
	}
}

func (c *HTTPClient) GetNodeByID(nodeID string) (*types.Node, error) {
	url := fmt.Sprintf("%s%s/%s", c.config.URL, api.NodePrefixPath, nodeID)
	r, err := c.client.Get(url)
	if err != nil {
		return nil, errors.FromError(fmt.Errorf("%v - url: %s", err, url)).ExtendComponent(component)
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

func (c *HTTPClient) GetNodes() ([]*types.Node, error) {
	url := fmt.Sprintf("%s%s", c.config.URL, api.NodesPrefixPath)
	r, err := c.client.Get(url)
	if err != nil {
		return nil, errors.FromError(fmt.Errorf("%v - url: %s", err, url)).ExtendComponent(component)
	}
	defer func() {
		if deferErr := r.Body.Close(); err != nil {
			log.WithError(deferErr).Errorf("%s: could close body", component)
		}
	}()

	if r.StatusCode != http.StatusOK {
		return nil, errors.FromError(fmt.Errorf("could not get nodes - got %d - url %s", r.StatusCode, url)).ExtendComponent(component)
	}

	var nodes []*types.Node
	if err := json.NewDecoder(r.Body).Decode(&nodes); err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	return nodes, nil
}

func (c *HTTPClient) UpdateBlockPosition(nodeID string, blockNumber int64) error {
	body := new(bytes.Buffer)
	_ = json.NewEncoder(body).Encode(&api.PatchBlockPositionRequest{
		BlockPosition: blockNumber,
	})

	url := fmt.Sprintf("%s%s/%s/%s", c.config.URL, api.NodePrefixPath, nodeID, api.BlockPositionPath)
	req, _ := http.NewRequest(http.MethodPatch, url, body)
	r, err := c.client.Do(req)
	if err != nil {
		return errors.FromError(fmt.Errorf("%v - url: %s", err, url)).ExtendComponent(component)
	}
	defer func() {
		if deferErr := r.Body.Close(); err != nil {
			log.WithError(deferErr).Errorf("%s: could close body", component)
		}
	}()

	if r.StatusCode != http.StatusOK {
		buf := new(bytes.Buffer)
		_, _ = buf.ReadFrom(r.Body)
		return errors.FromError(fmt.Errorf("could not update block position - got %d - url: %s - body: %s", r.StatusCode, url, buf.String())).ExtendComponent(component)
	}

	return nil
}
