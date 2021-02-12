package account

import (
	"context"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"strings"

	"github.com/scaleway/scaleway-cli/internal/args"
	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-cli/internal/interactive"
	account "github.com/scaleway/scaleway-sdk-go/api/account/v2alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func GetCommands() *core.Commands {
	commands := GetGeneratedCommands()

	commands.Merge(core.NewCommands(
		initCommand(),
	))
	return commands
}

func initCommand() *core.Command {
	return &core.Command{
		Short:     `Initialize SSH key`,
		Long:      `Initialize SSH key.`,
		Namespace: "account",
		Resource:  "ssh-key",
		Verb:      "init",
		ArgsType:  reflect.TypeOf(args.RawArgs{}),
		ArgSpecs:  core.ArgSpecs{},
		Run:       InitRun,
	}
}

func InitRun(ctx context.Context, argsI interface{}) (i interface{}, e error) {
	// Get default local SSH key
	var shortenedFilename string
	var err error
	var localSSHKeyContent []byte
	for _, keyName := range [3]string{"id_ecdsa.pub", "id_ed25519.pub", "id_rsa.pub"} {
		// element is the element from someSlice for where we are
		relativePath := path.Join(".ssh", keyName)
		filename := path.Join(core.ExtractUserHomeDir(ctx), relativePath)
		shortenedFilename = "~/" + relativePath
		localSSHKeyContent, err = ioutil.ReadFile(filename)
		// If we managed to load an ssh key, let's stop there
		if err == nil {
			break
		}
	}
	addKeyInstructions := `scw account ssh-key add name=my-key key="$(cat path/to/my/key.pub)"`
	if err != nil && os.IsNotExist(err) {
		return nil, sshKeyNotFound(shortenedFilename, addKeyInstructions)
	}

	// Get all SSH keys from Scaleway
	client := core.ExtractClient(ctx)
	api := account.NewAPI(client)
	listSSHKeysResponse, err := api.ListSSHKeys(&account.ListSSHKeysRequest{}, scw.WithAllPages())
	if err != nil {
		return nil, err
	}

	// Early exit if the SSH key is present locally and on Scaleway
	for _, SSHKey := range listSSHKeysResponse.SSHKeys {
		if strings.TrimSpace(SSHKey.PublicKey) == strings.TrimSpace(string(localSSHKeyContent)) {
			_, _ = interactive.Println("Looks like your local SSH key " + shortenedFilename + " is already present in your Scaleway account.")
			return nil, nil
		}
	}

	// Ask user
	_, _ = interactive.Println("An SSH key is required if you want to connect to a server. More info at https://www.scaleway.com/en/docs/configure-new-ssh-key")
	addSSHKey, err := interactive.PromptBoolWithConfig(&interactive.PromptBoolConfig{
		Ctx:          ctx,
		Prompt:       "We found an SSH key in " + shortenedFilename + ". Do you want to add it to your Scaleway account?",
		DefaultValue: true,
	})
	if err != nil {
		return nil, err
	}

	// Early exit if user doesn't want to add the key
	if !addSSHKey {
		return nil, installationCanceled(addKeyInstructions)
	}

	// Add key
	_, err = api.CreateSSHKey(&account.CreateSSHKeyRequest{
		PublicKey: string(localSSHKeyContent),
	})
	if err != nil {
		return nil, err
	}

	return &core.SuccessResult{
		Message: "Key " + shortenedFilename + " successfully added",
	}, nil
}
