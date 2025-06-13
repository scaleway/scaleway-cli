package inference

import (
	"github.com/scaleway/scaleway-cli/v2/core/human"
	inference "github.com/scaleway/scaleway-sdk-go/api/inference/v1beta1"
)

func ListModelMarshalerFunc(i any, opt *human.MarshalOpt) (string, error) {
	type tmp []*inference.Model
	model := tmp(i.([]*inference.Model))
	opt.Fields = []*human.MarshalFieldOpt{
		{
			FieldName: "ID",
			Label:     "ID",
		},
		{
			FieldName: "Name",
			Label:     "Name",
		},
		{
			FieldName: "Provider",
			Label:     "Provider",
		},
		{
			FieldName: "Tags",
			Label:     "Tags",
		},
	}
	str, err := human.Marshal(model, opt)
	if err != nil {
		return "", err
	}

	return str, nil
}
