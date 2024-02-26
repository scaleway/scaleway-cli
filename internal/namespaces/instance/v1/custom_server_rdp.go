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

	"golang.org/x/crypto/ssh"

	"github.com/scaleway/scaleway-cli/v2/internal/core"
	"github.com/scaleway/scaleway-cli/v2/internal/interactive"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type instanceServerGetRdpPasswordRequest struct {
	ServerID string
	Zone     scw.Zone
	Key      string
}

func instanceServerGetRdpPassword() *core.Command {
	return &core.Command{
		Short:     `Get your server rdp and decrypt it using your ssh key`,
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
	}
}

func instanceServerGetRdpPasswordRun(ctx context.Context, argsI interface{}) (i interface{}, e error) {
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
		return nil, fmt.Errorf("expected rsa private key, got %s", reflect.TypeOf(privateKey).String())
	}

	client := core.ExtractClient(ctx)
	apiInstance := instance.NewAPI(client)
	resp, err := apiInstance.GetEncryptedRdpPassword(&instance.GetEncryptedRdpPasswordRequest{
		Zone:     args.Zone,
		ServerID: args.ServerID,
	})
	if err != nil {
		return nil, err
	}
	if resp.Value == nil {
		return nil, fmt.Errorf("rdp password is nil")
	}

	encryptedRdpPassword, err := base64.RawStdEncoding.DecodeString(*resp.Value)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 encoded rdp password: %w", err)
	}

	password, err := rsa.DecryptPKCS1v15(nil, rsaKey, encryptedRdpPassword)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt rdp password: %w", err)
	}

	return struct {
		Username string
		Password string
	}{
		Username: "Administrator",
		Password: string(password),
	}, err
}

func parsePrivateKey(ctx context.Context, key []byte) (any, error) {
	privateKey, err := ssh.ParseRawPrivateKey(key)
	if err == nil {
		return privateKey, err
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
