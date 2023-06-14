package terraform

import (
	"reflect"

	"github.com/scaleway/scaleway-sdk-go/api/baremetal/v1"
	container "github.com/scaleway/scaleway-sdk-go/api/container/v1beta1"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
)

type associationSubResource struct {
	TerraformAttributeName string
	Command                string
	AsDataSource           bool
}

type association struct {
	ResourceName string
	ImportFormat string
	SubResources map[string]*associationSubResource
}

// const importFormatID = "{{ .Region }}/{{ .ID }}"
const importFormatZoneID = "{{ .Zone }}/{{ .ID }}"
const importFormatRegionID = "{{ .Region }}/{{ .ID }}"

var associations = map[interface{}]*association{
	&baremetal.Server{}: {
		ResourceName: "scaleway_baremetal_server",
		ImportFormat: importFormatZoneID,
	},
	&instance.Server{}: {
		ResourceName: "scaleway_instance_server",
		ImportFormat: importFormatZoneID,
	},
	&container.Container{}: {
		ResourceName: "scaleway_container",
		ImportFormat: importFormatRegionID,
		SubResources: map[string]*associationSubResource{
			"NamespaceID": {
				TerraformAttributeName: "namespace_id",
				Command:                "container namespace get {{ .NamespaceID }}",
			},
		},
	},
	&container.Namespace{}: {
		ResourceName: "scaleway_container_namespace",
		ImportFormat: importFormatRegionID,
		SubResources: map[string]*associationSubResource{
			"ProjectID": {
				TerraformAttributeName: "project_id",
				Command:                "container project get project-id={{ .ProjectID }}",
				AsDataSource:           true,
			},
		},
	},
}

func getAssociation(data interface{}) (*association, bool) {
	dataType := reflect.TypeOf(data)

	for i, association := range associations {
		if dataType == reflect.TypeOf(i) {
			return association, true
		}
	}

	return nil, false
}
