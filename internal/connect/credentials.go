// Package connect persists product credentials in Secret Manager so that
// scw connect commands can retrieve them without prompting the user.
package connect

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/interactive"
	secretSDK "github.com/scaleway/scaleway-sdk-go/api/secret/v1beta1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// SecretPath is the Secret Manager folder that holds all cli-connect
// credentials, one secret per resource. It must stay absolute: the by-path
// Secret Manager endpoints reject relative paths.
const SecretPath = "/cli-connect"

// CliConnectArgSpec is the shared --cli-connect flag spec.
func CliConnectArgSpec() *core.ArgSpec {
	return &core.ArgSpec{
		Name:       "cli-connect",
		Short:      `Persist the credentials in a Secret Manager secret so that the connect command can retrieve them without prompting`,
		Required:   false,
		Deprecated: false,
		Positional: false,
		Default:    core.DefaultValueSetter("false"),
	}
}

var description = "secret used by the cli to ease up connection"

// Resource identifies the product resource whose credentials cli-connect
// persists in Secret Manager: a secret named after the resource UUID, in the
// cli-connect folder, tagged by product and locality.
type Resource struct {
	Product  string
	Location string
	ID       string
	Username string
}

// Name is the name of the credentials secret of the resource: the resource
// UUID.
func (r Resource) Name() string {
	return r.ID
}

// Path is the folder that holds all cli-connect credentials secrets.
func (r Resource) Path() string {
	return SecretPath
}

// Tags let the credentials secret be filtered by product and locality.
func (r Resource) Tags() []string {
	return []string{r.Product, r.Location}
}

// Description returns the description persisted on the credentials secret:
// a copy-pastable scw connect command for the resource, so a user reading the
// secret in the console knows exactly how to reconnect.
// ponytail: only the products with a hand-written connect command get a
// copy-pastable command; add a case in connectCommand when a new product is
// wired to cli-connect.
func (r Resource) Description() string {
	cmd := connectCommand(r.Product, r.Location, r.ID, r.Username)
	if cmd == "" {
		return description
	}

	return "Secret used to enable cli-connect. Connect directly using: " + cmd
}

func connectCommand(product, location, resourceID, username string) string {
	switch product {
	case "rdb", "mongodb":
		user := ""
		if username != "" {
			user = " username=" + username
		}

		return fmt.Sprintf(
			"scw %s instance connect %s%s --region %s",
			product,
			resourceID,
			user,
			location,
		)
	case "redis":
		return fmt.Sprintf("scw redis cluster connect %s --zone %s", resourceID, location)
	case "datawarehouse":
		return fmt.Sprintf(
			"scw datawarehouse deployment connect %s --region %s",
			resourceID,
			location,
		)
	default:
		return ""
	}
}

// Credentials is the data stored in a credentials secret: username ->
// password, so all the users of a resource live in a single secret.
type Credentials = map[string]string

// Persist stores the username/password pair in the credentials secret of the
// resource, naming the owning resource on failure. When the secret already
// exists a new version with the updated password is created on it.
func (r Resource) Persist(
	ctx context.Context,
	region scw.Region,
	username, password, resourceType string,
) error {
	client := core.ExtractClient(ctx)
	api := secretSDK.NewAPI(client)

	fail := func(err error) error {
		return fmt.Errorf(
			"%s %s was created but its credentials could not be persisted: %w",
			resourceType,
			r.ID,
			err,
		)
	}

	secret, err := r.findOrCreateSecret(region, api)
	if err != nil {
		return fail(err)
	}

	creds := Credentials{}
	if secret.VersionCount > 0 {
		creds = latestCredentials(region, api, secret.ID)
	}
	creds[username] = password

	data, err := json.Marshal(creds)
	if err != nil {
		return err
	}

	if _, err := api.CreateSecretVersion(&secretSDK.CreateSecretVersionRequest{
		Region:   region,
		SecretID: secret.ID,
		Data:     data,
	}); err != nil {
		return fail(fmt.Errorf("failed to create a version of secret %s: %w", secret.ID, err))
	}

	interactive.Printf(
		"Credentials persisted in secret %s (path: %s, id: %s)\n",
		r.Name(),
		r.Path(),
		secret.ID,
	)

	return nil
}

