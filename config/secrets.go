package config

import (
	"context"
	"fmt"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

// loadSecrets pulls the necessary information from google's secret manager
func LoadSecrets() error {
	// some context to string this all together
	ctx := context.Background()

	// create the secret client
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return err
	}

	// load the secrests
	BotToken, err = loadSecret(ctx, client, "JEEVES_TOKEN")
	if err != nil {
		return err
	}
	DBUser, err = loadSecret(ctx, client, "JEEVES_DB_USER")
	if err != nil {
		return err
	}
	DBPassword, err = loadSecret(ctx, client, "JEEVES_DB_PASSWORD")
	if err != nil {
		return err
	}

	// nothing went wrong
	return nil
}

// loadSecret pulls the latest secret value from google's SecretManager
func loadSecret(ctx context.Context, client *secretmanager.Client, name string) (string, error) {
	result, err := client.AccessSecretVersion(ctx, &secretmanagerpb.AccessSecretVersionRequest{
		Name: fmt.Sprintf("projects/%s/secrets/%s/versions/latest", GoogleCloudProject, name),
	})
	if err != nil {
		return "", err
	}

	return string(result.Payload.Data), nil
}
