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

type lifecycleCreateArgs struct {
	Bucket                 string
	LifecycleConfiguration string
	Region                 scw.Region
	ProjectID              string
	S3Endpoint             string `json:"s3-endpoint"`
	S3UsePathStyle         string `json:"s3-use-path-style"`
}

func lifecycleCreateCommand() *core.Command {
	return &core.Command{
		Namespace: "object",
		Resource:  "bucket-lifecycle",
		Verb:      "create",
		Short:     "Create a lifecycle configuration for an S3 bucket's objects.",
		Long:      "Create a lifecycle configuration and apply to an Object Bucket's objects with the S3 protocol.",
		ArgsType:  reflect.TypeFor[lifecycleCreateArgs](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "bucket",
				Positional: true,
				Required:   true,
				Short:      "The name of the bucket to which assign the lifecycle configuration.",
			},
			{
				Name:       "lifecycle-configuration",
				Positional: false,
				Required:   true,
				Short:      "The path to the local JSON file containing the bucket lifecycle configuration.",
			},
			{
				Name:       "s3-endpoint",
				Positional: false,
				Required:   false,
				Short:      "Custom S3 endpoint to use instead of the default",
			},
			{
				Name:       "s3-use-path-style",
				Positional: false,
				Required:   false,
				Short:      "Whether to use path style addressing for S3 API calls or not",
			},
			core.ProjectIDArgSpec(),
			core.RegionArgSpec(),
		},
		Run: func(ctx context.Context, argsI any) (any, error) {
			args := argsI.(*lifecycleCreateArgs)
			if args.Bucket == "" {
				return nil, errors.New("bucket name cannot be empty")
			}
			if args.LifecycleConfiguration == "" {
				return nil, errors.New("configuration file path cannot be empty")
			}

			configPath, err := os.ReadFile(args.LifecycleConfiguration)
			if err != nil {
				return nil, fmt.Errorf(
					"could not open configuration file \"%s\": %w",
					args.LifecycleConfiguration,
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

			s3UsePathStyle, err := getS3UsePathStyle(ctx, args.S3UsePathStyle)
			if err != nil {
				return nil, fmt.Errorf("could not get S3 use path style flag: %w", err)
			}

			s3Endpoint, _, err := getS3Endpoints(
				ctx, args.Region.String(), args.S3Endpoint, "", s3UsePathStyle,
			)
			if err != nil {
				return nil, fmt.Errorf("could not get S3 endpoints: %w", err)
			}

			client := newS3Client(ctx, args.Region, args.ProjectID, s3Endpoint, s3UsePathStyle)

			params := s3.PutBucketLifecycleConfigurationInput{
				Bucket:                 &args.Bucket,
				LifecycleConfiguration: &config,
			}

			_, err = client.PutBucketLifecycleConfiguration(ctx, &params)
			if err != nil {
				return nil, fmt.Errorf("could not create bucket lifecycle configuration: %w", err)
			}

			bucket, err := getBucketInfo(
				ctx, args.Region, args.Bucket, args.ProjectID, s3UsePathStyle, s3Endpoint,
			)
			if err != nil {
				return nil, fmt.Errorf("could not get bucket's information: %w", err)
			}

			return &BucketResponse{
				BucketInfo: bucket,
				SuccessResult: &core.SuccessResult{
					Resource: "bucket-lifecycle",
					Verb:     "create",
				},
			}, nil
		},
	}
}

type lifecycleDeleteArgs struct {
	Bucket         string
	Region         scw.Region
	ProjectID      string
	S3Endpoint     string `json:"s3-endpoint"`
	S3UsePathStyle string `json:"s3-use-path-style"`
}

func lifecycleDeleteCommand() *core.Command {
	return &core.Command{
		Namespace: "object",
		Resource:  "bucket-lifecycle",
		Verb:      "delete",
		Short:     "Delete an S3 bucket's lifecycle configuration if it exists.",
		Long:      "Delete an Object Bucket's lifecycle configuration with the S3 protocol.",
		ArgsType:  reflect.TypeFor[lifecycleDeleteArgs](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:             "bucket",
				Positional:       true,
				Required:         true,
				Short:            "The unique name of the bucket",
				AutoCompleteFunc: autocompleteBucketName,
			},
			{
				Name:       "s3-endpoint",
				Positional: false,
				Required:   false,
				Short:      "Custom S3 endpoint to use instead of the default",
			},
			{
				Name:       "s3-use-path-style",
				Positional: false,
				Required:   false,
				Short:      "Whether to use path style addressing for S3 API calls or not",
			},
			core.ProjectIDArgSpec(),
			core.RegionArgSpec(),
		},
		Run: func(ctx context.Context, argsI any) (any, error) {
			args := argsI.(*lifecycleDeleteArgs)
			if args.Bucket == "" {
				return nil, errors.New("bucket name cannot be empty")
			}

			s3UsePathStyle, err := getS3UsePathStyle(ctx, args.S3UsePathStyle)
			if err != nil {
				return nil, fmt.Errorf("could not get S3 use path style flag: %w", err)
			}

			s3Endpoint, _, err := getS3Endpoints(
				ctx, args.Region.String(), args.S3Endpoint, "", s3UsePathStyle,
			)
			if err != nil {
				return nil, fmt.Errorf("could not get S3 endpoints: %w", err)
			}

			client := newS3Client(ctx, args.Region, args.ProjectID, s3Endpoint, s3UsePathStyle)

			params := s3.DeleteBucketLifecycleInput{
				Bucket: &args.Bucket,
			}

			_, err = client.DeleteBucketLifecycle(ctx, &params)
			if err != nil {
				return nil, fmt.Errorf("could not delete bucket lifecycle: %w", err)
			}

			return &core.SuccessResult{
				Resource: "bucket-lifecycle",
				Verb:     "delete",
			}, nil
		},
	}
}

