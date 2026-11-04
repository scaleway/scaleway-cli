package applesilicon

import (
	"reflect"

	humanv2 "github.com/scaleway/scaleway-cli/v2/core/human/v2"
	applesilicon "github.com/scaleway/scaleway-sdk-go/api/applesilicon/v1alpha1"
)

func init() {
	humanv2.SetLayout(reflect.TypeFor[applesilicon.OS](), humanv2.Layout{
		Sections: []humanv2.Section{
			{
				Path: "CompatibleServerTypes",
			},
		},
	})
}
