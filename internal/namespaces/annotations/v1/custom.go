package annotations

import (
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	annotationsSDK "github.com/scaleway/scaleway-sdk-go/api/annotations/v1"
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	human.RegisterMarshalerFunc(annotationsSDK.BindingKey{}, bindingKeyMarshalerFunc)
	human.RegisterMarshalerFunc(annotationsSDK.BindingValue{}, bindingValueMarshalerFunc)

	cmds.MustFind("annotations", "key-value", "list").Override(keyValueListBuilder)

	return cmds
}

func bindingKeyMarshalerFunc(i any, opt *human.MarshalOpt) (string, error) {
	key := i.(annotationsSDK.BindingKey)

	if opt.TableCell {
		return key.Name, nil
	}

	// To avoid recursion of human.Marshal we create a dummy type
	type tmp annotationsSDK.BindingKey

	return human.Marshal(tmp(key), opt)
}

func bindingValueMarshalerFunc(i any, opt *human.MarshalOpt) (string, error) {
	value := i.(annotationsSDK.BindingValue)

	if opt.TableCell {
		return value.Name, nil
	}

	// To avoid recursion of human.Marshal we create a dummy type
	type tmp annotationsSDK.BindingValue

	return human.Marshal(tmp(value), opt)
}
