package registry

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"text/template"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/interactive"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type emptyRequest struct{}

type registrySetupDockerHelperArgs struct {
	Path string
}

func registryInstallDockerHelperCommand() *core.Command {
	return &core.Command{
		Short: `Install a local Docker credential helper`,
		Long: `This command will install the Docker credential helper for your account.

It will create a new script named docker-credential-scw. 
This script will be called each time Docker needs the credentials and will return the correct credentials.
It avoid running docker login commands.
`,
		Namespace: "registry",
		Resource:  "install-docker-helper",
		ArgsType:  reflect.TypeOf(registrySetupDockerHelperArgs{}),
		ArgSpecs: []*core.ArgSpec{
			{
				Name:    "path",
				Short:   "Directory in which the Docker helper will be installed. This directory should be in your $PATH",
				Default: core.DefaultValueSetter("/usr/local/bin"),
				ValidateFunc: func(_ *core.ArgSpec, value interface{}) error {
					stat, err := os.Stat(value.(string))
					if err != nil || !stat.IsDir() {
						return fmt.Errorf("%s is not a directory", value)
					}

					return nil
				},
			},
		},
		Run: registrySetupDockerHelperRun,
	}
}

func registrySetupDockerHelperRun(ctx context.Context, argsI interface{}) (i interface{}, e error) {
	// TODO add windows support
	if runtime.GOOS == "windows" {
		return nil, core.WindowIsNotSupportedError()
	}

	binaryName := core.ExtractBinaryName(ctx)
	scriptDirArg := argsI.(*registrySetupDockerHelperArgs).Path

	tpl, err := template.New("script").Parse(helperScriptTemplate)
	if err != nil {
		return nil, err
	}
	helperScriptPath := filepath.Join(scriptDirArg, "docker-credential-"+binaryName)
	buf := bytes.Buffer{}
	tplData := map[string]string{
		"BinaryName": binaryName,
	}

	if profileName := core.ExtractProfileName(ctx); profileName != scw.DefaultProfileName {
		tplData["ProfileName"] = core.ExtractProfileName(ctx)
	}
	err = tpl.Execute(&buf, tplData)
	if err != nil {
		return nil, err
	}
	helperScriptContent := buf.String()

	// Warning
	_, _ = interactive.Println(
		"To enable the Docker credential helper we need to create the file " + helperScriptPath + " with the following lines:\n",
	)
	_, _ = interactive.Println(helperScriptContent)

	// Early exit if user disagrees
	_, _ = interactive.Println()
	continueInstallation, err := interactive.PromptBoolWithConfig(&interactive.PromptBoolConfig{
		Ctx:          ctx,
		Prompt:       "Do you want to proceed with these changes?",
		DefaultValue: true,
	})
	if err != nil {
		return nil, err
	}
	if !continueInstallation {
		return nil, errors.New("installation cancelled")
	}

	err = writeHelperScript(helperScriptPath, helperScriptContent)
	if err != nil {
		return nil, fmt.Errorf("failed to write helper script: %s", err)
	}

	registries := make([]string, 0, len(scw.AllRegions))
	for _, region := range scw.AllRegions {
		registries = append(registries, getRegistryEndpoint(region))
	}

	err = setupDockerConfigFile(ctx, registries, binaryName)
	if err != nil {
		return nil, fmt.Errorf("failed to write docker config file: %s", err)
	}

	_, _ = interactive.Println()
	_, err = exec.LookPath("docker-credential-" + binaryName)
	if err != nil {
		_, _ = interactive.Println(
			fmt.Sprintf(
				"docker-credential-%s is not present in your $PATH, you should add %s to your $PATH to make it work.",
				binaryName,
				path.Dir(helperScriptPath),
			),
		)
		_, _ = interactive.Println(
			fmt.Sprintf(
				"You can add it by adding `export PATH=$PATH:%s` to your `.bashrc`, `.fishrc` or `.zshrc`",
				path.Dir(helperScriptPath),
			),
		)
	} else {
		_, _ = interactive.PrintlnWithoutIndent("Docker credential helper successfully installed.")
		_, _ = interactive.PrintlnWithoutIndent("The Docker credential helper will now take care of the authentication for you.")
		_, _ = interactive.PrintlnWithoutIndent("You don't have to login to your registries anymore.")
	}

	return &core.SuccessResult{}, nil
}

