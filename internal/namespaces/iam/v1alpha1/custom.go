package iam

import (
	"context"
	"os"
	"path"
	"reflect"
	"strings"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	"github.com/scaleway/scaleway-cli/v2/internal/args"
	"github.com/scaleway/scaleway-cli/v2/internal/interactive"
	iam "github.com/scaleway/scaleway-sdk-go/api/iam/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

var logActionMarshalSpecs = human.EnumMarshalSpecs{
	iam.LogActionUnknownAction: &human.EnumMarshalSpec{Attribute: color.Faint},
	iam.LogActionCreated:       &human.EnumMarshalSpec{Attribute: color.FgGreen},
	iam.LogActionUpdated:       &human.EnumMarshalSpec{Attribute: color.FgBlue},
	iam.LogActionDeleted:       &human.EnumMarshalSpec{Attribute: color.FgRed},
}

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.MustFind("iam").Groups = []string{"security"}

	human.RegisterMarshalerFunc(iam.LogAction(""), human.EnumMarshalFunc(logActionMarshalSpecs))

	cmds.Merge(core.NewCommands(
		initWithSSHCommand(),
		iamRuleCreateCommand(),
		iamRuleDeleteCommand(),
	))

	// These commands have an "optional" organization-id that is required for now.
	for _, commandPath := range [][]string{
		{"iam", "group", "list"},
		{"iam", "api-key", "list"},
		{"iam", "ssh-key", "list"},
		{"iam", "user", "list"},
		{"iam", "policy", "list"},
		{"iam", "application", "list"},
	} {
		cmds.MustFind(commandPath...).Override(setOrganizationDefaultValue)
	}

	// Autocomplete permission set names using IAM API.
	cmds.MustFind("iam", "policy", "create").Override(iamPolicyCreateBuilder)
	cmds.MustFind("iam", "policy", "get").Override(iamPolicyGetBuilder)

	cmds.MustFind("iam", "ssh-key", "create").ArgSpecs.GetByName("public-key").CanLoadFile = true

	iamCmd := cmds.MustFind("iam", "api-key", "get")
	iamCmd.ArgsType = iamApiKeyCustomBuilder.argType
	iamCmd.ArgSpecs = iamApiKeyCustomBuilder.argSpecs
	iamCmd.Run = iamApiKeyCustomBuilder.run
	human.RegisterMarshalerFunc(apiKeyResponse{}, apiKeyMarshalerFunc)

	return cmds
}

func initWithSSHCommand() *core.Command {
	return &core.Command{
		Short:     `Initialize SSH key`,
		Long:      `Initialize SSH key.`,
		Namespace: "iam",
		Resource:  "ssh-key",
		Verb:      "init",
		Groups:    []string{"workflow"},
		ArgsType:  reflect.TypeOf(args.RawArgs{}),
		ArgSpecs:  core.ArgSpecs{},
		Run:       InitWithSSHKeyRun,
	}
}

func setOrganizationDefaultValue(c *core.Command) *core.Command {
	c.ArgSpecs.GetByName("organization-id").Default = func(ctx context.Context) (value string, doc string) {
		organizationID := core.GetOrganizationIDFromContext(ctx)

		return organizationID, "<retrieved from config>"
	}

	return c
}

func InitWithSSHKeyRun(ctx context.Context, _ any) (i any, e error) {
	// Get default local SSH key
	var shortenedFilename string
	var err error
	var localSSHKeyContent []byte
	for _, keyName := range [3]string{"id_ecdsa.pub", "id_ed25519.pub", "id_rsa.pub"} {
		// element is the element from someSlice for where we are
		relativePath := path.Join(".ssh", keyName)
		filename := path.Join(core.ExtractUserHomeDir(ctx), relativePath)
		shortenedFilename = "~/" + relativePath
		localSSHKeyContent, err = os.ReadFile(filename)
		// If we managed to load an ssh key, let's stop there
		if err == nil {
			break
		}
	}
	addKeyInstructions := `scw iam ssh-key create name=my-key key="$(cat path/to/my/key.pub)"`
	if err != nil && os.IsNotExist(err) {
		return nil, sshKeyNotFound(shortenedFilename, addKeyInstructions)
	}

	// Get all SSH keys from Scaleway
	client := core.ExtractClient(ctx)
	api := iam.NewAPI(client)
	listSSHKeysResponse, err := api.ListSSHKeys(&iam.ListSSHKeysRequest{}, scw.WithAllPages())
	if err != nil {
		return nil, err
	}

	// Early exit if the SSH key is present locally and on Scaleway
	for _, SSHKey := range listSSHKeysResponse.SSHKeys {
		if strings.TrimSpace(SSHKey.PublicKey) == strings.TrimSpace(string(localSSHKeyContent)) {
			_, _ = interactive.Println(
				"Looks like your local SSH key " + shortenedFilename + " is already present in your Scaleway account.",
			)

			return nil, nil
		}
	}

	// Ask user
	_, _ = interactive.Println(
		"An SSH key is required if you want to connect to a server. More info at https://www.scaleway.com/en/docs/identity-and-access-management/iam/how-to/create-api-keys/",
	)
	addSSHKey, err := interactive.PromptBoolWithConfig(&interactive.PromptBoolConfig{
		Ctx:          ctx,
		Prompt:       "We found an SSH key in " + shortenedFilename + ". Do you want to add it to your Scaleway project?",
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
	_, err = api.CreateSSHKey(&iam.CreateSSHKeyRequest{
		PublicKey: string(localSSHKeyContent),
	})
	if err != nil {
		return nil, err
	}

	return &core.SuccessResult{
		Message: "Key " + shortenedFilename + " successfully added",
	}, nil
}
