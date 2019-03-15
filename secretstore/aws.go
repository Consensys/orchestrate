package secretstore

import (
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"fmt"
)

// AWS can manage secrets on AWS secret manager
type AWS struct {
	client *secretsmanager.SecretsManager
	recoveryTimeInDays int
}

// NewAWS returns a default configured AWS secretstore
func NewAWS(recoveryTimeInDays int) (*AWS) {
	return &AWS{
		client: secretsmanager.New(session.New()),
		recoveryTimeInDays: recoveryTimeInDays,
	}
}

// Store set the new string value in the AWS secrets manager
func (ss *AWS) Store(key, value string) (err error) {

	err = ss.create(key, value)

	if err != nil {
		err = ss.update(key, value)
		if err != nil {
			return err
		}
	}

	return nil
}

// Delete removes a secret from the secret store
func (ss *AWS) Delete(key string) (err error) {

	if ss.client == nil { 
		return fmt.Errorf("Client not set") 
	}

	input := secretsmanager.DeleteSecretInput{
		RecoveryWindowInDays: aws.Int64(7),
		SecretId:             aws.String(key),
	}

	_, err = ss.client.DeleteSecret(&input); 
	if err != nil {
		return err
	}

	return nil	
}

// List returns a list of available secrets
func (ss *AWS) List() ([]string, error) {

	if ss.client == nil {
		return []string{}, fmt.Errorf("Client not set")
	}

	input := secretsmanager.ListSecretsInput{}

	res, err := ss.client.ListSecrets(&input)
	if err != nil {
		return []string{}, err
	}

	list := make([]string, len(res.SecretList))
	for i := 0; i<len(res.SecretList); i++ {
		list[i] = *res.SecretList[i].Name
	}

	return list, nil
}

// Load the secret value from the secret manager of AWS
func (ss *AWS) Load(key string) (string, error) {

	if ss.client == nil { return "", fmt.Errorf("Client not set")}

	input := secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(key),
		VersionStage: aws.String("AWSCURRENT"),
	}

	res, err := ss.client.GetSecretValue(&input)
	if err != nil {
		return "", err
	}

	return *res.SecretString, nil
}

func (ss *AWS) create(key, value string) (err error) {

	if ss.client == nil { 
		return fmt.Errorf("Client not set")
	}

	input := secretsmanager.CreateSecretInput{
		Description:        aws.String("Miscellaneous core-stack secret"),
		Name:               aws.String(key),
		SecretString:       aws.String(value),
	}

	_, err = ss.client.CreateSecret(&input)
	if err != nil {
		return err
	}

	return nil	
}

func (ss *AWS) update(key, value string) (err error) {

	if ss.client == nil { 
		return fmt.Errorf("Client not set")
	}

	input := secretsmanager.PutSecretValueInput{
		SecretId:           aws.String(key),
		SecretString:       aws.String(value),
	}

	_, err = ss.client.PutSecretValue(&input)
	if err != nil {
		return err
	}

	return nil	

}
