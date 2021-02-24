package client

import (
	"context"
	"fmt"

	"github.com/ConsenSys/orchestrate/pkg/errors"
	clientutils "github.com/ConsenSys/orchestrate/pkg/http/client-utils"
	"github.com/ConsenSys/orchestrate/pkg/http/httputil"
	"github.com/ConsenSys/orchestrate/pkg/types/keymanager"
	types "github.com/ConsenSys/orchestrate/pkg/types/keymanager/zk-snarks"
)

const (
	zksAccountsPath = "zk-snarks/accounts"
)

func (c *HTTPClient) ZKSCreateAccount(ctx context.Context, req *types.CreateZKSAccountRequest) (*types.ZKSAccountResponse, error) {
	resp := &types.ZKSAccountResponse{}
	err := c.createAccount(ctx, zksAccountType, req, resp)
	return resp, err
}

func (c *HTTPClient) ZKSListAccounts(ctx context.Context, namespace string) ([]string, error) {
	return c.listAccounts(ctx, zksAccountType, namespace)
}

func (c *HTTPClient) ZKSListNamespaces(ctx context.Context) ([]string, error) {
	return c.listNamespaces(ctx, zksAccountType)
}

func (c *HTTPClient) ZKSGetAccount(ctx context.Context, address, namespace string) (*types.ZKSAccountResponse, error) {
	resp := &types.ZKSAccountResponse{}
	err := c.getAccount(ctx, zksAccountType, address, namespace, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *HTTPClient) ZKSSign(ctx context.Context, address string, req *keymanager.SignPayloadRequest) (string, error) {
	reqURL := fmt.Sprintf("%v/%s/%v/sign", c.config.URL, zksAccountsPath, address)

	response, err := clientutils.PostRequest(ctx, c.client, reqURL, req)
	if err != nil {
		return "", errors.ServiceConnectionError("error while signing data with zk-snarks account")
	}

	defer clientutils.CloseResponse(response)
	return httputil.ParseStringResponse(ctx, response)
}

func (c *HTTPClient) ZKSVerifySignature(ctx context.Context, request *types.VerifyPayloadRequest) error {
	reqURL := fmt.Sprintf("%v/%s/verify-signature", c.config.URL, zksAccountsPath)

	response, err := clientutils.PostRequest(ctx, c.client, reqURL, request)
	if err != nil {
		return errors.ServiceConnectionError("error while verifying zks signature")
	}

	defer clientutils.CloseResponse(response)
	return httputil.ParseEmptyBodyResponse(ctx, response)
}
