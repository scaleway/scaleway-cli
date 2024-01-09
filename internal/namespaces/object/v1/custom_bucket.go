package object

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"reflect"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/scaleway/scaleway-cli/v2/internal/core"
	"github.com/scaleway/scaleway-cli/v2/internal/human"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type objectBucketBasicArgs struct {
	Region scw.Region
	Name   string
}

type objectBucketInfo struct {
	ID               string
	Region           scw.Region
	APIEndpoint      string
	BucketEndpoint   string
	EnableVersioning bool
	Tags             []types.Tag
	ACL              []customS3ACLGrant
	Owner            *string
}

func objectBucketInfoMarshalerFunc(i interface{}, opt *human.MarshalOpt) (string, error) {
	// To avoid recursion of human.Marshal we create a dummy type
	type tmp objectBucketInfo
	info := tmp(i.(objectBucketInfo))

	opt.Sections = []*human.MarshalSection{
		{
			FieldName:   "Tags",
			HideIfEmpty: true,
		},
		{
			FieldName:   "ACL",
			HideIfEmpty: true,
		},
	}
	str, err := human.Marshal(info, opt)
	if err != nil {
		return "", err
	}
	return str, nil
}

type objectBucketResponse struct {
	SuccessResult *core.SuccessResult
	BucketInfo    objectBucketInfo
}

func objectBucketResponseMarshalerFunc(i interface{}, opt *human.MarshalOpt) (string, error) {
	resp := i.(objectBucketResponse)

	messageStr, err := resp.SuccessResult.MarshalHuman()
	if err != nil {
		return "", err
	}
	bucketStr, err := objectBucketInfoMarshalerFunc(resp.BucketInfo, opt)
	if err != nil {
		return "", err
	}

	return strings.Join([]string{
		messageStr,
		bucketStr,
	}, "\n"), nil
}

type objectBucketConfigArgs struct {
	Region           scw.Region
	Name             string
	Tags             []string
	EnableVersioning bool `json:"enable-versioning"`
	ACL              string
}

func bucketCreateCommand() *core.Command {
	return &core.Command{
		Namespace: "object",
		Resource:  "bucket",
		Verb:      "create",
		Short:     "Create an S3 bucket",
		Long:      "Create an Object Storage Bucket with the S3 protocol. The namespace is shared between all S3 users, so its name must be unique.",
		ArgsType:  reflect.TypeOf(objectBucketConfigArgs{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "name",
				Positional: true,
				Required:   true,
				Short:      "The unique name of the bucket",
			},
			{
				Name:       "tags",
				Positional: false,
				Required:   false,
				Short:      "The new tags to set on the bucket",
			},
			{
				Name:       "enable-versioning",
				Positional: false,
				Required:   false,
				Default:    core.DefaultValueSetter("false"),
				Short:      "Whether or not objects in the bucket should have multiple versions",
			},
			{
				Name:             "acl",
				Positional:       false,
				Required:         false,
				Default:          core.DefaultValueSetter("private"),
				Short:            "The permissions given to users (grantees) to read or write objects",
				AutoCompleteFunc: autocompleteObjectBucketACL,
			},
			core.RegionArgSpec(),
		},
		Run: func(ctx context.Context, argsI interface{}) (interface{}, error) {
			args := argsI.(*objectBucketConfigArgs)
			client := newS3Client(ctx, args.Region)

			if ok, possibleValues := verifyACLInput(args.ACL); !ok {
				return nil, fmt.Errorf("ACL field must be one of %v", possibleValues)
			}

			_, err := client.CreateBucket(ctx, &s3.CreateBucketInput{
				Bucket: &args.Name,
				ACL:    types.BucketCannedACL(args.ACL),
				CreateBucketConfiguration: &types.CreateBucketConfiguration{
					LocationConstraint: types.BucketLocationConstraint(args.Region),
				},
				ObjectLockEnabledForBucket: scw.BoolPtr(false),
			})
			if err != nil {
				return nil, fmt.Errorf("could not create bucket: %w", err)
			}

			err = putBucketVersioning(ctx, client, args.Name, args.EnableVersioning)
			if err != nil {
				return nil, fmt.Errorf("could not put bucket versioning: %w", err)
			}
			err = putBucketTagging(ctx, client, args.Name, args.Tags)
			if err != nil {
				return nil, fmt.Errorf("could not put bucket tags: %w", err)
			}

			bucketResponse, err := getBucketInfo(ctx, args.Region, args.Name)
			if err != nil {
				return nil, fmt.Errorf("could not get bucket's information: %w", err)
			}

			return &objectBucketResponse{
				BucketInfo: *bucketResponse,
				SuccessResult: &core.SuccessResult{
					Resource: "bucket",
					Verb:     "create",
				},
			}, nil
		},
	}
}

