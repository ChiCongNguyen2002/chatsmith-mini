package configs

import (
	"context"
	"fmt"
	"os"
	"strings"

	"chatsmith-mini/pkg/failure"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"github.com/spf13/viper"
)

const (
	envGCPProjectID = "GCP_PROJECT_ID"
	secretVersion   = "latest"
	secretResource  = "projects/%s/secrets/%s/versions/%s"

	secretDatabaseURL        = "database-url"
	secretAnthropicAPIKey    = "anthropic-api-key"
	secretGoogleClientSecret = "google-client-secret"
	secretSessionSecret      = "session-secret"
)

var managedSecrets = map[string]string{
	KeyDatabaseURL:            secretDatabaseURL,
	KeyProviderAnthropicKey:   secretAnthropicAPIKey,
	KeyAuthGoogleClientSecret: secretGoogleClientSecret,
	KeyAuthSessionSecret:      secretSessionSecret,
}

func loadSecrets(ctx context.Context, v *viper.Viper) error {
	project := os.Getenv(envGCPProjectID)
	if project == "" {
		return nil
	}

	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return failure.ErrWithTrace(err)
	}
	defer client.Close()

	for key, name := range managedSecrets {
		value, err := accessSecret(ctx, client, fmt.Sprintf(secretResource, project, name, secretVersion))
		if err != nil {
			return err
		}
		v.Set(key, value)
	}

	return nil
}

func accessSecret(ctx context.Context, client *secretmanager.Client, resource string) (string, error) {
	resp, err := client.AccessSecretVersion(ctx, &secretmanagerpb.AccessSecretVersionRequest{Name: resource})
	if err != nil {
		return "", failure.ErrWithTrace(err)
	}

	return strings.TrimSpace(string(resp.GetPayload().GetData())), nil
}
