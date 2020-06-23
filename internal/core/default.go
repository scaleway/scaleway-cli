package core

import (
	"context"
	"strings"

	"github.com/scaleway/scaleway-cli/internal/args"
	"github.com/scaleway/scaleway-sdk-go/namegenerator"
)

// ApplyDefaultValues will hydrate args with default values.
func ApplyDefaultValues(ctx context.Context, argSpecs ArgSpecs, rawArgs args.RawArgs) args.RawArgs {
	for _, argSpec := range argSpecs {
		if argSpec.Default == nil {
			continue
		}
		defaultValue, _ := argSpec.Default(ctx)

		// If argSpec is not part of slice or map simply add it to raw args
		if !argSpec.IsPartOfMapOrSlice() {
			if _, exist := rawArgs.Get(argSpec.Name); !exist {
				rawArgs = rawArgs.Add(argSpec.Name, defaultValue)
			}
			continue
		}

		// If argSpec is part of a map or slice we must lookup for existing index in other args
		// Example:
		//    argSpec = { Name: "friends.{index}.Age", "Default": 42 }
		//    rawArgs = friends.0.name=bob friends.1.name=alice
		// In this case we should add friends.0.age=42 friends.1.age=42
		//
		// We will construct a slice prefixes that will contain all args prefixes
		// In the upper example prefix will be [friends.0 and friends.1]
		parts := strings.Split(argSpec.Name, ".")
		prefixes := []string{parts[0]}
		for _, part := range parts[1 : len(parts)-1] {
			switch part {
			case sliceSchema, mapSchema:
				newPrefixes := []string(nil)
				for _, prefix := range prefixes {
					for _, key := range rawArgs.GetSliceOrMapKeys(prefix) {
						newPrefixes = append(newPrefixes, prefix+"."+key)
					}
				}
				prefixes = newPrefixes
			default:
				for idx := range prefixes {
					prefixes[idx] = prefixes[idx] + "." + part
				}
			}
		}

		// Now that we compute the full list of existing prefix we can generate the default value.
		for _, prefix := range prefixes {
			argName := prefix + "." + parts[len(parts)-1]
			if _, exist := rawArgs.Get(argName); !exist {
				rawArgs = rawArgs.Add(argName, defaultValue)
			}
		}
	}

	return rawArgs
}

// GetRandomName returns a random name prefixed for the CLI.
func GetRandomName(prefix string) string {
	return namegenerator.GetRandomName("cli", prefix)
}

func RandomValueGenerator(prefix string) DefaultFunc {
	return func(context.Context) (value string, doc string) {
		return GetRandomName(prefix), "<generated>"
	}
}

func DefaultValueSetter(defaultValue string) DefaultFunc {
	return func(context.Context) (value string, doc string) {
		return defaultValue, defaultValue
	}
}