func bucketDeleteCommand() *core.Command {
	return &core.Command{
		Namespace: "object",
		Resource:  "bucket",
		Verb:      "delete",
		Short:     "Delete an S3 bucket",
		Long:      "Delete an S3 bucket with all its content.",
		ArgsType:  reflect.TypeOf(objectBucketBasicArgs{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "name",
				Positional: true,
				Required:   true,
				Short:      "The unique name of the bucket",
			},
			core.RegionArgSpec(),
		},
		Run: func(ctx context.Context, argsI interface{}) (interface{}, error) {
			args := argsI.(*objectBucketBasicArgs)
			client := newS3Client(ctx, args.Region)

			_, err := client.DeleteBucket(ctx, &s3.DeleteBucketInput{
				Bucket: &args.Name,
			})
			if err != nil {
				return nil, err
			}

			return &core.SuccessResult{
				Resource: "bucket",
				Verb:     "delete",
			}, nil
		},
	}
}

type objectBucketGetArgs struct {
	Region   scw.Region
	Name     string
	WithSize bool
}

type objectBucketGetResult struct {
	*objectBucketInfo
	Size      *scw.Size
	NbObjects *int64
	NbParts   *int64
}

func objectBucketGetResultMarshalerFunc(i interface{}, opt *human.MarshalOpt) (string, error) {
	type tmp objectBucketGetResult
	result := tmp(i.(objectBucketGetResult))
	opt.Sections = []*human.MarshalSection{
		{
			FieldName:   "Tags",
			HideIfEmpty: true,
		},
		{
			FieldName:   "ACL",
			HideIfEmpty: true,
		},
	}
	str, err := human.Marshal(result, opt)
	if err != nil {
		return "", err
	}
	return str, nil
}

func bucketGetCommand() *core.Command {
	return &core.Command{
		Namespace: "object",
		Resource:  "bucket",
		Verb:      "get",
		Short:     "Get information about an S3 bucket",
		Long:      "Get the properties of an S3 bucket like tags, endpoint, access control, versioning, size, etc.",
		ArgsType:  reflect.TypeOf(objectBucketGetArgs{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "name",
				Positional: true,
				Required:   true,
				Short:      "The unique name of the bucket",
			},
			{
				Name:       "with-size",
				Positional: false,
				Required:   false,
				Default:    core.DefaultValueSetter("false"),
				Short:      "Whether to calculate the total size of the bucket and the number of objects. This operation can take long for large buckets.",
			},
			core.RegionArgSpec(),
		},
		Run: func(ctx context.Context, argsI interface{}) (interface{}, error) {
			args := argsI.(*objectBucketGetArgs)
			client := newS3Client(ctx, args.Region)

			bucketInfo, err := getBucketInfo(ctx, args.Region, args.Name)
			if err != nil {
				return nil, fmt.Errorf("could not get bucket's information: %w", err)
			}

			result := objectBucketGetResult{
				bucketInfo,
				nil,
				nil,
				nil,
			}

			if args.WithSize {
				nbObjects, nbParts, totalSize, err := countBucketObjects(ctx, client, args.Name)
				if err != nil {
					return nil, fmt.Errorf("could not get bucket size: %w", err)
				}
				result.Size = &totalSize
				result.NbObjects = &nbObjects
				result.NbParts = &nbParts
			}

			return result, nil
		},
	}
}

