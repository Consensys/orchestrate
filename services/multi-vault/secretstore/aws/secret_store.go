package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
)

// SecretStore can manage secrets on AWS secret manager
type SecretStore struct {
	client *secretsmanager.SecretsManager
}

// NewSecretStore returns a default configured AWS secretstore
func NewSecretStore() *SecretStore {
	secretStoreSession, _ := session.NewSession()
	return &SecretStore{
		client: secretsmanager.New(secretStoreSession),
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
		return errors.InternalError("client not set").SetComponent(component)
	}

	input := secretsmanager.DeleteSecretInput{
		RecoveryWindowInDays: aws.Int64(7),
		SecretId:             aws.String(key),
	}

	_, err = s.client.DeleteSecret(&input)
	if err != nil {
		return FromAWSError(err).SetComponent(component)
	}

	return nil
}

// List returns a list of available secrets
func (s *SecretStore) List() ([]string, error) {
	if s.client == nil {
		return []string{}, errors.InternalError("client not set").SetComponent(component)
	}

	input := secretsmanager.ListSecretsInput{}

	res, err := s.client.ListSecrets(&input)
	if err != nil {
		return []string{}, FromAWSError(err).SetComponent(component)
	}

	list := make([]string, len(res.SecretList))
	for i := 0; i < len(res.SecretList); i++ {
		list[i] = *res.SecretList[i].Name
	}

	return list, nil
}

// Load the secret value from the secret manager of AWS
func (s *SecretStore) Load(key string) (value string, ok bool, err error) {
	if s.client == nil {
		return "", false, errors.InternalError("client not set").SetComponent(component)
	}

	input := secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(key),
		VersionStage: aws.String("AWSCURRENT"),
	}

	res, err := s.client.GetSecretValue(&input)
	if err != nil {
		return "", false, FromAWSError(err).SetComponent(component)
	}

	return *res.SecretString, true, nil
}

func (s *SecretStore) create(key, value string) (err error) {
	if s.client == nil {
		return errors.InternalError("client not set").SetComponent(component)
	}

	input := secretsmanager.CreateSecretInput{
		Description:  aws.String("Miscellaneous Orchestrate secret"),
		Name:         aws.String(key),
		SecretString: aws.String(value),
	}

	_, err = s.client.CreateSecret(&input)
	if err != nil {
		return FromAWSError(err).SetComponent(component)
	}

	return nil
}

func (s *SecretStore) update(key, value string) (err error) {
	if s.client == nil {
		return errors.InternalError("client not set").SetComponent(component)
	}

	input := secretsmanager.PutSecretValueInput{
		SecretId:     aws.String(key),
		SecretString: aws.String(value),
	}

	_, err = s.client.PutSecretValue(&input)
	if err != nil {
		return FromAWSError(err).SetComponent(component)
	}

	return nil
}