type lifecycleGetArgs struct {
	Bucket         string
	Region         scw.Region
	ProjectID      string
	S3Endpoint     string `json:"s3-endpoint"`
	S3UsePathStyle string `json:"s3-use-path-style"`
}

func lifecycleGetCommand() *core.Command {
	return &core.Command{
		Namespace: "object",
		Resource:  "bucket-lifecycle",
		Verb:      "get",
		Short:     "Get the lifecycle configuration of an S3 bucket.",
		Long:      "Retrieve an Object Bucket's list of lifecycle rules with the S3 protocol.",
		ArgsType:  reflect.TypeFor[lifecycleGetArgs](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:             "bucket",
				Positional:       true,
				Required:         true,
				Short:            "The unique name of the bucket",
				AutoCompleteFunc: autocompleteBucketName,
			},
			{
				Name:       "s3-endpoint",
				Positional: false,
				Required:   false,
				Short:      "Custom S3 endpoint to use instead of the default",
			},
			{
				Name:       "s3-use-path-style",
				Positional: false,
				Required:   false,
				Short:      "Whether to use path style addressing for S3 API calls or not",
			},
			core.ProjectIDArgSpec(),
			core.RegionArgSpec(),
		},
		Run: func(ctx context.Context, argsI any) (any, error) {
			args := argsI.(*lifecycleGetArgs)
			if args.Bucket == "" {
				return nil, errors.New("bucket name cannot be empty")
			}

			s3UsePathStyle, err := getS3UsePathStyle(ctx, args.S3UsePathStyle)
			if err != nil {
				return nil, fmt.Errorf("could not get S3 use path style flag: %w", err)
			}

			s3Endpoint, _, err := getS3Endpoints(
				ctx, args.Region.String(), args.S3Endpoint, "", s3UsePathStyle,
			)
			if err != nil {
				return nil, fmt.Errorf("could not get S3 endpoints: %w", err)
			}

			client := newS3Client(ctx, args.Region, args.ProjectID, s3Endpoint, s3UsePathStyle)

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
