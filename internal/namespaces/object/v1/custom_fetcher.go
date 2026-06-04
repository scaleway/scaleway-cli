package object

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/fetch"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type FetchBuckets struct{}

func (f FetchBuckets) Product() string {
	return "object"
}

func (f FetchBuckets) Resource() string {
	return "bucket"
}

func (f FetchBuckets) LocalityType() fetch.LocalityType {
	return fetch.LocalityTypeRegion
}

// Fetch fetches all object storage buckets in a given region.
func (f FetchBuckets) Fetch(
	ctx context.Context,
	region scw.Region,
	projectID string,
) ([]fetch.ResourceResult, error) {
	client := core.ExtractClient(ctx)

	// Create S3 client for use in the resources command.
	s3Client := newS3ClientForResources(client, region)
	if s3Client == nil {
		return nil, nil
	}

	// Note: S3 API doesn't support project filtering, so we fetch all buckets
	_ = projectID // mark as intentionally unused

	resp, err := s3Client.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		if fetch.ShouldIgnoreError(err) {
			return nil, nil
		}

		return nil, fmt.Errorf("could not list buckets: %w", err)
	}

	results := make([]fetch.ResourceResult, 0, len(resp.Buckets))
	for _, bucket := range resp.Buckets {
		bucketName := aws.ToString(bucket.Name)
		results = append(results, fetch.ResourceResult{
			Locality: region.String(),
			ID:       bucketName,
			Name:     bucketName,
		})
	}

	return results, nil
}

// newS3ClientForResources creates an S3 client for use in the resources command.
func newS3ClientForResources(scwClient *scw.Client, region scw.Region) *s3.Client {
	accessKey, ok := scwClient.GetAccessKey()
	if !ok {
		return nil
	}
	secretKey, ok := scwClient.GetSecretKey()
	if !ok {
		return nil
	}

	customEndpoint := "https://s3." + region.String() + ".scw.cloud"

	return s3.New(s3.Options{
		ClientLogMode: 0,
		Credentials: aws.CredentialsProviderFunc(func(_ context.Context) (aws.Credentials, error) {
			return aws.Credentials{
				AccessKeyID:     accessKey,
				SecretAccessKey: secretKey,
			}, nil
		}),
		BaseEndpoint: &customEndpoint,
		Region:       region.String(),
	})
}