// Fetch returns the password stored for the username in the credentials
// secret of the resource. found is false when the secret cannot be retrieved.
// ponytail: any lookup error (not just 404) degrades to "not found" so a broken
// Secret Manager access never blocks connecting; upgrade to erroring on
// non-404 if silent degradation proves confusing.
func (r Resource) Fetch(ctx context.Context, region scw.Region, username string) (string, bool) {
	client := core.ExtractClient(ctx)
	api := secretSDK.NewAPI(client)

	resp, err := api.AccessSecretVersionByPath(&secretSDK.AccessSecretVersionByPathRequest{
		Region:     region,
		Revision:   "latest",
		SecretPath: r.Path(),
		SecretName: r.Name(),
	})
	if err != nil {
		core.ExtractLogger(ctx).Debugf("no credentials secret for %s %s: %s", r.Product, r.ID, err)

		return "", false
	}

	var creds Credentials
	if err := json.Unmarshal(resp.Data, &creds); err != nil {
		core.ExtractLogger(ctx).Debugf("failed to parse credentials secret %s: %s", r.ID, err)

		return "", false
	}

	password, ok := creds[username]
	if !ok || password == "" {
		core.ExtractLogger(ctx).Debugf("no credentials for user %s in secret %s", username, r.ID)

		return "", false
	}

	return password, true
}

// FetchOrPrompt returns the password stored for the username, or prompts for
// one when no usable secret exists.
func (r Resource) FetchOrPrompt(
	ctx context.Context,
	region scw.Region,
	username string,
) (string, error) {
	if password, found := r.Fetch(ctx, region, username); found {
		return password, nil
	}

	return PromptPassword(ctx)
}

// RunUserUpdate runs the cli-connect flow around a user update: it prompts
// for a password when none is given and persists the credentials after a
// successful update.
func RunUserUpdate(
	ctx context.Context,
	cliConnect bool,
	product, resourceType string,
	region scw.Region,
	location, resourceID, name string,
	password **string,
	update func(ctx context.Context) (any, error),
) (any, error) {
	if cliConnect && *password == nil {
		prompted, err := PromptPassword(ctx)
		if err != nil {
			return nil, err
		}
		*password = &prompted
	}

	res, err := update(ctx)
	if err != nil {
		return nil, err
	}

	if cliConnect {
		resource := Resource{
			Product:  product,
			Location: location,
			ID:       resourceID,
			Username: name,
		}
		if err := resource.Persist(ctx, region, name, **password, resourceType); err != nil {
			return nil, err
		}
	}

	return res, nil
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

// findOrCreateSecret returns the credentials secret of the resource, creating
// it when missing. When the secret already exists its metadata is refreshed
// so that tags and the connect command stay up to date.
func (r Resource) findOrCreateSecret(
	region scw.Region,
	api *secretSDK.API,
) (*secretSDK.Secret, error) {
	path := r.Path()
	desc := r.Description()
	name := r.Name()

	secret, err := api.CreateSecret(&secretSDK.CreateSecretRequest{
		Region:      region,
		Name:        name,
		Path:        &path,
		Description: &desc,
		Tags:        r.Tags(),
	})
	if err == nil {
		return secret, nil
	}

	// The slot is taken (a version of the CLI already persisted these
	// credentials): reuse the existing secret.
	existing, listErr := api.ListSecrets(&secretSDK.ListSecretsRequest{
		Region: region,
		Name:   &name,
		Path:   &path,
	}, scw.WithAllPages())
	if listErr != nil || len(existing.Secrets) == 0 {
		return nil, fmt.Errorf("failed to create secret %s/%s: %w", path, name, err)
	}
	secret = existing.Secrets[0]

	needsUpdate := secret.Description == nil || *secret.Description != desc ||
		!sameTags(secret.Tags, r.Tags())
	if needsUpdate {
		tags := r.Tags()
		if _, err := api.UpdateSecret(&secretSDK.UpdateSecretRequest{
			Region:      region,
			SecretID:    secret.ID,
			Description: &desc,
			Tags:        &tags,
		}); err != nil {
			return nil, fmt.Errorf("failed to update secret %s: %w", secret.ID, err)
		}
	}

	return secret, nil
}

func sameTags(a, b []string) bool {
	a = append([]string(nil), a...)
	b = append([]string(nil), b...)
	slices.Sort(a)
	slices.Sort(b)

	return slices.Equal(a, b)
}

// latestCredentials returns the credentials stored in the latest version of
// the secret, or an empty map when it cannot be read.
func latestCredentials(region scw.Region, api *secretSDK.API, secretID string) Credentials {
	resp, err := api.AccessSecretVersion(&secretSDK.AccessSecretVersionRequest{
		Region:   region,
		SecretID: secretID,
		Revision: "latest",
	})
	if err != nil {
		return Credentials{}
	}

	var creds Credentials
	if err := json.Unmarshal(resp.Data, &creds); err != nil {
		return Credentials{}
	}

	return creds
}
