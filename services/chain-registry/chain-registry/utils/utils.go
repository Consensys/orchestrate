package utils

import (
	"context"
	"encoding/json"
	"html"
	"math/big"
	"net/http"
	"net/url"

	"github.com/containous/traefik/v2/pkg/log"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/ethereum/ethclient"
)

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

func VerifyURLs(ctx context.Context, ec ethclient.ChainSyncReader, uris []string) error {
	var prevChainID *big.Int

	for i, uri := range uris {
		chainID, err := ec.Network(ctx, uri)
		if err != nil {
			log.FromContext(ctx).WithError(err).Errorf("failed to fetch chain id for URL %s", uri)
			return err
		}

		if i > 0 && chainID != prevChainID {
			errMessage := "URLs in the list point to different networks"
			log.FromContext(ctx).Errorf(errMessage)
			return errors.InvalidParameterError(errMessage)
		}

		prevChainID = chainID
	}

	return nil
}
func GetChainID(ctx context.Context, ec ethclient.ChainSyncReader, uris []string) (*big.Int, error) {
	var chainID *big.Int
	var err error

	if len(uris) == 0 {
		errMessage := "invalid URLs list"
		log.FromContext(ctx).Errorf(errMessage)
		return nil, errors.InvalidParameterError(errMessage)
	}

	// All URLs must be valid and we return the head of the latest one
	for _, uri := range uris {
		chainID, err = ec.Network(ctx, uri)
		if err != nil {
			log.FromContext(ctx).WithError(err).Errorf("failed to fetch chain id for URL %s", uri)
			continue
		}

		return chainID, nil
	}

	errMessage := "All URLs in the list are unreachable"
	log.FromContext(ctx).WithError(err).Errorf(errMessage)
	return nil, errors.InvalidParameterError(errMessage)
}

func GetChainTip(ctx context.Context, ec ethclient.ChainLedgerReader, uris []string) uint64 {
	for _, uri := range uris {
		head, err := ec.HeaderByNumber(ctx, uri, nil)
		if err != nil {
			log.FromContext(ctx).WithError(err).Errorf("failed to fetch chain id for URL %s", uri)
			continue
		}

		return head.Number.Uint64()
	}

	return 0
}
