package client

import (
	"context"
	"fmt"

	"github.com/ConsenSys/orchestrate/pkg/errors"
	clientutils "github.com/ConsenSys/orchestrate/pkg/toolkit/app/http/client-utils"
	"github.com/ConsenSys/orchestrate/pkg/toolkit/app/http/httputil"
	"github.com/ConsenSys/orchestrate/pkg/types/keymanager"
	types "github.com/ConsenSys/orchestrate/pkg/types/keymanager/ethereum"
)

const (
	ethAccountsPath = "ethereum/accounts"
)

func (c *HTTPClient) ETHCreateAccount(ctx context.Context, req *types.CreateETHAccountRequest) (*types.ETHAccountResponse, error) {
	resp := &types.ETHAccountResponse{}
	err := c.createAccount(ctx, ethAccountType, req, resp)
	return resp, err
}

func (c *HTTPClient) ETHImportAccount(ctx context.Context, req *types.ImportETHAccountRequest) (*types.ETHAccountResponse, error) {
	reqURL := fmt.Sprintf("%v/%s/import", c.config.URL, ethAccountsPath)
	resp := &types.ETHAccountResponse{}

	response, err := clientutils.PostRequest(ctx, c.client, reqURL, req)
	if err != nil {
		return nil, errors.ServiceConnectionError("error while importing ethereum account")
	}

	defer clientutils.CloseResponse(response)
	if err := httputil.ParseResponse(ctx, response, resp); err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	return resp, nil
}

func (c *HTTPClient) ETHSign(ctx context.Context, address string, req *keymanager.SignPayloadRequest) (string, error) {
	reqURL := fmt.Sprintf("%v/%s/%v/sign", c.config.URL, ethAccountsPath, address)

	response, err := clientutils.PostRequest(ctx, c.client, reqURL, req)
	if err != nil {
		return "", errors.ServiceConnectionError("error while signing data with Ethereum account")
	}

	defer clientutils.CloseResponse(response)
	return httputil.ParseStringResponse(ctx, response)
}

func (c *HTTPClient) ETHSignTypedData(ctx context.Context, address string, req *types.SignTypedDataRequest) (string, error) {
	reqURL := fmt.Sprintf("%v/%s/%v/sign-typed-data", c.config.URL, ethAccountsPath, address)

	response, err := clientutils.PostRequest(ctx, c.client, reqURL, req)
	if err != nil {
		return "", errors.ServiceConnectionError("error while signing typed data with Ethereum account")
	}

	defer clientutils.CloseResponse(response)
	return httputil.ParseStringResponse(ctx, response)
}

func (c *HTTPClient) ETHSignTransaction(ctx context.Context, address string, req *types.SignETHTransactionRequest) (string, error) {
	reqURL := fmt.Sprintf("%v/%s/%v/sign-transaction", c.config.URL, ethAccountsPath, address)

	response, err := clientutils.PostRequest(ctx, c.client, reqURL, req)
	if err != nil {
		return "", errors.ServiceConnectionError("error while signing transaction")
	}

	defer clientutils.CloseResponse(response)
	return httputil.ParseStringResponse(ctx, response)
}

func (c *HTTPClient) ETHSignQuorumPrivateTransaction(ctx context.Context, address string, req *types.SignQuorumPrivateTransactionRequest) (string, error) {
	reqURL := fmt.Sprintf("%v/%s/%v/sign-quorum-private-transaction", c.config.URL, ethAccountsPath, address)

	response, err := clientutils.PostRequest(ctx, c.client, reqURL, req)
	if err != nil {
		return "", errors.ServiceConnectionError("error while signing quorum private transaction")
	}

	defer clientutils.CloseResponse(response)
	return httputil.ParseStringResponse(ctx, response)
}

func (c *HTTPClient) ETHSignEEATransaction(ctx context.Context, address string, req *types.SignEEATransactionRequest) (string, error) {
	reqURL := fmt.Sprintf("%v/%s/%v/sign-eea-transaction", c.config.URL, ethAccountsPath, address)

	response, err := clientutils.PostRequest(ctx, c.client, reqURL, req)
	if err != nil {
		return "", errors.ServiceConnectionError("error while signing eea private transaction")
	}

	defer clientutils.CloseResponse(response)
	return httputil.ParseStringResponse(ctx, response)
}

func (c *HTTPClient) ETHListAccounts(ctx context.Context, namespace string) ([]string, error) {
	return c.listAccounts(ctx, ethAccountType, namespace)
}

func (c *HTTPClient) ETHListNamespaces(ctx context.Context) ([]string, error) {
	return c.listNamespaces(ctx, ethAccountType)
}

func (c *HTTPClient) ETHGetAccount(ctx context.Context, address, namespace string) (*types.ETHAccountResponse, error) {
	resp := &types.ETHAccountResponse{}
	err := c.getAccount(ctx, ethAccountType, address, namespace, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *HTTPClient) ETHVerifySignature(ctx context.Context, request *types.VerifyPayloadRequest) error {
	reqURL := fmt.Sprintf("%v/%s/verify-signature", c.config.URL, ethAccountsPath)

	response, err := clientutils.PostRequest(ctx, c.client, reqURL, request)
	if err != nil {
		return errors.ServiceConnectionError("error while verifying ethereum signature")
	}

	defer clientutils.CloseResponse(response)
	return httputil.ParseEmptyBodyResponse(ctx, response)
}

func (c *HTTPClient) ETHVerifyTypedDataSignature(ctx context.Context, request *types.VerifyTypedDataRequest) error {
	reqURL := fmt.Sprintf("%v/%s/verify-typed-data-signature", c.config.URL, ethAccountsPath)

	response, err := clientutils.PostRequest(ctx, c.client, reqURL, request)
	if err != nil {
		return errors.ServiceConnectionError("error while verifying typed data signature").ExtendComponent(component)
	}

	defer clientutils.CloseResponse(response)
	return httputil.ParseEmptyBodyResponse(ctx, response)
}
