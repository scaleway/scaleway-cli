//go:build darwin || linux || windows

package object

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go/middleware"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/scaleway/scaleway-sdk-go/validation"
)

// newS3Client creates a new S3 client to interact with the S3 API of the passed
// region. If `projectID` is empty, the default one is used.
func newS3Client(
	ctx context.Context,
	region scw.Region,
	projectID string,
	s3Endpoint string,
	s3UsePathStyle bool,
) *s3.Client {
	httpClient := core.ExtractHTTPClient(ctx)
	scwClient := core.ExtractClient(ctx)
	buildInfo := core.ExtractBuildInfo(ctx)
	accessKey, ok := scwClient.GetAccessKey()
	if !ok {
		return nil
	}
	secretKey, ok := scwClient.GetSecretKey()
	if !ok {
		return nil
	}

	defaultProjectID, _ := scwClient.GetDefaultProjectID()
	accessKey = FormatAccessKey(accessKey, projectID, defaultProjectID)

	options := []func(*middleware.Stack) error{
		func(stack *middleware.Stack) error {
			return awsmiddleware.AddUserAgentKeyValue(
				"scaleway-cli",
				buildInfo.Version.String(),
			)(stack)
		},
	}

	return s3.New(s3.Options{
		APIOptions:    options,
		ClientLogMode: 0,
		Credentials: aws.CredentialsProviderFunc(func(_ context.Context) (aws.Credentials, error) {
			return aws.Credentials{
				AccessKeyID:     accessKey,
				SecretAccessKey: secretKey,
			}, nil
		}),
		BaseEndpoint: new(s3Endpoint),
		Region:       region.String(),
		HTTPClient:   httpClient,
		UsePathStyle: s3UsePathStyle,
	})
}

// FormatAccessKey formats the access key to the <KEY>@<PROJECT_ID> format,
// overriding the already present project ID with the "project-id" argument
// if present.
func FormatAccessKey(accessKey, argProjectID, defaultProjectID string) (formattedAccessKey string) {
	// The project ID from the CLI arguments takes precedence
	projectID := argProjectID

	if projectID == "" {
		projectID = defaultProjectID
	}

	if projectID != "" {
		if validation.IsAccessKeyWithProjectID(accessKey) {
			keySplit := strings.Split(accessKey, "@")
			formattedAccessKey = keySplit[0] + "@" + projectID
		} else {
			// Is a standard access key
			formattedAccessKey = accessKey + "@" + projectID
		}
	}

	return
}

// Caching BucketCannedACL values for shell completion
var completeBucketACLCache []types.BucketCannedACL

func autocompleteBucketACL(_ context.Context, prefix string, _ any) core.AutocompleteSuggestions {
	suggestions := core.AutocompleteSuggestions(nil)

	if len(completeBucketACLCache) == 0 {
		var awsCannedACL types.BucketCannedACL
		completeBucketACLCache = awsCannedACL.Values()
	}

	for _, acl := range completeBucketACLCache {
		if strings.HasPrefix(string(acl), prefix) {
			suggestions = append(suggestions, string(acl))
		}
	}

	return suggestions
}

func verifyACLInput(aclInput string) (bool, []types.BucketCannedACL) {
	var awsCannedACL types.BucketCannedACL
	possibleValues := awsCannedACL.Values()

	for _, possibleValue := range possibleValues {
		if string(possibleValue) == aclInput {
			return true, nil
		}
	}

	return false, possibleValues
}

func awsACLToCustomGrants(output *s3.GetBucketAclOutput) []CustomS3ACLGrant {
	customGrants := make([]CustomS3ACLGrant, 0, len(output.Grants))
	for _, grant := range output.Grants {
		var grantee *string
		switch grant.Grantee.Type {
		case types.TypeCanonicalUser:
			grantee = new(normalizeOwnerID(grant.Grantee.ID))
		case types.TypeGroup:
			split := strings.Split(*grant.Grantee.URI, "/")
			grantee = new(split[len(split)-1])
		}
		customGrants = append(customGrants, CustomS3ACLGrant{
			Grantee:    grantee,
			Permission: grant.Permission,
		})
	}

	return customGrants
}

func normalizeOwnerID(id *string) string {
	if id == nil {
		return ""
	}
	tab := strings.Split(*id, ":")
	if len(tab) != 2 {
		return ""
	}

	return tab[0]
}

