package utils

import (
	"context"
	"encoding/json"
	"html"
	"net/http"
	"net/url"

	"github.com/containous/traefik/v2/pkg/log"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/ethereum/ethclient"
	ethclientutils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/ethereum/ethclient/utils"
)

const component = "chain-registry.store.api"

type apiError struct {
	Message string `json:"message"`
}

func ToFilters(values url.Values) map[string]string {
	filters := make(map[string]string)
	for key := range values {
		k := html.EscapeString(key)
		v := html.EscapeString(values.Get(key))
		if k != "" && v != "" {
			filters[k] = v
		}
	}
	return filters
}

func HandleStoreError(rw http.ResponseWriter, err error) {
	switch {
	case errors.IsAlreadyExistsError(err):
		WriteError(rw, err.Error(), http.StatusConflict)
	case errors.IsNotFoundError(err):
		WriteError(rw, err.Error(), http.StatusNotFound)
	case errors.IsDataError(err):
		WriteError(rw, err.Error(), http.StatusBadRequest)
	case err != nil:
		WriteError(rw, err.Error(), http.StatusInternalServerError)
	}
}

func WriteError(rw http.ResponseWriter, msg string, code int) {
	data, _ := json.Marshal(apiError{Message: msg})
	http.Error(rw, string(data), code)
}

func GetChainTip(ctx context.Context, ec ethclient.ChainLedgerReader, urls []string) (uint64, error) {
	var tip uint64

	// All URLs must be valid and we return the head of the latest one
	for _, uri := range urls {
		head, err := ec.HeaderByNumber(ethclientutils.RetryNotFoundError(ctx, true), uri, nil)
		if err != nil {
			log.FromContext(ctx).WithError(err).Errorf("failed to fetch chain tip for URL %s", uri)
			return 0, err
		}

		tip = head.Number.Uint64()
	}

	return tip, nil
}
