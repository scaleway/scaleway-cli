// Package connect persists product credentials in Secret Manager connect so
// that scw connect commands can retrieve them without prompting the user.
package connect

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
	// SecretName is the name of the secret created by the CLI for cli-connect.
	SecretName = "cli-connect"
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

var description = "secret used by the cli to ease up connection"

// Resource identifies the product resource whose credentials cli-connect
// persists in Secret Manager: /<product>/<location>/<id>/<username>.
// location is the product's native scope (region or zone).
type Resource struct {
	Product  string
	Location string
	ID       string
	Username string
}

// Path returns the path of the credentials secret of the resource
// (the user SRN without the "srn://" prefix).
func (r Resource) Path() string {
	return fmt.Sprintf("/%s/%s/%s/%s", r.Product, r.Location, r.ID, r.Username)
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

// Credentials are the username/password pair stored by cli-connect.
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Persist stores the credentials in a secret named SecretName at r.Path(),
// naming the owning resource on failure. When a secret already exists at that
// path (e.g. the password of a user is updated after the credentials were
// persisted at resource creation) a new version is created on the existing
// secret instead.
func (r Resource) Persist(
	ctx context.Context,
	region scw.Region,
	password, resourceType string,
) error {
	err := persist(
		ctx,
		region,
		r.Path(),
		r.Description(),
		Credentials{
			Username: r.Username,
			Password: password,
		},
	)
	if err != nil {
		return fmt.Errorf(
			"%s %s was created but its credentials could not be persisted: %w",
			resourceType,
			r.ID,
			err,
		)
	}

	return nil
}

func persist(ctx context.Context, region scw.Region, path, desc string, creds Credentials) error {
	client := core.ExtractClient(ctx)
	api := secretSDK.NewAPI(client)

	secret, err := api.CreateSecret(&secretSDK.CreateSecretRequest{
		Region:      region,
		Name:        SecretName,
		Path:        &path,
		Description: &desc,
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
		// The secret may have been created by an older version of the CLI with
		// the generic description: refresh it so it always carries the
		// copy-pastable connect command.
		if secret.Description == nil || *secret.Description != desc {
			if _, err := api.UpdateSecret(&secretSDK.UpdateSecretRequest{
				Region:      region,
				SecretID:    secret.ID,
				Description: &desc,
			}); err != nil {
				return fmt.Errorf(
					"failed to update the description of secret %s: %w",
					secret.ID,
					err,
				)
			}
		}
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

	interactive.Printf(
		"Credentials persisted in secret %s (path: %s, id: %s)\n",
		SecretName,
		path,
		secret.ID,
	)

	return nil
}

// Fetch returns the credentials stored at the resource's path. found is false
// when the secret cannot be retrieved.
// ponytail: any lookup error (not just 404) degrades to "not found" so a broken
// Secret Manager access never blocks connecting; upgrade to erroring on
// non-404 if silent degradation proves confusing.
func (r Resource) Fetch(ctx context.Context, region scw.Region) (Credentials, bool) {
	path := r.Path()
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

// FetchOrPrompt returns the credentials stored for the resource, or prompts
// for a password when no usable secret exists.
func (r Resource) FetchOrPrompt(ctx context.Context, region scw.Region) (Credentials, error) {
	creds, found := r.Fetch(ctx, region)
	if found && creds.Password != "" {
		return creds, nil
	}

	password, err := PromptPassword(ctx)
	if err != nil {
		return Credentials{}, err
	}

	return Credentials{Password: password}, nil
}

// RunUserUpdate runs the cli-connect flow around a user update: it prompts
// for a password when none is given and persists the credentials after a
// successful update.
func RunUserUpdate(
	ctx context.Context,
	cliConnect bool,
	product, resourceType string,
	region scw.Region,
	resourceID, name string,
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
			Location: string(region),
			ID:       resourceID,
			Username: name,
		}
		if err := resource.Persist(ctx, region, **password, resourceType); err != nil {
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
