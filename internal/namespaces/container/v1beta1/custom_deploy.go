//go:build darwin || linux || windows

package container

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"time"

	pack "github.com/buildpacks/pack/pkg/client"
	"github.com/buildpacks/pack/pkg/logging"
	dockertypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/image"
	dockerregistry "github.com/docker/docker/api/types/registry"
	docker "github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/container/v1beta1/getorcreate"
	"github.com/scaleway/scaleway-cli/v2/internal/tasks"
	"github.com/scaleway/scaleway-cli/v2/internal/terminal"
	container "github.com/scaleway/scaleway-sdk-go/api/container/v1beta1"
	"github.com/scaleway/scaleway-sdk-go/api/registry/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type containerDeployRequest struct {
	Region scw.Region

	Name string

	Builder      string
	Dockerfile   string
	ForceBuilder bool

	BuildSource string
	Cache       bool
	BuildArgs   map[string]*string

	NamespaceID *string
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
				Short: "Name of the application (defaults to build-source's directory name)",
			},
			{
				Name:    "builder",
				Short:   "Builder image to use",
				Default: core.DefaultValueSetter("paketobuildpacks/builder-jammy-base:latest"),
			},
			{
				Name:    "dockerfile",
				Short:   "Path to the Dockerfile",
				Default: core.DefaultValueSetter("Dockerfile"),
			},
			{
				Name:    "force-builder",
				Short:   "Force the use of the builder image (even if a Dockerfile is present)",
				Default: core.DefaultValueSetter("false"),
			},
			{
				Name:    "build-source",
				Short:   "Path to the build context",
				Default: core.DefaultValueSetter("."),
			},
			{
				Name:    "cache",
				Short:   "Use cache when building the image",
				Default: core.DefaultValueSetter("true"),
			},
			{
				Name:     "build-args.{key}",
				Short:    "Build-time variables",
				Required: false,
			},
			{
				Name:    "port",
				Short:   "Port to expose",
				Default: core.DefaultValueSetter("8080"),
			},
			{
				Name:  "namespace-id",
				Short: "Container Namespace ID to deploy to",
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.RegionPlWaw,
				scw.Region(core.AllLocalities),
			),
		},
		Run: containerDeployRun,
	}
}

