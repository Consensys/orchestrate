package secretstore

import (
	"github.com/hashicorp/vault/api"
	"sync"
	"encoding/json"
	"fmt"
)

type credentials struct {
	Keys []string 			`json:"keys"`
	KeysBase64 []string 	`json:"keys_base_64"`
	Token string 			`json:"root_token"`

	retrieveSecretOnce *sync.Once // TODO
}


var once sync.Once

// Credentials is a singleton that contains the vault credentials.
// It ensures, that the secrets are not copied in the program
var Credentials credentials


func (c *credentials) FetchFromAWS(
	ss *AWS,
	name string,
) (err error) {

	once.Do(func() {
		secret, err := ss.Load(name)
		if err != nil {
			return
		}
		c.fromEncoded(secret)
	})

	return err
}

func (c *credentials) SendToAWS(
	ss *AWS,
	name string,
) (err error) {

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
func (c *credentials) FetchFromVaultInit(client *api.Client) (err error) {

	sys := client.Sys()

	initRequest := &api.InitRequest{
		SecretShares: 1,
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



func (c *credentials) fromEncoded(value string) (err error) {

	decoded := &credentials{}
	
	err = json.Unmarshal([]byte(value), decoded)
	if err != nil {
		return err
	}

	*c = *decoded
	return nil
}

func (c *credentials) encode() (encoded string, err error) {

	res, err := json.Marshal(*c)
	if err != nil {
		return "", err
	}

	encoded = (string)(res)
	return encoded, err

}

func (c *credentials) AttachTo(client *api.Client) {
	client.SetToken(c.Token)
}

func (c *credentials) Unseal(client *api.Client) (err error) {

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
		} else if status2.Sealed {
			return fmt.Errorf("Error, the vault was not properly unsealed")
		}
	}

	return nil

}

