package hashicorps

import (
	"github.com/hashicorp/vault/api"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	aws "gitlab.com/ConsenSys/client/fr/core-stack/infra/aws-secret-manager.git/aws"
	"sync"
	"encoding/json"
	"sync"
)

type credentials struct {
	Keys []string 			`json: "keys"`
	KeysBase64 []string 	`json: "keys_base_64"`
	Token string 		`json: "root_token"`
}


var once sync.Once
var Credentials credentials


func (c *credentials) FetchFromAWS(
	client *secretsmanager.SecretsManager,
	name string,
) {
	once.Do(func() {
		secret := aws.SecretFromKey(name)
			.SetClient(client)
			.GetValue()

		c.fromEncoded(secret.value)
	})
}

func (c *credentials) fromEncoded(value string) (err error) {

	decoded := &credentials{}
	
	err := json.Unmarchal([]byte(value), decoded)
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
	var status api.SealStatusResponse

	status, err = sys.SealStatus()
	if err != nil {
		return err
	}
	
	// Unseal is idemnpotent so no need to solve race conditions here
	if status.Sealed {
		status, err = sys.Unseal(c.Keys[0])

		if err != nil {
			return err
		} else if status.Sealed {
			return fmt.Errorf("Error, the vault was not properly unsealed")
		}
	}

	return nil

}

