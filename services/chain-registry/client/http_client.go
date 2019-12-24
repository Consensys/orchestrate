package client

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/types"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
)

type Config struct {
	url string
}

type HTTPClient struct {
	client http.Client

	config Config
}

func (c *HTTPClient) GetNodeByID(id string) (*types.Node, error) {
	r, err := http.Get(fmt.Sprintf("%s%s%s", c.config.url, "/api/nodes/", id))
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}
	defer func() {
		if deferErr := r.Body.Close(); err != nil {
			log.WithError(deferErr).Errorf("%s: could close body", component)
		}
	}()

	var node *types.Node
	if err := json.NewDecoder(r.Body).Decode(node); err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	return node, nil
}

func (c *HTTPClient) GetNodes() ([]*types.Node, error) {
	r, err := http.Get(fmt.Sprintf("%s%s", c.config.url, "/api/nodes"))
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}
	defer func() {
		if deferErr := r.Body.Close(); err != nil {
			log.WithError(deferErr).Errorf("%s: could close body", component)
		}
	}()

	var nodes []*types.Node
	if err := json.NewDecoder(r.Body).Decode(&nodes); err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	return nodes, nil
}
