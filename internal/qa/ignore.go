package qa

import (
	"reflect"

	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-cli/internal/namespaces"
)

func filterIgnore(unfilteredErrors []error) []error {
	var res []error
	for _, v := range unfilteredErrors {
		if !isIgnoredError(v) {
			res = append(res, v)
		}
	}
	return res
}

func isIgnoredError(err error) bool {
	for _, ignoredError := range ignoredErrors {
		isEqual := false
		switch typedError := err.(type) {
		case *DifferentLocalizationForNamespaceError:
			isEqual = areCommandsEquals(typedError.Command1, ignoredError.Command) || areCommandsEquals(typedError.Command2, ignoredError.Command)
		default:
			isEqual = (reflect.TypeOf(err) == ignoredError.Type) && areCommandsEquals(reflect.ValueOf(err).FieldByName("Command").Interface().(*core.Command), ignoredError.Command)
		}
		if isEqual {
			return true
		}
	}
	return false
}

type ignoredError struct {
	Type    reflect.Type
	Command *core.Command
}

func areCommandsEquals(c1 *core.Command, c2 *core.Command) bool {
	return (c1.Namespace == c2.Namespace) &&
		(c1.Verb == c2.Verb) &&
		(c1.Resource == c2.Resource)
}

var ignoredErrors = []ignoredError{
	{
		Type:    reflect.TypeOf(&DifferentLocalizationForNamespaceError{}),
		Command: namespaces.GetCommands().MustFind("k8s", "kubeconfig", "uninstall"),
	},
	{
		Type:    reflect.TypeOf(&DifferentLocalizationForNamespaceError{}),
		Command: namespaces.GetCommands().MustFind("registry", "logout"),
	},
	{
		Type:    reflect.TypeOf(&DifferentLocalizationForNamespaceError{}),
		Command: namespaces.GetCommands().MustFind("registry", "login"),
	},
}
