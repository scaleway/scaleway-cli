//go:build !wasm

package function

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/tasks"
	function "github.com/scaleway/scaleway-sdk-go/api/function/v1beta1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type functionDeployRequest struct {
	NamespaceID string                   `json:"namespace_id"`
	ZipFile     string                   `json:"zip_file"`
	Runtime     function.FunctionRuntime `json:"runtime"`
	Name        string                   `json:"name"`
	Region      scw.Region               `json:"region"`
}

func functionDeploy() *core.Command {
	functionCreate := functionFunctionCreate()

	return &core.Command{
		Short:     `Deploy a function`,
		Long:      `Create or fetch, upload and deploy your function`,
		Namespace: "function",
		Resource:  "deploy",
		Groups:    []string{"workflow"},
		ArgsType:  reflect.TypeOf(functionDeployRequest{}),
		ArgSpecs: []*core.ArgSpec{
			{
				Name:  "namespace-id",
				Short: "Function Namespace ID to deploy to",
			},
			{
				Name:     "name",
				Short:    "Name of the function to deploy, will be used in namespace's name if no ID is provided",
				Required: true,
			},
			{
				Name:       "runtime",
				EnumValues: functionCreate.ArgSpecs.GetByName("runtime").EnumValues,
				Required:   true,
			},
			{
				Name:     "zip-file",
				Short:    "Path of the zip file that contains your code",
				Required: true,
			},
			core.RegionArgSpec((&function.API{}).Regions()...),
		},
		Run: func(ctx context.Context, argsI interface{}) (interface{}, error) {
			args := argsI.(*functionDeployRequest)
			scwClient := core.ExtractClient(ctx)
			httpClient := core.ExtractHTTPClient(ctx)
			api := function.NewAPI(scwClient)

			if err := validateRuntime(api, args.Region, args.Runtime); err != nil {
				return nil, err
			}

			zipFileStat, err := os.Stat(args.ZipFile)
			if err != nil {
				return nil, fmt.Errorf("failed to stat zip-file: %w", err)
			}

			if zipFileStat.Size() < 0 {
				return nil, errors.New("invalid zip-file, invalid size")
			}

			ts := tasks.Begin()

			if args.NamespaceID != "" {
				tasks.Add(
					ts,
					"Fetching namespace",
					DeployStepFetchNamespace(api, args.Region, args.NamespaceID),
				)
			} else {
				tasks.Add(ts, "Creating or fetching namespace", DeployStepCreateNamespace(api, args.Region, args.Name))
			}
			tasks.Add(
				ts,
				"Creating or fetching function",
				DeployStepCreateFunction(api, args.Name, args.Runtime),
			)
			tasks.Add(
				ts,
				"Uploading function",
				DeployStepFunctionUpload(
					httpClient,
					scwClient,
					api,
					args.ZipFile,
					zipFileStat.Size(),
				),
			)
			tasks.Add(ts, "Deploying function", DeployStepFunctionDeploy(api, args.Runtime))

			return ts.Execute(ctx, nil)
		},
	}
}

func validateRuntime(api *function.API, region scw.Region, runtime function.FunctionRuntime) error {
	runtimeName := string(runtime)

	resp, err := api.ListFunctionRuntimes(&function.ListFunctionRuntimesRequest{
		Region: region,
	})
	if err != nil {
		return fmt.Errorf("failed to list available runtimes: %w", err)
	}
	for _, r := range resp.Runtimes {
		if r.Name == runtimeName {
			return nil
		}
	}

	return fmt.Errorf("invalid runtime %q", runtimeName)
}

