package hashicorps

import (
	"github.com/hashicorp/vault/api"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	aws "gitlab.com/ConsenSys/client/fr/core-stack/infra/aws-secret-manager.git/aws"
	"sync"
	"encoding/json"
	"fmt"
)

type credentials struct {
	Keys []string 			`json:"keys"`
	KeysBase64 []string 	`json:"keys_base_64"`
	Token string 			`json:"root_token"`
}


var once sync.Once

// Credentials is a singleton that contains the vault credentials.
// It ensures, that the secrets are not copied in the program
var Credentials credentials


func (c *credentials) FetchFromAWS(
	client *secretsmanager.SecretsManager,
	name string,
) (err error) {

	once.Do(func() {
		secret, err := aws.SecretFromKey(name).SetClient(client).GetValue()
		if err != nil {
			return
		}
		c.fromEncoded(secret)
	})

	return err
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

func (c *credentials) AttachTo(client *api.Client) {
	client.SetToken(c.Token)
}

var once2 sync.Once

func (c *credentials) Unseal(client *api.Client) (err error) {

	sys := client.Sys()

	status, err := sys.SealStatus()
	if err != nil {
		return err
	}
	
	// Unseal is idemnpotent so no need to solve race conditions here
	if status.Sealed {
		status2, err := sys.Unseal(c.Keys[0])

		if err != nil {
			return err
		} else if status2.Sealed {
			return fmt.Errorf("Error, the vault was not properly unsealed")
		}
	}

	return nil

}

