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
			isEqual = (reflect.TypeOf(err) == ignoredError.Type) && areCommandsEquals(reflect.ValueOf(err).FieldByName("Commmand").Interface().(*core.Command), ignoredError.Command)
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
		Command: namespaces.GetCommands().MustFind("k8s", "cluster", "list"),
	},
	{
		Type:    reflect.TypeOf(&DifferentLocalizationForNamespaceError{}),
		Command: namespaces.GetCommands().MustFind("k8s", "cluster", "create"),
	},
	{
		Type:    reflect.TypeOf(&DifferentLocalizationForNamespaceError{}),
		Command: namespaces.GetCommands().MustFind("k8s", "cluster", "get"),
	},
	{
		Type:    reflect.TypeOf(&DifferentLocalizationForNamespaceError{}),
		Command: namespaces.GetCommands().MustFind("k8s", "cluster", "update"),
	},
	{
		Type:    reflect.TypeOf(&DifferentLocalizationForNamespaceError{}),
		Command: namespaces.GetCommands().MustFind("k8s", "cluster", "delete"),
	},
	{
		Type:    reflect.TypeOf(&DifferentLocalizationForNamespaceError{}),
		Command: namespaces.GetCommands().MustFind("k8s", "cluster", "upgrade"),
	},
	{
		Type:    reflect.TypeOf(&DifferentLocalizationForNamespaceError{}),
		Command: namespaces.GetCommands().MustFind("k8s", "cluster", "list-available-versions"),
	},
	{
		Type:    reflect.TypeOf(&DifferentLocalizationForNamespaceError{}),
		Command: namespaces.GetCommands().MustFind("k8s", "cluster", "reset-admin-token"),
	},
	{
		Type:    reflect.TypeOf(&DifferentLocalizationForNamespaceError{}),
		Command: namespaces.GetCommands().MustFind("k8s", "pool", "list"),
	},
	{
		Type:    reflect.TypeOf(&DifferentLocalizationForNamespaceError{}),
		Command: namespaces.GetCommands().MustFind("k8s", "pool", "create"),
	},
	{
		Type:    reflect.TypeOf(&DifferentLocalizationForNamespaceError{}),
		Command: namespaces.GetCommands().MustFind("k8s", "pool", "get"),
	},
	{
		Type:    reflect.TypeOf(&DifferentLocalizationForNamespaceError{}),
		Command: namespaces.GetCommands().MustFind("k8s", "pool", "upgrade"),
	},
	{
		Type:    reflect.TypeOf(&DifferentLocalizationForNamespaceError{}),
		Command: namespaces.GetCommands().MustFind("k8s", "pool", "update"),
	},
	{
		Type:    reflect.TypeOf(&DifferentLocalizationForNamespaceError{}),
		Command: namespaces.GetCommands().MustFind("k8s", "pool", "delete"),
	},
	{
		Type:    reflect.TypeOf(&DifferentLocalizationForNamespaceError{}),
		Command: namespaces.GetCommands().MustFind("k8s", "node", "list"),
	},
	{
		Type:    reflect.TypeOf(&DifferentLocalizationForNamespaceError{}),
		Command: namespaces.GetCommands().MustFind("k8s", "node", "get"),
	},
	{
		Type:    reflect.TypeOf(&DifferentLocalizationForNamespaceError{}),
		Command: namespaces.GetCommands().MustFind("k8s", "node", "replace"),
	},
	{
		Type:    reflect.TypeOf(&DifferentLocalizationForNamespaceError{}),
		Command: namespaces.GetCommands().MustFind("k8s", "node", "reboot"),
	},
	{
		Type:    reflect.TypeOf(&DifferentLocalizationForNamespaceError{}),
		Command: namespaces.GetCommands().MustFind("k8s", "version", "list"),
	},
	{
		Type:    reflect.TypeOf(&DifferentLocalizationForNamespaceError{}),
		Command: namespaces.GetCommands().MustFind("k8s", "version", "get"),
	},
	{
		Type:    reflect.TypeOf(&DifferentLocalizationForNamespaceError{}),
		Command: namespaces.GetCommands().MustFind("k8s", "kubeconfig", "get"),
	},
	{
		Type:    reflect.TypeOf(&DifferentLocalizationForNamespaceError{}),
		Command: namespaces.GetCommands().MustFind("k8s", "kubeconfig", "install"),
	},
	{
		Type:    reflect.TypeOf(&DifferentLocalizationForNamespaceError{}),
		Command: namespaces.GetCommands().MustFind("k8s", "kubeconfig", "uninstall"),
	},
}
