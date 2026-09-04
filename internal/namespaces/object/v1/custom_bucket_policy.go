//go:build darwin || linux || windows

package object

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type bucketPolicyCreateArgs struct {
	Bucket         string
	Policy         string
	Region         scw.Region
	ProjectID      string
	S3Endpoint     string
	S3UsePathStyle string
}

func bucketPolicyCreateCommand() *core.Command {
	return &core.Command{
		Namespace: "object",
		Resource:  "bucket-policy",
		Verb:      "create",
		Short:     "Create a policy for an S3 bucket",
		Long:      "Create a policy and apply to an Object Bucket with the S3 protocol.",
		ArgsType:  reflect.TypeFor[bucketPolicyCreateArgs](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "bucket",
				Positional: true,
				Required:   true,
				Short:      "The name of the bucket to which assign the policy.",
			},
			{
				Name:       "policy",
				Positional: false,
				Required:   true,
				Short:      "The path to the local JSON file containing the bucket policy.",
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
			args := argsI.(*bucketPolicyCreateArgs)
			if args.Bucket == "" {
				return nil, errors.New("bucket name cannot be empty")
			}
			if args.Policy == "" {
				return nil, errors.New("policy file path cannot be empty")
			}

			policyBytes, err := os.ReadFile(args.Policy)
			if err != nil {
				return nil, fmt.Errorf(
					"could not open policy file \"%s\": %w",
					args.Policy,
					err,
				)
			}

			policy := string(policyBytes)

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

			params := s3.PutBucketPolicyInput{
				Bucket: &args.Bucket,
				Policy: &policy,
			}

			_, err = client.PutBucketPolicy(ctx, &params)
			if err != nil {
				return nil, fmt.Errorf("could not create bucket policy: %w", err)
			}

			bucket, err := getBucketInfo(
				ctx,
				args.Region,
				args.Bucket,
				args.ProjectID,
				s3UsePathStyle,
				s3Endpoint,
			)
			if err != nil {
				return nil, fmt.Errorf("could not get bucket's information: %w", err)
			}

			return &BucketResponse{
				BucketInfo: bucket,
				SuccessResult: &core.SuccessResult{
					Resource: "bucket",
					Verb:     "create",
				},
			}, nil
		},
	}
}

type bucketPolicyDeleteArgs struct {
	Bucket         string
	Region         scw.Region
	ProjectID      string
	S3Endpoint     string
	S3UsePathStyle string
}

func bucketPolicyDeleteCommand() *core.Command {
	return &core.Command{
		Namespace: "object",
		Resource:  "bucket-policy",
		Verb:      "delete",
		Short:     "Delete an S3 bucket's policy if it exists.",
		Long:      "Delete an Object Bucket's policy with the S3 protocol.",
		ArgsType:  reflect.TypeFor[bucketPolicyDeleteArgs](),
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
			args := argsI.(*bucketPolicyDeleteArgs)
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
			params := s3.DeleteBucketPolicyInput{
				Bucket: &args.Bucket,
			}

			_, err = client.DeleteBucketPolicy(ctx, &params)
			if err != nil {
				return nil, fmt.Errorf("could not delete bucket policy: %w", err)
			}

			return &core.SuccessResult{
				Resource: "bucket-policy",
				Verb:     "delete",
			}, nil
		},
	}
}

type bucketPolicyGetArgs struct {
	Bucket         string
	Region         scw.Region
	ProjectID      string
	S3Endpoint     string
	S3UsePathStyle string
}

func bucketPolicyGetCommand() *core.Command {
	return &core.Command{
		Namespace: "object",
		Resource:  "bucket-policy",
		Verb:      "get",
		Short:     "Retrieve an S3 bucket's policy.",
		Long:      "Retrieve an Object Bucket's policy with the S3 protocol.",
		ArgsType:  reflect.TypeFor[bucketPolicyGetArgs](),
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
			args := argsI.(*bucketPolicyGetArgs)
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

			params := s3.GetBucketPolicyInput{
				Bucket: &args.Bucket,
			}

			policyResponse, err := client.GetBucketPolicy(ctx, &params)
			if err != nil {
				return nil, fmt.Errorf("could not get bucket policy: %w", err)
			}

			var prettyJSON bytes.Buffer
			err = json.Indent(&prettyJSON, []byte(*policyResponse.Policy), "", "  ")
			if err != nil {
				return nil, fmt.Errorf("error indenting JSON: %w", err)
			}

			return prettyJSON.String(), nil
		},
	}
}
