package hashicorp

import (
	"path"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
)

func (c *OrchestrateVaultClient) ZKSListNamespaces() ([]string, error) {
	return c.listNamespaces(zksAccountType)
}

func (c *OrchestrateVaultClient) ZKSGetAccount(address, namespace string) (*entities.ZKSAccount, error) {
	account := &entities.ZKSAccount{}
	err := c.getAccount(zksAccountType, address, namespace, account)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (c *OrchestrateVaultClient) ZKSListAccounts(namespace string) ([]string, error) {
	return c.listAccounts(zksAccountType, namespace)
}

func (c *OrchestrateVaultClient) ZKSCreateAccount(namespace string) (*entities.ZKSAccount, error) {
	account := &entities.ZKSAccount{}
	err := c.createAccount(zksAccountType, namespace, account)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (c *OrchestrateVaultClient) ZKSSign(address, namespace, data string) (string, error) {
	c.client.SetNamespace(namespace)
	res, err := c.client.Logical().Write(path.Join(c.config.MountPoint, "zk-snarks/accounts", address, "sign"), map[string]interface{}{
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