func bucketUpdateCommand() *core.Command {
	return &core.Command{
		Namespace: "object",
		Resource:  "bucket",
		Verb:      "update",
		Short:     "Update an S3 bucket",
		Long:      "Update an S3 bucket's properties like tags, access control and versioning.",
		ArgsType:  reflect.TypeOf(objectBucketConfigArgs{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "name",
				Positional: true,
				Required:   true,
				Short:      "The unique name of the bucket",
			},
			{
				Name:       "tags",
				Positional: false,
				Required:   false,
				Short:      "The new tags to set on the bucket",
			},
			{
				Name:       "enable-versioning",
				Positional: false,
				Required:   false,
				Default:    core.DefaultValueSetter("false"),
				Short:      "Whether or not objects in the bucket should have multiple versions",
			},
			{
				Name:             "acl",
				Positional:       false,
				Required:         false,
				Default:          core.DefaultValueSetter("private"),
				Short:            "The permissions given to users (grantees) to read or write objects",
				AutoCompleteFunc: autocompleteObjectBucketACL,
			},
			core.RegionArgSpec(),
		},
		Run: func(ctx context.Context, argsI interface{}) (interface{}, error) {
			args := argsI.(*objectBucketConfigArgs)
			client := newS3Client(ctx, args.Region)

			err := putBucketVersioning(ctx, client, args.Name, args.EnableVersioning)
			if err != nil {
				return nil, fmt.Errorf("could not update bucket versioning: %w", err)
			}

			err = putBucketTagging(ctx, client, args.Name, args.Tags)
			if err != nil {
				return nil, fmt.Errorf("could not update bucket tags: %w", err)
			}

			err = putBucketACL(ctx, client, args.Name, args.ACL)
			if err != nil {
				return nil, fmt.Errorf("could not update bucket ACL: %w", err)
			}

			bucketResponse, err := getBucketInfo(ctx, args.Region, args.Name)
			if err != nil {
				return nil, fmt.Errorf("could not get bucket's information: %w", err)
			}

			return &objectBucketResponse{
				BucketInfo: *bucketResponse,
				SuccessResult: &core.SuccessResult{
					Resource: "bucket",
					Verb:     "update",
				},
			}, nil
		},
	}
}

func getBucketInfo(ctx context.Context, region scw.Region, name string) (*objectBucketInfo, error) {
	client := newS3Client(ctx, region)
	bucketInfo := &objectBucketInfo{
		ID:     name,
		Region: region,
	}

	// get versioning
	versioningOutput, err := client.GetBucketVersioning(ctx, &s3.GetBucketVersioningInput{
		Bucket: &name,
	})
	if err != nil {
		return nil, fmt.Errorf("could not get bucket versioning: %w", err)
	}
	switch versioningOutput.Status {
	case types.BucketVersioningStatusSuspended, "":
		bucketInfo.EnableVersioning = false
	case types.BucketVersioningStatusEnabled:
		bucketInfo.EnableVersioning = true
	}

	// get tagging
	tagging, err := client.GetBucketTagging(ctx, &s3.GetBucketTaggingInput{
		Bucket: &name,
	})
	if err != nil && !strings.Contains(err.Error(), "NoSuchTagSet") {
		return nil, fmt.Errorf("could not get bucket tagging: %w", err)
	} else if tagging != nil {
		bucketInfo.Tags = tagging.TagSet
	}

	// get ACL
	acl, err := client.GetBucketAcl(ctx, &s3.GetBucketAclInput{
		Bucket: &name,
	})
	if err != nil {
		return nil, fmt.Errorf("could not get bucket ACL: %w", err)
	}
	bucketInfo.Owner = normalizeOwnerID(acl.Owner.ID)
	bucketInfo.ACL = awsACLToCustomGrants(acl)

	// get endpoints
	bucketInfo.APIEndpoint = getAPIEndpoint(region.String())
	bucketInfo.BucketEndpoint, err = getBucketEndpoint(name, region.String())
	if err != nil {
		return nil, err
	}

	return bucketInfo, nil
}

