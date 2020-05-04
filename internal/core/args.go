package core

import "strings"

type RawArgs []string

func (args RawArgs) GetPositionalArgs() []string {

	positionalArgs := []string(nil)
	for _, arg := range args {
		if isPositionalArg(arg) {
			positionalArgs = append(positionalArgs, arg)
		}
	}
	return positionalArgs
}

func (args RawArgs) Get(argName string) (string, bool) {
	for _, arg := range args {
		name, value := splitArg(arg)
		if name == argName {
			return value, true
		}
	}
	return "", false
}

func (args RawArgs) RemoveAllPositional() RawArgs {
	return args.filter(func(arg string) bool {
		return !isPositionalArg(arg)
	})
}

func (args RawArgs) Add(name string, value string) RawArgs {
	return append(args, name+"="+value)
}

func (args RawArgs) Remove(argName string) RawArgs {
	return args.filter(func(arg string) bool {
		name, _ := splitArg(arg)
		return name != argName
	})
}

func (args RawArgs) filter(test func(string) bool) RawArgs {
	argsCopy := RawArgs{}
	for _, arg := range args {
		if test(arg) {
			argsCopy = append(argsCopy, arg)
		}
	}
	return argsCopy
}

func (args RawArgs) GetSliceOrMapKeys(prefix string) []string {
	keys := []string(nil)
	for _, arg := range args {
		name, _ := splitArg(arg)
		if !strings.HasPrefix(name, prefix+".") {
			continue
		}

		name = strings.TrimPrefix(name, prefix+".")
		keys = append(keys, strings.SplitN(name, ".", 2)[0])
	}
	return keys
}

func splitArg(arg string) (name string, value string) {
	part := strings.SplitN(arg, "=", 2)
	if len(part) == 1 {
		return "", part[0]
	}
	return part[0], part[1]
}

func isPositionalArg(arg string) bool {
	pos := strings.IndexRune(arg, '=')
	return pos == -1
}
