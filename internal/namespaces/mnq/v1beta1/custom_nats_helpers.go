package mnq

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/interactive"
	mnq "github.com/scaleway/scaleway-sdk-go/api/mnq/v1beta1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type NatsEntity struct {
	Name    string
	Content []byte
}

func makeDirectoryIfNotExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, os.ModeDir|0o755)
	}

	return nil
}

func wrapError(err error, message, name, path string) error {
	return &core.CliError{
		Err:     err,
		Message: fmt.Sprintf("%s into file %q", message, path),
		Details: fmt.Sprintf("You may want to delete created credentials %q", name),
		Code:    1,
	}
}

func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)

	return !os.IsNotExist(err)
}

func natsContextFrom(account *mnq.NatsAccount, credsPath string) ([]byte, error) {
	ctx := &natsContext{
		Description:     "Nats context created by Scaleway CLI",
		URL:             account.Endpoint,
		CredentialsPath: credsPath,
	}
	b, err := json.Marshal(ctx)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func writeFile(
	ctx context.Context,
	dir string,
	entity *NatsEntity,
	extension string,
) (string, error) {
	path := filepath.Join(dir, entity.Name+"."+extension)
	if err := makeDirectoryIfNotExists(dir); err != nil {
		return "", wrapError(err, "Failed to create directory", entity.Name, path)
	}
	if FileExists(path) {
		overWrite, err := promptOverWriteFile(ctx, path)
		if err != nil {
			return "", wrapError(err, "Failed to prompt for overwrite", entity.Name, path)
		}
		if !overWrite {
			return "", wrapError(nil, "File already exists", entity.Name, path)
		}
	}
	if err := os.WriteFile(path, entity.Content, 0o600); err != nil {
		return "", wrapError(err, "Failed to write file", entity.Name, path)
	}
	_, _ = interactive.Println(entity.Name + " file has been successfully written to " + path)

	return path, nil
}

func getNATSContextDir(ctx context.Context) (string, error) {
	xdgConfigHome := core.ExtractEnv(ctx, "XDG_CONFIG_HOME")
	interactive.Println("xdgConfigHome:", xdgConfigHome)
	if xdgConfigHome == "" {
		homeDir := core.ExtractEnv(ctx, "HOME")
		if homeDir == "" {
			return "", errors.New("both XDG_CONFIG_HOME and HOME are not set")
		}

		return filepath.Join(homeDir, ".config", "nats", "context"), nil
	}

	return xdgConfigHome, nil
}

func saveNATSCredentials(
	ctx context.Context,
	creds *mnq.NatsCredentials,
	natsAccount *mnq.NatsAccount,
) (string, error) {
	natsContextDir, err := getNATSContextDir(ctx)
	if err != nil {
		return "", err
	}
	credsEntity := &NatsEntity{
		Name:    creds.Name,
		Content: []byte(creds.Credentials.Content),
	}
	credsPath, err := writeFile(ctx, natsContextDir, credsEntity, "creds")
	if err != nil {
		return "", err
	}

	natsContent, err := natsContextFrom(natsAccount, credsPath)
	if err != nil {
		return "", err
	}

	contextEntity := &NatsEntity{
		Name:    natsAccount.Name,
		Content: natsContent,
	}

	contextPath, err := writeFile(ctx, natsContextDir, contextEntity, "json")
	if err != nil {
		return "", err
	}

	return contextPath, nil
}

func getNatsAccountID(
	ctx context.Context,
	args *CreateContextRequest,
	api *mnq.NatsAPI,
) (*mnq.NatsAccount, error) {
	var natsAccount *mnq.NatsAccount
	if args.NatsAccountID == "" {
		natsAccountsResp, err := api.ListNatsAccounts(&mnq.NatsAPIListNatsAccountsRequest{
			Region: args.Region,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to list nats account: %w", err)
		}
		natsAccount, err = promptNatsAccounts(
			ctx,
			natsAccountsResp.NatsAccounts,
			natsAccountsResp.TotalCount,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to list nats account: %w", err)
		}
	} else {
		var err error
		natsAccount, err = api.GetNatsAccount(&mnq.NatsAPIGetNatsAccountRequest{
			Region:        args.Region,
			NatsAccountID: args.NatsAccountID,
		}, scw.WithContext(ctx))
		if err != nil {
			return nil, fmt.Errorf("failed to get nats account: %w", err)
		}
	}

	return natsAccount, nil
}