func DeployStepCreateNamespace(
	api *function.API,
	region scw.Region,
	functionName string,
) tasks.TaskFunc[any, *function.Namespace] {
	return func(t *tasks.Task, _ any) (nextArgs *function.Namespace, err error) {
		namespaceName := functionName

		namespaces, err := api.ListNamespaces(&function.ListNamespacesRequest{
			Region: region,
			Name:   &namespaceName,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to list namespaces: %w", err)
		}
		for _, ns := range namespaces.Namespaces {
			if ns.Name == namespaceName {
				return ns, nil
			}
		}

		namespace, err := api.CreateNamespace(&function.CreateNamespaceRequest{
			Name:   namespaceName,
			Region: region,
		}, scw.WithContext(t.Ctx))
		if err != nil {
			return nil, fmt.Errorf("could not create namespace: %w", err)
		}

		t.AddToCleanUp(func(_ context.Context) error {
			_, err := api.DeleteNamespace(&function.DeleteNamespaceRequest{
				Region:      namespace.Region,
				NamespaceID: namespace.ID,
			})

			return err
		})

		namespace, err = api.WaitForNamespace(&function.WaitForNamespaceRequest{
			NamespaceID: namespace.ID,
			Region:      namespace.Region,
		})
		if err != nil {
			return nil, fmt.Errorf("could not fetch created namespace: %w", err)
		}

		return namespace, nil
	}
}

func DeployStepFetchNamespace(
	api *function.API,
	region scw.Region,
	namespaceID string,
) tasks.TaskFunc[any, *function.Namespace] {
	return func(_ *tasks.Task, _ any) (nextArgs *function.Namespace, err error) {
		namespace, err := api.WaitForNamespace(&function.WaitForNamespaceRequest{
			NamespaceID: namespaceID,
			Region:      region,
		})
		if err != nil {
			return nil, fmt.Errorf("could not fetch namespace: %w", err)
		}

		return namespace, nil
	}
}

func DeployStepCreateFunction(
	api *function.API,
	functionName string,
	runtime function.FunctionRuntime,
) tasks.TaskFunc[*function.Namespace, *function.Function] {
	return func(t *tasks.Task, namespace *function.Namespace) (*function.Function, error) {
		functions, err := api.ListFunctions(&function.ListFunctionsRequest{
			Name:        &functionName,
			NamespaceID: namespace.ID,
			Region:      namespace.Region,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to list functions: %w", err)
		}
		for _, fc := range functions.Functions {
			if fc.Name == functionName {
				return fc, err
			}
		}

		fc, err := api.CreateFunction(&function.CreateFunctionRequest{
			Name:        functionName,
			NamespaceID: namespace.ID,
			Runtime:     runtime,
			Region:      namespace.Region,
		}, scw.WithContext(t.Ctx))
		if err != nil {
			return nil, fmt.Errorf("could not create function: %w", err)
		}

		t.AddToCleanUp(func(_ context.Context) error {
			_, err := api.DeleteFunction(&function.DeleteFunctionRequest{
				FunctionID: fc.ID,
				Region:     fc.Region,
			})

			return err
		})

		return fc, nil
	}
}

func DeployStepFunctionUpload(
	httpClient *http.Client,
	scwClient *scw.Client,
	api *function.API,
	zipPath string,
	zipSize int64,
) tasks.TaskFunc[*function.Function, *function.Function] {
	return func(t *tasks.Task, fc *function.Function) (nextArgs *function.Function, err error) {
		uploadURL, err := api.GetFunctionUploadURL(&function.GetFunctionUploadURLRequest{
			Region:        fc.Region,
			FunctionID:    fc.ID,
			ContentLength: uint64(zipSize),
		})
		if err != nil {
			return nil, err
		}

		zip, err := os.Open(zipPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read zip file: %w", err)
		}
		defer zip.Close()

		req, err := http.NewRequest(http.MethodPut, uploadURL.URL, zip)
		if err != nil {
			return nil, fmt.Errorf("failed to init request: %w", err)
		}
		req = req.WithContext(t.Ctx)
		req.ContentLength = zipSize

		for headerName, headerList := range uploadURL.Headers {
			for _, header := range *headerList {
				req.Header.Add(headerName, header)
			}
		}

		secretKey, _ := scwClient.GetSecretKey()
		req.Header.Add("X-Auth-Token", secretKey)

		resp, err := httpClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to send upload request: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("failed to upload function (Status: %d)", resp.StatusCode)
		}

		return fc, nil
	}
}

func DeployStepFunctionDeploy(
	api *function.API,
	runtime function.FunctionRuntime,
) tasks.TaskFunc[*function.Function, *function.Function] {
	return func(_ *tasks.Task, fc *function.Function) (*function.Function, error) {
		fc, err := api.UpdateFunction(&function.UpdateFunctionRequest{
			Region:     fc.Region,
			FunctionID: fc.ID,
			Runtime:    runtime,
			Redeploy:   scw.BoolPtr(true),
		})
		if err != nil {
			return nil, err
		}

		return api.WaitForFunction(&function.WaitForFunctionRequest{
			FunctionID:    fc.ID,
			Region:        fc.Region,
			RetryInterval: core.DefaultRetryInterval,
		})
	}
}
