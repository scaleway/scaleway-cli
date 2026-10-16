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

type Product string

const (
	Redis         Product = "redis"
	RDB           Product = "rdb"
	MongoDB       Product = "mongodb"
	Datawarehouse Product = "datawarehouse"
)

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

// CliConnectDeleteArgSpec is the shared cli-connect arg spec of the delete
// commands: it removes the credentials secret of the deleted resource.
func CliConnectDeleteArgSpec() *core.ArgSpec {
	return &core.ArgSpec{
		Name:       "cli-connect",
		Short:      `Remove the credentials secret of the resource, if it was created with cli-connect`,
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
	Product  Product
	Location string
	ID       string
	Username string
}

// SecretName is the name of the credentials secret of the resource: the resource
// UUID.
func (r Resource) SecretName() string {
	return r.ID
}

// SecretPath is the folder that holds all cli-connect credentials secrets.
func (r Resource) SecretPath() string {
	return SecretPath
}

// SecretTags let the credentials secret be filtered by product and locality.
func (r Resource) SecretTags() []string {
	return []string{string(r.Product), r.Location}
}

// SecretDescription returns the description persisted on the credentials secret:
// a copy-pastable scw connect command for the resource, so a user reading the
// secret in the console knows exactly how to reconnect.
// ponytail: only the products in connectSpecs get a copy-pastable command;
// add a row there when a new product is wired to cli-connect.
func (r Resource) SecretDescription() string {
	cmd := connectCommand(r.Product, r.Location, r.ID, r.Username)
	if cmd == "" {
		return description
	}

	return "Secret used to enable cli-connect. Connect directly using: " + cmd
}

// connectSpecs is the data behind the copy-pastable connect command stored
// in the credentials secret description. Add a row when a new product is
// wired to cli-connect.
type connectSpec struct {
	resource    string // "instance", "cluster", "deployment"
	locationArg string // "--region" or "--zone"
	withUser    bool
}

var connectSpecs = map[Product]connectSpec{
	RDB:           {resource: "instance", locationArg: "--region", withUser: true},
	MongoDB:       {resource: "instance", locationArg: "--region", withUser: true},
	Redis:         {resource: "cluster", locationArg: "--zone", withUser: false},
	Datawarehouse: {resource: "deployment", locationArg: "--region", withUser: false},
}

func connectCommand(product Product, location, resourceID, username string) string {
	spec, ok := connectSpecs[product]
	if !ok {
		return ""
	}
	user := ""
	if spec.withUser && username != "" {
		user = " username=" + username
	}

	return fmt.Sprintf(
		"scw %s %s connect %s%s %s %s",
		product,
		spec.resource,
		resourceID,
		user,
		spec.locationArg,
		location,
	)
}

// Credentials is the data stored in a credentials secret: username ->
// password, so all the users of a resource live in a single secret.
type Credentials = map[string]string

// SaveSecret stores the username/password pair in the credentials secret of the
// resource, naming the owning resource on failure. When the secret already
// exists a new version with the updated password is created on it.
func (r Resource) SaveSecret(
	ctx context.Context,
	secretManagerRegion scw.Region,
	username, password string,
) error {
	client := core.ExtractClient(ctx)
	api := secretSDK.NewAPI(client)

	fail := func(err error) error {
		return fmt.Errorf(
			"%s %s was created but its credentials could not be persisted: %w",
			r.Product,
			r.ID,
			err,
		)
	}

	secret, err := r.findOrCreateSecret(api, secretManagerRegion)
	if err != nil {
		return fail(err)
	}

	creds := Credentials{}
	if secret.VersionCount > 0 {
		creds = latestCredentials(api, secretManagerRegion, secret.ID)
	}
	creds[username] = password

	data, err := json.Marshal(creds)
	if err != nil {
		return err
	}

	if _, err := api.CreateSecretVersion(&secretSDK.CreateSecretVersionRequest{
		Region:   secretManagerRegion,
		SecretID: secret.ID,
		Data:     data,
	}); err != nil {
		return fail(fmt.Errorf("failed to create a version of secret %s: %w", secret.ID, err))
	}

	interactive.Printf(
		"Credentials persisted in secret %s (path: %s, id: %s)\n",
		r.SecretName(),
		r.SecretPath(),
		secret.ID,
	)

	return nil
}

// Fetch returns the password stored for the username in the credentials
// secret of the resource. found is false when the secret cannot be retrieved.
// ponytail: any lookup error (not just 404) degrades to "not found" so a broken
// Secret Manager access never blocks connecting; upgrade to erroring on
// non-404 if silent degradation proves confusing.
func (r Resource) Fetch(
	ctx context.Context,
	secretManagerRegion scw.Region,
	username string,
) (string, bool) {
	client := core.ExtractClient(ctx)
	api := secretSDK.NewAPI(client)

	resp, err := api.AccessSecretVersionByPath(&secretSDK.AccessSecretVersionByPathRequest{
		Region:     secretManagerRegion,
		Revision:   "latest",
		SecretPath: r.SecretPath(),
		SecretName: r.SecretName(),
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

// DeleteSecret removes the credentials secret of the resource, so that deleting a
// resource with cli-connect does not leave orphaned secrets behind. It is a
// no-op when the secret does not exist.
//
// ponytail: the delete commands call this right after the delete API call
// succeeds, i.e. before --wait confirms the resource is gone; if the wait
// then fails the secret is already removed. Move the call into the delete
// WaitFuncs if that proves confusing.
func (r Resource) DeleteSecret(ctx context.Context, secretManagerRegion scw.Region) error {
	client := core.ExtractClient(ctx)
	api := secretSDK.NewAPI(client)

	secret, err := findSecret(api, secretManagerRegion, r.SecretName(), r.SecretPath())
	if err != nil {
		return err
	}
	if secret == nil {
		core.ExtractLogger(ctx).Debugf("no credentials secret to remove for %s %s", r.Product, r.ID)

		return nil
	}
	if err := api.DeleteSecret(&secretSDK.DeleteSecretRequest{
		Region:   secretManagerRegion,
		SecretID: secret.ID,
	}); err != nil {
		return fmt.Errorf("failed to delete secret %s: %w", secret.ID, err)
	}

	interactive.Printf("Credentials secret %s (path: %s) removed\n", r.SecretName(), r.SecretPath())

	return nil
}

// RunUserUpdate runs the cli-connect flow around a user update: it prompts
// for a password when none is given and persists the credentials after a
// successful update.
func RunUserUpdate(
	ctx context.Context,
	cliConnect bool,
	product Product,
	secretManagerRegion scw.Region,
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
		if err := resource.SaveSecret(ctx, secretManagerRegion, name, **password); err != nil {
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
	api *secretSDK.API,
	secretManagerRegion scw.Region,
) (*secretSDK.Secret, error) {
	path := r.SecretPath()
	desc := r.SecretDescription()
	name := r.SecretName()

	secret, err := api.CreateSecret(&secretSDK.CreateSecretRequest{
		Region:      secretManagerRegion,
		Name:        name,
		Path:        &path,
		Description: &desc,
		Tags:        r.SecretTags(),
	})
	if err == nil {
		return secret, nil
	}

	// The slot is taken (a version of the CLI already saved these
	// credentials): reuse the existing secret.
	secret, listErr := findSecret(api, secretManagerRegion, name, path)
	if listErr != nil || secret == nil {
		return nil, fmt.Errorf("failed to create secret %s/%s: %w", path, name, err)
	}

	needsUpdate := secret.Description == nil || *secret.Description != desc ||
		!sameTags(secret.Tags, r.SecretTags())
	if needsUpdate {
		tags := r.SecretTags()
		if _, err := api.UpdateSecret(&secretSDK.UpdateSecretRequest{
			Region:      secretManagerRegion,
			SecretID:    secret.ID,
			Description: &desc,
			Tags:        &tags,
		}); err != nil {
			return nil, fmt.Errorf("failed to update secret %s: %w", secret.ID, err)
		}
	}

	return secret, nil
}

// findSecret returns the secret named exactly name under path, or nil when it
// does not exist. The ListSecrets name filter is a prefix match, hence the
// exact comparison.
func findSecret(
	api *secretSDK.API,
	secretManagerRegion scw.Region,
	name, path string,
) (*secretSDK.Secret, error) {
	secrets, err := api.ListSecrets(&secretSDK.ListSecretsRequest{
		Region: secretManagerRegion,
		Name:   &name,
		Path:   &path,
	}, scw.WithAllPages())
	if err != nil {
		return nil, fmt.Errorf("failed to list secrets %s/%s: %w", path, name, err)
	}

	for _, secret := range secrets.Secrets {
		if secret.Name == name {
			return secret, nil
		}
	}

	return nil, nil
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
func latestCredentials(
	api *secretSDK.API,
	secretManagerRegion scw.Region,
	secretID string,
) Credentials {
	resp, err := api.AccessSecretVersion(&secretSDK.AccessSecretVersionRequest{
		Region:   secretManagerRegion,
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