func getBucketInfo(
	ctx context.Context,
	region scw.Region,
	name string,
	projectID string,
	s3UsePathStyle bool,
	s3Endpoint string,
) (*bucketInfo, error) {
	client := newS3Client(ctx, region, projectID, s3Endpoint, s3UsePathStyle)

	bucket := &bucketInfo{
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
		bucket.EnableVersioning = false
	case types.BucketVersioningStatusEnabled:
		bucket.EnableVersioning = true
	}

	// get tagging
	tagging, err := client.GetBucketTagging(ctx, &s3.GetBucketTaggingInput{
		Bucket: &name,
	})
	if err != nil && !strings.Contains(err.Error(), "NoSuchTagSet") {
		return nil, fmt.Errorf("could not get bucket tagging: %w", err)
	} else if tagging != nil {
		bucket.Tags = tagging.TagSet
	}

	// get ACL
	acl, err := client.GetBucketAcl(ctx, &s3.GetBucketAclInput{
		Bucket: &name,
	})
	if err != nil {
		return nil, fmt.Errorf("could not get bucket ACL: %w", err)
	}
	bucket.Owner = normalizeOwnerID(acl.Owner.ID)
	bucket.ACL = awsACLToCustomGrants(acl)

	// Get endpoints
	bucket.APIEndpoint, bucket.BucketEndpoint, err = getS3Endpoints(
		ctx, region.String(), s3Endpoint, name, s3UsePathStyle,
	)
	if err != nil {
		return nil, err
	}

	return bucket, nil
}

// getS3Endpoints retrieves the S3 API URL according to various configuration
// variables. The priority is following this pattern:
// - CLI arg (parameter of this function)
// - SCW Environment variable
// - AWS configuration
// - Profile field value
// - Default value, built with the provided region
// Also builds the bucket endpoint when provided with a bucket name.
func getS3Endpoints(
	ctx context.Context, region, customEndpoint, bucketName string, pathStyle bool,
) (s3Endpoint, bucketEndpoint string, err error) {
	// 1. Get S3 endpoint

	if customEndpoint != "" {
		// CLI argument
		s3Endpoint = customEndpoint
	} else if ep := os.Getenv(scw.ScwS3EndpointEnv); ep != "" {
		// SCW environment variable
		s3Endpoint = ep
	} else if ep := scw.GetS3EndpointFromAWSConf(); ep != "" {
		// AWS configuration, by environment variable or config file
		s3Endpoint = ep
	} else {
		// Profile field value
		scwClient := core.ExtractClient(ctx)
		profileS3Endpoint, s3EndpointOk := scwClient.GetS3Endpoint()

		if s3EndpointOk && profileS3Endpoint != "" {
			// Don't use the default S3 endpoint with the default region,
			// in case a different region was passed to this fonction
			defaultRegion, _ := scwClient.GetDefaultRegion()
			autoGeneratedEndpoint := fmt.Sprintf("https://s3.%s.scw.cloud", defaultRegion)
			if profileS3Endpoint != autoGeneratedEndpoint {
				s3Endpoint = profileS3Endpoint
			}
		}
	}

	// Default value
	if s3Endpoint == "" {
		s3Endpoint = fmt.Sprintf("https://s3.%s.scw.cloud", region)
	}

	// 2. Build bucket endpoint

	if bucketName != "" {
		u, err := url.Parse(s3Endpoint)
		if err != nil {
			return "", "", fmt.Errorf("could not parse custom endpoint %s: %w", s3Endpoint, err)
		}

		if pathStyle {
			u = u.JoinPath(bucketName)
		} else {
			u.Host = bucketName + "." + u.Host
		}

		bucketEndpoint = u.String()
	}

	return s3Endpoint, bucketEndpoint, nil
}

// getS3UsePathStyle retrieves the S3UsePathStyle flag value from various
// configuration variables. The priority is following this pattern:
// - CLI arg (parameter of this function)
// - SCW Environment variable
// - Profile field value
func getS3UsePathStyle(ctx context.Context, usePathStyle string) (bool, error) {
	// CLI argument
	if usePathStyle != "" {
		parsed, err := strconv.ParseBool(usePathStyle)

		return parsed, err
	}

	// SCW environment variable
	if flag := os.Getenv(scw.ScwS3UsePathStyleEnv); flag != "" {
		parsed, err := strconv.ParseBool(flag)

		return parsed, err
	}

	// Profile field value
	scwClient := core.ExtractClient(ctx)

	return scwClient.GetS3UsePathStyle(), nil
}