func getAPIEndpoint(region string) string {
	if customEndpoint := os.Getenv("SCW_S3_ENDPOINT"); customEndpoint != "" {
		return customEndpoint
	}
	return fmt.Sprintf("https://s3.%s.scw.cloud", region)
}

func getBucketEndpoint(name, region string) (string, error) {
	if customEndpoint := os.Getenv("SCW_S3_ENDPOINT"); customEndpoint != "" {
		u, err := url.Parse(customEndpoint)
		if err != nil {
			return "", fmt.Errorf("could not parse custom endpoint %s: %w", customEndpoint, err)
		}
		u = u.JoinPath(name, u.Path)
		return u.String(), nil
	}
	return fmt.Sprintf("https://%s.s3.%s.scw.cloud", name, region), nil
}

func countBucketObjects(ctx context.Context, client *s3.Client, name string) (nbObjects, nbParts int64, totalSize scw.Size, err error) {
	var size int64

	objectsList, err := client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: &name,
	})
	if err != nil {
		return 0, 0, 0, err
	}

	for _, object := range objectsList.Contents {
		nbObjects++
		size += *object.Size
	}

	multipartsList, err := client.ListMultipartUploads(ctx, &s3.ListMultipartUploadsInput{
		Bucket:         &name,
		KeyMarker:      nil,
		MaxUploads:     nil,
		Prefix:         nil,
		RequestPayer:   "",
		UploadIdMarker: nil,
	})
	if err != nil {
		return 0, 0, 0, err
	}

	for _, multipart := range multipartsList.Uploads {
		partsList, err := client.ListParts(ctx, &s3.ListPartsInput{
			Bucket:   &name,
			Key:      multipart.Key,
			UploadId: multipart.UploadId,
		})
		if err != nil {
			return 0, 0, 0, err
		}
		for _, part := range partsList.Parts {
			nbParts++
			size += *part.Size
		}
	}

	return nbObjects, nbParts, scw.Size(size), nil
}

func putBucketVersioning(ctx context.Context, client *s3.Client, name string, enabled bool) error {
	request := &s3.PutBucketVersioningInput{
		Bucket: &name,
	}
	if enabled {
		request.VersioningConfiguration = &types.VersioningConfiguration{
			Status: types.BucketVersioningStatusEnabled,
		}
	} else {
		request.VersioningConfiguration = &types.VersioningConfiguration{
			Status: types.BucketVersioningStatusSuspended,
		}
	}
	_, err := client.PutBucketVersioning(ctx, request)
	if err != nil {
		return err
	}
	return nil
}

func putBucketTagging(ctx context.Context, client *s3.Client, name string, tags []string) error {
	if len(tags) == 0 {
		return nil
	}
	newTags := []types.Tag(nil)
	for _, tag := range tags {
		trim := strings.Trim(tag, "\"")
		split := strings.Split(trim, "=")
		newTags = append(newTags, types.Tag{
			Key:   &split[0],
			Value: &split[1],
		})
	}
	_, err := client.PutBucketTagging(ctx, &s3.PutBucketTaggingInput{
		Bucket:  &name,
		Tagging: &types.Tagging{TagSet: newTags},
	})
	if err != nil {
		return err
	}
	return nil
}

func putBucketACL(ctx context.Context, client *s3.Client, name string, acl string) error {
	_, err := client.PutBucketAcl(ctx, &s3.PutBucketAclInput{
		Bucket: &name,
		ACL:    types.BucketCannedACL(acl),
	})
	if err != nil {
		return err
	}
	return nil
}
