package sshkey

import (
	"context"
	"io/ioutil"
	"os"
	"path"
	"reflect"

	"github.com/scaleway/scaleway-sdk-go/scw"

	account "github.com/scaleway/scaleway-sdk-go/api/account/v2alpha1"

	"github.com/scaleway/scaleway-cli/internal/args"
	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-cli/internal/interactive"
)

func GetCommands() *core.Commands {
	return core.NewCommands(
		sshKeyCommand(),
		initCommand(),
	)
}

func sshKeyCommand() *core.Command {
	return &core.Command{
		Short:     `Handle SHH key`,
		Long:      `Handle SHH key.`,
		Namespace: "ssh-key",
		NoClient:  true,
		ArgsType:  reflect.TypeOf(args.RawArgs{}),
	}
}

func initCommand() *core.Command {
	return &core.Command{
		Short:     `Initiliaze SHH key`,
		Long:      `Initiliaze SHH key.`,
		Namespace: "ssh-key",
		Resource:  "init",
		ArgsType:  reflect.TypeOf(args.RawArgs{}),
		Run:       InitRun,
	}
}

func InitRun(ctx context.Context, argsI interface{}) (i interface{}, e error) {
	// Explanation
	_, _ = interactive.Println("An SSH key is required if you want to connect to a server. More info at https://www.scaleway.com/en/docs/configure-new-ssh-key/")

	// Get all SSH keys from Scaleway
	client := core.ExtractClient(ctx)
	api := account.NewAPI(client)
	listSSHKeysResponse, err := api.ListSSHKeys(&account.ListSSHKeysRequest{}, scw.WithAllPages())
	if err != nil {
		return nil, err
	}

	// Get default SSH key locally
	relativePath := ".ssh/id_rsa.pub"
	filename := path.Join(os.Getenv("HOME"), relativePath)
	shortenedFilename := "~/" + relativePath
	localSHHKeyContent, err := ioutil.ReadFile(filename)

	addKeyInstructions := `scw account ssh-key add name=my-key key="($)(cat path/to/my/key.pub)"`

	// Early exit if key is not present locally
	if os.IsNotExist(err) {
		return "We did not find an ssh key at " + shortenedFilename + "\n" +
			"You can add one later using:\n" +
			"  " + addKeyInstructions, nil
	} else if err != nil {
		return nil, err
	}

	// Early exit if the SSH key is present locally and on Scaleway
	for _, SSHKey := range listSSHKeysResponse.SSHKeys {
		if SSHKey.PublicKey == string(localSHHKeyContent) {
			return "✅ Key " + shortenedFilename + " is already present on your scaleway account", nil
		}
	}

	// Ask user
	addSHHKey, err := interactive.PromptBoolWithConfig(&interactive.PromptBoolConfig{
		Prompt:       "We found an SSH key in " + shortenedFilename + ". Do you want to add it to your Scaleway account ?",
		DefaultValue: true,
	})
	if err != nil {
		return nil, err
	}

	// Early exit if user doesn't want to add the key
	if !addSHHKey {
		return "You can add it later using:\n" +
			"  " + addKeyInstructions, nil
	}

	// Add key
	_, err = api.CreateSSHKey(&account.CreateSSHKeyRequest{
		PublicKey: string(localSHHKeyContent),
	})
	if err != nil {
		return nil, err
	}

	return "✅ Key " + shortenedFilename + " successfully added", nil
}
