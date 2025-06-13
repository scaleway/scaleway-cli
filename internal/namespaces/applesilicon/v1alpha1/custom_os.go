package applesilicon

import (
	"github.com/scaleway/scaleway-cli/v2/core/human"
	applesilicon "github.com/scaleway/scaleway-sdk-go/api/applesilicon/v1alpha1"
)

func OSMarshalerFunc(i any, opt *human.MarshalOpt) (string, error) {
	type tmp applesilicon.OS
	os := tmp(i.(applesilicon.OS))
	opt.Sections = []*human.MarshalSection{
		{
			FieldName: "CompatibleServerTypes",
			Title:     "CompatibleServerTypes",
		},
	}
	str, err := human.Marshal(os, opt)
	if err != nil {
		return "", err
	}

	return str, nil
}
