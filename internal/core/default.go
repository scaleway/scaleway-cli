package core

import (
	"context"
	"strings"

	"github.com/scaleway/scaleway-cli/internal/args"
	"github.com/scaleway/scaleway-sdk-go/namegenerator"
)

// ApplyDefaultValues will hydrate args with default values.
func ApplyDefaultValues(ctx context.Context, argSpecs ArgSpecs, rawArgs []string) []string {
	argsMap := args.SplitRawMap(rawArgs)

	// following the format ["arg1=1", "arg2=2", "arg3"]
	additionalArgs := []string{}
	for _, arg := range argSpecs {
		// TODO: Handle nested, map and slices
		if strings.Contains(arg.Name, ".") {
			continue
		}

		_, exist := argsMap[arg.Name]
		if exist {
			continue
		}

		if arg.Default != nil {
			defaultValue, _ := arg.Default()
			additionalArgs = append(additionalArgs, arg.Name+"="+defaultValue)
		}
	}

	return append(rawArgs, additionalArgs...)
}

// GetRandomName returns a random name prefixed for the CLI.
func GetRandomName(prefix string) string {
	return namegenerator.GetRandomName("cli", prefix)
}

func RandomValueGenerator(prefix string) DefaultFunc {
	return func() (value string, doc string) {
		return GetRandomName(prefix), "<generated>"
	}
}

func DefaultValueSetter(defaultValue string) DefaultFunc {
	return func() (value string, doc string) {
		return defaultValue, defaultValue
	}
}
