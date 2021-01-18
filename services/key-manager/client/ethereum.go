package client

import (
	"context"
	"fmt"

	"github.com/containous/traefik/v2/pkg/log"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	clientutils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/http/client-utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/http/httputil"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/keymanager"
	types "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/keymanager/ethereum"
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
		errMessage := "error while importing ethereum account"
		log.FromContext(ctx).WithError(err).Error(errMessage)
		return nil, errors.ServiceConnectionError(errMessage).ExtendComponent(component)
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
		errMessage := "error while signing data with Ethereum account"
		log.FromContext(ctx).WithError(err).Error(errMessage)
		return "", errors.ServiceConnectionError(errMessage).ExtendComponent(component)
	}

	defer clientutils.CloseResponse(response)
	return httputil.ParseStringResponse(ctx, response)
}

func (c *HTTPClient) ETHSignTypedData(ctx context.Context, address string, req *types.SignTypedDataRequest) (string, error) {
	reqURL := fmt.Sprintf("%v/%s/%v/sign-typed-data", c.config.URL, ethAccountsPath, address)

	response, err := clientutils.PostRequest(ctx, c.client, reqURL, req)
	if err != nil {
		errMessage := "error while signing typed data with Ethereum account"
		log.FromContext(ctx).WithError(err).Error(errMessage)
		return "", errors.ServiceConnectionError(errMessage).ExtendComponent(component)
	}

	defer clientutils.CloseResponse(response)
	return httputil.ParseStringResponse(ctx, response)
}

func (c *HTTPClient) ETHSignTransaction(ctx context.Context, address string, req *types.SignETHTransactionRequest) (string, error) {
	reqURL := fmt.Sprintf("%v/%s/%v/sign-transaction", c.config.URL, ethAccountsPath, address)

	response, err := clientutils.PostRequest(ctx, c.client, reqURL, req)
	if err != nil {
		errMessage := "error while signing transaction"
		log.FromContext(ctx).WithError(err).Error(errMessage)
		return "", errors.ServiceConnectionError(errMessage).ExtendComponent(component)
	}

	defer clientutils.CloseResponse(response)
	return httputil.ParseStringResponse(ctx, response)
}

func (c *HTTPClient) ETHSignQuorumPrivateTransaction(ctx context.Context, address string, req *types.SignQuorumPrivateTransactionRequest) (string, error) {
	reqURL := fmt.Sprintf("%v/%s/%v/sign-quorum-private-transaction", c.config.URL, ethAccountsPath, address)

	response, err := clientutils.PostRequest(ctx, c.client, reqURL, req)
	if err != nil {
		errMessage := "error while signing quorum private transaction"
		log.FromContext(ctx).WithError(err).Error(errMessage)
		return "", errors.ServiceConnectionError(errMessage).ExtendComponent(component)
	}

	defer clientutils.CloseResponse(response)
	return httputil.ParseStringResponse(ctx, response)
}

func (c *HTTPClient) ETHSignEEATransaction(ctx context.Context, address string, req *types.SignEEATransactionRequest) (string, error) {
	reqURL := fmt.Sprintf("%v/%s/%v/sign-eea-transaction", c.config.URL, ethAccountsPath, address)

	response, err := clientutils.PostRequest(ctx, c.client, reqURL, req)
	if err != nil {
		errMessage := "error while signing eea private transaction"
		log.FromContext(ctx).WithError(err).Error(errMessage)
		return "", errors.ServiceConnectionError(errMessage).ExtendComponent(component)
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
		errMessage := "error while verifying ethereum signature"
		log.FromContext(ctx).WithError(err).Error(errMessage)
		return errors.ServiceConnectionError(errMessage).ExtendComponent(component)
	}

	defer clientutils.CloseResponse(response)
	return httputil.ParseEmptyBodyResponse(ctx, response)
}

func (c *HTTPClient) ETHVerifyTypedDataSignature(ctx context.Context, request *types.VerifyTypedDataRequest) error {
	reqURL := fmt.Sprintf("%v/%s/verify-typed-data-signature", c.config.URL, ethAccountsPath)

	response, err := clientutils.PostRequest(ctx, c.client, reqURL, request)
	if err != nil {
		errMessage := "error while verifying typed data signature"
		log.FromContext(ctx).WithError(err).Error(errMessage)
		return errors.ServiceConnectionError(errMessage).ExtendComponent(component)
	}

	defer clientutils.CloseResponse(response)
	return httputil.ParseEmptyBodyResponse(ctx, response)
}
