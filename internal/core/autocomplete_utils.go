package core

import (
	"context"
	"reflect"
	"strings"

	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/scaleway/scaleway-sdk-go/strcase"
)

func AutocompleteProfileName() AutoCompleteArgFunc {
	return func(ctx context.Context, prefix string) AutocompleteSuggestions {
		res := AutocompleteSuggestions(nil)
		configPath := ExtractConfigPath(ctx)
		config, err := scw.LoadConfigFromPath(configPath)
		if err != nil {
			return res
		}

		for profileName := range config.Profiles {
			if strings.HasPrefix(profileName, prefix) {
				res = append(res, profileName)
			}
		}

		if strings.HasPrefix(scw.DefaultProfileName, prefix) {
			res = append(res, scw.DefaultProfileName)
		}
		return res
	}
}

// AutocompleteGetArg tries to complete an argument by using the list verb if it exists for the same resource
// It will search for the same field in the response of the list
// Field name will be stripped of the resource name (ex: cluster-id -> id)
func AutocompleteGetArg(ctx context.Context, cmd *Command, argSpec *ArgSpec) []string {
	commands := ExtractCommands(ctx)

	// The argument we want to find (ex: server-id)
	argName := argSpec.Name
	argResource := cmd.Resource

	// if arg name does not start with resource
	// ex with "scw instance private-nic list server-id=<tab>"
	// we get server as resource instead of private-nic to find command "scw instance server list"
	if !strings.HasPrefix(argName, cmd.Resource) {
		dashIndex := strings.Index(argName, "-")
		if dashIndex > 0 {
			argResource = argName[:dashIndex]
		}
	}

	// remove resource from arg name (ex: server-id -> id)
	argName = strings.TrimPrefix(argName, argResource)
	argName = strings.TrimLeft(argName, "-")

	listCmd, hasList := commands.find(cmd.Namespace, argResource, "list")
	if !hasList {
		return nil
	}

	// Build empty arguments and run command
	// Has to use interceptor if it exists as ArgsType could be handled by interceptor
	listCmdArgs := reflect.New(listCmd.ArgsType).Interface()
	if listCmd.Interceptor == nil {
		listCmd.Interceptor = func(ctx context.Context, argsI interface{}, runner CommandRunner) (interface{}, error) {
			return runner(ctx, argsI)
		}
	}
	resp, err := listCmd.Interceptor(ctx, listCmdArgs, listCmd.Run)
	if err != nil {
		return nil
	}

	// As we run the "list" verb instead of using the sdk ListResource, response is already the slice
	// ex: ListServersResponse -> ListServersResponse.Servers
	resources := reflect.ValueOf(resp)
	if resources.Kind() != reflect.Slice {
		return nil
	}
	values := []string(nil)
	// Let's iterate over the struct in the response slice and get the searched field
	for i := 0; i < resources.Len(); i++ {
		resource := resources.Index(i).Elem()
		resourceField := resource.FieldByName(strcase.ToPublicGoName(argName))
		if resourceField.Kind() == reflect.String {
			values = append(values, resourceField.String())
		}
	}
	return values
}
