//go:build darwin || linux || windows

package object

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type lifecycleConfigurationCreateArgs struct {
	Bucket            string
	ConfigurationPath string
	Region            scw.Region
	ProjectID         string
}

func lifecycleConfigurationCreateCommand() *core.Command {
	return &core.Command{
		Namespace: "object",
		Resource:  "bucket-lifecycle-configuration",
		Verb:      "create",
		Short:     "Create a lifecycle configuration for an S3 bucket's objects.",
		Long:      "Create a lifecycle configuration and apply to an Object Bucket's objectswith the S3 protocol.",
		ArgsType:  reflect.TypeOf(lifecycleConfigurationCreateArgs{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "bucket",
				Positional: true,
				Required:   true,
				Short:      "The name of the bucket to which assign the policy.",
			},
			{
				Name:       "configuration-path",
				Positional: false,
				Required:   true,
				Short:      "The path to the local JSON file containing the bucket lifecycle configuration.",
			},
			core.ProjectIDArgSpec(),
			core.RegionArgSpec(),
		},
		Run: func(ctx context.Context, argsI any) (any, error) {
			args := argsI.(*lifecycleConfigurationCreateArgs)
			if args.Bucket == "" {
				return nil, errors.New("bucket name cannot be empty")
			}
			if args.ConfigurationPath == "" {
				return nil, errors.New("configuration file path cannot be empty")
			}

			configPath, err := os.ReadFile(args.ConfigurationPath)
			if err != nil {
				return nil, fmt.Errorf(
					"could not open configuration file \"%s\": %w",
					args.ConfigurationPath,
					err,
				)
			}

			configStr := string(configPath)

			var config types.BucketLifecycleConfiguration
			err = json.Unmarshal([]byte(configStr), &config)
			if err != nil {
				return nil, fmt.Errorf(
					"failed to unmarshal bucket lifecycle configuration: %v",
					err,
				)
			}

			client := newS3Client(ctx, args.Region, args.ProjectID)
			params := s3.PutBucketLifecycleConfigurationInput{
				Bucket:                 &args.Bucket,
				LifecycleConfiguration: &config,
			}

			_, err = client.PutBucketLifecycleConfiguration(ctx, &params)
			if err != nil {
				return nil, fmt.Errorf("could not create bucket lifecycle configuration: %w", err)
			}

			bucket, err := getBucketInfo(ctx, args.Region, args.Bucket, args.ProjectID)
			if err != nil {
				return nil, fmt.Errorf("could not get bucket's information: %w", err)
			}

			return &BucketResponse{
				BucketInfo: bucket,
				SuccessResult: &core.SuccessResult{
					Resource: "bucket-lifecycle-configuration",
					Verb:     "create",
				},
			}, nil
		},
	}
}

type lifecycleConfigurationDeleteArgs struct {
	Bucket    string
	Region    scw.Region
	ProjectID string
}

func lifecycleConfigurationDeleteCommand() *core.Command {
	return &core.Command{
		Namespace: "object",
		Resource:  "bucket-lifecycle-configuration",
		Verb:      "delete",
		Short:     "Delete an S3 bucket's lifecycle configuration if it exists.",
		Long:      "Delete an Object Bucket's lifecycle configuration with the S3 protocol.",
		ArgsType:  reflect.TypeOf(lifecycleConfigurationDeleteArgs{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:             "bucket",
				Positional:       true,
				Required:         true,
				Short:            "The unique name of the bucket",
				AutoCompleteFunc: autocompleteBucketName,
			},
			core.ProjectIDArgSpec(),
			core.RegionArgSpec(),
		},
		Run: func(ctx context.Context, argsI any) (any, error) {
			args := argsI.(*lifecycleConfigurationDeleteArgs)
			if args.Bucket == "" {
				return nil, errors.New("bucket name cannot be empty")
			}

			client := newS3Client(ctx, args.Region, args.ProjectID)
			params := s3.DeleteBucketLifecycleInput{
				Bucket: &args.Bucket,
			}

			_, err := client.DeleteBucketLifecycle(ctx, &params)
			if err != nil {
				return nil, fmt.Errorf("could not delete bucket lifecycle: %w", err)
			}

			return &core.SuccessResult{
				Resource: "bucket-lifecycle-configuration",
				Verb:     "delete",
			}, nil
		},
	}
}

type lifecycleConfigurationGetArgs struct {
	Bucket    string
	Region    scw.Region
	ProjectID string
}

func lifecycleConfigurationGetCommand() *core.Command {
	return &core.Command{
		Namespace: "object",
		Resource:  "bucket-lifecycle-configuration",
		Verb:      "get",
		Short:     "Get the lifecycle configuration of an S3 bucket.",
		Long:      "Retrieve an Object Bucket's list of lifecycle rules with the S3 protocol.",
		ArgsType:  reflect.TypeOf(lifecycleConfigurationGetArgs{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:             "bucket",
				Positional:       true,
				Required:         true,
				Short:            "The unique name of the bucket",
				AutoCompleteFunc: autocompleteBucketName,
			},
			core.ProjectIDArgSpec(),
			core.RegionArgSpec(),
		},
		Run: func(ctx context.Context, argsI any) (any, error) {
			args := argsI.(*lifecycleConfigurationGetArgs)
			if args.Bucket == "" {
				return nil, errors.New("bucket name cannot be empty")
			}

			client := newS3Client(ctx, args.Region, args.ProjectID)
			params := s3.GetBucketLifecycleConfigurationInput{
				Bucket: &args.Bucket,
			}

			conf, err := client.GetBucketLifecycleConfiguration(ctx, &params)
			if err != nil {
				return nil, fmt.Errorf("could not get bucket configuration: %w", err)
			}
			if conf == nil {
				return nil, errors.New("could not get bucket configuration (result is nil)")
			}

			prettyJSONBytes, err := CleanAndIndentJSON(conf, "", "    ")
			if err != nil {
				return nil, fmt.Errorf("could not format bucket configuration: %w", err)
			}

			return string(prettyJSONBytes), nil
		},
	}
}
