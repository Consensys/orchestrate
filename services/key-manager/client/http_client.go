package client

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/containous/traefik/v2/pkg/log"
	healthz "github.com/heptiolabs/healthcheck"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	clientutils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/http/client-utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/http/httputil"
)

const zksAccountType = "zk-snarks"
const ethAccountType = "ethereum"

func NewHTTPClient(h *http.Client, c *Config) KeyManagerClient {
	return &HTTPClient{
		client: h,
		config: c,
	}
}

type HTTPClient struct {
	client *http.Client
	config *Config
}

func (c *HTTPClient) Checker() healthz.Check {
	return healthz.HTTPGetCheck(fmt.Sprintf("%s/live", c.config.MetricsURL), time.Second)
}

func (c *HTTPClient) listNamespaces(ctx context.Context, accountType string) ([]string, error) {
	reqURL := fmt.Sprintf("%v/%s/namespaces", c.config.URL, accountType)

	response, err := clientutils.GetRequest(ctx, c.client, reqURL)
	if err != nil {
		errMessage := "error while listing namespaces"
		log.FromContext(ctx).WithError(err).Error(errMessage)
		return []string{}, errors.ServiceConnectionError(errMessage).ExtendComponent(component)
	}

	defer clientutils.CloseResponse(response)
	var resp []string
	err = httputil.ParseResponse(ctx, response, &resp)
	return resp, err
}

func (c *HTTPClient) listAccounts(ctx context.Context, accountType, namespace string) ([]string, error) {
	reqURL := fmt.Sprintf("%v/%s/accounts", c.config.URL, accountType)
	if namespace != "" {
		reqURL += fmt.Sprintf("?namespace=%s", namespace)
	}

	response, err := clientutils.GetRequest(ctx, c.client, reqURL)
	if err != nil {
		errMessage := "error while listing accounts"
		log.FromContext(ctx).WithError(err).Error(errMessage)
		return []string{}, errors.ServiceConnectionError(errMessage).ExtendComponent(component)
	}

	defer clientutils.CloseResponse(response)
	var resp []string
	err = httputil.ParseResponse(ctx, response, &resp)
	return resp, err
}

func (c *HTTPClient) getAccount(ctx context.Context, accountType, address, namespace string, resp interface{}) error {
	reqURL := fmt.Sprintf("%v/%s/accounts/%s", c.config.URL, accountType, address)
	if namespace != "" {
		reqURL += fmt.Sprintf("?namespace=%s", namespace)
	}

	response, err := clientutils.GetRequest(ctx, c.client, reqURL)
	if err != nil {
		errMessage := "error while getting account"
		log.FromContext(ctx).WithError(err).Error(errMessage)
		return errors.ServiceConnectionError(errMessage).ExtendComponent(component)
	}

	defer clientutils.CloseResponse(response)
	if err := httputil.ParseResponse(ctx, response, resp); err != nil {
		return errors.FromError(err).ExtendComponent(component)
	}

	return nil
}

func (c *HTTPClient) createAccount(ctx context.Context, accountType string, req, resp interface{}) error {
	reqURL := fmt.Sprintf("%v/%s/accounts", c.config.URL, accountType)
	response, err := clientutils.PostRequest(ctx, c.client, reqURL, req)
	if err != nil {
		errMessage := "error while creating account"
		log.FromContext(ctx).WithError(err).Error(errMessage)
		return errors.ServiceConnectionError(errMessage).ExtendComponent(component)
	}

	defer clientutils.CloseResponse(response)
	if err := httputil.ParseResponse(ctx, response, resp); err != nil {
		return errors.FromError(err).ExtendComponent(component)
	}

	return nil
}
