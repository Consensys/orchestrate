package hashicorp

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/hashicorp/vault/api"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/key-store.git/secretstore/aws"
)

type credentials struct {
	Keys       []string `json:"keys"`
	KeysBase64 []string `json:"keys_base_64"`
	Token      string   `json:"root_token"`

	fetchOnce *sync.Once // TODO
}

// Credentials is a singleton that contains the vault credentials.
// It ensures, that the secrets are not copied in the program
var Credentials credentials

func (c *credentials) FetchFromAWS(ss *aws.AWS, name string) error {
	secret, err := ss.Load(name)
	if err != nil {
		return err
	}
	c.fromEncoded(secret)
	return nil
}

func (c *credentials) SendToAWS(ss *aws.AWS, name string) error {
	encoded, err := c.encode()
	if err != nil {
		return err
	}

	err = ss.Store(name, encoded)
	if err != nil {
		return err
	}

	return nil
}

// FetchFromVaultInit runs a basic vault initialization and creates a credential object from there
func (c *credentials) FetchFromVaultInit(client *api.Client) error {
	sys := client.Sys()
	initRequest := &api.InitRequest{
		SecretShares:    1,
		SecretThreshold: 1,
	}

	initResponse, err := sys.Init(initRequest)
	if err != nil {
		return err
	}

	c.Keys = initResponse.Keys
	c.KeysBase64 = initResponse.KeysB64
	c.Token = initResponse.RootToken

	return nil
}

func (c *credentials) fromEncoded(value string) error {
	decoded := &credentials{}

	err := json.Unmarshal([]byte(value), decoded)
	if err != nil {
		return err
	}

	*c = *decoded
	return nil
}

func (c *credentials) encode() (string, error) {
	res, err := json.Marshal(*c)
	if err != nil {
		return "", err
	}
	return (string)(res), err
}

func (c *credentials) AttachTo(client *api.Client) {
	client.SetToken(c.Token)
}

func (c *credentials) Unseal(client *api.Client) error {
	sys := client.Sys()
	status, err := sys.SealStatus()
	if err != nil {
		return err
	}

	// Unseal is idemnpotent so no need to solve race conditions here
	if status.Sealed {
		if len(c.Keys) == 0 {
			return fmt.Errorf("The UnsealKey has not been imported. \n If you are running in dev, Consider verifying if the credentials have been correctly passed")
		}

		status2, err := sys.Unseal(c.Keys[0])
		if err != nil {
			return err
		}

		if status2.Sealed {
			return fmt.Errorf("Error, the vault was not properly unsealed")
		}
	}

	return nil

}
