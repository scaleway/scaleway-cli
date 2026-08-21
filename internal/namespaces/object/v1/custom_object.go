//go:build darwin || linux || windows

package object

import (
	"context"
	"fmt"
	"os"
	"reflect"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type objectCreateArgs struct {
	Region               scw.Region
	ProjectID            string
	Bucket               string
	Key                  string
	File                 string
	SSECustomerAlgorithm *string
	SSECustomerKey       *string
	SSECustomerKeyMD5    *string
}

func objectCreateCommand() *core.Command {
	return &core.Command{
		Namespace: "object",
		Resource:  "object",
		Verb:      "create",
		Short:     "Create an object in an S3 bucket",
		Long:      "Create an object inside an Object Storage bucket with the S3 protocol.",
		ArgsType:  reflect.TypeOf(objectCreateArgs{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:     "bucket",
				Short:    "Name of the destination bucket",
				Required: true,
			},
			{
				Name:     "key",
				Short:    "Key of the object to create",
				Required: true,
			},
			{
				Name:     "file",
				Short:    "Path of the local file to upload",
				Required: true,
			},
			{
				Name:  "sse-customer-algorithm",
				Short: "The algorithm to use for SSE-C (e.g., AES256)",
			},
			{
				Name:  "sse-customer-key",
				Short: "The customer-provided encryption key for SSE-C",
			},
			{
				Name:  "sse-customer-key-md5",
				Short: "The MD5 hash of the customer-provided encryption key",
			},
			core.RegionArgSpec(),
			core.ProjectIDArgSpec(),
		},
		Run: func(ctx context.Context, argsI any) (any, error) {
			args := argsI.(*objectCreateArgs)

			s3Client := newS3Client(ctx, args.Region, args.ProjectID)

			file, err := os.Open(args.File)
			if err != nil {
				return nil, fmt.Errorf("failed to open file %s: %w", args.File, err)
			}
			defer file.Close()

			input := &s3.PutObjectInput{
				Bucket:               &args.Bucket,
				Key:                  &args.Key,
				Body:                 file,
				SSECustomerAlgorithm: args.SSECustomerAlgorithm,
				SSECustomerKey:       args.SSECustomerKey,
				SSECustomerKeyMD5:    args.SSECustomerKeyMD5,
			}

			_, err = s3Client.PutObject(ctx, input)
			if err != nil {
				return nil, fmt.Errorf("failed to upload object: %w", err)
			}

			return &core.SuccessResult{
				Message: fmt.Sprintf("Object '%s' successfully uploaded to bucket '%s'", args.Key, args.Bucket),
			}, nil
		},
	}
}

type objectDeleteArgs struct {
	Region    scw.Region
	ProjectID string
	Bucket    string
	Key       string
}

func objectDeleteCommand() *core.Command {
	return &core.Command{
		Namespace: "object",
		Resource:  "object",
		Verb:      "delete",
		Short:     "Delete an object in an S3 bucket",
		Long:      "Delete an object inside an Object Storage bucket with the S3 protocol.",
		ArgsType:  reflect.TypeOf(objectDeleteArgs{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:     "bucket",
				Short:    "Name of the bucket",
				Required: true,
			},
			{
				Name:     "key",
				Short:    "Key of the object to delete",
				Required: true,
			},
			core.RegionArgSpec(),
			core.ProjectIDArgSpec(),
		},
		Run: func(ctx context.Context, argsI any) (any, error) {
			args := argsI.(*objectDeleteArgs)

			s3Client := newS3Client(ctx, args.Region, args.ProjectID)

			input := &s3.DeleteObjectInput{
				Bucket: &args.Bucket,
				Key:    &args.Key,
			}

			_, err := s3Client.DeleteObject(ctx, input)
			if err != nil {
				return nil, fmt.Errorf("failed to delete object: %w", err)
			}

			return &core.SuccessResult{
				Message: fmt.Sprintf("Object '%s' successfully deleted from bucket '%s'", args.Key, args.Bucket),
			}, nil
		},
	}
}

type objectCopyArgs struct {
	Region    scw.Region
	ProjectID string

	SourceBucket               string
	SourceKey                  string
	SourceSSECustomerAlgorithm *string
	SourceSSECustomerKey       *string
	SourceSSECustomerKeyMD5    *string

	DestBucket               string
	DestKey                  string
	DestSSECustomerAlgorithm *string
	DestSSECustomerKey       *string
	DestSSECustomerKeyMD5    *string
}

func objectCopyCommand() *core.Command {
	return &core.Command{
		Namespace: "object",
		Resource:  "object",
		Verb:      "copy",
		Short:     "Copy an object in an S3 bucket",
		Long:      "Copy an object from one Object Storage bucket to another, or within the same bucket, using the S3 protocol.",
		ArgsType:  reflect.TypeOf(objectCopyArgs{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:     "source-bucket",
				Short:    "Name of the source bucket",
				Required: true,
			},
			{
				Name:     "source-key",
				Short:    "Key of the source object",
				Required: true,
			},
			{
				Name:     "dest-bucket",
				Short:    "Name of the destination bucket",
				Required: true,
			},
			{
				Name:     "dest-key",
				Short:    "Key of the destination object",
				Required: true,
			},
			{
				Name:  "sse-customer-algorithm",
				Short: "The algorithm to use for encrypting the destination object (e.g., AES256)",
			},
			{
				Name:  "sse-customer-key",
				Short: "The customer-provided encryption key for the destination object",
			},
			{
				Name:  "sse-customer-key-md5",
				Short: "The MD5 hash of the destination encryption key",
			},
			{
				Name:  "copy-source-sse-customer-algorithm",
				Short: "The algorithm used to encrypt the source object (if any)",
			},
			{
				Name:  "copy-source-sse-customer-key",
				Short: "The decryption key for the source object (if it was encrypted with SSE-C)",
			},
			{
				Name:  "copy-source-sse-customer-key-md5",
				Short: "The MD5 hash of the source decryption key",
			},
			core.RegionArgSpec(),
			core.ProjectIDArgSpec(),
		},
		Run: func(ctx context.Context, argsI any) (any, error) {
			args := argsI.(*objectCopyArgs)

			s3Client := newS3Client(ctx, args.Region, args.ProjectID)

			copySource := fmt.Sprintf("%s/%s", args.SourceBucket, args.SourceKey)

			input := &s3.CopyObjectInput{
				Bucket:     &args.DestBucket,
				Key:        &args.DestKey,
				CopySource: &copySource,

				SSECustomerAlgorithm: args.DestSSECustomerAlgorithm,
				SSECustomerKey:       args.DestSSECustomerKey,
				SSECustomerKeyMD5:    args.DestSSECustomerKeyMD5,

				CopySourceSSECustomerAlgorithm: args.SourceSSECustomerAlgorithm,
				CopySourceSSECustomerKey:       args.SourceSSECustomerKey,
				CopySourceSSECustomerKeyMD5:    args.SourceSSECustomerKeyMD5,
			}

			_, err := s3Client.CopyObject(ctx, input)
			if err != nil {
				return nil, fmt.Errorf("failed to copy object: %w", err)
			}

			return &core.SuccessResult{
				Message: fmt.Sprintf("Object successfully copied to s3://%s/%s", args.DestBucket, args.DestKey),
			}, nil
		},
	}
}

type objectListArgs struct {
	Region    scw.Region
	ProjectID string
	Bucket    string
	Prefix    string
}

func objectListCommand() *core.Command {
	return &core.Command{
		Namespace: "object",
		Resource:  "object",
		Verb:      "list",
		Short:     "List objects in an S3 bucket",
		Long:      "List all objects inside an Object Storage bucket using the S3 protocol.",
		ArgsType:  reflect.TypeOf(objectListArgs{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:     "bucket",
				Short:    "Name of the bucket to list",
				Required: true,
			},
			{
				Name:  "prefix",
				Short: "Optional prefix to filter objects",
			},
			core.RegionArgSpec(),
			core.ProjectIDArgSpec(),
		},
		Run: func(ctx context.Context, argsI any) (any, error) {
			args := argsI.(*objectListArgs)

			s3Client := newS3Client(ctx, args.Region, args.ProjectID)

			input := &s3.ListObjectsV2Input{
				Bucket: &args.Bucket,
			}
			if args.Prefix != "" {
				input.Prefix = &args.Prefix
			}

			resp, err := s3Client.ListObjectsV2(ctx, input)
			if err != nil {
				return nil, fmt.Errorf("failed to list objects: %w", err)
			}

			return resp.Contents, nil
		},
	}
}
