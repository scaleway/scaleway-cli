package instance

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/interactive"
	iam "github.com/scaleway/scaleway-sdk-go/api/iam/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"golang.org/x/crypto/ssh"
)

type instanceServerGetRdpPasswordRequest struct {
	ServerID string
	Zone     scw.Zone
	Key      string
}

type ServerGetRdpPasswordResponse struct {
	Username          string
	Password          string
	SSHKeyID          *string
	SSHKeyDescription string
}

func instanceServerGetRdpPassword() *core.Command {
	return &core.Command{
		Short:     `Get your server rdp password and decrypt it using your ssh key`,
		Namespace: "instance",
		Verb:      "get-rdp-password",
		Resource:  "server",
		ArgsType:  reflect.TypeOf(instanceServerGetRdpPasswordRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "server-id",
				Short:      "Server ID to connect to",
				Required:   true,
				Positional: true,
			},
			{
				Name:  "key",
				Short: "Path of the SSH key used to decrypt the rdp password",
				Default: func(ctx context.Context) (value string, doc string) {
					homeDir := core.ExtractUserHomeDir(ctx)

					return filepath.Join(homeDir, ".ssh/id_rsa"), "~/.ssh/id_rsa"
				},
			},
			core.ZoneArgSpec(),
		},
		Run: instanceServerGetRdpPasswordRun,
		WaitFunc: func(ctx context.Context, argsI, respI any) (any, error) {
			// Wait only if response does not contain a password
			if _, isPasswd := respI.(*ServerGetRdpPasswordResponse); isPasswd {
				return respI, nil
			}

			args := argsI.(*instanceServerGetRdpPasswordRequest)
			apiInstance := instance.NewAPI(core.ExtractClient(ctx))
			_, err := apiInstance.WaitForServerRDPPassword(
				&instance.WaitForServerRDPPasswordRequest{
					Zone:          args.Zone,
					ServerID:      args.ServerID,
					RetryInterval: core.DefaultRetryInterval,
				},
			)
			if err != nil {
				return nil, err
			}

			// Retry command now that encrypted password is available
			return instanceServerGetRdpPasswordRun(ctx, argsI)
		},
	}
}

func instanceServerGetRdpPasswordRun(
	ctx context.Context,
	argsI any,
) (i any, e error) {
	args := argsI.(*instanceServerGetRdpPasswordRequest)

	if strings.HasPrefix(args.Key, "~") {
		args.Key = strings.Replace(args.Key, "~", core.ExtractUserHomeDir(ctx), 1)
	}

	rawKey, err := os.ReadFile(args.Key)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key file: %w", err)
	}

	privateKey, err := parsePrivateKey(ctx, rawKey)
	if err != nil {
		return nil, err
	}

	rsaKey, isRSA := privateKey.(*rsa.PrivateKey)
	if !isRSA {
		return nil, fmt.Errorf(
			"expected rsa private key, got %s",
			reflect.TypeOf(privateKey).String(),
		)
	}

	client := core.ExtractClient(ctx)
	apiInstance := instance.NewAPI(client)
	resp, err := apiInstance.GetServer(&instance.GetServerRequest{
		Zone:     args.Zone,
		ServerID: args.ServerID,
	})
	if err != nil {
		return nil, err
	}
	if resp.Server.AdminPasswordEncryptedValue == nil ||
		*resp.Server.AdminPasswordEncryptedValue == "" {
		return &core.CliError{
			Err:     errors.New("rdp password is empty"),
			Message: "RDP password is nil or empty in api response",
			Details: "Your server have no RDP password available",
			Hint:    "You may need to wait for your OS to start before having a generated RDP password, it can take more than 10 minutes.\nUse -w, --wait to wait for password to be available",
			Code:    1,
		}, nil
	}

	encryptedRdpPassword, err := base64.StdEncoding.DecodeString(
		*resp.Server.AdminPasswordEncryptedValue,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 encoded rdp password: %w", err)
	}

	password, err := rsa.DecryptPKCS1v15(nil, rsaKey, encryptedRdpPassword)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt rdp password: %w", err)
	}

	sshKeyDescription := ""
	if resp.Server.AdminPasswordEncryptionSSHKeyID != nil {
		iamAPI := iam.NewAPI(client)
		key, err := iamAPI.GetSSHKey(&iam.GetSSHKeyRequest{
			SSHKeyID: *resp.Server.AdminPasswordEncryptionSSHKeyID,
		})
		if err == nil {
			sshKeyDescription = key.Name
		}
	}

	return &ServerGetRdpPasswordResponse{
		Username:          "Administrator",
		Password:          string(password),
		SSHKeyID:          resp.Server.AdminPasswordEncryptionSSHKeyID,
		SSHKeyDescription: sshKeyDescription,
	}, err
}

func parsePrivateKey(ctx context.Context, key []byte) (any, error) {
	privateKey, err := ssh.ParseRawPrivateKey(key)
	if err == nil {
		return privateKey, nil
	}
	// Key may need a passphrase
	missingPassphraseError := &ssh.PassphraseMissingError{}
	if !errors.As(err, &missingPassphraseError) {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	passphrase, err := interactive.PromptPasswordWithConfig(&interactive.PromptPasswordConfig{
		Ctx:    ctx,
		Prompt: "passphrase",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to read input: %w", err)
	}

	privateKey, err = ssh.ParseRawPrivateKeyWithPassphrase(key, []byte(passphrase))
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	return privateKey, nil
}

func completeSSHKeyID(ctx context.Context, prefix string, _ any) core.AutocompleteSuggestions {
	resp, err := iam.NewAPI(core.ExtractClient(ctx)).
		ListSSHKeys(&iam.ListSSHKeysRequest{}, scw.WithAllPages())
	if err != nil {
		return nil
	}

	suggestion := make([]string, 0, len(resp.SSHKeys))

	for _, sshKey := range resp.SSHKeys {
		if strings.HasPrefix(sshKey.ID, prefix) {
			suggestion = append(suggestion, sshKey.ID)
		}
	}

	return suggestion
}
