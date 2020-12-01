package hashicorp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/hashicorp/vault/api"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
	types "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/keymanager/ethereum"
)

const (
	signatureLabel = "signature"
	dataLabel      = "data"
	chainIDLabel   = "chainID"
	toLabel        = "to"
	gasPriceLabel  = "gasPrice"
	nonceLabel     = "nonce"
	gasLimitLabel  = "gasLimit"
	amountLabel    = "amount"
)

// OrchestrateVaultClient wraps a HashiCorp client and manage the unsealing
type OrchestrateVaultClient struct {
	client *api.Client
	config *Config
}

// NewOrchestrateVaultClient construct a new OrchestrateVaultClient
func NewOrchestrateVaultClient(config *Config) (*OrchestrateVaultClient, error) {
	client, err := api.NewClient(ToVaultConfig(config))
	if err != nil {
		errMessage := "failed to instantiate Hashicorp Vault client"
		log.WithError(err).Error(errMessage)
		return nil, errors.HashicorpVaultConnectionError("errMessage")
	}

	orchestrateVaultClient := &OrchestrateVaultClient{
		client: client,
		config: config,
	}

	err = orchestrateVaultClient.setTokenFromConfig(config)
	if err != nil {
		return nil, err
	}

	err = orchestrateVaultClient.manageToken()
	if err != nil {
		return nil, err
	}

	log.Info("Hashicorp vault initialized")
	return orchestrateVaultClient, nil
}

func (c *OrchestrateVaultClient) ETHCreateAccount(namespace string) (*entities.ETHAccount, error) {
	log.WithField("token", c.client.Token()).Info("Token before HTTP call")

	c.client.SetNamespace(namespace)
	res, err := c.client.Logical().Write(path.Join(c.config.MountPoint, "ethereum/accounts"), nil)
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
	c.client.SetNamespace(namespace)
	res, err := c.client.Logical().List(path.Join(c.config.MountPoint, "ethereum/accounts"))
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

func (c *OrchestrateVaultClient) ETHListNamespaces() ([]string, error) {
	res, err := c.client.Logical().List(path.Join(c.config.MountPoint, "ethereum/namespaces"))
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

func (c *OrchestrateVaultClient) ETHGetAccount(address, namespace string) (*entities.ETHAccount, error) {
	c.client.SetNamespace(namespace)
	res, err := c.client.Logical().Read(path.Join(c.config.MountPoint, "ethereum/accounts", address))
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

func (c *OrchestrateVaultClient) HealthCheck() error {
	resp, err := c.client.Sys().Health()
	if err != nil {
		return parseErrorResponse(err)
	}

	if !resp.Initialized {
		errMessage := "Hashicorp Vault service is not initialized"
		log.WithError(err).Error(errMessage)
		return errors.HashicorpVaultConnectionError(errMessage)
	}

	return nil
}

func parseResponse(data map[string]interface{}, resp interface{}) error {
	jsonbody, err := json.Marshal(data)
	if err != nil {
		errMessage := "failed to marshal response data"
		log.WithError(err).Error(errMessage)
		return errors.EncodingError(errMessage)
	}

	if err := json.Unmarshal(jsonbody, &resp); err != nil {
		errMessage := "failed to unmarshal response data"
		log.WithError(err).Error(errMessage)
		return errors.EncodingError(errMessage)
	}

	return nil
}

func parseErrorResponse(err error) error {
	httpError, ok := err.(*api.ResponseError)
	if !ok {
		errMessage := "failed to parse error response"
		log.WithError(err).Error(errMessage)
		return errors.HashicorpVaultConnectionError(errMessage)
	}

	switch httpError.StatusCode {
	case http.StatusNotFound:
		return errors.NotFoundError(httpError.Error())
	case http.StatusBadRequest:
		return errors.InvalidFormatError(httpError.Error())
	case http.StatusUnprocessableEntity:
		return errors.InvalidParameterError(httpError.Error())
	case http.StatusConflict:
		return errors.AlreadyExistsError(httpError.Error())
	default:
		return errors.HashicorpVaultConnectionError(httpError.Error())
	}
}

func (c *OrchestrateVaultClient) manageToken() error {
	secret, err := c.client.Auth().Token().LookupSelf()
	if err != nil {
		errMessage := "initial token lookup failed"
		log.WithError(err).Error(errMessage)
		return errors.HashicorpVaultConnectionError(errMessage)
	}
	log.WithField("token_data", secret.Data).Debug("token data")

	tokenTTL64, err := secret.Data["creation_ttl"].(json.Number).Int64()
	if err != nil {
		errMessage := "failed to get creation_ttl field from token data"
		log.WithError(err).Error(errMessage)
		return errors.HashicorpVaultConnectionError(errMessage)
	}

	if int(tokenTTL64) == 0 {
		log.Debug("root token never expires")
		return nil
	}

	tokenExpireIn64, err := secret.Data["ttl"].(json.Number).Int64()
	if err != nil {
		return errors.InternalError("HashiCorp: Could not read vault ttl: %v", err)
	}
	log.WithField("expiration_duration", tokenExpireIn64).Debug("vault token expiration duration")

	rtl := newRenewTokenLoop(tokenExpireIn64, c.client)

	err = rtl.Refresh()
	if err != nil {
		return err
	}
	log.Info("initial token refresh succeeded")

	rtl.Run()
	return nil
}

func (c *OrchestrateVaultClient) setTokenFromConfig(config *Config) error {
	encoded, err := ioutil.ReadFile(config.TokenFilePath)
	if err != nil {
		errMessage := "token file path could not be found"
		log.WithError(err).Fatal(errMessage)
		return errors.ConfigError(errMessage)
	}

	decoded := strings.TrimSuffix(string(encoded), "\n") // Remove the newline if it exists
	decoded = strings.TrimSuffix(decoded, "\r")          // This one is for windows compatibility

	log.WithField("token", decoded).Info("Token first from file")
	c.client.SetToken(decoded)

	// Immediately delete the file after it was read
	err = os.Remove(config.TokenFilePath)
	if err != nil {
		log.WithError(err).Warn("could not delete token file")
	}

	return nil
}
