package terraform

import (
	"github.com/scaleway/scaleway-sdk-go/api/baremetal/v1"
	container "github.com/scaleway/scaleway-sdk-go/api/container/v1beta1"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
)

type Association struct {
	ResourceName string
	ImportFormat string
	SubResources []string
}

const ImportFormatID = "{{ .Region }}/{{ .ID }}"
const ImportFormatZoneID = "{{ .Zone }}/{{ .ID }}"
const ImportFormatRegionID = "{{ .Region }}/{{ .ID }}"

var Associations = map[interface{}]Association{
	&baremetal.Server{}: {
		ResourceName: "scaleway_baremetal_server",
		ImportFormat: ImportFormatZoneID,
	},
	&instance.Server{}: {
		ResourceName: "scaleway_instance_server",
		ImportFormat: ImportFormatZoneID,
	},
	&container.Container{}: {
		ResourceName: "scaleway_container",
		ImportFormat: ImportFormatRegionID,
	},
}
