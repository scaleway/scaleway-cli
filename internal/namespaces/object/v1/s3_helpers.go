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
	bucket.APIEndpoint = getS3Endpoint(ctx, region.String(), s3Endpoint)
	bucket.BucketEndpoint, err = getBucketEndpoint(
		name,
		region.String(),
		s3Endpoint,
		s3UsePathStyle,
	)
	if err != nil {
		return nil, err
	}

	return bucket, nil
}

// getS3Endpoint retrieves the S3 API URL according to various configuration
// variables. The priority is following this pattern:
// - CLI arg (parameter of this function)
// - SCW Environment variable
// - AWS configuration
// - Profile field value
// - Default value, built with the provided region
func getS3Endpoint(ctx context.Context, region string, customEndpoint string) string {
	// CLI argument
	if customEndpoint != "" {
		return customEndpoint
	}

	// SCW environment variable
	if ep := os.Getenv(scw.ScwS3EndpointEnv); ep != "" {
		return ep
	}

	// AWS configuration, by environment variable or config file
	ep := scw.GetS3EndpointFromAWSConf()
	if ep != "" {
		return ep
	}

	// Profile field value
	// The SDK auto-populates the S3 endpoint from the client's default region,
	// so we must check whether the profile endpoint is just that auto-generated
	// default. If it is, we fall through to use the passed region instead.
	scwClient := core.ExtractClient(ctx)
	profileS3Endpoint, s3EndpointOk := scwClient.GetS3Endpoint()

	if s3EndpointOk && profileS3Endpoint != "" {
		defaultRegion, _ := scwClient.GetDefaultRegion()
		autoGeneratedEndpoint := fmt.Sprintf("https://s3.%s.scw.cloud", defaultRegion)
		if profileS3Endpoint != autoGeneratedEndpoint {
			return profileS3Endpoint
		}
	}

	// Default value
	return fmt.Sprintf("https://s3.%s.scw.cloud", region)
}

// getS3UsePathStyle retrieves the S3UsePathStyle flag value from various
// configuration variables. The priority is following this pattern:
// - CLI arg (parameter of this function)
// - SCW Environment variable
// - Profile field value
func getS3UsePathStyle(ctx context.Context, usePathStyle string) bool {
	// CLI argument
	if usePathStyle != "" {
		if usePathStyle == "true" {
			return true
		} else {
			return false
		}
	}

	// SCW environment variable
	if flag := os.Getenv(scw.ScwS3UsePathStyleEnv); flag != "" {
		if parsed, err := strconv.ParseBool(flag); err == nil {
			return parsed
		}
	}

	// Profile field value
	scwClient := core.ExtractClient(ctx)

	return scwClient.GetS3UsePathStyle()
}

func getBucketEndpoint(
	name, region string,
	s3Endpoint string,
	s3UsePathStyle bool,
) (string, error) {
	if s3Endpoint != "" {
		u, err := url.Parse(s3Endpoint)
		if err != nil {
			return "", fmt.Errorf("could not parse custom endpoint %s: %w", s3Endpoint, err)
		}

		if s3UsePathStyle {
			u = u.JoinPath(name)
		} else {
			u.Host = name + "." + u.Host
		}

		return u.String(), nil
	}

	if s3UsePathStyle {
		return fmt.Sprintf("https://s3.%s.scw.cloud/%s", region, name), nil
	}

	return fmt.Sprintf("https://%s.s3.%s.scw.cloud", name, region), nil
}