const helperScriptTemplate = `#!/bin/sh
{{ if .ProfileName -}}
PROFILE_NAME="{{ .ProfileName }}"
if [[ ! -z "$SCW_PROFILE" ]]
then 
	PROFILE_NAME="$SCW_PROFILE"
fi
{{ end -}}
{{ .BinaryName }}{{ if .ProfileName }} --profile $PROFILE_NAME{{ end }} registry docker-helper "$@"
`

type registryDockerHelperGetResponse struct {
	Secret   string `json:"Secret"`
	Username string `json:"Username"`
}

func registryDockerHelperGetCommand() *core.Command {
	return &core.Command{
		Hidden:    true,
		Namespace: "registry",
		Resource:  "docker-helper",
		Verb:      "get",
		ArgsType:  reflect.TypeOf(emptyRequest{}),
		Run:       registryDockerHelperGetRun,
	}
}

func registryDockerHelperGetRun(ctx context.Context, _ interface{}) (i interface{}, e error) {
	var serverURL string
	serverURL, err := bufio.NewReader(core.ExtractStdin(ctx)).ReadString('\n')
	if err != nil && err != io.EOF {
		return nil, err
	}

	serverURL = strings.TrimRight(serverURL, "\n")

	serverFound := false
	for _, region := range scw.AllRegions {
		if serverURL == getRegistryEndpoint(region) {
			serverFound = true

			break
		}
	}

	if !serverFound {
		return nil, fmt.Errorf("endpoint %s does not exist", serverURL)
	}

	client := core.ExtractClient(ctx)

	secretKey, ok := client.GetSecretKey()
	if !ok {
		return nil, errors.New("could not get secret key")
	}

	raw, err := json.Marshal(registryDockerHelperGetResponse{
		Username: "scaleway",
		Secret:   secretKey,
	})
	if err != nil {
		return nil, err
	}

	return core.RawResult(raw), nil
}

func registryDockerHelperStoreCommand() *core.Command {
	return &core.Command{
		Hidden:    true,
		Namespace: "registry",
		Resource:  "docker-helper",
		Verb:      "store",
		ArgsType:  reflect.TypeOf(emptyRequest{}),
		Run: func(_ context.Context, _ interface{}) (i interface{}, e error) {
			return nil, nil
		},
	}
}

func registryDockerHelperEraseCommand() *core.Command {
	return &core.Command{
		Hidden:    true,
		Namespace: "registry",
		Resource:  "docker-helper",
		Verb:      "erase",
		ArgsType:  reflect.TypeOf(emptyRequest{}),
		ArgSpecs: []*core.ArgSpec{
			{},
		},
		Run: func(_ context.Context, _ interface{}) (i interface{}, e error) {
			return nil, nil
		},
	}
}

func registryDockerHelperListCommand() *core.Command {
	return &core.Command{
		Hidden:    true,
		Namespace: "registry",
		Resource:  "docker-helper",
		Verb:      "list",
		ArgsType:  reflect.TypeOf(emptyRequest{}),
		Run: func(_ context.Context, _ interface{}) (i interface{}, e error) {
			registryEndpoints := make(map[string]string)
			for _, region := range scw.AllRegions {
				registryEndpoints[getRegistryEndpoint(region)] = "scaleway"
			}
			raw, err := json.Marshal(registryEndpoints)
			if err != nil {
				return nil, err
			}

			return core.RawResult(raw), nil
		},
	}
}
