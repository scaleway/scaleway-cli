package registry

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"

	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-cli/internal/interactive"
	"github.com/scaleway/scaleway-cli/internal/printer"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

const endpointPrefix = "rg."
const endpointSuffix = ".scw.cloud"

type emptyRequest struct{}

type registrySetupDockerHelperArgs struct {
	HelperDirectory string
}

func registrySetupDockerHelperCommand() *core.Command {
	return &core.Command{
		Short: `Setup a local Docker credential helper`,
		Long: `This command will configure the Docker credential helper for your account.

It will create a new script named docker-credential-scw. 
This script will be called each time Docker needs the credentials and will return the correct credentials.
It avoid running docker login commands.
`,
		Namespace: "registry",
		Resource:  "setup-docker-helper",
		ArgsType:  reflect.TypeOf(registrySetupDockerHelperArgs{}),
		ArgSpecs: []*core.ArgSpec{
			{
				Name:  "helper-directory",
				Short: "Directory in which the Docker helper will be installed",
			},
		},
		Run: registrySetupDockerHelperRun,
	}
}

func registrySetupDockerHelperRun(ctx context.Context, argsI interface{}) (i interface{}, e error) {
	// TODO add windows support
	if runtime.GOOS == "windows" {
		return nil, fmt.Errorf("windows is not currently supported")
	}

	binaryName := core.ExtractBinaryName(ctx)

	_, _ = interactive.Println("To enable the Docker credential helper, scw needs to create a script inside of your $PATH.")
	scriptDirArg := argsI.(*registrySetupDockerHelperArgs).HelperDirectory
	if scriptDirArg == "" {
		defaultScriptDirArg := filepath.Join(filepath.Dir(scw.GetConfigPath()), "bin")

		promptedScriptDir, err := interactive.PromptStringWithConfig(&interactive.PromptStringConfig{
			Prompt:          "In which directory do you want to install the script?",
			DefaultValue:    defaultScriptDirArg,
			DefaultValueDoc: defaultScriptDirArg,
		})
		if err != nil {
			return nil, err
		}
		scriptDirArg = promptedScriptDir
	}
	profileFlag := ""
	profileName := core.ExtractProfileName(ctx)
	if profileName != "" {
		profileFlag = fmt.Sprintf(" --profile %s", profileName)
	}

	helperScriptPath := filepath.Join(scriptDirArg, fmt.Sprintf("docker-credential-%s", binaryName))
	helperScript := fmt.Sprintf("#!/bin/sh\n\n%s%s registry docker-helper \"$@\"", binaryName, profileFlag)

	// Warning
	_, _ = interactive.Println()
	_, _ = interactive.PrintlnWithoutIndent("To enable the Docker credential helper we need to create the file " + helperScriptPath + " with the following lines:")
	_, _ = interactive.Println(helperScript)

	// Early exit if user disagrees
	_, _ = interactive.Println()
	continueInstallation, err := interactive.PromptBoolWithConfig(&interactive.PromptBoolConfig{
		Prompt:       fmt.Sprintf("Do you want to proceed with these changes?"),
		DefaultValue: true,
	})
	if err != nil {
		return nil, err
	}
	if !continueInstallation {
		return nil, fmt.Errorf("installation cancelled")
	}

	err = writeHelperScript(helperScriptPath, helperScript)
	if err != nil {
		return nil, fmt.Errorf("failed to write helper script: %s", err)
	}

	var registries []string
	for _, region := range scw.AllRegions {
		registries = append(registries, endpointPrefix+region.String()+endpointSuffix)
	}

	err = setupDockerConfigFile(registries, binaryName)
	if err != nil {
		return nil, fmt.Errorf("failed to write docker config file: %s", err)
	}

	_, _ = interactive.Println()
	_, err = exec.LookPath(fmt.Sprintf("docker-credential-%s", binaryName))
	if err != nil {
		_, _ = interactive.Println(fmt.Sprintf("docker-credential-%s is not present in your $PATH, you should add %s to your $PATH to make it work.", binaryName, path.Dir(helperScriptPath)))
		_, _ = interactive.Println(fmt.Sprintf("You can add it by adding `export PATH=$PATH:%s` to your `.bashrc` or `.zshrc`", path.Dir(helperScriptPath)))
	} else {
		_, _ = interactive.PrintlnWithoutIndent("Docker credential helper successfully installed.")
		_, _ = interactive.PrintlnWithoutIndent("You can now pull/push without logging in to your registries.")
	}

	return nil, nil
}

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
		ArgSpecs: []*core.ArgSpec{
			{},
		},
		Run:         registryDockerHelperGetRun,
		PrinterType: &printer.JSON,
	}
}

func registryDockerHelperGetRun(ctx context.Context, argsI interface{}) (i interface{}, e error) {
	in := bufio.NewReader(os.Stdin)
	var serverURL string
	for {
		s, err := in.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			serverURL = serverURL + s
			break
		}
		serverURL = serverURL + s
	}

	serverURL = strings.TrimSpace(serverURL)

	serverFound := false

	for _, region := range scw.AllRegions {
		if serverURL == endpointPrefix+region.String()+endpointSuffix {
			serverFound = true
			break
		}
	}

	if !serverFound {
		return nil, fmt.Errorf("endpoint %s does not exists", serverURL)
	}

	client := core.ExtractClient(ctx)

	secretKey, ok := client.GetSecretKey()
	if !ok {
		return nil, fmt.Errorf("could not get secret key")
	}

	response := registryDockerHelperGetResponse{
		Username: "scaleway",
		Secret:   secretKey,
	}

	return response, nil
}

func registryDockerHelperStoreCommand() *core.Command {
	return &core.Command{
		Hidden:    true,
		Namespace: "registry",
		Resource:  "docker-helper",
		Verb:      "store",
		ArgsType:  reflect.TypeOf(emptyRequest{}),
		ArgSpecs: []*core.ArgSpec{
			{},
		},
		Run:         registryDockerHelperStoreRun,
		PrinterType: &printer.JSON,
	}
}

func registryDockerHelperStoreRun(ctx context.Context, argsI interface{}) (i interface{}, e error) {
	return nil, nil
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
		Run:         registryDockerHelperEraseRun,
		PrinterType: &printer.JSON,
	}
}

func registryDockerHelperEraseRun(ctx context.Context, argsI interface{}) (i interface{}, e error) {
	return nil, nil
}

func registryDockerHelperListCommand() *core.Command {
	return &core.Command{
		Hidden:    true,
		Namespace: "registry",
		Resource:  "docker-helper",
		Verb:      "list",
		ArgsType:  reflect.TypeOf(emptyRequest{}),
		ArgSpecs: []*core.ArgSpec{
			{},
		},
		Run:         registryDockerHelperListRun,
		PrinterType: &printer.JSON,
	}
}

func registryDockerHelperListRun(ctx context.Context, argsI interface{}) (i interface{}, e error) {
	registryEndpoints := make(map[string]string)
	for _, region := range scw.AllRegions {
		registryEndpoints[endpointPrefix+region.String()+endpointSuffix] = "scaleway"
	}
	return registryEndpoints, nil
}
