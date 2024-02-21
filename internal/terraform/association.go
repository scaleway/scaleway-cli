package terraform

import (
	"reflect"

	"github.com/scaleway/scaleway-sdk-go/api/account/v2"
	"github.com/scaleway/scaleway-sdk-go/api/baremetal/v1"
	container "github.com/scaleway/scaleway-sdk-go/api/container/v1beta1"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type associationParent struct {
	Fetcher      func(client *scw.Client, data interface{}) (interface{}, error)
	AsDataSource bool
}

type associationChild struct {
	// {
	//     [<child attribute>]: <parent attribute>
	// }
	ParentFieldMap map[string]string

	Fetcher func(client *scw.Client, data interface{}) (interface{}, error)
}

type association struct {
	ResourceName string
	ImportFormat string
	Parents      map[string]*associationParent
	Children     []*associationChild
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
		Parents: map[string]*associationParent{
			"namespace_id": {
				Fetcher: func(client *scw.Client, raw interface{}) (interface{}, error) {
					api := container.NewAPI(client)
					data := raw.(*container.Container)

					return api.GetNamespace(&container.GetNamespaceRequest{
						NamespaceID: data.NamespaceID,
						Region:      data.Region,
					})
				},
			},
		},
	},
	&container.Namespace{}: {
		ResourceName: "scaleway_container_namespace",
		ImportFormat: importFormatRegionID,
		Parents: map[string]*associationParent{
			"project_id": {
				AsDataSource: true,
				Fetcher: func(client *scw.Client, raw interface{}) (interface{}, error) {
					api := account.NewAPI(client)
					data := raw.(*container.Namespace)

					return api.GetProject(&account.GetProjectRequest{
						ProjectID: data.ProjectID,
					})
				},
			},
		},
		Children: []*associationChild{
			{
				ParentFieldMap: map[string]string{
					"namespace_id": "id",
				},
				Fetcher: func(client *scw.Client, raw interface{}) (interface{}, error) {
					api := container.NewAPI(client)
					data := raw.(*container.Namespace)

					res, err := api.ListContainers(&container.ListContainersRequest{
						NamespaceID: data.ID,
						Region:      data.Region,
					})
					if err != nil {
						return nil, err
					}

					return res.Containers, nil
				},
			},
		},
	},
	&account.Project{}: {
		ResourceName: "scaleway_account_project",
		ImportFormat: "{{ .ID }}",
		Children: []*associationChild{
			{
				ParentFieldMap: map[string]string{
					"project_id": "id",
				},
				Fetcher: func(client *scw.Client, raw interface{}) (interface{}, error) {
					api := container.NewAPI(client)
					data := raw.(*account.Project)

					res, err := api.ListNamespaces(&container.ListNamespacesRequest{
						ProjectID: &data.ID,
					})
					if err != nil {
						return nil, err
					}

					return res.Namespaces, nil
				},
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
