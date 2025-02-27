package backends

import (
	"encoding/json"
	"fmt"

	"github.com/IBM/argocd-vault-plugin/pkg/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/service/secretsmanager/secretsmanageriface"
)

// AWSSecretsManager is a struct for working with a AWS Secrets Manager backend
type AWSSecretsManager struct {
	Client secretsmanageriface.SecretsManagerAPI
}

// NewAWSSecretsManagerBackend initializes a new AWS Secrets Manager backend
func NewAWSSecretsManagerBackend(client secretsmanageriface.SecretsManagerAPI) *AWSSecretsManager {
	return &AWSSecretsManager{
		Client: client,
	}
}

// Login does nothing as a "login" is handled on the instantiation of the aws sdk
func (a *AWSSecretsManager) Login() error {
	return nil
}

// GetSecrets gets secrets from aws secrets manager and returns the formatted data
func (a *AWSSecretsManager) GetSecrets(path string, version string, annotations map[string]string) (map[string]interface{}, error) {
	input := &secretsmanager.GetSecretValueInput{
		SecretId:  aws.String(path),
		VersionId: aws.String(types.AWSCurrentSecretVersion),
	}

	if version != "" {
		input.VersionId = aws.String(version)
	}

	result, err := a.Client.GetSecretValue(input)
	if err != nil {
		return nil, err
	}

	var dat map[string]interface{}

	if result.SecretString != nil {
		err := json.Unmarshal([]byte(*result.SecretString), &dat)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("Could not find secret %s", path)
	}

	return dat, nil
}
