package aws

import (
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/aws"
	"fmt"
	//"github.com/google/uuid"
)

// Secret represent a key, value secret stored in AWS Secret Manager
type Secret struct {
	key string
	value string
	client *secretsmanager.SecretsManager
}

// List of the different availables constructors

// NewSecret is the default constructor for Secret
func NewSecret() *Secret {
	return &Secret{
		key: "",
		value: "",
		client: nil,
	}
}

// CreateSecret creates a Secret from key and value
func CreateSecret(key, value string) (*Secret) {
	return &Secret{
		key: key,
		value: value,
		client: nil,
	}
}

// SecretFromKey creates a secret from a key, it does not fetch the associated value.
func SecretFromKey(key string) (*Secret) {
	return &Secret{
		key: key,
		value: "",
		client: nil,
	}
}

// SetKey setter of attribute key for Secret struct object
func (sec *Secret) SetKey(key string) *Secret {
	sec.key = key
	return sec
}

// SetValue setter of attribute value for Secret struct object
func (sec *Secret) SetValue(value string) *Secret {
	sec.value = value
	return sec
}

// SetClient setter of attribute client for Secret struct object
func (sec *Secret) SetClient(client *secretsmanager.SecretsManager) *Secret {
	sec.client = client
	return sec
}

// SaveNew stores a new secret in the vault
func (sec *Secret) SaveNew() (*secretsmanager.CreateSecretOutput, error) {

	if sec.client == nil { return nil, fmt.Errorf("Client not set")}

	input := secretsmanager.CreateSecretInput{
		Description:        aws.String("Miscellaneous core-stack secret"),
		Name:               aws.String(sec.key),
		SecretString:       aws.String(sec.value),
	}

	res, err := sec.client.CreateSecret(&input)
	if err != nil {
		return nil, err
	}

	return res, nil	

}

// GetValue fetch the value from AWS SecretManager by key
func (sec *Secret) GetValue() (string, error) {

	if sec.client == nil { return "", fmt.Errorf("Client not set")}

	input := secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(sec.key),
		VersionStage: aws.String("AWSCURRENT"),
	}

	res, err := sec.client.GetSecretValue(&input)
	if err != nil {
		return "", err
	}

	sec.value = *res.SecretString

	return sec.value, nil
	
}

// Update the secret value stored in the aws secret manager
func (sec *Secret) Update() (*secretsmanager.PutSecretValueOutput, error) {

	if sec.client == nil { return nil, fmt.Errorf("Client not set")}

	input := secretsmanager.PutSecretValueInput{
		SecretId:           aws.String(sec.key),
		SecretString:       aws.String(sec.value),
	}

	res, err := sec.client.PutSecretValue(&input)
	if err != nil {
		return nil, err
	}

	return res, nil	

}


// Delete remove the key from the secret manager
func (sec *Secret) Delete() (*secretsmanager.DeleteSecretOutput, error) {

	if sec.client == nil { return nil, fmt.Errorf("Client not set") }

	input := secretsmanager.DeleteSecretInput{
		RecoveryWindowInDays: aws.Int64(7),
		SecretId:             aws.String(sec.key),
	}

	res, err := sec.client.DeleteSecret(&input); 
	if err != nil {
		return nil, err
	}

	return res, nil	

} 

// List retrieve all the keys availables in the secret manager
func (sec *Secret) List() (*secretsmanager.ListSecretsOutput, error) {

	if sec.client == nil { return nil, fmt.Errorf("Client not set") }

	input := secretsmanager.ListSecretsInput{}

	res, err := sec.client.ListSecrets(&input); 
	if err != nil {
		return nil, err
	}

	return res, nil	

} 

