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

type bucketCorsCreateArgs struct {
	Bucket            string
	CorsConfiguration string
	Region            scw.Region
	ProjectID         string
	S3Endpoint        string `json:"s3-endpoint"`
	S3UsePathStyle    string `json:"s3-use-path-style"`
}

func bucketCorsCreateCommand() *core.Command {
	return &core.Command{
		Namespace: "object",
		Resource:  "bucket-cors",
		Verb:      "create",
		Short:     "Add CORS rules to a bucket",
		Long:      "Add CORS rules to an Object Storage bucket with the S3 protocol.",
		ArgsType:  reflect.TypeFor[bucketCorsCreateArgs](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:             "bucket",
				Positional:       true,
				Required:         true,
				Short:            "The name of the bucket",
				AutoCompleteFunc: autocompleteBucketName,
			},
			{
				Name:       "cors-configuration",
				Positional: false,
				Required:   true,
				Short:      "The path to the local JSON file containing the CORS configuration.",
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
			args := argsI.(*bucketCorsCreateArgs)
			if args.Bucket == "" {
				return nil, errors.New("bucket name cannot be empty")
			}
			if args.CorsConfiguration == "" {
				return nil, errors.New("configuration file path cannot be empty")
			}

			configPath, err := os.ReadFile(args.CorsConfiguration)
			if err != nil {
				return nil, fmt.Errorf(
					"could not open configuration file \"%s\": %w",
					args.CorsConfiguration,
					err,
				)
			}

			configStr := string(configPath)

			var config types.CORSConfiguration
			err = json.Unmarshal([]byte(configStr), &config)
			if err != nil {
				return nil, fmt.Errorf(
					"failed to unmarshal bucket CORS configuration: %v",
					err,
				)
			}

			params := s3.PutBucketCorsInput{
				Bucket:            &args.Bucket,
				CORSConfiguration: &config,
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

			_, err = client.PutBucketCors(ctx, &params)
			if err != nil {
				return nil, fmt.Errorf("could not create bucket CORS configuration: %w", err)
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
					Resource: "bucket-cors",
					Verb:     "create",
				},
			}, nil
		},
	}
}

type bucketCorsDeleteArgs struct {
	Bucket         string
	Region         scw.Region
	ProjectID      string
	S3Endpoint     string `json:"s3-endpoint"`
	S3UsePathStyle string `json:"s3-use-path-style"`
}

func bucketCorsDeleteCommand() *core.Command {
	return &core.Command{
		Namespace: "object",
		Resource:  "bucket-cors",
		Verb:      "delete",
		Short:     "Remove CORS rules from a bucket",
		Long:      "Remove CORS rules from an Object Storage bucket with the S3 protocol.",
		ArgsType:  reflect.TypeFor[bucketCorsDeleteArgs](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:             "bucket",
				Positional:       true,
				Required:         true,
				Short:            "The name of the bucket",
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
			args := argsI.(*bucketCorsDeleteArgs)
			if args.Bucket == "" {
				return nil, errors.New("bucket name cannot be empty")
			}

			params := s3.DeleteBucketCorsInput{
				Bucket: &args.Bucket,
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

			_, err = client.DeleteBucketCors(ctx, &params)
			if err != nil {
				return nil, fmt.Errorf("could not delete bucket CORS: %w", err)
			}

			return &core.SuccessResult{
				Resource: "bucket-cors",
				Verb:     "delete",
			}, nil
		},
	}
}

type bucketCorsGetArgs struct {
	Bucket         string
	Region         scw.Region
	ProjectID      string
	S3Endpoint     string `json:"s3-endpoint"`
	S3UsePathStyle string `json:"s3-use-path-style"`
}

func bucketCorsGetCommand() *core.Command {
	return &core.Command{
		Namespace: "object",
		Resource:  "bucket-cors",
		Verb:      "get",
		Short:     "Get the CORS configuration of a bucket",
		Long:      "Get the CORS configuration of an Object Storage bucket with the S3 protocol.",
		ArgsType:  reflect.TypeFor[bucketCorsGetArgs](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:             "bucket",
				Positional:       true,
				Required:         true,
				Short:            "The name of the bucket",
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
			args := argsI.(*bucketCorsGetArgs)
			if args.Bucket == "" {
				return nil, errors.New("bucket name cannot be empty")
			}

			params := s3.GetBucketCorsInput{
				Bucket: &args.Bucket,
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

			conf, err := client.GetBucketCors(ctx, &params)
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
