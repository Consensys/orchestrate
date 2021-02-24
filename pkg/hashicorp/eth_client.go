package hashicorp

import (
	"path"

	"github.com/ConsenSys/orchestrate/pkg/types/entities"
	types "github.com/ConsenSys/orchestrate/pkg/types/keymanager/ethereum"
)

func (c *OrchestrateVaultClient) ETHCreateAccount(namespace string) (*entities.ETHAccount, error) {
	account := &entities.ETHAccount{}
	err := c.createAccount(ethAccountType, namespace, account)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (c *OrchestrateVaultClient) ETHImportAccount(namespace, privateKey string) (*entities.ETHAccount, error) {
	c.client.SetNamespace(namespace)
	res, err := c.client.Logical().Write(path.Join(c.config.MountPoint, "ethereum/accounts/import"), map[string]interface{}{
		"privateKey": privateKey,
	})
	if err != nil {
		return nil, parseErrorResponse(err)
	}

	if res == nil || res.Data == nil {
		return nil, nil
	}

	ethAccount := &entities.ETHAccount{}
	err = parseResponse(res.Data, ethAccount)
	if err != nil {
		return nil, err
	}

	return ethAccount, nil
}

func (c *OrchestrateVaultClient) ETHSign(address, namespace, data string) (string, error) {
	c.client.SetNamespace(namespace)
	res, err := c.client.Logical().Write(path.Join(c.config.MountPoint, "ethereum/accounts", address, "sign"), map[string]interface{}{
		dataLabel: data,
	})
	if err != nil {
		return "", parseErrorResponse(err)
	}

	if res == nil || res.Data == nil {
		return "", nil
	}

	return res.Data[signatureLabel].(string), nil
}

func (c *OrchestrateVaultClient) ETHSignTransaction(address string, request *types.SignETHTransactionRequest) (string, error) {
	c.client.SetNamespace(request.Namespace)
	res, err := c.client.Logical().Write(path.Join(c.config.MountPoint, "ethereum/accounts", address, "sign-transaction"), map[string]interface{}{
		dataLabel:     request.Data,
		chainIDLabel:  request.ChainID,
		toLabel:       request.To,
		gasPriceLabel: request.GasPrice,
		nonceLabel:    request.Nonce,
		gasLimitLabel: request.GasLimit,
		amountLabel:   request.Amount,
	})
	if err != nil {
		return "", parseErrorResponse(err)
	}

	if res == nil || res.Data == nil {
		return "", nil
	}

	return res.Data[signatureLabel].(string), nil
}

func (c *OrchestrateVaultClient) ETHSignQuorumPrivateTransaction(address string, request *types.SignQuorumPrivateTransactionRequest) (string, error) {
	c.client.SetNamespace(request.Namespace)
	res, err := c.client.Logical().Write(path.Join(c.config.MountPoint, "ethereum/accounts", address, "sign-quorum-private-transaction"), map[string]interface{}{
		dataLabel:     request.Data,
		toLabel:       request.To,
		gasPriceLabel: request.GasPrice,
		nonceLabel:    request.Nonce,
		gasLimitLabel: request.GasLimit,
		amountLabel:   request.Amount,
	})
	if err != nil {
		return "", parseErrorResponse(err)
	}

	if res == nil || res.Data == nil {
		return "", nil
	}

	return res.Data[signatureLabel].(string), nil
}

func (c *OrchestrateVaultClient) ETHSignEEATransaction(address string, request *types.SignEEATransactionRequest) (string, error) {
	c.client.SetNamespace(request.Namespace)
	res, err := c.client.Logical().Write(path.Join(c.config.MountPoint, "ethereum/accounts", address, "sign-eea-transaction"), map[string]interface{}{
		dataLabel:        request.Data,
		toLabel:          request.To,
		chainIDLabel:     request.ChainID,
		nonceLabel:       request.Nonce,
		"privateFrom":    request.PrivateFrom,
		"privateFor":     request.PrivateFor,
		"privacyGroupID": request.PrivacyGroupID,
	})
	if err != nil {
		return "", parseErrorResponse(err)
	}

	if res == nil || res.Data == nil {
		return "", nil
	}

	return res.Data[signatureLabel].(string), nil
}

func (c *OrchestrateVaultClient) ETHListAccounts(namespace string) ([]string, error) {
	return c.listAccounts(ethAccountType, namespace)
}

func (c *OrchestrateVaultClient) ETHListNamespaces() ([]string, error) {
	return c.listNamespaces(ethAccountType)
}

func (c *OrchestrateVaultClient) ETHGetAccount(address, namespace string) (*entities.ETHAccount, error) {
	account := &entities.ETHAccount{}
	err := c.getAccount(ethAccountType, address, namespace, account)
	if err != nil {
		return nil, err
	}

	return account, nil
}
