package core

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/scaleway/scaleway-cli/v2/internal/args"
	"github.com/scaleway/scaleway-cli/v2/internal/cache"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/scaleway/scaleway-sdk-go/strcase"
)

var autoCompleteCache *cache.Cache

// getGlobalFlags returns the list of flags that should be added to all commands
func getGlobalFlags(ctx context.Context) []FlagSpec {
	printerTypes := []string{
		PrinterTypeHuman.String(),
		PrinterTypeJSON.String(),
		PrinterTypeYAML.String(),
		PrinterTypeTemplate.String(),
	}
	profiles := []string(nil)
	cfg := extractConfig(ctx)
	if cfg != nil {
		for profile := range cfg.Profiles {
			profiles = append(profiles, profile)
		}
	}

	return []FlagSpec{
		{Name: "-c"},
		{Name: "--config"},
		{Name: "-D"},
		{Name: "--debug"},
		{Name: "-h"},
		{Name: "--help"},
		{
			Name:       "-o",
			EnumValues: printerTypes,
		},
		{
			Name:       "--output",
			EnumValues: printerTypes,
		},
		{
			Name:             "-p",
			HasVariableValue: true,
			EnumValues:       profiles,
		},
		{
			Name:             "--profile",
			HasVariableValue: true,
			EnumValues:       profiles,
		},
	}
}

func AutocompleteProfileName() AutoCompleteArgFunc {
	return func(ctx context.Context, prefix string, _ any) AutocompleteSuggestions {
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

		if strings.HasPrefix(scw.DefaultProfileName, prefix) { //nolint:gocritic
			res = append(res, scw.DefaultProfileName)
		}

		return res
	}
}

var autocompleteResourceToNamespace = map[string]string{
	"private-network": "vpc",
}

// AutocompleteGetArg tries to complete an argument by using the list verb if it exists for the same resource
// It will search for the same field in the response of the list
// Field name will be stripped of the resource name (ex: cluster-id -> id)
func AutocompleteGetArg(
	ctx context.Context,
	cmd *Command,
	argSpec *ArgSpec,
	completedArgs map[string]string,
) []string {
	commands := ExtractCommands(ctx)

	// The argument we want to find (ex: server-id)
	argName := argSpec.Name
	argResource := cmd.Resource

	// if arg name does not start with resource
	// ex with "scw instance private-nic list server-id=<tab>"
	// we get server as resource instead of private-nic to find command "scw instance server list"
	if !strings.HasPrefix(argName, cmd.Resource) {
		// Use last index as resource name might contain a dash (ex: security-group-id)
		dashIndex := strings.LastIndex(argName, "-")
		if dashIndex > 0 {
			argResource = argName[:dashIndex]
		}
	}

	// skip if creating a resource and the arg to complete is from the same resource
	// does not complete name in "scw instance server create name=<tab>"
	// but still complete for different resources ex: "scw container container create namespace-id=<tab>"
	if cmd.Verb == "create" && argResource == cmd.Resource {
		return nil
	}

	// remove resource from arg name (ex: server-id -> id)
	argName = strings.TrimPrefix(argName, argResource)
	argName = strings.TrimLeft(argName, "-")

	listCmd, hasList := commands.find(cmd.Namespace, argResource, "list")

	if !hasList {
		namespace, ok := autocompleteResourceToNamespace[argResource]
		if ok {
			listCmd, hasList = commands.find(namespace, argResource, "list")
		}
		if !hasList {
			return nil
		}
	}

	// Build empty arguments and run command
	// Has to use interceptor if it exists as ArgsType could be handled by interceptor
	listCmdArgs := reflect.New(listCmd.ArgsType).Interface()

	// Keep zone and region arguments
	listRawArgs := listRawArgsLocalities(completedArgs, listCmd)

	// Apply default arguments
	listRawArgs = ApplyDefaultValues(ctx, listCmd.ArgSpecs, listRawArgs)

	// Unmarshal args.
	// After that we are done working with rawArgs
	// and will be working with cmdArgs.
	err := args.UnmarshalStruct(listRawArgs, listCmdArgs)
	if err != nil {
		return nil
	}

	if listCmd.Interceptor == nil {
		listCmd.Interceptor = func(ctx context.Context, argsI interface{}, runner CommandRunner) (interface{}, error) {
			return runner(ctx, argsI)
		}
	}

	rawCommand := fmt.Sprintf("%s %s", listCmd.getPath(), strings.Join(listRawArgs, " "))
	resp := autoCompleteCache.Get(rawCommand)
	if resp == nil {
		resp, err = listCmd.Interceptor(ctx, listCmdArgs, listCmd.Run)
		if err != nil {
			return nil
		}
		autoCompleteCache.Set(rawCommand, resp)
	}

	// As we run the "list" verb instead of using the sdk ListResource, response is already the slice
	// ex: ListServersResponse -> ListServersResponse.Servers
	resources := reflect.ValueOf(resp)
	if resources.Kind() != reflect.Slice {
		return nil
	}
	values := []string(nil)
	// Let's iterate over the struct in the response slice and get the searched field
	for i := range resources.Len() {
		resource := resources.Index(i)
		if resource.Kind() == reflect.Ptr {
			resource = resource.Elem()
		}
		resourceField := resource.FieldByName(strcase.ToPublicGoName(argName))
		if resourceField.Kind() == reflect.String {
			values = append(values, resourceField.String())
		}
	}

	return values
}

func listRawArgsLocalities(completedArgs map[string]string, cmd *Command) []string {
	listRawArgs := []string(nil)
	specs := cmd.ArgSpecs
	for arg, value := range completedArgs {
		if strings.HasPrefix(arg, "zone") {
			if specs.GetByName("region") != nil {
				zone, _ := scw.ParseZone(value)
				region, _ := zone.Region()
				listRawArgs = append(listRawArgs, "region="+region.String())
			}
			if specs.GetByName("zone") != nil {
				listRawArgs = append(listRawArgs, arg+value)
			}
		}
		if strings.HasPrefix(arg, "region") && specs.GetByName("region") != nil {
			listRawArgs = append(listRawArgs, arg+value)
		}
	}

	return listRawArgs
}

func flattenCompleteArgs(completedArgs map[string]string) []string {
	rawArgs := make([]string, 0, len(completedArgs))
	for arg, value := range completedArgs {
		rawArgs = append(rawArgs, arg+value) // "region=" + "nl-ams"
	}

	return rawArgs
}

func requestFromCompletedArgs(cmd *Command, completedArgs map[string]string) any {
	request := reflect.New(cmd.ArgsType).Interface()
	_ = args.UnmarshalStruct(flattenCompleteArgs(completedArgs), request)

	return request
}
