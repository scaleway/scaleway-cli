package inference

import (
	"github.com/scaleway/scaleway-cli/v2/core/human"
	inference "github.com/scaleway/scaleway-sdk-go/api/inference/v1beta1"
)

func ListNodeTypeMarshaler(i any, opt *human.MarshalOpt) (string, error) {
	type tmp []*inference.NodeType
	node := tmp(i.([]*inference.NodeType))

	opt.Fields = []*human.MarshalFieldOpt{
		{
			FieldName: "Name",
			Label:     "Name",
		},
		{
			FieldName: "StockStatus",
			Label:     "Stock Status",
		},
		{
			FieldName: "Description",
			Label:     "Description",
		},
		{
			FieldName: "Vcpus",
			Label:     "VCPUs",
		},
		{
			FieldName: "Memory",
			Label:     "Memory",
		},
		{
			FieldName: "Vram",
			Label:     "VRAM",
		},
	}
	str, err := human.Marshal(node, opt)
	if err != nil {
		return "", err
	}

	return str, nil
}
