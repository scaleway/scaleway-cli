package init

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/scaleway/scaleway-cli/internal/core"
	account "github.com/scaleway/scaleway-cli/internal/namespaces/account/v2alpha1"
	accountsdk "github.com/scaleway/scaleway-sdk-go/api/account/v2alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func setUpSSHKeyLocallyWithKeyName(key string, name string) core.BeforeFunc {
	return func(ctx *core.BeforeFuncCtx) error {
		homeDir := ctx.OverrideEnv["HOME"]
		// TODO we persist the key as ~/.ssh/id_rsa.pub regardless of the type of key it is (rsa, ed25519)
		keyPath := path.Join(homeDir, ".ssh", name)
		ctx.Logger.Info("public key path set to: ", keyPath)

		// Ensure the subfolders for the configuration files are all created
		err := os.MkdirAll(filepath.Dir(keyPath), 0755)
		if err != nil {
			return err
		}

		// Write the configuration file
		err = ioutil.WriteFile(keyPath, []byte(key), 0600)
		if err != nil {
			return err
		}
		return nil
	}
}

func removeSSHKeyFromAccount(publicSSHKey string) core.AfterFunc {
	return func(ctx *core.AfterFuncCtx) error {
		api := accountsdk.NewAPI(ctx.Client)
		resp, err := api.ListSSHKeys(&accountsdk.ListSSHKeysRequest{},
			scw.WithAllPages())
		if err != nil {
			return err
		}
		id := ""
		for _, v := range resp.SSHKeys {
			if v.PublicKey == publicSSHKey {
				id = v.ID
			}
		}
		if id != "" {
			err = api.DeleteSSHKey(&accountsdk.DeleteSSHKeyRequest{SSHKeyID: id})
		}
		return err
	}
}

// add an ssh key with a given meta key
func addSSHKeyToAccount(metaKey string, name string, key string) core.BeforeFunc {
	return func(ctx *core.BeforeFuncCtx) error {
		cmd := []string{
			"scw", "account", "ssh-key", "add", "public-key=" + key, "name=" + name,
		}
		ctx.Meta[metaKey] = ctx.ExecuteCmd(cmd)
		return nil
	}
}

func Test_InitSSH(t *testing.T) {
	defaultSettings := map[string]string{
		"secret-key":           "{{ .SecretKey }}",
		"send-telemetry":       "false",
		"remove-v1-config":     "false",
		"install-autocomplete": "false",
	}
	cmds := GetCommands()
	cmds.Merge(account.GetCommands())

	// We create a key in each tests to be able to run those tests in parallel

	t.Run("KeyRegistered", func(t *testing.T) {
		dummySSHKey := "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAICd8ZxAm9mXQsRHhQ5iADEJuO+Ai8EbXMI7TIlsh9jbE foobar@foobar"
		core.Test(&core.TestConfig{
			Commands: cmds,
			BeforeFunc: core.BeforeFuncCombine(
				baseBeforeFunc(),
				setUpSSHKeyLocallyWithKeyName(dummySSHKey, "id_rsa.pub"),
				addSSHKeyToAccount("key", "test-cli-KeyRegistered", dummySSHKey),
			),
			Cmd:        appendArgs("scw init with-ssh-key=true", defaultSettings),
			Check:      core.TestCheckGolden(),
			AfterFunc:  removeSSHKeyFromAccount(dummySSHKey),
			TmpHomeDir: true,
		})(t)
	})

	t.Run("KeyUnregistered", func(t *testing.T) {
		dummySSHKey := "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIIQE67HxSRicWd4ol7ntM2jdeD/qEehPJxK/3thmMiZg foobar@foobar"
		core.Test(&core.TestConfig{
			Commands: cmds,
			BeforeFunc: core.BeforeFuncCombine(
				baseBeforeFunc(),
				setUpSSHKeyLocallyWithKeyName(dummySSHKey, "id_rsa.pub"),
			),
			Cmd:        appendArgs("scw init with-ssh-key=true", defaultSettings),
			Check:      core.TestCheckGolden(),
			TmpHomeDir: true,
			AfterFunc:  removeSSHKeyFromAccount(dummySSHKey),
		})(t)
	})

	t.Run("NoLocalKey", core.Test(&core.TestConfig{
		Commands:   cmds,
		BeforeFunc: baseBeforeFunc(),
		Cmd:        appendArgs("scw init with-ssh-key=true", defaultSettings),
		Check:      core.TestCheckGolden(),
		TmpHomeDir: true,
	}))

	t.Run("WithLocalEd25519Key", func(t *testing.T) {
		dummySSHKey := "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIJb42xh9l10D/u9imDFfLZ+U6KrZmr/F/qBClnmijCFF/qEehPJxK/3thmMiZg foobaz@foobaz"
		core.Test(&core.TestConfig{
			Commands: cmds,
			BeforeFunc: core.BeforeFuncCombine(
				baseBeforeFunc(),
				setUpSSHKeyLocallyWithKeyName(dummySSHKey, "id_ed25519.pub"),
			),
			Cmd:        appendArgs("scw init with-ssh-key=true", defaultSettings),
			Check:      core.TestCheckGolden(),
			TmpHomeDir: true,
			AfterFunc:  removeSSHKeyFromAccount(dummySSHKey),
		})(t)
	})
}
