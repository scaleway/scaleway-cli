package instance_test

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	iamCLI "github.com/scaleway/scaleway-cli/v2/internal/namespaces/iam/v1alpha1"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/instance/v1"
	"github.com/scaleway/scaleway-cli/v2/internal/testhelpers"
	iam "github.com/scaleway/scaleway-sdk-go/api/iam/v1alpha1"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/ssh"
)

// tryLoadKey will try to load an RSA SSH Key from given path.
// If not found, it will generate the key and create the file.
func loadRSASSHKey(path string) (*rsa.PrivateKey, error) {
	pemContent, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to read test key: %w", err)
	}

	if *core.UpdateCassettes {
		privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			return nil, fmt.Errorf("failed to generate key: %w", err)
		}
		privatePEM, err := ssh.MarshalPrivateKey(
			privateKey,
			"test-cli-instance-server-get-rdp-password",
		)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal private key: %w", err)
		}
		pemContent = pem.EncodeToMemory(privatePEM)
		err = os.WriteFile(path, pemContent, 0o600)
		if err != nil {
			return nil, fmt.Errorf("failed to save test key: %w", err)
		}
	} else if len(pemContent) == 0 {
		return nil, errors.New("empty test key, should be updated with cassettes")
	}

	key, err := ssh.ParseRawPrivateKey(pemContent)
	if err != nil {
		return nil, fmt.Errorf("failed to parse test key: %w", err)
	}

	privateKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf(
			"failed to assert private key type, expected *rsa.PrivateKey, got: %v",
			reflect.TypeOf(privateKey),
		)
	}

	return privateKey, nil
}

// generateRSASSHKey generates an RSA SSH Key and upload it to IAM.
// IAMSSHKey object is stored in metaKey.
func generateRSASSHKey(metaKey string) func(beforeFunc *core.BeforeFuncCtx) error {
	return func(ctx *core.BeforeFuncCtx) error {
		privateKey, err := loadRSASSHKey("testfixture/id_rsa")
		if err != nil {
			return fmt.Errorf("failed to load private key: %w", err)
		}
		privatePEM, err := ssh.MarshalPrivateKey(
			privateKey,
			"test-cli-instance-server-get-rdp-password",
		)
		if err != nil {
			return fmt.Errorf("failed to marshal private key: %w", err)
		}
		publicKey, err := ssh.NewPublicKey(privateKey.Public())
		if err != nil {
			return fmt.Errorf("failed to load public key: %w", err)
		}
		authorizedKey := ssh.MarshalAuthorizedKey(publicKey)

		sshDir := filepath.Join(ctx.OverrideEnv["HOME"], ".ssh")
		err = os.MkdirAll(sshDir, 0o700)
		if err != nil {
			return fmt.Errorf("failed to create directory %q: %w", sshDir, err)
		}

		err = os.WriteFile(filepath.Join(sshDir, "id_rsa"), pem.EncodeToMemory(privatePEM), 0o600)
		if err != nil {
			return fmt.Errorf("failed to write private key: %w", err)
		}
		err = os.WriteFile(filepath.Join(sshDir, "id_rsa.pub"), authorizedKey, 0o600)
		if err != nil {
			return fmt.Errorf("failed to write public key: %w", err)
		}

		api := iam.NewAPI(ctx.Client)
		projectID, exists := ctx.Client.GetDefaultProjectID()
		if !exists {
			return errors.New("missing project id")
		}
		key, err := api.CreateSSHKey(&iam.CreateSSHKeyRequest{
			Name:      "test-cli-instance-server-get-rdp-password",
			PublicKey: string(authorizedKey),
			ProjectID: projectID,
		})
		if err != nil {
			return fmt.Errorf("failed to create iam ssh key: %w", err)
		}

		ctx.Meta[metaKey] = key

		return nil
	}
}

func Test_ServerGetRdpPassword(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: core.NewCommandsMerge(
			instance.GetCommands(),
			iamCLI.GetCommands(),
		),
		BeforeFunc: core.BeforeFuncCombine(
			generateRSASSHKey("SSHKey"),
			core.ExecStoreBeforeCmd(
				"Server",
				"scw instance server create type=POP2-2C-8G-WIN image=windows_server_2022 admin-password-encryption-ssh-key-id={{.SSHKey.ID}}",
			),
		),
		Cmd: "scw instance server get-rdp-password {{.Server.ID}} --wait",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				assert.NotNil(t, ctx.Result)
				resp := testhelpers.Value[*instance.ServerGetRdpPasswordResponse](t, ctx.Result)

				assert.NotEmpty(t, resp.Password)
			},
		),
		AfterFunc: core.AfterFuncCombine(
			core.ExecAfterCmd(
				"scw instance server terminate {{.Server.ID}} with-ip=true with-block=true",
			),
			core.ExecAfterCmd("scw iam ssh-key delete {{.SSHKey.ID}}"),
		),
		TmpHomeDir: true,
	}))
}
