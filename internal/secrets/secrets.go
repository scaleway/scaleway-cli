// Package secrets persists product credentials in Secret Manager secrets so
// that scw connect commands can retrieve them without prompting the user.
package secrets

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/interactive"
	secretSDK "github.com/scaleway/scaleway-sdk-go/api/secret/v1beta1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

const (
	// SecretName is the name of the secrets created by the CLI.
	SecretName = "cli"

	// FastConnectArgSpecShort is the shared --help text for the fast-connect flag.
	FastConnectArgSpecShort = `Persist the credentials in a Secret Manager secret so that the connect command can retrieve them without prompting`
)

// FastConnectArgSpec is the shared --fast-connect flag spec.
func FastConnectArgSpec() *core.ArgSpec {
	return &core.ArgSpec{
		Name:       "fast-connect",
		Short:      FastConnectArgSpecShort,
		Required:   false,
		Deprecated: false,
		Positional: false,
		Default:    core.DefaultValueSetter("false"),
	}
}

var description = "secret used by the cli to ease up connection"

// Credentials are the username/password pair stored by fast-connect.
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Path returns the path of the credentials secret of a product resource:
// /<product>/<location>/<id>/<username> (the user SRN without the "srn://" prefix).
// location is the product's native scope (region or zone).
func Path(product, location, id, username string) string {
	return fmt.Sprintf("/%s/%s/%s/%s", product, location, id, username)
}

// Persist stores the credentials in a secret named SecretName at the given path.
// When a secret already exists at that path (e.g. the password of a user is
// updated after the credentials were persisted at resource creation) a new
// version is created on the existing secret instead.
func Persist(ctx context.Context, region scw.Region, path string, creds Credentials) error {
	client := core.ExtractClient(ctx)
	api := secretSDK.NewAPI(client)

	secret, err := api.CreateSecret(&secretSDK.CreateSecretRequest{
		Region:      region,
		Name:        SecretName,
		Path:        &path,
		Description: &description,
	})
	if err != nil {
		secretName := SecretName
		existing, listErr := api.ListSecrets(&secretSDK.ListSecretsRequest{
			Region: region,
			Name:   &secretName,
			Path:   &path,
		}, scw.WithAllPages())
		if listErr != nil || len(existing.Secrets) == 0 {
			return fmt.Errorf("failed to create secret %q: %w", path, err)
		}
		secret = existing.Secrets[0]
	}

	data, err := json.Marshal(creds)
	if err != nil {
		return err
	}

	_, err = api.CreateSecretVersion(&secretSDK.CreateSecretVersionRequest{
		Region:   region,
		SecretID: secret.ID,
		Data:     data,
	})
	if err != nil {
		return fmt.Errorf("failed to create a version of secret %s: %w", secret.ID, err)
	}

	fmt.Printf(
		"Credentials persisted in secret %s (path: %s, id: %s)\n",
		SecretName,
		path,
		secret.ID,
	)

	return nil
}

// PersistCredentials stores the credentials in a secret at
// Path(product, location, id, username), naming the owning resource on failure.
func PersistCredentials(
	ctx context.Context,
	region scw.Region,
	product, location, resourceID, username, password, resourceType string,
) error {
	if err := Persist(ctx, region, Path(product, location, resourceID, username), Credentials{
		Username: username,
		Password: password,
	}); err != nil {
		return fmt.Errorf(
			"%s %s was created but its credentials could not be persisted: %w",
			resourceType,
			resourceID,
			err,
		)
	}

	return nil
}

// Fetch returns the credentials stored at the given path. found is false when
// the secret cannot be retrieved.
// ponytail: any lookup error (not just 404) degrades to "not found" so a broken
// Secret Manager access never blocks connecting; upgrade to erroring on
// non-404 if silent degradation proves confusing.
func Fetch(ctx context.Context, region scw.Region, path string) (Credentials, bool) {
	client := core.ExtractClient(ctx)
	api := secretSDK.NewAPI(client)

	resp, err := api.AccessSecretVersionByPath(&secretSDK.AccessSecretVersionByPathRequest{
		Region:     region,
		Revision:   "latest",
		SecretPath: path,
		SecretName: SecretName,
	})
	if err != nil {
		core.ExtractLogger(ctx).Debugf("no credentials secret at %s: %s", path, err)

		return Credentials{}, false
	}

	var creds Credentials
	if err := json.Unmarshal(resp.Data, &creds); err != nil {
		core.ExtractLogger(ctx).Debugf("failed to parse credentials secret at %s: %s", path, err)

		return Credentials{}, false
	}

	return creds, true
}

// FetchOrPrompt returns the credentials stored at path, or prompts for a
// password when no usable secret exists.
func FetchOrPrompt(ctx context.Context, region scw.Region, path string) (Credentials, error) {
	creds, found := Fetch(ctx, region, path)
	if found && creds.Password != "" {
		return creds, nil
	}

	password, err := PromptPassword(ctx)
	if err != nil {
		return Credentials{}, err
	}

	return Credentials{Password: password}, nil
}

// PromptPassword asks the user for a password interactively so that it does
// not end up in the shell history.
func PromptPassword(ctx context.Context) (string, error) {
	password, err := interactive.PromptPasswordWithConfig(ctx, &interactive.PromptPasswordConfig{
		Prompt: "Password",
	})
	if err != nil {
		return "", fmt.Errorf("failed to get password: %w", err)
	}

	return password, nil
}
