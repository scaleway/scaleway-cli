package container

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"time"

	dockerTypes "github.com/docker/docker/api/types"
	docker "github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/scaleway/scaleway-cli/v2/internal/core"
	"github.com/scaleway/scaleway-cli/v2/internal/tasks"
	container "github.com/scaleway/scaleway-sdk-go/api/container/v1beta1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type containerDeployRequest struct {
	Region scw.Region

	Name string
}

func containerDeployCommand() *core.Command {
	return &core.Command{
		Short:     `Deploy a container`,
		Long:      `Automatically build and deploy a container.`,
		Namespace: "container",
		Resource:  "deploy",
		Groups:    []string{"workflow"},
		ArgsType:  reflect.TypeOf(containerDeployRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:  "name",
				Short: "Name of the application",
				Default: func(ctx context.Context) (value string, doc string) {
					currentDirection, err := os.Getwd()
					if err != nil {
						return "", ""
					}

					name := filepath.Base(currentDirection)
					if name == "." {
						return "", ""
					}

					name = "app-" + name
					return name, name
				},
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw, scw.Region(core.AllLocalities)),
		},
		Run: containerDeployRun,
	}
}

func containerDeployRun(ctx context.Context, argsI interface{}) (i interface{}, e error) {
	args := argsI.(*containerDeployRequest)

	client := core.ExtractClient(ctx)
	api := container.NewAPI(client)

	fileInfo, err := os.Stat("Dockerfile")
	if err != nil {
		return nil, err
	}

	if fileInfo.IsDir() {
		return nil, fmt.Errorf("'Dockerfile' is a directory")
	}

	cli, err := docker.NewClientWithOpts(docker.FromEnv, docker.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("could not connect to Docker: %v", err)
	}

	actions := tasks.Begin()

	tasks.Add(actions, "Creating namespace", func(t *tasks.Task, _ interface{}) (*container.Namespace, error) {
		namespace, err := getOrCreateNamespace(t.Ctx, api, args.Region, args.Name)
		if err != nil {
			return nil, err
		}

		return namespace, nil
	})

	type packingImageResponse struct {
		namespace *container.Namespace
		tar       io.ReadCloser
	}
	tasks.Add(actions, "Packing image", func(t *tasks.Task, namespace *container.Namespace) (*packingImageResponse, error) {
		tar, err := archive.TarWithOptions(".", &archive.TarOptions{})
		if err != nil {
			return nil, fmt.Errorf("could not create tar: %v", err)
		}

		return &packingImageResponse{
			namespace: namespace,
			tar:       tar,
		}, nil
	})

	type buildImageResponse struct {
		namespace *container.Namespace
		tag       string
	}
	tasks.Add(actions, "Building image", func(t *tasks.Task, packing *packingImageResponse) (*buildImageResponse, error) {
		tag := packing.namespace.RegistryEndpoint + "/" + args.Name + ":latest"

		imageBuildResponse, err := cli.ImageBuild(t.Ctx, packing.tar, dockerTypes.ImageBuildOptions{
			Dockerfile: "Dockerfile",
			Tags:       []string{tag},
		})
		if err != nil {
			return nil, fmt.Errorf("could not build image: %v", err)
		}
		defer imageBuildResponse.Body.Close()

		_, err = io.Copy(os.Stdout, imageBuildResponse.Body)
		if err != nil {
			return nil, fmt.Errorf("could not read image build response: %v", err)
		}

		return &buildImageResponse{
			namespace: packing.namespace,
			tag:       tag,
		}, nil
	})

	tasks.Add(actions, "Pushing image", func(t *tasks.Task, build *buildImageResponse) (*buildImageResponse, error) {
		accessKey, _ := client.GetAccessKey()
		secretKey, _ := client.GetSecretKey()
		authConfig := dockerTypes.AuthConfig{
			ServerAddress: build.namespace.RegistryEndpoint,
			Username:      accessKey,
			Password:      secretKey,
		}

		encodedJSON, err := json.Marshal(authConfig)
		if err != nil {
			return nil, fmt.Errorf("could not marshal auth config: %v", err)
		}

		authStr := base64.URLEncoding.EncodeToString(encodedJSON)

		imagePushResponse, err := cli.ImagePush(t.Ctx, build.tag, dockerTypes.ImagePushOptions{
			RegistryAuth: authStr,
		})
		if err != nil {
			return nil, fmt.Errorf("could not push image: %v", err)
		}
		defer imagePushResponse.Close()

		_, err = io.Copy(os.Stdout, imagePushResponse)
		if err != nil {
			return nil, fmt.Errorf("could not read image push response: %v", err)
		}

		return build, nil
	})

	tasks.Add(actions, "Creating container", func(t *tasks.Task, build *buildImageResponse) (*container.Container, error) {
		targetContainer, err := getOrCreateContainer(t.Ctx, api, args.Region, build.namespace.ID, args.Name)
		if err != nil {
			return nil, fmt.Errorf("could not get or create container: %v", err)
		}

		_, err = api.UpdateContainer(&container.UpdateContainerRequest{
			Region:        args.Region,
			ContainerID:   targetContainer.ID,
			RegistryImage: &build.tag,
			Redeploy:      scw.BoolPtr(false),
		}, scw.WithContext(t.Ctx))
		if err != nil {
			return nil, fmt.Errorf("could not update container: %v", err)
		}

		targetContainer, err = api.WaitForContainer(&container.WaitForContainerRequest{
			Region:      args.Region,
			ContainerID: targetContainer.ID,
			Timeout:     scw.TimeDurationPtr(12*time.Minute + 30*time.Second),
		}, scw.WithContext(t.Ctx))
		if err != nil {
			return nil, fmt.Errorf("failed to deploy container: %v", err)
		}

		return targetContainer, nil
	})

	tasks.Add(actions, "Deploying container", func(t *tasks.Task, targetContainer *container.Container) (*container.Container, error) {
		targetContainer, err := api.DeployContainer(&container.DeployContainerRequest{
			Region:      args.Region,
			ContainerID: targetContainer.ID,
		}, scw.WithContext(t.Ctx))
		if err != nil {
			return nil, fmt.Errorf("could not deploy container: %v", err)
		}

		targetContainer, err = api.WaitForContainer(&container.WaitForContainerRequest{
			Region:      args.Region,
			ContainerID: targetContainer.ID,
			Timeout:     scw.TimeDurationPtr(12*time.Minute + 30*time.Second),
		}, scw.WithContext(t.Ctx))
		if err != nil {
			return nil, fmt.Errorf("failed to deploy container: %v", err)
		}

		return targetContainer, nil
	})

	return actions.Execute(ctx, nil)
}

func getOrCreateNamespace(ctx context.Context, api *container.API, region scw.Region, name string) (*container.Namespace, error) {
	listNamespacesResponse, err := api.ListNamespaces(&container.ListNamespacesRequest{
		Region: region,
		Name:   &name,
	}, scw.WithContext(ctx))
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

	return namespace, nil
}

func getOrCreateContainer(ctx context.Context, api *container.API, region scw.Region, namespaceID string, name string) (*container.Container, error) {
	listContainersResponse, err := api.ListContainers(&container.ListContainersRequest{
		Region:      region,
		NamespaceID: namespaceID,
		Name:        &name,
	}, scw.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	containers := listContainersResponse.Containers

	var matchingContainer *container.Container
	for _, container := range containers {
		if container.Name == name {
			matchingContainer = container
			break
		}
	}

	if matchingContainer != nil {
		return matchingContainer, nil
	}

	container, err := api.CreateContainer(&container.CreateContainerRequest{
		Region:      region,
		NamespaceID: namespaceID,
		Name:        name,
		Port:        scw.Uint32Ptr(80),
	}, scw.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	return container, nil
}
