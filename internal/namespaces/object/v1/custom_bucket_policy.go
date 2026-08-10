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
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

/*
**
** Bucket policy structures and marshalling methods.
**
 */

// StringList is a simple slice of string with custom unmarshalling methods
type StringList []string

func (s *StringList) UnmarshalJSON(data []byte) error {
	var single string
	if err := json.Unmarshal(data, &single); err == nil {
		*s = StringList{single}

		return nil
	}

	var multi []string
	if err := json.Unmarshal(data, &multi); err != nil {
		return err
	}

	*s = multi

	return nil
}

type BucketPolicy struct {
	Version   string                   `json:"Version"`
	ID        string                   `json:"Id,omitempty"`
	Statement []*BucketPolicyStatement `json:"Statement,omitempty"`
}

type BucketPolicyStatement struct {
	Sid       string                    `json:"Sid,omitempty"`
	Effect    string                    `json:"Effect"`
	Principal BucketPolicyPrincipal     `json:"Principal"`
	Action    StringList                `json:"Action"`
	Resource  StringList                `json:"Resource"`
	Condition map[string]map[string]any `json:"Condition,omitempty"`
}

type BucketPolicyPrincipal struct {
	SCW StringList `json:"SCW,omitempty"`
	Raw string     `json:",omitempty"` // This interpolates the "*" notation
}

func (p *BucketPolicyPrincipal) UnmarshalJSON(data []byte) error {
	var wildcard string
	if err := json.Unmarshal(data, &wildcard); err == nil {
		p.Raw = wildcard

		return nil
	}

	type Alias BucketPolicyPrincipal
	var temp Alias
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	*p = BucketPolicyPrincipal(temp)

	return nil
}

func (p *BucketPolicyPrincipal) MarshalJSON() ([]byte, error) {
	if p.Raw != "" {
		return json.Marshal(p.Raw)
	}

	type Alias BucketPolicyPrincipal

	return json.Marshal((Alias)(*p))
}

/*
**
** CLI definition for the bucket policies.
**
 */

type bucketPolicyCreateArgs struct {
	Bucket     string
	PolicyPath string
	Region     scw.Region
}

func bucketPolicyCreateCommand() *core.Command {
	return &core.Command{
		Namespace: "object",
		Resource:  "bucket-policy",
		Verb:      "create",
		Short:     "Create a policy for an S3 bucket",
		Long:      "Create a policy and apply to an Object Bucket with the S3 protocol.",
		ArgsType:  reflect.TypeOf(bucketPolicyCreateArgs{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "bucket",
				Positional: true,
				Required:   true,
				Short:      "The name of the bucket to which assign the policy.",
			},
			{
				Name:       "policy",
				Positional: true,
				Required:   true,
				Short:      "The path to the local JSON file containing the bucket policy.",
			},
			core.RegionArgSpec(),
		},
		Run: func(ctx context.Context, argsI any) (any, error) {
			args := argsI.(*bucketPolicyCreateArgs)
			if args.Bucket == "" {
				return nil, errors.New("bucket name cannot be empty")
			}
			if args.PolicyPath == "" {
				return nil, errors.New("policy file path cannot be empty")
			}

			policyBytes, err := os.ReadFile(args.PolicyPath)
			if err != nil {
				return nil, fmt.Errorf(
					"could not open policy file \"%s\": %w",
					args.PolicyPath,
					err,
				)
			}

			policy := string(policyBytes)

			client := newS3Client(ctx, args.Region)
			params := s3.PutBucketPolicyInput{
				Bucket: &args.Bucket,
				Policy: &policy,
			}

			_, err = client.PutBucketPolicy(ctx, &params)
			if err != nil {
				return nil, fmt.Errorf("could not create bucket policy: %w", err)
			}

			bucket, err := getBucketInfo(ctx, args.Region, args.Bucket)
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
	Bucket string
	Region scw.Region
}

func bucketPolicyDeleteCommand() *core.Command {
	return &core.Command{
		Namespace: "object",
		Resource:  "bucket-policy",
		Verb:      "delete",
		Short:     "Delete an S3 bucket's policy if it exists.",
		Long:      "Delete an Object Bucket's policy with the S3 protocol.",
		ArgsType:  reflect.TypeOf(bucketPolicyDeleteArgs{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:             "name",
				Positional:       true,
				Required:         true,
				Short:            "The unique name of the bucket",
				AutoCompleteFunc: autocompleteBucketName,
			},
			core.RegionArgSpec(),
		},
		Run: func(ctx context.Context, argsI any) (any, error) {
			args := argsI.(*bucketPolicyDeleteArgs)
			if args.Bucket == "" {
				return nil, errors.New("bucket name cannot be empty")
			}

			client := newS3Client(ctx, args.Region)
			params := s3.DeleteBucketPolicyInput{
				Bucket: &args.Bucket,
			}

			_, err := client.DeleteBucketPolicy(ctx, &params)
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
	Bucket string
	Region scw.Region
}

func bucketPolicyGetCommand() *core.Command {
	return &core.Command{
		Namespace: "object",
		Resource:  "bucket-policy",
		Verb:      "Get",
		Short:     "Retrieve an S3 bucket's policy.",
		Long:      "Retrieve an Object Bucket's policy with the S3 protocol.",
		ArgsType:  reflect.TypeOf(bucketPolicyDeleteArgs{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:             "name",
				Positional:       true,
				Required:         true,
				Short:            "The unique name of the bucket",
				AutoCompleteFunc: autocompleteBucketName,
			},
			core.RegionArgSpec(),
		},
		Run: func(ctx context.Context, argsI any) (any, error) {
			args := argsI.(*bucketPolicyGetArgs)
			if args.Bucket == "" {
				return nil, errors.New("bucket name cannot be empty")
			}

			client := newS3Client(ctx, args.Region)
			params := s3.GetBucketPolicyInput{
				Bucket: &args.Bucket,
			}

			policyResponse, err := client.GetBucketPolicy(ctx, &params)
			if err != nil {
				return nil, fmt.Errorf("could not get bucket policy: %w", err)
			}

			// FIXME What to return here?
			return policyResponse.Policy, nil
		},
	}
}
