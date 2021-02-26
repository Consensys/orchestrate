package hashicorp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

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
func NewOrchestrateVaultClient(config *Config) (*OrchestrateVaultClient, error) {
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

	err = orchestrateVaultClient.setTokenFromConfig(config)
	if err != nil {
		return nil, err
	}

	err = orchestrateVaultClient.manageToken()
	if err != nil {
		return nil, err
	}

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

func (c *OrchestrateVaultClient) manageToken() error {
	secret, err := c.client.Auth().Token().LookupSelf()
	if err != nil {
		errMessage := "initial token lookup failed"
		c.logger.WithError(err).Error(errMessage)
		return errors.HashicorpVaultConnectionError(errMessage)
	}

	tokenTTL64, err := secret.Data["creation_ttl"].(json.Number).Int64()
	if err != nil {
		errMessage := "failed to get 'creation_ttl' field"
		c.logger.WithError(err).Error(errMessage)
		return errors.HashicorpVaultConnectionError(errMessage)
	}

	tokenRenewable := secret.Data["renewable"].(bool)
	if int(tokenTTL64) == 0 || !tokenRenewable {
		c.logger.Debug("token in use never expires or cannot be renewed")
		return nil
	}

	tokenExpireIn64, err := secret.Data["ttl"].(json.Number).Int64()
	if err != nil {
		return errors.InternalError("could not read vault ttl").AppendReason(err.Error())
	}
	if int(tokenExpireIn64) == 0 {
		return errors.InternalError("token is expired")
	}

	c.logger.WithField("expiration_duration", tokenExpireIn64).Debug("token expiration time")

	rtl := newRenewTokenLoop(tokenExpireIn64, c.client, c.logger)

	err = rtl.Refresh()
	if err != nil {
		return err
	}

	rtl.Run()
	return nil
}

func (c *OrchestrateVaultClient) setTokenFromConfig(config *Config) error {
	encoded, err := ioutil.ReadFile(config.TokenFilePath)
	if err != nil {
		errMessage := "token file path could not be found"
		c.logger.WithError(err).Fatal(errMessage)
		return errors.ConfigError(errMessage)
	}

	decoded := strings.TrimSuffix(string(encoded), "\n") // Remove the newline if it exists
	decoded = strings.TrimSuffix(decoded, "\r")          // This one is for windows compatibility
	c.client.SetToken(decoded)

	// Immediately delete the file after it was read
	err = os.Remove(config.TokenFilePath)
	if err != nil {
		c.logger.WithError(err).Warn("could not delete token file")
	}

	return nil
}
