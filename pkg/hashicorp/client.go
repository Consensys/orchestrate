package hashicorp

import (
	"context"
	"fmt"
	"path"

	"github.com/ConsenSys/orchestrate/pkg/errors"
	"github.com/ConsenSys/orchestrate/pkg/toolkit/app/log"
	"github.com/hashicorp/vault/api"
)

const component = "hashicorp-vault.client"

const (
	signatureLabel = "signature"
	dataLabel      = "data"
	chainIDLabel   = "chainID"
	toLabel        = "to"
	gasPriceLabel  = "gasPrice"
	nonceLabel     = "nonce"
	gasLimitLabel  = "gasLimit"
	amountLabel    = "amount"

	zksAccountType = "zk-snarks"
	ethAccountType = "ethereum"
)

// OrchestrateVaultClient wraps a HashiCorp client and manage the unsealing
type OrchestrateVaultClient struct {
	client *api.Client
	config *Config
	logger *log.Logger
}

// NewOrchestrateVaultClient construct a new OrchestrateVaultClient
func NewOrchestrateVaultClient(ctx context.Context, config *Config) (*OrchestrateVaultClient, error) {
	logger := log.NewLogger().SetComponent(component)

	client, err := api.NewClient(ToVaultConfig(config))
	if err != nil {
		errMessage := "failed to instantiate Hashicorp Vault client"
		logger.WithError(err).Error(errMessage)
		return nil, errors.HashicorpVaultConnectionError("errMessage")
	}

	orchestrateVaultClient := &OrchestrateVaultClient{
		client: client,
		config: config,
		logger: logger,
	}

	tokenWatcher, err := newRenewTokenWatcher(client, config.TokenFilePath, logger)
	if err != nil {
		return nil, err
	}

	go func() {
		err = tokenWatcher.Run(ctx)
		if err != nil {
			logger.WithError(err).Fatal("token watcher routine has exited with errors")
		}
		logger.Warn("token watcher routine has exited gracefully")
	}()

	logger.Info("client has been initialized successfully")
	return orchestrateVaultClient, nil
}

func (c *OrchestrateVaultClient) HealthCheck() error {
	resp, err := c.client.Sys().Health()
	if err != nil {
		return parseErrorResponse(err)
	}

	if !resp.Initialized {
		errMessage := "client is not initialized"
		c.logger.WithError(err).Error(errMessage)
		return errors.HashicorpVaultConnectionError(errMessage)
	}

	return nil
}

func (c *OrchestrateVaultClient) listNamespaces(accountType string) ([]string, error) {
	res, err := c.client.Logical().List(path.Join(c.config.MountPoint, accountType, "/namespaces"))
	if err != nil {
		return []string{}, parseErrorResponse(err)
	}

	if res == nil || res.Data == nil {
		return []string{}, nil
	}

	secrets, ok := res.Data["keys"].([]interface{})
	if !ok {
		return []string{}, nil
	}

	rv := make([]string, len(secrets))
	for i, elem := range secrets {
		rv[i] = fmt.Sprintf("%v", elem)
	}

	return rv, nil
}

func (c *OrchestrateVaultClient) listAccounts(accountType, namespace string) ([]string, error) {
	c.client.SetNamespace(namespace)
	res, err := c.client.Logical().List(path.Join(c.config.MountPoint, accountType, "/accounts"))
	if err != nil {
		return []string{}, parseErrorResponse(err)
	}

	if res == nil || res.Data == nil {
		return []string{}, nil
	}

	secrets, ok := res.Data["keys"].([]interface{})
	if !ok {
		return []string{}, nil
	}

	rv := make([]string, len(secrets))
	for i, elem := range secrets {
		rv[i] = fmt.Sprintf("%v", elem)
	}

	return rv, nil
}

func (c *OrchestrateVaultClient) getAccount(accountType, accID, namespace string, account interface{}) error {
	c.client.SetNamespace(namespace)
	res, err := c.client.Logical().Read(path.Join(c.config.MountPoint, accountType, "/accounts", accID))
	if err != nil {
		return parseErrorResponse(err)
	}

	if res == nil || res.Data == nil {
		return nil
	}

	err = parseResponse(res.Data, account)
	if err != nil {
		return err
	}

	return nil
}

func (c *OrchestrateVaultClient) createAccount(accountType, namespace string, account interface{}) error {
	c.client.SetNamespace(namespace)
	res, err := c.client.Logical().Write(path.Join(c.config.MountPoint, accountType, "/accounts"), nil)
	if err != nil {
		return parseErrorResponse(err)
	}

	if res == nil || res.Data == nil {
		return nil
	}

	err = parseResponse(res.Data, account)
	if err != nil {
		return err
	}

	return nil
}
