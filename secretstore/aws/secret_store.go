package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

// SecretStore can manage secrets on AWS secret manager
type SecretStore struct {
	client *secretsmanager.SecretsManager
}

// NewSecretStore returns a default configured AWS secretstore
func NewSecretStore() *SecretStore {
	return &SecretStore{
		client: secretsmanager.New(session.New()),
	}
}

// Store set the new string value in the AWS secrets manager
func (s *SecretStore) Store(key, value string) (err error) {
	err = s.create(key, value)
	if err != nil {
		err = s.update(key, value)
		if err != nil {
			return err
		}
	}

	return nil
}

// Delete removes a secret from the secret store
func (s *SecretStore) Delete(key string) (err error) {
	if s.client == nil {
		return fmt.Errorf("Client not set")
	}

	input := secretsmanager.DeleteSecretInput{
		RecoveryWindowInDays: aws.Int64(7),
		SecretId:             aws.String(key),
	}

	_, err = s.client.DeleteSecret(&input)
	if err != nil {
		return err
	}

	return nil
}

// List returns a list of available secrets
func (s *SecretStore) List() ([]string, error) {
	if s.client == nil {
		return []string{}, fmt.Errorf("Client not set")
	}

	input := secretsmanager.ListSecretsInput{}

	res, err := s.client.ListSecrets(&input)
	if err != nil {
		return []string{}, err
	}

	list := make([]string, len(res.SecretList))
	for i := 0; i < len(res.SecretList); i++ {
		list[i] = *res.SecretList[i].Name
	}

	return list, nil
}

// Load the secret value from the secret manager of AWS
func (s *SecretStore) Load(key string) (string, bool, error) {
	if s.client == nil {
		return "", false, fmt.Errorf("Client not set")
	}

	input := secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(key),
		VersionStage: aws.String("AWSCURRENT"),
	}

	res, err := s.client.GetSecretValue(&input)
	if err != nil {
		return "", false, err
	}

	return *res.SecretString, true, nil
}

func (s *SecretStore) create(key, value string) (err error) {
	if s.client == nil {
		return fmt.Errorf("Client not set")
	}

	input := secretsmanager.CreateSecretInput{
		Description:  aws.String("Miscellaneous core-stack secret"),
		Name:         aws.String(key),
		SecretString: aws.String(value),
	}

	_, err = s.client.CreateSecret(&input)
	if err != nil {
		return err
	}

	return nil
}

func (s *SecretStore) update(key, value string) (err error) {
	if s.client == nil {
		return fmt.Errorf("Client not set")
	}

	input := secretsmanager.PutSecretValueInput{
		SecretId:     aws.String(key),
		SecretString: aws.String(value),
	}

	_, err = s.client.PutSecretValue(&input)
	if err != nil {
		return err
	}

	return nil
}
