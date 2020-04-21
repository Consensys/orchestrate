package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"net/http"
	"net/url"

	"github.com/containous/traefik/v2/pkg/log"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/ethereum/ethclient"
	ethclientutils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/ethereum/ethclient/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"

	"github.com/go-playground/validator/v10"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
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

func UnmarshalBody(body io.Reader, req interface{}) error {
	dec := json.NewDecoder(body)
	dec.DisallowUnknownFields() // Force errors if unknown fields
	err := dec.Decode(req)
	if err != nil {
		return errors.FromError(err).ExtendComponent(component)
	}

	err = utils.GetValidator().Struct(req)
	if err != nil {
		if ves, ok := err.(validator.ValidationErrors); ok {
			var errorMessage string
			for _, fe := range ves {
				errorMessage += fmt.Sprintf(" field validation for '%s' failed on the '%s' tag", fe.Field(), fe.Tag())
			}
			return errors.InvalidParameterError("invalid body, with:%s", errorMessage).ExtendComponent(component)
		}
		return errors.FromError(fmt.Errorf("invalid body")).ExtendComponent(component)
	}

	return nil
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