func containerDeployRun(ctx context.Context, argsI interface{}) (i interface{}, e error) {
	args := argsI.(*containerDeployRequest)
	buildSource, err := filepath.Abs(args.BuildSource)
	if err != nil {
		return nil, err
	}
	args.BuildSource = buildSource

	if args.Name == "" {
		args.Name = filepath.Base(args.BuildSource)
		if args.Name == "." {
			return nil, errors.New(
				"unable to determine application name, please specify it with name=",
			)
		}

		args.Name = "app-" + args.Name
	}

	client := core.ExtractClient(ctx)
	api := container.NewAPI(client)

	actions := tasks.Begin()

	if args.NamespaceID != nil {
		tasks.Add(actions, "Fetching namespace", DeployStepFetchNamespace)
	} else {
		tasks.Add(actions, "Creating namespace", DeployStepCreateNamespace)
	}

	tasks.Add(actions, "Fetch or create image registry", DeployStepFetchOrCreateRegistry)

	hasDockerfile := false
	if _, err := os.Stat(filepath.Join(args.BuildSource, args.Dockerfile)); err == nil {
		hasDockerfile = true
	}

	if hasDockerfile && !args.ForceBuilder {
		tasks.Add(actions, "Packing image", DeployStepDockerPackImage)
		tasks.Add(actions, "Building image", DeployStepDockerBuildImage)
	} else {
		tasks.Add(actions, "Building image", DeployStepBuildpackBuildImage)
	}

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

	container := result.(*DeployStepDeployContainerResponse).Container

	return fmt.Sprintln(
		terminal.Style("Your application is now available at", color.FgGreen),
		terminal.Style("https://"+container.DomainName, color.FgGreen, color.Bold),
	), nil
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

func DeployStepFetchNamespace(
	t *tasks.Task,
	data *DeployStepData,
) (*DeployStepCreateNamespaceResponse, error) {
	namespace, err := data.API.GetNamespace(&container.GetNamespaceRequest{
		Region:      data.Args.Region,
		NamespaceID: *data.Args.NamespaceID,
	}, scw.WithContext(t.Ctx))
	if err != nil {
		return nil, fmt.Errorf("could not fetch namespace: %w", err)
	}

	return &DeployStepCreateNamespaceResponse{
		DeployStepData: data,
		Namespace:      namespace,
	}, nil
}

func DeployStepCreateNamespace(
	t *tasks.Task,
	data *DeployStepData,
) (*DeployStepCreateNamespaceResponse, error) {
	namespace, err := getorcreate.Namespace(t.Ctx, data.API, data.Args.Region, data.Args.Name)
	if err != nil {
		return nil, err
	}

	return &DeployStepCreateNamespaceResponse{
		DeployStepData: data,
		Namespace:      namespace,
	}, nil
}

type DeployStepFetchOrCreateResponse struct {
	*DeployStepData
	Namespace        *container.Namespace
	RegistryEndpoint string
}

func DeployStepFetchOrCreateRegistry(
	t *tasks.Task,
	data *DeployStepCreateNamespaceResponse,
) (*DeployStepFetchOrCreateResponse, error) {
	registryEndpoint := data.Namespace.RegistryEndpoint
	if registryEndpoint == "" {
		registryAPI := registry.NewAPI(data.Client)
		registryNamespace, err := getorcreate.Registry(
			t.Ctx,
			registryAPI,
			data.Args.Region,
			data.Namespace.Name,
		)
		if err != nil {
			return nil, err
		}

		registryEndpoint = registryNamespace.Endpoint
	}

	return &DeployStepFetchOrCreateResponse{
		DeployStepData:   data.DeployStepData,
		Namespace:        data.Namespace,
		RegistryEndpoint: registryEndpoint,
	}, nil
}

type DeployStepPackImageResponse struct {
	*DeployStepData
	Namespace        *container.Namespace
	RegistryEndpoint string
	Tar              io.Reader
}

func DeployStepDockerPackImage(
	_ *tasks.Task,
	data *DeployStepFetchOrCreateResponse,
) (*DeployStepPackImageResponse, error) {
	tar, err := archive.TarWithOptions(data.Args.BuildSource, &archive.TarOptions{})
	if err != nil {
		return nil, fmt.Errorf("could not create tar: %w", err)
	}

	return &DeployStepPackImageResponse{
		DeployStepData:   data.DeployStepData,
		Namespace:        data.Namespace,
		RegistryEndpoint: data.RegistryEndpoint,
		Tar:              tar,
	}, nil
}

type DeployStepBuildImageResponse struct {
	*DeployStepData
	Namespace    *container.Namespace
	Tag          string
	DockerClient DockerClient
}

func DeployStepDockerBuildImage(
	t *tasks.Task,
	data *DeployStepPackImageResponse,
) (*DeployStepBuildImageResponse, error) {
	tag := data.RegistryEndpoint + "/" + data.Args.Name + ":latest"

	httpClient := core.ExtractHTTPClient(t.Ctx)
	dockerClient, err := docker.NewClientWithOpts(
		docker.FromEnv,
		docker.WithAPIVersionNegotiation(),
		docker.WithHTTPClient(httpClient),
	)
	if err != nil {
		return nil, fmt.Errorf("could not connect to Docker: %w", err)
	}

	imageBuildResponse, err := dockerClient.ImageBuild(
		t.Ctx,
		data.Tar,
		dockertypes.ImageBuildOptions{
			Dockerfile: data.Args.Dockerfile,
			Tags:       []string{tag},
			NoCache:    !data.Args.Cache,
			BuildArgs:  data.Args.BuildArgs,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("could not build image: %w", errors.Unwrap(err))
	}
	defer imageBuildResponse.Body.Close()

	err = jsonmessage.DisplayJSONMessagesStream(
		imageBuildResponse.Body,
		t.Logs,
		t.Logs.Fd(),
		true,
		nil,
	)
	if err != nil {
		if jerr, ok := err.(*jsonmessage.JSONError); ok {
			// If no error code is set, default to 1
			if jerr.Code == 0 {
				jerr.Code = 1
			}

			return nil, fmt.Errorf(
				"docker build failed with error code %d: %s",
				jerr.Code,
				jerr.Message,
			)
		}

		return nil, err
	}

	return &DeployStepBuildImageResponse{
		DeployStepData: data.DeployStepData,
		Namespace:      data.Namespace,
		Tag:            tag,
		DockerClient:   dockerClient,
	}, nil
}

func DeployStepBuildpackBuildImage(
	t *tasks.Task,
	data *DeployStepFetchOrCreateResponse,
) (*DeployStepBuildImageResponse, error) {
	tag := data.RegistryEndpoint + "/" + data.Args.Name + ":latest"

	httpClient := core.ExtractHTTPClient(t.Ctx)
	dockerClient, err := NewCustomDockerClient(httpClient)
	if err != nil {
		return nil, err
	}

	packClient, err := pack.NewClient(
		pack.WithDockerClient(dockerClient),
		pack.WithLogger(logging.NewLogWithWriters(t.Logs, t.Logs)),
	)
	if err != nil {
		return nil, fmt.Errorf("could not create pack client: %w", err)
	}

	err = packClient.Build(t.Ctx, pack.BuildOptions{
		AppPath:      data.Args.BuildSource,
		Builder:      data.Args.Builder,
		Image:        tag,
		ClearCache:   !data.Args.Cache,
		TrustBuilder: func(string) bool { return true },
	})
	if err != nil {
		return nil, fmt.Errorf("could not build: %w", err)
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

func DeployStepPushImage(
	t *tasks.Task,
	data *DeployStepBuildImageResponse,
) (*DeployStepPushImageResponse, error) {
	accessKey, _ := data.Client.GetAccessKey()
	secretKey, _ := data.Client.GetSecretKey()
	authConfig := dockerregistry.AuthConfig{
		ServerAddress: data.Namespace.RegistryEndpoint,
		Username:      accessKey,
		Password:      secretKey,
	}

	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		return nil, fmt.Errorf("could not marshal auth config: %w", err)
	}

	authStr := base64.URLEncoding.EncodeToString(encodedJSON)

	imagePushResponse, err := data.DockerClient.ImagePush(t.Ctx, data.Tag, image.PushOptions{
		RegistryAuth: authStr,
	})
	if err != nil {
		return nil, fmt.Errorf("could not push image: %w", err)
	}
	defer imagePushResponse.Close()

	err = jsonmessage.DisplayJSONMessagesStream(imagePushResponse, t.Logs, t.Logs.Fd(), true, nil)
	if err != nil {
		if jerr, ok := err.(*jsonmessage.JSONError); ok {
			// If no error code is set, default to 1
			if jerr.Code == 0 {
				jerr.Code = 1
			}

			return nil, fmt.Errorf(
				"docker build failed with error code %d: %s",
				jerr.Code,
				jerr.Message,
			)
		}

		return nil, err
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

func DeployStepCreateContainer(
	t *tasks.Task,
	data *DeployStepPushImageResponse,
) (*DeployStepCreateContainerResponse, error) {
	targetContainer, err := getorcreate.Container(
		t.Ctx,
		data.API,
		data.Args.Region,
		data.Namespace.ID,
		data.Args.Name,
	)
	if err != nil {
		return nil, fmt.Errorf("could not get or create container: %w", err)
	}

	_, err = data.API.UpdateContainer(&container.UpdateContainerRequest{
		Region:        data.Args.Region,
		ContainerID:   targetContainer.ID,
		RegistryImage: &data.Tag,
		Port:          scw.Uint32Ptr(data.Args.Port),
		Redeploy:      scw.BoolPtr(false),
	}, scw.WithContext(t.Ctx))
	if err != nil {
		return nil, fmt.Errorf("could not update container: %w", err)
	}

	targetContainer, err = data.API.WaitForContainer(&container.WaitForContainerRequest{
		Region:        data.Args.Region,
		ContainerID:   targetContainer.ID,
		Timeout:       scw.TimeDurationPtr(12*time.Minute + 30*time.Second),
		RetryInterval: core.DefaultRetryInterval,
	}, scw.WithContext(t.Ctx))
	if err != nil {
		return nil, fmt.Errorf("failed to deploy container: %w", err)
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

func DeployStepDeployContainer(
	t *tasks.Task,
	data *DeployStepCreateContainerResponse,
) (*DeployStepDeployContainerResponse, error) {
	targetContainer, err := data.API.DeployContainer(&container.DeployContainerRequest{
		Region:      data.Args.Region,
		ContainerID: data.Container.ID,
	}, scw.WithContext(t.Ctx))
	if err != nil {
		return nil, fmt.Errorf("could not deploy container: %w", err)
	}

	targetContainer, err = data.API.WaitForContainer(&container.WaitForContainerRequest{
		Region:        data.Args.Region,
		ContainerID:   targetContainer.ID,
		Timeout:       scw.TimeDurationPtr(12*time.Minute + 30*time.Second),
		RetryInterval: core.DefaultRetryInterval,
	}, scw.WithContext(t.Ctx))
	if err != nil {
		return nil, fmt.Errorf("failed to deploy container: %w", err)
	}

	return &DeployStepDeployContainerResponse{
		DeployStepData: data.DeployStepData,
		Container:      targetContainer,
	}, nil
}
