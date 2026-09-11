package getorcreate

import (
	"context"

	"github.com/scaleway/scaleway-sdk-go/api/container/v1"
	"github.com/scaleway/scaleway-sdk-go/api/registry/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func Namespace(
	ctx context.Context,
	api *container.API,
	region scw.Region,
	name string,
) (*container.Namespace, error) {
	listNamespacesResponse, err := api.ListNamespaces(&container.ListNamespacesRequest{
		Region: region,
		Name:   &name,
	}, scw.WithContext(ctx), scw.WithAllPages())
	if err != nil {
		return nil, err
	}

	namespaces := listNamespacesResponse.Namespaces

	var matchingNamespace *container.Namespace
	for _, namespace := range namespaces {
		if namespace.Name == name {
			matchingNamespace = namespace

			break
		}
	}

	if matchingNamespace != nil {
		return matchingNamespace, nil
	}

	namespace, err := api.CreateNamespace(&container.CreateNamespaceRequest{
		Region: region,
		Name:   name,
	}, scw.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	namespace, err = api.WaitForNamespace(&container.WaitForNamespaceRequest{
		Region:      region,
		NamespaceID: namespace.ID,
	}, scw.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	return namespace, nil
}

func Container(
	ctx context.Context,
	api *container.API,
	region scw.Region,
	namespaceID string,
	name string,
	image string,
	port uint32,
) (*container.Container, error) {
	listContainersResponse, err := api.ListContainers(&container.ListContainersRequest{
		Region:      region,
		NamespaceID: &namespaceID,
		Name:        &name,
	}, scw.WithContext(ctx), scw.WithAllPages())
	if err != nil {
		return nil, err
	}

	containers := listContainersResponse.Containers

	var matchingContainer *container.Container
	for _, c := range containers {
		if c.Name == name {
			matchingContainer = c

			break
		}
	}

	if matchingContainer != nil {
		return matchingContainer, nil
	}

	ctn, err := api.CreateContainer(&container.CreateContainerRequest{
		Region:      region,
		NamespaceID: namespaceID,
		Name:        name,
		Image:       image,
		Port:        &port,
	}, scw.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	return ctn, nil
}

func Registry(
	ctx context.Context,
	api *registry.API,
	region scw.Region,
	name string,
) (*registry.Namespace, error) {
	listNamespacesResponse, err := api.ListNamespaces(&registry.ListNamespacesRequest{
		Region: region,
		Name:   &name,
	}, scw.WithContext(ctx), scw.WithAllPages())
	if err != nil {
		return nil, err
	}

	registries := listNamespacesResponse.Namespaces

	var matchingNamespace *registry.Namespace
	for _, namespace := range registries {
		if namespace.Name == name {
			matchingNamespace = namespace

			break
		}
	}

	if matchingNamespace != nil {
		return matchingNamespace, nil
	}

	namespace, err := api.CreateNamespace(&registry.CreateNamespaceRequest{
		Region: region,
		Name:   name,
	}, scw.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	return namespace, nil
}
