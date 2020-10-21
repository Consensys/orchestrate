package httputil

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/containous/traefik/v2/pkg/log"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
)

const (
	invalidResponseBody = "failed to decode response body"
)

func ParseResponse(ctx context.Context, response *http.Response, resp interface{}) error {
	if response.StatusCode == http.StatusAccepted || response.StatusCode == http.StatusOK {
		if resp == nil {
			return nil
		}

		if err := json.NewDecoder(response.Body).Decode(resp); err != nil {
			log.FromContext(ctx).WithError(err).Error(invalidResponseBody)
			return errors.ServiceConnectionError(invalidResponseBody)
		}

		return nil
	}

	errResp := ErrorResponse{}
	if err := json.NewDecoder(response.Body).Decode(&errResp); err != nil {
		log.FromContext(ctx).WithError(err).Error(invalidResponseBody)
		return errors.ServiceConnectionError(invalidResponseBody)
	}

	log.FromContext(ctx).Error(errResp.Message)
	return errors.Errorf(errResp.Code, errResp.Message)
}

func ParseStringResponse(ctx context.Context, response *http.Response) (string, error) {
	if response.StatusCode != http.StatusOK {
		errResp := ErrorResponse{}
		if err := json.NewDecoder(response.Body).Decode(&errResp); err != nil {
			errMessage := "failed to decode error response body"
			log.FromContext(ctx).WithError(err).Error(errMessage)
			return "", errors.ServiceConnectionError(errMessage)
		}

		return "", errors.Errorf(errResp.Code, errResp.Message)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		errMessage := "failed to decode response body"
		log.FromContext(ctx).WithError(err).Error(errMessage)
		return "", errors.ServiceConnectionError(errMessage)
	}

	return string(responseData), nil
}
