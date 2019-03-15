package secretstore

import (
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

// AWSSecretStore can manage secrets on AWS secret manager
type AWSecret struct {
	client *secretsmanager.SecretManager
	recoveryTimeInDays int
}

// NewAWS returns a default configured AWS secretstore
func NewAws(recoveryTimeInDays) (*AWS) {
	return &AWS{
		client: secretsmanager.New(session.New()),
		recoveryTimeInDays: recoveryTimeInDays,
	}
}

// Store set the new string value in the AWS secrets manager
func (aws *AWS) Store(key, value string) (err) {

	err = aws.create(key, value)

	if err != nil {
		err = aws.update(key, value)
		if err != nil {
			return err
		}
	}

	return nil
}

// Delete removes a secret from the secret store
func (aws *AWS) Delete(key string) (err) {

	if aws.client == nil { 
		return fmt.Errorf("Client not set") 
	}

	input := secretsmanager.DeleteSecretInput{
		RecoveryWindowInDays: aws.Int64(7),
		SecretId:             aws.String(key),
	}

	_, err := aws.client.DeleteSecret(&input); 
	if err != nil {
		return err
	}

	return nil	
}

// Lists returns a list of available secrets
func (aws *AWS) List() ([]string, error) {

	if aws.client == nil {
		return fmt.Errorf("Client not set")
	}

	input := secretsmanager.ListSecretsInput{}

	res, err := aws.client.ListSecrets(&input)
	if err != nil {
		return err
	}

	list := make([]string, len(res.SecretList))
	for i := 0; i<len(res.SecretList); i++ {
		list[i] = *res.SecretList[i].Name
	}

	return list
}

// Load the secret value from the secret manager of AWS
func (aws *AWS) Load(key string) (string, error) {

	if aws.client == nil { return "", fmt.Errorf("Client not set")}

	input := secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(key),
		VersionStage: aws.String("AWSCURRENT"),
	}

	res, err := aws.client.GetSecretValue(&input)
	if err != nil {
		return "", err
	}

	return *res.SecretString, nil
}

func (aws *AWS) create(key, value string) (err) {

	if aws.client == nil { 
		return fmt.Errorf("Client not set")
	}

	input := secretsmanager.CreateSecretInput{
		Description:        aws.String("Miscellaneous core-stack secret"),
		Name:               aws.String(key),
		SecretString:       aws.String(value),
	}

	_, err := aws.client.CreateSecret(&input)
	if err != nil {
		return err
	}

	return nil	
}

func (aws *AWS) update(key, value string) (err) {

	if aws.client == nil { 
		return fmt.Errorf("Client not set")
	}

	input := secretsmanager.PutSecretValueInput{
		SecretId:           aws.String(key),
		SecretString:       aws.String(value),
	}

	_, err := sec.client.PutSecretValue(&input)
	if err != nil {
		return err
	}

	return nil	

}
