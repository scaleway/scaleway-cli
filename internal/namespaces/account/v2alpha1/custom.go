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
	cmds := GetGeneratedCommands()
	cmds.Merge(core.NewCommands(
		sshKeyCommand(),
		initCommand(),
	))
	cmds.MustFind("account", "ssh-key", "create").Override(addSSHKeyCommand)
	cmds.MustFind("account", "ssh-key", "delete").Override(removeSSHKeyCommand)

	return cmds
}

func addSSHKeyCommand(c *core.Command) *core.Command {
	c.Short = "Add a SSH key to your Scaleway account"
	c.Long = "Add a SSH key to your Scaleway account"
	c.Verb = "add"
	c.Examples = []*core.Example{
		{
			Short:   "Add a given ssh key",
			Request: `{"name": "foobar", "public-key": "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQC+VZDiTJwQGqKmzx1NYduvxFNi+jw7X2SdG5DpTUkFNuPXfAmYzWOeF/iRe2YO5bWy95bQTsgoh44Ed55YF13a1i75HbIpPhQmUEQOb6MeYFHf6Tgg5KfII5Y7oacunlYsjLffGZH3Glxy6jyg/Sx8XOP4lXqMWZapDbnCzAZ15S8jnZBnlWNo1Z60gnX0QQnyiOPFJj+gDtjZG05qCK8Gdh2WlHpKJGB1tTRwYmqStKnpb3aBmzvU9D8eys2XlaJl0DVj63RXv1ej2+Xm0TKL/P7bw+gs+f+GKu5gjI8Sr50G/V8E4svjzlZfFnZcQzVt181QInO/YRWNSmMa1lv/wIPBGdCkWWNgDXieZojkLSe33KxmHYLaNvCoRP700pAHx/5X++nyylleMCRGZfJAGENPwdKNCiT/xqOOdpgbbwJT/baMu3qe83B5qPVFkMw4bI0htwF+gcbABAH/PJD95VFjb1dLf84rwq2pQcNXulJfR3bUqvZ1udx16sgY03s= foobar@foobar"}`,
		},
	}
	c.SeeAlsos = []*core.SeeAlso{
		{
			Short:   "Remove an SSH key",
			Command: "scw account ssh-key remove",
		},
		{
			Short:   "List all SSH keys",
			Command: "scw account ssh-key list",
		},
	}
	return c
}

func removeSSHKeyCommand(c *core.Command) *core.Command {
	c.Short = "Remove a SSH key from your Scaleway account"
	c.Long = "Remove a SSH key from your Scaleway account"
	c.Verb = "remove"
	c.Examples = []*core.Example{
		{
			Short:   "Remove a given ssh key",
			Request: `{"ssh-key-id": "11111111-1111-1111-1111-111111111111"}`,
		},
	}
	c.SeeAlsos = []*core.SeeAlso{
		{
			Short:   "Add an SSH key",
			Command: "scw account ssh-key add",
		},
		{
			Short:   "List all SSH keys",
			Command: "scw account ssh-key list",
		},
	}
	return c
}

func sshKeyCommand() *core.Command {
	return &core.Command{
		Short:     `Manage SSH key`,
		Long:      `Manage SSH key.`,
		Namespace: "account",
		Resource:  "ssh-key",
		NoClient:  true,
		ArgsType:  reflect.TypeOf(args.RawArgs{}),
	}
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
	relativePath := path.Join(".ssh", "id_rsa.pub")
	filename := path.Join(core.ExtractUserHomeDir(ctx), relativePath)
	shortenedFilename := "~/" + relativePath
	localSSHKeyContent, err := ioutil.ReadFile(filename)

	addKeyInstructions := `scw account ssh-key add name=my-key key="$(cat path/to/my/key.pub)"`

	// Early exit if key is not present locally
	if os.IsNotExist(err) {
		return nil, sshKeyNotFound(filename, addKeyInstructions)
	} else if err != nil {
		return nil, err
	}

	// Get all SSH keys from Scaleway
	client, err := core.ExtractClientOrCreate(ctx)
	if err != nil {
		return nil, err
	}
	api := account.NewAPI(client)
	listSSHKeysResponse, err := api.ListSSHKeys(&account.ListSSHKeysRequest{}, scw.WithAllPages())
	if err != nil {
		return nil, err
	}

	// Early exit if the SSH key is present locally and on Scaleway
	for _, SSHKey := range listSSHKeysResponse.SSHKeys {
		if strings.TrimSpace(SSHKey.PublicKey) == strings.TrimSpace(string(localSSHKeyContent)) {
			_, _ = interactive.Println("Look like your local SSH key " + shortenedFilename + " is already present on your Scaleway account.")
			return nil, nil
		}
	}

	// Ask user
	_, _ = interactive.Println("An SSH key is required if you want to connect to a server. More info at https://www.scaleway.com/en/docs/configure-new-ssh-key")
	addSSHKey, err := interactive.PromptBoolWithConfig(&interactive.PromptBoolConfig{
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
