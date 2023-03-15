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
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/container/v1beta1/container_utils"
	"github.com/scaleway/scaleway-cli/v2/internal/tasks"
	container "github.com/scaleway/scaleway-sdk-go/api/container/v1beta1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type containerDeployRequest struct {
	Region scw.Region

	Name        string
	Dockerfile  string
	BuildSource string
	Port        uint32
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
			{
				Name:    "dockerfile",
				Short:   "Path to the Dockerfile",
				Default: core.DefaultValueSetter("Dockerfile"),
			},
			{
				Name:    "build-source",
				Short:   "Path to the build context",
				Default: core.DefaultValueSetter("."),
			},
			{
				Name:    "port",
				Short:   "Port to expose",
				Default: core.DefaultValueSetter("8080"),
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

	fileInfo, err := os.Stat(args.Dockerfile)
	if err != nil {
		return nil, fmt.Errorf("could not stat '%s': %v", args.Dockerfile, err)
	}

	if fileInfo.IsDir() {
		return nil, fmt.Errorf("'%s' is a directory", args.Dockerfile)
	}

	actions := tasks.Begin()
	tasks.Add(actions, "Creating namespace", DeployStepCreateNamespace)
	tasks.Add(actions, "Packing image", DeployStepPackImage)
	tasks.Add(actions, "Building image", DeployStepBuildImage)
	tasks.Add(actions, "Pushing image", DeployStepPushImage)
	tasks.Add(actions, "Creating container", DeployStepCreateContainer)
	tasks.Add(actions, "Deploying container", DeployStepDeployContainer)

	result, err := actions.Execute(ctx, &DeployStepData{
		Client: client,
		API:    api,
		Args:   args,
	})
	if err != nil {
		return nil, err
	}

	return result.(*DeployStepDeployContainerResponse).Container, nil
}

type DeployStepData struct {
	Client *scw.Client
	API    *container.API
	Args   *containerDeployRequest
}

type DeployStepCreateNamespaceResponse struct {
	*DeployStepData
	Namespace *container.Namespace
}

func DeployStepCreateNamespace(t *tasks.Task, data *DeployStepData) (*DeployStepCreateNamespaceResponse, error) {
	namespace, err := container_utils.GetOrCreateNamespace(t.Ctx, data.API, data.Args.Region, data.Args.Name)
	if err != nil {
		return nil, err
	}

	return &DeployStepCreateNamespaceResponse{
		DeployStepData: data,
		Namespace:      namespace,
	}, nil
}

type DeployStepPackImageResponse struct {
	*DeployStepData
	Namespace *container.Namespace
	Tar       io.Reader
}

func DeployStepPackImage(t *tasks.Task, data *DeployStepCreateNamespaceResponse) (*DeployStepPackImageResponse, error) {
	tar, err := archive.TarWithOptions(data.Args.BuildSource, &archive.TarOptions{})
	if err != nil {
		return nil, fmt.Errorf("could not create tar: %v", err)
	}

	return &DeployStepPackImageResponse{
		DeployStepData: data.DeployStepData,
		Namespace:      data.Namespace,
		Tar:            tar,
	}, nil
}

type DeployStepBuildImageResponse struct {
	*DeployStepData
	Namespace    *container.Namespace
	Tag          string
	DockerClient *docker.Client
}

func DeployStepBuildImage(t *tasks.Task, data *DeployStepPackImageResponse) (*DeployStepBuildImageResponse, error) {
	tag := data.Namespace.RegistryEndpoint + "/" + data.Args.Name + ":latest"

	dockerClient, err := docker.NewClientWithOpts(docker.FromEnv, docker.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("could not connect to Docker: %v", err)
	}

	imageBuildResponse, err := dockerClient.ImageBuild(t.Ctx, data.Tar, dockerTypes.ImageBuildOptions{
		Dockerfile: data.Args.Dockerfile,
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

	return &DeployStepBuildImageResponse{
		DeployStepData: data.DeployStepData,
		Namespace:      data.Namespace,
		Tag:            tag,
		DockerClient:   dockerClient,
	}, nil
}

type DeployStepPushImageResponse struct {
	*DeployStepData
	Namespace *container.Namespace
	Tag       string
}

func DeployStepPushImage(t *tasks.Task, data *DeployStepBuildImageResponse) (*DeployStepPushImageResponse, error) {
	accessKey, _ := data.Client.GetAccessKey()
	secretKey, _ := data.Client.GetSecretKey()
	authConfig := dockerTypes.AuthConfig{
		ServerAddress: data.Namespace.RegistryEndpoint,
		Username:      accessKey,
		Password:      secretKey,
	}

	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		return nil, fmt.Errorf("could not marshal auth config: %v", err)
	}

	authStr := base64.URLEncoding.EncodeToString(encodedJSON)

	imagePushResponse, err := data.DockerClient.ImagePush(t.Ctx, data.Tag, dockerTypes.ImagePushOptions{
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

	return &DeployStepPushImageResponse{
		DeployStepData: data.DeployStepData,
		Namespace:      data.Namespace,
		Tag:            data.Tag,
	}, nil
}

type DeployStepCreateContainerResponse struct {
	*DeployStepData
	Container *container.Container
}

func DeployStepCreateContainer(t *tasks.Task, data *DeployStepPushImageResponse) (*DeployStepCreateContainerResponse, error) {
	targetContainer, err := container_utils.GetOrCreateContainer(t.Ctx, data.API, data.Args.Region, data.Namespace.ID, data.Args.Name)
	if err != nil {
		return nil, fmt.Errorf("could not get or create container: %v", err)
	}

	_, err = data.API.UpdateContainer(&container.UpdateContainerRequest{
		Region:        data.Args.Region,
		ContainerID:   targetContainer.ID,
		RegistryImage: &data.Tag,
		Port:          scw.Uint32Ptr(data.Args.Port),
		Redeploy:      scw.BoolPtr(false),
	}, scw.WithContext(t.Ctx))
	if err != nil {
		return nil, fmt.Errorf("could not update container: %v", err)
	}

	targetContainer, err = data.API.WaitForContainer(&container.WaitForContainerRequest{
		Region:      data.Args.Region,
		ContainerID: targetContainer.ID,
		Timeout:     scw.TimeDurationPtr(12*time.Minute + 30*time.Second),
	}, scw.WithContext(t.Ctx))
	if err != nil {
		return nil, fmt.Errorf("failed to deploy container: %v", err)
	}

	return &DeployStepCreateContainerResponse{
		DeployStepData: data.DeployStepData,
		Container:      targetContainer,
	}, nil
}

type DeployStepDeployContainerResponse struct {
	*DeployStepData
	Container *container.Container
}

func DeployStepDeployContainer(t *tasks.Task, data *DeployStepCreateContainerResponse) (*DeployStepDeployContainerResponse, error) {
	targetContainer, err := data.API.DeployContainer(&container.DeployContainerRequest{
		Region:      data.Args.Region,
		ContainerID: data.Container.ID,
	}, scw.WithContext(t.Ctx))
	if err != nil {
		return nil, fmt.Errorf("could not deploy container: %v", err)
	}

	targetContainer, err = data.API.WaitForContainer(&container.WaitForContainerRequest{
		Region:      data.Args.Region,
		ContainerID: targetContainer.ID,
		Timeout:     scw.TimeDurationPtr(12*time.Minute + 30*time.Second),
	}, scw.WithContext(t.Ctx))
	if err != nil {
		return nil, fmt.Errorf("failed to deploy container: %v", err)
	}

	return &DeployStepDeployContainerResponse{
		DeployStepData: data.DeployStepData,
		Container:      targetContainer,
	}, nil
}
