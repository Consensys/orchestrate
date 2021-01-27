package httputil

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/containous/traefik/v2/pkg/log"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
)

const (
	cannotReadResponseBody = "failed to read response body"
	invalidResponseBody    = "failed to decode response body"
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

	// Read body
	respMsg, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.FromContext(ctx).WithError(err).Error(cannotReadResponseBody)
		return errors.ServiceConnectionError(cannotReadResponseBody)
	}

	if string(respMsg) != "" {
		errResp := ErrorResponse{}
		if err = json.Unmarshal(respMsg, &errResp); err == nil {
			return errors.Errorf(errResp.Code, errResp.Message)
		}
	}

	return parseResponseError(response.StatusCode, string(respMsg))
}

func parseResponseError(statusCode int, errMsg string) error {
	switch statusCode {
	case http.StatusBadRequest:
		if errMsg == "" {
			errMsg = "invalid request data"
		}
		return errors.InvalidFormatError(errMsg)
	case http.StatusConflict:
		if errMsg == "" {
			errMsg = "invalid data message"
		}
		return errors.StorageError(errMsg)
	case http.StatusNotFound:
		if errMsg == "" {
			errMsg = "cannot find entity"
		}
		return errors.NotFoundError(errMsg)
	case http.StatusUnauthorized:
		if errMsg == "" {
			errMsg = "not authorized"
		}
		return errors.UnauthorizedError(errMsg)
	case http.StatusUnprocessableEntity:
		if errMsg == "" {
			errMsg = "invalid request format"
		}
		return errors.InvalidParameterError(errMsg)
	default:
		if errMsg == "" {
			errMsg = "server error"
		}
		return errors.ServiceConnectionError(errMsg)
	}
}

func ParseStringResponse(ctx context.Context, response *http.Response) (string, error) {
	if response.StatusCode != http.StatusOK {
		errResp := ErrorResponse{}
		if err := json.NewDecoder(response.Body).Decode(&errResp); err != nil {
			log.FromContext(ctx).WithError(err).Error(invalidResponseBody)
			return "", errors.ServiceConnectionError(invalidResponseBody)
		}

		return "", errors.Errorf(errResp.Code, errResp.Message)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.FromContext(ctx).WithError(err).Error(invalidResponseBody)
		return "", errors.ServiceConnectionError(invalidResponseBody)
	}

	return string(responseData), nil
}

func ParseEmptyBodyResponse(ctx context.Context, response *http.Response) error {
	if response.StatusCode != http.StatusNoContent && response.StatusCode != http.StatusAccepted {
		errResp := ErrorResponse{}
		if err := json.NewDecoder(response.Body).Decode(&errResp); err != nil {
			log.FromContext(ctx).WithError(err).Error(invalidResponseBody)
			return errors.ServiceConnectionError(invalidResponseBody)
		}

		return errors.Errorf(errResp.Code, errResp.Message)
	}

	return nil
}
